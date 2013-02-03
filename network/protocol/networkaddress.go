package protocol

import (
	"github.com/spearson78/guardian/dataio"
	"io"
	"net"
	"time"
)

type VersionNetworkAddress struct {
	Services uint64
	Address  net.TCPAddr
}

type NetworkAddress struct {
	Time time.Time
	VersionNetworkAddress
}

func NewNetworkAddress(address *net.TCPAddr) *NetworkAddress {

	return &NetworkAddress{
		Time: time.Now(),
		VersionNetworkAddress: VersionNetworkAddress{
			Services: 0,
			Address:  *address,
		},
	}
}

func NewVersionNetworkAddress(address *net.TCPAddr) *VersionNetworkAddress {

	return &VersionNetworkAddress{
		Services: 0,
		Address:  *address,
	}
}

func (this *VersionNetworkAddress) WriteTo(w io.Writer) (int64, error) {

	var dw dataio.DataWriter
	dw.Init(w)

	err := dw.WriteUint64(this.Services)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteTCPAddr(&this.Address)
	if err != nil {
		return dw.Count(), err
	}

	return dw.Count(), nil
}

func (this *NetworkAddress) WriteTo(w io.Writer) (int64, error) {

	var dw dataio.DataWriter
	dw.Init(w)

	err := dw.WriteTime32(this.Time)
	if err != nil {
		return dw.Count(), err
	}

	err = dw.WriteWriterTo(&this.VersionNetworkAddress)
	if err != nil {
		return dw.Count(), err
	}

	return dw.Count(), nil
}

func (this *NetworkAddress) ReadFrom(r io.Reader) (int64, error) {

	var err error
	var dr dataio.DataReader
	dr.Init(r)

	this.Time, err = dr.ReadTime32()
	if err != nil {
		return dr.Count(), err
	}

	err = dr.ReadReaderFrom(&this.VersionNetworkAddress)
	return dr.Count(), err
}

func (this *VersionNetworkAddress) ReadFrom(r io.Reader) (int64, error) {

	var err error
	var dr dataio.DataReader
	dr.Init(r)

	this.Services, err = dr.ReadUint64()
	if err != nil {
		return dr.Count(), err
	}

	this.Address, err = dr.ReadTCPAddr()
	return dr.Count(), err
}

func (this *VersionNetworkAddress) Equal(other *VersionNetworkAddress) bool {

	if other == nil {
		return false
	}

	return this.Services == other.Services &&
		this.Address.IP.Equal(other.Address.IP) &&
		this.Address.Port == other.Address.Port
}

func (this *NetworkAddress) Equal(other *NetworkAddress) bool {

	if other == nil {
		return false
	}

	return this.Time.Equal(other.Time) &&
		this.VersionNetworkAddress.Equal(&other.VersionNetworkAddress)
}
