// +build linux

package fs

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func Chall(fi *FileInfo) error {
	if err := os.Lchown(fi.Name, int(fi.Uid), int(fi.Gid)); err != nil {
		return fmt.Errorf("chown %q, %s", fi.Name, err)
	}

	// TODO utime not work for symlink

	if (fi.Mode & os.ModeType) != os.ModeSymlink {
		err := unix.Utimes(fi.Name, []unix.Timeval{
			{
				Sec:  fi.Atime.Unix(),
				Usec: (fi.Atime.UnixNano() / 1000000) % 1000,
			},
			{
				Sec:  fi.Mtime.Unix(),
				Usec: (fi.Mtime.UnixNano() / 1000000) % 1000,
			},
		})
		if err != nil {
			return fmt.Errorf("utimes %q, atime %s, mtime %s, %s", fi.Name, fi.Atime, fi.Mtime, err)
		}
	}
	return nil
}

func MkdirAll(fi *FileInfo) error {
	if err := os.MkdirAll(fi.Name, fi.Mode); err != nil {
		return fmt.Errorf("mkdirall %q, %s", fi.Name, err)
	}
	return Chall(fi)
}

func Symlink(fi *FileInfo) error {
	if err := os.Symlink(fi.Linkname, fi.Name); err != nil {
		return fmt.Errorf("symlink %s -> %s, %s", fi.Name, fi.Linkname, err)
	}
	return Chall(fi)
}

// Mknod creates a filesystem node (file, device special file or named pipe) named path
// with attributes specified by mode and dev.
func Mknod(fi *FileInfo) error {
	var mode uint32
	switch fi.Mode & os.ModeType {
	case os.ModeDevice:
		mode = unix.S_IFBLK
	case os.ModeDevice | os.ModeCharDevice:
		mode = unix.S_IFCHR
	}
	if err := unix.Mknod(fi.Name, mode, int(mkdev(fi.Major, fi.Minor))); err != nil {
		return fmt.Errorf("mknod %q with (%d,%d), %s", fi.Name, fi.Major, fi.Minor, err)
	}
	if err := os.Chmod(fi.Name, fi.Mode); err != nil {
		return fmt.Errorf("chmod %q with %o, %s", fi.Name, fi.Mode, err)
	}
	return Chall(fi)
}

// Mkdev is used to build the value of linux devices (in /dev/) which specifies major
// and minor number of the newly created device special file.
// Linux device nodes are a bit weird due to backwards compat with 16 bit device nodes.
// They are, from low to high: the lower 8 bits of the minor, then 12 bits of the major,
// then the top 12 bits of the minor.
func mkdev(major uint32, minor uint32) uint32 {
	return uint32(unix.Mkdev(major, minor))
}
