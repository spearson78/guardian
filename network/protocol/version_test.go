package protocol

import (
	"bytes"
	"encoding/hex"
	"net"
	"testing"
	"time"
)

func TestVersion(t *testing.T) {

	nulladdress, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:0")
	testaddress, _ := net.ResolveTCPAddr("tcp4", "78.42.135.86:19000")
	testtime, _ := time.Parse(time.UnixDate, "Sun Feb 03 10:37:13 CET 2013")

	tests := []struct {
		raw     Version
		encoded string
	}{
		{
			Version{
				Version:  60002,
				Services: 1,
				Time:     testtime,
				Receive: VersionNetworkAddress{
					Services: 1,
					Address:  *nulladdress,
				},
				From: VersionNetworkAddress{
					Services: 1,
					Address:  *testaddress,
				},
				Nonce:       0xc33dedd79ce54e68,
				UserAgent:   "/Satoshi:0.7.2/",
				StartHeight: 747,
			},
			"62ea00000100000000000000c92f0e5100000000010000000000000000000000000000000000ffff000000000000010000000000000000000000000000000000ffff4e2a87564a38684ee59cd7ed3dc30f2f5361746f7368693a302e372e322feb020000",
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer

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

		var decoded Version
		i, err = decoded.ReadFrom(bytes.NewReader(expected))
		if err != nil {
			t.Errorf("ReadFrom Failed %v", err)
		}

		if test.raw.Version != decoded.Version {
			t.Errorf("Version mismatch %v", decoded.Version)
		}

		if test.raw.Services != decoded.Services {
			t.Errorf("Services mismatch %v", decoded.Services)
		}

		if !test.raw.Time.Equal(decoded.Time) {
			t.Errorf("Time mismatch %v", decoded.Time)
		}

		if test.raw.Nonce != decoded.Nonce {
			t.Errorf("Nonce mismatch %v", decoded.Nonce)
		}

		if test.raw.UserAgent != decoded.UserAgent {
			t.Errorf("UserAgent mismatch %v", decoded.UserAgent)
		}

		if test.raw.StartHeight != decoded.StartHeight {
			t.Errorf("StartHeight mismatch %v", decoded.StartHeight)
		}

		if !test.raw.Receive.Equal(&decoded.Receive) {
			t.Errorf("Recieve mismatch %v", decoded.Receive)
		}

		if !test.raw.From.Equal(&decoded.From) {
			t.Errorf("From mismatch %v", decoded.From)
		}

	}
}
