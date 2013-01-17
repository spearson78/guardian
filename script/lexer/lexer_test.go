package lexer

import (
	"bytes"
	"encoding/hex"
	"github.com/spearson78/guardian/script/opcode"
	"github.com/spearson78/guardian/script/token"
	"math/big"
	"strings"
	"testing"
)

func TestComments(t *testing.T) {
	script := `
//Duplicate the public key
DUP
//Hash the public key
HASH160
//Push the required public key hash to the stack
0x89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA
//Verify the 2 public key hashes match
EQUALVERIFY
//Check the transaction signature
CHECKSIG
`

	checkData, _ := hex.DecodeString("89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA")

	var l Lexer

	errorReported := false

	l.Init(strings.NewReader(script), func(pos int, msg string) {
		errorReported = true
	})

	tok := l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.DUP {
		t.Errorf("Failed DUP tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.HASH160 {
		t.Errorf("Failed HASH160 tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.DATA || !bytes.Equal(l.Data(), checkData) {
		t.Errorf("Failed DUP tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.EQUALVERIFY {
		t.Errorf("Failed EQUALVERIFY tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.CHECKSIG {
		t.Errorf("Failed CHECKSIG tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed ENDOFSCRIPT tok %s", tok)
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if l.ErrorCount() != 0 {
		t.Errorf("Failed ErrorCount %d", l.ErrorCount())
	}
}

func TestStandardTransactionToBitcoinAddress(t *testing.T) {
	script := `
DUP
HASH160
0x89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA
EQUALVERIFY
CHECKSIG
`

	checkData, _ := hex.DecodeString("89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA")

	var l Lexer

	errorReported := false

	l.Init(strings.NewReader(script), func(pos int, msg string) {
		errorReported = true
	})

	tok := l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.DUP {
		t.Errorf("Failed DUP tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.HASH160 {
		t.Errorf("Failed HASH160 tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.DATA || !bytes.Equal(l.Data(), checkData) {
		t.Errorf("Failed DUP tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.EQUALVERIFY {
		t.Errorf("Failed EQUALVERIFY tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.CHECKSIG {
		t.Errorf("Failed CHECKSIG tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed ENDOFSCRIPT tok %s", tok)
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if l.ErrorCount() != 0 {
		t.Errorf("Failed ErrorCount() %d", l.ErrorCount())
	}
}

func TestMinusOne(t *testing.T) {
	script := "-1"

	var l Lexer

	errorReported := false

	l.Init(strings.NewReader(script), func(pos int, msg string) {
		errorReported = true
	})

	tok := l.Scan()
	if tok != token.NUMBER || l.Number().Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Failed -1 tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if l.ErrorCount() != 0 {
		t.Errorf("Failed ErrorCount() %d", l.ErrorCount())
	}
}

func TestMinusOperation(t *testing.T) {
	script := "-DUP"

	var l Lexer

	errorReported := false

	l.Init(strings.NewReader(script), func(pos int, msg string) {
		errorReported = true
	})

	tok := l.Scan()
	if tok != token.INVALID {
		t.Errorf("Failed -DUP tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if l.ErrorCount() != 0 {
		t.Errorf("Failed ErrorCount() %d", l.ErrorCount())
	}
}

func TestMinus(t *testing.T) {
	script := "-"

	var l Lexer

	errorReported := false

	l.Init(strings.NewReader(script), func(pos int, msg string) {
		errorReported = true
	})

	tok := l.Scan()
	if tok != token.INVALID {
		t.Errorf("Failed - tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if l.ErrorCount() != 0 {
		t.Errorf("Failed ErrorCount() %d", l.ErrorCount())
	}
}

func TestInvalidToken(t *testing.T) {
	script := `
DUP
HASH160
0x89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA
BAD
CHECKSIG
`

	checkData, _ := hex.DecodeString("89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA")

	var l Lexer

	errorReported := false

	l.Init(strings.NewReader(script), func(pos int, msg string) {
		errorReported = true
	})

	tok := l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.DUP {
		t.Errorf("Failed DUP tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.HASH160 {
		t.Errorf("Failed HASH160 tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.DATA || !bytes.Equal(l.Data(), checkData) {
		t.Errorf("Failed DUP tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.INVALID {
		t.Errorf("Failed EQUALVERIFY tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.OPERATION || l.Op() != opcode.CHECKSIG {
		t.Errorf("Failed CHECKSIG tok %s op %s", tok, l.Op())
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed ENDOFSCRIPT tok %s", tok)
	}

	tok = l.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if !errorReported {
		t.Errorf("Failed No ErrorReported")
	}

	if l.ErrorCount() != 1 {
		t.Errorf("Failed ErrorCount() %d", l.ErrorCount())
	}
}

func BenchmarkStandardTransactionToBitcoinAddress(b *testing.B) {
	script := `
DUP
HASH160
0x89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA
EQUALVERIFY
CHECKSIG
`
	var l Lexer

	for i := 0; i < b.N; i++ {
		l.Init(strings.NewReader(script), nil)
		l.Scan()
		l.Scan()
		l.Scan()
		l.Scan()
		l.Scan()
		l.Scan()
	}
}
