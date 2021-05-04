package hdwallet

import (
	"github.com/liyaojian/hxwallet/api/jsonrpc"
	"github.com/liyaojian/hxwallet/rpc/http"
	"github.com/liyaojian/hxwallet/types"
	"github.com/liyaojian/hxwallet/util"
)

const hxUrl = "http://47.88.171.179:8090"

func init() {
	coins[HX] = newHX
}

type hx struct {
	*btc
}

func newHX(key *Key) Wallet {
	token := newBTC(key).(*btc)
	token.name = "HyperExchange"
	token.symbol = "HX"

	return &hx{btc: token}
}

func (c *hx) SignRawTransaction(signIn *SignInput) (*SignResult, error) {
	db := jsonrpc.NewAPI("database", http.NewTransport(hxUrl))
	dynamicGlobalProperties, err := db.GetDynamicGlobalProperties()
	if err != nil {
		return nil, err
	}
	refInfo := util.CalRefInfo(dynamicGlobalProperties.HeadBlockID)        // https://docs.gxchain.org/zh/guide/apis.html#get-dynamic-global-properties 获取 head_block_id
	chainId, err := db.GetChainId() // 链ID可通过 https://docs.gxchain.org/zh/guide/apis.html#get-chain-id 获取
	if err != nil {
		return nil, err
	}

	trxData, err := types.Transfer(refInfo, signIn.PrivateKey, chainId, signIn.SrcAddr, signIn.DestAddr, signIn.Coin, signIn.Amount, signIn.Fee, signIn.Memo, "")
	if err != nil {
		return nil, err
	}
	txHex := string(trxData)

	return &SignResult{
		Res:   1,
		RawTX: txHex,
	}, nil
}

func (c *hx) GetWalletAccountFromWif() (*WalletAccount, error) {
	wif := c.GetKey().Wif
	if len(wif) > 0 {
		wallet, err := types.NewPrivateKeyFromWif(wif)
		if err != nil {
			return nil, err
		}

		return &WalletAccount{
			Res:        1,
			PrivateKey: wif,
			PublicKey:  wallet.PublicKey().String(),
			Address:    wallet.ToWIF(),
		}, nil
	}
	return &WalletAccount{
		Res: 0,
	}, nil
}
