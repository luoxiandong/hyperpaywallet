package tronadpter

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/blocktree/tron-adapter/tron/grpc-gateway/core"

	"github.com/blocktree/go-owcrypt"
	"github.com/gogo/protobuf/proto"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// A Client is a Tron RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL string
	// AccessToken string
	Debug  bool
	client *req.Req
}

// NewClient create new client to connect
func NewClient(url, token string, debug bool) *Client {
	c := Client{
		BaseURL: url,
		// AccessToken: token,
		Debug: debug,
	}

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *Client) Call(path string, param interface{}) (*gjson.Result, error) {

	if c == nil || c.client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	url := c.BaseURL + path
	authHeader := req.Header{"Accept": "application/json"}

	r, err := req.Post(url, req.BodyJSON(&param), authHeader)
	if err != nil {
		return nil, err
	}

	if c.Debug {

	}

	if r.Response().StatusCode != http.StatusOK {
		message := gjson.ParseBytes(r.Bytes()).String()
		message = fmt.Sprintf("[%s]%s", r.Response().Status, message)
		return nil, errors.New(message)
	}

	res := gjson.ParseBytes(r.Bytes())
	return &res, nil
}

type Block struct {
	/*
		 {
			 "blockID":"000000000035e1c0f60afaa8387fd17fd9b84fe4381265ff084d739f814558ea",
			 "block_header":{
				 "raw_data":{"number":3531200,
						 "txTrieRoot":"9d98f6cbbde8302774ab87c003831333e132c89e00009cb1f7da35e1e59ae8ca",
						 "witness_address":"41d1dbde8b8f71b48655bec4f6bb532a0142b88bc0",
						 "parentHash":"000000000035e1bfb3b6f244e5316ce408aa8cea4c348eabe2545247f5a4600c",
						 "version":3,
						 "timestamp":1540545240000},
				 "witness_signature":"6ceedcacd8d0111b48eb4131484de3d13f27a2f4bd8156279d03d4208690158e20a641b77d900e026ee33adc328f9ec674f6483ea7b1ca5a27fa24d7fb23964100"
			 },
			 "transactions":[
				 {"ret":[{"contractRet":"SUCCESS"}],
				  "signature":["40aa520f01cebf12948615b9c5a5df5fe7d57e1a7f662d53907b4aa14f647a3a47be2a097fdb58159d0bee7eb1ff0a15ac738f24643fe5114cab8ec0d52cc04d01"],
				  "txID":"ac005b0a195a130914821a6c28db1eec44b4ec3a2358388ceb6c87b677866f1f",
				  "raw_data":{
					  "contract":[
						 {"parameter":{"value":{"amount":1,"asset_name":"48756f6269546f6b656e","owner_address":"416b201fb7b9f2b97bbdaf5e0920191229767c30ee","to_address":"412451d09536fca47760ea6513372bbbbef8583105"},
								   "type_url":"type.googleapis.com/protocol.TransferAssetContract"},
						  "type":"TransferAssetContract"}
					  ],
					  "ref_block_bytes":"e1be",
					  "ref_block_hash":"8dbf5f0cf4c324f2",
					  "expiration":1540545294000,
					  "timestamp":1540545235358}
				 },
				 ...]
		 }
	*/
	Hash   string // 这里采用 BlockID
	Time   int64
	Height uint64 `storm:"id"`
	// Merkleroot        string
	//Confirmations     uint64
}

func (block *Block) GetHeight() uint64 {
	return block.Height
}

// GetNowBlock Done!
// Function：Query the latest block
// 	demo: curl -X POST http://127.0.0.1:8090/wallet/getnowblock
// Parameters：None
// Return value：Latest block on full node
func GetNowBlock() (block *Block, err error) {
	client := NewClient("https://api.trongrid.io", "", true)
	r, err := client.Call("/wallet/getnowblock", nil)

	if err != nil {
		return nil, err
	}

	block = NewBlock(r, false)
	if block.GetBlockHashID() == "" || block.GetHeight() <= 0 {
		return nil, errors.New("GetNowBlock failed: No found <block>")
	}

	// Check for TX
	currstamp := time.Now().UnixNano() / (1000 * 1000) // Unit: ms
	timestamp := int64(block.Time)
	//if timestamp < currstamp-(5*1000) {
	//	wm.Log.Warningf(fmt.Sprintf("Get block timestamp: %d [%+v]", timestamp, time.Unix(timestamp/1000, 0)))
	//	wm.Log.Warningf(fmt.Sprintf("Current d timestamp: %d [%+v]", currstamp, time.Unix(currstamp/1000, 0)))
	//	wm.Log.Warningf("Diff seconds: %ds ", (currstamp-timestamp)/1000)
	//}
	if timestamp < currstamp-(5*60*1000) {
		return nil, errors.New("GetNowBlock returns with unsynced")
	}

	return block, nil
}

func NewBlock(json *gjson.Result, isTestnet bool) *Block {

	header := gjson.Get(json.Raw, "block_header").Get("raw_data")
	// 解析json
	b := &Block{}
	b.Hash = gjson.Get(json.Raw, "blockID").String()
	b.Time = header.Get("timestamp").Int()
	b.Height = header.Get("number").Uint()
	return b
}

func (block *Block) GetBlockHashID() string {
	return block.Hash
}

// according to RFC6979 standard
func tron_hmac_rfc6979_init(key []byte, keylen int) ([]byte, int) {
	k := make([]byte, 32)
	v := make([]byte, 32)
	out := make([]byte, 64)
	tempbuf := make([]byte, 33+keylen)
	//step b in RFC6979
	//copy(v,0x01,32)
	for i := 0; i < 32; i++ {
		v[i] = 0x1
	}
	//step c in RFC6979
	for i := 0; i < 32; i++ {
		k[i] = 0x0
	}
	//step d in RFC6979
	copy(tempbuf[:32], v[:])
	tempbuf[32] = 0
	copy(tempbuf[33:33+keylen], key[:])
	k = owcrypt.Hmac(k, tempbuf, owcrypt.HMAC_SHA256_ALG)

	//step e in RFC6979
	v = owcrypt.Hmac(k, v, owcrypt.HMAC_SHA256_ALG)
	//step f in RFC6979
	copy(tempbuf[:32], v[:])
	tempbuf[32] = 0x01
	k = owcrypt.Hmac(k, tempbuf, owcrypt.HMAC_SHA256_ALG)
	//step g in RFC6979
	v = owcrypt.Hmac(k, v, owcrypt.HMAC_SHA256_ALG)
	retry := 0
	copy(out[:32], k[:])
	copy(out[32:64], v[:])
	//返回k||v,retry
	return out, retry
}

func tron_hmac_rfc6979_gnerate(k, v []byte, retry, nouncelen int) ([]byte, int) {
	nounce := make([]byte, nouncelen)
	j := 0
	if retry == 1 {
		tempbuf := make([]byte, 33)
		copy(tempbuf[:32], v[:])
		//memset(tempbuf,0,1)
		tempbuf[32] = 0
		k = owcrypt.Hmac(k, tempbuf, owcrypt.HMAC_SHA256_ALG)
		v = owcrypt.Hmac(k, v, owcrypt.HMAC_SHA256_ALG)
	}
	for i := 0; i < nouncelen; i += 32 {
		v = owcrypt.Hmac(k[:], v[:], owcrypt.HMAC_SHA256_ALG)
		copy(nounce[(j*32):((j+1)*32)], v[:])
		j++
	}
	retry = 1
	return nounce, retry
}

func tron_nonce_function_rfc6979(msg, key, algo, extradata []byte, counter uint32) []byte {
	keydata := make([]byte, 112)
	nounce := make([]byte, (counter+1)*32)
	copy(keydata[:32], key[:])
	copy(keydata[32:64], msg[:])
	keylen := 64
	if extradata != nil {
		copy(keydata[64:96], extradata[:])
		keylen += 32
	}
	if algo != nil {
		copy(keydata[keylen:keylen+16], algo[:])
	}
	ret, retry := tron_hmac_rfc6979_init(keydata, keylen)

	for i := uint32(0); i <= counter; i++ {
		ReTry := retry
		temp, retry := tron_hmac_rfc6979_gnerate(ret[:32], ret[32:], ReTry, 32)
		copy(nounce[i*32:(i+1)*32], temp[:])
		ReTry = retry
	}

	return nounce

}

func tron_set_uint64(a []byte) []uint64 {
	b := make([]uint64, len(a)>>3)
	b[0] = (uint64(a[0]) << 56) | (uint64(a[1]) << 48) | (uint64(a[2]) << 40) | (uint64(a[3]) << 32) | (uint64(a[4]) << 24) | (uint64(a[5]) << 16) | (uint64(a[6]) << 8) | (uint64(a[7]))
	b[1] = (uint64(a[8]) << 56) | (uint64(a[9]) << 48) | (uint64(a[10]) << 40) | (uint64(a[11]) << 32) | (uint64(a[12]) << 24) | (uint64(a[13]) << 16) | (uint64(a[14]) << 8) | (uint64(a[15]))
	b[2] = (uint64(a[16]) << 56) | (uint64(a[17]) << 48) | (uint64(a[18]) << 40) | (uint64(a[19]) << 32) | (uint64(a[20]) << 24) | (uint64(a[21]) << 16) | (uint64(a[22]) << 8) | (uint64(a[23]))
	b[3] = (uint64(a[24]) << 56) | (uint64(a[25]) << 48) | (uint64(a[26]) << 40) | (uint64(a[27]) << 32) | (uint64(a[28]) << 24) | (uint64(a[29]) << 16) | (uint64(a[30]) << 8) | (uint64(a[31]))
	return b
}

func tron_set_uint32(a []byte) []uint32 {
	b := make([]uint32, len(a)>>2)
	b[0] = (uint32(a[0]) << 24) | (uint32(a[1]) << 16) | (uint32(a[2]) << 8) | (uint32(a[3]))
	b[1] = (uint32(a[4]) << 24) | (uint32(a[5]) << 16) | (uint32(a[6]) << 8) | (uint32(a[7]))
	b[2] = (uint32(a[8]) << 24) | (uint32(a[9]) << 16) | (uint32(a[10]) << 8) | (uint32(a[11]))
	b[3] = (uint32(a[12]) << 24) | (uint32(a[13]) << 16) | (uint32(a[14]) << 8) | (uint32(a[15]))
	b[4] = (uint32(a[16]) << 24) | (uint32(a[17]) << 16) | (uint32(a[18]) << 8) | (uint32(a[19]))
	b[5] = (uint32(a[20]) << 24) | (uint32(a[21]) << 16) | (uint32(a[22]) << 8) | (uint32(a[23]))
	b[6] = (uint32(a[24]) << 24) | (uint32(a[25]) << 16) | (uint32(a[26]) << 8) | (uint32(a[27]))
	b[7] = (uint32(a[28]) << 24) | (uint32(a[29]) << 16) | (uint32(a[30]) << 8) | (uint32(a[31]))
	return b
}

func tron_check_overflow_uint64(a []byte) bool {
	var yes bool
	var no bool
	yes = false
	no = false
	curveOrder := owcrypt.GetCurveOrder(owcrypt.ECC_CURVE_SECP256K1)
	a_uint64 := tron_set_uint64(a)
	curveOrder_uint64 := tron_set_uint64(curveOrder)
	no = no || (a_uint64[0] < curveOrder_uint64[0]) /*no need check for a > check*/
	no = no || (a_uint64[1] < curveOrder_uint64[1])
	yes = yes || (a_uint64[1] > curveOrder_uint64[1]) && (!no)
	no = no || (a_uint64[2] < curveOrder_uint64[2])
	yes = yes || (a_uint64[2] > curveOrder_uint64[2]) && (!no)
	yes = yes || (a_uint64[3] >= curveOrder_uint64[3]) && (!no)
	return yes
}

func tron_check_overflow_uint32(a []byte) bool {
	var yes bool
	var no bool
	yes = false
	no = false
	curveOrder := owcrypt.GetCurveOrder(owcrypt.ECC_CURVE_SECP256K1)
	a_uint32 := tron_set_uint32(a)
	curveOrder_uint32 := tron_set_uint32(curveOrder)
	no = no || (a_uint32[0] < curveOrder_uint32[0]) /*no need check for a > check.*/
	no = no || (a_uint32[1] < curveOrder_uint32[1]) /*no need check for a check. */
	no = no || (a_uint32[2] < curveOrder_uint32[2]) /*no need check for a check.*/
	no = no || (a_uint32[3] < curveOrder_uint32[3])
	yes = yes || (a_uint32[3] > curveOrder_uint32[3]) && (!no)
	no = no || (a_uint32[4] < curveOrder_uint32[4]) && (!yes)
	yes = yes || (a_uint32[4] > curveOrder_uint32[4]) && (!no)
	no = no || (a_uint32[5] < curveOrder_uint32[5]) && (!yes)
	yes = yes || (a_uint32[5] > curveOrder_uint32[5]) && (!no)
	no = no || (a_uint32[6] < curveOrder_uint32[6]) && (!yes)
	yes = yes || (a_uint32[6] > curveOrder_uint32[6]) && (!no)
	yes = yes || (a_uint32[7] >= curveOrder_uint32[7]) && (!no)

	return yes
}

func tron_check_is_zero(a []byte) bool {
	b := make([]uint64, len(a)>>3)
	b = tron_set_uint64(a)
	return ((b[3] == 0) && (b[2] == 0) && (b[1] == 0) && (b[0] == 0))
}

func tron_signatureInner(prikey []byte, hash []byte, nounce []byte) ([]byte, uint16) {
	var recid byte
	signature := make([]byte, 65)
	ret := owcrypt.PreprocessRandomNum(nounce)

	if ret != owcrypt.SUCCESS {
		return nil, ret
	}

	//外部传入随机数，外部已经计算哈希值
	sig, ret := owcrypt.Signature(prikey, nil, 0, hash, 32, owcrypt.ECC_CURVE_SECP256K1|owcrypt.NOUNCE_OUTSIDE_FLAG)
	if ret != owcrypt.SUCCESS {
		return nil, ret
	}
	//判断[nounce]G(G is base point) Y-coordinate 的奇偶性,如果为奇数，recid=0x0;如果为奇数，recid=0x01.
	//这里应该添加判断（签名值r>order,发生的概率接近于1/2^127，几乎为0.这里不再判断，因为底层的C库输出的签名值r已经对order求模数，排除了这种情况）
	yPoint, ret1 := owcrypt.GenPubkey(nounce, owcrypt.ECC_CURVE_SECP256K1)
	if ret1 != owcrypt.SUCCESS {
		return nil, ret1
	}
	if yPoint[63]%2 == 1 {
		recid |= 0x01
	} else {
		recid |= 0x00
	}
	curveOrder := new(big.Int).SetBytes(owcrypt.GetCurveOrder(owcrypt.ECC_CURVE_SECP256K1))
	halfcurveorder := big.NewInt(0)
	s := new(big.Int).SetBytes(sig[32:64])
	divider := big.NewInt(2)
	halfcurveorder.Div(curveOrder, divider)
	sign := s.Cmp(halfcurveorder)
	if sign > 0 {
		s.Sub(curveOrder, s)
		sByte := s.Bytes()
		if len(sByte) < 32 {
			for i := 0; i < 32-len(sByte); i++ {
				sByte = append([]byte{0x00}, sByte...)
			}
		}
		copy(sig[32:64], sByte)
		recid ^= 1
	}
	copy(signature[:64], sig[:])
	signature[64] = recid
	return signature, ret
}

/*
@function:Tron signature(ECDSA&&secp256k1)
@paramter[in]prikey pointer to private key
@paramter[in]hash pointer to the hash of message(Transaction txt)
@parameter[out]the first part is signature(r||s||v,total 65 byte);
the second part
*/
func TronSignature(prikey []byte, hash []byte) ([]byte, uint16) {
	signature := make([]byte, 65)
	//	var recid byte
	var ret uint16
	var counter uint32
	counter = 0
	if len(hash) != 32 {
		return nil, owcrypt.FAILURE
	}
	prikey_overflow := tron_check_overflow_uint64(prikey)
	prikey_IsZero := tron_check_is_zero(prikey)
	if !prikey_overflow && !prikey_IsZero {
		for {
			nounce := tron_nonce_function_rfc6979(hash, prikey, nil, nil, counter)
			nounce_overflow := tron_check_overflow_uint32(nounce)
			nounce_IsZero := tron_check_is_zero(nounce)
			if !nounce_overflow && !nounce_IsZero {
				signature, ret = tron_signatureInner(prikey, hash, nounce)
				if ret == owcrypt.SUCCESS {
					break
				}
			}
			counter++
		}
	}

	return signature, ret
}

func ELFSignature(prikey []byte, hash []byte) ([]byte, uint16) {
	signature := make([]byte, 65)
	//	var recid byte
	var ret uint16
	var counter uint32
	counter = 0
	// if len(hash) != 32 {
	// 	return nil, owcrypt.FAILURE
	// }
	prikey_overflow := tron_check_overflow_uint64(prikey)
	prikey_IsZero := tron_check_is_zero(prikey)
	if !prikey_overflow && !prikey_IsZero {
		for {
			nounce := tron_nonce_function_rfc6979(hash, prikey, nil, nil, counter)
			nounce_overflow := tron_check_overflow_uint32(nounce)
			nounce_IsZero := tron_check_is_zero(nounce)
			if !nounce_overflow && !nounce_IsZero {
				signature, ret = tron_signatureInner(prikey, hash, nounce)
				if ret == owcrypt.SUCCESS {
					break
				}
			}
			counter++
		}
	}

	return signature, ret
}

func InsertSignatureIntoRawTransaction(txHex string, signature string) (string, error) {

	tx := &core.Transaction{}
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return "", err
	}
	if err := proto.Unmarshal(txBytes, tx); err != nil {
		return "", err
	}
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		//log.Errorf("invalid transaction signature hex data;unexpected err:%v", err)
		return "", fmt.Errorf("invalid signature hex data")
	}

	tx.Signature = append(tx.Signature, signatureBytes)
	x, err := proto.Marshal(tx)
	if err != nil {
		//wm.Log.Info("marshal tx failed;unexpected error:%v", err)
		return "", err
	}

	mergeTxHex := hex.EncodeToString(x)
	return mergeTxHex, nil

}
