package wasm

import (
    "github.com/wasmerio/wasmer-go/wasmer"
)

type HostFunction = func([]wasmer.Value, []wasmer.Value) error

type ImportsManager struct {
    qasMeter *QasMeter
    api      *BlockchainAPI
    imports  *wasmer.ImportObject
}

func NewImportsManager(qasMeter *QasMeter, api *BlockchainAPI) *ImportsManager {
    return &ImportsManager{
        qasMeter: qasMeter,
        api:      api,
        imports:  wasmer.NewImportObject(),
    }
}

func (im *ImportsManager) RegisterHostFunctions() {
    im.imports.Register("env", map[string]interface{}{
        "write_state": im.wrapWithGas(func(args []wasmer.Value, results []wasmer.Value) error {
            key := int32(args[0].I32())
            value := int32(args[1].I32())
            im.api.WriteState(im.int32ToString(key), im.int32ToString(value))
            return nil
        }, 200),
        "read_state": im.wrapWithGas(func(args []wasmer.Value, results []wasmer.Value) error {
            key := int32(args[0].I32())
            stateValue, err := im.api.ReadState(im.int32ToString(key))
            if err != nil {
                return err
            }
            results[0] = wasmer.NewI32(im.stringToInt32(stateValue))
            return nil
        }, 100),
        "log_message": im.wrapWithGas(func(args []wasmer.Value, results []wasmer.Value) error {
            message := im.int32ToString(int32(args[0].I32()))
            im.api.Log(message)
            return nil
        }, 50),
    })
}

func (im *ImportsManager) GetImports() *wasmer.ImportObject {
    return im.imports
}

func (im *ImportsManager) wrapWithGas(hostFunc HostFunction, gasCost uint64) HostFunction {
    return func(args []wasmer.Value, results []wasmer.Value) error {
        if err := im.qasMeter.ConsumeQas(gasCost); err != nil {
            return err
        }
        return hostFunc(args, results)
    }
}

func (im *ImportsManager) int32ToString(ptr int32) string {
    return "example"
}

func (im *ImportsManager) stringToInt32(value string) int32 {
    return 1234
}

