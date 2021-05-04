package wavesTransaction

import (
	"bytes"

	"github.com/pkg/errors"
)

type B58Bytes []byte

const (
	headerSize     = 2
	bodySize       = 20
	checksumSize   = 4
	AddressSize    = headerSize + bodySize + checksumSize
	aliasFixedSize = 4

	addressVersion byte = 0x01
	aliasVersion   byte = 0x02

	AliasMinLength = 4
	AliasMaxLength = 30
	AliasAlphabet  = "-.0123456789@_abcdefghijklmnopqrstuvwxyz"
	AliasPrefix    = "alias"

	MainNetScheme byte = 'W'
	TestNetScheme byte = 'T'
	DevNetScheme  byte = 'D'
)

// Address is the transformed Public Key with additional bytes of the version, a blockchain scheme and a checksum.
type Address [AddressSize]byte

// String produces the BASE58 string representation of the Address.
func (a Address) String() string {
	return Encode(a[:])
}

// NewAddressFromPublicKey produces an Address from given scheme and Public Key bytes.
func NewAddressFromPublicKey(scheme byte, publicKey PublicKey) (Address, error) {
	var a Address
	a[0] = addressVersion
	a[1] = scheme
	h, err := SecureHash(publicKey[:])
	if err != nil {
		return a, errors.Wrap(err, "failed to produce Digest from PublicKey")
	}
	copy(a[headerSize:], h[:bodySize])
	cs, err := addressChecksum(a[:headerSize+bodySize])
	if err != nil {
		return a, errors.Wrap(err, "failed to calculate Address checksum")
	}
	copy(a[headerSize+bodySize:], cs)
	return a, nil
}

func MustAddressFromPublicKey(scheme byte, publicKey PublicKey) Address {
	rs, err := NewAddressFromPublicKey(scheme, publicKey)
	if err != nil {
		panic(err)
	}
	return rs
}

func RebuildAddress(scheme byte, body []byte) (Address, error) {
	var a Address
	a[0] = addressVersion
	a[1] = scheme
	if l := len(body); l != bodySize {
		return Address{}, errors.Errorf("%d is unexpected address' body size", l)
	}
	copy(a[headerSize:], body[:bodySize])
	cs, err := addressChecksum(a[:headerSize+bodySize])
	if err != nil {
		return a, errors.Wrap(err, "failed to calculate Address checksum")
	}
	copy(a[headerSize+bodySize:], cs)
	return a, nil
}

// NewAddressFromString creates an Address from its string representation. This function checks that the address is valid.
func NewAddressFromString(s string) (Address, error) {
	var a Address
	b, err := Decode(s)
	if err != nil {
		return a, errors.Wrap(err, "invalid Base58 string")
	}
	a, err = NewAddressFromBytes(b)
	if err != nil {
		return a, errors.Wrap(err, "failed to create an Address from Base58 string")
	}
	return a, nil
}

// NewAddressFromBytes creates an Address from the slice of bytes and checks that the result address is valid address.
func NewAddressFromBytes(b []byte) (Address, error) {
	var a Address
	if l := len(b); l < AddressSize {
		return a, errors.Errorf("insufficient array length %d, expected at least %d", l, AddressSize)
	}
	copy(a[:], b[:AddressSize])
	if ok, err := a.Valid(); !ok {
		return a, errors.Wrap(err, "invalid address")
	}
	return a, nil
}

// Valid checks that version and checksum of the Address are correct.
func (a *Address) Valid() (bool, error) {
	if a[0] != addressVersion {
		return false, errors.Errorf("unsupported address version %d", a[0])
	}
	hb := a[:headerSize+bodySize]
	ec, err := addressChecksum(hb)
	if err != nil {
		return false, errors.Wrap(err, "failed to calculate Address checksum")
	}
	ac := a[headerSize+bodySize:]
	if !bytes.Equal(ec, ac) {
		return false, errors.New("invalid Address checksum")
	}
	return true, nil
}

// Bytes converts the fixed-length byte array of the Address to a slice of bytes.
func (a Address) Bytes() []byte {
	return a[:]
}

func addressChecksum(b []byte) ([]byte, error) {
	h, err := SecureHash(b)
	if err != nil {
		return nil, err
	}
	c := make([]byte, checksumSize)
	copy(c, h[:checksumSize])
	return c, nil
}

func correctAlphabet(s string) bool {
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'z') && c != '_' && c != '@' && c != '-' && c != '.' {
			return false
		}
	}
	return true
}

type Recipient struct {
	Address *Address
	len     int
}

func (r *Recipient) MarshalBinary() ([]byte, error) {
	return r.Address[:], nil
}
func NewRecipientFromAddress(a Address) Recipient {
	return Recipient{Address: &a, len: AddressSize}
}
