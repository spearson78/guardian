package transaction

import (
	"bytes"
	"encoding/binary"
	"github.com/spearson78/guardian/crypto/sha256d"
	"github.com/spearson78/guardian/encoding/varint"
	"io"
)

type TxHash [32]byte

type OutPoint struct {
	Hash  TxHash
	Index uint32
}

type Tx struct {
	Version uint32

	Inputs  []TxIn
	Outputs []TxOut

	LockTime uint32
}

type UnsupportedTxVersionError uint32

func (e UnsupportedTxVersionError) Error() string {
	return "Unsupported Tx Version" + string(e)
}

func (this *Tx) WriteTo(w io.Writer) (int64, error) {
	if this.Version > 2 {
		return 0, UnsupportedTxVersionError(this.Version)
	}

	var totalWritten int64

	var bufarray [varint.MaxBufferSize]byte
	buf := bufarray[:]

	binary.LittleEndian.PutUint32(buf, this.Version)
	i, err := w.Write(buf[:4])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	i = varint.PutVarInt(buf, uint64(len(this.Inputs)))
	i, err = w.Write(buf[:i])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	for _, txin := range this.Inputs {
		i, err := txin.WriteTo(w)
		totalWritten += i
		if err != nil {
			return totalWritten, err
		}
	}

	i = varint.PutVarInt(buf, uint64(len(this.Outputs)))
	i, err = w.Write(buf[:i])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	for _, txout := range this.Outputs {
		i, err := txout.WriteTo(w)
		totalWritten += i
		if err != nil {
			return totalWritten, err
		}
	}

	binary.LittleEndian.PutUint32(buf, this.LockTime)
	i, err = w.Write(buf[:4])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	return totalWritten, nil
}

func (this *Tx) ReadFrom(r io.Reader) (int64, error) {

	var totalRead int64

	var bufarray [varint.MaxBufferSize]byte
	buf := bufarray[:]

	i, err := r.Read(buf[0:4])
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}
	this.Version = binary.LittleEndian.Uint32(buf)

	if this.Version > 2 {
		return 0, UnsupportedTxVersionError(this.Version)
	}

	inputCount, i, err := varint.ReadVarInt(r)
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	this.Inputs = make([]TxIn, inputCount)

	for txInNum := uint64(0); txInNum < inputCount; txInNum++ {
		i, err := this.Inputs[txInNum].ReadFrom(r)
		totalRead += i
		if err != nil {
			return totalRead, err
		}
	}

	outputCount, i, err := varint.ReadVarInt(r)
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	this.Outputs = make([]TxOut, outputCount)

	for txOutNum := uint64(0); txOutNum < outputCount; txOutNum++ {
		i, err := this.Outputs[txOutNum].ReadFrom(r)
		totalRead += i
		if err != nil {
			return totalRead, err
		}
	}

	i, err = r.Read(buf[0:4])
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}
	this.LockTime = binary.LittleEndian.Uint32(buf)

	return totalRead, nil
}

func (this *Tx) Bytes() ([]byte, error) {

	var buf bytes.Buffer
	_, err := this.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *Tx) Set(data []byte) error {

	buf := bytes.NewReader(data)
	_, err := this.ReadFrom(buf)
	return err
}

func (this *Tx) Hash() (TxHash, error) {

	var txHash TxHash

	h := sha256d.New()
	txData, err := this.Bytes()
	if err != nil {
		return txHash, err
	}
	_, err = h.Write(txData)
	if err != nil {
		return txHash, err
	}

	h.Sum(txHash[:0])

	return txHash, nil
}
