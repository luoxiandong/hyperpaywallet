package hdwallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	wiccwallet "github.com/WaykiChain/wicc-wallet-utils-go"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons"
	"github.com/btcsuite/btcutil"
	"io/ioutil"
	"net/http"
	"strings"
)

func init() {
	coins[WICC] = newWICC
}

const NET_TYPE = commons.MAINNET

type wicc struct {
	*btc
}

type blockCount struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data int64  `json:"data"`
}

type account struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data accountData `json:"data"`
}

type accountData struct {
	Address string `json:"address"`
	Regid   string `json:"regid"`
}

func newWICC(key *Key) Wallet {
	token := newBTC(key).(*btc)
	token.name = "WaykiChain"
	token.symbol = "WICC"

	return &wicc{btc: token}
}

func (c *wicc) GetKey() *Key {
	return c.key
}

func (c *wicc) GetAddress() (string, error) {
	pWif, err := c.GetKey().PrivateWIF(true)
	if err != nil {
		return "", err
	}
	address := wiccwallet.GetAddressFromPrivateKey(pWif, int(NET_TYPE))

	return address, nil
}

func (c *wicc) GetWalletAccount() *WalletAccount {
	if c.GetKey().Extended == nil {
		return &WalletAccount{
			Res: 0,
		}
	}

	address, err := c.GetAddress()
	if err != nil {
		return &WalletAccount{
			Res:    0,
			ErrMsg: err.Error(),
		}
	}
	pWif, err1 := c.GetKey().PrivateWIF(true)
	if err1 != nil {
		return &WalletAccount{
			Res:    0,
			ErrMsg: err1.Error(),
		}
	}
	publicKey := c.GetKey().PublicHex(true)

	return &WalletAccount{
		Res:        1,
		PrivateKey: pWif,
		Address:    address,
		PublicKey:  publicKey,
		Seed:       c.GetKey().Seed,
	}
}

func (c *wicc) GetWalletAccountFromWif() (*WalletAccount, error) {
	hexPri := c.GetKey().Wif
	if len(hexPri) > 0 {
		pri, err := btcutil.DecodeWIF(hexPri)
		if err != nil {
			return nil, err
		}

		pub := hex.EncodeToString(pri.SerializePubKey())
		address := commons.GetAddressFromPrivateKey(hexPri, commons.Network(int(NET_TYPE)))
		return &WalletAccount{
			Res:        1,
			Address:    address,
			PrivateKey: hexPri,
			PublicKey:  pub,
		}, nil
	}
	return &WalletAccount{
		Res:        0,
		PrivateKey: hexPri,
	}, nil
}

// 获取当前区块的高度
func GetWICCBlockCount() (int64, error) {
	var url string
	switch NET_TYPE {
	case commons.MAINNET:
		url = "https://baas.wiccdev.org/v2/api/block/getblockcount"
	default:
		url = "https://baas-test.wiccdev.org/v2/api/block/getblockcount"
	}
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return 0, err

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var blockCount blockCount
	err = json.Unmarshal(body, &blockCount)
	if err != nil {
		return 0, err
	}

	if blockCount.Code > 0 {
		return 0, errors.New(blockCount.Msg)
	}

	return blockCount.Data, nil
}

// 获取地址对应的Reg ID
func GetAccountRegId(addr string) (string, error) {
	var url string
	switch NET_TYPE {
	case commons.MAINNET:
		url = "https://baas.wiccdev.org/v2/api/account/getaccountinfo"
	default:
		url = "https://baas-test.wiccdev.org/v2/apiaccount/getaccountinfo"
	}
	req := make(map[string]string)
	req["address"] = addr
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(reqData)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var accoutResponse account
	err = json.Unmarshal(body, &accoutResponse)
	if err != nil {
		return "", err
	}

	if accoutResponse.Code > 0 {
		return "", errors.New(accoutResponse.Msg)
	}

	if accoutResponse.Data.Regid == "" {
		return "", errors.New("WICC Address: " + addr + " RegID Empty")
	}

	return accoutResponse.Data.Regid, nil
}

func (c *wicc) SignRawTransaction(signIn *SignInput) (*SignResult, error) {
	var signedTrans string
	height, err := GetWICCBlockCount() // 获取当前区块高度
	if err != nil {
		return &SignResult{
			Res:    0,
			ErrMsg: "WICC Get BlockCount err",
		}, err
	}
	pubKey, _ := wiccwallet.GetPubKeyFromPrivateKey(signIn.PrivateKey)

	var txParam wiccwallet.UCoinTransferTxParam
	txParam.FeeSymbol = strings.ToUpper(string(commons.WICC)) //Fee Type (WICC/WUSD)
	txParam.Fees = signIn.Fee                                 // fees for mining
	txParam.ValidHeight = height                              // valid height Within the height of the latest block
	txParam.SrcRegId = ""
	txParam.Dests = wiccwallet.NewDestArr()
	dest1 := wiccwallet.Dest{
		CoinSymbol: strings.ToUpper(signIn.Symbol), // From Coin Type
		CoinAmount: uint64(signIn.Amount),          // the values send to the contract app
		DestAddr:   signIn.DestAddr,                // To address
	}
	txParam.Dests.Add(&dest1)
	txParam.PubKey = pubKey
	txParam.Memo = ""
	signedTrans, err = wiccwallet.SignUCoinTransferTx(signIn.PrivateKey, &txParam)
	if err != nil {
		return &SignResult{
			Res:    0,
			ErrMsg: "WICC UCoinTransferTx err",
		}, err
	}

	return &SignResult{
		Res:   1,
		RawTX: signedTrans,
	}, nil
}
