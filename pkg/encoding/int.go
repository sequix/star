package encoding

import (
	"encoding/binary"
	"fmt"
)

var (
	bigEndian = binary.BigEndian
)

func PutUint16(dst []byte, d uint16) []byte {
	ld, cd := len(dst), cap(dst)
	if nn := 2 - (cd - ld); nn > 0 {
		dst = append(dst[:cd], make([]byte, nn)...)
	}
	bigEndian.PutUint16(dst[ld:ld+2], d)
	return dst[:ld+2]
}

func PutUint32(dst []byte, d uint32) []byte {
	ld, cd := len(dst), cap(dst)
	if nn := 4 - (cd - ld); nn > 0 {
		dst = append(dst[:cd], make([]byte, nn)...)
	}
	bigEndian.PutUint32(dst[ld:ld+4], d)
	return dst[:ld+4]
}

func PutUint64(dst []byte, d uint64) []byte {
	ld, cd := len(dst), cap(dst)
	if nn := 8 - (cd - ld); nn > 0 {
		dst = append(dst[:cd], make([]byte, nn)...)
	}
	bigEndian.PutUint64(dst[ld:ld+8], d)
	return dst[:ld+8]
}

func GetUint16(src []byte) ([]byte, uint16, error) {
	if len(src) < 2 {
		return nil, 0, fmt.Errorf("GetUint16 need 2 bytes at least, got %d bytes", len(src))
	}
	return src[2:], bigEndian.Uint16(src), nil
}

func GetUint32(src []byte) ([]byte, uint32, error) {
	if len(src) < 4 {
		return nil, 0, fmt.Errorf("GetUint32 need 4 bytes at least, got %d bytes", len(src))
	}
	return src[4:], bigEndian.Uint32(src), nil
}

func GetUint64(src []byte) ([]byte, uint64, error) {
	if len(src) < 8 {
		return nil, 0, fmt.Errorf("GetUint64 need 8 bytes at least, got %d bytes", len(src))
	}
	return src[8:], bigEndian.Uint64(src), nil
}