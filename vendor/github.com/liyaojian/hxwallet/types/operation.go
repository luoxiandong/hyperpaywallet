package types

import "strings"

type Asset struct {
	HxAmount  int64  `json:"amount"`
	HxAssetId string `json:"asset_id"`
}

// hx  --- "1.3.0"
// btc --- "1.3.1"
// ltc --- "1.3.2"
// hc  --- "1.3.3"
func DefaultAsset() Asset {
	return Asset{
		0,
		"1.3.0",
	}
}

func (asset *Asset) SetAssetBySymbol(symbol string) {
	symbol = strings.ToUpper(symbol)

	if symbol == "HX" {
		asset.HxAssetId = "1.3.0"
	} else if symbol == "BTC" {
		asset.HxAssetId = "1.3.1"
	} else if symbol == "LTC" {
		asset.HxAssetId = "1.3.2"
	} else if symbol == "HC" {
		asset.HxAssetId = "1.3.3"
	}
}

type Memo struct {
	HxFrom    string `json:"from"` //public_key_type  33
	HxTo      string `json:"to"`   //public_key_type  33
	HxNonce   uint64 `json:"nonce"`
	HxMessage string `json:"message"`

	IsEmpty bool   `json:"-"`
	Message string `json:"-"`
}

func DefaultMemo() Memo {
	return Memo{
		"HX1111111111111111111111111111111114T1Anm",
		"HX1111111111111111111111111111111114T1Anm",
		0,
		"",
		true,
		"",
	}
}

type Authority struct {
	HxWeightThreshold uint32          `json:"weight_threshold"`
	HxAccountAuths    []interface{}   `json:"account_auths"`
	HxKeyAuths        [][]interface{} `json:"key_auths"`
	HxAddressAuths    []interface{}   `json:"address_auths"`

	KeyAuths string `json:"-"`
}

func DefaultAuthority() Authority {

	return Authority{
		1,
		[]interface{}{},
		[][]interface{}{{"", 1}},
		[]interface{}{},
		"",
	}
}

type AccountOptions struct {
	HxMemoKey            string        `json:"memo_key"`
	HxVotingAccount      string        `json:"voting_account"`
	HxNumWitness         uint16        `json:"num_witness"`
	HxNumCommittee       uint16        `json:"num_committee"`
	HxVotes              []interface{} `json:"votes"`
	HxMinerPledgePayBack byte          `json:"miner_pledge_pay_back"`
	HxExtensions         []interface{} `json:"extensions"`
}

func DefaultAccountOptions() AccountOptions {

	return AccountOptions{
		"",
		"1.2.5",
		0,
		0,
		[]interface{}{},
		10,
		[]interface{}{},
	}

}

// transfer operation tag is  0
type TransferOperation struct {
	HxFee         Asset  `json:"fee"`
	HxGuaranteeId string `json:"guarantee_id,omitempty"`
	HxFrom        string `json:"from"`
	HxTo          string `json:"to"`

	HxFromAddr string `json:"from_addr"`
	HxToAddr   string `json:"to_addr"`

	HxAmount Asset `json:"amount"`
	HxMemo   *Memo `json:"memo,omitempty"`

	HxExtensions []interface{} `json:"extensions"`
}

func DefaultTransferOperation() *TransferOperation {

	return &TransferOperation{
		DefaultAsset(),
		"",
		"1.2.0",
		"1.2.0",
		"",
		"",
		DefaultAsset(),
		nil,
		make([]interface{}, 0),
	}
}

// account bind operation tag is 10
type AccountBindOperation struct {
	HxFee              Asset  `json:"fee"`
	HxCrosschainType   string `json:"crosschain_type"`
	HxAddr             string `json:"addr"`
	HxAccountSignature string `json:"account_signature"`
	HxTunnelAddress    string `json:"tunnel_address"`
	HxTunnelSignature  string `json:"tunnel_signature"`
	HxGuaranteeId      string `json:"guarantee_id,omitempty"`
}

func DefaultAccountBindOperation() *AccountBindOperation {

	return &AccountBindOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

// account unbind operation tag is 11
type AccountUnBindOperation struct {
	HxFee              Asset  `json:"fee"`
	HxCrosschainType   string `json:"crosschain_type"`
	HxAddr             string `json:"addr"`
	HxAccountSignature string `json:"account_signature"`
	HxTunnelAddress    string `json:"tunnel_address"`
	HxTunnelSignature  string `json:"tunnel_signature"`
}

func DefaultAccountUnBindOperation() *AccountUnBindOperation {

	return &AccountUnBindOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
	}
}

// withdraw cross chain operation tag is 61
type WithdrawCrosschainOperation struct {
	HxFee             Asset  `json:"fee"`
	HxWithdrawAccount string `json:"withdraw_account"`
	HxAmount          string `json:"amount"`
	HxAssetSymbol     string `json:"asset_symbol"`

	HxAssetId           string `json:"asset_id"`
	HxCrosschainAccount string `json:"crosschain_account"`
	HxMemo              string `json:"memo"`
}

func DefaultWithdrawCrosschainOperation() *WithdrawCrosschainOperation {

	return &WithdrawCrosschainOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

//register account operation tag is 5
type RegisterAccountOperation struct {
	HxFee             Asset     `json:"fee"`
	HxRegistrar       string    `json:"registrar"`
	HxReferrer        string    `json:"referrer"`
	HxReferrerPercent uint16    `json:"referrer_percent"`
	HxName            string    `json:"name"`
	HxOwner           Authority `json:"owner"`
	HxActive          Authority `json:"active"`
	HxPayer           string    `json:"payer"`

	HxOptions     AccountOptions `json:"options"`
	HxExtensions  interface{}    `json:"extensions"`
	HxGuaranteeId string         `json:"guarantee_id,omitempty"`
}

func DefaultRegisterAccountOperation() *RegisterAccountOperation {

	return &RegisterAccountOperation{
		DefaultAsset(),
		"1.2.0",
		"1.2.0",
		0,
		"",
		DefaultAuthority(),
		DefaultAuthority(),
		"",

		DefaultAccountOptions(),
		make(map[string]interface{}, 0),
		"",
	}

}

//lock balance operation tag is 55
type LockBalanceOperation struct {
	HxLockAssetId     string `json:"lock_asset_id"`
	HxLockAssetAmount int64  `json:"lock_asset_amount"`
	HxContractAddr    string `json:"contract_addr"`

	HxLockBalanceAccount string `json:"lock_balance_account"`
	HxLocktoMinerAccount string `json:"lockto_miner_account"`
	HxLockBalanceAddr    string `json:"lock_balance_addr"`

	HxFee Asset `json:"fee"`
}

func DefaultLockBalanceOperation() *LockBalanceOperation {

	return &LockBalanceOperation{
		"1.3.0",
		0,
		"",
		"",
		"",
		"",
		DefaultAsset(),
	}
}

//foreclose balance operation tag is 56
type ForecloseBalanceOperation struct {
	HxFee Asset `json:"fee"`

	HxForecloseAssetId     string `json:"foreclose_asset_id"`
	HxForecloseAssetAmount int64  `json:"foreclose_asset_amount"`

	HxForecloseMinerAccount string `json:"foreclose_miner_account"`
	HxForecloseContractAddr string `json:"foreclose_contract_addr"`

	HxForecloseAccount string `json:"foreclose_account"`
	HxForecloseAddr    string `json:"foreclose_addr"`
}

func DefaultForecloseBalanceOperation() *ForecloseBalanceOperation {

	return &ForecloseBalanceOperation{
		DefaultAsset(),
		"1.3.0",
		0,
		"",
		"",
		"",
		"",
	}
}

//obtain pay back operation tag is 73
type ObtainPaybackOperation struct {
	HxPayBackOwner   string          `json:"pay_back_owner"`
	HxPayBackBalance [][]interface{} `json:"pay_back_balance"`
	HxGuaranteeId    string          `json:"guarantee_id,omitempty"`
	HxFee            Asset           `json:"fee"`

	citizenName []string
	obtainAsset []Asset
}

func DefaultObtainPaybackOperation() *ObtainPaybackOperation {

	return &ObtainPaybackOperation{
		"",
		[][]interface{}{{"", DefaultAsset()}},
		"",
		DefaultAsset(),
		nil,
		nil,
	}
}

// contract invoke operation tag is 79
type ContractInvokeOperation struct {
	HxFee          Asset  `json:"fee"`
	HxInvokeCost   uint64 `json:"invoke_cost"`
	HxGasPrice     uint64 `json:"gas_price"`
	HxCallerAddr   string `json:"caller_addr"`
	HxCallerPubkey string `json:"caller_pubkey"`
	HxContractId   string `json:"contract_id"`
	HxContractApi  string `json:"contract_api"`
	HxContractArg  string `json:"contract_arg"`
	//Hx_extension     []interface{} `json:"extensions"`
	HxGuaranteeId string `json:"guarantee_id,omitempty"`
}

func DefaultContractInvokeOperation() *ContractInvokeOperation {
	return &ContractInvokeOperation{
		DefaultAsset(),
		0,
		0,
		"",
		"",
		"",
		"",
		"",
		//make([]interface{}, 0),
		"",
	}
}

// transfer to contract operation tag is 81
type ContractTransferOperation struct {
	HxFee          Asset  `json:"fee"`
	HxInvokeCost   uint64 `json:"invoke_cost"`
	HxGasPrice     uint64 `json:"gas_price"`
	HxCallerAddr   string `json:"caller_addr"`
	HxCallerPubkey string `json:"caller_pubkey"`
	HxContractId   string `json:"contract_id"`
	HxAmount       Asset  `json:"amount"`
	HxParam        string `json:"param"`
	//Hx_extension     []interface{} `json:"extensions"`
	HxGuaranteeId string `json:"guarantee_id,omitempty"`
}

func DefaultContractTransferOperation() *ContractTransferOperation {
	return &ContractTransferOperation{
		DefaultAsset(),
		0,
		0,
		"",
		"",
		"",
		DefaultAsset(),
		"",
		//make([]interface{}, 0),
		"",
	}
}
