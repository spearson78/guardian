package base58

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestEncodeDecodeBitcoin(t *testing.T) {

	tests := []struct {
		raw     string
		encoded string
	}{

		{"", ""},
		{"61", "2g"},
		{"626262", "a3gV"},
		{"636363", "aPEr"},
		{"73696d706c792061206c6f6e6720737472696e67", "2cFupjhnEsSn59qHXstmK2ffpLv2"},
		{"00eb15231dfceb60925886b67d065299925915aeb172c06647", "1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L"},
		{"516b6fcd0f", "ABnLTmg"},
		{"bf4f89001e670274dd", "3SEo3LWLoPntC"},
		{"572e4794", "3EFU7m"},
		{"ecac89cad93923c02321", "EJDM8drfXA6uyA"},
		{"10c8511e", "Rt5zm"},
		{"00000000000000000000", "1111111111"},
		{"00010966776006953D5567439E5E39F86A0D273BEED61967F6", "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM"},
	}

	for i, test := range tests {

		rawBytes, err := hex.DecodeString(test.raw)
		if err != nil {
			panic(err)
		}

		encoded, err := BitcoinEncoding.Encode(rawBytes)
		if err != nil {
			panic(err)
		}

		if string(encoded) != test.encoded {
			t.Errorf("Bitcoin Encode Failed %d Encoded %s Expected %s", i, string(encoded), test.encoded)
			return
		}

		decoded, err := BitcoinEncoding.Decode(encoded)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(rawBytes, decoded) {
			t.Errorf("Bitcoin Decode Failed %d Decoded %s Expected %s", i, hex.EncodeToString(decoded), test.raw)
			return
		}
	}
}

func TestEncodeDecodeFlickr(t *testing.T) {

	tests := []struct {
		raw     int64
		encoded string
	}{
		{4847765341, "8oo4K4"},
		{6857269519, "brXijP"},
		{0, ""},
	}

	for i, test := range tests {

		rawBigInt := big.NewInt(test.raw)

		encoded, err := FlickrEncoding.EncodeInt(rawBigInt)
		if err != nil {
			panic(err)
		}

		if string(encoded) != test.encoded {
			t.Errorf("Encode Failed %d Encoded %s Expected %s", i, string(encoded), test.encoded)
			return
		}

		decoded, err := FlickrEncoding.DecodeToInt(encoded)
		if err != nil {
			panic(err)
		}

		if decoded.Int64() != test.raw {
			t.Errorf("Bitcoin Decode Failed %d Decoded %d Expected %d", i, decoded.Int64(), test.raw)
			return
		}
	}
}

func TestDecodeIllegalChar(t *testing.T) {
	in := []byte("16UwLL9RisO3QfPqBUvKofHmBQ7wMtjvM")

	_, err := BitcoinEncoding.Decode(in)
	if err == nil {
		t.Errorf("Illegal char error not returned")
	}
	_, ok := err.(CorruptDecodeInputError)
	if !ok {
		t.Errorf("Wrong illegal char error type returned")
	}
}

func TestEncodeIllegalInt(t *testing.T) {
	in := big.NewInt(-1)

	_, err := BitcoinEncoding.EncodeInt(in)
	if err == nil {
		t.Errorf("Illegal int error not returned")
	}
	_, ok := err.(CorruptEncodeInputError)
	if !ok {
		t.Errorf("Wrong illegal int error type returned")
	}
}

func BenchmarkEncodeBitcoin(b *testing.B) {

	in, err := hex.DecodeString("00010966776006953D5567439E5E39F86A0D273BEED61967F6")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		BitcoinEncoding.Encode(in)
	}
}

func BenchmarkEncodeFlickr(b *testing.B) {

	in := big.NewInt(4847765341)

	for i := 0; i < b.N; i++ {
		FlickrEncoding.EncodeInt(in)
	}
}

func BenchmarkDecodeBitcoin(b *testing.B) {
	in := []byte("16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM")

	for i := 0; i < b.N; i++ {
		BitcoinEncoding.Decode(in)
	}
}

func BenchmarkDecodeFlickr(b *testing.B) {
	in := []byte("8oo4K4")

	for i := 0; i < b.N; i++ {
		FlickrEncoding.DecodeToInt(in)
	}
}
