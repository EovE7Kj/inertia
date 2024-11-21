package wasm

import (
	"github.com/wasmerio/wasmer-go/wasmer"
	"fmt"
	"io/ioutil"
)

type Contract struct {
	Module  *wasmer.Module
	Store   *wasmer.Store
	Imports *wasmer.ImportObject
}

func NewContract(wasmPath string) (*Contract, error) {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	wasmBytes, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read wasm file: %w", err)
	}

	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile wasm module: %w", err)
	}

	return &Contract{
		Module:  module,
		Store:   store,
		Imports: wasmer.NewImportObject(),
	}, nil
}

