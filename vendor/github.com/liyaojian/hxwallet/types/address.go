package types

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"github.com/liyaojian/hxwallet/util"

	"github.com/btcsuite/btcutil/base58"
)

type Address struct {
	prefix   string
	data     []byte
	checksum []byte
}

func (p Address) String() string {
	b := append(p.data, p.checksum...)
	return fmt.Sprintf("%s%s", p.prefix, base58.Encode(b))
}

func (p Address) Bytes() []byte {
	return p.data
}

func NewAddress(pub *PublicKey) (*Address, error) {
	buf512 := sha512.Sum512(pub.Bytes())
	data, err := util.Ripemd160(buf512[:])
	if err != nil {
		return nil, fmt.Errorf("Ripemd160, %+v\n", err)
	}
	binAddy := append([]byte{walletVersion}, data...)

	chk1, err := util.Ripemd160Checksum(binAddy)
	if err != nil {
		return nil, fmt.Errorf("Ripemd160Checksum, %+v\n", err)
	}

	ad := Address{
		prefix:   addrPrefix,
		data:     binAddy,
		checksum: chk1,
	}

	return &ad, nil
}

//NewAddress creates a new Address from string
func NewAddressFromString(add string) (*Address, error) {
	prefixChain := addrPrefix

	prefix := add[:len(prefixChain)]

	if prefix != prefixChain {
		return nil, ErrAddressChainPrefixMismatch
	}

	b58 := base58.Decode(add[len(prefixChain):])
	if len(b58) < 5 {
		return nil, ErrInvalidAddress
	}

	chk1 := b58[len(b58)-4:]
	data := b58[:len(b58)-4]

	chk2, err := util.Ripemd160Checksum(data)
	if err != nil {
		return nil, fmt.Errorf("Ripemd160Checksum, %+v\n", err)
	}

	if !bytes.Equal(chk1, chk2) {
		return nil, ErrInvalidAddress
	}

	a := Address{
		data:     data,
		prefix:   prefix,
		checksum: chk1,
	}

	return &a, nil
}
