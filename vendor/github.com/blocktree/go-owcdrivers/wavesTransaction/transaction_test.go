package wavesTransaction

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_transaction(t *testing.T) {
	tests := []struct {
		id          string
		sig         string
		spk         string
		rcp         string
		amountAsset string
		amount      uint64
		feeAsset    string
		fee         uint64
		timestamp   uint64
		attachment  string
	}{
		{"93H1i2jgP21Eh4Q5uzwmCYCVfGHZcAMzpC6PPbwvCSTs", "4jbfvXGiqsaKkZso6ykNMQZDARewvBjKxgaz55jF4g6tBVPgxv5qChSvBYdHRGHjUdXbG3CZ3PUNBBK3eoiuRfVt", "6tbTkJukCZ4qX13ucXVeUV2aN88t9ypi1MADZb9PfFQD", "3P5mTiUpUnb1eM19udtM8QyNLBGf6VjS19j", "GryqKQBmTZGZnbZ4efrQvNGNpeLM83djWSNJBWuhZg5H", 566, "9PVyxDPUjauYafvq83JTXvHQ8nPnxwKA7siUFcqthCDJ", 1000000000, 1541593367281, ""},
		{"HE7jA4xjRiqdNVEP4jSXAY8FEy412MGKTD1hqWtpnYrZ", "33VU7yYd6bLrHf5VBCtR5iMpxDhMwWWdfJMSzBUw38WsSjCsKhfjerBiatDtJNpPPx8FK8cqX4kKeb5XiUhSv7ev", "6cDEgFTH9mZwuJjRubE4ZzSjrCvofzn24M7jmw5oTu5p", "3PKi8kvBCMUZnPFVRDBMzYaY49wLf7TurEe", "FGJQGTG13wKXSaYB4JJ6But7Ui3iRq5ZA9DTFsNTYJvt", 9788200000000, "", 100000, 1541593585115, "0x09F7f8d4f0e4BCC89073318759179EB1e5cFC500"},
		{"ERhAQmKArX6Yy2iC5N9S9aV9xPhnabyhwMBXufSZkgEw", "36u5TudkkRDE6V67jydUCAE3Vh8xf7Yv5FM8M53DJHtyhjdEPJtFqLBAnEFaXUEyyV2s7qPV8DyL3nyhDJ7YBCcH", "9zMXKmq3tWJJmezrkaYjmpiD3LZkFbhiwy9AP2r4CBnC", "3PGxhF7LtybhRdZfErxBTN4ZDDJjLUQW8Rb", "", 6500000, "", 100000, 1541593775634, "Send"},
	}
	for _, tc := range tests {
		id, _ := NewDigestFromBase58(tc.id)
		//sig, _ := NewSignatureFromBase58(tc.sig)
		spk, _ := NewPublicKeyFromBase58(tc.spk)
		addr, _ := NewAddressFromString(tc.rcp)
		rcp := NewRecipientFromAddress(addr)
		aa, _ := NewOptionalAssetFromString(tc.amountAsset)
		fa, _ := NewOptionalAssetFromString(tc.feeAsset)
		fmt.Println("spk = ", tc.spk)
		fmt.Println("timestamp = ", tc.timestamp)
		fmt.Println("amount = ", tc.amount)
		fmt.Println("fee = ", tc.fee)

		fmt.Println("amountAsset = ", tc.amountAsset)
		fmt.Println("feeAsset = ", tc.feeAsset)
		fmt.Println("rcp = ", rcp)
		fmt.Println("attachment = ", tc.attachment)

		tx := NewUnsignedTransferV2(spk, *aa, *fa, tc.timestamp, tc.amount, tc.fee, rcp, tc.attachment)
		b, _ := tx.BodyMarshalBinary()
		if hex.EncodeToString(b) == "0402534f9917bf4361da0060f4514ca8b98dc16a15e699983dcbf0e542b15888615701d3ef0130448df654b209722a00bc84d1492211dfd5c1a8baaa8b08f09ab815b10000000166ee2355db000008e6fe2f1a0000000000000186a00157c3089138a7430a89996c29b3a34ba10172ec718a6fb6bb7b002a307830394637663864346630653442434338393037333331383735393137394542316535634643353030" {
			fmt.Println("相等")
		}
		h, _ := FastHash(b)
		if id == h {

		}
	}
}
