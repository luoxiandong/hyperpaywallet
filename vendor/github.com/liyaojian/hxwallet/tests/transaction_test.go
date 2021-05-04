package tests

import (
	"fmt"
	"github.com/liyaojian/hxwallet/api/jsonrpc"
	"github.com/liyaojian/hxwallet/rpc/http"
	"github.com/liyaojian/hxwallet/types"
	"github.com/liyaojian/hxwallet/util"
	"testing"
)

const url = "http://47.88.171.179:8090"

func Test_Transfer(t *testing.T) {
	refInfo := util.CalRefInfo("00a8f587e97017317eefd151da7742cedb5dc29b") // https://docs.gxchain.org/zh/guide/apis.html#get-dynamic-global-properties 获取 head_block_id
	wif := "5KNkNu1GYrehgLcSz2GMznVoQHfmce8t24ZTPEGoygvHbtgGW7k"           // 私钥
	fromAddr := "HXNb7KtLSX5yaj4111beUh486uKiodUZkL9J"
	toAddr := "HXNV81CkWqjivD2L3bBvdDpcmo7qS9nq44q4"
	chainId := "2e13ba07b457f2e284dcfcbd3d4a3e4d78a6ed89a61006cdb7fdad6d67ef0b12" // 链ID可通过 https://docs.gxchain.org/zh/guide/apis.html#get-chain-id 获取
	coinType := "HX"
	amount := int64(1000)
	fee := int64(200) // 手续费可通过 https://docs.gxchain.org/zh/guide/apis.html#get-required-fees 获取
	memo := "hello world"

	trxData, err := types.Transfer(refInfo, wif, chainId, fromAddr, toAddr, coinType, amount, fee, memo, "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(trxData))

	// 广播地址：http://nodeapi.hxlab.org/
}

func Test_GetChainId(t *testing.T) {
	db := jsonrpc.NewAPI("database", http.NewTransport(url))
	id, err := db.GetChainId()
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}

func Test_GetDGP(t *testing.T) {
	db := jsonrpc.NewAPI("database", http.NewTransport(url))
	re, err := db.GetDynamicGlobalProperties()
	if err != nil {
		panic(err)
	}
	fmt.Println(re.HeadBlockID)
}
