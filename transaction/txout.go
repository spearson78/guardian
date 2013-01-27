package transaction

import (
	"encoding/binary"
	"github.com/spearson78/guardian/encoding/varint"
	"io"
)

type TxOut struct {
	Value  int64
	Script []byte
}

func (this *TxOut) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	var bufarray [varint.MaxBufferSize]byte
	buf := bufarray[:]

	binary.LittleEndian.PutUint64(buf, uint64(this.Value))
	i, err := w.Write(buf[:8])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	i = varint.PutVarInt(buf, uint64(len(this.Script)))
	i, err = w.Write(buf[:i])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	i, err = w.Write(this.Script)
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	return totalWritten, nil
}

func (this *TxOut) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	var bufarray [varint.MaxBufferSize]byte
	buf := bufarray[:]

	i, err := r.Read(buf[:8])
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}
	this.Value = int64(binary.LittleEndian.Uint64(buf))

	scriptLength, i, err := varint.ReadVarInt(r)
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	this.Script = make([]byte, scriptLength)

	i, err = r.Read(this.Script)
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	return totalRead, nil
}
