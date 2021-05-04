package types

import (
	"fmt"
)

var (
	ErrRPCClientNotInitialized      = fmt.Errorf("RPC client is not initialized")
	ErrNotImplemented               = fmt.Errorf("not implemented")
	ErrInvalidInputType             = fmt.Errorf("invalid input type")
	ErrInvalidInputLength           = fmt.Errorf("invalid input length")
	ErrInvalidPublicKey             = fmt.Errorf("invalid PublicKey")
	ErrInvalidAddress               = fmt.Errorf("invalid Address")
	ErrPublicKeyChainPrefixMismatch = fmt.Errorf("PublicKey database prefix mismatch")
	ErrAddressChainPrefixMismatch   = fmt.Errorf("Address database prefix mismatch")
	ErrInvalidChecksum              = fmt.Errorf("invalid checksum")
	ErrNoSigningKeyFound            = fmt.Errorf("no signing key found")
	ErrNoVerifyingKeyFound          = fmt.Errorf("no verifying key found")
	ErrInvalidDigestLength          = fmt.Errorf("invalid digest length")
	ErrInvalidPrivateKeyCurve       = fmt.Errorf("invalid PrivateKey curve")
	ErrCurrentChainConfigIsNotSet   = fmt.Errorf("current database config is not set")
)

type Int8 int8
type UInt8 uint8
type UInt16 uint16
type UInt32 uint32
type UInt64 uint64
type Int64 int64

const addrPrefix = "HX"          // 地址前缀
const walletVersion = byte(0x35) // 版本
