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

	tok, op, bcPos, endbcPos, data, number := s.Scan()
	if tok != token.OPERATION || op != opcode.DUP || bcPos != 0 || endbcPos != 1 || data != nil || number != nil {
		t.Errorf("Failed DUP tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.HASH160 || bcPos != 1 || endbcPos != 2 || data != nil || number != nil {
		t.Errorf("Failed HASH160 tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.DATA || op != opcode.INVALID || bcPos != 2 || endbcPos != 23 || !bytes.Equal(data, checkData) || number != nil {
		t.Errorf("Failed DATA tok %s op %s bcPos %d endbcPos %d data %s number %s", tok, op, bcPos, endbcPos, hex.EncodeToString(data), number)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.EQUALVERIFY || bcPos != 23 || endbcPos != 24 || data != nil || number != nil {
		t.Errorf("Failed EQUALVERIFY tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.CHECKSIG || bcPos != 24 || endbcPos != 25 || data != nil || number != nil {
		t.Errorf("Failed CHECKSIG tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.ENDOFSCRIPT || op != opcode.INVALID || bcPos != 25 || endbcPos != 26 || data != nil || number != nil {
		t.Errorf("Failed ENDOFSCRIPT tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.ENDOFSCRIPT || op != opcode.INVALID || bcPos != 25 || endbcPos != 26 || data != nil || number != nil {
		t.Errorf("Failed double ENDOFSCRIPT tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if s.ErrorCount != 0 {
		t.Errorf("Failed ErrorCount !=0 %d", s.ErrorCount)
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

	tok, op, bcPos, endbcPos, data, number := s.Scan()
	if tok != token.OPERATION || op != opcode.DUP || bcPos != 0 || endbcPos != 1 || data != nil || number != nil {
		t.Errorf("Failed DUP tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.HASH160 || bcPos != 1 || endbcPos != 2 || data != nil || number != nil {
		t.Errorf("Failed HASH160 tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.DATA || op != opcode.INVALID || bcPos != 2 || endbcPos != 23 || !bytes.Equal(data, checkData) || number != nil {
		t.Errorf("Failed DATA tok %s op %s bcPos %d endbcPos %d data %s number %s", tok, op, bcPos, endbcPos, hex.EncodeToString(data), number)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.ENDOFSCRIPT || op != opcode.INVALID || bcPos != 23 || endbcPos != 24 || data != nil || number != nil {
		t.Errorf("Failed ENDOFSCRIPT tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.ENDOFSCRIPT || op != opcode.INVALID || bcPos != 23 || endbcPos != 24 || data != nil || number != nil {
		t.Errorf("Failed double ENDOFSCRIPT tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	if !errorReported {
		t.Errorf("Failed No ErrorReported")
	}

	if s.ErrorCount != 1 {
		t.Errorf("Failed ErrorCount !=1 %d", s.ErrorCount)
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

	tok, op, bcPos, endbcPos, data, number := s.Scan()
	if tok != token.OPERATION || op != opcode.DUP || bcPos != 0 || endbcPos != 1 || data != nil || number != nil {
		t.Errorf("Failed DUP tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.HASH160 || bcPos != 1 || endbcPos != 2 || data != nil || number != nil {
		t.Errorf("Failed HASH160 tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.DATA || op != opcode.INVALID || bcPos != 2 || endbcPos != 23 || !bytes.Equal(data, checkData) || number != nil {
		t.Errorf("Failed DATA tok %s op %s bcPos %d endbcPos %d data %s number %s", tok, op, bcPos, endbcPos, hex.EncodeToString(data), number)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.INVALID || op != opcode.INVALID || bcPos != 23 || endbcPos != 24 || data != nil || number != nil {
		t.Errorf("Failed EQUALVERIFY tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.CHECKSIG || bcPos != 24 || endbcPos != 25 || data != nil || number != nil {
		t.Errorf("Failed CHECKSIG tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.ENDOFSCRIPT || op != opcode.INVALID || bcPos != 25 || endbcPos != 26 || data != nil || number != nil {
		t.Errorf("Failed ENDOFSCRIPT tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.ENDOFSCRIPT || op != opcode.INVALID || bcPos != 25 || endbcPos != 26 || data != nil || number != nil {
		t.Errorf("Failed double ENDOFSCRIPT tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	if !errorReported {
		t.Errorf("Failed No ErrorReported")
	}

	if s.ErrorCount != 1 {
		t.Errorf("Failed ErrorCount !=0 %d", s.ErrorCount)
	}

}

func TestInvalidScript(t *testing.T) {
	script := []byte("script")

	var s Scanner
	errorReported := false

	s.Init(script, func(pos int, msg string) {
		errorReported = true
	})

	tok, op, bcPos, endbcPos, data, number := s.Scan()
	if tok != token.OPERATION || op != opcode.IFDUP || bcPos != 0 || endbcPos != 1 || data != nil || number != nil {
		t.Errorf("Failed IFDUP tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.IF || op != opcode.INVALID || bcPos != 1 || endbcPos != 2 || data != nil || number != nil {
		t.Errorf("Failed IF tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.TWOSWAP || bcPos != 2 || endbcPos != 3 || data != nil || number != nil {
		t.Errorf("Failed TWOSWAP tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.VERIFY || bcPos != 3 || endbcPos != 4 || data != nil || number != nil {
		t.Errorf("Failed EQUALVERIFY tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.TWOOVER || bcPos != 4 || endbcPos != 5 || data != nil || number != nil {
		t.Errorf("Failed TWOOVER tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.OPERATION || op != opcode.DEPTH || bcPos != 5 || endbcPos != 6 || data != nil || number != nil {
		t.Errorf("Failed TWOOVER tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.ENDOFSCRIPT || op != opcode.INVALID || bcPos != 6 || endbcPos != 7 || data != nil || number != nil {
		t.Errorf("Failed ENDOFSCRIPT tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	tok, op, bcPos, endbcPos, data, number = s.Scan()
	if tok != token.ENDOFSCRIPT || op != opcode.INVALID || bcPos != 6 || endbcPos != 7 || data != nil || number != nil {
		t.Errorf("Failed double ENDOFSCRIPT tok %s op %s bcPos %d endbcPos %d", tok, op, bcPos, endbcPos)
	}

	if errorReported {
		t.Errorf("Failed ErrorReported")
	}

	if s.ErrorCount != 0 {
		t.Errorf("Failed ErrorCount !=0 %d", s.ErrorCount)
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
