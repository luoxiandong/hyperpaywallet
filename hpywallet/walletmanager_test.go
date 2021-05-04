package hpywallet_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/blocktree/go-owcdrivers/bitcoincashTransaction"
	"github.com/blocktree/go-owcdrivers/btcTransaction"
	"github.com/blocktree/go-owcdrivers/omniTransaction"
	"github.com/btcsuite/btcutil"
	"github.com/liyaojian/hxwallet/keypair"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"hyperpaywallet/hpywallet"
	"testing"
)

var (
	testmnemonic = "hint fatigue scale tired venture regular vicious echo satoshi gun cash spy"
	mnemonic     = "zero guard grass mandate invest anger actress moral gasp easy way student"
	ontmnemonic  = "deposit syrup useful goat enroll snow goddess year acid gravity rocket fatigue"
)

func TestGenerateMnemonic(t *testing.T) {
	mnemonic := hpywallet.GenerateMnemonic()
	fmt.Println("Mnemonic: ", mnemonic)
}
func TestGenerateSeed(t *testing.T) {
	seed := hpywallet.GenerateSeed(mnemonic, "")
	fmt.Println("Seed1 = ", seed)
	seed1 := hpywallet.GenerateSeed(mnemonic, "123")

	fmt.Println("Seed2 = ", seed1)
	wallet1 := hpywallet.GenerateSeedWallet(seed, "btc")
	wallet2 := hpywallet.GenerateSeedWallet(seed1, "btc")
	fmt.Println("Wallet1 = ", wallet1)
	fmt.Println("Wallet2 = ", wallet2)

	wallet3 := hpywallet.GenerateMnemonicWallet(mnemonic, "123456", "btc")
	wallet4 := hpywallet.GenerateMnemonicWallet(mnemonic, "12345", "btc")

	fmt.Println("Wallet3 = ", wallet3)
	fmt.Println("Wallet4 = ", wallet4)

}

func TestKeystore(t *testing.T) {
	privateKey := "KxkUeg9G2ghaZHewELhjuDUtKtnbtP95pH3g8Siswvu5uNxxFPRc"
	pwd := "11111"
	udid := "AOIJF-QWEQR-VDFBET-YTAWWE"

	// Encode
	enResult := hpywallet.EnKeystore(privateKey, pwd, udid)
	if enResult.ResCode != 1 {
		fmt.Println("Error: ", enResult.ErrMsg)
		return
	}
	fmt.Println("Keystore : \n", enResult.Result)

	fmt.Println("************************************************")
	// Decode
	deResult := hpywallet.DeKeystore(enResult.Result, pwd, udid)
	if deResult.ResCode != 1 {
		fmt.Println("Error: ", deResult.ErrMsg)
		return
	}
	fmt.Println("PrivateKey : ", deResult.Result)

}

func TestCreateDoge(t *testing.T) {
	dogeWallet := hpywallet.GenerateWallet(mnemonic, "doge")
	toWallet := hpywallet.GenerateWallet(testmnemonic, "doge")
	importWallet := hpywallet.ImportPrivateWIF(dogeWallet.PrivateKey, "doge")

	fmt.Println("dogeWallet: ", dogeWallet)
	fmt.Println("toWallet: ", toWallet)
	fmt.Println("importWallet: ", importWallet)

	// need utxo
	signInput := &hpywallet.SignInput{
		Coin:       "doge",
		Symbol:     "doge",
		PrivateKey: dogeWallet.PrivateKey,
		SrcAddr:    dogeWallet.Address,
		DestAddr:   toWallet.Address,
		Fee:        10000,
		Amount:     2000000000,
	}

	tranferResult := hpywallet.SignRawTransaction(signInput)
	if tranferResult.ResCode == 0 {
		fmt.Println("Fail! Transfer doge Msg: ", tranferResult.ErrMsg)
	} else {
		fmt.Println("Success! Transfer doge RawTX: ", tranferResult.RawTX)
	}

	// https://dogechain.info/address/DNzF31pAuNjofu27ukV3NjWCcy6uB9haFR
	// https://dogechain.info/api/v1/unspent/DNzF31pAuNjofu27ukV3NjWCcy6uB9haFR?nsukey=ANpE25lxuDndL1c%2BTKD4rJcpG58RT9EkO6FB%2FNh1xAJNDX8coA8WdflJ%2BQgp%2FkUwYxiA1aTajX0b67nYpihZI7%2BaMoR4s7oS5VW8pbOQ4ky3MMDyy%2BJ9yFyGTeFLZcdsG6yzmoXTk4V6hYZEBY9zJLe3xgle76BdIzLyIfZ%2FKLB097prwUVV6QEk3sv%2BWH3iYylvxma%2BwPFMN2YGZZSR2Q%3D%3D
}

func TestCreateDASH(t *testing.T) {
	wallet := hpywallet.GenerateWallet(mnemonic, "dash") //
	toWallet := hpywallet.GenerateWallet(testmnemonic, "dash")
	// importWallet := hpywallet.ImportPrivateWIF(wallet.PrivateKey, "dash")

	fmt.Println("DASH Wallet: ", wallet)
	fmt.Println("toWallet: ", toWallet)

	fmt.Println("From PrivateKey：", wallet.PrivateKey)

	fmt.Println("To Address：", toWallet.Address)

	item1 := hpywallet.OutPutItem{
		TxHash:   "f620269227d91079cc1c499047b5cfd9c4a0ec6c3ddf220e20dd5761a9c26b7c",
		Value:    1000000,
		Vout:     0,
		Pkscript: "76a91426739443ce47332bb1d21cd8eae9039da3f9602e88ac",
	}

	outputs := []hpywallet.OutPutItem{item1}
	// fmt.Println("ltc outputs: ", outputs)

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		fmt.Println("outputs err: ", err.Error())
	}

	signInput := &hpywallet.SignInput{
		Coin:       "dash",
		Symbol:     "dash",
		PrivateKey: wallet.PrivateKey,
		SrcAddr:    wallet.Address,
		DestAddr:   toWallet.Address,
		Change:     11000,
		Amount:     10000,
		Inputs:     jsonInputs,
	}
	// fmt.Println("Transfer dash signInput: ", signInput)
	tranferResult := hpywallet.SignRawTransaction(signInput)
	if tranferResult.ResCode == 0 {
		fmt.Println("错误 Transfer dash Msg: ", tranferResult.ErrMsg)
	} else {
		fmt.Println("成功 Transfer dash RawTX: ", tranferResult.RawTX)
	}

	//  https://chainz.cryptoid.info/dash/address.dws?XhqADVUCcqN93C92JToCerMnqWTAZdTgoM.htm
	//  https://explorer.dash.org/chain/Dash
	//  https://www.thepolyglotdeveloper.com/2018/03/create-bitcoin-hardware-wallet-golang-raspberry-pi-zero/

}

func Test_TransferBCH(t *testing.T) {
	wallet := hpywallet.GenerateWallet(mnemonic, "bch") //

	hpwallet := hpywallet.ImportPrivateWIF(wallet.PrivateKey, "bch")
	fmt.Println("bch Address: ", hpwallet.Address)
	fmt.Println("bch Wif: ", hpwallet.PrivateKey)
	item1 := hpywallet.OutPutItem{
		TxHash:   "b211224d3a773dd7566ba2c221125547df3bbb011c7626021ab5e9f6de7fc112",
		Value:    1890000,
		Vout:     1,
		Pkscript: "76a9143a97cd827522fd88d97ee1b44eaa0ed37cb0585b88ac",
	}

	item2 := hpywallet.OutPutItem{
		TxHash:   "b211224d3a773dd7566ba2c221125547df3bbb011c7626021ab5e9f6de7fc112",
		Value:    10000,
		Vout:     0,
		Pkscript: "76a9143a97cd827522fd88d97ee1b44eaa0ed37cb0585b88ac",
	}

	item3 := hpywallet.OutPutItem{
		TxHash:   "c02179b60aefeed32fe003ddc699f2e9aeaf2c8d8d495fa0cbc0d5f04b8fea49",
		Value:    10000,
		Vout:     0,
		Pkscript: "76a9143a97cd827522fd88d97ee1b44eaa0ed37cb0585b88ac",
	}

	outputs := []hpywallet.OutPutItem{item1, item2, item3}
	// fmt.Println("ltc outputs: ", outputs)

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		//log.Fatal("Cannot encode to JSON ", err)
		fmt.Println("outputs err: ", err.Error())

	}
	// fmt.Println("ltc outputs: ", jsonInputs)

	signInput := &hpywallet.SignInput{
		Coin:       "bch",
		Symbol:     "bch",
		PrivateKey: wallet.PrivateKey,
		SrcAddr:    "16Lp3ZvcusRGtqy7DgF5gz6PGutvuEtwRt",
		DestAddr:   "16Lp3ZvcusRGtqy7DgF5gz6PGutvuEtwRt",
		Amount:     900000,
		Change:     1000000,
		Inputs:     jsonInputs,
	}
	tranferResult := hpywallet.SignRawTransaction(signInput)
	if tranferResult.ResCode == 0 {
		fmt.Println("Transfer bch Msg: ", tranferResult.ErrMsg)
	} else {
		fmt.Println("Transfer bch RawTX: ", tranferResult.RawTX)
	}

}

func Test_TransferQTUMQRC20(t *testing.T) {
	qtumWallet := hpywallet.GenerateWallet(mnemonic, "qtum")
	item1 := hpywallet.OutPutItem{
		TxHash:   "c774983ae03dd2fd3d29022899ba9d26ae7792ae0b939e1909b6389b261fc109",
		Value:    43000000,
		Vout:     1,
		Pkscript: "76a914507d2234de017230c7cd4e9971c13496bde771c488ac",
	}

	outputs := []hpywallet.OutPutItem{item1}
	// fmt.Println("ltc outputs: ", outputs)

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		//log.Fatal("Cannot encode to JSON ", err)
		fmt.Println("outputs err: ", err.Error())

	}
	// fmt.Println("ltc outputs: ", jsonInputs)

	signInput := &hpywallet.SignInput{
		Coin:         "qtum",
		Symbol:       "hpy",
		PrivateKey:   qtumWallet.PrivateKey,
		SrcAddr:      qtumWallet.Address,
		DestAddr:     "QTa3opXhHQSD1kwJmHiFS2TY3DV1RPPP4n",
		Fee:          10000000,
		Amount:       1000000,
		Change:       33000000,
		ContractAddr: "f2703e93f87b846a7aacec1247beaec1c583daa4",
		Inputs:       jsonInputs,
	}
	tranferResult := hpywallet.SignRawTransaction(signInput)
	if tranferResult.ResCode == 0 {
		fmt.Println("Transfer qtum Msg: ", tranferResult.ErrMsg)
	} else {
		fmt.Println("Transfer qtum RawTX: ", tranferResult.RawTX)
	}

}
func Test_TransferBTC(t *testing.T) {
	btcWallet := hpywallet.GenerateWallet(mnemonic, "btc")
	fmt.Println("btc Address: ", btcWallet.Address)
	fmt.Println("btc Wif: ", btcWallet.PrivateKey)
	// fmt.Println("ltc PublicKey: ", hcWallet.PublicKey)
	item1 := hpywallet.OutPutItem{
		TxHash:   "921784b1e11fcbfe267a04b9e54a45e597710d0f9413e658813737c06f44a987",
		Value:    1200000,
		Vout:     1,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	item2 := hpywallet.OutPutItem{
		TxHash:   "62409f73492acc32aa7760423e56127f78f62745f07a68c4107ea69394dce485",
		Value:    400000,
		Vout:     1,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	item3 := hpywallet.OutPutItem{
		TxHash:   "b886364bf351ac4ae47e4eea8ab4e1039c4fd92fed04899bb360d0570a0f92a2",
		Value:    100000,
		Vout:     0,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	item4 := hpywallet.OutPutItem{
		TxHash:   "ca2695a680d6e7bea7dc5155f9cbe13f3b230a648dd48135ff5b090ae9c89194",
		Value:    100000,
		Vout:     1,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	item5 := hpywallet.OutPutItem{
		TxHash:   "761faedb198081b97478ae3bc6f85deefc43ab5d1e116fe2a83ff33bd9dbea14",
		Value:    15500000,
		Vout:     1,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	outputs := []hpywallet.OutPutItem{item1, item2, item3, item4, item5}
	// fmt.Println("ltc outputs: ", outputs)

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		//log.Fatal("Cannot encode to JSON ", err)
		fmt.Println("outputs err: ", err.Error())

	}
	// fmt.Println("ltc outputs: ", jsonInputs)

	signInput := &hpywallet.SignInput{
		Coin:       "btc",
		Symbol:     "btc",
		PrivateKey: btcWallet.PrivateKey,
		SrcAddr:    btcWallet.Address,
		DestAddr:   btcWallet.Address,
		Fee:        200000,
		Amount:     100000,
		Change:     900000,
		Inputs:     jsonInputs,
	}
	tranferResult := hpywallet.SignRawTransaction(signInput)
	if tranferResult.ResCode == 0 {
		fmt.Println("Transfer btc Msg: ", tranferResult.ErrMsg)
	} else {
		fmt.Println("Transfer btc RawTX: ", tranferResult.RawTX)
	}

}
func Test_TransferBCD(t *testing.T) {
	//bcdWallet := hpywallet.GenerateWallet(mnemonic, "bcd")
	//fmt.Println("bcd Address: ", bcdWallet.Address)
	//fmt.Println("bcd Wif: ", bcdWallet.PrivateKey)
	bcdWallet := hpywallet.ImportPrivateWIF("Ky3n4DNba9iyyXPHyjsYKiKYCyJjtZiEDM3LVq7ShHVx8JzSzvDn", "bcd")
	fmt.Println("bcd Address: ", bcdWallet.Address)
	item1 := hpywallet.OutPutItem{
		TxHash:   "72e9c6e7490ba50cded735ec0f1d48d324c9509d2ef6a237b45f351638a1a28a",
		Value:    1000000,
		Vout:     0,
		Pkscript: "76a9140f343240a5c2b4c532be07835716b2b46483f31188ac",
	}
	outputs := []hpywallet.OutPutItem{item1}

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		fmt.Println("outputs err: ", err.Error())

	}

	signInput := &hpywallet.SignInput{
		Coin:       "bcd",
		Symbol:     "bcd",
		PrivateKey: bcdWallet.PrivateKey,
		SrcAddr:    bcdWallet.Address,
		DestAddr:   "1Er5oNTJFiweAm7HSLhMikGzh6ZhudQeFP",
		Fee:        1000,
		Amount:     999000,
		Change:     0,
		Inputs:     jsonInputs,
	}
	tranferResult := hpywallet.SignRawTransaction(signInput)
	if tranferResult.ResCode == 0 {
		fmt.Println("Transfer bcd Msg: ", tranferResult.ErrMsg)
	} else {
		fmt.Println("Transfer bcd RawTX: ", tranferResult.RawTX)
		// 0c0000000100000000000000000000000000000000000000000000000000000000000000018aa2a13816355fb437a2f62e9d50c924d3481d0fec35d7de0ca50b49e7c6e972000000006b483045022100a98cc99cdfc33e1a10d0b24dfb97b36291131ad2eb324822dd69bb4de5ec562b022019c556e86952027ac4480b9f48553dd72d309cd25f424d1210f2ec9760ea82f8012103bec0cc32c8e2117488e9262ce7c99f6d7a56301c3420ae585ece49ec035d9ed7ffffffff01583e0f00000000001976a91497e206129e477018a56b6266af7f1dacbc061fe388ac00000000
		// txid: 018d24ff27e442c64ac2e24bc7930db6828bc3ef4264d4ebdb26052132a6b9be
	}

}

func Test_TransferUSDT(t *testing.T) {
	wallet := hpywallet.GenerateWallet(mnemonic, "usdt")
	usdtWallet := hpywallet.ImportPrivateWIF(wallet.PrivateKey, "usdt")
	fmt.Println("usdt : ", wallet)
	fmt.Println("usdtWallet : ", usdtWallet)
	item1 := hpywallet.OutPutItem{
		TxHash:   "0f90f5e98cfbcbeb23e84e9275540d853aa062630ae8b2a30ac5f1308a30fe14",
		Value:    40172,
		Vout:     1,
		Pkscript: "76a914ccfbd451e79ed9dc55356179548f18fa95c500f288ac",
	}
	// item2 := hpywallet.OutPutItem{
	// 	TxHash:   "42d8fa535fd889abe51efe633882cf4399e1f39ea346fd98967b9af9da8f9754",
	// 	Value:    546,
	// 	Vout:     2,
	// 	Pkscript: "76a9149d5fb37c2ac97ec80bcacba917cb4fccfca9f1ea88ac",
	// }
	outputs := []hpywallet.OutPutItem{item1}
	// fmt.Println("ltc outputs: ", outputs)

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		//log.Fatal("Cannot encode to JSON ", err)
		fmt.Println("outputs err: ", err.Error())

	}
	// fmt.Println("ltc outputs: ", jsonInputs)

	signInput := &hpywallet.SignInput{
		Coin:       "usdt",
		Symbol:     "usdt",
		PrivateKey: wallet.PrivateKey,
		SrcAddr:    wallet.Address,
		DestAddr:   "1FQ6Sv1yDi6AnC8n8BQnLLNAaQ1esGT5oN",
		Fee:        2000,
		Amount:     10000,
		Change:     40172 - 1000 - 2000,
		Inputs:     jsonInputs,
	}
	tranferResult := hpywallet.SignRawTransaction(signInput)
	if tranferResult.ResCode == 0 {
		fmt.Println("Transfer USDT Msg: ", tranferResult.ErrMsg)
	} else {
		fmt.Println("Transfer USDT RawTX: ", tranferResult.RawTX)
	}

}

func Test_TransferWICC(t *testing.T) {
	//wallet := hpywallet.GenerateWallet(mnemonic, "wicc")
	//privateKey := wallet.PrivateKey
	privateKey := "Y66sU3zeS8aPF45tdSw5XUyPCmGTcJCd8gR3rPmynCLm2hDJYKkg"
	wallet := hpywallet.ImportPrivateWIF(privateKey, "wicc")
	fmt.Println("wicc wallet: ", wallet)
	fmt.Println("wicc address: ", wallet.Address)
	fmt.Println("wicc private key: ", wallet.PrivateKey)
	fmt.Println("wicc public key: ", wallet.PublicKey)

	signInput := &hpywallet.SignInput{
		Coin:       "wicc",
		Symbol:     "wicc",
		PrivateKey: wallet.PrivateKey,
		SrcAddr:    wallet.Address,
		DestAddr:   "wKf3QkwNa2t7cSL6MVrFU12eVMEFU1YDx8",
		Fee:        1000000,
		Amount:     2000000,
	}
	transferResult := hpywallet.SignRawTransaction(signInput)
	if transferResult.ResCode == 0 {
		fmt.Println("Transfer WICC Msg: ", transferResult.ErrMsg)
	} else {
		fmt.Println("Transfer WICC RawTX: ", transferResult.RawTX)
	}
}

func Test_TransferLTC(t *testing.T) {
	ltcWallet := hpywallet.GenerateWallet(mnemonic, "ltc")
	// fmt.Println("ltc Address: ", ltcWallet.Address)
	// fmt.Println("ltc Wif: ", ltcWallet.PrivateKey)
	// fmt.Println("ltc PublicKey: ", ltcWallet.PublicKey)
	item1 := hpywallet.OutPutItem{
		TxHash:   "921784b1e11fcbfe267a04b9e54a45e597710d0f9413e658813737c06f44a987",
		Value:    1200000,
		Vout:     1,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	item2 := hpywallet.OutPutItem{
		TxHash:   "62409f73492acc32aa7760423e56127f78f62745f07a68c4107ea69394dce485",
		Value:    400000,
		Vout:     1,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	item3 := hpywallet.OutPutItem{
		TxHash:   "b886364bf351ac4ae47e4eea8ab4e1039c4fd92fed04899bb360d0570a0f92a2",
		Value:    100000,
		Vout:     0,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	item4 := hpywallet.OutPutItem{
		TxHash:   "ca2695a680d6e7bea7dc5155f9cbe13f3b230a648dd48135ff5b090ae9c89194",
		Value:    100000,
		Vout:     1,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	item5 := hpywallet.OutPutItem{
		TxHash:   "761faedb198081b97478ae3bc6f85deefc43ab5d1e116fe2a83ff33bd9dbea14",
		Value:    15500000,
		Vout:     1,
		Pkscript: "76a914bc68c7efc2f672c3ea028e10ec321a9c68d5da7788ac",
	}
	outputs := []hpywallet.OutPutItem{item1, item2, item3, item4, item5}
	// fmt.Println("ltc outputs: ", outputs)

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		//log.Fatal("Cannot encode to JSON ", err)
		fmt.Println("outputs err: ", err.Error())

	}
	// fmt.Println("ltc outputs: ", jsonInputs)

	signInput := &hpywallet.SignInput{
		Coin:       "ltc",
		Symbol:     "ltc",
		PrivateKey: ltcWallet.PrivateKey,
		SrcAddr:    ltcWallet.Address,
		DestAddr:   "LezCFqxL7NHKzAKiuwzE12QyNsF6rTe3rU",
		Fee:        200000,
		Amount:     100000,
		Change:     900000,
		Inputs:     jsonInputs,
	}
	tranferResult := hpywallet.SignRawTransaction(signInput)
	if tranferResult.ResCode == 0 {
		fmt.Println("Transfer ltc Msg: ", tranferResult.ErrMsg)
	} else {
		fmt.Println("Transfer ltc RawTX: ", tranferResult.RawTX)
	}
}

func Test_TransferADA(t *testing.T) {
	//mnemonic = "fit today convince mixture crumble design identify swim jelly charge vast expect below science stove"
	//wallet := hpywallet.GenerateMnemonicWallet(mnemonic, "qwert", "ada")
	wallet := hpywallet.GenerateWallet(mnemonic, "ada")
	fmt.Println("ada Address: ", wallet.Address)
	fmt.Println("ada privateKey: ", wallet.PrivateKey)
	fmt.Println("ada Seed: ", wallet.Seed)

	item1 := hpywallet.OutPutItem{
		TxHash: "b19ec911686d6e34a064b26848c4dd487948d82e4f3e9b2d85662744eb39a552",
		Value:  689808512,
		Vout:   1,
	}
	outputs := []hpywallet.OutPutItem{item1}

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		fmt.Println("outputs err: ", err.Error())
	}

	signInput := &hpywallet.SignInput{
		Coin:       "ada",
		Symbol:     "ada",
		PrivateKey: wallet.PrivateKey,
		SrcAddr:    wallet.Address,
		DestAddr:   "37btjrVyb4KBEC2HAhrGPMDncPxe4noSnDeQsAtqpZecEHZ5bZKAcH8zzffE9w9EipF7BzU2WpECvqRFF7FJZDA6NyMQ5Daxd3usTbTLNLcUapnrPy",
		Fee:        200000,
		Amount:     1000000,
		Change:     688608512,
		Inputs:     jsonInputs,
	}
	transferResult := hpywallet.SignRawTransaction(signInput)
	if transferResult.ResCode == 0 {
		fmt.Println("Transfer ada Msg: ", transferResult.ErrMsg)
	} else {
		fmt.Println("Transfer ada Hash: ", transferResult.TxHash)
		fmt.Println("Transfer ada RawTX: ", transferResult.RawTX)
	}

}

func Test_importBCH(t *testing.T) {
	bchWallet := hpywallet.GenerateWallet(mnemonic, "btc")
	fmt.Println("bch Address: ", bchWallet.Address)
	fmt.Println("bch Wif: ", bchWallet.PrivateKey)

	fmt.Println("bch PublicKey: ", bchWallet.PublicKey)

	importWallet := hpywallet.ImportPrivateWIF(bchWallet.PrivateKey, "bch")
	if importWallet.ResCode == 0 {
		fmt.Println("import bch Msg: ", importWallet.ErrMsg)
	} else {
		fmt.Println("import bch Address: ", importWallet.Address)
	}

}

func Test_importETH(t *testing.T) {
	ethWallet := hpywallet.GenerateWallet(testmnemonic, "eth")
	fmt.Println("eth Address: ", ethWallet.Address)

	destWallet := hpywallet.GenerateWallet(mnemonic, "eth")
	fmt.Println("收方地址：", destWallet.Address)

	signInput := &hpywallet.SignInput{
		Coin:       "eth",
		Symbol:     "eth",
		PrivateKey: ethWallet.PrivateKey,
		SrcAddr:    ethWallet.Address,
		DestAddr:   destWallet.Address,
		Amount:     10000000000000, //
		GasLimit:   25200,
		GasPrice:   20000000000,
	}

	tranferResult := hpywallet.SignRawTransaction(signInput)
	fmt.Println("ETH rawTx ", tranferResult.RawTX)

}

func Test_transferVSYS(t *testing.T) {
	mnemonic := "cloud firm joy outside story divide broom vivid fatigue tag vast bless arrange detect inch"
	recAddress := "ARDmUHRPbd55E7bfqPDdYbV8kxRB3Y29c7Z"
	amount := int64(4890000000)

	acc := vsys.InitAccount(vsys.Mainnet)
	acc.BuildFromSeed(mnemonic, 0)
	fmt.Println("地址：", acc.Address()) // ARFNxs5uyfa9bCbRisqV1xYKNUis6NBhmCR

	signInput := &hpywallet.SignInput{
		Coin:       "vsys",
		Symbol:     "vsys",
		PrivateKey: acc.PrivateKey(),
		SrcAddr:    acc.Address(),
		DestAddr:   recAddress,
		Amount:     amount, //
	}

	tranferResult := hpywallet.SignRawTransaction(signInput)
	fmt.Println("Vsys rawTx ", tranferResult.RawTX)
	// api文档: http://test.v.systems:9922/api-docs/index.html#/vsys/broadcastPayment
	// 测试网广播接口：http://test.v.systems:9922/vsys/broadcast/payment
	// 主网广播接口：https://wallet.v.systems/api/vsys/broadcast/payment
	// 测试网广播并生成交易: https://testexplorer.v.systems/transactions/4RPztJy8p2NpyoUJdiYgcsbTPUdQHKho1z4VmQgwmWBw
	// 主网广播并生成交易: https://explorer.v.systems/transactions/6vpHuh9vrDeBKVtSrVbGQcsTmfvzd8D3VSUd1CtYQRR9
}

func Test_transferVSYSContract(t *testing.T) {
	recAddress := "AR3zYe4jBVsk9HNvEJZ85VCAT9CgDQhFLmg"
	contractId := "CC8Jx8aLkKVQmzuHBWNnhCSkn1GBLcjZ32k" // ipx的contractId 接口文档：http://test.v.systems:9922/api-docs/index.html#/contract/tokenInfo
	amount := int64(1000000000)
	acc := vsys.InitAccount(vsys.Mainnet)
	acc.BuildFromSeed("cloud firm joy outside story divide broom vivid fatigue tag vast bless arrange detect inch", 0)
	fmt.Println("地址：", acc.Address()) // ARFNxs5uyfa9bCbRisqV1xYKNUis6NBhmCR

	signInput := &hpywallet.SignInput{
		Coin:         "vsys",
		Symbol:       "ipx",
		PrivateKey:   acc.PrivateKey(),
		SrcAddr:      acc.Address(),
		DestAddr:     recAddress,
		Amount:       amount,
		ContractAddr: contractId,
	}

	tranferResult := hpywallet.SignRawTransaction(signInput)
	fmt.Println("Vsys rawTx ", tranferResult.RawTX)
	// api文档: http://test.v.systems:9922/api-docs/index.html#/contract/signedExecute
	// 测试网广播接口：http://test.v.systems:9922/contract/broadcast/execute
	// 主网广播接口：https://wallet.v.systems/api/contract/broadcast/execute
	// 测试网广播并生成交易：https://testexplorer.v.systems/transactions/HLAG5BeyWwD82offGuRPdGtXesitYr7PYHSkdnn9AwK4
	// 主网广播并生成交易：https://explorer.v.systems/transactions/7ue8DUvnebYNM7wxcLtC7DT9cLn1aGVhbobzTyNoH4JC
}

func Test_importETC(t *testing.T) {
	etcWallet := hpywallet.GenerateWallet(testmnemonic, "etc") // 0x95573e2ffD61A6c5e08Fc321A7e8754f41b6C471
	fmt.Println("etc Wallet: ", etcWallet)
	// importWallet := hpywallet.ImportPrivateWIF(etcWallet.PrivateKey, "etc")
	// fmt.Println("导入钱包：", importWallet)
	destWallet := hpywallet.GenerateWallet(mnemonic, "etc")
	fmt.Println("To Wallet：", destWallet)

	param := hpywallet.ETCParams{Nonce: 7} // Nonce 接口获取，每次取值为上次交易的 nonce+1
	jsonParam, err := json.Marshal(param)
	if err != nil {
		fmt.Println("jsonParam err: ", err.Error())
	}

	signInput := &hpywallet.SignInput{
		Coin:       "etc",
		Symbol:     "etc",
		PrivateKey: etcWallet.PrivateKey,
		SrcAddr:    etcWallet.Address,
		DestAddr:   destWallet.Address,
		Amount:     1e18 * 0.007, //
		GasLimit:   25200,
		GasPrice:   2e10, // 1 0025 2000 0000 0000
		Params:     jsonParam,
	}

	tranferResult := hpywallet.SignRawTransaction(signInput)
	fmt.Println("ETC 构造结果：", tranferResult)

	// https://gastracker.io/addr/0x95573e2ffd61a6c5e08fc321a7e8754f41b6c471
	// https://etherscan.io/pushTx?%3flang=zh-CN

	//ETC交易hash 等待结果 0xa97c24efe6abc6efc0e8ca7fbe44df830899420764b2180d1de3d0e1975eec0e
	// f86709808262709495573e2ffd61a6c5e08fc321a7e8754f41b6c471872386f26fc1000080819ea02ee84a3079b5ecbc97c30f8e2dbe06e4df4607ad42bfa0bb21e769a3a50eb775a063ec17d76f60198a3086f20b9d42e3163dda913bb2cd226d7912073d3c5a814f

}

func Test_importQTUM(t *testing.T) {
	hcWallet := hpywallet.GenerateWallet(mnemonic, "qtum")
	fmt.Println("qtum Address: ", hcWallet.Address)
	fmt.Println("qtum Wif: ", hcWallet.PrivateKey)

	fmt.Println("qtum PublicKey: ", hcWallet.PublicKey)

	importWallet := hpywallet.ImportPrivateWIF(hcWallet.PrivateKey, "qtum")
	if importWallet.ResCode == 0 {
		fmt.Println("import qtum Msg: ", importWallet.ErrMsg)
	} else {
		fmt.Println("import qtum Address: ", importWallet.Address)
	}

}

func Test_importLTC(t *testing.T) {
	ltcWallet := hpywallet.GenerateWallet(mnemonic, "ltc")
	fmt.Println("ltc Address: ", ltcWallet.Address)
	fmt.Println("ltc Wif: ", ltcWallet.PrivateKey)

	fmt.Println("ltc PublicKey: ", ltcWallet.PublicKey)

	importWallet := hpywallet.ImportPrivateWIF(ltcWallet.PrivateKey, "ltc")
	if importWallet.ResCode == 0 {
		fmt.Println("import ltc Msg: ", importWallet.ErrMsg)
	} else {
		fmt.Println("import ltc Address: ", importWallet.Address)
	}

}

func Test_MultiSigBTC(t *testing.T) {
	// 生成助记词
	mnemonicA := "ecology sponsor icon fox spot slot mirror business crumble laundry angle theory"
	mnemonicB := "ketchup salon shock alone track glimpse absent category grid nut fragile inside"
	mnemonicC := "cigar swarm silk imitate bind still blood carry shuffle endless impact onion"
	mnemonicRec := "pizza employ digital prison begin alcohol mail brother almost weekend march isolate"

	// 生成三个普通地址钱包
	btcWalletA := hpywallet.GenerateWallet(mnemonicA, "btc")
	fmt.Println("btc AddressA: ", btcWalletA.Address) // mx7CHaDSwPyFXMRQGGbtusH211UZ2fo8Fo
	pubA, _ := hex.DecodeString(btcWalletA.PublicKey)

	btcWalletB := hpywallet.GenerateWallet(mnemonicB, "btc")
	fmt.Println("btc AddressB: ", btcWalletB.Address) // mnJ8xAtJ8Hvqh8gZhNZsmH9KNfpi75h2c9
	pubB, _ := hex.DecodeString(btcWalletB.PublicKey)

	btcWalletC := hpywallet.GenerateWallet(mnemonicC, "btc")
	fmt.Println("btc AddressC: ", btcWalletC.Address) // mqjU3A2Tf5DYza9afgQq6Z6ttM3ZdtN6bq
	pubC, _ := hex.DecodeString(btcWalletC.PublicKey)

	// 生成接受地址钱包
	btcWalletRec := hpywallet.GenerateWallet(mnemonicRec, "btc")
	fmt.Println("btc AddressRec: ", btcWalletRec.Address) // mqCmhcEWaJGwg9EBUJQ7K2ihttAG4PD8t8

	// step 1
	//填充成为二维数组，获取多重签名地址
	required := byte(2)
	segwit := true
	addressPrefix := btcTransaction.AddressPrefix{P2PKHPrefix: []byte{0x6f}, P2WPKHPrefix: []byte{0xc4}, Bech32Prefix: "tb"}
	address, redeem, err := btcTransaction.CreateMultiSig(required, [][]byte{pubA, pubB, pubC}, segwit, addressPrefix)
	if err != nil {
		t.Error("创建多签地址失败！", err)
	} else {
		fmt.Println("多重签名地址为：")
		fmt.Println(address) // 2N4Gh4KCrUake53WPPJzjNHvwgSxhrgcJGT
		fmt.Println("赎回脚本为：")
		fmt.Println(redeem)
	}

	//step 2
	// 向该多重签名地址转入一定数额的比特币
	//txid 2dcd710076077def249485d5f938d8815daa0084c0ef9097b48cb11a794f1e94
	//vout 0
	//amount 0.01130647 BTC / 1130647 satoshi
	//ScriptPubkey a91478f0656432998d00d4224d15db4b43a8cb25e41787

	//step3
	// 构建空交易单
	fee := 10000                                  // 手续费
	totalAmount := 1130647                        // 总金额
	outAmount := 100647                           // 转账金额
	changeAmount := totalAmount - fee - outAmount // 找零
	in := btcTransaction.Vin{TxID: "2dcd710076077def249485d5f938d8815daa0084c0ef9097b48cb11a794f1e94", Vout: uint32(0)}
	out := btcTransaction.Vout{Address: btcWalletRec.Address, Amount: uint64(outAmount)}
	change := btcTransaction.Vout{Address: address, Amount: uint64(changeAmount)}
	//锁定时间
	lockTime := uint32(0)
	//追加手续费支持
	replaceable := false
	emptyTrans, err := btcTransaction.CreateEmptyRawTransaction([]btcTransaction.Vin{in}, []btcTransaction.Vout{out, change}, lockTime, replaceable, addressPrefix)
	if err != nil {
		t.Error("构建空交易单失败", err)
	} else {
		fmt.Println("空交易单：")
		fmt.Println(emptyTrans)
	}

	// 构建交易单签名哈希
	inLock := "a91478f0656432998d00d4224d15db4b43a8cb25e41787"
	inRedeem := redeem
	inAmount := uint64(totalAmount)
	unlockData := btcTransaction.TxUnlock{LockScript: inLock, RedeemScript: inRedeem, Amount: inAmount, SigType: btcTransaction.SigHashAll}

	// 计算待签名交易单哈希
	transHash, err := btcTransaction.CreateRawTransactionHashForSig(emptyTrans, []btcTransaction.TxUnlock{unlockData}, segwit, addressPrefix)
	if err != nil {
		t.Error("创建待签交易单哈希失败", err)
	} else {
		for i, t := range transHash {
			fmt.Println("第", i+1, "个交易单哈希值为")
			fmt.Printf("%+v \n", t)
		}
	}

	//////签名哈希
	// 获取到的transHash数组只有一个元素，该哈希值是所有多签参与方的签名哈希
	// 根据required值，选择足够数量的签名方，发送哈希值

	// A的私钥
	wifA, _ := btcutil.DecodeWIF(btcWalletA.PrivateKey)
	priA := wifA.PrivKey.Serialize()

	// B的私钥
	wifB, _ := btcutil.DecodeWIF(btcWalletB.PrivateKey)
	priB := wifB.PrivKey.Serialize()

	// A 签名
	sigPubA, err := btcTransaction.SignRawTransactionHash(transHash[0].Hash, priA)
	if err != nil {
		t.Error("A签名失败", err)
	} else {
		fmt.Println("A的签名结果为")
		fmt.Println(hex.EncodeToString(sigPubA.Signature))
	}

	// B 签名
	sigPubB, err := btcTransaction.SignRawTransactionHash(transHash[0].Hash, priB)
	if err != nil {
		t.Error("B签名失败", err)
	} else {
		fmt.Println("B的签名结果为")
		fmt.Println(hex.EncodeToString(sigPubB.Signature))
	}

	// 接收到签名结果后，回填TxHash结构体数组
	transHash[0].Multi[0].SigPub = *sigPubA
	transHash[0].Multi[1].SigPub = *sigPubB

	// 合并交易单
	signedTrans, err := btcTransaction.InsertSignatureIntoEmptyTransaction(emptyTrans, transHash, []btcTransaction.TxUnlock{unlockData}, segwit)
	if err != nil {
		t.Error("插入交易单失败", err)
	} else {
		fmt.Println("合并之后的交易单")
		fmt.Println(signedTrans)
	}

	// 验证交易单
	pass := btcTransaction.VerifyRawTransaction(signedTrans, []btcTransaction.TxUnlock{unlockData}, segwit, addressPrefix)
	if pass {
		fmt.Println("验证通过!")
	} else {
		t.Error("验证失败!")
	}

	// 广播交易 https://live.blockcypher.com/btc-testnet/pushtx/
	// https://live.blockcypher.com/btc-testnet/tx/61899e695457eae66e702b9308916cb6c1f4de212b46b7eb3b8fb4a41605dc2c/
}

func Test_MultiSigBCH(t *testing.T) {
	// 生成助记词
	mnemonicA := "ecology sponsor icon fox spot slot mirror business crumble laundry angle theory"
	mnemonicB := "ketchup salon shock alone track glimpse absent category grid nut fragile inside"
	mnemonicC := "cigar swarm silk imitate bind still blood carry shuffle endless impact onion"
	mnemonicRec := "pizza employ digital prison begin alcohol mail brother almost weekend march isolate"

	// 生成三个普通地址钱包
	bchWalletA := hpywallet.GenerateWallet(mnemonicA, "bch")
	fmt.Println("bch AddressA: ", bchWalletA.Address) // mttrbo7kZ68B1DzCHA9EEW235QyrEYfPuB

	//os.Exit(1)
	//btcWalletA := hpywallet.GenerateWallet(mnemonicA, "btc")
	//fmt.Println("btc AddressA: ", btcWalletA.Address) // 1ENuJk2mk4gvE7WaZbArQaoiDRP9LfmxrL
	//
	//os.Exit(2)
	pubA, _ := hex.DecodeString(bchWalletA.PublicKey)

	bchWalletB := hpywallet.GenerateWallet(mnemonicB, "bch")
	fmt.Println("bch AddressB: ", bchWalletB.Address) // mo7pGxZ9VFtCg1u8VPEzy2HwoBjC5jZyJE
	pubB, _ := hex.DecodeString(bchWalletB.PublicKey)

	bchWalletC := hpywallet.GenerateWallet(mnemonicC, "bch")
	fmt.Println("bch AddressC: ", bchWalletC.Address) // mq8EtmRWqBPmrhyLu8Vw9X8k4PTJXA6JsL
	pubC, _ := hex.DecodeString(bchWalletC.PublicKey)

	// 生成接受地址钱包
	bchWalletRec := hpywallet.GenerateWallet(mnemonicRec, "bch")
	fmt.Println("bch AddressRec: ", bchWalletRec.Address) //mgH28rvvy92ereBTsePDjyRian5xEBZ3nx

	// step 1
	//填充成为二维数组，获取多重签名地址
	required := byte(2)
	segwit := true
	addressPrefix := bitcoincashTransaction.AddressPrefix{P2PKHPrefix: []byte{0x6f}, P2WPKHPrefix: []byte{0xc4}, Bech32Prefix: "bc"}
	address, redeem, err := bitcoincashTransaction.CreateMultiSig(required, [][]byte{pubA, pubB, pubC}, segwit, addressPrefix)
	if err != nil {
		t.Error("创建多签地址失败！", err)
	} else {
		fmt.Println("多重签名地址为：")
		fmt.Println(address) // 2N2VfKzSzLxRxu1aXML5grdEQKeMhhinFf3  => bchtest:ppjhgxav4a43987jq5d9rkxm0sydewkssq36dggeh7
		fmt.Println("赎回脚本为：")
		fmt.Println(redeem)
	}

	//step 2
	// 向该多重签名地址转入一定数额的比特币
	//txid 8c1685fe00842d7f16b9631c8a9e62a7b38636bb091030b4a8308cfb4bbbf60b
	//vout 1
	//amount 0.1 BCH
	//ScriptPubkey a91465741bacaf6b129fd2051a51d8db7c08dcbad08087

	//step3
	// 构建空交易单
	in := bitcoincashTransaction.Vin{TxID: "8c1685fe00842d7f16b9631c8a9e62a7b38636bb091030b4a8308cfb4bbbf60b", Vout: uint32(1)}
	out := bitcoincashTransaction.Vout{Address: bchWalletRec.Address, Amount: uint64(9800000)}

	//锁定时间
	lockTime := uint32(0)

	//追加手续费支持
	replaceable := false

	emptyTrans, err := bitcoincashTransaction.CreateEmptyRawTransaction([]bitcoincashTransaction.Vin{in}, []bitcoincashTransaction.Vout{out}, lockTime, replaceable, addressPrefix)
	if err != nil {
		t.Error("构建空交易单失败", err)
	} else {
		fmt.Println("空交易单：")
		fmt.Println(emptyTrans)
	}

	// 构建交易单签名哈希
	inLock := "a91465741bacaf6b129fd2051a51d8db7c08dcbad08087"
	inRedeem := redeem
	inAmount := uint64(10000000)

	unlockData := bitcoincashTransaction.TxUnlock{LockScript: inLock, RedeemScript: inRedeem, Amount: inAmount, SigType: bitcoincashTransaction.SigHashAll}

	/////////计算待签名交易单哈希
	transHash, err := bitcoincashTransaction.CreateRawTransactionHashForSig(emptyTrans, []bitcoincashTransaction.TxUnlock{unlockData}, segwit, addressPrefix)
	if err != nil {
		t.Error("创建待签交易单哈希失败", err)
	} else {
		for i, t := range transHash {
			fmt.Println("第", i+1, "个交易单哈希值为")
			fmt.Printf("%+v \n", t)
		}
	}

	//////签名哈希
	// 获取到的transHash数组只有一个元素，该哈希值是所有多签参与方的签名哈希
	// 根据required值，选择足够数量的签名方，发送哈希值

	// A的私钥
	wifA, _ := btcutil.DecodeWIF(bchWalletA.PrivateKey)
	priA := wifA.PrivKey.Serialize()

	// B的私钥
	wifB, _ := btcutil.DecodeWIF(bchWalletB.PrivateKey)
	priB := wifB.PrivKey.Serialize()

	// A 签名
	sigPubA, err := bitcoincashTransaction.SignRawTransactionHash(transHash[0].Hash, priA)
	if err != nil {
		t.Error("A签名失败", err)
	} else {
		fmt.Println("A的签名结果为")
		fmt.Println(hex.EncodeToString(sigPubA.Signature))
	}

	// B 签名
	sigPubB, err := bitcoincashTransaction.SignRawTransactionHash(transHash[0].Hash, priB)
	if err != nil {
		t.Error("B签名失败", err)
	} else {
		fmt.Println("B的签名结果为")
		fmt.Println(hex.EncodeToString(sigPubB.Signature))
	}

	// 接收到签名结果后，回填TxHash结构体数组
	transHash[0].Multi[0].SigPub = *sigPubA
	transHash[0].Multi[1].SigPub = *sigPubB

	// 合并交易单
	signedTrans, err := bitcoincashTransaction.InsertSignatureIntoEmptyTransaction(emptyTrans, transHash, []bitcoincashTransaction.TxUnlock{unlockData}, segwit)
	if err != nil {
		t.Error("插入交易单失败", err)
	} else {
		fmt.Println("合并之后的交易单")
		fmt.Println(signedTrans)
	}

	// 验证交易单
	pass := bitcoincashTransaction.VerifyRawTransaction(signedTrans, []bitcoincashTransaction.TxUnlock{unlockData}, segwit, addressPrefix)
	if pass {
		fmt.Println("验证通过!")
	} else {
		t.Error("验证失败!")
	}
}

/**
 * USDT手续费代付
 * 原理：增加btc input
 * 交易广播后：https://blockexplorer.one/omni/testnet/tx/0f54114a7158357182beb443b0f137947dcba8e9a8353502a9df499fd9daca8d
 */
func Test_OmniUsdt(t *testing.T) {
	mnemonicA := "ecology sponsor icon fox spot slot mirror business crumble laundry angle theory"
	btcWallet := hpywallet.GenerateWallet(mnemonicA, "usdt")
	fmt.Println("btc address : ", btcWallet.Address) // mx7CHaDSwPyFXMRQGGbtusH211UZ2fo8Fo
	btcAddress := btcWallet.Address

	usdtWallet := hpywallet.ImportPrivateWIF("cViSUdw7nbD85zZ4FbDVVx1FicV5kLT9k8DCVyCAy2xPCF3guVoC", "usdt")
	fmt.Println("usdt address : ", usdtWallet.Address) // mv4MrxiBYQFHPgNuGL2TrY4mnvSsZzybB6
	toUsdtAddress := "mxAkuEoSG9UN2NJ7yxwZ6DnL26sry22g65"

	// 第一个输入 omni
	usdtIn := omniTransaction.Vin{"0730abe4d3380c1f1a2949034ff524de57376035f31656398005958cbfcb5314", uint32(0)} // mv4MrxiBYQFHPgNuGL2TrY4mnvSsZzybB6
	// 第二个输入 fee
	feeIn := omniTransaction.Vin{"c0db42974d205255c1af2297353073422188a582c1ebbcc6a58a1a80ee75999b", uint32(0)} // mx7CHaDSwPyFXMRQGGbtusH211UZ2fo8Fo

	// 目标地址与数额
	// 向 mxAkuEoSG9UN2NJ7yxwZ6DnL26sry22g65 发送
	// out 单位为聪
	change := omniTransaction.Vout{btcAddress, uint64(9000)} // btc
	to := omniTransaction.Vout{toUsdtAddress, uint64(1000)}  // usdt mxAkuEoSG9UN2NJ7yxwZ6DnL26sry22g65 地址私钥：cPRzxcdJqWeF4jjsCDu1aaZLiPpwfKcWbU5kxzUSTRZNW6g38ApU

	omniDetail := omniTransaction.OmniStruct{omniTransaction.SimpleSend, uint32(2), 10000000, 0, "", toUsdtAddress}
	//锁定时间
	lockTime := uint32(0)

	//追加手续费支持
	replaceable := false

	/////////构建空交易单
	emptyTrans, err := omniTransaction.CreateEmptyRawTransaction([]omniTransaction.Vin{usdtIn, feeIn}, []omniTransaction.Vout{change, to}, omniDetail, lockTime, replaceable, omniTransaction.BTCTestnetAddressPrefix)

	if err != nil {
		t.Error("构建空交易单失败")
	} else {
		fmt.Println("空交易单：")
		fmt.Println(emptyTrans)
	}

	// 获取usdtIn 和 feeIn 的锁定脚本
	// 填充TxUnlock结构体
	usdtInLock := "76a9149f82542e98780434e43a5e401f34f500fcca031f88ac"
	feeInLock := "76a914b5fbd5be687da46adf1495b9608aafdea222f71788ac"

	//针对此类指向公钥哈希地址的UTXO，此处仅需要锁定脚本即可计算待签交易单
	usdtUnlockData := omniTransaction.TxUnlock{usdtInLock, "", uint64(0), omniTransaction.SigHashAll}
	feeUnlockData := omniTransaction.TxUnlock{feeInLock, "", uint64(0), omniTransaction.SigHashAll}

	////////构建用于签名的交易单哈希
	transHash, err := omniTransaction.CreateRawTransactionHashForSig(emptyTrans, []omniTransaction.TxUnlock{usdtUnlockData, feeUnlockData}, omniTransaction.BTCTestnetAddressPrefix)
	if err != nil {
		t.Error("获取待签名交易单哈希失败")
	} else {
		for i, t := range transHash {
			fmt.Println("第", i+1, "个交易单哈希值为")
			fmt.Println(t)
		}
	}

	// 获取私钥
	// usdtInLock address mv4MrxiBYQFHPgNuGL2TrY4mnvSsZzybB6
	usdtWif, _ := btcutil.DecodeWIF(usdtWallet.PrivateKey)
	usdtPri := usdtWif.PrivKey.Serialize()

	// feeIn address mx7CHaDSwPyFXMRQGGbtusH211UZ2fo8Fo
	btcWif, _ := btcutil.DecodeWIF(btcWallet.PrivateKey)
	btcPri := btcWif.PrivKey.Serialize()

	// 客户端对第一条hash进行签名
	sigPub1, err := omniTransaction.SignRawTransactionHash(transHash[0].Hash, usdtPri)
	if err != nil {
		t.Error("第一条hash签名失败!")
	} else {
		fmt.Println("第一条hash签名结果")
		fmt.Println(hex.EncodeToString(sigPub1.Signature))
		fmt.Println("对应公钥")
		fmt.Println(hex.EncodeToString(sigPub1.Pubkey))
	}

	// 客户端对第二条hash进行签名
	sigPub2, err := omniTransaction.SignRawTransactionHash(transHash[1].Hash, btcPri)
	if err != nil {
		t.Error("第二条hash签名失败!")
	} else {
		fmt.Println("第二条hash签名结果")
		fmt.Println(hex.EncodeToString(sigPub2.Signature))
		fmt.Println("对应公钥")
		fmt.Println(hex.EncodeToString(sigPub2.Pubkey))
	}

	// 服务器收到签名结果后，使用签名结果填充TxHash结构体
	transHash[0].Normal.SigPub = *sigPub1
	transHash[1].Normal.SigPub = *sigPub2

	//回填后，将签名插入空交易单
	signedTrans, err := omniTransaction.InsertSignatureIntoEmptyTransaction(emptyTrans, transHash, []omniTransaction.TxUnlock{usdtUnlockData, feeUnlockData})
	if err != nil {
		t.Error("插入交易单失败")
	} else {
		fmt.Println("合并之后的交易单")
		fmt.Println(signedTrans)
	}

	// 验证交易单
	pass := omniTransaction.VerifyRawTransaction(signedTrans, []omniTransaction.TxUnlock{usdtUnlockData, feeUnlockData}, omniTransaction.BTCMainnetAddressPrefix)
	if pass {
		fmt.Println("验证通过!")
	} else {
		t.Error("验证失败!")
	}

	// 广播交易: https://blockexplorer.one/omni/testnet/tx/0f54114a7158357182beb443b0f137947dcba8e9a8353502a9df499fd9daca8d
}

func Test_transferHX(t *testing.T) {
	mnemonic := "embody balcony whisper arctic elephant method grace essay process magic trumpet sport"
	recAddress := "HXNdkvn13x3rB4cVEombnSx6654ZoMXNeeL4"
	amount := int64(1000)
	fee := int64(200)
	memo := "hello world"

	wallet, err := keypair.GenerateKeyPair(mnemonic)
	if err != nil {
		t.Error(err)
	}
	address, _ := wallet.PrivateKey.PublicKey().ToAddress()
	fmt.Println("地址：", address.String()) // ARFNxs5uyfa9bCbRisqV1xYKNUis6NBhmCR

	signInput := &hpywallet.SignInput{
		Coin:       "HX",
		Symbol:     "HX",
		PrivateKey: wallet.PrivateKey.ToWIF(),
		SrcAddr:    address.String(),
		DestAddr:   recAddress,
		Amount:     amount,
		Fee:        fee,
		Memo:       memo,
	}

	tranferResult := hpywallet.SignRawTransaction(signInput)
	fmt.Println("HX rawTx ", tranferResult.RawTX)
	// 广播地址：http://nodeapi.hxlab.org/
}

func TestCreateTRX(t *testing.T) {
	trx := hpywallet.GenerateWallet(mnemonic, "trx")
	// importTrx := hpywallet.ImportWIF("6c5e3afd3d6c0394dc922d9bcaf98fd9c972aa226948b44e14a7e4b0566c69ca", "trx")
	destTrx := hpywallet.GenerateWallet(ontmnemonic, "trx")
	// fmt.Println("TRX:", importTrx.Address, trx.PrivateKey)
	fmt.Println("助记词解析结果：\n", "地址："+trx.Address, "\n", "私钥：", trx.PrivateKey, "\n", "公钥：", trx.PublicKey)
	fmt.Println("********************* *************************")
	// wif := hpywallet.ImportWIF(trx.PrivateKey, "trx")
	// fmt.Println("私钥解析结果", "\n地址："+wif.Address, "\n私钥："+wif.PrivateKey, "\n公钥："+wif.PublicKey)

	signInput := &hpywallet.SignInput{
		Coin:         "trx",
		Symbol:       "trx",
		PrivateKey:   "9cbbc1a83de7e5f02c598cc398ea6f09fd3f015d1b87b3057f764836a0ceee0b",
		SrcAddr:      "TJB6FQwkWw399q2UoyNCN9iqJv1zqiMwUR",
		DestAddr:     destTrx.Address,
		ContractAddr: "",
		Type:         "",
		Amount:       1000000,
		LargeAmount:  "2000000",
		Change:       0,
	}
	// signInput := &hpywallet.SignInput{
	// 	Coin:         "trx",
	// 	Symbol:       "COSMO",
	// 	PrivateKey:   "9cbbc1a83de7e5f02c598cc398ea6f09fd3f015d1b87b3057f764836a0ceee0b",
	// 	SrcAddr:      "TJB6FQwkWw399q2UoyNCN9iqJv1zqiMwUR",
	// 	DestAddr:     destTrx.Address,
	// 	ContractAddr: "1002636",
	// 	Type:         "trc10",
	// 	Amount:       1000000,
	// 	Change:       0,
	// }
	result := hpywallet.SignRawTransaction(signInput)
	if result.ResCode == 0 {
		fmt.Println("构造错误：", result.ErrMsg)
	} else {
		fmt.Println("构造成功： ", result.RawTX)
	}
	fmt.Println("result = ", result)

}
