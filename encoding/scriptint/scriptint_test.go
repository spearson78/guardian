package scriptint

import (
	"math/big"
	"testing"
)

func Test(t *testing.T) {

	//TODO: negative zero is lost
	number := big.NewInt(0)
	encoded := Encode(number)
	decoded := Decode(encoded)

	if decoded.Cmp(number) != 0 {
		t.Errorf("Failed")
	}

}
