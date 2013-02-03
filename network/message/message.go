package message

import (
	"bytes"
	"encoding/binary"
	"github.com/spearson78/guardian/crypto/sha256d"
	"io"
)

type Message struct {
	Magic   uint32
	Command string
	Payload []byte
}

//Maximum Payload of 2MB
var MaxPayloadSize = 2 * 1024 * 1024
var MaxMessageSize = MaxPayloadSize + 4 + 12 + 4 + 4

var MagicMainNet = uint32(0xD9B4BEF9)
var MagicTestNet = uint32(0x0709110B)

type ChecksumError uint32

func (e ChecksumError) Error() string {
	return "ChecksumError " + string(e)
}

type PayloadOverflowError uint32

func (e PayloadOverflowError) Error() string {
	return "PayloadOverflowError " + string(e)
}

func (this *Message) Checksum() uint32 {
	h := sha256d.New()
	_, err := h.Write(this.Payload)
	if err != nil {
		panic(err)
	}
	hash := h.Sum(nil)

	var checksum uint32

	checksum = binary.LittleEndian.Uint32(hash)

	return checksum
}

func (this *Message) WriteTo(w io.Writer) (int64, error) {

	var bufarray [4]byte
	buf := bufarray[:]
	totalWritten := int64(0)

	binary.LittleEndian.PutUint32(buf, this.Magic)
	i, err := w.Write(buf)
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	var command [12]byte
	copy(command[:], []byte(this.Command))
	i, err = w.Write(command[:])
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	binary.LittleEndian.PutUint32(buf, uint32(len(this.Payload)))
	i, err = w.Write(buf)
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	binary.LittleEndian.PutUint32(buf, this.Checksum())
	i, err = w.Write(buf)
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	i, err = w.Write(this.Payload)
	totalWritten += int64(i)
	if err != nil {
		return totalWritten, err
	}

	return totalWritten, err
}

func (this *Message) ReadFrom(r io.Reader) (int64, error) {

	totalRead := int64(0)

	var bufarray [4]byte
	buf := bufarray[:]

	i, err := io.ReadFull(r, buf)
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	this.Magic = binary.LittleEndian.Uint32(buf)

	var command [12]byte
	i, err = io.ReadFull(r, command[:])
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	endOfCommand := bytes.IndexByte(command[:], 0)
	this.Command = string(command[:endOfCommand])

	i, err = io.ReadFull(r, buf)
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	payloadLength := binary.LittleEndian.Uint32(buf)

	if payloadLength > uint32(MaxPayloadSize) {
		return totalRead, PayloadOverflowError(payloadLength)
	}

	i, err = io.ReadFull(r, buf)
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	checksum := binary.LittleEndian.Uint32(buf)

	this.Payload = make([]byte, payloadLength)
	i, err = io.ReadFull(r, this.Payload)
	totalRead += int64(i)
	if err != nil {
		return totalRead, err
	}

	if this.Checksum() != checksum {
		return totalRead, ChecksumError(checksum)
	}

	return totalRead, nil

}
