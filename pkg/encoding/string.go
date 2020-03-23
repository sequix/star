package encoding

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PutStr(dst []byte, s string) []byte {
	dst = PutUint32(dst, uint32(len(s)))
	dst = append(dst, ToUnsafeBytes(s)...)
	return dst
}

func GetStr(src []byte) ([]byte, string, error) {
	var (
		n uint32
		err error
	)
	src, n, err = GetUint32(src)
	if err != nil {
		return nil, "", fmt.Errorf("GetStr getting length, %w", err)
	}
	if uint32(len(src)) < n {
		return nil, "", fmt.Errorf("GetStr need %d bytes, got %d byte", n, len(src))
	}
	return src[n:], string(src[:n]), nil
}

// ToUnsafeString converts b to string without memory allocations.
//
// The returned string is valid only until b is reachable and unmodified.
func ToUnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToUnsafeBytes converts s to a byte slice without memory allocations.
//
// The returned byte slice is valid only until s is reachable and unmodified.
func ToUnsafeBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var slh reflect.SliceHeader
	slh.Data = sh.Data
	slh.Len = sh.Len
	slh.Cap = sh.Len
	return *(*[]byte)(unsafe.Pointer(&slh))
}

// Resize resizes b to n bytes and returns b (which may be newly allocated).
func Resize(b []byte, n int) []byte {
	if nn := n - cap(b); nn > 0 {
		b = append(b[:cap(b)], make([]byte, nn)...)
	}
	return b[:n]
}