package fs

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type TarReader struct {
	tr *tar.Reader
}

func NewTarReader(tr *tar.Reader) Reader {
	return &TarReader{
		tr: tr,
	}
}

func (r *TarReader) Next() (*File, error) {
	tr := r.tr

	th, err := tr.Next()
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	// TODO this will take up huge memory when the file is large.
	content, err := ioutil.ReadAll(tr)
	if err != nil {
		return nil, fmt.Errorf("reading %q from tar: %w", th.Name, err)
	}

	mode := os.FileMode(th.Mode)

	switch th.Typeflag {
	//case tar.TypeReg:
	//case tar.TypeLink:
	case tar.TypeSymlink:
		mode |= os.ModeSymlink
	case tar.TypeChar:
		mode |= os.ModeCharDevice
		mode |= os.ModeDevice
	case tar.TypeBlock:
		mode |= os.ModeDevice
	case tar.TypeDir:
		mode |= os.ModeDir
	case tar.TypeFifo:
		mode |= os.ModeNamedPipe
	}

	f := &File{
		FileInfo: FileInfo{
			Name:     th.Name,
			Size:     uint64(th.Size),
			Uid:      uint32(th.Uid),
			Gid:      uint32(th.Gid),
			Mtime:    th.ModTime,
			Atime:    th.AccessTime,
			Ctime:    th.ChangeTime,
			Mode:     mode,
			Xattrs:   th.PAXRecords,
			Linkname: th.Linkname,
			Major:    uint32(th.Devmajor),
			Minor:    uint32(th.Devminor),
		},
		Data: ioutil.NopCloser(bytes.NewReader(content)),
	}
	return f, nil
}
