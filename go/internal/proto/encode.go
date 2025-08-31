package proto

import "encoding/binary"

// Varint encodes a uint64 into protobuf varint form.
func Varint(x uint64) []byte {
	buf := make([]byte, 0, 10)
	for x >= 0x80 {
		buf = append(buf, byte(x)|0x80)
		x >>= 7
	}
	buf = append(buf, byte(x))
	return buf
}

// VarintField encodes a varint field with number num.
func VarintField(num int, value uint64) []byte {
	key := uint64(num<<3 | 0)
	b := Varint(key)
	b = append(b, Varint(value)...)
	return b
}

// Fixed64Field encodes a 64-bit field with number num.
func Fixed64Field(num int, value uint64) []byte {
	key := uint64(num<<3 | 1)
	b := Varint(key)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, value)
	return append(b, buf...)
}

// BytesField encodes a length-delimited field with number num.
func BytesField(num int, data []byte) []byte {
	key := uint64(num<<3 | 2)
	b := Varint(key)
	b = append(b, Varint(uint64(len(data)))...)
	b = append(b, data...)
	return b
}
