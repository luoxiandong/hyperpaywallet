package zecutil

import (
	"fmt"

	"github.com/blake2b-simd"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func blake2bHash(data, key []byte) (h chainhash.Hash, err error) {
	config := &blake2b.Config{
		Size:   32,
		Person: key,
	}
	bHash, err1 := blake2b.New(config)
	if err1 != nil {
		fmt.Printf("Initial blake2b hash failure, %s\n", err1)
		return h, err1
	}

	// Hash tx
	bHash.Write(data)

	txHash := bHash.Sum(nil)

	err = (&h).SetBytes(txHash)
	return h, err
}

// // blake2bHash zcash hash func
// func blake2bHash(data, key []byte) (h chainhash.Hash, err error) {
// 	bHash := blake2.New(&blake2.Config{
// 		Size:     32,
// 		Personal: key,
// 	})

// 	if _, err = bHash.Write(data); err != nil {
// 		return h, err
// 	}

// 	err = (&h).SetBytes(bHash.Sum(nil))
// 	return h, err
// }
