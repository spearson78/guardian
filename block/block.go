package block

import (
	"encoding/binary"
	"errors"
	"github.com/spearson78/guardian/crypto/sha256d"
)

type BlockHash [32]byte
type MerkleHash [32]byte

type Block struct {
	Version       uint32     //Block version information, based upon the software version creating this block 
	PrevBlockHash BlockHash  // 	The hash value of the previous block this particular block references
	MerkleRoot    MerkleHash // 	The reference to a Merkle tree collection which is a hash of all transactions related to this block
	Timestamp     uint32     // 	A unix timestamp recording when this block was created (Currently limited to dates before the year 2106!)
	Bits          uint32     // 	The calculated difficulty target being used for this block
	Nonce         uint32     // 	The nonce used to generate this block… to allow variations of the header and compute different hashes
}

type UnsupportedBlockVersionError uint32

func (e UnsupportedBlockVersionError) Error() string {
	return "Unsupported Block Version" + string(e)
}

func (this *Block) Set(data []byte) error {

	if len(data) != 80 {
		return errors.New("block data must be length 80")
	}

	this.Version = binary.LittleEndian.Uint32(data[0:4])

	if this.Version > 2 {
		return UnsupportedBlockVersionError(this.Version)
	}

	copy(this.PrevBlockHash[:], data[4:36])
	copy(this.MerkleRoot[:], data[36:68])
	this.Timestamp = binary.LittleEndian.Uint32(data[68:72])
	this.Bits = binary.LittleEndian.Uint32(data[72:76])
	this.Nonce = binary.LittleEndian.Uint32(data[76:80])

	return nil
}

func (this *Block) Bytes() ([]byte, error) {
	if this.Version > 2 {
		return nil, UnsupportedBlockVersionError(this.Version)
	}

	encoded := make([]byte, 80)

	binary.LittleEndian.PutUint32(encoded[0:4], this.Version)
	copy(encoded[4:36], this.PrevBlockHash[:])
	copy(encoded[36:68], this.MerkleRoot[:])
	binary.LittleEndian.PutUint32(encoded[68:72], this.Timestamp)
	binary.LittleEndian.PutUint32(encoded[72:76], this.Bits)
	binary.LittleEndian.PutUint32(encoded[76:80], this.Nonce)

	return encoded, nil
}

func (this *Block) Hash() (BlockHash, error) {

	var blockHash BlockHash

	h := sha256d.New()
	blockData, err := this.Bytes()
	if err != nil {
		return blockHash, err
	}
	_, err = h.Write(blockData)
	if err != nil {
		return blockHash, err
	}

	h.Sum(blockHash[:0])

	return blockHash, nil
}
