package keypair

import (
	"github.com/liyaojian/hxwallet/types"
	"github.com/tyler-smith/go-bip39"
)

type KeyPair struct {
	BrainKey   string
	PrivateKey *types.PrivateKey
}

//Generates the key pair
func GenerateKeyPair(brainKey string) (*KeyPair, error) {
	if len(brainKey) == 0 {
		// Generate a mnemonic for memorization or user-friendly seeds
		entropy, err := bip39.NewEntropy(128)
		if err != nil {
			return nil, err
		}
		if brainKey, err = bip39.NewMnemonic(entropy); err != nil {
			return nil, err
		}
	}

	pri, err := types.NewPrivateKeyFromBrainKey(brainKey)
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		BrainKey:   brainKey,
		PrivateKey: pri,
	}, nil
}

//Export public key from private key
func PrivateToPublic(priWif string) (*types.PublicKey, error) {
	pri, err := types.NewPrivateKeyFromWif(priWif)
	if err != nil {
		return nil, err
	}
	return pri.PublicKey(), nil
}

//Check if privateKey is valid or not
func IsValidPrivate(priWif string) bool {
	_, err := types.NewPrivateKeyFromWif(priWif)
	if err != nil {
		return false
	}
	return true
}

//Check if publicKey is valid or not
func IsValidPublic(pub string) bool {
	_, err := types.NewPublicKeyFromString(pub)
	if err != nil {
		return false
	}
	return true
}
