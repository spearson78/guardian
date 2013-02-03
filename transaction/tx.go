package transaction

import (
	"bytes"
	"github.com/spearson78/guardian/crypto/sha256d"
	"github.com/spearson78/guardian/dataio"
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

	var dw dataio.DataWriter
	dw.Init(w)

	err := dw.WriteUint32(this.Version)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteVarInt(uint64(len(this.Inputs)))
	if err != nil {
		return dw.Count(), err
	}

	for i, _ := range this.Inputs {
		err := dw.WriteWriterTo(&this.Inputs[i])
		if err != nil {
			return dw.Count(), err
		}
	}

	err = dw.WriteVarInt(uint64(len(this.Outputs)))
	if err != nil {
		return dw.Count(), err
	}

	for i, _ := range this.Outputs {
		err := dw.WriteWriterTo(&this.Outputs[i])
		if err != nil {
			return dw.Count(), err
		}
	}

	err = dw.WriteUint32(this.LockTime)
	if err != nil {
		return dw.Count(), err
	}

	return dw.Count(), nil
}

func (this *Tx) ReadFrom(r io.Reader) (int64, error) {

	var err error
	var dr dataio.DataReader
	dr.Init(r)

	this.Version, err = dr.ReadUint32()
	if err != nil {
		return dr.Count(), err
	}

	if this.Version > 2 {
		return 0, UnsupportedTxVersionError(this.Version)
	}

	inputCount, err := dr.ReadVarInt()
	if err != nil {
		return dr.Count(), err
	}

	this.Inputs = make([]TxIn, inputCount)

	for txInNum := uint64(0); txInNum < inputCount; txInNum++ {
		err = dr.ReadReaderFrom(&this.Inputs[txInNum])
		if err != nil {
			return dr.Count(), err
		}
	}

	outputCount, err := dr.ReadVarInt()
	if err != nil {
		return dr.Count(), err
	}

	this.Outputs = make([]TxOut, outputCount)

	for txOutNum := uint64(0); txOutNum < outputCount; txOutNum++ {
		err = dr.ReadReaderFrom(&this.Outputs[txOutNum])
		if err != nil {
			return dr.Count(), err
		}
	}

	this.LockTime, err = dr.ReadUint32()
	if err != nil {
		return dr.Count(), err
	}

	return dr.Count(), nil
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
	_, err := this.WriteTo(h)
	if err != nil {
		return txHash, err
	}

	h.Sum(txHash[:0])

	return txHash, nil
}
