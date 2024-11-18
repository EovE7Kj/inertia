package wasm

import (
	"fmt"
	"github.com/zenon-network/go-zenon/node"
)

// BlockchainAPI exposes blockchain functionalities to WASM
type BlockchainAPI struct {
	node *node.Node // Access to blockchain state
}

// NewBlockchainAPI initializes the API layer
func NewBlockchainAPI(n *node.Node) *BlockchainAPI {
	return &BlockchainAPI{node: n}
}

// Host function: Retrieve account balance
func (api *BlockchainAPI) GetBalance(address string) (uint64, error) {
	account := api.node.Ledger.GetAccount(address)
	if account == nil {
		return 0, fmt.Errorf("account not found")
	}
	return account.Balance, nil
}

// Host function: Write data to state
func (api *BlockchainAPI) WriteState(key, value string) error {
	return api.node.Ledger.UpdateState(key, value)
}

