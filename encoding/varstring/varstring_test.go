package varstring

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestVarInt(t *testing.T) {

	tests := []struct {
		raw     string
		encoded string
	}{

		{"", "00"},
		{"A", "0141"},
	}

	for _, test := range tests {
		var b bytes.Buffer

		i, err := WriteVarString(&b, test.raw)
		if err != nil {
			t.Errorf("WriteVarString error %v", err)
		}

		buf := b.Bytes()
		expected, _ := hex.DecodeString(test.encoded)
		if i != len(expected) {
			t.Errorf("Written byte count mismatch for %s expected %d got %d", test.raw, len(expected), i)
		}
		if !bytes.Equal(buf, expected) {
			t.Errorf("bytes mismatch for %s expected %s got %s", test.raw, hex.EncodeToString(expected), hex.EncodeToString(buf))
		}
		decoded, x, err := ReadVarString(bytes.NewReader(buf))
		if err != nil {
			t.Errorf("ReadVarString error %v", err)
		}
		if x != i {
			t.Errorf("Read byte count mismatch for %s expected %d got %d", test.raw, i, x)
		}

		if decoded != test.raw {
			t.Errorf("Decoded mismatch for %s got %s", test.raw, decoded)
		}
	}
}
