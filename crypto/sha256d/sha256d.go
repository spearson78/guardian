//Package sha256d implements the double SHA256 hash algorithm as used in the Bitcoin protocol
package sha256d

import (
	"crypto"
	"crypto/sha256"
	"hash"
)

// The size of a SHA256d checksum in bytes
const Size = sha256.Size

type sha256d struct {
	sha256 hash.Hash
	round2 hash.Hash
	dirty  bool
}

// New returns a new hash.Hash computing the SHA256d checksum.
func New() hash.Hash {
	var h = new(sha256d)
	h.sha256 = crypto.SHA256.New()
	h.round2 = crypto.SHA256.New()
	return h
}

func (h *sha256d) Write(p []byte) (n int, err error) {
	return h.sha256.Write(p)
}

func (h *sha256d) Sum(b []byte) []byte {
	var result [Size]byte
	firstRound := h.sha256.Sum(result[:0])
	if h.dirty {
		h.round2.Reset()
	}
	h.round2.Write(firstRound)
	h.dirty = true
	return h.round2.Sum(b)
}

func (h *sha256d) Reset() {
	h.sha256.Reset()
}

func (h *sha256d) Size() int {
	return h.sha256.Size()
}

func (h *sha256d) BlockSize() int {
	return h.sha256.BlockSize()
}
