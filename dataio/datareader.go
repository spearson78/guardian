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

type DataReader struct {
	r io.Reader

	buf          [32]byte
	slice2       []byte
	slice4       []byte
	slice8       []byte
	sliceTCPAddr []byte
	sliceVarInt  []byte

	count int64
}

func (this *DataReader) Init(r io.Reader) {
	this.r = r
	this.slice2 = this.buf[:2]
	this.slice4 = this.buf[:4]
	this.slice8 = this.buf[:8]
	this.sliceTCPAddr = this.buf[:tcpaddr.Size]
	this.sliceVarInt = this.buf[:varint.MaxBufferSize]

}

func (this *DataReader) Count() int64 {
	return this.count
}

func (this *DataReader) ReadInt32() (int32, error) {
	value, err := this.ReadUint32()
	return int32(value), err
}

func (this *DataReader) ReadUint32() (uint32, error) {

	i, err := io.ReadFull(this.r, this.slice4)
	this.count += int64(i)
	value := binary.LittleEndian.Uint32(this.slice4)
	return value, err
}

func (this *DataReader) ReadInt64() (int64, error) {
	value, err := this.ReadUint64()
	return int64(value), err
}

func (this *DataReader) ReadUint64() (uint64, error) {

	i, err := io.ReadFull(this.r, this.slice8)
	this.count += int64(i)
	value := binary.LittleEndian.Uint64(this.slice8)
	return value, err
}

func (this *DataReader) ReadTime32() (time.Time, error) {
	value, err := this.ReadUint32()
	return time.Unix(int64(value), 0), err
}

func (this *DataReader) ReadTime64() (time.Time, error) {
	value, err := this.ReadUint64()
	return time.Unix(int64(value), 0), err
}

func (this *DataReader) ReadReaderFrom(value io.ReaderFrom) error {
	i, err := value.ReadFrom(this.r)
	this.count += int64(i)
	return err
}

func (this *DataReader) ReadVarInt() (uint64, error) {

	value, i, err := varint.ReadVarInt(this.r)
	this.count += int64(i)
	return value, err
}

func (this *DataReader) ReadVarString() (string, error) {
	value, i, err := varstring.ReadVarString(this.r)
	this.count += int64(i)
	return value, err
}

func (this *DataReader) ReadTCPAddr() (net.TCPAddr, error) {
	i, err := io.ReadFull(this.r, this.sliceTCPAddr)
	this.count += int64(i)
	value, _ := tcpaddr.TCPAddr(this.sliceTCPAddr)
	return value, err
}

func (this *DataReader) ReadFull(buf []byte) error {
	i, err := io.ReadFull(this.r, buf)
	this.count += int64(i)
	return err
}

func (this *DataReader) ReadVarBytes() ([]byte, error) {
	length, err := this.ReadVarInt()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length)

	err = this.ReadFull(buf)

	return buf, err
}
