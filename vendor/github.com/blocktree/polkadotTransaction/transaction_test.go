package polkadotTransaction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	//"github.com/blocktree/go-owcdrivers/polkadotTransaction/codec"
	"testing"

	"github.com/blocktree/go-owcrypt/eddsa"
	"github.com/blocktree/polkadotTransaction/codec"
	"github.com/polkadot-adapter/address"
)

func Test_transaction1(t *testing.T) {
	tx := TxStruct{
		//发送方公钥
		SenderPubkey: "86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f", //"88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//接收方公钥
		RecipientPubkey: "88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//发送金额（最小单位）
		Amount: 12,
		//nonce
		Nonce: 1,
		//手续费（最小单位）
		Fee: 20,
		//当前高度
		BlockHeight: 1470121,
		//当前高度区块哈希
		BlockHash: "cbd176fab2c85663d66fdc1a503f8c3b65b81e0eb76db35c9ad955d3045f9c54",
		//创世块哈希
		GenesisHash: "b0a8d493285c2df73290dfb7e61f870f17b41801197a149ca93654499ea3dafe",
		//spec版本
		SpecVersion: 1059,
		//Transaction版本
		TxVersion: 1,
	}

	// 创建空交易单和待签消息
	emptyTrans, message, err := tx.CreateEmptyTransactionAndMessage()
	if err != nil {
		t.Error("create failed : ", err)
		return
	}
	fmt.Println("空交易单 ： ", emptyTrans)
	fmt.Println("待签消息 ： ", message)

	// 签名
	prikey, _ := hex.DecodeString("e86bcaaab0a5aa5e3f3b0885db7e932e34eddb5a620b6bcc097a4b236a5a0354")
	pubkey, _ := eddsa.ED25519_genPub(prikey)
	fmt.Println("pubkey = ", hex.EncodeToString(pubkey))
	signature, err := SignTransaction(message, prikey)
	if err != nil {
		t.Error("sign failed")
		return
	}
	fmt.Println("签名结果 ： ", hex.EncodeToString(signature))

	// signature, _ := hex.DecodeString("1cc69f7ba50ee37793c83d74b21f50239894e8733cdf7fd13565eded13ba97d8229fc51174035be6d4543908f58b016efd0aae137f8ad584c5540002326bc809")

	// 验签与交易单合并
	signedTrans, pass := VerifyAndCombineTransaction(emptyTrans, hex.EncodeToString(signature))
	if pass {
		fmt.Println("验签成功")
		fmt.Println("签名交易单 ： ", signedTrans)
	} else {
		t.Error("验签失败")
	}
}

func Test_WNDtransaction(t *testing.T) {
	tx := TxStruct{
		//发送方公钥
		//SenderPubkey: "ec93d475a57f0997256259540cf22f9cc5c0a0da27f93ebb7c237feaafda8cfb", //"88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",

		SenderPubkey: "86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f", //"88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//接收方公钥
		RecipientPubkey: "88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//发送金额（最小单位）
		Amount: 100000000000,
		//nonce
		Nonce: 0,
		//手续费（最小单位）
		Fee: 10000000000,
		//当前高度
		BlockHeight: 1497453,
		//当前高度区块哈希
		BlockHash: "a92fd876d9378771980831b040da82d5f3196aa96e2d8f74741c628a9b51bbc9",
		//创世块哈希
		GenesisHash: "e143f23803ac50e8f6f8e62695d1ce9e4e1d68aa36c1cd2cfd15340213f3423e",
		//spec版本
		SpecVersion: 41,
		//Transaction版本
		TxVersion: 2,
	}
	pkByte, _ := hex.DecodeString("86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f")
	fmt.Println("TestNet Address = ", address.PublicKeyToAddress(pkByte, "testnet"))

	//ec93d475a57f0997256259540cf22f9cc5c0a0da27f93ebb7c237feaafda8cfb

	// 创建空交易单和待签消息
	emptyTrans, message, err := tx.CreateEmptyTransactionAndMessage()
	if err != nil {
		t.Error("create failed : ", err)
		return
	}
	fmt.Println("空交易单 ： ", emptyTrans)
	fmt.Println("待签消息 ： ", message)

	// 签名
	prikey, _ := hex.DecodeString("e86bcaaab0a5aa5e3f3b0885db7e932e34eddb5a620b6bcc097a4b236a5a0354")
	pk, _ := eddsa.ED25519_genPub(prikey)

	fmt.Println("pk ： ", "86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f")

	fmt.Println("pk ： ", hex.EncodeToString(pk))

	fmt.Println("TestNet Address = ", address.PublicKeyToAddress(pk, "testnet"))

	signature, err := SignTransaction(message, prikey)
	if err != nil {
		t.Error("sign failed")
		return
	}
	fmt.Println("签名结果 ： ", hex.EncodeToString(signature))

	// signature, _ := hex.DecodeString("1cc69f7ba50ee37793c83d74b21f50239894e8733cdf7fd13565eded13ba97d8229fc51174035be6d4543908f58b016efd0aae137f8ad584c5540002326bc809")

	// 验签与交易单合并
	signedTrans, pass := VerifyAndCombineTransaction(emptyTrans, hex.EncodeToString(signature))
	if pass {
		fmt.Println("验签成功")
		fmt.Println("签名交易单 ： ", signedTrans)
	} else {
		t.Error("验签失败")
	}
}

func Test_transaction(t *testing.T) {
	tx := TxStruct{
		//发送方公钥
		//SenderPubkey: "ec93d475a57f0997256259540cf22f9cc5c0a0da27f93ebb7c237feaafda8cfb", //"88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",

		SenderPubkey: "86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f", //"88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//接收方公钥
		RecipientPubkey: "88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee",
		//发送金额（最小单位）
		Amount: 100000000000,
		//nonce
		Nonce: 0,
		//手续费（最小单位）
		Fee: 10000000000,
		//当前高度
		BlockHeight: 1470562,
		//当前高度区块哈希
		BlockHash: "7f8dbca5c2745cc4da23a31780f5fec3fb1dcbc0bf0b406ed91a7b5ec3e06438",
		//创世块哈希
		GenesisHash: "b0a8d493285c2df73290dfb7e61f870f17b41801197a149ca93654499ea3dafe",
		//spec版本
		SpecVersion: 1059,
		//Transaction版本
		TxVersion: 1,
	}
	pkByte, _ := hex.DecodeString("86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f")
	fmt.Println("TestNet Address = ", address.PublicKeyToAddress(pkByte, "testnet"))

	//ec93d475a57f0997256259540cf22f9cc5c0a0da27f93ebb7c237feaafda8cfb

	// 创建空交易单和待签消息
	emptyTrans, message, err := tx.CreateEmptyTransactionAndMessage()
	if err != nil {
		t.Error("create failed : ", err)
		return
	}
	fmt.Println("空交易单 ： ", emptyTrans)
	fmt.Println("待签消息 ： ", message)

	// 签名
	prikey, _ := hex.DecodeString("e86bcaaab0a5aa5e3f3b0885db7e932e34eddb5a620b6bcc097a4b236a5a0354")
	signature, err := SignTransaction(message, prikey)
	if err != nil {
		t.Error("sign failed")
		return
	}
	fmt.Println("签名结果 ： ", hex.EncodeToString(signature))

	// signature, _ := hex.DecodeString("1cc69f7ba50ee37793c83d74b21f50239894e8733cdf7fd13565eded13ba97d8229fc51174035be6d4543908f58b016efd0aae137f8ad584c5540002326bc809")

	// 验签与交易单合并
	signedTrans, pass := VerifyAndCombineTransaction(emptyTrans, hex.EncodeToString(signature))
	if pass {
		fmt.Println("验签成功")
		fmt.Println("签名交易单 ： ", signedTrans)
	} else {
		t.Error("验签失败")
	}
}

func Test_json(t *testing.T) {
	ts := TxStruct{
		SenderPubkey:    "123",
		RecipientPubkey: "",
		Amount:          0,
		Nonce:           0,
		Fee:             0,
		BlockHeight:     0,
		BlockHash:       "234",
		GenesisHash:     "345",
		SpecVersion:     0,
	}

	js, _ := json.Marshal(ts)

	fmt.Println(string(js))
}

func Test_decode(t *testing.T) {
	en, _ := codec.Encode(Compact_U32, uint64(139))
	fmt.Println(en)
}
