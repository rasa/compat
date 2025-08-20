//nolint:mnd // quiet linter
// SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
// SPDX-License-Identifier: MIT

//go:build !wasip1

// package updater checks if there are any updates to snippetted code, and
// provides diffs of any changes found.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/mattn/go-colorable"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

type Snip struct {
	// Raw GitHub URL found in the source file comment.
	URL string

	Owner   string
	Repo    string
	BaseRef string // commit or ref from the URL (for base)
	ThemRef string // commit or ref from the URL (for them)
	Path    string // repo-relative path (e.g., src/os/tempfile.go)

	StartLine int // from #Lstart
	EndLine   int // inclusive

	// File-safe joined name of Path (after /src/ if present), e.g. os_tempfile.go
	Joined string

	// Ours text (section between this Snippet: and next Snippet: or EOF)
	Ours []byte
}

var (
	baseDir string
	themDir string
	snipDir string
	meldDir string
)

func main() {
	dirs := []string{"../../golang"}
	if len(os.Args) > 1 {
		dirs = os.Args[:1]
	}

	rootDir := xdg.CacheHome + "/goupdater"
	baseDir = rootDir + "/base"
	themDir = rootDir + "/them"
	snipDir = rootDir + "/snip"
	meldDir = rootDir + "/meld"

	must(os.MkdirAll(baseDir, 0o755))
	must(os.MkdirAll(themDir, 0o755))
	must(os.MkdirAll(snipDir, 0o755))
	must(os.MkdirAll(meldDir, 0o755))

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		must(err)
		for _, entry := range entries {
			name := filepath.Join(dir, entry.Name())
			if !strings.HasSuffix(name, ".go") {
				continue
			}
			doFile(name)
		}
	}
}

func doFile(src string) {
	log.Printf("Reading %v", src)
	b, err := os.ReadFile(src)
	must(err)

	snips, err := findSnips(b)
	must(err)
	if len(snips) == 0 {
		log.Printf("No snips found in %v", src)
		return
	}
	log.Printf("Found %d snips in %v", len(snips), src)

	must(DownloadBase(snips))

	snips, err = DownloadThem(snips)
	must(err)

	must(SnipOurs(snips))

	must(SnipBase(snips))

	must(SnipThem(snips))

	must(Meld(snips))

	must(Compare(snips))

	fmt.Println("Done.")
}

// DownloadBase downloads each "base" file at the exact ref from the Source URL
// into cache/base/<joined>.
func DownloadBase(snips []Snip) error {
	for i, s := range snips {
		snip_id := i + 1
		dst := cacheName(baseDir, s)
		if fileExists(dst) {
			continue
		}
		log.Printf("[base] Downloading snip %2d/%2d: %v", snip_id, len(snips), dst)

		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", s.Owner, s.Repo, s.BaseRef, s.Path)
		body, err := httpGet(rawURL)
		if err != nil {
			return fmt.Errorf("DownloadBase %s: %w", rawURL, err)
		}
		if err := writeFileAtomic(dst, body, 0o644); err != nil {
			return err
		}
	}

	return nil
}

// DownloadThem downloads the same path at repository HEAD into cache/them/<joined>.
func DownloadThem(snips []Snip) ([]Snip, error) {
	for i, s := range snips {
		snip_id := i + 1
		defBranch, err := getDefaultBranch(s.Owner, s.Repo)
		must(err)
		themRef, err := getCommitSHA(s.Owner, s.Repo, defBranch)
		must(err)
		snips[i].ThemRef = themRef
		s.ThemRef = themRef
		dst := cacheNameThem(themDir, s)
		if fileExists(dst) {
			continue
		}
		log.Printf("[base] Downloading snip %2d/%2d: %v", snip_id, len(snips), dst)
		// Use HEAD to track the latest on default branch.
		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/HEAD/%s", s.Owner, s.Repo, s.Path)
		body, err := httpGet(rawURL)
		if err != nil {
			return snips, fmt.Errorf("DownloadThem %s: %w", rawURL, err)
		}
		err = writeFileAtomic(dst, body, 0o644)
		if err != nil {
			return snips, err
		}
	}

	return snips, nil
}

// SnipOurs writes the section between Source lines as *.ours.
func SnipOurs(snips []Snip) error {
	for i, s := range snips {
		snip_id := i + 1
		p := snipName(snipDir, "ours", s)
		log.Printf("[ours] Creating snip %2d/%2d: %v (%d bytes)", snip_id, len(snips), p, len(s.Ours))
		if fileExists(p) {
			continue
		}

		err := writeFileAtomic(p, s.Ours, 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}

// SnipBase uses the start/end lines from the URL to slice the base file into *.base.
func SnipBase(snips []Snip) error {
	for i, s := range snips {
		snip_id := i + 1
		basePath := cacheName(baseDir, s)
		content, err := os.ReadFile(basePath)
		if err != nil {
			return fmt.Errorf("SnipBase read base %s: %w", basePath, err)
		}
		part, err := extractLines(content, s.StartLine, s.EndLine)
		if err != nil {
			return fmt.Errorf("SnipBase extract %s L%d-L%d: %w", s.Joined, s.StartLine, s.EndLine, err)
		}
		out := snipName(snipDir, "base", s)
		log.Printf("[base] Creating snip %2d/%2d: %v (%d bytes)", snip_id, len(snips), out, len(part))
		if fileExists(out) {
			continue
		}
		err = writeFileAtomic(out, part, 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	funcRE  = regexp.MustCompile(`^\s*func\s+([^(]+)\(`)
	braceRE = regexp.MustCompile(`^}`)
)

// SnipThem finds in "them" the region most similar to the base snippet and writes *.them.
// If an exact substring match exists, it is used; otherwise, choose the best window by
// minimal normalized edit distance.
func SnipThem(snips []Snip) error { //nolint:gocyclo // quiet linter
	for i, s := range snips {
		snip_id := i + 1
		log.Printf("[them] Creating snip %2d/%2d: %v@%d", snip_id, len(snips), s.Path, s.StartLine)
		themPath := cacheNameThem(themDir, s)
		themContent, err := os.ReadFile(themPath)
		if err != nil {
			return fmt.Errorf("SnipThem read them %s: %w", themPath, err)
		}
		basePath := snipName(snipDir, "base", s)
		baseContent, err := os.ReadFile(basePath)
		if err != nil {
			return fmt.Errorf("SnipThem read base snip %s: %w", basePath, err)
		}

		// Try exact match first.
		if idx := bytes.Index(themContent, baseContent); idx >= 0 {
			out := snipName(snipDir, "them", s)
			if !fileExists(out) {
				err = writeFileAtomic(out, baseContent, 0o644)
				if err != nil {
					return err
				}
				// log.Printf("[them] Created snip %2d: %v (%d bytes)", snip_id, out, len(baseContent))
			}
			continue
		}

		// Approximate: sliding window over "them" by line count.
		baseLines := splitLines(baseContent)
		themLines := splitLines(themContent)

		funcName := ""
		for _, line := range baseLines {
			matches := funcRE.FindStringSubmatch(line)
			if len(matches) > 1 {
				funcName = matches[1]
				break
			}
		}
		if funcName == "" {
			log.Printf("[them] snip %2d/%2d: Cannot find function name in %v (%d lines)", snip_id, len(snips), basePath, len(baseLines))
			os.Exit(1)
		}

		var best []byte
		lines := []string{}

		for _, line := range themLines {
			// log.Printf("%4d: %v", i, line)
			if len(lines) == 0 {
				matches := funcRE.FindStringSubmatch(line)
				if len(matches) > 0 && matches[1] == funcName {
					lines = append(lines, line)
				}
				continue
			}
			lines = append(lines, line)
			if braceRE.MatchString(line) {
				break
			}
		}
		best = []byte(strings.Join(lines, "\n"))
		log.Printf("[them] snip %2d/%2d: found function %v (%d lines) (%d bytes)", snip_id, len(snips), funcName, len(lines), len(best))

		out := snipName(snipDir, "them", s)
		if !fileExists(out) {
			// log.Printf("[them] Created snip %2d/%2d: %v (%d bytes)", snip_id, len(snips), out, len(best))
			err = writeFileAtomic(out, best, 0o644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var conflictMarkersRE = regexp.MustCompile(`(?s)<<<<<<.*?=======.*?>>>>>>>`)

// Meld runs:
//
//	git merge-file snip/X@NNNN.ours snip/X@NNNN.them snip/X@NNNN.base
//
// and writes meld/X (no @NNNN). On conflicts (exit 1), keep markers and
// create meld/X.rej containing just the conflicting regions.
func Meld(snips []Snip) error {
	for i, s := range snips {
		snip_id := i + 1
		ours := snipName(snipDir, "ours", s)
		them := snipName(snipDir, "them", s)
		base := snipName(snipDir, "base", s)

		meldPath := filepath.Join(meldDir, fmt.Sprintf("%s@%s.%s", s.Joined, lineTag(s.StartLine), "meld"))

		// Run: git merge-file -p -L ours -L base -L them <ours> <them> <base>
		out, exit, runErr := runGitMergeFile(ours, them, base)
		if runErr != nil && exit < 0 {
			// git not available or failed unexpectedly; synthesize a conflict block
			out = synthesizeConflict(ours, them, base)
			exit = 1
		}

		conflicts := ""
		conflictLocations := conflictMarkersRE.FindAllIndex(out, -1)
		if len(conflictLocations) > 0 {
			conflicts = fmt.Sprintf("%d CONFLICTS FOUND: ", len(conflictLocations))
		}

		log.Printf("[meld] snip %2d/%2d: %sWriting %v", snip_id, len(snips), conflicts, meldPath)

		must(writeFileAtomic(meldPath, out, 0o644))
		if exit == 1 {
			// conflicts present, extract rejects
			rej := extractRejects(out)
			rejPath := meldPath + ".rej"
			if len(rej) == 0 {
				// if we couldn't detect blocks, at least drop full merged
				rej = out
			}
			log.Printf("[meld] snip %2d/%2d: REJECTS: Writing %v", snip_id, len(snips), rejPath)
			err := writeFileAtomic(rejPath, rej, 0o644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Compare(snips []Snip) error {
	out := colorable.NewColorableStdout()
	updates := 0
	for i, s := range snips {
		snip_id := i + 1
		oursPath := snipName(snipDir, "ours", s)

		meldPath := filepath.Join(meldDir, fmt.Sprintf("%s@%s.%s", s.Joined, lineTag(s.StartLine), "meld"))

		oursContent, err := os.ReadFile(oursPath)
		if err != nil {
			return err
		}
		meldContent, err := os.ReadFile(meldPath)
		if err != nil {
			return err
		}

		diff := ""
		updated := ""
		if !bytes.Equal(oursContent, meldContent) {
			a := string(oursContent)
			b := string(meldContent)
			d := dmp.New()
			diffs := d.DiffMain(a, b, false)
			diff = d.DiffPrettyText(diffs)
			updated = ": UPDATED:"
			updates++
		}

		log.Printf("[comp] snip %2d/%2d: %v%v", snip_id, len(snips), meldPath, updated)
		if diff != "" {
			fmt.Fprintln(out, diff)
		}
	}
	log.Printf("%d of %d snippets have been updated", updates, len(snips))
	return nil
}

var srcLineRE = regexp.MustCompile(`^//\s*Snippet:\s*(https://github\.com/[^ \t\r\n]+)`)

func findSnips(file []byte) ([]Snip, error) {
	lines := splitLinesBytes(file)
	var out []Snip

	for i := 0; i < len(lines); i++ {
		m := srcLineRE.FindSubmatch(lines[i])
		if m == nil {
			continue
		}
		url := string(m[1])
		meta, err := parseGitHubURL(url)
		if err != nil {
			return nil, fmt.Errorf("parse url %q: %w", url, err)
		}

		// ours content: from this Snippet: line up to (but excluding) the next Snippet: line
		start := i
		j := i + 1
		for ; j < len(lines); j++ {
			if srcLineRE.Match(lines[j]) {
				break
			}
		}
		ours := bytes.Join(lines[start:j], []byte("\n"))

		out = append(out, Snip{
			URL:       url,
			Owner:     meta.Owner,
			Repo:      meta.Repo,
			BaseRef:   meta.Ref,
			Path:      meta.Path,
			StartLine: meta.Start,
			EndLine:   meta.End,
			Joined:    joinPathForName(meta.Path),
			Ours:      ours,
		})
		// continue after the block
		i = j - 1
	}
	return out, nil
}

type urlMeta struct {
	Owner string
	Repo  string
	Ref   string
	Path  string
	Start int
	End   int
}

func parseGitHubURL(u string) (urlMeta, error) {
	// Expect: https://github.com/<owner>/<repo>/blob/<ref>/<path>#Lstart[-Lend]
	// Example: https://github.com/golang/go/blob/e282cbb1/src/os/tempfile.go#L22-L24
	var m urlMeta
	noFrag, frag, _ := strings.Cut(u, "#")
	if frag == "" {
		return m, errors.New("missing #Lstart[-Lend] fragment")
	}
	parts := strings.Split(noFrag, "/")
	// ... https: '' github.com owner repo blob ref <path...>
	if len(parts) < 7 || parts[2] != "github.com" || parts[5] != "blob" {
		return m, fmt.Errorf("unsupported github URL format: %s", u)
	}
	m.Owner = parts[3]
	m.Repo = parts[4]
	m.Ref = parts[6]
	m.Path = strings.Join(parts[7:], "/")

	// Lines
	// Accept LNNN or LNNN-LMMM
	if !strings.HasPrefix(frag, "L") {
		return m, fmt.Errorf("unexpected fragment %q", frag)
	}
	segs := strings.Split(frag, "-")
	start, err := strconv.Atoi(strings.TrimPrefix(segs[0], "L"))
	if err != nil {
		return m, fmt.Errorf("bad start line in %q: %w", frag, err)
	}
	end := start
	if len(segs) == 2 {
		if !strings.HasPrefix(segs[1], "L") {
			return m, fmt.Errorf("bad end fragment %q", segs[1])
		}
		end, err = strconv.Atoi(strings.TrimPrefix(segs[1], "L"))
		if err != nil {
			return m, fmt.Errorf("bad end line in %q: %w", frag, err)
		}
	}
	m.Start = start
	m.End = end
	return m, nil
}

func joinPathForName(repoPath string) string {
	// Prefer trimming up to and including "src/" if present (per examples),
	// then join remaining components with underscores.
	p := repoPath
	if idx := strings.Index(p, "/src/"); idx >= 0 {
		p = p[idx+len("/src/"):]
	}
	parts := strings.Split(p, "/")
	return strings.Join(parts, "_")
}

func lineTag(n int) string {
	return fmt.Sprintf("%04d", n)
}

func extractLines(content []byte, start, end int) ([]byte, error) {
	if start <= 0 || end < start {
		return nil, fmt.Errorf("invalid range %d-%d", start, end)
	}
	sc := bufio.NewScanner(bytes.NewReader(content))
	var buf bytes.Buffer
	line := 0
	for sc.Scan() {
		line++
		if line < start {
			continue
		}
		if line > end {
			break
		}
		buf.Write(sc.Bytes())
		buf.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func httpGet(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Set a UA to avoid 403 by some CDNs
	req.Header.Set("User-Agent", "scanmerge/1.0 (+https://github.com)")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 4<<10))
		return nil, fmt.Errorf("GET %s: %s: %s", url, res.Status, string(b))
	}
	return io.ReadAll(res.Body)
}

func writeFileAtomic(path string, data []byte, mode os.FileMode) error { //nolint:unparam // quiet linter
	tmp := path + ".tmp"
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(tmp, data, mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func splitLines(s []byte) []string {
	// Normalize to \n
	return strings.Split(strings.ReplaceAll(string(s), "\r\n", "\n"), "\n")
}

func splitLinesBytes(s []byte) [][]byte {
	// Returns lines without trailing newline; keeps content stable.
	s = bytes.ReplaceAll(s, []byte("\r\n"), []byte("\n"))
	return bytes.Split(s, []byte("\n"))
}

func runGitMergeFile(ours, them, base string) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // always cancel to release resources
	// Use -p to print to stdout so we can capture it. Label streams for clarity.
	cmd := exec.CommandContext(ctx, "git", "merge-file", "-p",
		"-L", "ours", "-L", "base", "-L", "them",
		ours, base, them,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exit := 0
	if err != nil {
		// git merge-file returns exit code 1 on conflicts (which is OK for us).
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			exit = ee.ExitCode()
		} else {
			exit = -1
		}
	}
	// If stdout is empty (some versions write in-place), fall back to reading ours (merged into it),
	// but our command uses -p so stdout should have the merge result.
	out := stdout.Bytes()
	if len(out) == 0 {
		// Unexpected; synthesize something
		out = []byte(fmt.Sprintf("<<<<<<< ours\n%s\n=======\n%s\n>>>>>>> them\n",
			readOrEmpty(ours), readOrEmpty(them)))
		exit = 1
	}
	if exit > 1 {
		return out, exit, fmt.Errorf("git merge-file error: %s", strings.TrimSpace(stderr.String()))
	}
	return out, exit, nil
}

func readOrEmpty(p string) string {
	b, _ := os.ReadFile(p)
	return string(b)
}

// extractRejects pulls out the conflict regions (between <<<<<<< and >>>>>>>).
func extractRejects(merged []byte) []byte {
	var rej bytes.Buffer
	sc := bufio.NewScanner(bytes.NewReader(merged))
	in := false
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "<<<<<<<") {
			in = true
			rej.WriteString(line)
			rej.WriteByte('\n')
			continue
		}
		if strings.HasPrefix(line, ">>>>>>>") {
			rej.WriteString(line)
			rej.WriteByte('\n')
			in = false
			rej.WriteByte('\n')
			continue
		}
		if in {
			rej.WriteString(line)
			rej.WriteByte('\n')
		}
	}
	return rej.Bytes()
}

func synthesizeConflict(ours, them, _ string) []byte {
	return []byte(fmt.Sprintf(
		"<<<<<<< ours\n%s\n=======\n%s\n>>>>>>> them\n",
		readOrEmpty(ours), readOrEmpty(them)))
}

func cacheName(dir string, s Snip) string {
	sha8 := s.BaseRef
	if len(sha8) > 8 {
		sha8 = sha8[:8]
	}
	return filepath.Join(dir, sha8, s.Joined)
}

func cacheNameThem(dir string, s Snip) string {
	sha8 := s.ThemRef
	if len(sha8) > 8 {
		sha8 = sha8[:8]
	}
	return filepath.Join(dir, sha8, s.Joined)
}

func snipName(dir string, ext string, s Snip) string {
	return filepath.Join(dir, fmt.Sprintf("%s@%s.%s", s.Joined, lineTag(s.StartLine), ext))
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type repoInfo struct {
	DefaultBranch string `json:"default_branch"`
}
type commitInfo struct {
	SHA string `json:"sha"`
}

var branches = map[string]string{}

func getDefaultBranch(owner, repo string) (string, error) {
	key := owner + "/" + repo

	branch, ok := branches[key]
	if ok {
		return branch, nil
	}
	var info repoInfo
	if err := ghJSON(fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), &info); err != nil {
		return "", err
	}
	branches[key] = info.DefaultBranch
	return info.DefaultBranch, nil
}

var shas = map[string]string{}

func getCommitSHA(owner, repo, ref string) (string, error) {
	key := owner + "/" + repo + "/" + ref

	sha, ok := shas[key]
	if ok {
		return sha, nil
	}
	var c commitInfo
	if err := ghJSON(fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, ref), &c); err != nil {
		return "", err
	}
	shas[key] = c.SHA
	return c.SHA, nil
}

func addAuth(req *http.Request) {
	if tok := os.Getenv("GITHUB_TOKEN"); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
}

func ghJSON(u string, out any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel() // always cancel to release resources
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	addAuth(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("GET %s: %s", u, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
