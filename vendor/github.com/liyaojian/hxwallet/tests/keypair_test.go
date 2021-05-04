package tests

import (
	"encoding/hex"
	"fmt"
	"github.com/liyaojian/hxwallet/keypair"
	"testing"
)

// 助记词生成钱包
func Test_GenerateKeyPair(t *testing.T) {
	mnemonic := "THIS IS A TERRIBLE BRAINKEY SEED WORD SEQUENCE"
	keyPair, err := keypair.GenerateKeyPair(mnemonic)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(keyPair.BrainKey)
	fmt.Println("私钥:", keyPair.PrivateKey.ToWIF())                                    // 私钥
	fmt.Println("公钥:", keyPair.PrivateKey.PublicKey().String())                       // 公钥
	fmt.Println("hex公钥:", hex.EncodeToString(keyPair.PrivateKey.PublicKey().Bytes())) // 公钥(hex格式)
	address, _ := keyPair.PrivateKey.PublicKey().ToAddress()
	fmt.Println("地址:", address) // 地址

	//embody balcony whisper arctic elephant method grace essay process magic trumpet sport
	//私钥: 5KNkNu1GYrehgLcSz2GMznVoQHfmce8t24ZTPEGoygvHbtgGW7k
	//公钥: HX68Xt87frzo8DJmPpJqvbwpAmDxTXf3dQXLRz3A1DTKLwcvAjFP
	//hex公钥: 02a39573e0065febf9690e46d60f63d5cda8ff3c62fb8f8a2dd458c09045fce3f5
	//地址: HXNb7KtLSX5yaj4111beUh486uKiodUZkL9J
}

// 私钥生成钱包(公钥和地址)
func Test_PrivateToPublic(t *testing.T) {
	testPri := "5KAffU3Pw7RNJAJ3d1qUrJ6QPVb6UFx6CJ4MhgfoHL7YwYspHhs"
	pub, err := keypair.PrivateToPublic(testPri)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("公钥:", pub.String())                       // 公钥
	fmt.Println("hex公钥:", hex.EncodeToString(pub.Bytes())) // 公钥
	address, _ := pub.ToAddress()
	fmt.Println("地址:", address) // 地址
}

// 验证私钥合法性
func Test_IsValidPrivate(t *testing.T) {
	testPri := "5KAffU3Pw7RNJAJ3d1qUrJ6QPVb6UFx6CJ4MhgfoHL7YwYspHhs"
	isValidPrivate := keypair.IsValidPrivate(testPri)
	if !isValidPrivate {
		t.Error("私钥不合法")
	}
}

// 验证公钥合法性
func Test_IsValidPublic(t *testing.T) {
	testPub := "HX8mT7XvtTARjdZQ9bqHRoJRMf7P7azFqTQACckaVenM2GmJyxLh"
	isValidPublic := keypair.IsValidPublic(testPub)
	if !isValidPublic {
		t.Errorf("公钥不合法")
	}
}
