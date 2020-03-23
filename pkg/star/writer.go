package star

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/sequix/star/pkg/encoding"
	"github.com/sequix/star/pkg/fs"
)

const (
	Magic    = 0xab72617473cd
	Version1 = 0x00
)

// <magic>(8) <version>(1) <payload-length>(8) <info-length>(4)
// <payload1> <payload2> ... <payloadN>
// <Info1> <Info2> .... <InfoN>
func WriteTo(w io.WriterAt, fsr fs.Reader) error {
	var (
		infos   []*Info
		infoBuf = make([]byte, 0, 128)
		offset  = uint64(8 + 1 + 8 + 8)
		wto     = &writerToOffset{w: w}
	)

	n, err := w.WriteAt(encoding.PutUint64(nil, Magic), 0)
	if err != nil {
		return fmt.Errorf("writing star magic, written %d bytes, err %w", n, err)
	}
	n, err = w.WriteAt([]byte{Version1}, 8)
	if err != nil {
		return fmt.Errorf("writing star vetsion, written %d bytes, err %w", n, err)
	}

	for {
		f, err := fsr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("reading next file, %w", err)
		}

		info := &Info{
			FileInfo: &f.FileInfo,
			Offset:   offset,
		}
		infos = append(infos, info)

		if !info.Mode.IsRegular() {
			continue
		}

		wto.offset = int64(offset)
		if n, err := io.Copy(wto, f.Data); err != nil {
			return fmt.Errorf("copying file %q to offset %d, written %d, err %w", f.Name, offset, n, err)
		}
		offset += f.Size
	}

	payloadLength := offset - 21
	n, err = w.WriteAt(encoding.PutUint64(nil, payloadLength), 9)
	if err != nil {
		return fmt.Errorf("writing payload length %d, written %d, err %w", payloadLength, n, err)
	}

	offsetBeforeInfo := offset
	for _, info := range infos {
		infoBuf = marshalInfoTo(infoBuf[:0], info)
		n, err = w.WriteAt(infoBuf, int64(offset))
		if err != nil {
			return fmt.Errorf("writing info for %q, written %d, err %w", info.Name, n, err)
		}
		offset += uint64(len(infoBuf))
	}

	infoLength := offset - offsetBeforeInfo
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})
	n, err = w.WriteAt(encoding.PutUint32(nil, uint32(infoLength)), 17)
	if err != nil {
		return fmt.Errorf("writing info length %d, written %d, err %w", infoLength, n, err)
	}
	return nil
}

type writerToOffset struct {
	w      io.WriterAt
	offset int64
}

func (w *writerToOffset) Write(p []byte) (n int, err error) {
	return w.w.WriteAt(p, w.offset)
}
