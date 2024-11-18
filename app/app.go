package app

import (
	"log"
	"github.com/zenon-network/go-zenon/wasm"
)

func (n *Node) InitializeWASM() {
	// Initialize Blockchain API
	api := wasm.NewBlockchainAPI(n)

	// Initialize WASM runtime
	n.WASMRuntime = wasm.NewWASMRuntime()
	n.WASMRuntime.ExposeAPI(api)

	log.Println("WASM runtime initialized successfully")
}

