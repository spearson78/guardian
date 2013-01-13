package sha256d

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestNewSHA256d(t *testing.T) {

	in, err := hex.DecodeString("00010966776006953D5567439E5E39F86A0D273BEE")
	if err != nil {
		panic(err)
	}

	out, err := hex.DecodeString("D61967F63C7DD183914A4AE452C9F6AD5D462CE3D277798075B107615C1A8A30")
	if err != nil {
		panic(err)
	}
	h := New()
	h.Write(in)
	result := h.Sum(nil)
	if !bytes.Equal(result, out) {
		t.Errorf("Failed")
	}
}

func TestResetSHA256d(t *testing.T) {

	in, err := hex.DecodeString("00010966776006953D5567439E5E39F86A0D273BEE")
	if err != nil {
		panic(err)
	}

	out, err := hex.DecodeString("D61967F63C7DD183914A4AE452C9F6AD5D462CE3D277798075B107615C1A8A30")
	if err != nil {
		panic(err)
	}

	h := New()
	for i := 0; i < 10; i++ {
		h.Reset()
		h.Write(in)
		result := h.Sum(nil)
		if !bytes.Equal(result, out) {
			t.Errorf("Failed %d", i)
		}
	}

}

func TestSmallWritesSHA256d(t *testing.T) {

	in, err := hex.DecodeString("00010966776006953D5567439E5E39F86A0D273BEE")
	if err != nil {
		panic(err)
	}
	out, err := hex.DecodeString("D61967F63C7DD183914A4AE452C9F6AD5D462CE3D277798075B107615C1A8A30")
	if err != nil {
		panic(err)
	}

	h := New()
	for i := 0; i < len(in); i++ {
		h.Write(in[i : i+1])
	}
	result := h.Sum(nil)
	if !bytes.Equal(result, out) {
		t.Errorf("Failed")
	}

}

func TestSmallWritesIntermediateSumsSHA256d(t *testing.T) {

	in, err := hex.DecodeString("00010966776006953D5567439E5E39F86A0D273BEE")
	if err != nil {
		panic(err)
	}
	out, err := hex.DecodeString("D61967F63C7DD183914A4AE452C9F6AD5D462CE3D277798075B107615C1A8A30")
	if err != nil {
		panic(err)
	}

	h := New()
	for i := 0; i < len(in); i++ {
		h.Write(in[i : i+1])
		h.Sum(nil)
	}
	result := h.Sum(nil)
	if !bytes.Equal(result, out) {
		t.Errorf("Failed")
	}

}

func BenchmarkNewSHA256d(b *testing.B) {

	in, err := hex.DecodeString("00010966776006953D5567439E5E39F86A0D273BEE")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		h := New()
		h.Write(in)
		_ = h.Sum(nil)

	}
}

func BenchmarkResetSHA256d(b *testing.B) {

	in, err := hex.DecodeString("00010966776006953D5567439E5E39F86A0D273BEE")
	if err != nil {
		panic(err)
	}
	h := New()

	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(in)
		_ = h.Sum(nil)

	}
}

func Benchmark160(b *testing.B) {
	in := make([]byte, 160)
	h := New()

	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(in)
		_ = h.Sum(nil)

	}
}

func Benchmark256(b *testing.B) {
	in := make([]byte, 256)
	h := New()

	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(in)
		_ = h.Sum(nil)

	}
}

func Benchmark1k(b *testing.B) {

	in := make([]byte, 1024)
	h := New()

	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(in)
		_ = h.Sum(nil)

	}
}

func Benchmark8k(b *testing.B) {

	in := make([]byte, 8192)
	h := New()

	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(in)
		_ = h.Sum(nil)
	}
}
