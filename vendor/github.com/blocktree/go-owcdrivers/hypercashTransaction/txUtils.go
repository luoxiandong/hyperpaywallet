package hypercashTransaction

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
)

//reverseBytes endian reverse
func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//reverseHexToBytes decode a hex string to an byte array,then change the endian
func reverseHexToBytes(hexVar string) ([]byte, error) {
	if len(hexVar)%2 == 1 {
		return nil, errors.New("Invalid TxHash!")
	}
	ret, err := hex.DecodeString(hexVar)
	if err != nil {
		return nil, err
	}
	return reverseBytes(ret), nil
}

//reverseBytesToHex change the endian of the input byte array then encode it to hex string
func reverseBytesToHex(bytesVar []byte) string {
	return hex.EncodeToString(reverseBytes(bytesVar))
}

func uint16ToBigEndianBytes(data uint16) []byte {
	tmp := [2]byte{}
	binary.BigEndian.PutUint16(tmp[:], data)
	return tmp[:]
}

//uint16ToLittleEndianBytes
func uint16ToLittleEndianBytes(data uint16) []byte {
	tmp := [2]byte{}
	binary.LittleEndian.PutUint16(tmp[:], data)
	return tmp[:]
}

//littleEndianBytesToUint16
func littleEndianBytesToUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

func uint32ToBigEndianBytes(data uint32) []byte {
	tmp := [4]byte{}
	binary.BigEndian.PutUint32(tmp[:], data)
	return tmp[:]
}

//uint32ToLittleEndianBytes
func uint32ToLittleEndianBytes(data uint32) []byte {
	tmp := [4]byte{}
	binary.LittleEndian.PutUint32(tmp[:], data)
	return tmp[:]
}

//littleEndianBytesToUint32
func littleEndianBytesToUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func uint64ToBigEndianBytes(data uint64) []byte {
	tmp := [8]byte{}
	binary.BigEndian.PutUint64(tmp[:], data)
	return tmp[:]
}

//uint64ToLittleEndianBytes
func uint64ToLittleEndianBytes(data uint64) []byte {
	tmp := [8]byte{}
	binary.LittleEndian.PutUint64(tmp[:], data)
	return tmp[:]
}

//littleEndianBytesToUint64
func littleEndianBytesToUint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func varIntToBytes(val uint64) []byte {
	if val < 0xfd {
		return []byte{byte(val)}
	}

	if val <= math.MaxUint16 {
		return append([]byte{0xfd}, uint16ToLittleEndianBytes(uint16(val))...)
	}

	if val <= math.MaxUint32 {
		return append([]byte{0xfe}, uint32ToLittleEndianBytes(uint32(val))...)
	}

	return append([]byte{0xff}, uint64ToLittleEndianBytes(val)...)
}