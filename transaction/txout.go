package transaction

import (
	"github.com/spearson78/guardian/dataio"
	"io"
)

type TxOut struct {
	Value  int64
	Script []byte
}

func (this *TxOut) WriteTo(w io.Writer) (int64, error) {
	var dw dataio.DataWriter
	dw.Init(w)

	err := dw.WriteInt64(this.Value)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteVarBytes(this.Script)
	if err != nil {
		return dw.Count(), err
	}

	return dw.Count(), nil
}

func (this *TxOut) ReadFrom(r io.Reader) (int64, error) {
	var err error
	var dr dataio.DataReader
	dr.Init(r)

	this.Value, err = dr.ReadInt64()
	if err != nil {
		return dr.Count(), err
	}

	this.Script, err = dr.ReadVarBytes()
	if err != nil {
		return dr.Count(), err
	}

	return dr.Count(), nil
}
