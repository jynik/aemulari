package aemulari

type Endianness int // Describes of bit or byte order

const (
	BigEndian    = iota // Big endian bit or byte order
	LittleEndian        // Little endian bit or byte order
)
