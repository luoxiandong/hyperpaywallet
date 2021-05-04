package util

import (
	secp "github.com/bitnexty/secp256k1-go"
)

func SignCompact(msg []byte, seckey []byte, requireCanonical bool) ([]byte, error) {
	// BtsSign has check IsCanonical
	return secp.BtsSign(msg, seckey, requireCanonical)
}
