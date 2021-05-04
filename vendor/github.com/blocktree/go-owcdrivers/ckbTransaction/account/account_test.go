package account

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_AccountDecode(t *testing.T) {
	fmt.Println(len("ckb1qyq85a5m8aqv5mn03y6h4n4d4wuk36vn289sr7y7m9"))

	res, pkhash, err := DecodeSegWitAddress("ckb1qyq85a5m8aqv5mn03y6h4n4d4wuk36vn289sr7y7m9")

	fmt.Println("res = ", res, " from hash = ", hex.EncodeToString(pkhash), "  err = ", err)

	res, pkhash, err = DecodeSegWitAddress("ckb1qyq9gwznvqxjdsnk832cxps9gnf4nxmv3zfqmj6kh2")

	fmt.Println("res = ", res, " to hash = ", hex.EncodeToString(pkhash), "  err = ", err)

}

func Test_AccountAddress(t *testing.T) {

	hash, _ := hex.DecodeString("b39bbc0b3673c7d36450bc14cfcdad2d559c6c64")
	prefix := "ckb"
	bech32Addr, err := EncodeSegWitAddress(prefix, hash)
	if err != nil {
		panic(err)
	}
	fmt.Println(bech32Addr)

	fmt.Println("ckb1qyqt8xaupvm8837nv3gtc9x0ekkj64vud3jqfwyw5v")
	if bech32Addr == "ckb1qyqt8xaupvm8837nv3gtc9x0ekkj64vud3jqfwyw5v" {
		fmt.Println("相等")
	} else {
		fmt.Println("不相等")
	}

}
