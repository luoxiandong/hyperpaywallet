package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

var (
	ErrInvalidCurve               = fmt.Errorf("invalid elliptic curve")
	ErrSharedKeyTooBig            = fmt.Errorf("shared key params are too big")
	ErrSharedKeyIsPointAtInfinity = fmt.Errorf("shared key is point at infinity")
)

type PrivateKeys []PrivateKey

type PrivateKey struct {
	priv *btcec.PrivateKey
	pub  *PublicKey
	raw  []byte
}

func NewPrivateKeyFromBrainKey(brainKey string) (*PrivateKey, error) {
	hashByte := sha256.Sum256([]byte(brainKey))
	privateKey, err := NewDeterministicPrivateKey(bytes.NewBuffer(hashByte[:]))
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func NewDeterministicPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	return newRandomPrivateKey(randSource)
}

func newRandomPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	rawPrivKey := make([]byte, 32)
	written, err := io.ReadFull(randSource, rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("error feeding crypto-rand numbers to seed ephemeral private key: %s", err)
	}
	if written != 32 {
		return nil, fmt.Errorf("couldn't write 32 bytes of randomness to seed ephemeral private key")
	}

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), rawPrivKey)

	pub, err := NewPublicKey(privKey.PubKey())
	if err != nil {
		return nil, fmt.Errorf("NewPublicKey, %+v\n", err)
	}

	raw := append([]byte{128}, privKey.D.Bytes()...)
	raw = append(raw, checksum(raw)...)

	return &PrivateKey{
		priv: privKey,
		pub:  pub,
		raw:  raw,
	}, nil
}

func checksum(data []byte) []byte {
	c1 := sha256.Sum256(data)
	c2 := sha256.Sum256(c1[:])
	return c2[0:4]
}

func NewPrivateKeyFromWif(wifPrivateKey string) (*PrivateKey, error) {
	w, err := btcutil.DecodeWIF(wifPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("DecodeWIF, %+v\n", err)
	}

	priv := w.PrivKey
	raw := base58.Decode(wifPrivateKey)
	pub, err := NewPublicKey(priv.PubKey())
	if err != nil {
		return nil, fmt.Errorf("NewPublicKey, %+v\n", err)
	}

	k := PrivateKey{
		priv: priv,
		raw:  raw,
		pub:  pub,
	}

	return &k, nil
}

func (p PrivateKey) PublicKey() *PublicKey {
	return p.pub
}

func (p PrivateKey) ECPrivateKey() *btcec.PrivateKey {
	return p.priv
}

func (p PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return p.priv.ToECDSA()
}

func (p PrivateKey) Bytes() []byte {
	return p.priv.Serialize()
}

func (p PrivateKey) ToHex() string {
	return hex.EncodeToString(p.Bytes())
}

func (p PrivateKey) ToWIF() string {
	return base58.Encode(p.raw)
}

func (p PrivateKey) SignCompact(hash []byte) (sig []byte, err error) {
	sig, err = btcec.SignCompact(btcec.S256(), p.ECPrivateKey(), hash, true)
	return
}

func (p PrivateKey) SharedSecret(pub *PublicKey, skLen, macLen int) (sk []byte, err error) {
	puk := pub.ToECDSA()
	pvk := p.priv

	if pvk.PublicKey.Curve != puk.Curve {
		return nil, ErrInvalidCurve
	}

	if skLen+macLen > pub.MaxSharedKeyLength() {
		return nil, ErrSharedKeyTooBig
	}

	x, _ := puk.Curve.ScalarMult(puk.X, puk.Y, pvk.D.Bytes())
	if x == nil {
		return nil, ErrSharedKeyIsPointAtInfinity
	}

	sk = make([]byte, skLen+macLen)
	skBytes := x.Bytes()
	copy(sk[len(sk)-len(skBytes):], skBytes)
	return sk, nil
}
