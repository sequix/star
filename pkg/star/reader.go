package star

import (
	"fmt"
	"io"

	"github.com/sequix/star/pkg/encoding"
)

type Reader struct {
	r          io.ReaderAt
	payloadLen uint64
	infoLen    uint32
	infos      []*Info
	name2Info  map[string]*Info
}

func NewReader(r io.ReaderAt) (*Reader, error) {
	var (
		ifo *Info
		src = make([]byte, 0, 512)
		sr  = &Reader{
			r:         r,
			name2Info: map[string]*Info{},
		}
	)

	n, err := r.ReadAt(src[:8], 0)
	if err != nil {
		return nil, fmt.Errorf("reading star magic, read %d bytes, err %w", n, err)
	}
	_, magic, err := encoding.GetUint64(src[:8])
	if err != nil || magic != Magic {
		return nil, fmt.Errorf("parsing star magic, want %x, got %x, err %w", Magic, magic, err)
	}

	// If need upgrade the star format version, read version flag here.

	n, err = r.ReadAt(src[:8], 9)
	if err != nil {
		return nil, fmt.Errorf("reading payload length, read %d bytes, err %w", n, err)
	}
	_, sr.payloadLen, err = encoding.GetUint64(src[:8])
	if err != nil {
		return nil, fmt.Errorf("parsing payload length, %w", err)
	}

	n, err = r.ReadAt(src[:4], 17)
	if err != nil {
		return nil, fmt.Errorf("reading info length, read %d bytes, err %w", n, err)
	}
	_, sr.infoLen, err = encoding.GetUint32(src[:4])
	if err != nil {
		return nil, fmt.Errorf("parsing info length, %w", err)
	}

	src = encoding.Resize(src, int(sr.infoLen))
	n, err = r.ReadAt(src, int64(21+sr.payloadLen))
	if err != nil {
		return nil, fmt.Errorf("reading infos, reda %d bytes, err %w", n, err)
	}

	for len(src) > 0 {
		src, ifo, err = unmarshalInfoFrom(src)
		if err != nil {
			return nil, fmt.Errorf("parsing info, %w", err)
		}
		sr.infos = append(sr.infos, ifo)
		sr.name2Info[ifo.Name] = ifo
	}
	return sr, nil
}

func (r *Reader) ListFiles() []*Info {
	return r.infos
}

func (r *Reader) ListNames() []string {
	names := make([]string, 0, len(r.infos))
	for _, ifo := range r.infos {
		names = append(names, ifo.Name)
	}
	return names
}

func (r *Reader) ReaderAtFor(name string) (io.ReaderAt, error) {
	fi, ok := r.name2Info[name]
	if !ok || fi == nil {
		return nil, fmt.Errorf("not found info with name %q", name)
	}
	fr := &fileReaderAt{
		r:     r.r,
		start: int64(fi.Offset),
		size:  int64(fi.Size),
	}
	return fr, nil
}

func (r *Reader) ReaderFor(name string) (io.Reader, error) {
	fi, ok := r.name2Info[name]
	if !ok || fi == nil {
		return nil, fmt.Errorf("not found info with name %q", name)
	}
	fr := &fileReader{
		r:      r.r,
		offset: int64(fi.Offset),
		bound:  int64(fi.Offset + fi.Size),
	}
	return fr, nil
}

func (r *Reader) Mount(mountpoint string) error {
	panic("todo")
}

type fileReaderAt struct {
	r     io.ReaderAt
	start int64
	size  int64
}

func (r *fileReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off > (r.size-r.start+1) {
		return 0, fmt.Errorf("ReaderAt want off within [0, %d], got %d", r.size-r.start+1, off)
	}
	return r.r.ReadAt(p, off+r.start)
}

type fileReader struct {
	r      io.ReaderAt
	offset int64
	bound  int64
}

func (f *fileReader) Read(p []byte) (n int, err error) {
	if f.offset >= f.bound {
		return 0, io.EOF
	}
	if f.offset+int64(len(p)) > f.bound {
		p = p[:f.bound-f.offset]
	}
	n, err = f.r.ReadAt(p, f.offset)
	if err != nil {
		return n, fmt.Errorf("Reader read %d bytes, err %w", n, err)
	}
	f.offset += int64(n)
	return n, nil
}
