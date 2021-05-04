package jsonrpc

import (
	"encoding/json"
	"github.com/liyaojian/hxwallet/types"
)

type Asset struct {
	ID                 string `json:"id"`
	Symbol             string `json:"symbol"`
	Precision          uint8  `json:"precision"`
	Issuer             string `json:"issuer"`
	DynamicAssetDataID string `json:"dynamic_asset_data_id"`
}

type BlockHeader struct {
	TransactionMerkleRoot string            `json:"transaction_merkle_root"`
	Previous              string            `json:"previous"`
	Timestamp             types.Time        `json:"timestamp"`
	Witness               string            `json:"witness"`
	Extensions            []json.RawMessage `json:"extensions"`
}

//todo operation_results
type Block struct {
	TransactionMerkleRoot string              `json:"transaction_merkle_root"`
	Previous              string              `json:"previous"`
	Timestamp             types.Time          `json:"timestamp"`
	Witness               string              `json:"witness"`
	Extensions            []json.RawMessage   `json:"extensions"`
	WitnessSignature      string              `json:"witness_signature"`
	Transactions          []types.Transaction `json:"transactions"`
	BlockId               string              `json:"block_id"`
	SigningKey            string              `json:"signing_key"`
	TransactionIds        []string            `json:"transaction_ids"`
	RefBlockPrefix        uint32              `json:"ref_block_prefix"`
}

type DynamicGlobalProperties struct {
	ID                             string     `json:"id"`
	HeadBlockNumber                uint32     `json:"head_block_number"`
	HeadBlockID                    string     `json:"head_block_id"`
	Time                           types.Time `json:"time"`
	CurrentWitness                 string     `json:"current_witness"`
	NextMaintenanceTime            types.Time `json:"next_maintenance_time"`
	LastBudgetTime                 types.Time `json:"last_budget_time"`
	AccountsRegisteredThisInterval int        `json:"accounts_registered_this_interval"`
	DynamicFlags                   int        `json:"dynamic_flags"`
	RecentSlotsFilled              string     `json:"recent_slots_filled"`
	LastIrreversibleBlockNum       uint32     `json:"last_irreversible_block_num"`
	CurrentAslot                   int64      `json:"current_aslot"`
	WitnessBudget                  int64      `json:"witness_budget"`
	RecentlyMissedCount            int64      `json:"recently_missed_count"`
	Parameters                     string     `json:"parameters"`
}

type Config struct {
	GrapheneSymbol               string `json:"GRAPHENE_SYMBOL"`
	GrapheneAddressPrefix        string `json:"GRAPHENE_ADDRESS_PREFIX"`
	GrapheneMinAccountNameLength uint8  `json:"GRAPHENE_MIN_ACCOUNT_NAME_LENGTH"`
	GrapheneMaxAccountNameLength uint8  `json:"GRAPHENE_MAX_ACCOUNT_NAME_LENGTH"`
	GrapheneMinAssetSymbolLength uint8  `json:"GRAPHENE_MIN_ASSET_SYMBOL_LENGTH"`
	GrapheneMaxAssetSymbolLength uint8  `json:"GRAPHENE_MAX_ASSET_SYMBOL_LENGTH"`
	GrapheneMaxShareSupply       string `json:"GRAPHENE_MAX_SHARE_SUPPLY"`
}

type GlobalProperties struct {
	Properties string
}

func (o *GlobalProperties) UnmarshalJSON(b []byte) error {
	o.Properties = string(b[:])
	return nil
}
