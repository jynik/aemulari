package aemulari

type archBase struct {
	processor   processorType
	mode        processorMode
	maxInstrLen uint
	registerMap
}

func (b *archBase) id() processorType {
	return b.processor
}

func (b *archBase) initialMode() processorMode {
	return b.mode
}

func (b *archBase) maxInstructionSize() uint {
	return b.maxInstrLen
}
