package types

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/blake2b-simd"
	//	"github.com/blocktree/go-owcrypt"
)

// CkbBlake2BHashPersonalization personal
const CkbBlake2BHashPersonalization = "ckb-default-hash"

func CreateEmptyRawTransaction(cellInput []CellInput, cellOutput []CellOutput) (*Transaction, error) {

	return &Transaction{}, nil
}

func newEmptyTransaction(cellInput []CellInput, cellOutput []CellOutput) (*Transaction, error) {
	return &Transaction{}, nil
}

func ComputeHash(txBolb []byte) ([]byte, error) {

	// Prepare hash
	config := &blake2b.Config{
		Size:   32,
		Person: []byte(CkbBlake2BHashPersonalization),
	}
	h, err := blake2b.New(config)
	if err != nil {
		fmt.Printf("Initial blake2b hash failure, %s\n", err)
		return []byte{}, err
	}

	// Hash tx
	h.Write(txBolb)
	txHash := h.Sum(nil)

	return txHash, nil
}

func ComputeWitnessHash(txBolb []byte, BobSecKey []byte) ([]Bytes, error) {

	// Prepare hash
	config := &blake2b.Config{
		Size:   32,
		Person: []byte(CkbBlake2BHashPersonalization),
	}
	h, err := blake2b.New(config)
	if err != nil {
		fmt.Printf("Initial blake2b hash failure, %s\n", err)
		return []Bytes{}, err
	}

	// Hash tx
	h.Write(txBolb)
	txHash := h.Sum(nil)
	fmt.Printf("txHash =  %s\n", hex.EncodeToString(txHash))

	lock := hex.EncodeToString(make([]byte, 65))
	witnessLock := Bytes(fmt.Sprintf("0x%s", lock))
	wa := WitnessArgs{
		Lock: &witnessLock,
	}

	wab, err := wa.Serialize()
	if err != nil {
		fmt.Printf("Serialize witness args failure, %s\n", err)
		return []Bytes{}, err
	}

	wabLen := make([]byte, 8)
	binary.LittleEndian.PutUint64(wabLen, uint64(len(wab)))

	// Hash again
	h.Reset()
	h.Write(txHash)
	h.Write(wabLen)
	h.Write(wab)
	witnessMessage := h.Sum(nil)

	bobSecKey := BobSecKey[:32]
	//_, err = hex.Decode(bobSecKey, BobSecKey)
	if err != nil {
		fmt.Println("Invalid bob private key")
		return []Bytes{}, err
	}
	// witnessSig, retCode := owcrypt.Signature(bobSecKey, nil, 0, witnessMessage, 32, owcrypt.ECC_CURVE_SECP256K1)
	// if retCode != owcrypt.SUCCESS {
	// 	return nil, errors.New("Sign transaction failed!")
	// }

	witnessSig, err := secp256k1.Sign(witnessMessage[:], bobSecKey)
	if err != nil {
		fmt.Printf("Calc witness failure: %s\n", err)
		return []Bytes{}, err
	}

	// Hex witness signature
	// 130 is hexed secp256k1 with recoverable public length
	witness := make([]byte, 130)
	hex.Encode(witness, witnessSig[:])

	// Refill witness args
	witnessLock = Bytes(fmt.Sprintf("0x%s", witness))
	wa = WitnessArgs{
		Lock: &witnessLock,
	}
	// Serialize witness args
	witness, err = wa.Serialize()
	if err != nil {
		fmt.Printf("witness serialize failure: %s\n", err)
		return []Bytes{}, err
	}

	// Update transaction with witness

	return []Bytes{Bytes(fmt.Sprintf("0x%s", hex.EncodeToString(witness)))}, nil
}

func SignRawTransaction(transaction *Transaction, prikey []byte) (*Transaction, error) {

	transactionSerize, err := transaction.Serialize()

	if err != nil {
		return nil, err
	}

	// hashBytes, err := ComputeHash(transactionSerize)
	// if err != nil || len(hashBytes) != 32 {
	// 	return nil, errors.New("Invalid tansaction hash!")
	// }
	// message := []byte{}
	// message = append(message[:], hashBytes[:]...)
	// usignItemCount := len(transaction.UnsignWitnesses)
	// if usignItemCount > 0 {
	// 	firstUnsignWitness := transaction.UnsignWitnesses[0]
	// 	lock, _ := hex.DecodeString(SIGNATURE_PLACEHOLDER)
	// 	emptyWitness := &WitnessArgs{
	// 		Lock:       (*Bytes)(unsafe.Pointer(&lock)),
	// 		InputType:  firstUnsignWitness.InputType,
	// 		OutputType: firstUnsignWitness.OutputType,
	// 	}
	// 	emptiedWitnessData, err := emptyWitness.Serialize()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	emptiedLen := len(emptiedWitnessData[:])
	// 	emptiedWitnessDataLen := serializeUint64(uint64(emptiedLen))
	// 	message = append(message[:], emptiedWitnessDataLen[:]...)
	// 	message = append(message[:], emptiedWitnessData[:]...)
	// 	for i := 1; i < usignItemCount; i++ {
	// 		unsignWitness := transaction.UnsignWitnesses[i]
	// 		emptiedWitnessData, err := unsignWitness.Serialize()
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		emptiedLen := len(emptiedWitnessData[:])
	// 		emptiedWitnessDataLen := serializeUint64(uint64(emptiedLen))
	// 		message = append(message[:], emptiedWitnessDataLen[:]...)
	// 		message = append(message[:], emptiedWitnessData[:]...)
	// 	}
	// } else {
	// 	return transaction, errors.New("invalidNumberOfWitnesses")
	// }

	witness, err := ComputeWitnessHash(transactionSerize, prikey)
	if err != nil {
		return nil, err
	}
	// signature, retCode := owcrypt.Signature(prikey, nil, 0, hashMessage, 32, owcrypt.ECC_CURVE_SECP256K1)
	// if retCode != owcrypt.SUCCESS {
	// 	return nil, errors.New("Sign transaction failed!")
	// }

	// transaction.UnsignWitnesses[0].Lock = (*Bytes)(unsafe.Pointer(&signature))
	// witness := make([]Bytes, usignItemCount)

	// for i := 0; i < usignItemCount; i++ {
	// 	unsignWitness := transaction.UnsignWitnesses[i]
	// 	emptiedWitnessData, err := unsignWitness.Serialize()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	signWitness := hex.EncodeToString(emptiedWitnessData)
	// 	witness[i] = Bytes(signWitness)
	// }

	transaction.Witnesses = witness
	//signature = serilizeS(signature)

	return transaction, nil
}

func VerifyAndCombineRawTransaction(emptyTrans, signature, pubkey string) (bool, string) {
	return true, "ab"
}
