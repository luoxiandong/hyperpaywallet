package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/aBitcoinDiamond/chaincfg"
	"github.com/aBitcoinDiamond/chaincfg/chainhash"
	"github.com/aBitcoinDiamond/txscript"
	"github.com/aBitcoinDiamond/wire"
	"github.com/liyaojian/btcutil"
)

func main() {
	mtx := wire.NewMsgTx(wire.TxVersion)

	dest_addr, err := btcutil.DecodeAddress("1Er5oNTJFiweAm7HSLhMikGzh6ZhudQeFP", &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("dest_addr : ", dest_addr)
	dest_pkScript, err := txscript.PayToAddrScript(dest_addr)
	if err != nil {
		panic(err)
	}

	output := wire.NewTxOut(999000, dest_pkScript)
	mtx.AddTxOut(output)
	// Add all outputs as inputs
	txHash, err := chainhash.NewHashFromStr("72e9c6e7490ba50cded735ec0f1d48d324c9509d2ef6a237b45f351638a1a28a")
	if err != nil {
		panic(err)
	}
	prevOut := wire.NewOutPoint(txHash, 0)
	txIn := wire.NewTxIn(prevOut, []byte{}, [][]byte{})
	mtx.AddTxIn(txIn)

	a, err := btcutil.DecodeWIF("Ky3n4DNba9iyyXPHyjsYKiKYCyJjtZiEDM3LVq7ShHVx8JzSzvDn")

	pk, _ := hex.DecodeString("76a9140f343240a5c2b4c532be07835716b2b46483f31188ac")
	if err != nil {
		panic(err)
	}
	sigScript, err := txscript.SignatureScript(mtx, 0, pk, txscript.SigHashAll, a.PrivKey, true)
	//sigScript, err := txscript.RawTxInSignature(mtx, 0, pk, txscript.SigHashAll, a.PrivKey)

	if err != nil {
		panic(err)
	}
	mtx.TxIn[0].SignatureScript = sigScript

	fmt.Println("sigScript", hex.EncodeToString(sigScript))
	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, mtx.SerializeSize()))
	if err := mtx.Serialize(buf); err != nil {
		panic(err)
	}
	txHex := hex.EncodeToString(buf.Bytes())
	fmt.Println("txHex :", txHex)
	//0c0000000100000000000000000000000000000000000000000000000000000000000000018aa2a13816355fb437a2f62e9d50c924d3481d0fec35d7de0ca50b49e7c6e972000000006b483045022100a98cc99cdfc33e1a10d0b24dfb97b36291131ad2eb324822dd69bb4de5ec562b022019c556e86952027ac4480b9f48553dd72d309cd25f424d1210f2ec9760ea82f8012103bec0cc32c8e2117488e9262ce7c99f6d7a56301c3420ae585ece49ec035d9ed7ffffffff01583e0f00000000001976a91497e206129e477018a56b6266af7f1dacbc061fe388ac00000000
	// txid: 018d24ff27e442c64ac2e24bc7930db6828bc3ef4264d4ebdb26052132a6b9be
}
