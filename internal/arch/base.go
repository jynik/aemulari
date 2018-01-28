package arch

type ArchBase struct {
	archType    Type
	maxInstrLen uint
	defaults    Defaults

	RegisterMap
}

func (b *ArchBase) Defaults() Defaults {
	return b.defaults
}

func (b *ArchBase) Type() Type {
	return b.archType
}

func (b *ArchBase) MaxInstrLen() uint {
	return b.maxInstrLen
}
