package transaction

import (
	_ "encoding/hex"
	_ "fmt"

	"github.com/blocktree/go-owcdrivers/ed25519WalletKey"
	"golang.org/x/crypto/ed25519"
)

type ALGOTransaction struct {
	From        []byte
	To          []byte
	Fee         uint64
	Amount      uint64
	FirstRound  uint64
	LastRound   uint64
	Type        string
	Note        []byte
	GenesisId   string
	GenesisHash []byte
}

func SignData(data []byte, sk []byte) []byte {
	rawSignData := []byte{}
	rawSignData = append(rawSignData[:], []byte{84, 88}...)
	rawSignData = append(rawSignData[:], data[:]...)
	publicKey := ed25519WalletKey.WalletPubKeyFromKeyBytes(sk)
	privateByte := []byte{}

	privateByte = append(privateByte[:], sk[:]...)
	privateByte = append(privateByte[:], publicKey[:]...)

	signature := ed25519.Sign(privateByte[:], rawSignData[:])
	//fmt.Println("ed25519 prv = ", hex.EncodeToString(sk))
	return signature[:]
}
func (algo *ALGOTransaction) SerializeRawTx() []byte {
	resData := []byte{}
	size := 9
	if len(algo.Note) > 0 {
		size = size + 1
	}
	resData = append(resData[:], byte(0x80+size))
	resData = append(resData[:], ALGOEcodeString("amt")[:]...)
	resData = append(resData[:], ALGOEcodeNumber(algo.Amount)[:]...)

	resData = append(resData[:], ALGOEcodeString("fee")[:]...)
	resData = append(resData[:], ALGOEcodeNumber(algo.Fee)[:]...)

	resData = append(resData[:], ALGOEcodeString("fv")[:]...)
	resData = append(resData[:], ALGOEcodeNumber(algo.FirstRound)[:]...)

	resData = append(resData[:], ALGOEcodeString("gen")[:]...)
	resData = append(resData[:], ALGOEcodeString(algo.GenesisId)[:]...)

	resData = append(resData[:], ALGOEcodeString("gh")[:]...)
	resData = append(resData[:], ALGOEcodeByte(algo.GenesisHash)[:]...)
	//fmt.Println(hex.EncodeToString(resData))
	resData = append(resData[:], ALGOEcodeString("lv")[:]...)
	resData = append(resData[:], ALGOEcodeNumber(algo.LastRound)[:]...)
	//fmt.Println(hex.EncodeToString(resData))

	if len(algo.Note) > 0 {
		resData = append(resData[:], ALGOEcodeString("note")[:]...)
		resData = append(resData[:], ALGOEcodeByte(algo.Note)[:]...)
	}
	//fmt.Println(hex.EncodeToString(resData))
	resData = append(resData[:], ALGOEcodeString("rcv")[:]...)
	resData = append(resData[:], ALGOEcodeByte(algo.To)[:]...)
	//fmt.Println(hex.EncodeToString(resData))
	resData = append(resData[:], ALGOEcodeString("snd")[:]...)
	resData = append(resData[:], ALGOEcodeByte(algo.From)[:]...)
	//fmt.Println(hex.EncodeToString(resData))
	resData = append(resData[:], ALGOEcodeString("type")[:]...)
	resData = append(resData[:], ALGOEcodeString(algo.Type)[:]...)

	//resData = append(resData[:], data[:]...)
	return resData
}

func (algo *ALGOTransaction) SerializeRawSignTx(signature []byte) []byte {
	resData := []byte{}

	resData = append(resData[:], byte(0x80+2))

	resData = append(resData[:], ALGOEcodeString("sig")[:]...)
	resData = append(resData[:], ALGOEcodeByte(signature)[:]...)

	resData = append(resData[:], ALGOEcodeString("txn")[:]...)
	resData = append(resData[:], algo.SerializeRawTx()[:]...)
	return resData
}

func ALGOEcodeString(branch string) []byte {
	forgeData := []byte{}
	length := len([]byte(branch))
	if length < 0x20 {
		forgeData = append(forgeData[:], byte(uint8(0xa0+length)))
		//forgeData = append(forgeData[:], data[:]...)
	} else if length < 0x100 {
		forgeData = append(forgeData[:], byte(uint8(0xd9)))
		forgeData = append(forgeData[:], byte(uint8(length)))
	} else if length < 0x10000 {
		forgeData = append(forgeData[:], byte(uint8(0xda)))
		forgeData = append(forgeData[:], encode16BE(uint16(length))[:]...)
	} else if length < 0x100000000 {
		forgeData = append(forgeData[:], byte(uint8(0xdb)))
		forgeData = append(forgeData[:], encode32BE(uint32(length))[:]...)
	}
	forgeData = append(forgeData[:], []byte(branch)[:]...)

	return forgeData
}

func ALGOEcodeByte(data []byte) []byte {
	forgeData := []byte{}
	length := len(data)
	if length < 0x100 {
		forgeData = append(forgeData[:], byte(uint8(0xc4)))
		forgeData = append(forgeData[:], byte(uint8(length)))
	} else if length < 0x10000 {
		forgeData = append(forgeData[:], byte(uint8(0xc5)))
		forgeData = append(forgeData[:], encode16BE(uint16(length))[:]...)
	} else if length < 0x100000000 {
		forgeData = append(forgeData[:], byte(uint8(0xc6)))
		forgeData = append(forgeData[:], encode32BE(uint32(length))[:]...)
	}
	forgeData = append(forgeData[:], data[:]...)

	return forgeData
}
func ALGOEcodeNumber(number uint64) []byte {
	forgeData := []byte{}
	length := number
	if length < 0x80 {
		forgeData = append(forgeData[:], byte(uint8(length)))
		//forgeData = append(forgeData[:], data[:]...)
	} else if length < 0x100 {
		forgeData = append(forgeData[:], byte(uint8(0xcc)))
		forgeData = append(forgeData[:], byte(uint8(length)))
	} else if length < 0x10000 {
		forgeData = append(forgeData[:], byte(uint8(0xcd)))
		forgeData = append(forgeData[:], encode16BE(uint16(length))[:]...)
	} else if length < 0x100000000 {
		forgeData = append(forgeData[:], byte(uint8(0xce)))
		forgeData = append(forgeData[:], encode32BE(uint32(length))[:]...)
	} else {
		forgeData = append(forgeData[:], byte(uint8(0xcf)))
		forgeData = append(forgeData[:], encode64BE(uint64(length))[:]...)
	}

	return forgeData
}
func encode16BE(val uint16) []byte {
	var b [2]byte
	b[0] = byte(val >> 8)
	b[1] = byte(val)
	return b[:]
}

func encode32BE(i uint32) []byte {
	var b [4]byte
	b[0] = byte(i >> 24)
	b[1] = byte(i >> 16)
	b[2] = byte(i >> 8)
	b[3] = byte(i)
	return b[:]
}
func encode64BE(i uint64) []byte {
	var b [8]byte
	b[0] = byte(i >> 56)
	b[1] = byte(i >> 48)
	b[2] = byte(i >> 40)
	b[3] = byte(i >> 32)
	b[4] = byte(i >> 24)
	b[5] = byte(i >> 16)
	b[6] = byte(i >> 8)
	b[7] = byte(i)
	return b[:]
}

// func putint(b []byte, i uint64) (size int) {
// 	switch {
// 	case i < (1 << 8):
// 		b[0] = byte(i)
// 		return 1
// 	case i < (1 << 16):
// 		b[0] = byte(i >> 8)
// 		b[1] = byte(i)
// 		return 2
// 	case i < (1 << 24):
// 		b[0] = byte(i >> 16)
// 		b[1] = byte(i >> 8)
// 		b[2] = byte(i)
// 		return 3
// 	case i < (1 << 32):
// 		b[0] = byte(i >> 24)
// 		b[1] = byte(i >> 16)
// 		b[2] = byte(i >> 8)
// 		b[3] = byte(i)
// 		return 4
// 	case i < (1 << 40):
// 		b[0] = byte(i >> 32)
// 		b[1] = byte(i >> 24)
// 		b[2] = byte(i >> 16)
// 		b[3] = byte(i >> 8)
// 		b[4] = byte(i)
// 		return 5
// 	case i < (1 << 48):
// 		b[0] = byte(i >> 40)
// 		b[1] = byte(i >> 32)
// 		b[2] = byte(i >> 24)
// 		b[3] = byte(i >> 16)
// 		b[4] = byte(i >> 8)
// 		b[5] = byte(i)
// 		return 6
// 	case i < (1 << 56):
// 		b[0] = byte(i >> 48)
// 		b[1] = byte(i >> 40)
// 		b[2] = byte(i >> 32)
// 		b[3] = byte(i >> 24)
// 		b[4] = byte(i >> 16)
// 		b[5] = byte(i >> 8)
// 		b[6] = byte(i)
// 		return 7
// 	default:
// 		b[0] = byte(i >> 56)
// 		b[1] = byte(i >> 48)
// 		b[2] = byte(i >> 40)
// 		b[3] = byte(i >> 32)
// 		b[4] = byte(i >> 24)
// 		b[5] = byte(i >> 16)
// 		b[6] = byte(i >> 8)
// 		b[7] = byte(i)
// 		return 8
// 	}
// }
