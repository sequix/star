// +build linux

package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type LocalReader struct {
	fis []os.FileInfo
}

func NewLocalReader(files ...string) (Reader, error) {
	var (
		dedup = map[string]struct{}{}
		lr    = &LocalReader{
			fis: make([]os.FileInfo, 0, len(files)),
		}
	)

	for _, f := range files {
		dedup[f] = struct{}{}
	}

	for f := range dedup {
		fi, err := os.Lstat(f)
		if err != nil {
			return nil, fmt.Errorf("lstat file %q, %w", f, err)
		}
		lr.fis = append(lr.fis, fi)
	}
	return lr, nil
}

func (r *LocalReader) Next() (*File, error) {
	if len(r.fis) == 0 {
		return nil, io.EOF
	}
	var (
		err  error
		data io.ReadCloser
		fi   = r.fis[0]
		name = fi.Name()
	)
	r.fis = r.fis[1:]

	mode := fi.Mode()
	size := uint64(0)
	if mode.IsRegular() {
		data, err = os.Open(name)
		if err != nil {
			return nil, fmt.Errorf("opening file %q, %w", name, err)
		}
		size = uint64(fi.Size())
	} else if mode.IsDir() {
		nfis, err := ioutil.ReadDir(name)
		if err != nil {
			return nil, fmt.Errorf("reading dir %q, %w", name, err)
		}
		for _, nfi := range nfis {
			r.fis = append(r.fis, &fullPathFileInfo{
				FileInfo: nfi,
				prefix:   fi.Name(),
			})
		}
	}

	var linkname string
	if (fi.Mode() & os.ModeType) == os.ModeSymlink {
		linkname, err = os.Readlink(name)
		if err != nil {
			return nil, fmt.Errorf("reading link %q")
		}
	}

	fsi, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || fsi == nil {
		return nil, fmt.Errorf("getting Stat_t")
	}

	fh := &File{
		FileInfo: FileInfo{
			Name:     name,
			Size:     size,
			Uid:      fsi.Uid,
			Gid:      fsi.Gid,
			Mtime:    time.Unix(fsi.Mtim.Sec, fsi.Mtim.Nsec),
			Atime:    time.Unix(fsi.Atim.Sec, fsi.Atim.Nsec),
			Ctime:    time.Unix(fsi.Ctim.Sec, fsi.Ctim.Nsec),
			Mode:     fi.Mode(),
			Xattrs:   map[string]string{},
			Linkname: linkname,
			Major:    uint32(fsi.Rdev) >> 8,
			Minor:    uint32(fsi.Rdev & 0xFF),
		},
		Data: data,
	}
	return fh, nil
}

type fullPathFileInfo struct {
	os.FileInfo
	prefix string
}

func (fi *fullPathFileInfo) Name() string {
	return filepath.Join(fi.prefix, fi.FileInfo.Name())
}