package arch

const VARIABLE_INSTR_LEN = 0

type ArchBase struct {
	archType    Type
	archMode	Mode
	maxInstrLen uint
	RegisterMap
}

func (b *ArchBase) Type() Type {
	return b.archType
}

func (b *ArchBase) MaxInstrLen() uint {
	return b.maxInstrLen
}
