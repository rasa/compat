package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/text/unicode/norm"
)

type CaseFoldingMode int

const (
	CaseFoldingNone    CaseFoldingMode = iota // strict byte comparison (e.g., ext4)
	CaseFoldingASCII                          // ASCII-only folding (e.g., FAT/exFAT)
	CaseFoldingUnicode                        // full Unicode folding (e.g., NTFS, APFS)
)

const (
	perm600 os.FileMode = 0o600
	perm755 os.FileMode = 0o755
)

func (m CaseFoldingMode) String() string {
	switch m {
	case CaseFoldingNone:
		return "None"
	case CaseFoldingASCII:
		return "ASCII"
	case CaseFoldingUnicode:
		return "Unicode"
	default:
		return "Unknown" //nolint:goconst
	}
}

type CaseProcessing struct {
	CaseSensitive     bool
	CasePreserving    bool
	AccentInsensitive bool
	NormalizationForm string // "NFC", "NFD", "None", or "Unknown"
	SortOrder         string // "ASCII", "Unicode", "Unsorted"
	CaseFoldingMode   CaseFoldingMode
}

func main() {
	dir := "."

	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	var (
		cp  CaseProcessing
		err error
	)

	cp.CaseSensitive, err = DetectCaseSensitivity(dir)
	if err != nil {
		fmt.Println("CaseSensitivity error:", err)
	}

	cp.CasePreserving, err = DetectCasePreserving(dir)
	if err != nil {
		fmt.Println("CasePreserving error:", err)
	}

	cp.AccentInsensitive, err = DetectAccentInsensitivity(dir)
	if err != nil {
		fmt.Println("AccentInsensitivity error:", err)
	}

	cp.NormalizationForm, err = DetectNormalization(dir)
	if err != nil {
		fmt.Println("Normalization error:", err)
	}

	cp.SortOrder, _, err = DetectSortOrder(dir)
	if err != nil {
		fmt.Println("SortOrder error:", err)
	}

	cp.CaseFoldingMode, err = DetectCaseFoldingMode(dir)
	if err != nil {
		fmt.Println("CaseFolding error:", err)
	}

	fmt.Printf("\nFilesystem Case Processing:\n")
	fmt.Printf("  CaseSensitive:     %v\n", cp.CaseSensitive)
	fmt.Printf("  CasePreserving:    %v\n", cp.CasePreserving)
	fmt.Printf("  AccentInsensitive: %v\n", cp.AccentInsensitive)
	fmt.Printf("  NormalizationForm: %s\n", cp.NormalizationForm)
	fmt.Printf("  SortOrder:         %s\n", cp.SortOrder)
	fmt.Printf("  CaseFoldingMode:   %s\n", cp.CaseFoldingMode)
}

//
// ────────────────────────────── CASE SENSITIVITY ──────────────────────────────
//

func DetectCaseSensitivity(dir string) (bool, error) {
	tmp := filepath.Join(dir, ".case_sense_test")
	os.RemoveAll(tmp) //nolint:gosec

	err := os.MkdirAll(tmp, perm755) //nolint:gosec
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(tmp)

	lower := filepath.Join(tmp, "file.txt")
	upper := filepath.Join(tmp, "FILE.txt")

	f1, err := os.OpenFile(lower, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm600) //nolint:gosec
	if err != nil {
		return false, fmt.Errorf("create lower: %w", err)
	}

	f1.Close()

	f2, err := os.OpenFile(upper, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm600) //nolint:gosec

	switch {
	case err == nil:
		f2.Close()

		return true, nil
	case os.IsExist(err):
		return false, nil
	default:
		return false, fmt.Errorf("create upper: %w", err)
	}
}

//
// ────────────────────────────── CASE PRESERVATION ──────────────────────────────
//

func DetectCasePreserving(dir string) (bool, error) {
	tmp := filepath.Join(dir, ".case_preserve_test")
	os.RemoveAll(tmp) //nolint:gosec

	err := os.MkdirAll(tmp, perm755) //nolint:gosec
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(tmp)

	name := "MiXeDcAsE.txt"
	path := filepath.Join(tmp, name)

	err = os.WriteFile(path, []byte("data"), perm600) //nolint:gosec
	if err != nil {
		return false, err
	}

	entries, err := os.ReadDir(tmp)
	if err != nil {
		return false, err
	}

	for _, e := range entries {
		if e.Name() == name {
			return true, nil
		}
	}

	return false, nil
}

//
// ────────────────────────────── ACCENT INSENSITIVITY ──────────────────────────────
//

func DetectAccentInsensitivity(dir string) (bool, error) {
	tmp := filepath.Join(dir, ".accent_test")
	os.RemoveAll(tmp) //nolint:gosec

	err := os.MkdirAll(tmp, perm755) //nolint:gosec
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(tmp)

	a := filepath.Join(tmp, "e.txt")
	b := filepath.Join(tmp, "é.txt")

	fa, err := os.OpenFile(a, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm600) //nolint:gosec
	if err != nil {
		return false, fmt.Errorf("create a: %w", err)
	}

	fa.Close()

	fb, err := os.OpenFile(b, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm600) //nolint:gosec

	switch {
	case err == nil:
		fb.Close()

		return false, nil
	case os.IsExist(err):
		return true, nil
	default:
		return false, fmt.Errorf("create b: %w", err)
	}
}

//
// ────────────────────────────── NORMALIZATION FORM ──────────────────────────────
//

// DetectNormalization writes both NFC and NFD forms of "é.txt" and checks what survives.
func DetectNormalization(dir string) (string, error) {
	tmp := filepath.Join(dir, ".normalize_test")
	os.RemoveAll(tmp) //nolint:gosec

	err := os.MkdirAll(tmp, perm755) //nolint:gosec
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)

	nfcName := "é"       // U+00E9
	nfdName := "e\u0301" // U+0065 + U+0301

	nfcPath := filepath.Join(tmp, nfcName)
	nfdPath := filepath.Join(tmp, nfdName)

	err = os.WriteFile(nfcPath, nil, perm600) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("write NFC: %w", err)
	}

	err = os.WriteFile(nfdPath, nil, perm600) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("write NFD: %w", err)
	}

	entries, err := os.ReadDir(tmp)
	if err != nil {
		return "", fmt.Errorf("readdir: %w", err)
	}

	switch len(entries) {
	case 1:
		name := entries[0].Name()

		switch {
		case norm.NFC.IsNormalString(name):
			return "NFC", nil
		case norm.NFD.IsNormalString(name):
			return "NFD", nil
		default:
			return "Unknown", nil
		}
	case 2: //nolint:mnd
		return "None", nil
	default:
		return "Unknown", nil
	}
}

//
// ────────────────────────────── SORT ORDER ──────────────────────────────
//

func DetectSortOrder(dir string) (sortOrder string, unicodeAware bool, err error) {
	tmp := filepath.Join(dir, ".sort_test")
	os.RemoveAll(tmp) //nolint:gosec

	err = os.MkdirAll(tmp, perm755) //nolint:gosec
	if err != nil {
		return "", false, err
	}
	defer os.RemoveAll(tmp)

	names := []string{"a.txt", "B.txt", "é.txt", "Ω.txt"}
	for _, n := range names {
		_ = os.WriteFile(filepath.Join(tmp, n), nil, perm600) //nolint:gosec
	}

	entries, err := os.ReadDir(tmp)
	if err != nil {
		return "", false, err
	}

	raw := make([]string, len(entries))
	for i, e := range entries {
		raw[i] = e.Name()
	}

	ascii := slices.Clone(raw)
	slices.Sort(ascii)

	unicode := slices.Clone(raw)
	slices.SortFunc(unicode, func(a, b string) int {
		ar, br := []rune(a), []rune(b)
		for i := 0; i < len(ar) && i < len(br); i++ {
			if ar[i] == br[i] {
				continue
			}

			if ar[i] < br[i] {
				return -1
			}

			return 1
		}

		return len(ar) - len(br)
	})

	switch {
	case slices.Equal(raw, ascii):
		return "ASCII", false, nil
	case slices.Equal(raw, unicode):
		return "Unicode", true, nil
	default:
		return "Unsorted", false, nil
	}
}

//
// ────────────────────────────── CASE FOLDING (ASCII vs UNICODE) ──────────────────────────────
//

// DetectCaseFoldingMode detects whether the FS performs Unicode or ASCII-only case folding.
func DetectCaseFoldingMode(dir string) (CaseFoldingMode, error) {
	tmp := filepath.Join(dir, ".fold_test")
	os.RemoveAll(tmp) //nolint:gosec

	err := os.MkdirAll(tmp, perm755) //nolint:gosec
	if err != nil {
		return CaseFoldingNone, err
	}
	defer os.RemoveAll(tmp)

	// 1. Write lowercase ä.txt
	lower := filepath.Join(tmp, "ä.txt")

	err = os.WriteFile(lower, []byte("x"), perm600) //nolint:gosec,lll
	if err != nil {
		return CaseFoldingNone, err
	}

	// 2. Try creating uppercase form (Ä.TXT)
	upper := filepath.Join(tmp, "Ä.TXT")

	f, err := os.OpenFile(upper, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm600) //nolint:gosec
	if err == nil {
		f.Close()
		// both allowed → no folding
	} else if os.IsExist(err) {
		return CaseFoldingUnicode, nil
	}

	// 3. Test ASCII-only (z.txt vs Z.TXT)
	err = os.WriteFile(filepath.Join(tmp, "z.txt"), []byte("z"), perm600) //nolint:gosec
	if err != nil {
		return CaseFoldingNone, err
	}

	_, err = os.OpenFile(filepath.Join(tmp, "Z.TXT"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm600) //nolint:gosec,lll
	if err != nil {
		if os.IsExist(err) {
			// ASCII folds but ä/Ä did not
			return CaseFoldingASCII, nil
		}
	}

	return CaseFoldingNone, nil
}
