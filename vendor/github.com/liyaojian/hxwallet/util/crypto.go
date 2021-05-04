package util

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

func Ripemd160(in []byte) ([]byte, error) {
	h := ripemd160.New()

	if _, err := h.Write(in); err != nil {
		return nil, fmt.Errorf("Write, %+v\n", err)
	}

	sum := h.Sum(nil)
	return sum, nil
}

func Ripemd160Checksum(in []byte) ([]byte, error) {
	buf, err := Ripemd160(in)
	if err != nil {
		return nil, fmt.Errorf("Ripemd160, %+v\n", err)
	}

	return buf[:4], nil
}
func Sha512Checksum(in []byte) ([]byte, error) {
	buf := sha512.Sum512(in)
	return buf[:4], nil
}

func Sha256(in []byte) []byte {
	buf := sha256.Sum256(in)
	return buf[:]
}

func CharToSymbol(c byte) uint64 {
	if c >= 'a' && c <= 'z' {
		return uint64((c - 'a') + 6)
	}
	if c >= '1' && c <= '5' {
		return uint64((c - '1') + 1)
	}
	return 0
}
