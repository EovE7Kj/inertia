package wasm

import (
	"fmt"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type Runtime struct {
	Instance *wasmer.Instance
	Store    *wasmer.Store
	Imports  *wasmer.ImportObject
	QasMeter *QasMeter // Gas metering
	API      BlockchainAPI
}

func (c *Contract) CreateRuntime(api BlockchainAPI, qasLimit uint64) (*Runtime, error) {
	// Initialize the gas meter
	qasMeter := NewQasMeter(qasLimit)

	// Add API functions as imports
	c.AddBlockchainAPI(api)

	instance, err := wasmer.NewInstance(c.Module, c.Imports)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %w", err)
	}

	return &Runtime{
		Instance: instance,
		Store:    c.Store,
		Imports:  c.Imports,
		QasMeter: qasMeter,
		API:      api,
	}, nil
}

func (r *Runtime) Execute(functionName string, args ...interface{}) ([]byte, error) {
	// Deduct Qas for execution
	if err := r.QasMeter.ConsumeQas(100); err != nil {
		return nil, fmt.Errorf("gas limit exceeded: %w", err)
	}

	// Retrieve the function from exports
	function, err := r.Instance.Exports.GetFunction(functionName)
	if err != nil {
		return nil, fmt.Errorf("function '%s' not found: %w", functionName, err)
	}

	// Execute the function with provided arguments
	result, err := function(args...)
	if err != nil {
		return nil, fmt.Errorf("execution of '%s' failed: %w", functionName, err)
	}

	// Handle result conversion
	switch v := result.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("unexpected result type: %T", v)
	}
}

func (c *Contract) AddBlockchainAPI(api BlockchainAPI) {
	c.Imports.Register("env", map[string]interface{}{
		"write_state": func(keyPtr, valuePtr int32) {
			key := int32ToString(keyPtr)
			value := int32ToString(valuePtr)
			if err := api.WriteState(key, value); err != nil {
				panic("state write failed")
			}
		},
		"read_state": func(keyPtr, valuePtr int32) int32 {
			key := int32ToString(keyPtr)
			value, err := api.ReadState(key)
			if err != nil {
				panic("state read failed")
			}
			return stringToInt32(value)
		},
	})
}

