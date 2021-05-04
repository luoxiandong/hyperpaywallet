package types

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/liyaojian/hxwallet/util"
)

type PublicKeys []PublicKey

type PublicKey struct {
	key      *btcec.PublicKey
	prefix   string
	checksum []byte
}

func (p PublicKey) String() string {
	b := append(p.Bytes(), p.checksum...)
	return fmt.Sprintf("%s%s", p.prefix, base58.Encode(b))
}

func (p *PublicKey) ToAddress() (*Address, error) {
	return NewAddress(p)
}

func (p PublicKey) Bytes() []byte {
	return p.key.SerializeCompressed()
}

func (p PublicKey) Equal(pub *PublicKey) bool {
	return p.key.IsEqual(pub.key)
}

func (p PublicKey) ToECDSA() *ecdsa.PublicKey {
	return p.key.ToECDSA()
}

// MaxSharedKeyLength returns the maximum length of the shared key the
// public key can produce.
func (p PublicKey) MaxSharedKeyLength() int {
	return (p.key.ToECDSA().Curve.Params().BitSize + 7) / 8
}

func (p PublicKey) IsNul() bool {
	return p.key == nil
}

//NewPublicKey creates a new PublicKey from string
func NewPublicKeyFromString(key string) (*PublicKey, error) {
	prefixChain := addrPrefix

	prefix := key[:len(prefixChain)]

	if prefix != prefixChain {
		return nil, ErrPublicKeyChainPrefixMismatch
	}

	b58 := base58.Decode(key[len(prefixChain):])
	if len(b58) < 5 {
		return nil, ErrInvalidPublicKey
	}

	chk1 := b58[len(b58)-4:]

	keyBytes := b58[:len(b58)-4]
	chk2, err := util.Ripemd160Checksum(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("Ripemd160Checksum, %+v\n", err)
	}

	if !bytes.Equal(chk1, chk2) {
		return nil, ErrInvalidPublicKey
	}

	pub, err := btcec.ParsePubKey(keyBytes, btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("ParsePubKey, %+v\n", err)
	}

	k := PublicKey{
		key:      pub,
		prefix:   prefix,
		checksum: chk1,
	}

	return &k, nil
}

func NewPublicKey(pub *btcec.PublicKey) (*PublicKey, error) {
	buf := pub.SerializeCompressed()
	chk, err := util.Ripemd160Checksum(buf)
	if err != nil {
		return nil, fmt.Errorf("Ripemd160Checksum, %+v\n", err)
	}

	k := PublicKey{
		key:      pub,
		prefix:   addrPrefix,
		checksum: chk,
	}

	return &k, nil
}
