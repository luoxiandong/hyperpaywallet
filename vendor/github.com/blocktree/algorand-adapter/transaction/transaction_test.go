package transaction

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/blocktree/algorand-adapter/address_decode"
	"github.com/blocktree/go-owcdrivers/ed25519WalletKey"
	//transaction "github.com/blocktree/algorand-adapter/transaction"
)

// TEST(AlgorandSigner, Sign) {
//     auto key = PrivateKey(parse_hex("c9d3cc16fecabe2747eab86b81528c6ed8b65efc1d6906d86aabc27187a1fe7c"));
//     auto publicKey = key.getPublicKey(TWPublicKeyTypeED25519);
//     auto from = Address(publicKey);
//     auto to = Address("UCE2U2JC4O4ZR6W763GUQCG57HQCDZEUJY4J5I6VYY4HQZUJDF7AKZO5GM");
//     Data note;
//     std::string genesisId = "mainnet-v1.0";
//     auto genesisHash = Base64::decode("wGHE2Pwdvd7S12BL5FaOP20EGYesN73ktiC1qzkkit8=");
//     auto transaction = Transaction(
//         /* from */ from,
//         /* to */ to,
//         /* fee */ 488931,
//         /* amount */ 847,
//         /* first round */ 51,
//         /* last round */ 61,
//         /* note */ note,
//         /* type */ "pay",
//         /* genesis id*/ genesisId,
//         /* genesis hash*/ genesisHash
//     );

//     auto serialized = transaction.serialize();
//     auto signature = Signer::sign(key, transaction);
//     auto result = transaction.serialize(signature);

//     ASSERT_EQ(hex(serialized), "89a3616d74cd034fa3666565ce000775e3a2667633a367656eac6d61696e6e65742d76312e30a26768c420c061c4d8fc1dbdded2d7604be4568e3f6d041987ac37bde4b620b5ab39248adfa26c763da3726376c420a089aa6922e3b998fadff6cd4808ddf9e021e4944e389ea3d5c638786689197ea3736e64c42074b000b6368551a6066d713e2866002e8dab34b69ede09a72e85a39bbb1f7928a474797065a3706179");
//     ASSERT_EQ(hex(signature), "de73363dbdeda0682adca06f6268a16a6ec47253c94d5692dc1c49a84a05847812cf66d7c4cf07c7e2f50f143ec365d405e30b35117b264a994626054d2af604");
//     ASSERT_EQ(hex(result), "82a3736967c440de73363dbdeda0682adca06f6268a16a6ec47253c94d5692dc1c49a84a05847812cf66d7c4cf07c7e2f50f143ec365d405e30b35117b264a994626054d2af604a374786e89a3616d74cd034fa3666565ce000775e3a2667633a367656eac6d61696e6e65742d76312e30a26768c420c061c4d8fc1dbdded2d7604be4568e3f6d041987ac37bde4b620b5ab39248adfa26c763da3726376c420a089aa6922e3b998fadff6cd4808ddf9e021e4944e389ea3d5c638786689197ea3736e64c42074b000b6368551a6066d713e2866002e8dab34b69ede09a72e85a39bbb1f7928a474797065a3706179");
// }

func Test_SerializedRawTx(t *testing.T) {
	parse_hex := "c9d3cc16fecabe2747eab86b81528c6ed8b65efc1d6906d86aabc27187a1fe7c"
	sk, _ := hex.DecodeString(parse_hex)
	pkByte := ed25519WalletKey.WalletPubKeyFromKeyBytes(sk)
	algoDecoder := &address_decode.AddressDecoderV2{}
	from, _ := algoDecoder.AddressEncode(pkByte)
	to := "UCE2U2JC4O4ZR6W763GUQCG57HQCDZEUJY4J5I6VYY4HQZUJDF7AKZO5GM"
	genesisId := "mainnet-v1.0"
	genesisHash, _ := base64.StdEncoding.DecodeString("wGHE2Pwdvd7S12BL5FaOP20EGYesN73ktiC1qzkkit8=")

	//         /* from */ from,
	//         /* to */ to,
	//         /* fee */ 488931,
	//         /* amount */ 847,
	//         /* first round */ 51,
	//         /* last round */ 61,
	//         /* note */ note,
	//         /* type */ "pay",
	//         /* genesis id*/ genesisId,
	//         /* genesis hash*/ genesisHash
	//     );
	fromAddress, err := address_decode.DecodeAddress(from)
	if err != nil {
		fmt.Println(err.Error())
	}
	toAddress, err := address_decode.DecodeAddress(to)
	if err != nil {
		fmt.Println(err.Error())
	}
	trans := &ALGOTransaction{
		From:        fromAddress[:],
		To:          toAddress[:],
		Fee:         488931,
		Amount:      847,
		FirstRound:  51,
		LastRound:   61,
		Type:        "pay",
		GenesisId:   genesisId,
		GenesisHash: genesisHash,
	}

	fmt.Println(trans)
	signature := SignData(trans.SerializeRawTx(), sk)
	sigHex := hex.EncodeToString(signature)
	if sigHex == "de73363dbdeda0682adca06f6268a16a6ec47253c94d5692dc1c49a84a05847812cf66d7c4cf07c7e2f50f143ec365d405e30b35117b264a994626054d2af604" {
		fmt.Println("sigHex 相等")
	} else {
		fmt.Println("sigHex 不相等")
	}
	fmt.Println(sigHex)
	fmt.Println("de73363dbdeda0682adca06f6268a16a6ec47253c94d5692dc1c49a84a05847812cf66d7c4cf07c7e2f50f143ec365d405e30b35117b264a994626054d2af604")

	serializeHex := hex.EncodeToString(trans.SerializeRawTx())
	if serializeHex == "89a3616d74cd034fa3666565ce000775e3a2667633a367656eac6d61696e6e65742d76312e30a26768c420c061c4d8fc1dbdded2d7604be4568e3f6d041987ac37bde4b620b5ab39248adfa26c763da3726376c420a089aa6922e3b998fadff6cd4808ddf9e021e4944e389ea3d5c638786689197ea3736e64c42074b000b6368551a6066d713e2866002e8dab34b69ede09a72e85a39bbb1f7928a474797065a3706179" {
		fmt.Println("serializeHex 相等")

	} else {
		fmt.Println("serializeHex 不相等")

	}
	fmt.Println(serializeHex)
	fmt.Println("89a3616d74cd034fa3666565ce000775e3a2667633a367656eac6d61696e6e65742d76312e30a26768c420c061c4d8fc1dbdded2d7604be4568e3f6d041987ac37bde4b620b5ab39248adfa26c763da3726376c420a089aa6922e3b998fadff6cd4808ddf9e021e4944e389ea3d5c638786689197ea3736e64c42074b000b6368551a6066d713e2866002e8dab34b69ede09a72e85a39bbb1f7928a474797065a3706179")

	serializeSignHex := hex.EncodeToString(trans.SerializeRawSignTx(signature))
	if serializeSignHex == "82a3736967c440de73363dbdeda0682adca06f6268a16a6ec47253c94d5692dc1c49a84a05847812cf66d7c4cf07c7e2f50f143ec365d405e30b35117b264a994626054d2af604a374786e89a3616d74cd034fa3666565ce000775e3a2667633a367656eac6d61696e6e65742d76312e30a26768c420c061c4d8fc1dbdded2d7604be4568e3f6d041987ac37bde4b620b5ab39248adfa26c763da3726376c420a089aa6922e3b998fadff6cd4808ddf9e021e4944e389ea3d5c638786689197ea3736e64c42074b000b6368551a6066d713e2866002e8dab34b69ede09a72e85a39bbb1f7928a474797065a3706179" {
		fmt.Println("serializeSignHex 相等")

	} else {
		fmt.Println("serializeSignHex 不相等")

	}
	fmt.Println(serializeSignHex)
	fmt.Println("82a3736967c440de73363dbdeda0682adca06f6268a16a6ec47253c94d5692dc1c49a84a05847812cf66d7c4cf07c7e2f50f143ec365d405e30b35117b264a994626054d2af604a374786e89a3616d74cd034fa3666565ce000775e3a2667633a367656eac6d61696e6e65742d76312e30a26768c420c061c4d8fc1dbdded2d7604be4568e3f6d041987ac37bde4b620b5ab39248adfa26c763da3726376c420a089aa6922e3b998fadff6cd4808ddf9e021e4944e389ea3d5c638786689197ea3736e64c42074b000b6368551a6066d713e2866002e8dab34b69ede09a72e85a39bbb1f7928a474797065a3706179")
}
