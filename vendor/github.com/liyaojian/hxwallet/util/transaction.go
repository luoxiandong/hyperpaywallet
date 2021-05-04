package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"strconv"
	"strings"
	"time"
)

func CalRefInfo(blockHash string) string {
	blockNum := blockHash[:8]
	refBlockId := blockHash[8:16]
	var refBlockNumLittle uint16
	blockNumBytes, _ := hex.DecodeString(blockNum)

	refBlockNumLittle = binary.BigEndian.Uint16(blockNumBytes[2:4])
	refBlockPrefix, _ := hex.DecodeString(refBlockId)

	endRefBlockPrefix := binary.LittleEndian.Uint32(refBlockPrefix)
	refStr := fmt.Sprintf("%d,%d", refBlockNumLittle, endRefBlockPrefix)

	return refStr
}

func GetRefblockInfo(info string) (uint16, uint32, error) {

	refinfo := strings.Split(info, ",")
	// refinfo := []string{"21771", "761216631"}

	if len(refinfo) != 2 {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, get refblockinfo failed")
	}
	refBlockNumStr, refBlockPrefixStr := refinfo[0], refinfo[1]
	refBlockNum, err := strconv.ParseUint(refBlockNumStr, 10, 16)
	if err != nil {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, convert ref_block_num failed: %v", err)
	}

	refBlockPrefix, err := strconv.ParseUint(refBlockPrefixStr, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, convert ref_block_prefix failed: %v", err)
	}

	return uint16(refBlockNum), uint32(refBlockPrefix), nil
}

func GetId(id string) (uint32, error) {

	idSlice := strings.Split(id, ".")

	if len(idSlice) != 3 {
		return 0, fmt.Errorf("in GetId function, get account id failed")
	}

	res, err := strconv.ParseUint(idSlice[2], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("in GetId function, Parse id error %v", err)
	}

	return uint32(res), nil

}

func Str2Time(str string) int64 {

	str += "Z"
	t, err := time.Parse(time.RFC3339, str)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return t.Unix()

}

func Time2Str(t int64) string {

	l_time := time.Unix(t, 0).UTC()
	timestr := l_time.Format(time.RFC3339)

	timestr = timestr[:len(timestr)-1]

	return timestr
}

// in multiple precision mode
func CalculateFee(basic_op_fee int64, len_memo int64) int64 {

	var basic_memo_fee int64 = 1
	return basic_op_fee + len_memo*basic_memo_fee
}

func GetAddressBytes(addr string) ([]byte, error) {

	if len(addr) <= 2 {
		return nil, fmt.Errorf("in GetAddressBytes function, wrong addr format")
	}

	base58Addr := addr[2:]

	addrBytes := base58.Decode(base58Addr)

	return addrBytes[:len(addrBytes)-4], nil
}

func GetPubkeyBytes(pub string) ([]byte, error) {

	if len(pub) <= 2 {
		return nil, fmt.Errorf("in GetAddressBytes function, wrong addr format")
	}

	base58Addr := pub[2:]

	pubBytes := base58.Decode(base58Addr)

	return pubBytes[:len(pubBytes)-4], nil
}
