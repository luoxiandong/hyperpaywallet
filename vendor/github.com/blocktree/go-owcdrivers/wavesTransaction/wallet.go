package wavesTransaction


func GenPair(seed []byte) (SecretKey, PublicKey, error) {
	prefix := make([]byte, 4)
	s := append(prefix, seed[:]...)

	d, err := SecureHash(s)
	if err != nil {
		return SecretKey{}, PublicKey{}, err
	}

	priv, pub, err := GenerateKeyPair(d.Bytes())
	if err != nil {
		return SecretKey{}, PublicKey{}, err
	}

	return priv, pub, nil
}
