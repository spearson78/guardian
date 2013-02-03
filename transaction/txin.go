package transaction

import (
	"github.com/spearson78/guardian/dataio"
	"io"
)

type TxIn struct {
	Previous OutPoint
	Script   []byte
	Sequence uint32
}

func (this *TxIn) WriteTo(w io.Writer) (int64, error) {
	var dw dataio.DataWriter
	dw.Init(w)

	err := dw.Write(this.Previous.Hash[:])
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteUint32(this.Previous.Index)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteVarBytes(this.Script)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteUint32(this.Sequence)
	if err != nil {
		return dw.Count(), err
	}

	return dw.Count(), nil
}

func (this *TxIn) ReadFrom(r io.Reader) (int64, error) {

	var err error
	var dr dataio.DataReader
	dr.Init(r)

	err = dr.ReadFull(this.Previous.Hash[:])
	if err != nil {
		return dr.Count(), err
	}

	this.Previous.Index, err = dr.ReadUint32()
	if err != nil {
		return dr.Count(), err
	}

	this.Script, err = dr.ReadVarBytes()
	if err != nil {
		return dr.Count(), err
	}

	this.Sequence, err = dr.ReadUint32()
	if err != nil {
		return dr.Count(), err
	}

	return dr.Count(), nil
}
