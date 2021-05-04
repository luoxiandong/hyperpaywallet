package hdwallet

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ed25519"
	//"encoding/hex"
	//"encoding/json"
	//"fmt"
	//"github.com/WaykiChain/wicc-wallet-utils-go/commons"
	//"github.com/btcsuite/btcutil"
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ugorji/go/codec"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
	"hash/crc32"
)

func init() {
	coins[ADA] = newADA
}

const NetType = MainNet

const (
	MainNet                    = 1
	TestNet                    = 2
	ProtocolMagicsMainnet      = 764824073
	ProtocolMagicsTestnet      = 1097911063
	txSignMessagePrefixTestnet = "011a4170cb175820"
	txSignMessagePrefixMainnet = "011a2d964a095820"
)

type ada struct {
	name   string
	symbol string
	key    *Key
}

func newADA(key *Key) Wallet {
	return &ada{
		name:   "Cardano",
		symbol: "ADA",
		key:    key,
	}
}

type TxInput struct {
	TxHash      string
	OutputIndex uint32
	Coins       uint64
}

func (ti *TxInput) EncodeCBOR() []byte {
	txHashHex, _ := hex.DecodeString(ti.TxHash)
	return EncodeCBOR([]interface{}{
		0,
		&codec.RawExt{
			Tag: 24,
			Value: EncodeCBOR([]interface{}{
				txHashHex,
				ti.OutputIndex,
			}),
		},
	})
}

type TxAux struct {
	Inputs     []TxInput
	Outputs    []TxOutput
	Attributes interface{}
}

func (ta *TxAux) GetID() []byte {
	hash := blake2b.Sum256(ta.EncodeCBOR())
	return hash[:]
}

func (ta *TxAux) EncodeCBOR() []byte {
	inputs := make([]interface{}, len(ta.Inputs))
	for i, v := range ta.Inputs {
		inputs[i] = v
	}
	outputs := make([]interface{}, len(ta.Outputs))
	for i, v := range ta.Outputs {
		outputs[i] = v
	}

	inputsHex := CborIndefiniteLengthArray(inputs)
	outputsHex := CborIndefiniteLengthArray(outputs)

	var data []byte = []byte{0x83}
	data = append(data, inputsHex...)
	data = append(data, outputsHex...)
	data = append(data, 0xa0)

	return data
}

type TxWitness struct {
	Signature         []byte
	ExtendedPublicKey []byte
}

func (tw *TxWitness) EncodeCBOR() []byte {
	return EncodeCBOR([]interface{}{
		0,
		&codec.RawExt{
			Tag: 24,
			Value: EncodeCBOR([]interface{}{
				tw.ExtendedPublicKey,
				tw.Signature,
			}),
		},
	})
}

type TxOutput struct {
	Address string
	Coins   uint64
}

func (to *TxOutput) EncodeCBOR() []byte {
	addrHex := base58.Decode(to.Address)
	coinsHex := EncodeCBOR(to.Coins)

	var data []byte = []byte{0x82}
	data = append(data, addrHex...)
	data = append(data, coinsHex...)

	return data
}

func (c *ada) GetKey() *Key {
	return c.key
}

func (c *ada) GetType() uint32 {
	return c.key.opt.CoinType
}

func (c *ada) GetName() string {
	return c.name
}

func (c *ada) GetSymbol() string {
	return c.symbol
}

func (c *ada) GetAddress() (string, error) {

	public := *c.key.PublicED25519
	priv := *c.key.PrivateED25519
	privateKey := hex.EncodeToString(priv)
	cc := privGetChainCode(privateKey)
	xpub := append(public, cc[:]...)

	address := addressString(xpub)

	return address, nil
}

func (c *ada) GetWalletAccountFromWif() (*WalletAccount, error) {
	hexPri := c.GetKey().Wif
	if len(hexPri) > 0 {
		privb, _ := hex.DecodeString(c.key.Wif)
		privateKey := ed25519.PrivateKey(privb)
		public := privateKey.Public().(ed25519.PublicKey)

		chainCode := privGetChainCode(c.key.Wif)
		xpub := append(public, chainCode[:]...)

		address := addressString(xpub)
		return &WalletAccount{
			Res:        1,
			Address:    address,
			PrivateKey: hexPri,
			PublicKey:  hex.EncodeToString(public),
		}, nil
	}
	return &WalletAccount{
		Res:        0,
		PrivateKey: hexPri,
	}, nil
}

func (c *ada) CreateRawTransaction(signIn *SignInput) (*SignResult, error) {
	return &SignResult{
		Res: 0,
	}, nil
}

func (c *ada) GenerateTxHash(signIn *SignInput) (*TxHashResult, error) {
	return &TxHashResult{}, nil
}

func (c *ada) SignTxHash(signIn *SignTxHashInput) (*TxHashResult, error) {
	return &TxHashResult{}, nil
}

func (c *ada) GetWalletAccount() *WalletAccount {
	if c.GetKey().Extended == nil {
		return &WalletAccount{
			Res: 0,
		}
	}

	address, err := c.GetAddress()
	if err != nil {
		return &WalletAccount{
			Res:    0,
			ErrMsg: err.Error(),
		}
	}
	pri := *c.GetKey().PrivateED25519
	publicKey := *c.GetKey().PublicED25519

	return &WalletAccount{
		Res:        1,
		PrivateKey: hex.EncodeToString(pri),
		Address:    address,
		PublicKey:  hex.EncodeToString(publicKey),
		Seed:       c.GetKey().Seed,
	}
}

func (c *ada) SignRawTransaction(signIn *SignInput) (*SignResult, error) {
	var (
		txInputs       []TxInput
		txOutputs      []TxOutput
		balance, coins uint64
	)

	var vins []OutPutItem
	json.Unmarshal(signIn.Inputs, &vins)
	fmt.Println("vins : ", vins)

	amount := uint64(signIn.Amount)
	fee := uint64(signIn.Fee)
	change := uint64(signIn.Change)
	var utxo []TxInput
	for _, vin := range vins {
		var input TxInput
		input.TxHash = vin.TxHash
		input.OutputIndex = vin.Vout
		input.Coins = uint64(vin.Value)
		utxo = append(utxo, input)
	}

	for _, in := range utxo {
		if coins < amount+fee+change {
			coins = coins + in.Coins
			txInputs = append(txInputs, in)
		}
		balance = balance + in.Coins
	}

	if coins < amount+fee+change {
		return nil, fmt.Errorf("balance is not enough")
	}

	txOutputs = append(txOutputs, TxOutput{Address: signIn.DestAddr, Coins: amount})

	if change > 0 {
		txOutputs = append(txOutputs, TxOutput{Address: signIn.SrcAddr, Coins: change})
	}

	txAux := TxAux{
		Inputs:     txInputs,
		Outputs:    txOutputs,
		Attributes: nil,
	}

	privb, _ := hex.DecodeString(signIn.PrivateKey)
	privateKey := ed25519.PrivateKey(privb)
	public := privateKey.Public().(ed25519.PublicKey)

	chainCode := privGetChainCode(signIn.PrivateKey)

	xpub := append(public, chainCode[:]...)
	txSignedStructured := signTxGetStructured(&txAux, privateKey, xpub)

	txHash := txSignedStructured.GetID()
	txBody := txSignedStructured.EncodeCBOR()

	return &SignResult{
		Res:    1,
		RawTX:  hex.EncodeToString(txBody),
		TxHash: hex.EncodeToString(txHash),
	}, nil
}

func seedGetChainCode(seed []byte) []byte {
	mac := hmac.New(sha512.New, []byte("Bitcoin seed"))
	mac.Write(seed)
	I := mac.Sum(nil)
	chainCode := I[len(I)/2:]

	return chainCode
}

func privGetChainCode(priv string) []byte {
	seed, _ := hex.DecodeString(priv[:64])
	return seedGetChainCode(seed)
}

func addressString(xpub []byte) string {
	addrType := 0
	addrAttributes := make(map[interface{}]interface{})
	//addrAttributes[1] = cborEncode("fa0500f0d3512d248b2ec3dd945a38dc2d65fae844f4bcb3df4927cb")
	switch NetType {
	case TestNet:
		addrAttributes[2] = cborEncode(ProtocolMagicsTestnet)
	default:

	}
	// variables
	v := []interface{}{
		addrType,
		[]interface{}{
			addrType,
			xpub,
		},
		addrAttributes,
	}

	encAddr := cborEncode(v)

	// compute addrHash
	h := sha3.Sum256(encAddr)

	b, _ := blake2b.New(28, nil)
	b.Write(h[:])
	addrHash := b.Sum(nil)

	// crc encoding
	addr := []interface{}{
		addrHash,
		addrAttributes,
		addrType,
	}

	s := cborEncode(addr)
	crc := crc32.ChecksumIEEE(s)

	cwid := cborEncode([]interface{}{
		//cbor.Tag
		&codec.RawExt{
			Tag: 24,
			//Data: s,
			Value: s,
		},
		crc,
	})

	return base58.Encode(cwid)
}

func cborEncode(v interface{}) []byte {
	var buf bytes.Buffer

	ch := &codec.CborHandle{}
	e := codec.NewEncoder(&buf, ch)
	e.MustEncode(v)

	return buf.Bytes()
}

type Transaction struct {
	Aux       *TxAux
	Witnesses []TxWitness
}

func (tx *Transaction) GetID() []byte {
	return tx.Aux.GetID()
}

func (tx *Transaction) EncodeCBOR() []byte {
	size := len(tx.Witnesses)
	buf := make([]interface{}, size)
	for i := 0; i < size; i++ {
		buf[i] = 0
	}
	bufEnc := EncodeCBOR(buf)
	prefix := bufEnc[:len(bufEnc)-size]

	witnesses := prefix
	for _, witness := range tx.Witnesses {
		witnesses = append(witnesses, witness.EncodeCBOR()...)
	}

	aux := tx.Aux.EncodeCBOR()

	var data []byte = []byte{0x82}
	data = append(data, aux...)
	data = append(data, witnesses...)

	return data
}

func signTxGetStructured(txAux *TxAux, privateKey ed25519.PrivateKey, xpub []byte) *Transaction {
	var txSignMessagePrefix string
	switch NetType {
	case TestNet:
		txSignMessagePrefix = txSignMessagePrefixTestnet
	default:
		txSignMessagePrefix = txSignMessagePrefixMainnet
	}
	prefix, _ := hex.DecodeString(txSignMessagePrefix)
	txHash := txAux.GetID()
	message := append(prefix, txHash[:]...)

	var (
		witnesses []TxWitness
	)

	for _, _ = range txAux.Inputs {
		signature := ed25519.Sign(privateKey, message)
		witnesses = append(witnesses, TxWitness{
			Signature:         signature,
			ExtendedPublicKey: xpub,
		})
	}

	return &Transaction{
		Aux:       txAux,
		Witnesses: witnesses,
	}
}

func EncodeCBOR(v interface{}) []byte {
	var buf bytes.Buffer
	ch := &codec.CborHandle{}
	e := codec.NewEncoder(&buf, ch)
	e.MustEncode(v)
	return buf.Bytes()
}

func CborIndefiniteLengthArray(elements []interface{}) []byte {
	var (
		data = []byte{0x9f} // indefinite array prefix
	)
	for _, e := range elements {
		if v, ok := e.(TxInput); ok {
			data = append(data, v.EncodeCBOR()...)
		} else if v, ok := e.(TxOutput); ok {
			data = append(data, v.EncodeCBOR()...)
		}
	}
	data = append(data, 0xff) // end of array
	return data
}
