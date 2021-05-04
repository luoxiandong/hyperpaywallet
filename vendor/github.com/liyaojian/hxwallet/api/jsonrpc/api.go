package jsonrpc

import (
	"encoding/json"
	"fmt"
	"github.com/liyaojian/hxwallet/rpc"
)

type API struct {
	caller rpc.Caller
	id     rpc.APIID
}

func NewAPI(id rpc.APIID, caller rpc.Caller) *API {
	return &API{id: id, caller: caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	err := api.caller.Connect()
	if err != nil {
		return err
	}
	return api.caller.Call(api.id, method, args, reply)
}

func (api *API) setCallback(method string, callback func(raw json.RawMessage)) error {
	return api.caller.SetCallback(api.id, method, callback)
}

// GET ChainId of entry point
func (api *API) GetChainId() (string, error) {
	var resp string
	err := api.call("get_chain_id", rpc.EmptyParams, &resp)
	return resp, err
}

// Gets dynamic global properties of current blockchain
func (api *API) GetDynamicGlobalProperties() (*DynamicGlobalProperties, error) {
	var resp DynamicGlobalProperties
	err := api.call("get_dynamic_global_properties", rpc.EmptyParams, &resp)
	return &resp, err
}

// Get block by block height
func (api *API) GetBlock(blockNum uint32) (*Block, error) {
	var resp Block
	err := api.call("get_block", []interface{}{blockNum}, &resp)
	return &resp, err
}

//lookup_asset_symbols
func (api *API) GetAssets(symbols ...string) ([]*Asset, error) {
	var resp []*Asset
	err := api.call("lookup_asset_symbols", []interface{}{symbols}, &resp)
	return resp, err
}

// LookupAssetSymbols get assets corresponding to the provided symbol or IDs
func (api *API) GetAsset(symbol string) (*Asset, error) {
	var resp []*Asset
	if err := api.call("lookup_asset_symbols", []interface{}{[]string{symbol}}, &resp); err != nil {
		return nil, err
	}
	if resp[0] == nil {
		return nil, fmt.Errorf("assets %s not exist", symbol)
	}
	return resp[0], nil
}

func (api *API) getGlobalProperties() (*GlobalProperties, error) {
	var resp *GlobalProperties
	err := api.call("get_global_properties", rpc.EmptyParams, &resp)
	return resp, err
}
