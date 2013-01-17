package scanner

import (
	"bytes"
	"encoding/hex"
	"github.com/spearson78/guardian/script/opcode"
	"github.com/spearson78/guardian/script/token"
	"testing"
)

func TestStandardTransactionToBitcoinAddress(t *testing.T) {
	script, _ := hex.DecodeString("76A91489ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA88AC")
	checkData, _ := hex.DecodeString("89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA")

	var s Scanner

	errorReported := false

	s.Init(script, func(pos int, msg string) {
		errorReported = true
	})

	tok := s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.DUP {
		t.Errorf("Failed DUP tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.HASH160 {
		t.Errorf("Failed HASH160 tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.DATA || !bytes.Equal(s.Data(), checkData) {
		t.Errorf("Failed DUP tok %s op %s bcPos %d endbcPos %d", tok, s.Op(), s.Pos(), s.EndPos())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.EQUALVERIFY {
		t.Errorf("Failed EQUALVERIFY tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.CHECKSIG {
		t.Errorf("Failed CHECKSIG tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed ENDOFSCRIPT tok %s", tok)
	}

	tok = s.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if s.ErrorCount() != 0 {
		t.Errorf("Failed ErrorCount() %d", s.ErrorCount())
	}

}

func TestUnderflow(t *testing.T) {
	script, _ := hex.DecodeString("76A91489ABCDEFABBAABBAABBAABBAABBAABBAABBAAB")
	checkData, _ := hex.DecodeString("89ABCDEFABBAABBAABBAABBAABBAABBAABBAAB")

	var s Scanner

	errorReported := false

	s.Init(script, func(pos int, msg string) {
		errorReported = true
	})

	tok := s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.DUP {
		t.Errorf("Failed DUP tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.HASH160 {
		t.Errorf("Failed HASH160 tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.DATA || !bytes.Equal(s.Data(), checkData) {
		t.Errorf("Failed DUP tok %s op %s bcPos %d endbcPos %d", tok, s.Op(), s.Pos(), s.EndPos())
	}

	tok = s.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed ENDOFSCRIPT tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if !errorReported {
		t.Errorf("Failed Not ErrorReported")
	}

	if s.ErrorCount() != 1 {
		t.Errorf("Failed ErrorCount() %d", s.ErrorCount())
	}

}

func TestInvalidToken(t *testing.T) {
	script, _ := hex.DecodeString("76A91489ABCDEFABBAABBAABBAABBAABBAABBAABBAABBAFFAC")
	checkData, _ := hex.DecodeString("89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA")

	var s Scanner
	errorReported := false

	s.Init(script, func(pos int, msg string) {
		errorReported = true
	})

	tok := s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.DUP {
		t.Errorf("Failed DUP tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.HASH160 {
		t.Errorf("Failed HASH160 tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.DATA || !bytes.Equal(s.Data(), checkData) {
		t.Errorf("Failed DUP tok %s op %s bcPos %d endbcPos %d", tok, s.Op(), s.Pos(), s.EndPos())
	}

	tok = s.Scan()
	if tok != token.INVALID {
		t.Errorf("Failed INVALID tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.CHECKSIG {
		t.Errorf("Failed CHECKSIG tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed ENDOFSCRIPT tok %s", tok)
	}

	tok = s.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if !errorReported {
		t.Errorf("Failed Not ErrorReported")
	}

	if s.ErrorCount() != 1 {
		t.Errorf("Failed ErrorCount() %d", s.ErrorCount())
	}

}

func TestInvalidScript(t *testing.T) {
	script := []byte("script")

	var s Scanner
	errorReported := false

	s.Init(script, func(pos int, msg string) {
		errorReported = true
	})

	tok := s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.IFDUP {
		t.Errorf("Failed IFDUP tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.IF {
		t.Errorf("Failed IF tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.TWOSWAP {
		t.Errorf("Failed TWOSWAP tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.VERIFY {
		t.Errorf("Failed VERIFY tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.TWOOVER {
		t.Errorf("Failed TWOOVER tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.OPERATION || s.Op() != opcode.DEPTH {
		t.Errorf("Failed DEPTH tok %s op %s", tok, s.Op())
	}

	tok = s.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed ENDOFSCRIPT tok %s", tok)
	}

	tok = s.Scan()
	if tok != token.ENDOFSCRIPT {
		t.Errorf("Failed double ENDOFSCRIPT tok %s", tok)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if s.ErrorCount() != 0 {
		t.Errorf("Failed ErrorCount() %d", s.ErrorCount())
	}
}

func BenchmarkStandardTransactionToBitcoinAddress(b *testing.B) {
	script, _ := hex.DecodeString("76A91489ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA88AC")

	var s Scanner

	for i := 0; i < b.N; i++ {
		s.Init(script, nil)
		s.Scan()
		s.Scan()
		s.Scan()
		s.Scan()
		s.Scan()
		s.Scan()
	}
}
