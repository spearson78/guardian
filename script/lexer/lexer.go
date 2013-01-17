package lexer

import (
	"encoding/hex"
	"github.com/spearson78/guardian/script/opcode"
	"github.com/spearson78/guardian/script/token"
	"io"
	"math/big"
	"strings"
	"text/scanner"
)

type ErrorHandler func(byteCodePos int, msg string)

type Lexer struct {
	s          scanner.Scanner
	errorCount int
	err        ErrorHandler

	op     opcode.OpCode
	data   []byte
	number *big.Int
}

func (l *Lexer) Init(r io.Reader, err ErrorHandler) {
	l.s.Init(r)
	l.s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments
	l.err = err
}

func (l *Lexer) raiseError(msg string) {
	l.errorCount++
	if l.err != nil {
		l.err(0, msg)
	}
}

func (l *Lexer) Pos() int {
	return l.s.Pos().Line
}

func (l *Lexer) Op() opcode.OpCode {
	return l.op
}

func (l *Lexer) Data() []byte {
	return l.data
}

func (l *Lexer) Number() *big.Int {
	return l.number
}

func (l *Lexer) ErrorCount() int {
	return l.errorCount
}

func (l *Lexer) Scan() (tok token.Token) {

	var err error

	tok = token.INVALID
	l.op = opcode.INVALID
	l.data = nil
	l.number = nil

	stok := l.s.Scan()
	switch stok {
	case scanner.Ident:
		switch l.s.TokenText() {
		case "IF":
			tok = token.IF
		case "NOTIF":
			tok = token.NOTIF
		case "ELSE":
			tok = token.ELSE
		case "ENDIF":
			tok = token.ENDIF
		case "CODESEPARATOR":
			tok = token.CODESEPARATOR
		default:
			tok = token.OPERATION
			l.op = opcode.Parse(l.s.TokenText())
			if l.op == opcode.INVALID {
				l.raiseError("Invalid Operation " + l.s.TokenText())
			}
		}
	case scanner.Int:
		if strings.HasPrefix(l.s.TokenText(), "0x") {
			tok = token.DATA
			l.data, err = hex.DecodeString(l.s.TokenText()[2:])
			if err != nil {
				l.raiseError("Hex Decode Failed " + l.s.TokenText() + err.Error())
			}
		} else {
			tok = token.NUMBER
			l.number = new(big.Int)
			l.number.SetString(l.s.TokenText(), 10)
		}
	case scanner.Float:
		l.raiseError("Floating Point Not Supported " + l.s.TokenText())
	case scanner.Char:
		tok = token.DATA
		l.data = []byte(l.s.TokenText())
	case scanner.String:
		tok = token.DATA
		l.data = []byte(l.s.TokenText())
	case scanner.RawString:
		tok = token.DATA
		l.data = []byte(l.s.TokenText())
	case scanner.EOF:
		tok = token.ENDOFSCRIPT
	case '-':
		tok = l.Scan()
		if tok == token.NUMBER {
			l.number.Neg(l.number)
		} else {
			tok = token.INVALID
		}
	}

	return
}
