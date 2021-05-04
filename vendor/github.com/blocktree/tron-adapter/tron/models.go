/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package tron

import (
	"encoding/hex"
	"fmt"
	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"strconv"
)

type Block struct {
	/*
		 {
			 "blockID":"000000000035e1c0f60afaa8387fd17fd9b84fe4381265ff084d739f814558ea",
			 "block_header":{
				 "raw_data":{"number":3531200,
						 "txTrieRoot":"9d98f6cbbde8302774ab87c003831333e132c89e00009cb1f7da35e1e59ae8ca",
						 "witness_address":"41d1dbde8b8f71b48655bec4f6bb532a0142b88bc0",
						 "parentHash":"000000000035e1bfb3b6f244e5316ce408aa8cea4c348eabe2545247f5a4600c",
						 "version":3,
						 "timestamp":1540545240000},
				 "witness_signature":"6ceedcacd8d0111b48eb4131484de3d13f27a2f4bd8156279d03d4208690158e20a641b77d900e026ee33adc328f9ec674f6483ea7b1ca5a27fa24d7fb23964100"
			 },
			 "transactions":[
				 {"ret":[{"contractRet":"SUCCESS"}],
				  "signature":["40aa520f01cebf12948615b9c5a5df5fe7d57e1a7f662d53907b4aa14f647a3a47be2a097fdb58159d0bee7eb1ff0a15ac738f24643fe5114cab8ec0d52cc04d01"],
				  "txID":"ac005b0a195a130914821a6c28db1eec44b4ec3a2358388ceb6c87b677866f1f",
				  "raw_data":{
					  "contract":[
						 {"parameter":{"value":{"amount":1,"asset_name":"48756f6269546f6b656e","owner_address":"416b201fb7b9f2b97bbdaf5e0920191229767c30ee","to_address":"412451d09536fca47760ea6513372bbbbef8583105"},
								   "type_url":"type.googleapis.com/protocol.TransferAssetContract"},
						  "type":"TransferAssetContract"}
					  ],
					  "ref_block_bytes":"e1be",
					  "ref_block_hash":"8dbf5f0cf4c324f2",
					  "expiration":1540545294000,
					  "timestamp":1540545235358}
				 },
				 ...]
		 }
	*/
	Hash              string // 这里采用 BlockID
	tx                []*Transaction
	Previousblockhash string
	Height            uint64 `storm:"id"`
	Version           uint64
	Time              int64
	Fork              bool
	txHash            []string
	// Merkleroot        string
	//Confirmations     uint64
}

func NewBlock(json *gjson.Result, isTestnet bool) *Block {

	header := gjson.Get(json.Raw, "block_header").Get("raw_data")
	// 解析json
	b := &Block{}
	b.Hash = gjson.Get(json.Raw, "blockID").String()
	b.Previousblockhash = header.Get("parentHash").String()
	b.Height = header.Get("number").Uint()
	b.Version = header.Get("version").Uint()
	b.Time = header.Get("timestamp").Int()

	txs := make([]*Transaction, 0)
	for _, x := range gjson.Get(json.Raw, "transactions").Array() {
		tx := NewTransaction(&x, b.Hash, b.Height, b.Time, isTestnet)
		txs = append(txs, tx)
	}
	b.tx = txs
	return b
}

func (block *Block) GetHeight() uint64 {
	return block.Height
}

func (block *Block) GetBlockHashID() string {
	return block.Hash
}

func (block *Block) GetTransactions() []*Transaction {
	return block.tx
}

//状态
const (
	SUCCESS = "SUCCESS"
)

//交易单类型
const (
	TransferContract      = "TransferContract"
	TransferAssetContract = "TransferAssetContract"
	TriggerSmartContract  = "TriggerSmartContract"
)

type Result struct {
	Ret         string
	Fee         int64
	ContractRet string
}

func NewResult(json gjson.Result) *Result {
	// 解析json
	b := &Result{}
	b.Ret = gjson.Get(json.Raw, "ret").String()
	b.Fee = gjson.Get(json.Raw, "fee").Int()
	b.ContractRet = gjson.Get(json.Raw, "contractRet").String()
	return b
}

type Contract struct {
	TxID            string
	BlockHash       string
	BlockHeight     uint64
	BlockTime       int64
	Type            string
	Parameter       gjson.Result
	Provider        []byte
	ContractName    []byte
	From            string
	To              string
	Amount          int64
	ContractAddress string
	SourceKey       string
	ContractRet     string
	Protocol        string
}

func NewContract(json gjson.Result, isTestnet bool) *Contract {

	var err error

	// 解析json
	b := &Contract{}
	b.Type = gjson.Get(json.Raw, "type").String()
	b.Parameter = gjson.Get(json.Raw, "parameter")
	switch b.Type {
	case TransferContract: //TRX

		b.From, err = EncodeAddress(b.Parameter.Get("value.owner_address").String(), isTestnet)
		if err != nil {
			return &Contract{}
		}
		b.To, err = EncodeAddress(b.Parameter.Get("value.to_address").String(), isTestnet)
		if err != nil {
			return &Contract{}
		}
		b.Amount = b.Parameter.Get("value.amount").Int()
		b.Protocol = ""
	case TransferAssetContract: //TRC10

		b.From, err = EncodeAddress(b.Parameter.Get("value.owner_address").String(), isTestnet)
		if err != nil {
			return &Contract{}
		}
		b.To, err = EncodeAddress(b.Parameter.Get("value.to_address").String(), isTestnet)
		if err != nil {
			return &Contract{}
		}
		b.Amount = b.Parameter.Get("value.amount").Int()
		assetsByte, err := hex.DecodeString(b.Parameter.Get("value.asset_name").String())
		if err != nil {
			return &Contract{}
		}
		b.ContractAddress = string(assetsByte)
		b.Protocol = TRC10
	case TriggerSmartContract: //TRC20

		b.From, err = EncodeAddress(b.Parameter.Get("value.owner_address").String(), isTestnet)
		if err != nil {
			return &Contract{}
		}

		b.ContractAddress, err = EncodeAddress(b.Parameter.Get("value.contract_address").String(), isTestnet)
		if err != nil {
			return &Contract{}
		}

		data := b.Parameter.Get("value.data").String()
		to, amount, err := ParseTransferEvent(data)
		if err != nil {
			return &Contract{}
		}
		b.To, err = EncodeAddress(to, isTestnet)
		if err != nil {
			return &Contract{}
		}

		b.Amount = amount
		b.Protocol = TRC20
	}

	return b
}

func ParseTransferEvent(data string) (string, int64, error) {

	const (
		TransferDataLength = 68 * 2
	)

	var (
		to     string
		amount int64
	)

	if len(data) != TransferDataLength {
		return "", 0, fmt.Errorf("call data is not transfer")
	}
	if data[0:8] != TRC20_TRANSFER_METHOD_ID {
		return "", 0, fmt.Errorf("call method is not transfer")
	}
	to = data[32:72]
	amount, err := strconv.ParseInt(data[72:], 16, 64)
	if err != nil {
		return "", 0, err
	}
	return to, amount, nil
}

type Transaction struct {
	/*
		{
		    "ret": [
		        {
		            "contractRet": "SUCCESS"
		        }
		    ],
		    "signature": [
		        "b6fbefec182a1db9759104e3d5709343f25acc1866c77deb050d9bf0f6a5c3342237fbbf076377ea26d1a09edd46e599534a8499998814da0e9187e07260382001"
		    ],
		    "txID": "86b5a123b5cc50047532f1a55ed627f29012bba41e6590b0545f903289e7099a",
		    "raw_data": {
		        "contract": [
		            {
		                "parameter": {
		                    "value": {
		                        "amount": 10000,
		                        "owner_address": "415bdf283199369adb124f39dda845ae02c5d3eb5d",
		                        "to_address": "41bb65606e20dbdd7f0cfe6a66ae2858a6534f2d45"
		                    },
		                    "type_url": "type.googleapis.com/protocol.TransferContract"
		                },
		                "type": "TransferContract"
		            }
		        ],
		        "ref_block_bytes": "2298",
		        "ref_block_hash": "26478922abba48a0",
		        "expiration": 1551957595231
		    }
		}
	*/
	TxID        string
	BlockHash   string
	BlockHeight uint64
	BlockTime   int64
	IsCoinBase  bool
	Ret         []*Result
	Contract    []*Contract
}

func NewTransaction(json *gjson.Result, blockHash string, blockHeight uint64, blocktime int64, isTestnet bool) *Transaction {

	rawData := gjson.Get(json.Raw, "raw_data")
	// 解析json
	b := &Transaction{}
	b.TxID = gjson.Get(json.Raw, "txID").String()
	b.BlockHash = blockHash
	b.BlockHeight = blockHeight
	b.BlockTime = blocktime

	b.Ret = make([]*Result, 0)
	if rets := gjson.Get(json.Raw, "ret"); rets.IsArray() {
		for _, r := range rets.Array() {
			ret := NewResult(r)
			b.Ret = append(b.Ret, ret)
		}
	}

	b.Contract = make([]*Contract, 0)
	if contracts := rawData.Get("contract"); contracts.IsArray() {
		for i, c := range contracts.Array() {
			contract := NewContract(c, isTestnet)
			contract.TxID = b.TxID
			contract.BlockHash = blockHash
			contract.BlockHeight = blockHeight
			contract.BlockTime = blocktime
			if len(b.Ret) > i {
				contract.ContractRet = b.Ret[i].ContractRet
			}
			b.Contract = append(b.Contract, contract)
		}
	}

	return b
}

type TransactionExtention struct {
	Transaction    gjson.Result `json:"transaction" `
	Txid           string       `json:"txid"`
	ConstantResult []string     `json:"constant_result"`
	Result         *Return      `json:"result"`
}

type Return struct {
	Result  bool   `json:"result"`
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func NewReturn(json *gjson.Result) *Return {
	b := &Return{}
	b.Result = json.Get("result").Bool()
	b.Code = json.Get("code").Int()
	msg, _ := hex.DecodeString(json.Get("message").String())
	b.Message = string(msg)
	return b
}

func NewTransactionExtention(json *gjson.Result) *TransactionExtention {

	// 解析json
	b := &TransactionExtention{}
	b.Transaction = json.Get("transaction")
	result := json.Get("result")
	b.Result = NewReturn(&result)
	b.Txid = json.Get("txid").String()

	b.ConstantResult = make([]string, 0)
	if constant_result := json.Get("constant_result"); constant_result.IsArray() {
		for _, c := range constant_result.Array() {
			b.ConstantResult = append(b.ConstantResult, c.String())
		}
	}

	return b
}

//UnscanRecords 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	TxID        string
	Reason      string
}

func NewUnscanRecord(height uint64, txID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.TxID = txID
	obj.Reason = reason
	obj.ID = hex.EncodeToString(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, txID))))
	return &obj
}

//BlockHeader 区块链头
func (b *Block) Blockheader() *openwallet.BlockHeader {
	obj := openwallet.BlockHeader{}
	//解析josn
	obj.Merkleroot = ""
	obj.Hash = b.Hash
	obj.Previousblockhash = b.Previousblockhash
	obj.Height = b.Height
	obj.Version = uint64(b.Version)
	obj.Time = uint64(b.Time)
	obj.Symbol = Symbol
	return &obj
}

type AccountNet struct {
	FreeNetUsed    int64
	FreeNetLimit   int64
	NetUsed        int64
	NetLimit       int64
	AssetNetUsed   map[string]int64
	AssetNetLimit  map[string]int64
	TotalNetLimit  int64
	TotalNetWeight int64
}

func NewAccountNet(json *gjson.Result) *AccountNet {
	obj := &AccountNet{}
	obj.FreeNetUsed = json.Get("freeNetUsed").Int()
	obj.FreeNetLimit = json.Get("freeNetLimit").Int()
	obj.NetUsed = json.Get("NetUsed").Int()
	obj.TotalNetLimit = json.Get("TotalNetLimit").Int()
	obj.TotalNetWeight = json.Get("TotalNetWeight").Int()
	return obj
}

type AccountResource struct {
	FreeNetUsed  int64
	FreeNetLimit int64
	NetUsed      int64
	NetLimit     int64
	EnergyUsed   int64
	EnergyLimit  int64
	StorageUsed  int64
	StorageLimit int64
}

func NewAccountResource(json *gjson.Result) *AccountResource {
	obj := &AccountResource{}
	obj.FreeNetUsed = json.Get("freeNetUsed").Int()
	obj.FreeNetLimit = json.Get("freeNetLimit").Int()
	obj.NetUsed = json.Get("NetUsed").Int()
	obj.NetLimit = json.Get("NetLimit").Int()
	obj.EnergyUsed = json.Get("EnergyUsed").Int()
	obj.EnergyLimit = json.Get("EnergyLimit").Int()
	obj.StorageUsed = json.Get("StorageUsed").Int()
	obj.StorageLimit = json.Get("StorageLimit").Int()
	return obj
}

type txFeeInfo struct {
	GasUsed  int64
	GasPrice decimal.Decimal
	Fee      decimal.Decimal
}

func (f *txFeeInfo) CalcFee() error {
	fee := f.GasPrice.Mul(decimal.New(f.GasUsed, 0))
	f.Fee = fee
	return nil
}

type ContractInfo struct {
	Bytecode                   string
	Name                       string
	ConsumeUserResourcePercent uint64
	ContractAddress            string
	ABI                        string
}

func NewContractInfo(json *gjson.Result) *ContractInfo {
	obj := &ContractInfo{}
	obj.Bytecode = json.Get("bytecode").String()
	obj.Name = json.Get("name").String()
	obj.ConsumeUserResourcePercent = json.Get("consume_user_resource_percent").Uint()
	obj.ContractAddress = json.Get("contract_address").String()
	obj.ABI = json.Get("abi.entrys").Raw
	return obj
}

type Account struct {
	AddressHex          string
	Balance             int64
	FreeNetUsage        int64
	AssetV2             map[string]int64
	FreeAssetNetUsageV2 map[string]int64

	/*
		{
		    "address": "41efb6d8a02f4b639605d71ff8dc78c97329759d70",
		    "balance": 1058035120983,
		    "create_time": 1546868937000,
		    "latest_opration_time": 1553592429000,
		    "free_net_usage": 4763,
		    "latest_consume_free_time": 1553590113000,
		    "account_resource": {
		        "latest_consume_time_for_energy": 1553590494000
		    },
		    "assetV2": [
		        {
		            "key": "1001064",
		            "value": 10
		        }
		    ],
		    "free_asset_net_usageV2": [
		        {
		            "key": "1001064",
		            "value": 0
		        }
		    ]
		}

	*/
}

func NewAccount(json *gjson.Result) *Account {
	obj := &Account{}
	obj.AddressHex = json.Get("address").String()
	obj.Balance = json.Get("balance").Int()
	obj.FreeNetUsage = json.Get("free_net_usage").Int()

	obj.AssetV2 = make(map[string]int64, 0)
	assetV2 := json.Get("assetV2")
	if assetV2.IsArray() {
		for _, as := range assetV2.Array() {
			obj.AssetV2[as.Get("key").String()] = as.Get("value").Int()
		}
	}

	obj.FreeAssetNetUsageV2 = make(map[string]int64, 0)
	freeAssetNetUsageV2 := json.Get("free_asset_net_usageV2")
	if freeAssetNetUsageV2.IsArray() {
		for _, as := range freeAssetNetUsageV2.Array() {
			obj.FreeAssetNetUsageV2[as.Get("key").String()] = as.Get("value").Int()
		}
	}
	return obj
}
