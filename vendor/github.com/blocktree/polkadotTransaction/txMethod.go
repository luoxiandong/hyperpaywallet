package polkadotTransaction

import (
	"encoding/hex"
	"errors"

	"github.com/blocktree/polkadotTransaction/codec"
	"math/big"
)

type MethodTransfer struct {
	DestPubkey []byte
	Amount     []byte
}

func NewMethodTransfer(pubkey string, amount string) (*MethodTransfer, error) {
	pubBytes, err := hex.DecodeString(pubkey)
	if err != nil || len(pubBytes) != 32 {
		return nil, errors.New("invalid dest public key")
	}

	if amount == "0" || amount == "" {
		return nil, errors.New("zero amount")
	}

	n := big.NewInt(0)
	n, _ = n.SetString(amount, 10)

	amountStr, err := codec.Encode("Compact<u32>", *n)
	if err != nil {
		return nil, errors.New("invalid amount")
	}

	amountBytes, _ := hex.DecodeString(amountStr)
	return &MethodTransfer{
		DestPubkey: pubBytes,
		Amount:     amountBytes,
	}, nil
}

func (mt MethodTransfer) ToBytes(transferCode string) ([]byte, error) {

	if mt.DestPubkey == nil || len(mt.DestPubkey) != 32 || mt.Amount == nil || len(mt.Amount) == 0 {
		return nil, errors.New("invalid method")
	}

	ret, _ := hex.DecodeString(transferCode)
	if transferCode == DOT_Balannce_Transfer || transferCode == EDG_Balannce_Transfer{
		ret = append(ret, 0x00)
	} else if AccounntIDFollow || transferCode == PLM_Balannce_Transfer || transferCode == FIS_Balannce_Transfer {
		ret = append(ret, 0xff)
	}

	ret = append(ret, mt.DestPubkey...)
	ret = append(ret, mt.Amount...)

	return ret, nil
}
