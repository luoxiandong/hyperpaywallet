package tron

import (
	"fmt"
	"math/big"
	"strings"
)

/*

在TRON中检测TRX或TRC10事务涉及4种类型的合同：

TransferContract（系统合同类型：TRX转账）
TransferAssetContract（系统合同类型：TRC10转账）
CreateSmartContract（智能合约类型）
TriggerSmartContract（智能合约类型：TRC20转账）
Transaction，TransactionInfo 和 Block 的数据包含所有智能合约交易信息。

技术细节
https://cn.developers.tron.network/docs/%E4%BA%A4%E6%8D%A2%E4%B8%AD%E7%9A%84trc10%E5%92%8Ctrx%E8%BD%AC%E7%A7%BB

TRX转账示例
https://tronscan.org/#/transaction/f8f8ac5b4b0df34dad410147231061806c9fa8c207e7f3107cadc6d00925ccbc

TRC10转账示例
https://tronscan.org/#/transaction/c0edfc83e3535700b46598444f2425696686d20566101d8b5b2aa95c0915a2a0

TRC20转账示例
https://tronscan.org/#/transaction/a5614f60e7d3b9d8859abe89968d81007c321c5ad83cb9c7abaa736a20401a11

*/

const (
	TRC10 = "trc10"
	TRC20 = "trc20"

	FeeLimit = 10000000
)

const (
	TRC20_BALANCE_OF_METHOD  = "balanceOf(address)"
	TRC20_TRANSFER_METHOD_ID = "a9059cbb"
	TRX_TRANSFER_EVENT_ID    = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

const (
	SOLIDITY_TYPE_ADDRESS = "address"
	SOLIDITY_TYPE_UINT256 = "uint256"
	SOLIDITY_TYPE_UINT160 = "uint160"
)

type SolidityParam struct {
	ParamType  string
	ParamValue interface{}
}

func makeRepeatString(c string, count uint) string {
	cs := make([]string, 0)
	for i := 0; i < int(count); i++ {
		cs = append(cs, c)
	}
	return strings.Join(cs, "")
}

func MakeTransactionParameter(methodId string, params []SolidityParam) (string, error) {

	data := methodId
	for i, _ := range params {
		var param string
		if params[i].ParamType == SOLIDITY_TYPE_ADDRESS {
			param = strings.ToLower(params[i].ParamValue.(string))
			param = strings.TrimPrefix(param, "0x")
			if len(param) != 42 {
				return "", fmt.Errorf("length of address error.")
			}
			param = makeRepeatString("0", 22) + param
		} else if params[i].ParamType == SOLIDITY_TYPE_UINT256 {
			intParam := params[i].ParamValue.(*big.Int)
			param = intParam.Text(16)
			l := len(param)
			if l > 64 {
				return "", fmt.Errorf("integer overflow.")
			}
			param = makeRepeatString("0", uint(64-l)) + param
			//fmt.Println("makeTransactionData intParam:", intParam.String(), " param:", param)
		} else {
			return "", fmt.Errorf("not support solidity type")
		}

		data += param
	}
	return data, nil
}
