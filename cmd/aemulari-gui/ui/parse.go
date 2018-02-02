package ui

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	ae "../../../aemulari"
)

func hexSeq(s string) ([]byte, bool) {
	if !strings.HasPrefix(s, "{") || !strings.HasSuffix(s, "}") {
		return []byte{}, false
	}

	ret, err := hex.DecodeString(s[1 : len(s)-1])
	return ret, err == nil
}

func u8Bytes(s string) ([]byte, bool) {
	if !strings.HasPrefix(s, "u8(") || !strings.HasSuffix(s, ")") {
		return []byte{}, false
	}

	val, err := strconv.ParseInt(s[3:len(s)-1], 0, 8)
	if err != nil {
		return []byte{}, false
	}

	return []byte{byte(val & 0xff)}, true
}

func i8Bytes(s string) ([]byte, bool) {
	if !strings.HasPrefix(s, "i8(") || !strings.HasSuffix(s, ")") {
		return []byte{}, false
	}

	val, err := strconv.ParseUint(s[3:len(s)-1], 0, 8)
	if err != nil {
		return []byte{}, false
	}

	return []byte{uint8(val)}, true
}

func u16Endian(val uint16, e ae.Endianness) []byte {
	var ret []byte = make([]byte, 2)

	if e == ae.BigEndian {
		binary.BigEndian.PutUint16(ret, val)
	} else {
		binary.LittleEndian.PutUint16(ret, val)
	}

	return ret
}

func u16Bytes(s string, e ae.Endianness) ([]byte, bool) {
	if !strings.HasPrefix(s, "u16(") || !strings.HasSuffix(s, ")") {
		return []byte{}, false
	}

	val, err := strconv.ParseUint(s[4:len(s)-1], 0, 16)
	if err != nil {
		return []byte{}, false
	}

	return u16Endian(uint16(val), e), true
}

func i16Bytes(s string, e ae.Endianness) ([]byte, bool) {
	if !strings.HasPrefix(s, "i16(") || !strings.HasSuffix(s, ")") {
		return []byte{}, false
	}

	val, err := strconv.ParseInt(s[4:len(s)-1], 0, 16)
	if err != nil {
		return []byte{}, false
	}

	return u16Endian(uint16(val), e), true
}

func u32Endian(val uint32, e ae.Endianness) []byte {
	var ret []byte = make([]byte, 4)

	if e == ae.BigEndian {
		binary.BigEndian.PutUint32(ret, val)
	} else {
		binary.LittleEndian.PutUint32(ret, val)
	}

	return ret
}

func u32Bytes(s string, e ae.Endianness) ([]byte, bool) {
	if !strings.HasPrefix(s, "u32(") || !strings.HasSuffix(s, ")") {
		return []byte{}, false
	}

	val, err := strconv.ParseUint(s[4:len(s)-1], 0, 32)
	if err != nil {
		return []byte{}, false
	}

	var ret []byte = make([]byte, 4)

	if e == ae.BigEndian {
		binary.BigEndian.PutUint32(ret, uint32(val))
	} else {
		binary.LittleEndian.PutUint32(ret, uint32(val))
	}

	return ret, true
}

func i32Bytes(s string, e ae.Endianness) ([]byte, bool) {
	if !strings.HasPrefix(s, "i32(") || !strings.HasSuffix(s, ")") {
		return []byte{}, false
	}

	val, err := strconv.ParseInt(s[4:len(s)-1], 0, 32)
	if err != nil {
		return []byte{}, false
	}

	var ret []byte = make([]byte, 4)

	if e == ae.BigEndian {
		binary.BigEndian.PutUint32(ret, uint32(val))
	} else {
		binary.LittleEndian.PutUint32(ret, uint32(val))
	}

	return ret, true
}

func u64Endian(val uint64, e ae.Endianness) []byte {
	var ret []byte = make([]byte, 8)

	if e == ae.BigEndian {
		binary.BigEndian.PutUint64(ret, val)
	} else {
		binary.LittleEndian.PutUint64(ret, val)
	}

	return ret
}

func u64Bytes(s string, e ae.Endianness) ([]byte, bool) {
	if !strings.HasPrefix(s, "u64(") || !strings.HasSuffix(s, ")") {
		return []byte{}, false
	}

	val, err := strconv.ParseUint(s[4:len(s)-1], 0, 64)
	if err != nil {
		return []byte{}, false
	}

	return u64Endian(val, e), true
}

func i64Bytes(s string, e ae.Endianness) ([]byte, bool) {
	if !strings.HasPrefix(s, "i64(") || !strings.HasSuffix(s, ")") {
		return []byte{}, false
	}

	val, err := strconv.ParseInt(s[4:len(s)-1], 0, 64)
	if err != nil {
		return []byte{}, false
	}

	return u64Endian(uint64(val), e), true
}

func UBytes(s string, e ae.Endianness) ([]byte, bool) {
	val, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return []byte{}, false
	}

	if val < 256 {
		return []byte{byte(val)}, true
	} else if val < 65535 {
		return u16Endian(uint16(val), e), true
	} else if val < 4294967296 {
		return u32Endian(uint32(val), e), true
	}

	return u64Endian(val, e), true
}

func IBytes(s string, e ae.Endianness) ([]byte, bool) {
	val, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return []byte{}, false
	}

	if val >= -128 && val <= 127 {
		return []byte{byte(val)}, true
	} else if val >= -32768 && val <= 32767 {
		return u16Endian(uint16(val), e), true
	} else if val >= -2147483648 && val <= 2147483647 {
		return u32Endian(uint32(val), e), true
	}

	return u64Endian(uint64(val), e), true
}

func parseValue(valStr string, e ae.Endianness) ([]byte, error) {
	if ret, ok := hexSeq(valStr); ok {
		return ret, nil
	} else if ret, ok = u8Bytes(valStr); ok {
		return ret, nil
	} else if ret, ok = i8Bytes(valStr); ok {
		return ret, nil
	} else if ret, ok = u16Bytes(valStr, e); ok {
		return ret, nil
	} else if ret, ok = i16Bytes(valStr, e); ok {
		return ret, nil
	} else if ret, ok = u32Bytes(valStr, e); ok {
		return ret, nil
	} else if ret, ok = i32Bytes(valStr, e); ok {
		return ret, nil
	} else if ret, ok = u64Bytes(valStr, e); ok {
		return ret, nil
	} else if ret, ok = i64Bytes(valStr, e); ok {
		return ret, nil
	}

	if strings.HasPrefix(valStr, "-") {
		if ret, ok := IBytes(valStr, e); ok {
			return ret, nil
		}
	} else {
		if ret, ok := UBytes(valStr, e); ok {
			return ret, nil
		}
	}

	return []byte{}, fmt.Errorf("\"%s\" is not a valid value.", valStr)
}
