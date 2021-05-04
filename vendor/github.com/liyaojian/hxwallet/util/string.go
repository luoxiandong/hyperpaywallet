package util

import (
	"fmt"
	"math"
	"strconv"
)

const (
	HXBind              = "bind"
	HXUbind             = "ubind"
	HXTransfer          = "transfer"
	HXWithdraw          = "withdraw"
	HXRegister          = "register"
	HXMining            = "mining"
	HXRewards           = "rewards"
	HxRedeem            = "redeem"
	HXContractInvoke    = "contractinvoke"
	HXContractTransfer  = "contracttransfer"
	minerInvalidAddress = "InvalidAddress"

	fieldCrossWif           = "crossWif"
	fieldFromAddr           = "fromAddr"
	fieldToAddr             = "toAddr"
	fieldAccountName        = "accountName"
	fieldOrigAddr           = "origAddr"
	fieldRefBlock           = "refBlock"
	fieldChainID            = "chainId"
	fieldPrecision          = "precision"
	fieldTranFee            = "tranFee"
	fieldTranAmt            = "tranAmt"
	fieldCoinAddr           = "coinAddr"
	fieldCoinType           = "coinType"
	fieldAccountId          = "accountId"
	fieldAssetId            = "assetId"
	fieldPayBackList        = "payBackList"
	fieldCitizenName        = "citizenName"
	fieldCitizenAmount      = "amount"
	fieldMainAssetPrecision = "mainCoinPrecision"
	fieldContractAPI        = "contractApi"
	fieldContractID         = "contractId"
	fieldBasicFee           = "basicFee"
	fieldGasPrice           = "gasPrice"
	fieldGasLimit           = "gasLimit"

	hxPrecisionBits = 5
	hxPrecision     = 100000
)

// convert string to float64, multiple precision
func GetInt64(jmap map[string]interface{}, field string, precision int64) (int64, error) {
	is, ok := jmap[field]
	if !ok {
		return 0, fmt.Errorf("not found field param %v", field)
	}
	s, ok := is.(string)
	if !ok {
		return 0, fmt.Errorf("field %v is invalid format", field)
	}

	ii, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(ii * float64(precision))), nil
}
func GetHxInt64(value string) (int64, error) {
	ii, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(ii * float64(hxPrecision))), nil
}

func GetIntField(jmap map[string]interface{}, field string) (int64, error) {
	is, ok := jmap[field]
	if !ok {
		return 0, fmt.Errorf("not found field param %v", field)
	}
	fi, ok := is.(float64)
	if !ok {
		return 0, fmt.Errorf("convert param %s failed", field)
	}
	return int64(fi), nil
}

func GetStringField(jmap map[string]interface{}, field string) string {
	is, ok := jmap[field]
	if !ok {
		fmt.Printf("not found field %v in json map %v\n", field, jmap)
		return ""
	}
	s := is.(string)
	return s
}

func GetStringFieldWithDefault(jmap map[string]interface{}, field, def string) string {
	is, ok := jmap[field]
	if !ok {
		fmt.Printf("not found field %v in json map %v, return default %v\n", field, jmap, def)
		return def
	}
	s := is.(string)
	return s
}
