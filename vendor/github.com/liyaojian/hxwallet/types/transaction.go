package types

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/liyaojian/hxwallet/util"
	"strings"
	"time"
)

const expireTimeout = 86000

type Transaction struct {
	HxRefBlockNum    uint16          `json:"ref_block_num"`
	HxRefBlockPrefix uint32          `json:"ref_block_prefix"`
	HxExpiration     string          `json:"expiration"`
	HxOperations     [][]interface{} `json:"operations"`
	HxExtensions     []interface{}   `json:"extensions"`
	HxSignatures     []string        `json:"signatures"`

	Expiration uint32        `json:"-"`
	Operations []interface{} `json:"-"`
}

func DefaultTransaction() *Transaction {

	return &Transaction{
		0,
		0,
		"",
		nil,
		nil,
		nil,
		0,
		nil,
	}
}

func BuildTransferTransaction(refinfo, wif string, from, to, memo, assetId string, amount, fee int64,
	symbol string, guarantee_id, chain_id string) (b []byte, err error) {

	asset_amount := DefaultAsset()
	asset_amount.HxAmount = amount
	asset_amount.HxAssetId = assetId // SetAssetBySymbol(symbol)

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.HxAmount = fee
	asset_fee.SetAssetBySymbol("HX")

	transferOp := DefaultTransferOperation()
	transferOp.HxFee = asset_fee
	transferOp.HxFromAddr = from
	transferOp.HxToAddr = to
	transferOp.HxAmount = asset_amount

	if memo == "" {
		transferOp.HxMemo = nil
	} else {
		memo_trx := DefaultMemo()
		memo_trx.Message = memo
		memo_trx.IsEmpty = false
		memo_trx.HxMessage = hex.EncodeToString(append(make([]byte, 4), []byte(memo_trx.Message)...))
		transferOp.HxMemo = &memo_trx
	}

	if guarantee_id != "" {
		transferOp.HxGuaranteeId = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := util.Time2Str(expir_sec)
	//expir_str := "2018-09-26T09:14:40"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := util.GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{0, transferOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*transferOp},
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.HxSignatures = append(transferTrx.HxSignatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

func GetSignature(wif string, hash []byte) ([]byte, error) {
	ecPrivkey, err := NewPrivateKeyFromWif(wif)
	if err != nil {
		return nil, fmt.Errorf("in GetSignature function, get ecprivkey failed: %v", err)
	}
	ecPrivkeyByte := ecPrivkey.priv.Serialize()

	return util.SignCompact(hash, ecPrivkeyByte, true)
}

func GetAssetId(coinType string) string {
	coinType = strings.ToUpper(coinType)
	switch coinType {
	case "HX":
		return "1.3.0"
	case "BTC":
		return "1.3.1"
	case "LTC":
		return "1.3.2"
	case "HC":
		return "1.3.3"
	case "ETH":
		return "1.3.4"
	case "ERCPAX":
		return "1.3.5"
	case "ERCELF":
		return "1.3.6"
	case "USDT":
		return "1.3.7"
	case "BCH":
		return "1.3.8"
	case "ERCTITAN":
		return "1.3.9"
	default:
		return "1.3.999"
	}
}
