package types

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_TransactionSerializeriOSComHash(t *testing.T) {
	transaction := &Transaction{Version: "0x0"}

	headerDeps := []Hash{}
	transaction.HeaderDeps = headerDeps

	cellDep := []CellDep{}
	outpoint := &OutPoint{
		TxHash: "0xbffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a",
		Index:  "0x00",
	}

	cellDep1 := &CellDep{
		OutPoint: *outpoint,
		DepType:  Code,
	}

	cellDep = append(cellDep, *cellDep1)
	transaction.CellDeps = cellDep
	outputs := []CellOutput{}
	lock := &Script{
		Args:     "0xe2193df51d78411601796b35b17b4f8f2cd85bd0",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	output := &CellOutput{
		Capacity: "0x174876e800",
		Lock:     *lock,
		Type:     nil,
	}

	lock1 := &Script{
		Args:     "0x36c329ed630d6ce750712a477543672adab57f4c",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	output1 := &CellOutput{
		Capacity: "0x474dec26800",
		Lock:     *lock1,
		Type:     nil,
	}

	outputs = append(outputs, *output)
	outputs = append(outputs, *output1)

	transaction.Outputs = outputs
	outputData := []Bytes{}
	outputData = append(outputData, "0x")
	outputData = append(outputData, "0x")

	transaction.OutputsData = outputData
	witnesses := []Bytes{}
	witnesses = append(witnesses, "0x")

	transaction.Witnesses = witnesses

	fmt.Println(transaction)
	serializeTx, err := transaction.Serialize()
	if err != nil {
		fmt.Println(err)
	}
	comHash, err := ComputeHash(serializeTx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ComputeHash = ", hex.EncodeToString(comHash))

	fmt.Println("ComputeHash = ", "4c02905db773301f73bbc6cd5a400c928caf410bbb13136f6f48bec0a79c22e4")
}
func Test_TransactionJSon(t *testing.T) {

	expectHex := "5f0100001c00000020000000490000004d0000007d0000004b0100000000000001000000b815a396c5226009670e89ee514850dcde452bca746cdd6b41c104b50e559c70000000000100000000010000000000000000000000ee046ce2baeda575266d4164f394c53f66009f64759f7a9f12a014c692e7939006000000ce0000000c0000006d0000006100000010000000180000006100000000406352bfc60100490000001000000030000000310000009bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce80114000000470dcdc5e44064909650113a274b3b36aecb6dc76100000010000000180000006100000000406352bfc60100490000001000000030000000310000009bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce80114000000c8328aabcd9b9e8e64fbc566c4385c3bdeb219d7140000000c000000100000000000000000000000"

	transaction := &Transaction{Version: "0x0"}
	headerDeps := []Hash{}
	transaction.HeaderDeps = headerDeps
	cellDep := []CellDep{}
	outpoint := &OutPoint{
		TxHash: "0xb815a396c5226009670e89ee514850dcde452bca746cdd6b41c104b50e559c70",
		Index:  "0x0",
	}

	cellDep1 := &CellDep{
		OutPoint: *outpoint,
		DepType:  "dep_group",
	}

	cellDep = append(cellDep, *cellDep1)
	transaction.CellDeps = cellDep

	inputs := []CellInput{}

	inpoint := &OutPoint{
		TxHash: "0xee046ce2baeda575266d4164f394c53f66009f64759f7a9f12a014c692e79390",
		Index:  "0x6",
	}

	input := &CellInput{
		PreviousOutput: *inpoint,
		Since:          "0x0",
	}

	inputs = append(inputs, *input)
	transaction.Inputs = inputs

	outputs := []CellOutput{}
	lock := &Script{
		Args:     "0x470dcdc5e44064909650113a274b3b36aecb6dc7",
		CodeHash: "0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8",
		HashType: Type,
	}

	output := &CellOutput{
		Capacity: "0x1c6bf52634000",
		Lock:     *lock,
		Type:     nil,
	}

	lock1 := &Script{
		Args:     "0xc8328aabcd9b9e8e64fbc566c4385c3bdeb219d7",
		CodeHash: "0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8",
		HashType: Type,
	}

	output1 := &CellOutput{
		Capacity: "0x1c6bf52634000",
		Lock:     *lock1,
		Type:     nil,
	}

	outputs = append(outputs, *output)
	outputs = append(outputs, *output1)

	transaction.Outputs = outputs
	outputData := []Bytes{}
	outputData = append(outputData, "0x")
	outputData = append(outputData, "0x")

	transaction.OutputsData = outputData
	witnesses := []Bytes{}
	// witnesses = append(witnesses, "0x")

	transaction.Witnesses = witnesses

	fmt.Println(transaction)
	serializeTx, err := transaction.Serialize()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(serializeTx)

	// err := json.Unmarshal([]byte(transaction), &tx)
	// if err != nil {
	// 	t.Errorf("fail to unmarshal test transaction json: %s\n", err)
	// 	return
	// }

	got, err := transaction.Serialize()
	if err != nil {
		t.Errorf("fail to serialize: %s\n", err)
		return
	}

	gotHex := hex.EncodeToString(got)
	if gotHex == expectHex {
		fmt.Println("got hex 相等")
	} else {
		fmt.Println("got hex 不相等")
	}
	fmt.Println(gotHex)
	fmt.Println(expectHex)
}

func Test_TransactionSerializer(t *testing.T) {
	transaction := &Transaction{Version: "0x0"}
	cellDep := []CellDep{}
	outpoint := &OutPoint{
		TxHash: "0xbffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a",
		Index:  "0x0",
	}

	headerDeps := []Hash{}
	transaction.HeaderDeps = headerDeps
	cellDep1 := &CellDep{
		OutPoint: *outpoint,
		DepType:  "code",
	}

	cellDep = append(cellDep, *cellDep1)
	transaction.CellDeps = cellDep
	outputs := []CellOutput{}
	lock := &Script{
		Args:     "0xe2193df51d78411601796b35b17b4f8f2cd85bd0",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	output := &CellOutput{
		Capacity: "0x174876e800",
		Lock:     *lock,
		Type:     nil,
	}

	lock1 := &Script{
		Args:     "0x36c329ed630d6ce750712a477543672adab57f4c",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	output1 := &CellOutput{
		Capacity: "0x474dec26800",
		Lock:     *lock1,
		Type:     nil,
	}

	outputs = append(outputs, *output)
	outputs = append(outputs, *output1)

	transaction.Outputs = outputs
	outputData := []Bytes{}
	outputData = append(outputData, "0x")
	outputData = append(outputData, "0x")

	transaction.OutputsData = outputData
	witnesses := []Bytes{}
	witnesses = append(witnesses, "0x")

	transaction.Witnesses = witnesses

	fmt.Println(transaction)
	serializeTx, err := transaction.Serialize()
	if err != nil {
		fmt.Println(err)
	}
	comHash, err := ComputeHash(serializeTx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ComputeHash = ", hex.EncodeToString(comHash))

	fmt.Println("ComputeHash = ", "4c02905db773301f73bbc6cd5a400c928caf410bbb13136f6f48bec0a79c22e4")
}

func Test_TransactionSignTestsiOS(t *testing.T) {
	transaction := &Transaction{Version: "0x0"}
	cellDep := []CellDep{}
	outpoint := &OutPoint{
		TxHash: "0xa76801d09a0eabbfa545f1577084b6f3bafb0b6250e7f5c89efcfd4e3499fb55",
		Index:  "0x1",
	}
	cellDep1 := &CellDep{
		OutPoint: *outpoint,
		DepType:  Code,
	}
	cellDep = append(cellDep, *cellDep1)

	transaction.CellDeps = cellDep

	// headerDeps := []Hash{}
	// transaction.HeaderDeps = headerDeps

	inputs := []CellInput{}

	inpoint := &OutPoint{
		TxHash: "0xa80a8e01d45b10e1cbc8a2557c62ba40edbdc36cd63a31fc717006ca7b157b50",
		Index:  "0x0",
	}

	input := &CellInput{
		PreviousOutput: *inpoint,
		Since:          "0x0",
	}

	inputs = append(inputs, *input)
	transaction.Inputs = inputs

	outputs := []CellOutput{}
	lock := &Script{
		Args:     "0xe2193df51d78411601796b35b17b4f8f2cd85bd0",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}
	// typeScript := &Script{
	// 	Args:     "0xe2193df51d78411601796b35b17b4f8f2cd85bd0",
	// 	CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
	// 	HashType: Data,
	// }
	// fmt.Println(typeScript)
	output := &CellOutput{
		Capacity: "0x174876e800",
		Lock:     *lock,
		Type:     nil,
	}
	outputs = append(outputs, *output)

	lock1 := &Script{
		Args:     "0x36c329ed630d6ce750712a477543672adab57f4c",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	output1 := &CellOutput{
		Capacity: "0x474dec26800",
		Lock:     *lock1,
		Type:     nil,
	}

	outputs = append(outputs, *output1)

	transaction.Outputs = outputs
	outputData := []Bytes{}
	outputData = append(outputData, "0x")
	outputData = append(outputData, "0x")

	transaction.OutputsData = outputData

	witnesses := []Bytes{}
	// witnesses = append(witnesses, "0x82df73581bcd08cb9aa270128d15e79996229ce8ea9e4f985b49fbf36762c5c37936caf3ea3784ee326f60b8992924fcf496f9503c907982525a3436f01ab32900")

	transaction.Witnesses = witnesses
	fmt.Println(transaction)
	serializeTx, err := transaction.Serialize()
	if err != nil {
		fmt.Println(err)
	}
	comHash, err := ComputeHash(serializeTx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ComputeHash = ", hex.EncodeToString(comHash))

	fmt.Println("ComputeHash = ", "d928abe5d7baab12837dbae5273175d260b8a0ebb7f24087d588c08779e25d31")

	fmt.Println("serializeTx")

	fmt.Println(hex.EncodeToString(serializeTx))
	fmt.Println("serializeTx")
	privateKey, _ := hex.DecodeString("e79f3207ea4980b7fed79956d5934249ceac4751a4fae01a0f7c4a96884bc4e3")
	witness, err := ComputeWitnessHash(serializeTx, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	transaction.Witnesses = witness
	SerializeRawTx, _ := transaction.SerializeRawTx()
	fmt.Println("transaction Hash = ", hex.EncodeToString(SerializeRawTx))
	fmt.Println("Witness Hash = ", witness[0])

	fmt.Println("Witness Hash = ", "0x55000000100000005500000055000000410000007a360306c20f1f0081d27feff5c59fb9b4307b25876543848010614fb78ea21d165f48f67ae3357eeafbad2033b1e53cd737d4e670de60e1081d514b1e05cf5100")

}

func Test_TransactionSerializerIOS(t *testing.T) {
	transaction := &Transaction{Version: "0x0"}
	cellDep := []CellDep{}
	outpoint := &OutPoint{
		TxHash: "0xbffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a",
		Index:  "0x0",
	}

	headerDeps := []Hash{}
	transaction.HeaderDeps = headerDeps
	cellDep1 := &CellDep{
		OutPoint: *outpoint,
		DepType:  Code,
	}

	cellDep = append(cellDep, *cellDep1)
	transaction.CellDeps = cellDep
	outputs := []CellOutput{}
	lock := &Script{
		Args:     "0xe2193df51d78411601796b35b17b4f8f2cd85bd0",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	i := 500000000000000
	h := fmt.Sprintf("%x", i)
	fmt.Printf("Hex conv of '%d' is '%s'\n", i, h)

	output := &CellOutput{
		Capacity: "0x174876e800",
		Lock:     *lock,
		Type:     nil,
	}

	lock1 := &Script{
		Args:     "0x36c329ed630d6ce750712a477543672adab57f4c",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	output1 := &CellOutput{
		Capacity: "0x474dec26800",
		Lock:     *lock1,
		Type:     nil,
	}

	outputs = append(outputs, *output)
	outputs = append(outputs, *output1)

	transaction.Outputs = outputs
	outputData := []Bytes{}
	outputData = append(outputData, "0x")
	outputData = append(outputData, "0x")

	transaction.OutputsData = outputData
	witnesses := []Bytes{}
	witnesses = append(witnesses, "0x")
	transaction.Witnesses = witnesses
	fmt.Println(transaction)

	serializeTx, err := transaction.Serialize()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("serializeTx")

	fmt.Println(hex.EncodeToString(serializeTx))
	fmt.Println("serializeTx")
	rawSerializeTx, _ := hex.DecodeString("330100001c00000020000000490000004d000000510000001f0100000000000001000000bffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a00000000000000000000000000ce0000000c0000006d0000006100000010000000180000006100000000e8764817000000490000001000000030000000310000009e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a080014000000e2193df51d78411601796b35b17b4f8f2cd85bd0610000001000000018000000610000000068c2de74040000490000001000000030000000310000009e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08001400000036c329ed630d6ce750712a477543672adab57f4c140000000c000000100000000000000000000000")
	if hex.EncodeToString(serializeTx) == "330100001c00000020000000490000004d000000510000001f0100000000000001000000bffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a00000000000000000000000000ce0000000c0000006d0000006100000010000000180000006100000000e8764817000000490000001000000030000000310000009e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a080014000000e2193df51d78411601796b35b17b4f8f2cd85bd0610000001000000018000000610000000068c2de74040000490000001000000030000000310000009e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08001400000036c329ed630d6ce750712a477543672adab57f4c140000000c000000100000000000000000000000" {
		fmt.Println("serializeTx 相等")
	} else {
		fmt.Println("serializeTx 不相等")
	}

	ComputeHash(rawSerializeTx)

	comHash, err := ComputeHash(serializeTx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ComputeHash = ", hex.EncodeToString(comHash))

	fmt.Println("ComputeHash = ", "4c02905db773301f73bbc6cd5a400c928caf410bbb13136f6f48bec0a79c22e4")

	// Calc witness args
	// 65 is secp256k1 signature length
}

func Test_TransactionSignIOS(t *testing.T) {
	transaction := &Transaction{Version: "0x0"}
	cellDep := []CellDep{}
	outpoint := &OutPoint{
		TxHash: "0xa76801d09a0eabbfa545f1577084b6f3bafb0b6250e7f5c89efcfd4e3499fb55",
		Index:  "0x1",
	}

	headerDeps := []Hash{}
	transaction.HeaderDeps = headerDeps
	cellDep1 := &CellDep{
		OutPoint: *outpoint,
		DepType:  Code,
	}

	cellDep = append(cellDep, *cellDep1)
	transaction.CellDeps = cellDep

	inputs := []CellInput{}

	inpoint := &OutPoint{
		TxHash: "0xa80a8e01d45b10e1cbc8a2557c62ba40edbdc36cd63a31fc717006ca7b157b50",
		Index:  "0x00",
	}

	input := &CellInput{
		PreviousOutput: *inpoint,
		Since:          "0x00",
	}

	inputs = append(inputs, *input)
	transaction.Inputs = inputs

	outputs := []CellOutput{}
	lock := &Script{
		Args:     "0xe2193df51d78411601796b35b17b4f8f2cd85bd0",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	i := 100000000000
	h := fmt.Sprintf("0x%x", i)
	fmt.Printf("Hex conv of '%d' is '%s'\n", i, h)

	output := &CellOutput{
		Capacity: Uint64(h),
		Lock:     *lock,
		Type:     nil,
	}

	lock1 := &Script{
		Args:     "0x36c329ed630d6ce750712a477543672adab57f4c",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}
	i1 := 4900000000000
	h1 := fmt.Sprintf("0x%x", i1)
	fmt.Printf("Hex conv of '%d' is '%s'\n", i1, h1)
	output1 := &CellOutput{
		Capacity: Uint64(h1),
		Lock:     *lock1,
		Type:     nil,
	}

	outputs = append(outputs, *output)
	outputs = append(outputs, *output1)

	transaction.Outputs = outputs
	outputData := []Bytes{}
	outputData = append(outputData, "0x")
	outputData = append(outputData, "0x")

	transaction.OutputsData = outputData
	witnesses := []Bytes{}
	witnesses = append(witnesses, "0x")
	transaction.Witnesses = witnesses
	fmt.Println(transaction)

	serializeTx, err := transaction.Serialize()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("serializeTx")

	fmt.Println(hex.EncodeToString(serializeTx))
	fmt.Println("serializeTx")
	privateKey, _ := hex.DecodeString("e79f3207ea4980b7fed79956d5934249ceac4751a4fae01a0f7c4a96884bc4e3")
	witness, err := ComputeWitnessHash(serializeTx, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Witness Hash = ", witness[0])

	fmt.Println("Witness Hash = ", "0x55000000100000005500000055000000410000007a360306c20f1f0081d27feff5c59fb9b4307b25876543848010614fb78ea21d165f48f67ae3357eeafbad2033b1e53cd737d4e670de60e1081d514b1e05cf5100")

	// Calc witness args
	// 65 is secp256k1 signature length
}

func Test_TransactionSignJava(t *testing.T) {
	transaction := &Transaction{Version: "0x0"}
	cellDep := []CellDep{}
	outpoint := &OutPoint{
		TxHash: "0xbffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a",
		Index:  "0x1",
	}

	headerDeps := []Hash{}
	// headerDeps = append(headerDeps, "0x")
	transaction.HeaderDeps = headerDeps
	cellDep1 := &CellDep{
		OutPoint: *outpoint,
		DepType:  DepGroup,
	}

	cellDep = append(cellDep, *cellDep1)
	transaction.CellDeps = cellDep

	inputs := []CellInput{}

	inpoint := &OutPoint{
		TxHash: "0xa80a8e01d45b10e1cbc8a2557c62ba40edbdc36cd63a31fc717006ca7b157b50",
		Index:  "0x0",
	}

	input := &CellInput{
		PreviousOutput: *inpoint,
		Since:          "0x0",
	}

	inputs = append(inputs, *input)
	transaction.Inputs = inputs

	outputs := []CellOutput{}
	lock := &Script{
		Args:     "0xe2193df51d78411601796b35b17b4f8f2cd85bd0",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	i := 100000000000
	h := fmt.Sprintf("0x%x", i)
	fmt.Printf("Hex conv of '%d' is '%s'\n", i, h)

	output := &CellOutput{
		Capacity: Uint64(h),
		Lock:     *lock,
		Type:     nil,
	}

	lock1 := &Script{
		Args:     "0x36c329ed630d6ce750712a477543672adab57f4c",
		CodeHash: "0xe3b513a2105a5d4f833d1fad3d968b96b4510687234cd909f86b3ac450d8a2b5",
		HashType: Data,
	}
	i1 := 4900000000000
	h1 := fmt.Sprintf("0x%x", i1)
	fmt.Printf("Hex conv of '%d' is '%s'\n", i1, h1)
	output1 := &CellOutput{
		Capacity: Uint64(h1),
		Lock:     *lock1,
		Type:     nil,
	}

	outputs = append(outputs, *output)
	outputs = append(outputs, *output1)

	transaction.Outputs = outputs
	outputData := []Bytes{}
	outputData = append(outputData, "0x")
	outputData = append(outputData, "0x")

	transaction.OutputsData = outputData
	witnesses := []Bytes{}
	// witnesses = append(witnesses, "0x")
	transaction.Witnesses = witnesses
	//fmt.Println(transaction)

	serializeTx, err := transaction.Serialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("serializeTx")

	fmt.Println(hex.EncodeToString(serializeTx))
	fmt.Println("serializeRawTx")

	fmt.Println("7f0100001c00000020000000490000006d0000009d0000006b0100000000000001000000bffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a0100000001010000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000a80a8e01d45b10e1cbc8a2557c62ba40edbdc36cd63a31fc717006ca7b157b5000000000ce0000000c0000006d0000006100000010000000180000006100000000e8764817000000490000001000000030000000310000009e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a080014000000e2193df51d78411601796b35b17b4f8f2cd85bd0610000001000000018000000610000000068c2de7404000049000000100000003000000031000000e3b513a2105a5d4f833d1fad3d968b96b4510687234cd909f86b3ac450d8a2b5001400000036c329ed630d6ce750712a477543672adab57f4c140000000c000000100000000000000000000000")

	txHash, _ := ComputeHash(serializeTx)
	fmt.Println("txHash = ", hex.EncodeToString(txHash))
	privateKey, _ := hex.DecodeString("e79f3207ea4980b7fed79956d5934249ceac4751a4fae01a0f7c4a96884bc4e3")
	witness, err := ComputeWitnessHash(serializeTx, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	transaction.Witnesses = witness
	serializeRawTx, err := transaction.SerializeRawTx()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("serializeTx")
	fmt.Println("serializeTransaction  = ", hex.EncodeToString(serializeRawTx))

	fmt.Println("serializeTransaction  = ", "ec0100000c0000008b0100007f0100001c00000020000000490000006d0000009d0000006b0100000000000001000000bffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a0100000001010000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000a80a8e01d45b10e1cbc8a2557c62ba40edbdc36cd63a31fc717006ca7b157b5000000000ce0000000c0000006d0000006100000010000000180000006100000000e8764817000000490000001000000030000000310000009e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a080014000000e2193df51d78411601796b35b17b4f8f2cd85bd0610000001000000018000000610000000068c2de7404000049000000100000003000000031000000e3b513a2105a5d4f833d1fad3d968b96b4510687234cd909f86b3ac450d8a2b5001400000036c329ed630d6ce750712a477543672adab57f4c140000000c000000100000000000000000000000610000000800000055000000550000001000000055000000550000004100000010f86974898b2f3685facb78741801bf2b932c7c548afe5bbc5d06ee135aeb792d700a02b62c492f1fd6e88afd655ffe305489fe9a76670a8999c641c8e2b16701")

	// Calc witness args
	// 65 is secp256k1 signature length
}

func Test_TransactionComHashIOS(t *testing.T) {
	transaction := &Transaction{Version: "0x0"}
	cellDep := []CellDep{}
	outpoint := &OutPoint{
		TxHash: "0xbffab7ee0a050e2cb882de066d3dbf3afdd8932d6a26eda44f06e4b23f0f4b5a",
		Index:  "0x0",
	}

	// headerDeps := []Hash{}
	// //headerDeps = append(headerDeps, "0x")
	// transaction.HeaderDeps = headerDeps
	cellDep1 := &CellDep{
		OutPoint: *outpoint,
		DepType:  Code,
	}

	cellDep = append(cellDep, *cellDep1)
	transaction.CellDeps = cellDep

	inputs := []CellInput{}

	// inpoint := &OutPoint{
	// 	TxHash: "0xa80a8e01d45b10e1cbc8a2557c62ba40edbdc36cd63a31fc717006ca7b157b50",
	// 	Index:  "0x0",
	// }

	// input := &CellInput{
	// 	PreviousOutput: *inpoint,
	// 	Since:          "0x0",
	// }

	// inputs = append(inputs, *input)
	transaction.Inputs = inputs

	outputs := []CellOutput{}
	lock := &Script{
		Args:     "0xe2193df51d78411601796b35b17b4f8f2cd85bd0",
		CodeHash: "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08",
		HashType: Data,
	}

	i := 100000000000
	h := fmt.Sprintf("0x%x", i)
	fmt.Printf("Hex conv of '%d' is '%s'\n", i, h)

	output := &CellOutput{
		Capacity: Uint64(h),
		Lock:     *lock,
		Type:     nil,
	}

	lock1 := &Script{
		Args:     "0x36c329ed630d6ce750712a477543672adab57f4c",
		CodeHash: "0xe3b513a2105a5d4f833d1fad3d968b96b4510687234cd909f86b3ac450d8a2b5",
		HashType: Data,
	}
	i1 := 4900000000000
	h1 := fmt.Sprintf("0x%x", i1)
	fmt.Printf("Hex conv of '%d' is '%s'\n", i1, h1)
	output1 := &CellOutput{
		Capacity: Uint64(h1),
		Lock:     *lock1,
		Type:     nil,
	}

	outputs = append(outputs, *output)
	outputs = append(outputs, *output1)

	transaction.Outputs = outputs
	outputData := []Bytes{}
	outputData = append(outputData, "0x")
	outputData = append(outputData, "0x")

	transaction.OutputsData = outputData
	// witnesses := []Bytes{}
	// witnesses = append(witnesses, "0x")
	// transaction.Witnesses = witnesses
	//fmt.Println(transaction)

	serializeTx, err := transaction.Serialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(hex.EncodeToString(serializeTx))
	fmt.Println("txHash = ")

	txHash, _ := ComputeHash(serializeTx)
	fmt.Println(hex.EncodeToString(txHash))
	fmt.Println("4c02905db773301f73bbc6cd5a400c928caf410bbb13136f6f48bec0a79c22e4")

	// Calc witness args
	// 65 is secp256k1 signature length
}
