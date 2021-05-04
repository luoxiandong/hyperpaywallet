package types

import (
	"encoding/hex"
	"fmt"

	"github.com/blocktree/go-owcdrivers/ckbTransaction/account"
)

// func PrvToPubAddress(prv string) (string, string, error) {
// 	privateKey, err := hex.DecodeString(prv)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	pub := PrvToPubKey(privateKey)
// 	address, err := PubKeyToAddress(pub)
// 	return hex.EncodeToString(pub), address, nil
// }
// func PrvToPubKey(prv []byte) (pub []byte) {

// 	return prv
// }
func AddressToPubKey(address string) (pub []byte, err error) {
	// prefix := "ckb"
	// hash := owcrypt.Hash(pub, 20, owcrypt.HASH_ALG_BLAKE2B)
	_, bech32Addr, err := account.DecodeSegWitAddress(address)
	if err != nil {
		return []byte{}, err
	}
	fmt.Println("publickey = " + hex.EncodeToString(bech32Addr))
	return bech32Addr, nil
}

// func ComputeHash(txBolb []byte) ([]byte, error) {

// 	// Prepare hash
// 	config := &blake2b.Config{
// 		Size:   32,
// 		Person: []byte(CkbBlake2BHashPersonalization),
// 	}
// 	h, err := blake2b.New(config)
// 	if err != nil {
// 		fmt.Printf("Initial blake2b hash failure, %s\n", err)
// 		return []byte{}, err
// 	}

// 	// Hash tx
// 	h.Write(txBolb)
// 	txHash := h.Sum(nil)

// 	return txHash, nil
// }
func PubKeyToAddress(pub []byte) (address string, err error) {
	prefix := "ckb"

	hash, err := ComputeHash(pub)
	if err != nil {
		return "", err
	}
	bech32Addr, err := account.EncodeSegWitAddress(prefix, hash[:20])
	if err != nil {
		return "", err
	}
	fmt.Println(bech32Addr)
	return bech32Addr, nil
}

func SignMessage(privateString string, hashMessage []byte) (rawTx string, err error) {
	privateKey, err := hex.DecodeString(privateString)
	// publicKey := PrvToPubKey(privateKey)
	fmt.Println(privateKey)
	return "", err
}
