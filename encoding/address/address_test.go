package address

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncodeDecode(t *testing.T) {

	tests := []struct {
		address string
		pk      string
		pkhash  string
		version byte
	}{
		{
			"16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
			"0450863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B23522CD470243453A299FA9E77237716103ABC11A1DF38855ED6F2EE187E9C582BA6",
			"010966776006953D5567439E5E39F86A0D273BEE",
			0,
		},
	}

	for i, test := range tests {

		pkhash, _ := hex.DecodeString(test.pkhash)
		pk, _ := hex.DecodeString(test.pk)

		decodedPkhash, decodedVersion, err := Decode(test.address)
		if err != nil {
			panic(err)
		}

		if decodedVersion != test.version {
			t.Errorf("TestAddressToPublicKeyHash Decode Version Failed %d", i)
		}

		if !bytes.Equal(decodedPkhash, pkhash) {
			t.Errorf("TestAddressToPublicKeyHash Decode PKHash Failed %d", i)
		}

		encodedAddress := EncodePublicKeyHash(pkhash, test.version)

		if encodedAddress != test.address {
			t.Errorf("TestAddressToPublicKeyHash Encode Address Failed %d", i)
		}

		encodedAddress = EncodePublicKey(pk, test.version)

		if encodedAddress != test.address {
			t.Errorf("TestAddressToPublicKeyHash Encode Address Failed %d", i)
		}

		hashedPk := HashPublicKey(pk)

		if !bytes.Equal(hashedPk, pkhash) {
			t.Errorf("TestAddressToPublicKeyHash Hash Public Key Failed %d", i)
		}

	}
}

func TestAddressChecksumFailure(t *testing.T) {
	in := "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvMX"

	_, _, err := Decode(in)
	if err == nil {
		t.Errorf("Checksum error not detected")
	}
	_, ok := err.(AddressChecksumError)
	if !ok {
		t.Errorf("Checksum error wrong type")
	}
}

func TestAddressUnderrunFailure(t *testing.T) {
	in := "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjv"

	_, _, err := Decode(in)
	if err == nil {
		t.Errorf("Underun error not detected")
	}
	_, ok := err.(AddressUnderunError)
	if !ok {
		t.Errorf("Underrun error wrong type")
	}
}

func BenchmarkEncodePublicKeyHash(b *testing.B) {
	in, _ := hex.DecodeString("010966776006953d5567439e5e39f86a0d273bee")

	for i := 0; i < b.N; i++ {
		EncodePublicKeyHash(in, 0)
	}
}

func BenchmarkDecode(b *testing.B) {
	in := "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM"

	for i := 0; i < b.N; i++ {
		Decode(in)
	}
}
