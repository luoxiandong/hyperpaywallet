package zcash

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/minio/blake2b-simd"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	txHeaderBytes          = []byte{0x04, 0x00, 0x00, 0x80}
	txNVersionGroupIDBytes = []byte{0x85, 0x20, 0x2f, 0x89}

	hashPrevOutPersonalization  = []byte("ZcashPrevoutHash")
	hashSequencePersonalization = []byte("ZcashSequencHash")
	hashOutputsPersonalization  = []byte("ZcashOutputsHash")
	sigHashPersonalization      = []byte("ZcashSigHash")
)

const (
	sigHashMask = 0x1f
	branchID    = 0x2BB40E60
)

// rawTxInSignature returns the serialized ECDSA signature for the input idx of
// the given transaction, with hashType appended to it.
func RawTxInSignature(tx *wire.MsgTx, idx int, prevScriptBytes []byte,
	hashType txscript.SigHashType, key *btcec.PrivateKey, amt int64) ([]byte, error) {

	hash, err := calcSignatureHash(prevScriptBytes, hashType, tx, idx, amt, 0)
	if err != nil {
		return nil, err
	}
	signature, err := key.Sign(hash)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx input: %s", err)
	}

	return append(signature.Serialize(), byte(hashType)), nil
}

func calcSignatureHash(prevScriptBytes []byte, hashType txscript.SigHashType, tx *wire.MsgTx, idx int, amt int64, expiry uint32) ([]byte, error) {

	// As a sanity check, ensure the passed input index for the transaction
	// is valid.
	if idx > len(tx.TxIn)-1 {
		return nil, fmt.Errorf("idx %d but %d txins", idx, len(tx.TxIn))
	}

	// We'll utilize this buffer throughout to incrementally calculate
	// the signature hash for this transaction.
	var sigHash bytes.Buffer

	// Write header
	_, err := sigHash.Write(txHeaderBytes)
	if err != nil {
		return nil, err
	}

	// Write group ID
	_, err = sigHash.Write(txNVersionGroupIDBytes)
	if err != nil {
		return nil, err
	}

	// Next write out the possibly pre-calculated hashes for the sequence
	// numbers of all inputs, and the hashes of the previous outs for all
	// outputs.
	var zeroHash chainhash.Hash

	// If anyone can pay isn't active, then we can use the cached
	// hashPrevOuts, otherwise we just write zeroes for the prev outs.
	if hashType&txscript.SigHashAnyOneCanPay == 0 {
		sigHash.Write(calcHashPrevOuts(tx))
	} else {
		sigHash.Write(zeroHash[:])
	}

	// If the sighash isn't anyone can pay, single, or none, the use the
	// cached hash sequences, otherwise write all zeroes for the
	// hashSequence.
	if hashType&txscript.SigHashAnyOneCanPay == 0 &&
		hashType&sigHashMask != txscript.SigHashSingle &&
		hashType&sigHashMask != txscript.SigHashNone {
		sigHash.Write(calcHashSequence(tx))
	} else {
		sigHash.Write(zeroHash[:])
	}

	// If the current signature mode isn't single, or none, then we can
	// re-use the pre-generated hashoutputs sighash fragment. Otherwise,
	// we'll serialize and add only the target output index to the signature
	// pre-image.
	if hashType&sigHashMask != txscript.SigHashSingle &&
		hashType&sigHashMask != txscript.SigHashNone {
		sigHash.Write(calcHashOutputs(tx))
	} else if hashType&sigHashMask == txscript.SigHashSingle && idx < len(tx.TxOut) {
		var b bytes.Buffer
		wire.WriteTxOut(&b, 0, 0, tx.TxOut[idx])
		sigHash.Write(chainhash.DoubleHashB(b.Bytes()))
	} else {
		sigHash.Write(zeroHash[:])
	}

	// Write hash JoinSplits
	sigHash.Write(make([]byte, 32))

	// Write hash ShieldedSpends
	sigHash.Write(make([]byte, 32))

	// Write hash ShieldedOutputs
	sigHash.Write(make([]byte, 32))

	// Write out the transaction's locktime, and the sig hash
	// type.
	var bLockTime [4]byte
	binary.LittleEndian.PutUint32(bLockTime[:], tx.LockTime)
	sigHash.Write(bLockTime[:])

	// Write expiry
	var bExpiryTime [4]byte
	binary.LittleEndian.PutUint32(bExpiryTime[:], expiry)
	sigHash.Write(bExpiryTime[:])

	// Write valueblance
	sigHash.Write(make([]byte, 8))

	// Write the hash type
	var bHashType [4]byte
	binary.LittleEndian.PutUint32(bHashType[:], uint32(hashType))
	sigHash.Write(bHashType[:])

	// Next, write the outpoint being spent.
	sigHash.Write(tx.TxIn[idx].PreviousOutPoint.Hash[:])
	var bIndex [4]byte
	binary.LittleEndian.PutUint32(bIndex[:], tx.TxIn[idx].PreviousOutPoint.Index)
	sigHash.Write(bIndex[:])

	// Write the previous script bytes
	wire.WriteVarBytes(&sigHash, 0, prevScriptBytes)

	// Next, add the input amount, and sequence number of the input being
	// signed.
	var bAmount [8]byte
	binary.LittleEndian.PutUint64(bAmount[:], uint64(amt))
	sigHash.Write(bAmount[:])
	var bSequence [4]byte
	binary.LittleEndian.PutUint32(bSequence[:], tx.TxIn[idx].Sequence)
	sigHash.Write(bSequence[:])

	leBranchID := make([]byte, 4)
	binary.LittleEndian.PutUint32(leBranchID, branchID)
	bl, _ := blake2b.New(&blake2b.Config{
		Size:   32,
		Person: append(sigHashPersonalization, leBranchID...),
	})
	bl.Write(sigHash.Bytes())
	h := bl.Sum(nil)
	return h[:], nil
}

// serializeVersion4Transaction serializes a wire.MsgTx into the zcash version four
// wire transaction format.
func SerializeVersion4Transaction(tx *wire.MsgTx, expiryHeight uint32) ([]byte, error) {
	var buf bytes.Buffer

	// Write header
	_, err := buf.Write(txHeaderBytes)
	if err != nil {
		return nil, err
	}

	// Write group ID
	_, err = buf.Write(txNVersionGroupIDBytes)
	if err != nil {
		return nil, err
	}

	// Write varint input count
	count := uint64(len(tx.TxIn))
	err = wire.WriteVarInt(&buf, wire.ProtocolVersion, count)
	if err != nil {
		return nil, err
	}

	// Write inputs
	for _, ti := range tx.TxIn {
		// Write outpoint hash
		_, err := buf.Write(ti.PreviousOutPoint.Hash[:])
		if err != nil {
			return nil, err
		}
		// Write outpoint index
		index := make([]byte, 4)
		binary.LittleEndian.PutUint32(index, ti.PreviousOutPoint.Index)
		_, err = buf.Write(index)
		if err != nil {
			return nil, err
		}
		// Write sigscript
		err = wire.WriteVarBytes(&buf, wire.ProtocolVersion, ti.SignatureScript)
		if err != nil {
			return nil, err
		}
		// Write sequence
		sequence := make([]byte, 4)
		binary.LittleEndian.PutUint32(sequence, ti.Sequence)
		_, err = buf.Write(sequence)
		if err != nil {
			return nil, err
		}
	}
	// Write varint output count
	count = uint64(len(tx.TxOut))
	err = wire.WriteVarInt(&buf, wire.ProtocolVersion, count)
	if err != nil {
		return nil, err
	}
	// Write outputs
	for _, to := range tx.TxOut {
		// Write value
		val := make([]byte, 8)
		binary.LittleEndian.PutUint64(val, uint64(to.Value))
		_, err = buf.Write(val)
		if err != nil {
			return nil, err
		}
		// Write pkScript
		err = wire.WriteVarBytes(&buf, wire.ProtocolVersion, to.PkScript)
		if err != nil {
			return nil, err
		}
	}
	// Write nLocktime
	nLockTime := make([]byte, 4)
	binary.LittleEndian.PutUint32(nLockTime, tx.LockTime)
	_, err = buf.Write(nLockTime)
	if err != nil {
		return nil, err
	}

	// Write nExpiryHeight
	expiry := make([]byte, 4)
	binary.LittleEndian.PutUint32(expiry, expiryHeight)
	_, err = buf.Write(expiry)
	if err != nil {
		return nil, err
	}

	// Write nil value balance
	_, err = buf.Write(make([]byte, 8))
	if err != nil {
		return nil, err
	}

	// Write nil value vShieldedSpend
	_, err = buf.Write(make([]byte, 1))
	if err != nil {
		return nil, err
	}

	// Write nil value vShieldedOutput
	_, err = buf.Write(make([]byte, 1))
	if err != nil {
		return nil, err
	}

	// Write nil value vJoinSplit
	_, err = buf.Write(make([]byte, 1))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func calcHashPrevOuts(tx *wire.MsgTx) []byte {
	var b bytes.Buffer
	for _, in := range tx.TxIn {
		// First write out the 32-byte transaction ID one of whose
		// outputs are being referenced by this input.
		b.Write(in.PreviousOutPoint.Hash[:])

		// Next, we'll encode the index of the referenced output as a
		// little endian integer.
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], in.PreviousOutPoint.Index)
		b.Write(buf[:])
	}
	bl, _ := blake2b.New(&blake2b.Config{
		Size:   32,
		Person: hashPrevOutPersonalization,
	})
	bl.Write(b.Bytes())
	h := bl.Sum(nil)
	return h[:]
}

func calcHashSequence(tx *wire.MsgTx) []byte {
	var b bytes.Buffer
	for _, in := range tx.TxIn {
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], in.Sequence)
		b.Write(buf[:])
	}
	bl, _ := blake2b.New(&blake2b.Config{
		Size:   32,
		Person: hashSequencePersonalization,
	})
	bl.Write(b.Bytes())
	h := bl.Sum(nil)
	return h[:]
}

func calcHashOutputs(tx *wire.MsgTx) []byte {
	var b bytes.Buffer
	for _, out := range tx.TxOut {
		wire.WriteTxOut(&b, 0, 0, out)
	}
	bl, _ := blake2b.New(&blake2b.Config{
		Size:   32,
		Person: hashOutputsPersonalization,
	})
	bl.Write(b.Bytes())
	h := bl.Sum(nil)
	return h[:]
}
