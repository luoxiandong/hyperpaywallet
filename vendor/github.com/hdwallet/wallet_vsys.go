package hdwallet

import (
	"encoding/json"
	vsyswallet "github.com/virtualeconomy/go-v-sdk/vsys"
)

func init() {
	coins[VSYS] = newVSYS
}

type vsys struct {
	*btc
}

func newVSYS(key *Key) Wallet {
	token := newBTC(key).(*btc)
	token.name = "VSystems"
	token.symbol = "VSYS"

	return &vsys{btc: token}
}

func (c *vsys) SignRawTransaction(signIn *SignInput) (*SignResult, error) {
	var rawTxByte []byte
	acc := vsyswallet.InitAccount(vsyswallet.Mainnet)
	acc.BuildFromPrivateKey(signIn.PrivateKey)
	if signIn.Symbol == signIn.Coin {
		tx := acc.BuildPayment(signIn.DestAddr, signIn.Amount, "")
		rawTxByte, _ = json.Marshal(tx)
	} else {
		a := &vsyswallet.Contract{
			Amount:    signIn.Amount,
			Recipient: signIn.DestAddr,
		}
		funcData := a.BuildSendData()
		tx := acc.BuildExecuteContract(
			signIn.ContractAddr,
			vsyswallet.FuncidxSend,
			funcData,
			"")
		rawTxByte, _ = json.Marshal(tx)
	}

	return &SignResult{
		Res:   1,
		RawTX: string(rawTxByte),
	}, nil
}
