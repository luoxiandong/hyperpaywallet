package address_decode

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"errors"
	"fmt"
)

// DigestSize is the number of bytes in the preferred hash Digest used here.
const DigestSize = sha512.Size256
const PublicKeySize = 32
const ChecksumLength = 4

// Digest represents a 32-byte value holding the 256-bit Hash digest.
type Digest [DigestSize]byte
type PublicKey [PublicKeySize]byte

type (
	Address Digest
)

type ChecksumAddress struct {
	shortAddress Address
	checksum     []byte
}

// AddressDecoderV2
type AddressDecoderV2 struct {
}

var (
	Default = AddressDecoderV2{}
)

var (
	errWrongAddressLen     = fmt.Errorf("decoded address is the wrong length, should be %d bytes", DigestSize+ChecksumLength)
	errWrongChecksum       = fmt.Errorf("address checksum is incorrect, did you copy the address correctly?")
	ErrorInvalidHashLength = errors.New("Invalid hash length!")
	ErrorInvalidAddress    = errors.New("Invalid address!")
)

//AddressEncode encode address bytes
func (dec *AddressDecoderV2) AddressEncode(address []byte) (string, error) {
	var pk PublicKey

	if len(pk) != len(address) {
		return "", ErrorInvalidHashLength
	}

	for i := range pk {
		pk[i] = address[i]
	}

	publicKeyChecksummed := Address(pk).GetChecksumAddress().String()
	return publicKeyChecksummed, nil
}

// DecodeAddress turns a checksum address string into an Address object. It
// checks that the checksum is correct, and returns an error if it's not.
func DecodeAddress(addr string) (a Address, err error) {
	// Interpret the address as base32
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(addr)
	if err != nil {
		return
	}

	// Ensure the decoded address is the correct length
	if len(decoded) != len(a)+ChecksumLength {
		err = errWrongAddressLen
		return
	}

	// Split into address + checksum
	addressBytes := decoded[:len(a)]
	checksumBytes := decoded[len(a):]

	// Compute the expected checksum
	checksumHash := sha512.Sum512_256(addressBytes)
	expectedChecksumBytes := checksumHash[DigestSize-ChecksumLength:]

	// Check the checksum
	if !bytes.Equal(expectedChecksumBytes, checksumBytes) {
		err = errWrongChecksum
		return
	}

	// Checksum is good, copy address bytes into output
	copy(a[:], addressBytes)
	return a, nil
}

// GetChecksumAddress returns the short address with its checksum as a string
// Checksum in Algorand are the last 4 bytes of the shortAddress Hash. H(Address)[28:]
func (addr Address) GetChecksumAddress() *ChecksumAddress {
	shortAddressHash := Hash(addr[:])
	return &ChecksumAddress{addr, shortAddressHash[len(shortAddressHash)-ChecksumLength:]}
}

// String returns a string representation of ChecksumAddress
func (addr *ChecksumAddress) String() string {
	var addrWithChecksum []byte
	addrWithChecksum = append(addr.shortAddress[:], addr.checksum...)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(addrWithChecksum)
}

// Hash computes the SHASum512_256 hash of an array of bytes
func Hash(data []byte) Digest {
	return sha512.Sum512_256(data)
}
