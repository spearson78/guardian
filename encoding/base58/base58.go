//Package base58 implements base58 encoding.and decoding

package base58

import (
	"math/big"
)

//An Encoding is a radix 58 encoding/decoding scheme, defined by a 58 character alphabet.
type Encoding struct {
	encode    string
	decodeMap [256]byte
}

const encodeBitcoin = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
const encodeFlickr = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

// NewEncoding returns a new Encoding defined by the given alphabet, which must be a 58-byte string.
func NewEncoding(encoder string) *Encoding {
	e := new(Encoding)
	e.encode = encoder

	for i := 0; i < len(e.decodeMap); i++ {
		e.decodeMap[i] = 0xFF
	}
	for i := 0; i < len(encoder); i++ {
		e.decodeMap[encoder[i]] = byte(i)
	}

	return e
}

//BitcoinEncoding is the base58 encoding used in the Bitcoin protocol
var BitcoinEncoding = NewEncoding(encodeBitcoin)

//FlickrEncoding is the base58 encoding used in Flickr short urls.
var FlickrEncoding = NewEncoding(encodeFlickr)

type CorruptEncodeInputError big.Int

func (e CorruptEncodeInputError) Error() string {
	return "Illegal base58 big int " + (*big.Int)(&e).String()
}

type CorruptDecodeInputError byte

func (e CorruptDecodeInputError) Error() string {
	return "Illegal base58 byte " + string(e)
}

var c58 = big.NewInt(58)
var c0 = big.NewInt(0)

//Encode encodes src using the encoding enc.
//Leading 0 bytes in the src are encoded as 0 bytes in the encoded alphabet.
//The ensure that leading 0 bytes are not lost during encoding.
func (enc *Encoding) Encode(src []byte) ([]byte, error) {
	i := new(big.Int).SetBytes(src)
	result, err := enc.encodeInternal(i)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(src) && src[i] == 0; i++ {
		result = append(result, enc.encode[0])
	}
	reverseInPlace(result)

	return result, nil

}

//Encode encodes src using the encoding enc.
//Leading 0 bytes in the src are encoded as 0 bytes in the encoded alphabet.
//The ensure that leading 0 bytes are not lost during encoding.
func (enc *Encoding) EncodeToString(src []byte) (string, error) {
	buf, err := enc.Encode(src)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func reverseInPlace(in []byte) {

	l := len(in)
	hl := l / 2

	for i := 0; i < hl; i++ {
		j := l - i - 1
		in[i], in[j] = in[j], in[i]
	}
}

func (enc *Encoding) encodeInternal(src *big.Int) ([]byte, error) {
	if src.Sign() == 0 {
		return make([]byte, 0), nil
	}

	if src.Sign() < 0 {
		return nil, CorruptEncodeInputError(*src)
	}

	val := new(big.Int).Set(src)

	result := make([]byte, 0, 32)

	var mod58 big.Int
	for val.Sign() > 0 {
		mod58.Mod(val, c58)
		val.Div(val, c58)
		result = append(result, enc.encode[mod58.Int64()])
	}

	return result, nil
}

//Encode encodes src using the encoding enc.
//src must be positive
func (enc *Encoding) EncodeInt(src *big.Int) ([]byte, error) {

	result, err := enc.encodeInternal(src)
	if err != nil {
		return nil, err
	}

	reverseInPlace(result)

	return result, nil
}

func (enc *Encoding) DecodeFromString(src string) ([]byte, error) {
	return enc.Decode([]byte(src))
}

func (enc *Encoding) Decode(src []byte) ([]byte, error) {

	if len(src) == 0 {
		return make([]byte, 0), nil
	}

	result, err := enc.DecodeToInt(src)
	if err != nil {
		return nil, err
	}

	resultBytes := result.Bytes()

	prefixNullCount := 0
	encodedNull := enc.encode[0]
	for prefixNullCount < len(src) && src[prefixNullCount] == encodedNull {
		prefixNullCount++
	}

	prefixedBytes := make([]byte, prefixNullCount, len(resultBytes)+prefixNullCount)

	prefixedBytes = append(prefixedBytes, resultBytes...)

	return prefixedBytes, nil
}

func (enc *Encoding) DecodeToInt(src []byte) (*big.Int, error) {
	result := new(big.Int)
	a := new(big.Int)
	for _, b := range src {

		result.Mul(result, c58)
		decoded := enc.decodeMap[b]
		if decoded == 0xFF {
			return nil, CorruptDecodeInputError(b)
		}
		a.SetInt64(int64(decoded))
		result.Add(result, a)
	}

	return result, nil
}
