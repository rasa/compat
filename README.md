# compat

[![Build](https://github.com/rasa/compat/actions/workflows/build.yml/badge.svg)](https://github.com/rasa/compat/actions/workflows/build.yml)
[![CodeQL](https://github.com/rasa/compat/actions/workflows/codeql.yml/badge.svg)](https://github.com/rasa/compat/actions/workflows/codeql.yml)
[![CodeFactor Grade](https://img.shields.io/codefactor/grade/github/rasa/compat)](https://www.codefactor.io/repository/github/rasa/compat)
[![Codecov](https://codecov.io/gh/rasa/compat/branch/main/graph/badge.svg)](https://codecov.io/gh/rasa/compat)
[![MegaLinter](https://github.com/rasa/compat/actions/workflows/mega-linter.yml/badge.svg)](https://github.com/rasa/compat/actions/workflows/mega-linter.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/rasa/compat.svg)](https://pkg.go.dev/github.com/rasa/compat)<br/>
[![Release](https://img.shields.io/github/v/release/rasa/compat.svg)](https://github.com/rasa/compat/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rasa/compat)](go.mod)
[![Last Commit](https://img.shields.io/github/last-commit/rasa/compat.svg)](https://github.com/rasa/compat/commits)
[![Lines of Code](https://img.shields.io/badge/LoC-16.4k-blue)](https://github.com/rasa/compat)
[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![License](https://img.shields.io/github/license/rasa/compat)](LICENSE)

**Portable file identity, extended metadata, and reliable file operations for Go.**

`compat` provides a consistent API for file information that Go's standard
`os.FileInfo` interface does not expose portably, including:

- access, birth, modification, and metadata-change times;
- unique file and partition identifiers;
- hard-link counts;
- user and group IDs and names;
- cross-platform file permission setting;
- atomic file replacement;
- filesystem and operating-system compatibility helpers.

The package supports Linux, macOS, Windows, BSD systems, Android, iOS, Plan 9,
Solaris, illumos, JavaScript/WebAssembly, WASI, TinyGo, and other targets
supported by Go.

## Contents

<!--ts-->
* [Why compat?](#why-compat)
* [Use cases](#use-cases)
* [Installation](#installation)
* [Quick start](#quick-start)
* [Read portable file metadata](#read-portable-file-metadata)
* [Write a file atomically](#write-a-file-atomically)
* [Compare file identity](#compare-file-identity)
* [Core API](#core-api)
* [Metadata and identity](#metadata-and-identity)
* [File operations](#file-operations)
* [Behavior and limitations](#behavior-and-limitations)
* [Platform support](#platform-support)
* [Extended file metadata](#extended-file-metadata)
* [Other operations](#other-operations)
* [Environment variables](#environment-variables)
* [Runtime configuration](#runtime-configuration)
* [Test configuration](#test-configuration)
* [COMPAT_DEBUG](#compat_debug)
* [COMPAT_DEBUG_FS](#compat_debug_fs)
* [COMPAT_DEBUG_FS_PATH](#compat_debug_fs_path)
* [COMPAT_DEBUG_FS_SIZE](#compat_debug_fs_size)
* [Comparing Stat Across Linux and Windows](#comparing-stat-across-linux-and-windows)
* [Related cross-platform Go packages](#related-cross-platform-go-packages)
* [Contributing](#contributing)
* [License](#license)
<!--te-->

## Why compat?

The Go standard library intentionally exposes a small, portable file interface.
Applications such as backup tools, synchronization engines, file indexes, and
duplicate detectors often need additional information that otherwise requires
different code for each operating system.

| Capability | Go standard library | `compat` |
|---|:---:|:---:|
| Basic name, size, mode, and modification time | Yes | Yes |
| Access time | Platform-specific | Portable API where supported |
| Birth or creation time | Platform-specific | Portable API where supported |
| Metadata-change time | Platform-specific | Portable API where supported |
| Hard-link count | Platform-specific | Portable API where supported |
| User and group IDs | Platform-specific | Portable API |
| User and group names | No unified API | Portable API |
| File and partition identity | Platform-specific | Portable API |
| Windows SID-to-POSIX ID mapping | No | Yes |
| Atomic file replacement | No unified API | Yes, where supported |
| Filesystem capability checks | Limited | Yes |

`compat` does not replace the standard library for ordinary file access. Use
`os`, `io`, `io/fs`, and `path/filepath` when their APIs already provide
everything your application needs.

## Use cases

`compat` is intended for software that must interpret file metadata consistently
across operating systems, including:

- backup and restore tools;
- file synchronization software;
- duplicate-file detectors;
- file indexes and search tools;
- installers and self-updaters;
- archival and migration utilities;
- cross-platform development tools;
- applications that need stable file identity on Windows and Unix-like systems.

## Installation

```console
go get github.com/rasa/compat
```

Documentation for the public API is available on
[pkg.go.dev](https://pkg.go.dev/github.com/rasa/compat).

## Quick start

### Read portable file metadata

```go
package main

import (
	"fmt"
	"log"

	"github.com/rasa/compat"
)

func main() {
	info, err := compat.Stat("example.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name:       %s\n", info.Name())
	fmt.Printf("Size:       %d bytes\n", info.Size())
	fmt.Printf("Mode:       0o%o (%v)\n", info.Mode(), info.Mode())
	fmt.Printf("Modified:   %s\n", info.ModTime())
	fmt.Printf("Accessed:   %s\n", info.ATime())

	if compat.SupportsBTime() {
		fmt.Printf("Created:    %s\n", info.BTime())
	}

	if compat.SupportsCTime() {
		fmt.Printf("Changed:    %s\n", info.CTime())
	}

	fmt.Printf("Owner:      %s (%d)\n", info.User(), info.UID())
	fmt.Printf("Group:      %s (%d)\n", info.Group(), info.GID())
	fmt.Printf("Links:      %d\n", info.Links())
	fmt.Printf("File ID:    %d:%d\n", info.PartitionID(), info.FileID())

	// Some extended values may require an additional system call.
	if err := info.Error(); err != nil {
		log.Printf("some extended metadata could not be read: %v", err)
	}
}
```

`Stat` follows symbolic links. Use `Lstat` when the returned metadata should
describe the symbolic link itself.

### Write a file atomically

```go
package main

import (
	"log"

	"github.com/rasa/compat"
)

func main() {
	data := []byte("configuration data\n")

	err := compat.WriteFile(
		"config.txt",
		data,
		0o600,
		compat.WithAtomicity(true), // optional
	)
	if err != nil {
		log.Fatal(err)
	}
}
```

With `WithAtomicity(true)`, the destination is either fully replaced or left
unchanged when the operating system and filesystem support atomic replacement.

On Plan 9, creating a new file atomically is supported, but atomically replacing
an existing file is not. `WriteFile` returns an error matching
`errors.ErrUnsupported` in that case. Applications that explicitly accept a
non-atomic fallback can also pass:

```go
compat.WithNonAtomicReplace(true)
```

Atomic replacement and durable persistence are different guarantees. An atomic
rename prevents readers from observing a partially written destination; it does
not necessarily guarantee that the data has reached permanent storage after a
power loss.

### Compare file identity

Use the file and partition identifiers returned by `Stat` to determine whether
two paths refer to the same underlying file:

```go
package main

import (
	"fmt"
	"log"

	"github.com/rasa/compat"
)

func main() {
	first, err := compat.Stat("first.txt")
	if err != nil {
		log.Fatal(err)
	}
	second, err := compat.Stat("second.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("same file:      %v\n", compat.SameFile(first, second))
	fmt.Printf("same partition: %v\n", compat.SamePartition(first, second))

	// or alternatively:
	same1, err := compat.SameFiles("first.txt", "second.txt")
	if err != nil {
		log.Fatal(err)
	}
	same2, err := compat.SamePartitions("first.txt", "second.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("same file:      %v\n", same1)
	fmt.Printf("same partition: %v\n", same2)
}
```

## Core API

### Metadata and identity

| Function or type | Purpose |
|---|---|
| `Stat` | Returns extended information for a path and follows symbolic links |
| `Lstat` | Returns extended information about a path without following its final symbolic link |
| `Fstat` | Returns extended information for an open file where supported |
| `FileInfo` | Extends `os.FileInfo` with portable metadata and identity methods |
| `SameFile` | Reports whether two `compat.FileInfo` values describe the same file |
| `SamePartition` | Reports whether two files reside on the same partition |
| `PartitionType` | Reports the underlying filesystem or partition type |
| `SupportsATime` | Reports operating-system support for access time |
| `SupportsATimeSetting` | Reports support for setting access time |
| `SupportsBTime` | Reports operating-system support for birth time |
| `SupportsCTime` | Reports operating-system support for metadata-change time |
| `SupportsFstat` | Reports support for `Fstat` |
| `SupportsLinks` | Reports support for hard-link counts |
| `SupportsSymlinks` | Reports operating-system support for symbolic links |
| `SupportsUmask` | Reports support for `Umask` |
| `UserIDSource` | Describes how user IDs are represented on the current platform |

### File operations

The package includes cross-platform variants of common file operations,
including:

- `Chmod` and `Fchmod`
- `Create`, `CreateTemp`
- `Link` and `Symlink`
- `Mkdir`, `Mkdirall` and `MkdirTemp`
- `Nice` and `Renice`
- `Open`, and `OpenFile`
- `ReadDir` and `WalkDir`
- `Rename`
- `Remove`, and `RemoveAll`
- `Stat`, `Fstat` and `LStat`
- `Umask`
- `WriteFile` and `WriteReader`

Several operations accept functional options:

| Option | Purpose |
|---|---|
| `WithAtomicity` | Requests atomic file creation or replacement |
| `WithNonAtomicReplace` | Allows a non-atomic replacement where atomic replacement is unavailable |
| `WithFileMode` | Sets the requested file mode |
| `WithDefaultFileMode` | Changes the default mode used when no explicit mode is supplied |
| `WithKeepFileMode` | Preserves the mode of an existing destination |
| `WithFlags` | Adds file-open flags |
| `WithReadOnlyMode` | Controls Windows read-only attribute handling |
| `WithRetrySeconds` | Retries selected operations for a bounded period |
| `WithSetSymlinkOwner` | Requests ownership adjustment for Windows symbolic links |

An option may apply only to the operations documented for it. Options that
control Windows-specific behavior are ignored on other operating systems.

## Behavior and limitations

Cross-platform filesystem behavior cannot be made identical in every case.
Callers should account for the following:

- **Filesystem support varies.** An operating system may support a metadata
  field even when a particular local, network, virtual, or removable filesystem
  does not.
- **Capability functions describe platform support.** A `Supports*` result does
  not override filesystem-specific limitations, mount options, permissions, or
  runtime restrictions.
- **Unavailable values use documented zero or sentinel values.** Check
  `FileInfo.Error()` after reading extended values that may require additional
  system calls.
- **Windows permissions are represented through ACLs.** POSIX mode bits are
  mapped to Windows access-control entries and cannot express every Windows ACL
  configuration.
- **User and group identity differ by operating system.** Windows SIDs are
  mapped to integer values compatible with Cygwin, MSYS2, and Git for Windows.
  Plan 9 uses hashes of user and group names.
- **Atomic replacement is not durability.** A successful atomic rename does not
  by itself guarantee persistence after a system or storage failure.
- **Symbolic-link support may require privileges or configuration.** This is
  especially relevant on Windows and restricted mobile environments.
- **Network and virtual filesystems may have different semantics.** Test the
  actual filesystems your application supports.

The project is currently pre-v1. Review [CHANGELOG.md](CHANGELOG.md) for
breaking changes before upgrading between releases.

## Platform support

### Extended file metadata

`Stat`, `Fstat`, and `Lstat` return a `FileInfo` value. Support for extended metadata
depends on both the operating system and the underlying filesystem.

| OS           | PartitionID()/ <br/>FileID()* | Links()* | ATime()*<br/>(last<br/>*A*ccessed) | BTime()*<br/>(*B*irthed/<br/>created) | CTime()*<br/>(metadata<br/>last *C*hanged) | UID()/GID() |
|--------------|--------|--------|------|--------|------|-------|
| AIX          | ✅     | ✅    | ✅   | ❌    | ✅   | ✅   |
| Android      | ✅     | ✅§   | ✅   | ❌    | ✅   | ✅   |
| Dragonfly    | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| FreeBSD      | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| Illumos      | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| iOS          | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| Js/<br/>WASM | ✅     | ✅    | ✅   | ❌    | ✅   | ✅   |
| Linux        | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| macOS        | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| NetBSD       | ✅     | ✅    | ✅   | ✅    | ✅   | ✅   |
| OpenBSD      | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| Plan9        | ✅     | ❌    | ✅   | ❌    | ❌   | ☑️   |
| Solaris      | ✅     | ✅    | ✅   | ✖️    | ✅   | ✅   |
| Wasip1/<br/>WASM | ✅ | ✅    | ✅†  | ❌    | ✅   | ✅   |
| Windows      | ✅     | ✅    | ✅   | ✅    | ✅   | ✅‡  |
<!--           | PartID+ | Links | ATime | BTime | CTime | UID+ | -->

**Key**

- ✅ Fully supported by the operating-system implementation.
- ☑️ The `UID()` and `GID()` values are 32-bit hashes of the user and group names.
- ✖️ Not currently implemented although the operating system may support it.
- ❌ Not implemented because the operating system appears not to provide it.
<!-- 🚧 planned to be implemented.<br/> -->

\* Actual support depends on the underlying filesystem. See [Comparison of file systems](https://wikipedia.org/wiki/Comparison_of_file_systems#Metadata) for details.<br/>
† Not supported when compiled with TinyGo.<br/>
‡ Windows values use the same SID-to-integer mapping as Cygwin, MSYS2, and Git
for Windows.<br/>
§ Android 7/API 24 and later appear to disallow hard-link creation by default.

### Other operations

| OS           | Chmod()* | Fstat() | Nice()/<br/>Renice() | PartitionType() | Symlink() | Umask() |
|--------------|----------|----------|--------|------|-------|------|
| AIX          | ✅       | ✅‖     | ✅    | ✅*  | ✅   | ✅   |
| Android      | ✅       | ✅      | ✅    | ✅   | ✅   | ✅   |
| Dragonfly    | ✅       | ✅‖     | ✅    | ✅   | ✅   | ✅   |
| FreeBSD      | ✅       | ✅      | ✅    | ✅‡  | ✅   | ✅   |
| Illumos      | ✅       | ✅‖     | ✅    | ✅   | ✅   | ✅   |
| iOS          | ✅       | ✅      | ☑️    | ✅   | ✅   | ✅   |
| Js/<br/>WASM | ❌       | ✅‖     | ☑️    | ✅   | ❌   | ✅   |
| Linux        | ✅       | ✅      | ✅    | ✅   | ✅   | ✅   |
| macOS        | ✅       | ✅      | ✅    | ✅   | ✅   | ✅   |
| NetBSD       | ✅       | ✅‖     | ✅    | ✅‡  | ✅   | ✅   |
| OpenBSD      | ✅       | ✅‖     | ✅    | ✅‡  | ✅   | ✅   |
| Plan9        | ✅       | ✅‖     | ✅    | ✅   | ❌   | ❌   |
| Solaris      | ✅       | ✅‖     | ✅    | ✅   | ✅   | ✅   |
| Wasip1/<br/>WASM | ❌   | ✅‖     | ☑️    | ✅   | ❌   | ✅†  |
| Windows      | ✅       | ✅      | ✅    | ✅   | ✅   | ✅§  |
<!--           | Chmod    | Fstat   | Nice  | Part | Symln | Umask | -->

**Key**

- ✅ Fully supported.
- ☑️ `Nice` always returns `0`, and `Renice` performs no operation.
- ✖️ Not currently implemented although the operating system may support it.
- ❌ Not implemented because the operating system appears not to provide it.

\* Actual support depends on the underlying filesystem. See [Comparison of file systems](https://wikipedia.org/wiki/Comparison_of_file_systems#Metadata) for details.<br/>
† Not supported when compiled with TinyGo.<br/>
‡ Not supported on `openbsd/ppc64`, `netbsd/386`, `freebsd/riscv64`, or
`aix/ppc64` with CGO because of compilation limitations.<br/>
§ Implemented through the `UMASK=0NNN` environment variable.
‖ Will fall on relative filenames where the current directory is changed after
opening the file.

## Environment variables

### Runtime configuration

#### `UMASK` — Windows only

Sets the initial default umask. When unset, `0o022` is used.

### Test configuration

The following variables are used only by the project's test suite.

#### `COMPAT_DEBUG`

- Set to `DEBUG` to emit additional debugging information.
- Set to `DUMP` to dump Windows ACL information.

#### `COMPAT_DEBUG_FS`

Selects a filesystem to test. Set it to `All` to test every available
filesystem.

Supported test filesystems include:

**Linux**

- Btrfs
- exFAT
- ext2
- ext3
- ext4
- F2FS
- FAT
- FAT32
- NTFS
- ReiserFS
- XFS

**macOS**

- APFS
- exFAT
- FAT32
- HFS+
- HFS+J
- HFSX
- JHFS+
- JHFS+X
- UDF

**Windows**

- exFAT
- FAT32
- FAT — currently disabled
- NTFS
- ReFS

#### `COMPAT_DEBUG_FS_PATH`

Sets the path used to mount test filesystems.

Defaults:

| OS | Path |
|---|---|
| Linux | `/mnt/` |
| macOS | `/Volumes/` |
| Windows | `Z:\` |

#### `COMPAT_DEBUG_FS_SIZE`

Sets the size of the virtual test volume. The default is `2GB`.

#### `COMPAT_DEBUG_PERM` — Windows only

Sets a specific file-permission value to test. When unset, the suite tests all
511 nonzero combinations of the user, group, and other permission bits.


## Comparing Stat Across Linux and Windows


```go
package main

import (
	"log"
	"os"

	"github.com/rasa/compat"
)

func main() {
	hello := []byte("Hello World!")
	err := compat.WriteFile("hello.txt", hello, 0o654)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := compat.Stat("hello.txt")
	if err != nil {
		log.Fatal(err)
	}

	print(fi.String())
}
```
On Linux, the above program prints:
```text
Name:   hello.txt
Size:   12
Mode:   0o654 (-rw-r-xr--)
ModTime:2025-08-14 09:25:27.190602462 -0700 PDT // last Modified
ATime:  2025-08-14 09:25:27.190602462 -0700 PDT // last Accessed
BTime:  2025-08-14 09:25:27.190602462 -0700 PDT // Birthed/created
CTime:  2025-08-14 09:25:27.190602462 -0700 PDT // metadata last Changed
IsDir:  false
Links:  1                                       // number of hard links
UID:    1000 (ross)                             // user ID
GID:    1000 (ross)                             // group ID
PartID: 64512                                   // unique partition (device) ID
FileID: 18756713                                // unique file ID on the partition
```
On Windows, the program prints:
```text
Name:   hello.txt
Size:   12
Mode:   0o654 (-rw-r-xr--)
ModTime:2025-08-14 09:28:50.4214934 -0700 PDT
ATime:  2025-08-14 09:28:50.4214934 -0700 PDT
BTime:  2025-08-14 09:28:50.4209614 -0700 PDT
CTime:  2025-08-14 09:28:50.4214934 -0700 PDT
IsDir:  false
Links:  1
UID:    197609 (domain\ross)
GID:    197121 (domain\None)
PartID: 8
FileID: 844424931319952
```
icacls shows:
```
icacls hello.txt

hello.txt domain\ross:(R,W,D)
          domain\None:(RX)
          Everyone:(R)
```
Powershell shows:
```
powershell -command "Get-Acl hello.txt | Format-List"

Path   : Microsoft.PowerShell.Core\FileSystem::C:\path\to\hello.txt
Owner  : domain\ross
Group  : domain\None
Access : Everyone Allow  Read, Synchronize
         domain\None Allow  ReadAndExecute, Synchronize
         domain\ross Allow  Write, Delete, Read, Synchronize
Audit  :
Sddl   : O:S-1-5-21-2970224322-3395479738-1485484954-1001G:S-1-5-21-2970224322-3395479738-1485484954-513D:P(A;;FR;;;WD)(A;;0x1200a9;;;S-1-5-21-2970224322-3395479738-1485484954-5
         19f;;;S-1-5-21-2970224322-3395479738-1485484954-1001)
```
Cygwin's stat shows:
```
$ stat hello.txt
  File: hello.txt
  Size: 12              Blocks: 1          IO Block: 65536  regular file
Device: 0,8     Inode: 844424931319952  Links: 1
Access: (0754/-rwxr-xr--)  Uid: (197609/    ross)   Gid: (197121/    None)
Access: 2025-08-14 09:28:50.421493400 -0700
Modify: 2025-08-14 09:28:50.421493400 -0700
Change: 2025-08-14 09:28:50.421493400 -0700
 Birth: 2025-08-14 09:28:50.420961400 -0700
```
Git for Windows's stat shows:
```
$ stat hello.txt
  File: hello.txt
  Size: 11              Blocks: 1          IO Block: 65536  regular file
Device: 8h/8d   Inode: 844424931319952  Links: 1
Access: (0644/-rw-r--r--)  Uid: (197609/    ross)   Gid: (197121/ UNKNOWN)
Access: 2025-08-14 09:30:06.729683700 -0700
Modify: 2025-08-14 09:28:50.421493400 -0700
Change: 2025-08-14 09:28:50.421493400 -0700
 Birth: 2025-08-14 09:28:50.420961400 -0700
```

## Related cross-platform Go packages

Go applications commonly rely on focused packages that normalize one category
of operating-system behavior. The following projects solve problems adjacent to,
but generally separate from, portable file metadata.

| Package | Capability | Purpose |
|---|---|---|
| [`99designs/keyring`](https://github.com/99designs/keyring) | Credential storage | Provides a common interface for macOS Keychain, Windows Credential Manager, Secret Service, KWallet, `pass`, KeyCtl, and other credential stores |
| [`adrg/xdg`](https://github.com/adrg/xdg) | Application directories | Provides standard configuration, data, cache, state, runtime, and user-directory paths across Unix, Windows, macOS, and Plan 9 |
| [`distatus/battery`](https://github.com/distatus/battery) | Battery and AC status | Reports battery capacity, charging state, and power information across supported operating systems |
| [`ebitengine/purego`](https://github.com/ebitengine/purego) | Native-library interoperability | Calls functions in native shared libraries without requiring conventional Cgo wrappers |
| [`fsnotify/fsnotify`](https://github.com/fsnotify/fsnotify) | Filesystem change notifications | Wraps facilities such as inotify, kqueue, `ReadDirectoryChangesW`, and FEN behind one event API |
| [`gofrs/flock`](https://github.com/gofrs/flock) | File locking | Provides a portable file-locking interface over platform-specific locking mechanisms |
| [`hymkor/trash-go`](https://github.com/hymkor/trash-go) | Move files to trash | Uses the Windows Recycle Bin and experimentally supports Freedesktop-style trash directories on non-Windows systems |
| [`mattn/go-colorable`](https://github.com/mattn/go-colorable) | Portable colored output | Supports ANSI-colored output through Windows console handling and ordinary writers on other systems |
| [`mattn/go-isatty`](https://github.com/mattn/go-isatty) | Terminal detection | Determines whether a file descriptor refers to a terminal or character device |
| [`pkg/browser`](https://github.com/pkg/browser) | Open URLs and files | Opens URLs, files, or reader contents through the platform's browser facilities |
| [`shirou/gopsutil`](https://github.com/shirou/gopsutil) | System and process information | Provides cross-platform CPU, memory, host, disk, network, and process information |
| [`shirou/gopsutil/process`](https://pkg.go.dev/github.com/shirou/gopsutil/v4/process) | Process information and control | Lists processes and provides operations such as suspend, resume, foreground or background detection, and CPU affinity where supported |
| [`tklauser/numcpus`](https://github.com/tklauser/numcpus) | Available CPU count | Determines the number of CPUs available to the current process using platform-specific facilities |

## Contributing

Issues and pull requests are welcome.

Filesystem behavior is highly platform-specific. Bug reports are most useful
when they include:

- the Go version;
- `GOOS` and `GOARCH`;
- the operating-system version;
- the filesystem or volume type;
- whether CGO or TinyGo was used;
- a minimal reproducer;
- the expected and actual results.

Before submitting a pull request, run the applicable tests and format the code
with the standard Go tools.

## License

`compat` is available under the [MIT License](LICENSE).
