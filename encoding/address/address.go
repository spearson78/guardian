package address

import (
	"bytes"
	_ "code.google.com/p/go.crypto/ripemd160"
	"crypto"
	_ "crypto/sha256"
	"github.com/spearson78/guardian/crypto/sha256d"
	"github.com/spearson78/guardian/encoding/base58"
)

type AddressChecksumError string

func (e AddressChecksumError) Error() string {
	return "Checksum error in address " + string(e)
}

type AddressUnderunError string

func (e AddressUnderunError) Error() string {
	return "Address is too short  " + string(e)
}

func Decode(address string) (pkHash []byte, version byte, err error) {

	binaryAddress, err := base58.BitcoinEncoding.DecodeFromString(address)
	if err != nil {
		return nil, 0, err
	}

	if len(binaryAddress) != 25 {
		return nil, 0, AddressUnderunError(address)
	}

	extended := binaryAddress[0:21]
	version = binaryAddress[0]
	pkHash = binaryAddress[1:21]
	extractedChecksum := binaryAddress[21:]

	s256d := sha256d.New()
	s256d.Write(extended)
	res := s256d.Sum(nil)

	calculatedChecksum := res[0:4]

	if !bytes.Equal(calculatedChecksum, extractedChecksum) {
		return nil, 0, AddressChecksumError(address)
	}

	return pkHash, version, nil
}

func EncodePublicKeyHash(publicKeyHash []byte, version byte) string {
	var err error

	extended := make([]byte, 1, len(publicKeyHash)+5)
	extended[0] = version
	extended = append(extended, publicKeyHash...)

	s256d := sha256d.New()
	s256d.Write(extended)
	res := s256d.Sum(nil)

	checksum := res[0:4]

	binaryAddress := append(extended, checksum...)
	address, err := base58.BitcoinEncoding.EncodeToString(binaryAddress)
	if err != nil {
		panic(err)
	}

	return address

}

func HashPublicKey(publicKey []byte) []byte {
	s256 := crypto.SHA256.New()
	r160 := crypto.RIPEMD160.New()

	s256.Write(publicKey)
	res := s256.Sum(nil)

	r160.Write(res)
	res = r160.Sum(nil)

	return res
}

func EncodePublicKey(publicKey []byte, version byte) string {
	pkHash := HashPublicKey(publicKey)
	return EncodePublicKeyHash(pkHash, version)
}
