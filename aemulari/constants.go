package aemulari

type Endianness int // Bit or byte order

const (
	BigEndian    = iota // Big endian bit or byte order
	LittleEndian        // Little endian bit or byte order
)
