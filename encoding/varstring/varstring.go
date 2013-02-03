package varstring

import (
	"github.com/spearson78/guardian/encoding/varint"
	"io"
)

func WriteVarString(w io.Writer, value string) (int, error) {

	totalWritten := 0

	length := uint64(len(value))
	i, err := varint.WriteVarInt(w, length)
	totalWritten += i
	if err != nil {
		return totalWritten, err
	}

	bytes := []byte(value)
	i, err = w.Write(bytes)
	totalWritten += i
	if err != nil {
		return totalWritten, err
	}

	return totalWritten, nil
}

func ReadVarString(r io.Reader) (string, int, error) {

	totalRead := 0

	length, i, err := varint.ReadVarInt(r)
	totalRead += i
	if err != nil {
		return "", totalRead, err
	}

	bytes := make([]byte, length)
	i, err = io.ReadFull(r, bytes)
	totalRead += i
	if err != nil {
		return err.Error(), totalRead, err
	}

	value := string(bytes)

	return value, totalRead, nil

}
