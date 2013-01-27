package varint

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"
)

func TestVarInt(t *testing.T) {

	tests := []struct {
		raw     uint64
		encoded string
	}{

		{0, "00"},
		{1, "01"},
		{0xfc, "fc"},
		{0xfd, "fdfd00"},
		{0xfe, "fdfe00"},
		{65534, "fdfeff"},
		{65535, "fdffff"},
		{65536, "fe00000100"},
		{4294967294, "fefeffffff"},
		{4294967295, "feffffffff"},
		{4294967296, "ff0000000001000000"},
		{math.MaxUint64, "ffffffffffffffffff"},
	}

	buf := make([]byte, MaxBufferSize)

	for _, test := range tests {
		i := PutVarInt(buf, test.raw)
		expected, _ := hex.DecodeString(test.encoded)
		if i != len(expected) {
			t.Errorf("Written byte count mismatch for %d expected %d got %d", test.raw, len(expected), i)
		}
		if !bytes.Equal(buf[:i], expected) {
			t.Errorf("bytes mismatch for %d expected %s got %s", test.raw, hex.EncodeToString(expected), hex.EncodeToString(buf[:i]))
		}
		decoded, x := VarInt(buf)
		if x != i {
			t.Errorf("Read byte count mismatch for %d expected %d got %d", test.raw, i, x)
		}

		if decoded != test.raw {
			t.Errorf("Decoded mismatch for %d got %d", test.raw, decoded)
		}

		decoded, x, err := ReadVarInt(bytes.NewBuffer(buf))
		if err != nil {
			t.Errorf("ReadVarInt error %v", err)
		}
		if x != i {
			t.Errorf("Read byte count mismatch for %d expected %d got %d", test.raw, i, x)
		}

		if decoded != test.raw {
			t.Errorf("Decoded mismatch for %d got %d", test.raw, decoded)
		}
	}
}

func BenchmarkPutVarInt(b *testing.B) {

	buf := make([]byte, MaxBufferSize)

	for i := 0; i < b.N; i++ {
		PutVarInt(buf, math.MaxUint64)
	}
}

func BenchmarkVarInt(b *testing.B) {

	data, _ := hex.DecodeString("ffffffffffffffffff")

	for i := 0; i < b.N; i++ {
		VarInt(data)
	}
}

func BenchmarkReadVarInt(b *testing.B) {

	data, _ := hex.DecodeString("ffffffffffffffffff")
	buf := bytes.NewReader(data)

	for i := 0; i < b.N; i++ {
		buf.Seek(0, 0)
		ReadVarInt(buf)
	}
}
