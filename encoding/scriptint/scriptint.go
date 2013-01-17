package scriptint

import (
	"math/big"
)

func Encode(i *big.Int) []byte {

	sign := i.Sign()
	data := i.Bytes()

	if len(data) == 0 {
		if sign == -1 {
			data = []byte{0x80}
		}
	} else {
		if data[0]&0x80 == 0x80 {
			extended := make([]byte, len(data)+1)
			copy(extended[1:], data)
			data = extended
		}

		if sign == -1 {
			data[0] = data[0] | 0x80
		}

		reverseInPlace(data)
	}

	return data
}

func Decode(b []byte) *big.Int {
	if len(b) == 0 {
		return big.NewInt(0)
	}

	var bytes []byte
	bytes = append(bytes, b...)
	reverseInPlace(bytes)

	neg := false

	result := big.NewInt(0)

	if bytes[0]&0x80 == 0x80 {
		neg = true
	}
	bytes[0] = bytes[0] & 0x7f
	result.SetBytes(bytes)
	if neg {
		result.Neg(result)
	}

	return result
}

func reverseInPlace(in []byte) {

	l := len(in)
	hl := l / 2

	for i := 0; i < hl; i++ {
		j := l - i - 1
		in[i], in[j] = in[j], in[i]
	}
}
