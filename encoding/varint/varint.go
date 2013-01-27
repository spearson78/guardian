package varint

import (
	"encoding/binary"
	"io"
)

const MaxBufferSize = 9

func PutVarInt(buf []byte, value uint64) int {

	i := 0

	switch {
	case value < uint64(0xfd):
		buf[0] = byte(value)
		i = 1
	case value <= uint64(65535):
		buf[0] = 0xfd

		binary.LittleEndian.PutUint16(buf[1:3], uint16(value))
		i = 3
	case value <= uint64(4294967295):
		buf[0] = 0xfe

		binary.LittleEndian.PutUint32(buf[1:5], uint32(value))
		i = 5
	default:
		buf[0] = 0xff

		binary.LittleEndian.PutUint64(buf[1:9], value)
		i = 9
	}

	return i
}

func WriteVarInt(w io.Writer, value uint64) (int, error) {
	buf := make([]byte, MaxBufferSize)

	i := PutVarInt(buf, value)

	return w.Write(buf[:i])
}

func ReadVarInt(r io.Reader) (uint64, int, error) {

	var bufarray [MaxBufferSize]byte
	buf := bufarray[:]

	var value uint64
	i := 0

	i, err := r.Read(buf[0:1])
	if err != nil {
		return 0, i, err
	}

	switch buf[0] {
	default:
		value = uint64(buf[0])
		i = 1
	case 0xfd:
		_, err := r.Read(buf[0:2])
		if err != nil {
			return 0, i, err
		}

		value = uint64(binary.LittleEndian.Uint16(buf))
		i = 3
	case 0xfe:
		_, err := r.Read(buf[0:4])
		if err != nil {
			return 0, i, err
		}
		value = uint64(binary.LittleEndian.Uint32(buf))
		i = 5
	case 0xff:
		_, err := r.Read(buf[0:8])
		if err != nil {
			return 0, i, err
		}
		value = binary.LittleEndian.Uint64(buf)
		i = 9
	}

	return value, i, nil
}

func VarInt(buf []byte) (uint64, int) {

	var value uint64
	i := 0

	switch buf[0] {
	default:
		value = uint64(buf[0])
		i = 1
	case 0xfd:
		value = uint64(binary.LittleEndian.Uint16(buf[1:3]))
		i = 3
	case 0xfe:
		value = uint64(binary.LittleEndian.Uint32(buf[1:5]))
		i = 5
	case 0xff:
		value = binary.LittleEndian.Uint64(buf[1:9])
		i = 9
	}

	return value, i
}
