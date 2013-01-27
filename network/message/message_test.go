package message

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestMessage(t *testing.T) {

	tests := []struct {
		raw      Message
		checksum uint32
		encoded  string
	}{

		{
			Message{
				Magic:   MagicMainNet,
				Command: "addr",
				Payload: []byte{0x01, 0xe2, 0x15, 0x10, 0x4d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x0A, 0x00, 0x00, 0x01, 0x20, 0x8D},
			},
			0x9B3952ED,
			"f9beb4d96164647200000000000000001F000000ED52399B01E215104D010000000000000000000000000000000000FFFF0A000001208D",
		},
		{
			Message{
				Magic:   MagicTestNet,
				Command: "version",
				Payload: nil,
			},
			3806393949,
			"0b11090776657273696f6e0000000000000000005df6e0e2",
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer

		if test.raw.Checksum() != test.checksum {
			t.Errorf("Checksum Mismatch %d", test.raw.Checksum())
		}

		i, err := test.raw.WriteTo(&buf)
		if err != nil {
			t.Errorf("WriteTo Failed %v", err)
		}

		encoded := buf.Bytes()

		if int64(len(encoded)) != i {
			t.Errorf("WriteTo written bytes mismatch len encoded %d writtenBytes %d", len(encoded), i)
		}

		expected, _ := hex.DecodeString(test.encoded)

		if !bytes.Equal(encoded, expected) {
			t.Errorf("WriteTo bytes mismatch expected %s encoded %s", hex.EncodeToString(expected), hex.EncodeToString(encoded))
		}

		var decoded Message
		i, err = decoded.ReadFrom(bytes.NewReader(encoded))

		if decoded.Magic != test.raw.Magic {
			t.Errorf("Magic mismatch %d", decoded.Magic)
		}

		if decoded.Checksum() != test.checksum {
			t.Errorf("Checksum Mismatch %d", test.raw.Checksum())
		}

		if decoded.Command != test.raw.Command {
			t.Errorf("Command mismatch %d", decoded.Magic)
		}

		if !bytes.Equal(decoded.Payload, test.raw.Payload) {
			t.Errorf("Command mismatch %d", decoded.Magic)
		}
	}
}
