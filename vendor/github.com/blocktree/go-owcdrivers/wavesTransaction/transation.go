package wavesTransaction

import (
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/blocktree/go-owcrypt"
	"github.com/pkg/errors"
)

const (
	proofsVersion  byte = 1
	proofsMinLen        = 1 + 2
	proofsMaxCount      = 8
	proofMaxSize        = 64

	transferLen            = PublicKeySize + 1 + 1 + 8 + 8 + 8 + 2
	transferV2FixedBodyLen = 1 + 1 + transferLen
	transferV2MinLen       = 1 + transferV2FixedBodyLen + proofsMinLen

	WavesAssetName = "WAVES"
)

type Attachment string

func (a Attachment) String() string {
	return string(a)
}

type TransferV2 struct {
	Version byte      `json:"version,omitempty"`
	ID      *Digest   `json:"id,omitempty"`
	Proofs  *ProofsV1 `json:"proofs,omitempty"`
	Transfer
}

type Transfer struct {
	SenderPK    PublicKey     `json:"senderPublicKey"`
	AmountAsset OptionalAsset `json:"assetId"`
	FeeAsset    OptionalAsset `json:"feeAssetId"`
	Timestamp   uint64        `json:"timestamp,omitempty"`
	Amount      uint64        `json:"amount"`
	Fee         uint64        `json:"fee"`
	Recipient   Recipient     `json:"recipient"`
	Attachment  Attachment    `json:"attachment,omitempty"`
}

func CreateEmptyRawTransactionAndHash() (string, string, error) {

	return "", "", nil
}
func (tx *TransferV2) Sign(secretKey SecretKey) error {
	b, err := tx.BodyMarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to sign TransferV2 transaction")
	}
	if tx.Proofs == nil {
		tx.Proofs = &ProofsV1{proofsVersion, make([]B58Bytes, 0)}
	}
	err = tx.Proofs.Sign(0, secretKey, b)
	if err != nil {
		return errors.Wrap(err, "failed to sign TransferV2 transaction")
	}
	d, err := FastHash(b)
	tx.ID = &d
	if err != nil {
		return errors.Wrap(err, "failed to sign TransferV2 transaction")
	}
	return nil
}
func SignRawTransaction(hash string, prikey []byte) ([]byte, error) {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, errors.New("invalid hash message")
	}
	sig, retCode := owcrypt.Signature(prikey, nil, 0, hashBytes, 32, owcrypt.ECC_CURVE_SECP256K1)

	if retCode != owcrypt.SUCCESS {
		return nil, errors.New("sign failed!")
	}

	return sig, nil
}
func VerifyAndCombinRawTransaction(emptyTrans string, signature, publicKey string) (bool, string) {

	return true, ""
}
func NewUnsignedTransferV2(senderPK PublicKey, amountAsset, feeAsset OptionalAsset, timestamp, amount, fee uint64, recipient Recipient, attachment string) *TransferV2 {
	t := Transfer{
		SenderPK:    senderPK,
		Recipient:   recipient,
		AmountAsset: amountAsset,
		Amount:      amount,
		FeeAsset:    feeAsset,
		Fee:         fee,
		Timestamp:   timestamp,
		Attachment:  Attachment(attachment),
	}
	return &TransferV2{Version: 2, Transfer: t}

}

func (tx *TransferV2) MarshalBinary() ([]byte, error) {
	bb, err := tx.BodyMarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal TransferV2 transaction to bytes")
	}
	bl := len(bb)
	if tx.Proofs == nil {
		return nil, errors.New("failed to marshal TransferV2 transaction to bytes: no proofs")
	}
	pb, err := tx.Proofs.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal TransferV2 transaction to bytes")
	}
	buf := make([]byte, 1+bl+len(pb))
	buf[0] = 0
	copy(buf[1:], bb)
	copy(buf[1+bl:], pb)
	return buf, nil
}
func (tx *TransferV2) BodyMarshalBinary() ([]byte, error) {
	b, err := tx.Transfer.marshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal TransferV2 body")
	}
	buf := make([]byte, 2+len(b))
	buf[0] = byte(4)
	buf[1] = tx.Version
	copy(buf[2:], b)
	return buf, nil
}

func (tr *Transfer) marshalBinary() ([]byte, error) {
	p := 0
	aal := 0
	if tr.AmountAsset.Present {
		aal += DigestSize
	}
	fal := 0
	if tr.FeeAsset.Present {
		fal += DigestSize
	}
	rb, err := tr.Recipient.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal Transfer body")
	}
	rl := len(rb)
	atl := len(tr.Attachment)
	buf := make([]byte, transferLen+aal+fal+atl+rl)
	copy(buf[p:], tr.SenderPK[:])
	p += PublicKeySize
	aab, err := tr.AmountAsset.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal Transfer body")
	}
	copy(buf[p:], aab)
	p += 1 + aal
	fab, err := tr.FeeAsset.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal Transfer body")
	}
	copy(buf[p:], fab)
	p += 1 + fal
	binary.BigEndian.PutUint64(buf[p:], tr.Timestamp)
	p += 8
	binary.BigEndian.PutUint64(buf[p:], tr.Amount)
	p += 8
	binary.BigEndian.PutUint64(buf[p:], tr.Fee)
	p += 8
	copy(buf[p:], rb)
	p += rl
	PutStringWithUInt16Len(buf[p:], tr.Attachment.String())
	return buf, nil
}
func PutStringWithUInt16Len(buf []byte, s string) {
	sl := uint16(len(s))
	binary.BigEndian.PutUint16(buf, sl)
	copy(buf[2:], s)
}

type ProofsV1 struct {
	Version byte
	Proofs  []B58Bytes
}

func (p *ProofsV1) binarySize() int {
	pl := 0
	if p != nil {
		for _, e := range p.Proofs {
			pl += len(e) + 2
		}
	}
	return proofsMinLen + pl
}
func (p *ProofsV1) MarshalBinary() ([]byte, error) {
	buf := make([]byte, p.binarySize())
	pos := 0
	buf[pos] = proofsVersion
	pos++
	binary.BigEndian.PutUint16(buf[pos:], uint16(len(p.Proofs)))
	pos += 2
	for _, e := range p.Proofs {
		el := len(e)
		binary.BigEndian.PutUint16(buf[pos:], uint16(el))
		pos += 2
		copy(buf[pos:], e)
		pos += el
	}
	return buf, nil
}

//Sign creates a signature and stores it as a proof at given position.
func (p *ProofsV1) Sign(pos int, key SecretKey, data []byte) error {
	if pos < 0 || pos > proofsMaxCount {
		return errors.Errorf("failed to create proof at position %d, allowed positions from 0 to %d", pos, proofsMaxCount-1)
	}
	if len(p.Proofs)-1 < pos {
		s, err := Sign(key, data)
		if err != nil {
			return errors.Errorf("crypto.Sign(): %v", err)
		}
		p.Proofs = append(p.Proofs[:pos], append([]B58Bytes{s[:]}, p.Proofs[pos:]...)...)
	} else {
		pr := p.Proofs[pos]
		if len(pr) > 0 {
			return errors.Errorf("unable to overwrite non-empty proof at position %d", pos)
		}
		s, err := Sign(key, data)
		if err != nil {
			return errors.Errorf("crypto.Sign(): %v", err)
		}
		copy(pr[:], s[:])
	}
	return nil
}

type OptionalAsset struct {
	Present bool
	ID      Digest
}

func (a OptionalAsset) binarySize() int {
	s := 1
	if a.Present {
		s += DigestSize
	}
	return s
}

//MarshalBinary marshals the optional asset to its binary representation.
func (a OptionalAsset) MarshalBinary() ([]byte, error) {
	buf := make([]byte, a.binarySize())
	PutBool(buf, a.Present)
	if a.Present {
		copy(buf[1:], a.ID[:])
	}
	return buf, nil
}
func PutBool(buf []byte, b bool) {
	if b {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
}

func NewOptionalAssetFromString(s string) (*OptionalAsset, error) {
	switch strings.ToUpper(s) {
	case WavesAssetName, "":
		return &OptionalAsset{Present: false}, nil
	default:
		a, err := NewDigestFromBase58(s)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create OptionalAsset from Base58 string")
		}
		return &OptionalAsset{Present: true, ID: a}, nil
	}
}
