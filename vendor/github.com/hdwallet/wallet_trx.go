package hdwallet

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	//"crypto/ecdsa"
	//"crypto/elliptic"
	//"crypto/rand"

	owcrypt "github.com/blocktree/go-owcrypt"

	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/tron-adapter/tron/grpc-gateway/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/tronadpter"
	"github.com/tronadpter/address"

	"github.com/tronadpter/tron"
)

func init() {
	coins[TRX] = newTRX
}

type trx struct {
	*btc
}

const (
	// Symbol 币种
	Symbol = "TRX"
	//MasterKey key for master
	MasterKey = "Troncoin seed"
	//CurveType to generate ChildKey by BIP32
	CurveType               = owcrypt.ECC_CURVE_SECP256K1
	Decimals                = 0             // Decimals
	SUN               int64 = 1             //最小单位
	TRXValue          int64 = SUN * 1000000 //1 TRX = 1000000 * sun
	GasPrice                = SUN * 10
	CreateAccountCost       = SUN * 100000 //0.1 TRX = 100000 * sun
)

func newTRX(key *Key) Wallet {

	token := newBTC(key).(*btc)
	token.name = "Tron"
	token.symbol = "TRX"
	token.key.opt.Params = &BTCParams

	return &trx{btc: token}
}

func (c *trx) GetWalletAccountFromWif() (*WalletAccount, error) {

	wif := c.GetKey().Wif

	if len(wif) == 0 {
		return &WalletAccount{
			Res:    0,
			ErrMsg: "私钥不正确",
		}, nil
	}

	priv, err := crypto.HexToECDSA(wif)
	if err != nil {
		fmt.Println("Wif err : ", err.Error())
		return nil, err
	}
	pub := crypto.FromECDSAPub(&priv.PublicKey)
	addr := address.FromPublicKey(&priv.PublicKey).ToBase58()

	return &WalletAccount{
		Res:        1,
		PrivateKey: wif,
		PublicKey:  hex.EncodeToString(pub),
		Address:    addr,
	}, nil
}

func (c *trx) GetAddress() (string, error) {
	return crypto.PubkeyToAddress(*c.key.PublicECDSA).Hex(), nil
}

func (c *trx) GetWalletAccount() *WalletAccount {
	if c.GetKey().Extended == nil {
		return &WalletAccount{
			Res: 0,
		}
	}

	addr := address.FromPublicKey(c.key.PublicECDSA).ToBase58()
	pri := crypto.FromECDSA(c.key.PrivateECDSA)
	pub := crypto.FromECDSAPub(c.key.PublicECDSA)
	// pub := c.key.PublicHex(true)
	return &WalletAccount{
		Res:        1,
		Address:    addr,
		PrivateKey: hex.EncodeToString(pri),
		PublicKey:  hex.EncodeToString(pub),
		Seed:       c.GetKey().Seed,
	}
}

// 构造交易
func (c *trx) SignRawTransaction(signIn *SignInput) (*SignResult, error) {

	// priv, err := crypto.HexToECDSA(signIn.PrivateKey)
	// if err != nil {
	// 	fmt.Println("Wif err : ", err.Error())
	// 	return nil, err
	// }
	//amount := strconv.FormatInt(signIn.Amount, 10)
	amount := signIn.LargeAmount
	toAddressBytes, err := addressEncoder.AddressDecode(signIn.DestAddr, addressEncoder.TRON_mainnetAddress)
	if err != nil {
		fmt.Printf("toAddress decode failed failed;unexpected error:%v", err)
		return nil, err
	}
	toAddressBytes = append([]byte{0x41}, toAddressBytes...)
	//fmt.Printf("Prefix  =  %s\n", hex.EncodeToString(addressEncoder.TRON_mainnetAddress.Prefix))

	ownerAddressBytes, err := addressEncoder.AddressDecode(signIn.SrcAddr, addressEncoder.TRON_mainnetAddress)
	if err != nil {
		fmt.Printf("ownerAddress decode failed failed;unexpected error:%v", err)
		return nil, err
	}
	ownerAddressBytes = append([]byte{0x41}, ownerAddressBytes...)
	fmt.Printf("ownerAddressBytes  =  %s\n", hex.EncodeToString(ownerAddressBytes))
	amountDec := address.StringNumToBigIntWithExp(amount, Decimals)

	block, err := tronadpter.GetNowBlock()
	if err != nil {
		return nil, err
	}
	blockID, err := hex.DecodeString(block.GetBlockHashID())

	refBlockBytes, refBlockHash := blockID[6:8], blockID[8:16]
	timestamp := time.Now().UnixNano() / 1000000 // <int64
	contracts := []*core.Transaction_Contract{}

	if signIn.ContractAddr == "" {
		// Check amount: amount * 1000000
		// ******** Generate TX Contract ********
		tc := &core.TransferContract{
			OwnerAddress: ownerAddressBytes,
			ToAddress:    toAddressBytes,
			Amount:       amountDec.Int64(),
		}

		tcRaw, err := proto.Marshal(tc)
		if err != nil {
			return nil, err
		}

		txContact := &core.Transaction_Contract{
			Type:         core.Transaction_Contract_TransferContract,
			Parameter:    &any.Any{Value: tcRaw, TypeUrl: "type.googleapis.com/protocol.TransferContract"},
			Provider:     nil,
			ContractName: nil,
		}
		contracts = []*core.Transaction_Contract{txContact}

	} else if strings.ToLower(signIn.Type) == "trc10" {
		tc := &core.TransferAssetContract{
			OwnerAddress: ownerAddressBytes,
			ToAddress:    toAddressBytes,
			Amount:       amountDec.Int64(),
			AssetName:    []byte(signIn.ContractAddr),
		}
		tcRaw, err := proto.Marshal(tc)
		if err != nil {
			return nil, err
		}

		txContact := &core.Transaction_Contract{
			Type:         core.Transaction_Contract_TransferAssetContract,
			Parameter:    &any.Any{Value: tcRaw, TypeUrl: "type.googleapis.com/protocol.TransferAssetContract"},
			Provider:     nil,
			ContractName: nil,
		}
		contracts = []*core.Transaction_Contract{txContact}

	} else if strings.ToLower(signIn.Type) == "trc20" {
		contractAddressBytes, err := addressEncoder.AddressDecode(signIn.ContractAddr, addressEncoder.TRON_mainnetAddress)
		contractAddressBytes = append([]byte{0x41}, contractAddressBytes...)

		if err != nil {
			return nil, err
		}

		var funcParams []tron.SolidityParam
		funcParams = append(funcParams, tron.SolidityParam{
			ParamType:  tron.SOLIDITY_TYPE_ADDRESS,
			ParamValue: hex.EncodeToString(toAddressBytes),
		})

		funcParams = append(funcParams, tron.SolidityParam{
			ParamType:  tron.SOLIDITY_TYPE_UINT256,
			ParamValue: amountDec,
		})

		//fmt.Println("make token transfer data, amount:", amount.String())
		dataHex, err := tron.MakeTransactionParameter(tron.TRC20_TRANSFER_METHOD_ID, funcParams)
		if err != nil {
			return nil, err
		}

		data, err := hex.DecodeString(dataHex)
		if err != nil {
			return nil, err
		}

		tc := &core.TriggerSmartContract{
			OwnerAddress:    ownerAddressBytes,
			ContractAddress: contractAddressBytes,
			Data:            data,
		}

		tcRaw, err := proto.Marshal(tc)
		if err != nil {
			return nil, err
		}

		txContact := &core.Transaction_Contract{
			Type:         core.Transaction_Contract_TriggerSmartContract,
			Parameter:    &any.Any{Value: tcRaw, TypeUrl: "type.googleapis.com/protocol.TriggerSmartContract"},
			Provider:     nil,
			ContractName: nil,
		}
		contracts = []*core.Transaction_Contract{txContact}
	}

	// ******** Create Traction ********
	txRaw := &core.TransactionRaw{
		RefBlockBytes: refBlockBytes,
		RefBlockHash:  refBlockHash,
		Contract:      contracts,
		Expiration:    timestamp + 10*60*60*60,
		// Timestamp:     timestamp,
	}
	txRaw.FeeLimit = 100000
	if signIn.Fee > 0 {
		txRaw.FeeLimit = signIn.Fee
	}
	tx := &core.Transaction{
		RawData: txRaw,
		// Signature: nil,
		// Ret:       nil,
	}
	// ******** TX Encoding ********
	x, err := proto.Marshal(tx)
	if err != nil {
		fmt.Printf("marshal tx failed;unexpected error:%v", err)
		return nil, err
	}
	txRawHex := hex.EncodeToString(x)
	txHash, err := getTxHash1(txRawHex)
	pk, err := hex.DecodeString(signIn.PrivateKey)
	fmt.Println("txHash = ", hex.EncodeToString(txHash))

	if err != nil {
		return nil, err
	}
	//sign, err := owcrypt.Signature(pk, nil, 0, txHash, 32, owcrypt.ECC_CURVE_SECP256R1)

	sign, ret := tronadpter.TronSignature(pk, txHash)
	if ret != owcrypt.SUCCESS {
		return nil, fmt.Errorf("sign txHash failed")
	}
	// var signs [][]byte
	// signs = append(signs, sign)
	// tx = &core.Transaction{
	// 	RawData:   txRaw,
	// 	Signature: signs,
	// 	Ret:       nil,
	// }
	// signedTx, err := proto.Marshal(tx)
	// if err != nil {
	// 	fmt.Printf("marshal tx failed;unexpected error:%v", err)
	// 	return nil, err
	// }
	// fmt.Printf("signedTx = %s\n", hex.EncodeToString(signedTx))
	rawTx, err := tronadpter.InsertSignatureIntoRawTransaction(txRawHex, hex.EncodeToString(sign))
	// fmt.Printf("rawTx = %s\n", rawTx)

	// if rawTx == hex.EncodeToString(signedTx) {
	// 	fmt.Printf("rawTx = %s", "相等")
	// } else {
	// 	fmt.Printf("rawTx = %s", "不相等")

	// }
	if err != nil {
		return nil, err
	}
	return &SignResult{
		Res:    1,
		RawTX:  rawTx,
		TxHash: hex.EncodeToString(txHash),
	}, nil

	// ************************

}

func (c *trx) GenerateTxHash(signIn *SignInput) (*TxHashResult, error) {
	amount := strconv.FormatInt(signIn.Amount, 10)
	toAddressBytes, err := addressEncoder.AddressDecode(signIn.DestAddr, addressEncoder.TRON_mainnetAddress)
	if err != nil {
		fmt.Printf("toAddress decode failed failed;unexpected error:%v", err)
		return nil, err
	}
	toAddressBytes = append([]byte{0x41}, toAddressBytes...)
	//fmt.Printf("Prefix  =  %s\n", hex.EncodeToString(addressEncoder.TRON_mainnetAddress.Prefix))

	ownerAddressBytes, err := addressEncoder.AddressDecode(signIn.SrcAddr, addressEncoder.TRON_mainnetAddress)
	if err != nil {
		fmt.Printf("ownerAddress decode failed failed;unexpected error:%v", err)
		return nil, err
	}
	ownerAddressBytes = append([]byte{0x41}, ownerAddressBytes...)
	fmt.Printf("ownerAddressBytes  =  %s\n", hex.EncodeToString(ownerAddressBytes))
	amountDec := address.StringNumToBigIntWithExp(amount, Decimals)

	block, err := tronadpter.GetNowBlock()
	if err != nil {
		return nil, err
	}
	blockID, err := hex.DecodeString(block.GetBlockHashID())

	refBlockBytes, refBlockHash := blockID[6:8], blockID[8:16]
	timestamp := time.Now().UnixNano() / 1000000 // <int64
	contracts := []*core.Transaction_Contract{}

	if signIn.ContractAddr == "" {
		// Check amount: amount * 1000000
		// ******** Generate TX Contract ********
		tc := &core.TransferContract{
			OwnerAddress: ownerAddressBytes,
			ToAddress:    toAddressBytes,
			Amount:       amountDec.Int64(),
		}

		tcRaw, err := proto.Marshal(tc)
		if err != nil {
			return nil, err
		}

		txContact := &core.Transaction_Contract{
			Type:         core.Transaction_Contract_TransferContract,
			Parameter:    &any.Any{Value: tcRaw, TypeUrl: "type.googleapis.com/protocol.TransferContract"},
			Provider:     nil,
			ContractName: nil,
		}
		contracts = []*core.Transaction_Contract{txContact}

	} else if strings.ToLower(signIn.Type) == "trc10" {
		tc := &core.TransferAssetContract{
			OwnerAddress: ownerAddressBytes,
			ToAddress:    toAddressBytes,
			Amount:       amountDec.Int64(),
			AssetName:    []byte(signIn.ContractAddr),
		}
		tcRaw, err := proto.Marshal(tc)
		if err != nil {
			return nil, err
		}

		txContact := &core.Transaction_Contract{
			Type:         core.Transaction_Contract_TransferAssetContract,
			Parameter:    &any.Any{Value: tcRaw, TypeUrl: "type.googleapis.com/protocol.TransferAssetContract"},
			Provider:     nil,
			ContractName: nil,
		}
		contracts = []*core.Transaction_Contract{txContact}

	} else if strings.ToLower(signIn.Type) == "trc20" {
		contractAddressBytes, err := addressEncoder.AddressDecode(signIn.ContractAddr, addressEncoder.TRON_mainnetAddress)
		contractAddressBytes = append([]byte{0x41}, contractAddressBytes...)

		if err != nil {
			return nil, err
		}

		var funcParams []tron.SolidityParam
		funcParams = append(funcParams, tron.SolidityParam{
			ParamType:  tron.SOLIDITY_TYPE_ADDRESS,
			ParamValue: hex.EncodeToString(toAddressBytes),
		})

		funcParams = append(funcParams, tron.SolidityParam{
			ParamType:  tron.SOLIDITY_TYPE_UINT256,
			ParamValue: amountDec,
		})

		//fmt.Println("make token transfer data, amount:", amount.String())
		dataHex, err := tron.MakeTransactionParameter(tron.TRC20_TRANSFER_METHOD_ID, funcParams)
		if err != nil {
			return nil, err
		}

		data, err := hex.DecodeString(dataHex)
		if err != nil {
			return nil, err
		}

		tc := &core.TriggerSmartContract{
			OwnerAddress:    ownerAddressBytes,
			ContractAddress: contractAddressBytes,
			Data:            data,
		}

		tcRaw, err := proto.Marshal(tc)
		if err != nil {
			return nil, err
		}

		txContact := &core.Transaction_Contract{
			Type:         core.Transaction_Contract_TriggerSmartContract,
			Parameter:    &any.Any{Value: tcRaw, TypeUrl: "type.googleapis.com/protocol.TriggerSmartContract"},
			Provider:     nil,
			ContractName: nil,
		}
		contracts = []*core.Transaction_Contract{txContact}
	}

	// ******** Create Traction ********
	txRaw := &core.TransactionRaw{
		RefBlockBytes: refBlockBytes,
		RefBlockHash:  refBlockHash,
		Contract:      contracts,
		Expiration:    timestamp + 10*60*60*60,
		// Timestamp:     timestamp,
	}
	txRaw.FeeLimit = 100000
	if signIn.Fee > 0 {
		txRaw.FeeLimit = signIn.Fee
	}
	tx := &core.Transaction{
		RawData: txRaw,
		// Signature: nil,
		// Ret:       nil,
	}
	// ******** TX Encoding ********
	x, err := proto.Marshal(tx)
	if err != nil {
		fmt.Printf("marshal tx failed;unexpected error:%v", err)
		return nil, err
	}
	txRawHex := hex.EncodeToString(x)
	txHash, err := getTxHash1(txRawHex)
	if err != nil {
		fmt.Printf("marshal tx failed;unexpected error:%v", err)
		return nil, err
	}
	fmt.Println("txHash = ", hex.EncodeToString(txHash))
	return &TxHashResult{
		ResCode:  1,
		TxHash:   hex.EncodeToString(txHash),
		TxRawHex: txRawHex,
	}, nil
}
func (c *trx) SignTxHash(signIn *SignTxHashInput) (*TxHashResult, error) {
	rawTx, err := tronadpter.InsertSignatureIntoRawTransaction(signIn.TxRawHex, signIn.Signature)
	if err != nil {
		return nil, err
	}
	return &TxHashResult{
		ResCode: 1,
		RawTX:   rawTx,
	}, nil
}

func getTxHash1(txHex string) ([]byte, error) {
	tx := &core.Transaction{}
	/*
		 if txRawBts, err := hex.DecodeString(tx.GetRawData()); err != nil {
			 return nil, err
		 } else {
			 if err := proto.Unmarshal(txRawBts, tx); err != nil {
				 return signedTxRaw, err
			 }
		 }
	*/
	txByte, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, fmt.Errorf("get Tx hex failed;unexpected error: %v", err)
	}
	if err := proto.Unmarshal(txByte, tx); err != nil {
		return nil, fmt.Errorf("unmarshal RawData failed; unexpected error: %v", err)
	}
	txRaw, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("marshal RawData failed;unexpected error:%v", err)
	}
	txHash := owcrypt.Hash(txRaw, 0, owcrypt.HASH_ALG_SHA256)
	return txHash, nil
}
