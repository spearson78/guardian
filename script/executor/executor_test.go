package executor

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/spearson78/guardian/script/compiler"
	"github.com/spearson78/guardian/script/lexer"
	"strings"
	"testing"
)

type MockCheckSig struct {
	Pk        []byte
	HashType  uint32
	Sig       []byte
	SubScript []byte
}

func (this *MockCheckSig) CheckSig(pk []byte, hashType uint32, sig []byte, subScript []byte) error {
	this.Pk = pk
	this.HashType = hashType
	this.Sig = sig
	this.SubScript = subScript

	return nil
}

func TestStandardTransactionToBitcoinAddress(t *testing.T) {

	fmt.Println("TestStandardTransactionToBitcoinAddress")

	//Transaction dd11322ddb4487a0a7a597838df02fbfd2070870721cd648f3b37f0a164efed9
	script := `
0x30440220694ff325724a4f4b0f3f0c36bf8e94cac58ad7c9b4d5bd8c7286c0da623f0b2c02206ae94680a8f31f30cd846da258e919c94afe2dd629b4f4ce11bbe8165ff99a5f01
0x04fc60372d27b067ca306ba812ced9c8cd69296b83a40b9b57c593258c1b9e0ee1c0c621ca558b878395f9645a4b67a96e51843e9c060d43a3833fdd29a91f4f31
DUP
HASH160
0x340cfcffe029e6935f4e4e5839a2ff5f29c7a571
EQUALVERIFY
CHECKSIG
`

	sig, _ := hex.DecodeString("30440220694ff325724a4f4b0f3f0c36bf8e94cac58ad7c9b4d5bd8c7286c0da623f0b2c02206ae94680a8f31f30cd846da258e919c94afe2dd629b4f4ce11bbe8165ff99a5f")
	hashType := uint32(1)
	pk, _ := hex.DecodeString("04fc60372d27b067ca306ba812ced9c8cd69296b83a40b9b57c593258c1b9e0ee1c0c621ca558b878395f9645a4b67a96e51843e9c060d43a3833fdd29a91f4f31")
	subscript := `
0x04fc60372d27b067ca306ba812ced9c8cd69296b83a40b9b57c593258c1b9e0ee1c0c621ca558b878395f9645a4b67a96e51843e9c060d43a3833fdd29a91f4f31
DUP
HASH160
0x340cfcffe029e6935f4e4e5839a2ff5f29c7a571
EQUALVERIFY
CHECKSIG
`

	l := new(lexer.Lexer)
	l.Init(strings.NewReader(script), nil)
	compiled, _ := compiler.Compile(l)

	l.Init(strings.NewReader(subscript), nil)
	compiledSubScript, _ := compiler.Compile(l)

	var checkSig MockCheckSig
	e := new(Executor)
	e.Init(&checkSig)

	err := e.Execute(compiled)
	if err != nil {
		t.Errorf("TestStandardTransactionToBitcoinAddressNoCheckSig Failed %v", err)
	}

	if !bytes.Equal(sig, checkSig.Sig) {
		t.Errorf("TestStandardTransactionToBitcoinAddressNoCheckSig Sig Mismatch %s", hex.EncodeToString(checkSig.Sig))
	}

	if !bytes.Equal(pk, checkSig.Pk) {
		t.Errorf("TestStandardTransactionToBitcoinAddressNoCheckSig Pk Mismatch %s", hex.EncodeToString(checkSig.Pk))
	}

	if hashType != checkSig.HashType {
		t.Errorf("TestStandardTransactionToBitcoinAddressNoCheckSig Pk Mismatch %s", hex.EncodeToString(checkSig.Pk))
	}

	if !bytes.Equal(compiledSubScript, checkSig.SubScript) {
		t.Errorf("TestStandardTransactionToBitcoinAddressNoCheckSig SubScript Mismatch %s", hex.EncodeToString(checkSig.SubScript))
	}
}

func TestIf(t *testing.T) {

	tests := []struct {
		script string
		result byte
	}{
		{
			script: "0x01 IF 0x02 ELSE 0x03 ENDIF",
			result: 0x02,
		},
		{
			script: "0x00 IF 0x02 ELSE 0x03 ENDIF",
			result: 0x03,
		},
		{
			script: "0x00 IF 0x02 ELSE 0x03 ENDIF",
			result: 0x03,
		},
		{
			script: "0x01 IF 0x02 IF 0x03 ENDIF ELSE 0x04 ENDIF",
			result: 0x03,
		},
		{
			script: "0x01 IF 0x02 IF 0x03 ENDIF 0x04 ELSE 0x05 ENDIF",
			result: 0x04,
		},
		{
			script: "0x01 NOTIF 0x02 IF 0x03 ENDIF 0x04 ELSE 0x05 ENDIF",
			result: 0x05,
		},
		{
			script: "0x01 NOTIF 0x02 IF 0x03 ENDIF 0x04 ELSE 0x05 IF 0x06 ELSE 0x07 ENDIF ENDIF",
			result: 0x06,
		},
	}

	for _, test := range tests {
		l := new(lexer.Lexer)
		l.Init(strings.NewReader(test.script), nil)
		compiled, _ := compiler.Compile(l)

		e := new(Executor)
		e.Init(nil)

		e.Execute(compiled)

		top := e.Pop()
		if top[0] != test.result {
			t.Errorf("Wrong Top %d", top[0])
		}
	}

}

func BenchmarkNops(b *testing.B) {
	script := `
NOP
NOP
NOP
NOP
NOP
NOP
NOP
NOP
NOP
NOP
`

	l := new(lexer.Lexer)
	l.Init(strings.NewReader(script), nil)
	compiled, _ := compiler.Compile(l)

	e := new(Executor)
	e.Init(nil)

	for i := 0; i < b.N; i++ {
		e.Execute(compiled)
	}

}

func BenchmarkStandardTransactionToBitcoinAddressNopCheckSig(b *testing.B) {

	//Transaction dd11322ddb4487a0a7a597838df02fbfd2070870721cd648f3b37f0a164efed9
	script := `
0x30440220694ff325724a4f4b0f3f0c36bf8e94cac58ad7c9b4d5bd8c7286c0da623f0b2c02206ae94680a8f31f30cd846da258e919c94afe2dd629b4f4ce11bbe8165ff99a5f01
0x04fc60372d27b067ca306ba812ced9c8cd69296b83a40b9b57c593258c1b9e0ee1c0c621ca558b878395f9645a4b67a96e51843e9c060d43a3833fdd29a91f4f31
DUP
HASH160
0x340cfcffe029e6935f4e4e5839a2ff5f29c7a571
EQUALVERIFY
CHECKSIG
`

	l := new(lexer.Lexer)
	l.Init(strings.NewReader(script), nil)
	compiled, _ := compiler.Compile(l)

	e := new(Executor)
	e.Init(nil)

	for i := 0; i < b.N; i++ {
		e.Execute(compiled)
	}

}
