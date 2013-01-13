package compiler

import (
	"bytes"
	"encoding/hex"
	"github.com/spearson78/guardian/script/lexer"
	"github.com/spearson78/guardian/script/scanner"
	"strings"
	"testing"
)

func TestStandardTransactionToBitcoinAddressFromLexer(t *testing.T) {
	script := `
DUP
HASH160
0x89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA
EQUALVERIFY
CHECKSIG
`
	result, _ := hex.DecodeString("76A91489ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA88AC")

	l := new(lexer.Lexer)
	l.Init(strings.NewReader(script), nil)

	compiled, _ := Compile(l)

	if !bytes.Equal(compiled, result) {
		t.Errorf("Failed")
	}

}

func TestStandardTransactionToBitcoinAddressFromScanner(t *testing.T) {
	script, _ := hex.DecodeString("76A91489ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA88AC")

	s := new(scanner.Scanner)
	s.Init(script, nil)

	compiled, _ := Compile(s)

	if !bytes.Equal(compiled, script) {
		t.Errorf("Failed")
	}

}

func TestNumbers(t *testing.T) {
	script := `
-14011978
-2
-1
0
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
14011978
`
	result, _ := hex.DecodeString("044aced58001824f005152535455565758595a5b5c5d5e5f600111044aced500")

	l := new(lexer.Lexer)
	l.Init(strings.NewReader(script), nil)

	compiled, err := Compile(l)
	if err != nil {
		t.Errorf("Failed %v", err)
	}

	if !bytes.Equal(compiled, result) {
		t.Errorf("Failed compiled %s expected %s", hex.EncodeToString(compiled), hex.EncodeToString(result))
	}

}
