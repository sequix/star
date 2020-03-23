package fs

import (
	"io"
	"os"
	"time"
)

// FilesystemReader is an interface for source filesystem to be used during
// tar operations. Next() is expected to return files and directories in a
// consistent and stable order and return io.EOF when no further files are available.
type Reader interface {
	Next() (*File, error)
}

// File represents a filesystem object such as directoy, file, symlink or device.
// It's used when creating archives from a source filesystem which can be a real
// OS filesystem, or another archive stream such as tar.
type File struct {
	FileInfo

	// File content. Nil for non-regular files.
	Data io.ReadCloser
}

type FileInfo struct {
	Name string
	Size uint64

	Uid uint32
	Gid uint32

	// Modification & Access time
	Mtime time.Time
	Atime time.Time
	Ctime time.Time

	// Mode & Extended Mode
	Mode os.FileMode
	Xattrs map[string]string

	// Link target for symlinks
	Linkname string

	// Major/Minor for character or block devices
	Major uint32
	Minor uint32
}

type OptFunc func(ifo *FileInfo)

func WithName(name string) OptFunc {
	return func(ifo *FileInfo) {
		ifo.Name = name
	}
}

func WithSize(size uint64) OptFunc {
	return func(ifo *FileInfo) {
		ifo.Size = size
	}
}

func WithFileInfo(nifo *FileInfo) OptFunc {
	return func(ifo *FileInfo) {
		ifo.Name = nifo.Name
		ifo.Size = nifo.Size
		ifo.Uid = nifo.Uid
		ifo.Gid = nifo.Gid
		ifo.Mtime = nifo.Mtime
		ifo.Atime = nifo.Atime
		ifo.Ctime = nifo.Ctime
		ifo.Mode = nifo.Mode

	}
}