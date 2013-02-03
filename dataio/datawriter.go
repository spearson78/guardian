package dataio

import (
	"encoding/binary"
	"github.com/spearson78/guardian/encoding/tcpaddr"
	"github.com/spearson78/guardian/encoding/varint"
	"github.com/spearson78/guardian/encoding/varstring"
	"io"
	"net"
	"time"
)

type DataWriter struct {
	w io.Writer

	buf          [18]byte
	slice2       []byte
	slice4       []byte
	slice8       []byte
	sliceTCPAddr []byte
	sliceVarInt  []byte

	count int64
}

func (this *DataWriter) Init(w io.Writer) {
	this.w = w
	this.slice2 = this.buf[:2]
	this.slice4 = this.buf[:4]
	this.slice8 = this.buf[:8]
	this.sliceTCPAddr = this.buf[:tcpaddr.Size]
	this.sliceVarInt = this.buf[:varint.MaxBufferSize]
}

func (this *DataWriter) Count() int64 {
	return this.count
}

func (this *DataWriter) WriteInt32(value int32) error {
	return this.WriteUint32(uint32(value))
}

func (this *DataWriter) WriteUint32(value uint32) error {

	binary.LittleEndian.PutUint32(this.slice4, value)
	i, err := this.w.Write(this.slice4)
	this.count += int64(i)
	return err
}

func (this *DataWriter) WriteInt64(value int64) error {
	return this.WriteUint64(uint64(value))
}

func (this *DataWriter) WriteUint64(value uint64) error {

	binary.LittleEndian.PutUint64(this.slice8, value)
	i, err := this.w.Write(this.slice8)
	this.count += int64(i)
	return err
}

func (this *DataWriter) WriteTime32(value time.Time) error {
	return this.WriteUint32(uint32(value.Unix()))
}

func (this *DataWriter) WriteTime64(value time.Time) error {
	return this.WriteUint64(uint64(value.Unix()))
}

func (this *DataWriter) WriteWriterTo(value io.WriterTo) error {
	i, err := value.WriteTo(this.w)
	this.count += int64(i)
	return err
}

func (this *DataWriter) WriteArray(value []io.WriterTo) error {

	err := this.WriteVarInt(uint64(len(value)))
	if err != nil {
		return err
	}

	for i, _ := range value {
		err := this.WriteWriterTo(value[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *DataWriter) WriteVarInt(value uint64) error {
	i := varint.PutVarInt(this.sliceVarInt, value)
	i, err := this.w.Write(this.buf[:i])
	this.count += int64(i)
	return err
}

func (this *DataWriter) WriteVarString(value string) error {
	i, err := varstring.WriteVarString(this.w, value)
	this.count += int64(i)
	return err
}

func (this *DataWriter) WriteTCPAddr(value *net.TCPAddr) error {
	i := tcpaddr.PutTCPAddr(this.sliceTCPAddr, value)
	i, err := this.w.Write(this.sliceTCPAddr)
	this.count += int64(i)
	return err
}

func (this *DataWriter) Write(data []byte) error {
	i, err := this.w.Write(data)
	this.count += int64(i)
	return err
}

func (this *DataWriter) WriteVarBytes(data []byte) error {
	err := this.WriteVarInt(uint64(len(data)))
	if err != nil {
		return err
	}

	err = this.Write(data)
	return err
}
