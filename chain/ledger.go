package ledger

type Contract struct {
	Binary  []byte
	Address string
}

func (l *Ledger) StoreContract(address string, binary []byte) error {
	// Persist the contract binary
	l.contractStorage[address] = &Contract{
		Binary:  binary,
		Address: address,
	}
	return nil
}

func (l *Ledger) GetContract(address string) (*Contract, error) {
	contract, exists := l.contractStorage[address]
	if !exists {
		return nil, fmt.Errorf("contract not found")
	}
	return contract, nil
}

