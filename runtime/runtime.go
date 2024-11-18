package wasm

import (
        "fmt"
        "io/ioutil"

        "github.com/wasmerio/wasmer-go/wasmer"
)

type WASMRuntime struct {
        engine  *wasmer.Engine
        store   *wasmer.Store
        module  *wasmer.Module
        imports *wasmer.ImportObject
        qasMeter *QasMeter
}

var gasCosts = map[string]uint64{
    "write_state": 200,
}

// initialize runtime
func NewWASMRuntime(qasLimit uint64) *WASMRuntime {
        engine := wasmer.NewEngine()
        store := wasmer.NewStore(engine)
        return &WASMRuntime{
                engine: engine,
                store:  store,
                qasMeter: NewQasMeter(qasLimit),
        }
}

func (r *WASMRuntime) LoadContract(filepath string) error {
        wasmBytes, err := ioutil.ReadFile(filepath)
        if err != nil {
                return fmt.Errorf("failed to read WASM file: %w", err)
        }

        module, err := wasmer.NewModule(r.store, wasmBytes)
        if err != nil {
                return fmt.Errorf("failed to compile WASM module: %w", err)
        }

        r.module = module
        return nil
}

func HandleContractInvocation(tx *Transaction, api *BlockchainAPI, runtime *WASMRuntime) error {
        // Load contract
        err := runtime.LoadContract(tx.ContractBinary)
        if err != nil {
                return fmt.Errorf("failed to load contract: %w", err)
        }

        // Expose API
        runtime.ExposeAPI(api)

        // Execute function
        _, err = runtime.ExecuteContract(tx.FunctionName, tx.Args...)
        if err != nil {
                return fmt.Errorf("contract execution failed: %w", err)
        }

        return nil
}

func (r *WASMRuntime) ExposeAPI(api *BlockchainAPI) {
    r.imports.Register("env", map[string]wasmer.Int32ToVoid{
        "write_state": func(key, value int32) {
            cost := gasCosts["write_state"]
            if err := r.qasMeter.ConsumeQas(cost); err != nil {
                fmt.Printf("Error: %v\n", err)
                return
            }
            api.WriteState(int32ToString(key), int32ToString(value))
        },
    })
}


func (r *WASMRuntime) ExecuteContract(functionName string, args ...interface{}) ([]byte, error) {
    if err := r.qasMeter.ConsumeQas(100); err != nil {
        return nil, err
    }

    instance, err := wasmer.NewInstance(r.module, r.imports)
    if err != nil {
        return nil, fmt.Errorf("failed to instantiate WASM module: %w", err)
    }
    defer instance.Close()

    function, err := instance.Exports.GetFunction(functionName)
    if err != nil {
        return nil, fmt.Errorf("function '%s' not found: %w", functionName, err)
    }

    result, err := function(args...)
    if err != nil {
        return nil, fmt.Errorf("failed to execute function '%s': %w", functionName, err)
    }

    if err := r.qasMeter.ConsumeQas(uint64(len(result.([]byte)))); err != nil {
        return nil, err
    }

    return result.([]byte), nil
}

