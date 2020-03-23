package star

import (
	"fmt"
	"os"
	"time"

	"github.com/sequix/star/pkg/encoding"
	"github.com/sequix/star/pkg/fs"
)

type Info struct {
	*fs.FileInfo
	Offset uint64
}

func marshalInfoTo(dst []byte, f *Info) []byte {
	dst = encoding.PutStr(dst, f.Name)
	dst = encoding.PutStr(dst, f.Linkname)
	dst = encoding.PutUint64(dst, f.Offset)
	dst = encoding.PutUint64(dst, f.Size)
	dst = encoding.PutUint32(dst, f.Uid)
	dst = encoding.PutUint32(dst, f.Gid)
	dst = encoding.PutUint64(dst, uint64(f.Mtime.UnixNano()))
	dst = encoding.PutUint64(dst, uint64(f.Atime.UnixNano()))
	dst = encoding.PutUint64(dst, uint64(f.Ctime.UnixNano()))
	dst = encoding.PutUint32(dst, uint32(f.Mode))
	dst = encoding.PutUint32(dst, f.Major)
	dst = encoding.PutUint32(dst, f.Minor)
	return dst
}

func unmarshalInfoFrom(src []byte) ([]byte, *Info, error) {
	var (
		u64 uint64
		u32 uint32
		err error
		f = &Info{
			FileInfo: &fs.FileInfo{},
		}
	)
	src, f.Name, err = encoding.GetStr(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting name, %w", err)
	}
	
	src, f.Linkname, err = encoding.GetStr(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting linkname, %w", err)
	}
	
	src, f.Offset, err = encoding.GetUint64(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting offset, %w", err)
	}

	src, f.Size, err = encoding.GetUint64(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting size, %w", err)
	}
	
	src, f.Uid, err = encoding.GetUint32(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting uid, %w", err)
	}
	
	src, f.Gid, err = encoding.GetUint32(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting gid, %w", err)
	}
	
	src, u64, err = encoding.GetUint64(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting mtime, %w", err)
	}
	f.Mtime = time.Unix(0, int64(u64))

	src, u64, err = encoding.GetUint64(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting atime, %w", err)
	}
	f.Ctime = time.Unix(0, int64(u64))

	src, u64, err = encoding.GetUint64(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting ctime, %w", err)
	}
	f.Atime = time.Unix(0, int64(u64))
	
	src, u32, err = encoding.GetUint32(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting mode, %w", err)
	}
	f.Mode = os.FileMode(u32)
	
	src, f.Major, err = encoding.GetUint32(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting major, %w", err)
	}

	src, f.Minor, err = encoding.GetUint32(src)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshalInfoFrom getting minor, %w", err)
	}
	return src, f, nil
}