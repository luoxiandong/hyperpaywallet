package polkadotTransaction

import (
	//"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/blocktree/go-owcdrivers/ed25519WalletKey"
	"golang.org/x/crypto/ed25519"
)

func (ts TxStruct) CreateEmptyTransactionAndMessage(transferCode string) (string, string, error) {

	tp, err := ts.NewTxPayLoad(transferCode)
	if err != nil {
		return "", "", err
	}

	return ts.ToJSONString(), tp.ToBytesString(transferCode), nil
}

func SignTransaction(msgStr string, prikey []byte) ([]byte, error) {
	msg, err := hex.DecodeString(msgStr)
	if err != nil || len(msg) == 0 {
		return nil, errors.New("invalid message to sign")
	}

	if prikey == nil || len(prikey) != 32 {
		return nil, errors.New("invalid private key")
	}

	publicKey := ed25519WalletKey.WalletPubKeyFromKeyBytes(prikey)
	privateByte := []byte{}
	fmt.Println(hex.EncodeToString(publicKey))
	privateByte = append(privateByte[:], prikey[:]...)
	privateByte = append(privateByte[:], publicKey[:]...)
	signature := ed25519.Sign(privateByte[:], msg)

	//signature, err := eddsa.ED25519_sign(prikey, msg)

	if err != nil {
		return nil, err
	}
	//signature, retCode := owcrypt.Signature(prikey, nil, 0, msg, 32, owcrypt.ECC_CURVE_ED25519)

	//signature, retCode := owcrypt.Signature(prikey, nil, msg, owcrypt.ECC_CURVE_ED25519)
	// if retCode != owcrypt.SUCCESS {
	// 	return nil, errors.New("sign failed")
	// }

	return signature[:], nil
}

func VerifyAndCombineTransaction(transferCode, emptyTrans, signature string) (string, bool) {
	ts, err := NewTxStructFromJSON(emptyTrans)
	if err != nil {
		return "", false
	}

	tp, err := ts.NewTxPayLoad(transferCode)
	if err != nil {
		return "", false
	}

	msg, _ := hex.DecodeString(tp.ToBytesString(transferCode))

	pubkey, _ := hex.DecodeString(ts.SenderPubkey)

	sig, err := hex.DecodeString(signature)
	if err != nil || len(sig) != 64 {
		return "", false
	}
	// if owcrypt.SUCCESS != owcrypt.Verify(pubkey, nil, msg, sig, owcrypt.ECC_CURVE_ED25519) {
	// 	return "", false
	// }
	suc := ed25519.Verify(ed25519.PublicKey(pubkey), msg, sig)

	//suc := eddsa.ED25519_verify(pubkey, msg, sig)
	if suc == false {
		return "", false
	}
	signned, err := ts.GetSignedTransaction(transferCode, signature)
	if err != nil {
		return "", false
	}

	return signned, true
}
