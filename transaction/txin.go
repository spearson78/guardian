package transaction

import (
	"encoding/binary"
	"github.com/spearson78/guardian/encoding/varint"
	"io"
)

type TxIn struct {
	Previous OutPoint
	Script   []byte
	Sequence uint32
}

func (this *TxIn) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	var bufarray [varint.MaxBufferSize]byte
	buf := bufarray[:]

	i, err := w.Write(this.Previous.Hash[:])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	binary.LittleEndian.PutUint32(buf, this.Previous.Index)
	i, err = w.Write(buf[:4])
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

	binary.LittleEndian.PutUint32(buf, this.Sequence)
	i, err = w.Write(buf[:4])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	return totalWritten, nil
}

func (this *TxIn) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	var bufarray [varint.MaxBufferSize]byte
	buf := bufarray[:]

	i, err := r.Read(this.Previous.Hash[:])
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	i, err = r.Read(buf[:4])
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}
	this.Previous.Index = binary.LittleEndian.Uint32(buf)

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

	i, err = r.Read(buf[:4])
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}
	this.Sequence = binary.LittleEndian.Uint32(buf)

	return totalRead, nil
}
