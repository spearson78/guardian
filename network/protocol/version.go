package protocol

import (
	"github.com/spearson78/guardian/dataio"
	"io"
	"time"
)

type Version struct {
	Version     int32
	Services    uint64
	Time        time.Time
	Receive     VersionNetworkAddress
	From        VersionNetworkAddress
	Nonce       uint64
	UserAgent   string
	StartHeight int32
}

func (this *Version) ReadFrom(r io.Reader) (int64, error) {

	var err error
	var dr dataio.DataReader
	dr.Init(r)

	this.Version, err = dr.ReadInt32()
	if err != nil {
		return dr.Count(), err
	}

	this.Services, err = dr.ReadUint64()
	if err != nil {
		return dr.Count(), err
	}

	this.Time, err = dr.ReadTime64()
	if err != nil {
		return dr.Count(), err
	}

	err = dr.ReadReaderFrom(&this.Receive)
	if err != nil {
		return dr.Count(), err
	}

	err = dr.ReadReaderFrom(&this.From)
	if err != nil {
		return dr.Count(), err
	}

	this.Nonce, err = dr.ReadUint64()
	if err != nil {
		return dr.Count(), err
	}

	this.UserAgent, err = dr.ReadVarString()
	if err != nil {
		return dr.Count(), err
	}

	this.StartHeight, err = dr.ReadInt32()
	if err != nil {
		return dr.Count(), err
	}

	return dr.Count(), nil
}

func (this *Version) WriteTo(w io.Writer) (int64, error) {

	var dw dataio.DataWriter
	dw.Init(w)

	err := dw.WriteInt32(this.Version)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteUint64(this.Services)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteTime64(this.Time)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteWriterTo(&this.Receive)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteWriterTo(&this.From)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteUint64(this.Nonce)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteVarString(this.UserAgent)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteInt32(this.StartHeight)
	if err != nil {
		return dw.Count(), err
	}

	return dw.Count(), nil

}
