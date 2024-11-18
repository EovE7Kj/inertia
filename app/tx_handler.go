package app

import (
	"github.com/zenon-network/go-zenon/types"
	"github.com/zenon-network/go-zenon/wasm"
)

func (n *Node) ProcessTransaction(tx *types.Transaction) error {
	switch tx.Type {
	case types.TxTypeDeployContract:
		return n.handleContractDeployment(tx)
	case types.TxTypeInvokeContract:
		return n.handleContractInvocation(tx)
	default:
		// Existing transaction types...
		return nil
	}
}

func (n *Node) handleContractDeployment(tx *types.Transaction) error {
	// Store the WASM binary in the ledger
	err := n.Ledger.StoreContract(tx.Sender, tx.Data)
	if err != nil {
		return err
	}

	log.Printf("Contract deployed by %s", tx.Sender)
	return nil
}

func (n *Node) handleContractInvocation(tx *types.Transaction) error {
	// Retrieve the contract binary
	contract, err := n.Ledger.GetContract(tx.ContractAddress)
	if err != nil {
		return err
	}

	// Load and execute the contract
	err = n.WASMRuntime.LoadContract(contract.Binary)
	if err != nil {
		return err
	}

	_, err = n.WASMRuntime.ExecuteContract(tx.FunctionName, tx.Args...)
	if err != nil {
		return err
	}

	log.Printf("Contract %s invoked successfully", tx.ContractAddress)
	return nil
}

