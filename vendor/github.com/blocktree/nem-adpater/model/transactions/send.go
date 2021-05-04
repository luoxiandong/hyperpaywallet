package transactions

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/blocktree/nem-adpater/extras"
	"github.com/blocktree/nem-adpater/model"
	"github.com/blocktree/nem-adpater/utils"

	"github.com/pkg/errors"
)

type Common struct {
	Password, PrivateKey string
	IsHW                 bool
}
type RequestAnnounce struct {
	Data      string `json:"data"`
	Signature string `json:"signature"`
}

// Serialize a transaction and broadcast it to the network
// param common - A common struct
// param entity - A prepared transaction struct
// param endpoint - An NIS endpoint struct
// return - An announce transaction promise of the com.requests service
func Send(common Common, entity interface{}) (string, error) {
	if extras.IsEmpty(common) || extras.IsEmpty(entity) {
		return "", errors.New("Missing parameter !")
	}
	if len(common.PrivateKey) != 64 && len(common.PrivateKey) != 66 {
		return "", errors.New("Invalid private key, length must be 64 or 66 characters !")
	}
	if !utils.IsHexadecimal(common.PrivateKey) {
		return "", errors.New("Private key must be hexadecimal only !")
	}
	kp := model.KeyPairCreate(common.PrivateKey)

	result := utils.SerializeTransaction(entity)
	signature := kp.Sign(string(result))
	obj := &RequestAnnounce{
		Data:      utils.Bt2Hex([]byte(result)),
		Signature: utils.Bt2Hex(signature),
	}
	payload, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	fmt.Println("Data = ", utils.Bt2Hex([]byte(result)), " Signature = ", utils.Bt2Hex(signature))
	fmt.Println("payload = ", hex.EncodeToString(payload[:]))

	return utils.Bt2Hex(payload[:]), nil
}
