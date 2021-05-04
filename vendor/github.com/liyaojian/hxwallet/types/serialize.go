/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: serialize.go
 * Date: 2018-09-07
 *
 */

package types

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/liyaojian/hxwallet/util"
)

// inferface for serialize hx transaction
type HxSearilze interface {
	Serialize() []byte
}

func PackUint16(val uint16, isLittleEndian bool) []byte {

	res := make([]byte, 2)

	if isLittleEndian {
		binary.LittleEndian.PutUint16(res, val)
	} else {
		binary.BigEndian.PutUint16(res, val)
	}

	return res

}

func UnPackUint16(bytes []byte, isLittleEndian bool) uint16 {

	var res uint16

	if isLittleEndian {
		res = binary.LittleEndian.Uint16(bytes)
	} else {
		res = binary.BigEndian.Uint16(bytes)
	}

	return res
}

func PackUint32(val uint32, isLittleEndian bool) []byte {

	res := make([]byte, 4)

	if isLittleEndian {
		binary.LittleEndian.PutUint32(res, val)
	} else {
		binary.BigEndian.PutUint32(res, val)
	}

	return res

}

func UnPackUint32(bytes []byte, isLittleEndian bool) uint32 {

	var res uint32

	if isLittleEndian {
		res = binary.LittleEndian.Uint32(bytes)
	} else {
		res = binary.BigEndian.Uint32(bytes)
	}

	return res
}

func PackInt64(val int64, isLittleEndian bool) []byte {

	res := make([]byte, 8)

	if isLittleEndian {
		binary.LittleEndian.PutUint64(res, uint64(val))
	} else {
		binary.BigEndian.PutUint64(res, uint64(val))
	}

	return res
}

func UnPackInt64(bytes []byte, isLittleEndian bool) int64 {

	var res int64

	if isLittleEndian {
		res = int64(binary.LittleEndian.Uint64(bytes))
	} else {
		res = int64(binary.BigEndian.Uint64(bytes))
	}

	return res
}

func PackVarUint32(val uint32) []byte {

	res := make([]byte, 0)

	//one byte
	if val < 0x80 {

		res = append(res, byte(val))

		return res
	} else if val < 0x4000 { //two byte

		byte1 := val / 0x80
		byte2 := val%0x80 + 0x80

		res = append(res, byte(byte2))
		res = append(res, byte(byte1))

	} else if val < 0x200000 { //three byte

		byte1 := val / 0x4000
		byte2 := val%0x4000/0x80 + 0x80
		byte3 := val%0x80 + 0x80

		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))

	} else if val < 0x10000000 { //four byte

		byte1 := val / 0x200000
		byte2 := val%0x200000/0x4000 + 0x80
		byte3 := val%0x4000/0x80 + 0x80
		byte4 := val%0x80 + 0x80

		res = append(res, byte(byte4))
		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))
	} else {

		byte1 := val / 0x10000000
		byte2 := val%0x10000000/0x200000 + 0x80
		byte3 := val%0x200000/0x4000 + 0x80
		byte4 := val%0x4000/0x80 + 0x80
		byte5 := val%0x80 + 0x80

		res = append(res, byte(byte5))
		res = append(res, byte(byte4))
		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))

	}

	return res
}

func (asset *Asset) Serialize() []byte {

	byte_int64 := PackInt64(asset.HxAmount, true)

	//byte for asset_id_type, default to zero
	tmp_id, err := util.GetId(asset.HxAssetId)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 := PackVarUint32(tmp_id)
	byte_int64 = append(byte_int64, byte_uint32...)

	return byte_int64
}

func (memo *Memo) Serialize() []byte {

	if memo == nil {
		return []byte{0}
	} else {

		//byte for optional, have element default to one
		var res []byte
		res = append(res, byte(1))
		byte_pub := make([]byte, 74)
		res = append(res, byte_pub...)
		// memo message
		res = append(res, byte(len(memo.Message)+4))
		byte_pub = make([]byte, 4)
		res = append(res, byte_pub...)
		res = append(res, []byte(memo.Message)...)
		return res

	}

}

func (authority *Authority) Serialize() []byte {

	var res []byte
	res = append(res, PackUint32(authority.HxWeightThreshold, true)...)
	res = append(res, byte(0))
	res = append(res, byte(len(authority.HxKeyAuths)))
	tmpByte, _ := util.GetPubkeyBytes(authority.KeyAuths)
	res = append(res, tmpByte...)
	res = append(res, PackUint16(1, true)...)
	res = append(res, byte(0))

	return res
}

func (acc *AccountOptions) Serialize() []byte {

	var res []byte
	tmpByte, _ := util.GetPubkeyBytes(acc.HxMemoKey)
	res = append(res, tmpByte...)
	res = append(res, byte(5))
	res = append(res, PackUint16(0, true)...)
	res = append(res, PackUint16(0, true)...)
	res = append(res, byte(0))
	res = append(res, byte(10))
	res = append(res, byte(0))

	return res
}

func (tranferOp *TransferOperation) Serialize() []byte {

	res := tranferOp.HxFee.Serialize()

	if tranferOp.HxGuaranteeId != "" {
		res = append(res, byte(1))
		tmp_id, err := util.GetId(tranferOp.HxGuaranteeId)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)

		byteTmp := make([]byte, 2)
		res = append(res, byteTmp...)

	} else {
		byteTmp := make([]byte, 3)
		res = append(res, byteTmp...)
	}

	byteTmp, _ := util.GetAddressBytes(tranferOp.HxFromAddr)
	res = append(res, byteTmp...)
	byteTmp, _ = util.GetAddressBytes(tranferOp.HxToAddr)
	res = append(res, byteTmp...)

	byteTmp = tranferOp.HxAmount.Serialize()
	res = append(res, byteTmp...)

	byteTmp = tranferOp.HxMemo.Serialize()
	res = append(res, byteTmp...)
	res = append(res, byte(0))

	return res

}

func (bindOp *AccountBindOperation) Serialize() []byte {

	res := bindOp.HxFee.Serialize()
	res = append(res, byte(len(bindOp.HxCrosschainType)))
	res = append(res, []byte(bindOp.HxCrosschainType)...)
	tmpByte, _ := util.GetAddressBytes(bindOp.HxAddr)
	res = append(res, tmpByte...)
	tmpByte, _ = hex.DecodeString(bindOp.HxAccountSignature)
	res = append(res, tmpByte...)
	res = append(res, byte(len(bindOp.HxTunnelAddress)))
	res = append(res, []byte(bindOp.HxTunnelAddress)...)

	tmpByte = PackVarUint32(uint32(len(bindOp.HxTunnelSignature)))
	res = append(res, tmpByte...)
	res = append(res, []byte(bindOp.HxTunnelSignature)...)

	if bindOp.HxGuaranteeId != "" {
		res = append(res, byte(1))
		tmp_id, err := util.GetId(bindOp.HxGuaranteeId)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	return res

}

func (unbindOp *AccountUnBindOperation) Serialize() []byte {

	res := unbindOp.HxFee.Serialize()
	res = append(res, byte(len(unbindOp.HxCrosschainType)))
	res = append(res, []byte(unbindOp.HxCrosschainType)...)
	tmpByte, _ := util.GetAddressBytes(unbindOp.HxAddr)
	res = append(res, tmpByte...)
	tmpByte, _ = hex.DecodeString(unbindOp.HxAccountSignature)
	res = append(res, tmpByte...)
	res = append(res, byte(len(unbindOp.HxTunnelAddress)))
	res = append(res, []byte(unbindOp.HxTunnelAddress)...)
	tmpByte = PackVarUint32(uint32(len(unbindOp.HxTunnelSignature)))
	res = append(res, tmpByte...)
	res = append(res, []byte(unbindOp.HxTunnelSignature)...)

	res = append(res, byte(0))

	return res

}

func (withdraw *WithdrawCrosschainOperation) Serialize() []byte {

	var res []byte
	res = append(res, withdraw.HxFee.Serialize()...)
	tmpByte, _ := util.GetAddressBytes(withdraw.HxWithdrawAccount)
	res = append(res, tmpByte...)
	res = append(res, byte(len(withdraw.HxAmount)))
	res = append(res, []byte(withdraw.HxAmount)...)
	res = append(res, byte(len(withdraw.HxAssetSymbol)))
	res = append(res, []byte(withdraw.HxAssetSymbol)...)

	//byte for asset_id_type, default to zero
	tmp_id, err := util.GetId(withdraw.HxAssetId)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 := PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)

	res = append(res, byte(len(withdraw.HxCrosschainAccount)))
	res = append(res, []byte(withdraw.HxCrosschainAccount)...)
	res = append(res, byte(len(withdraw.HxMemo)))
	res = append(res, []byte(withdraw.HxMemo)...)

	return res
}

func (register *RegisterAccountOperation) Serialize() []byte {

	var res []byte
	res = append(res, register.HxFee.Serialize()...)

	tmpByte := make([]byte, 2)
	res = append(res, tmpByte...)
	tmpByte = PackUint16(0, true)
	res = append(res, tmpByte...)
	res = append(res, byte(len(register.HxName)))
	res = append(res, []byte(register.HxName)...)

	res = append(res, register.HxOwner.Serialize()...)
	res = append(res, register.HxActive.Serialize()...)

	tmpByte, _ = util.GetAddressBytes(register.HxPayer)
	res = append(res, tmpByte...)

	res = append(res, register.HxOptions.Serialize()...)
	res = append(res, byte(0))

	if register.HxGuaranteeId != "" {
		res = append(res, byte(1))
		tmp_id, err := util.GetId(register.HxGuaranteeId)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	return res
}

func (lockOp *LockBalanceOperation) Serialize() []byte {

	var res []byte
	tmp_id, err := util.GetId(lockOp.HxLockAssetId)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 := PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)
	res = append(res, PackInt64(lockOp.HxLockAssetAmount, true)...)

	//tmpByte, _ := GetAddressBytes(lockOp.Hx_contract_addr)
	//res = append(res, tmpByte...)
	var invalid_address_byte []byte
	invalid_address_byte = append(invalid_address_byte, byte(0x35))
	tmpByte := make([]byte, 20)
	invalid_address_byte = append(invalid_address_byte, tmpByte...)
	res = append(res, invalid_address_byte...)

	tmp_id, err = util.GetId(lockOp.HxLockBalanceAccount)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 = PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)

	tmp_id, err = util.GetId(lockOp.HxLocktoMinerAccount)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 = PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)

	tmpByte, _ = util.GetAddressBytes(lockOp.HxLockBalanceAddr)
	res = append(res, tmpByte...)

	res = append(res, lockOp.HxFee.Serialize()...)

	return res
}

func (obtainOp *ObtainPaybackOperation) Serialize() []byte {

	var res []byte

	tmpByte, _ := util.GetAddressBytes(obtainOp.HxPayBackOwner)
	res = append(res, tmpByte...)

	res = append(res, byte(len(obtainOp.HxPayBackBalance)))
	for i := 0; i < len(obtainOp.HxPayBackBalance); i++ {
		res = append(res, byte(len(obtainOp.citizenName[i])))
		res = append(res, []byte(obtainOp.citizenName[i])...)
		res = append(res, obtainOp.obtainAsset[i].Serialize()...)
	}

	if obtainOp.HxGuaranteeId != "" {
		res = append(res, byte(1))
		tmp_id, err := util.GetId(obtainOp.HxGuaranteeId)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	res = append(res, obtainOp.HxFee.Serialize()...)

	return res
}

func (forecloseOp *ForecloseBalanceOperation) Serialize() []byte {

	var res []byte
	res = append(res, forecloseOp.HxFee.Serialize()...)

	tmp_id, err := util.GetId(forecloseOp.HxForecloseAssetId)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 := PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)
	res = append(res, PackInt64(forecloseOp.HxForecloseAssetAmount, true)...)

	tmp_id, err = util.GetId(forecloseOp.HxForecloseMinerAccount)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 = PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)

	//tmpByte, _ := GetAddressBytes(forecloseOp.Hx_foreclose_contract_addr)
	//res = append(res, tmpByte...)
	var invalid_address_byte []byte
	invalid_address_byte = append(invalid_address_byte, byte(0x35))
	tmpByte := make([]byte, 20)
	invalid_address_byte = append(invalid_address_byte, tmpByte...)
	res = append(res, invalid_address_byte...)

	tmp_id, err = util.GetId(forecloseOp.HxForecloseAccount)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 = PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)
	tmpByte, _ = util.GetAddressBytes(forecloseOp.HxForecloseAccount)
	res = append(res, tmpByte...)

	return res
}

func (contractOp *ContractInvokeOperation) Serialize() []byte {

	var res []byte
	res = append(res, contractOp.HxFee.Serialize()...)

	byte_int64 := PackInt64(int64(contractOp.HxInvokeCost), true)
	res = append(res, byte_int64...)
	byte_int64 = PackInt64(int64(contractOp.HxGasPrice), true)
	res = append(res, byte_int64...)

	tmpByte, _ := util.GetAddressBytes(contractOp.HxCallerAddr)
	res = append(res, tmpByte...)
	tmpByte, _ = hex.DecodeString(contractOp.HxCallerPubkey)
	res = append(res, tmpByte...)
	tmpByte, _ = util.GetAddressBytes(contractOp.HxContractId)
	res = append(res, tmpByte...)
	res = append(res, byte(len(contractOp.HxContractApi)))
	res = append(res, []byte(contractOp.HxContractApi)...)
	res = append(res, byte(len(contractOp.HxContractArg)))
	res = append(res, []byte(contractOp.HxContractArg)...)

	if contractOp.HxGuaranteeId != "" {
		res = append(res, byte(1))
		tmp_id, err := util.GetId(contractOp.HxGuaranteeId)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	return res
}

func (contractOp *ContractTransferOperation) Serialize() []byte {
	var res []byte
	res = append(res, contractOp.HxFee.Serialize()...)

	byte_int64 := PackInt64(int64(contractOp.HxInvokeCost), true)
	res = append(res, byte_int64...)
	byte_int64 = PackInt64(int64(contractOp.HxGasPrice), true)
	res = append(res, byte_int64...)

	tmpByte, _ := util.GetAddressBytes(contractOp.HxCallerAddr)
	res = append(res, tmpByte...)
	tmpByte, _ = hex.DecodeString(contractOp.HxCallerPubkey)
	res = append(res, tmpByte...)
	tmpByte, _ = util.GetAddressBytes(contractOp.HxContractId)
	res = append(res, tmpByte...)

	res = append(res, contractOp.HxAmount.Serialize()...)
	res = append(res, byte(len(contractOp.HxParam)))
	res = append(res, []byte(contractOp.HxParam)...)

	if contractOp.HxGuaranteeId != "" {
		res = append(res, byte(1))
		tmp_id, err := util.GetId(contractOp.HxGuaranteeId)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	return res

}

func (trx *Transaction) Serialize() []byte {

	var res []byte
	res = append(res, PackUint16(trx.HxRefBlockNum, true)...)
	res = append(res, PackUint32(trx.HxRefBlockPrefix, true)...)
	res = append(res, PackUint32(trx.Expiration, true)...)

	//operations
	res = append(res, byte(len(trx.Operations)))
	for _, v := range trx.Operations {

		if transferOp, ok := v.(TransferOperation); ok {
			res = append(res, byte(0))
			res = append(res, transferOp.Serialize()...)
		} else if bindOp, ok := v.(AccountBindOperation); ok {
			res = append(res, byte(10))
			res = append(res, bindOp.Serialize()...)
		} else if unbindOp, ok := v.(AccountUnBindOperation); ok {
			res = append(res, byte(11))
			res = append(res, unbindOp.Serialize()...)
		} else if withdrawOp, ok := v.(WithdrawCrosschainOperation); ok {
			res = append(res, byte(61))
			res = append(res, withdrawOp.Serialize()...)
		} else if registerOp, ok := v.(RegisterAccountOperation); ok {
			res = append(res, byte(5))
			res = append(res, registerOp.Serialize()...)
		} else if lockOp, ok := v.(LockBalanceOperation); ok {
			res = append(res, byte(55))
			res = append(res, lockOp.Serialize()...)
		} else if forecloseOp, ok := v.(ForecloseBalanceOperation); ok {
			res = append(res, byte(56))
			res = append(res, forecloseOp.Serialize()...)
		} else if obtainOp, ok := v.(ObtainPaybackOperation); ok {
			res = append(res, byte(73))
			res = append(res, obtainOp.Serialize()...)
		} else if contractOp, ok := v.(ContractInvokeOperation); ok {
			res = append(res, byte(79))
			res = append(res, contractOp.Serialize()...)
		} else if contractOp, ok := v.(ContractTransferOperation); ok {
			res = append(res, byte(81))
			res = append(res, contractOp.Serialize()...)
		}

	}

	//extension
	res = append(res, byte(0))

	//signature
	if len(trx.HxSignatures) > 0 {
		res = append(res, byte(len(trx.HxSignatures)))
	}

	return res
}
