package account

import (
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcdrivers/bytomTransaction/chainkd"
	"golang.org/x/crypto/ripemd160"
)

func CreateP2PKH(xPubs []chainkd.XPub, path [][]byte) (string, error) {
	derivedXPubs := chainkd.DeriveXPubs(xPubs, path)
	derivedPK := derivedXPubs[0].PublicKey()
	pubHash := Ripemd160(derivedPK)

	address := addressEncoder.AddressEncode(pubHash, addressEncoder.BTM_mainnetAddressBech32V0)

	return address, nil
}
func GenerateDerivedXprv(xPrv chainkd.XPrv, path [][]byte) chainkd.XPrv {
	derivedXPrv := xPrv.Derive(path)

	return derivedXPrv
}

func GenerateDerivedXpubAddress(xPrv chainkd.XPrv) (chainkd.XPub, string) {
	derivedXPub := xPrv.XPub()
	pubHash := Ripemd160(derivedXPub.PublicKey())

	address := addressEncoder.AddressEncode(pubHash, addressEncoder.BTM_mainnetAddressBech32V0)
	return derivedXPub, address
}

func GenerateP2PKH(xPub chainkd.XPub) (string, error) {
	derivedPK := xPub.PublicKey()
	pubHash := Ripemd160(derivedPK)

	address := addressEncoder.AddressEncode(pubHash, addressEncoder.BTM_mainnetAddressBech32V0)

	return address, nil
}

func Ripemd160(data []byte) []byte {
	ripemd := ripemd160.New()
	ripemd.Write(data)

	return ripemd.Sum(nil)
}
