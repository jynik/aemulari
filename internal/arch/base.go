package arch

const VARIABLE_INSTR_LEN = 0

type ArchBase struct {
	processor   Type
	mode        Mode
	maxInstrLen uint
	RegisterMap
}

func (b *ArchBase) Type() Type {
	return b.processor
}

func (b *ArchBase) InitialMode() Mode {
	return b.mode
}

func (b *ArchBase) MaxInstrLen() uint {
	return b.maxInstrLen
}
