package wasm

type QasMeter struct {
	qasLimit uint64
	qasUsed  uint64
}

func NewQasMeter(limit uint64) *QasMeter {
	return &QasMeter{qasLimit: limit, qasUsed: 0}
}

func (gm *QasMeter) ConsumeQas(amount uint64) error {
	if gm.qasUsed+amount > gm.qasLimit {
		return fmt.Errorf("out of qas")
	}
	gm.qasUsed += amount
	return nil
}

func (gm *QasMeter) RemainingQas() uint64 {
	return gm.qasLimit - gm.qasUsed
}

func (gm *QasMeter) Reset(limit uint64) {
    gm.qasLimit = limit
    gm.qasUsed = 0
}

