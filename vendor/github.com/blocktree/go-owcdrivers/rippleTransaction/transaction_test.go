package rippleTransaction

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestXRPTransaction(t *testing.T) {
	from := "rht6f38DukqqsbcVbdtJfFLYJJ1LVS27rQ"
	//from = "rGjAhunofbk83nUpcKUrYLP5Fp132w4a1u"
	to := "rGuvnj4uiG16oGuvNUfDTBvafYGWLsAksY"
	pubkey := "03185cfd62490fe166d2ff79b089584df6a1bda9f01d74a9f9a82629859d21e298"
	sequence := uint32(1)
	amount := uint64(2000000000)
	fee := uint64(10000)
	lastLedgerSequence := uint32(353535)
	memoType := "client"
	memoData := "111"
	memoFormat := "text/plain"

	emptyTrans, hash, err := CreateEmptyRawTransactionAndHash(from, pubkey, sequence, to, amount, fee, lastLedgerSequence, memoType, memoData, memoFormat)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("empty transaction : \n", emptyTrans)
		fmt.Println("transaction hash  : \n", hash)
	}

	prikey, _ := hex.DecodeString("4f482bfbd8a640d823f3b178b7ef637a25aa07b1ceeb563c2f2ea834402bedc2")
	fmt.Println("prikey ", prikey)
	signature, err := SignRawTransaction(hash, prikey)
	if err != nil {
		t.Error(err)
	} else {
		//signature = "168a76d7bef92f30761a03c0d039f9f018d32af756548f6e9ef41d1098d94ab5343513ca97f6437beca53b2e2ef36ee79a8fd2f4397473dd0031c12df11e6b83"
		fmt.Println("signature data : \n", signature)
	}

	//
	pass, signedTrans := VerifyAndCombinRawTransaction(emptyTrans, signature, pubkey)
	if pass {
		fmt.Println("signed transaction : \n", signedTrans)
	} else {
		t.Error("Verify transaction failed!")
	}
	if signedTrans == "12000022800000002400000001201b000564ff614000000077359400684000000000002710732103185cfd62490fe166d2ff79b089584df6a1bda9f01d74a9f9a82629859d21e29874473045022100fdd4058c94a0a9305bafbafabbf6d201e221534de145339d52f21f1de43b5c7c02203917b6f3274223a254141a98444d5ccf9ced9bf2127a38db47d614897cb9162281142aaf11ebe6cd01f9b7f28c8633c0a5486df187b78314ae8bd3b80ad0f8aa66e4625afa3d922f8f48a3ebf9ea7c06636c69656e747d033131317e0a746578742f706c61696ee1f1" {
		fmt.Println("signedTrans 相等")
	} else {
		fmt.Println("signedTrans 不相等")

	}

}
func TestTransaction(t *testing.T) {
	from := "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
	pubkey := "0330e7fc9d56bb25d6893ba3f317ae5bcf33b3291bd63db32654a313222f7fd020"
	to := "rb1fWuuAEtPUaeEWxocV3h4x5JwDTFZzH"
	sequence := uint32(44196)
	amount := uint64(2000000000)
	fee := uint64(2000000000)
	lastLedgerSequence := uint32(50897423)
	memoType := "client"
	memoData := "111"
	memoFormat := "text/plain"

	emptyTrans, hash, err := CreateEmptyRawTransactionAndHash(from, pubkey, sequence, to, amount, fee, lastLedgerSequence, memoType, memoData, memoFormat)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("empty transaction : \n", emptyTrans)
		fmt.Println("transaction hash  : \n", hash)
	}

	prikey, _ := hex.DecodeString("1acaaedece405b2a958212629e16f2eb46b153eee94cdd350fdeff52795525b7")
	signature, err := SignRawTransaction(hash, prikey)
	if err != nil {
		t.Error(err)
	} else {
		//signature = "168a76d7bef92f30761a03c0d039f9f018d32af756548f6e9ef41d1098d94ab5343513ca97f6437beca53b2e2ef36ee79a8fd2f4397473dd0031c12df11e6b83"
		fmt.Println("signature data : \n", signature)
	}

	//
	pass, signedTrans := VerifyAndCombinRawTransaction(emptyTrans, signature, pubkey)
	if pass {
		fmt.Println("signed transaction : \n", signedTrans)
	} else {
		t.Error("Verify transaction failed!")
	}
}
