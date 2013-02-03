package protocol

import (
	"bytes"
	"encoding/hex"
	"net"
	"testing"
	"time"
)

func TestNetworkAddress(t *testing.T) {

	testaddress, _ := net.ResolveTCPAddr("tcp4", "10.0.0.1:8333")
	testtime, _ := time.Parse(time.UnixDate, "Mon Dec 21 02:50:10 UTC 2010")

	tests := []struct {
		raw     *NetworkAddress
		encoded string
	}{
		{
			&NetworkAddress{
				Time: testtime,
				VersionNetworkAddress: VersionNetworkAddress{
					Services: 1,
					Address:  *testaddress,
				},
			},
			"E215104D010000000000000000000000000000000000FFFF0A000001208D",
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

		var decoded NetworkAddress
		i, err = decoded.ReadFrom(bytes.NewReader(encoded))

		if !test.raw.Time.Equal(decoded.Time) {
			t.Errorf("Time mismatch %v", decoded.Time)
		}

		if !bytes.Equal(test.raw.Address.IP, decoded.Address.IP) {
			t.Errorf("IP mismatch %d", hex.EncodeToString(decoded.Address.IP))
		}

		if test.raw.Address.Port != decoded.Address.Port {
			t.Errorf("Port mismatch %d", decoded.Address.Port)
		}

		if !test.raw.Equal(&decoded) {
			t.Errorf("Equals failed")
		}

	}
}
