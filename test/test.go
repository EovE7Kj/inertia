package wasm_test

import (
    "fmt"
    "testing"
    "github.com/wasmerio/wasmer-go/wasmer"
)

// Mock for "write_state"
type MockAPI struct{}

func (m *MockAPI) WriteState(key, value string) {
    fmt.Printf("Mock WriteState called with key: %s, value: %s\n", key, value)
}

func TestWASMContractExecution(t *testing.T) {
    engine := wasmer.NewEngine()
    store := wasmer.NewStore(engine)

    // Load the precompiled WASM contract
    wasmBytes, err := ioutil.ReadFile("simple_contract.wasm")
    if err != nil {
        t.Fatalf("Failed to read WASM file: %v", err)
    }

    module, err := wasmer.NewModule(store, wasmBytes)
    if err != nil {
        t.Fatalf("Failed to compile WASM module: %v", err)
    }

    // Create an import object (mock the host function)
    importObject := wasmer.NewImportObject()
    importObject.Register("env", map[string]wasmer.Int32ToVoid{
        "write_state": func(key, value int32) {
            fmt.Printf("Mock write_state called with key: %d, value: %d\n", key, value)
        },
    })

    // Create an instance of the module
    instance, err := wasmer.NewInstance(module, importObject)
    if err != nil {
        t.Fatalf("Failed to instantiate WASM module: %v", err)
    }
    defer instance.Close()

    // Call the "run" function from the WASM module
    runFunc, err := instance.Exports.GetFunction("run")
    if err != nil {
        t.Fatalf("Failed to get 'run' function: %v", err)
    }

    // Execute the function
    result, err := runFunc()
    if err != nil {
        t.Fatalf("Failed to execute 'run' function: %v", err)
    }

    // Check result
    if result != 42 {
        t.Fatalf("Expected result 42, got %d", result)
    }
}

