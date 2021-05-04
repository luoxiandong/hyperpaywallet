package hdwallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aBitcoinDiamond/btcec"
	"github.com/aBitcoinDiamond/chaincfg"
	"github.com/aBitcoinDiamond/chaincfg/chainhash"
	"github.com/aBitcoinDiamond/txscript"
	"github.com/aBitcoinDiamond/wire"
	"github.com/liyaojian/btcutil"
)

var BCDParams = chaincfg.MainNetParams //

func init() {
	coins[BCD] = newBCD
}

type bcd struct {
	*btc
}

func newBCD(key *Key) Wallet {
	token := newBTC(key).(*btc)
	token.name = "Bitcoin Diamond"
	token.symbol = "BCD"
	//token.key.opt.Params = &BCHParams

	return &bcd{btc: token}
}

func (c *bcd) SignRawTransaction(signIn *SignInput) (*SignResult, error) {
	var vins []OutPutItem
	var vouts []OutPutItem
	json.Unmarshal(signIn.Inputs, &vins)
	fmt.Println("vins : ", vins)
	mtx := wire.NewMsgTx(wire.TxVersion)

	dest_addr, err := btcutil.DecodeAddress(signIn.DestAddr, &BCDParams)
	if err != nil {
		return nil, err
	}
	fmt.Println("dest_addr : ", dest_addr)
	dest_pkScript, err := txscript.PayToAddrScript(dest_addr)
	if err != nil {
		return nil, err
	}

	output := wire.NewTxOut(signIn.Amount, dest_pkScript)
	mtx.AddTxOut(output)
	// Add all outputs as inputs
	var spendValue int64 = 0
	for _, input := range vins {
		txHash, err := chainhash.NewHashFromStr(input.TxHash)
		if err != nil {
			return nil, fmt.Errorf("txid error")
		}
		prevOut := wire.NewOutPoint(txHash, input.Vout)
		txIn := wire.NewTxIn(prevOut, []byte{}, [][]byte{})
		mtx.AddTxIn(txIn)
		vouts = append(vouts, input)
		spendValue = spendValue + input.Value
		if spendValue >= signIn.Change+signIn.Amount+signIn.Fee {
			break
		}
	}

	addrSrc, err := btcutil.DecodeAddress(signIn.SrcAddr, &BCDParams)
	if err != nil {
		return nil, err
	}

	if signIn.Change > 0 {
		pkScriptSrc, err := txscript.PayToAddrScript(addrSrc)
		if err != nil {
			return nil, err
		}

		output := wire.NewTxOut(signIn.Change, pkScriptSrc)
		mtx.AddTxOut(output)
	}
	a, err := btcutil.DecodeWIF(signIn.PrivateKey)

	for i, input := range vouts {
		pk, _ := hex.DecodeString(input.Pkscript)
		if err != nil {
			return nil, err
		}
		//sigScript, err := txscript.SignatureScript(mtx, i, input.Value, pk, txscript.SigHashAll, a.PrivKey, true)
		sigScript, err := txscript.SignatureScript(mtx, i, pk, txscript.SigHashAll, (*btcec.PrivateKey)(a.PrivKey), true)

		if err != nil {
			return nil, err
		}
		mtx.TxIn[i].SignatureScript = sigScript
	}

	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, mtx.SerializeSize()))
	if err := mtx.Serialize(buf); err != nil {
		return nil, err
	}
	txHex := hex.EncodeToString(buf.Bytes())
	fmt.Println("txHex :", txHex)
	return &SignResult{
		Res:   1,
		RawTX: txHex,
	}, nil
}

func (c *bcd) GetWalletAccountFromWif() (*WalletAccount, error) {
	wif := c.GetKey().Wif
	if len(wif) > 0 {
		btcWif, err := btcutil.DecodeWIF(wif)
		if err != nil {
			fmt.Println("Wif err : ", err.Error())
			return nil, err
		}
		isBcd := btcWif.IsForNet(&BCDParams)
		if isBcd == false {
			return nil, errors.New("key type error")
		}
		pk := btcWif.SerializePubKey()
		fmt.Println("pk : ", hex.EncodeToString(pk))
		address, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pk), &BCDParams)
		if err != nil {
			fmt.Println("Wif err : ", err.Error())
			return nil, err
		}
		btcAddress := address.EncodeAddress()

		return &WalletAccount{
			Res:        1,
			PrivateKey: wif,
			PublicKey:  hex.EncodeToString(pk),
			Address:    btcAddress,
		}, nil
	}
	return &WalletAccount{
		Res: 0,
	}, nil
}
