package types

func Transfer(ref, wif, chainId, fromAddr, toAddr, coinType string, transferAmount int64, fee int64, memo string, guaranteeId string) (buf []byte, err error) {
	gid := guaranteeId
	assetId := GetAssetId(coinType)

	return BuildTransferTransaction(ref, wif, fromAddr, toAddr, memo, assetId, transferAmount, fee, coinType, gid, chainId)
}
