package model

import (
	"fmt"
	"testing"
)

func TestAccountAddress(t *testing.T) {

	//keyPair := model.KeyPairCreate("bad677392e447b736d1814def5d48260cdff082efc5bbcf6be854ca7056ed58a")
	keyPair := KeyPairCreate("be052e8d052f25c5950375f520c39a5d80d0885a57aa98c2d30a50286734b099")

	publicKey := keyPair.PublicString()
	privateKey := keyPair.PrivateString()

	address := ToAddress(publicKey, Data.Mainnet.ID)

	fmt.Println("PrivateKey:\t", privateKey)

	fmt.Println("PublicKey:\t", publicKey)

	fmt.Println("Address:", address)
}
