package scanner

import (
	"github.com/spearson78/guardian/script/opcode"
	"github.com/spearson78/guardian/script/token"
	"math/big"
)

type ErrorHandler func(byteCodePos int, msg string)

type Scanner struct {
	script     []byte
	pos        int
	ErrorCount int
	err        ErrorHandler
}

var _EMPTY_SLICE = []byte{}

func (s *Scanner) Init(script []byte, err ErrorHandler) {
	s.script = script
	s.pos = 0
	s.err = err
}

func (s *Scanner) raiseError(msg string) {
	s.ErrorCount++
	if s.err != nil {
		s.err(s.pos, msg)
	}
}

func (s *Scanner) Scan() (tok token.Token, op opcode.OpCode, byteCodePos int, endByteCodePos int, data []byte, number *big.Int) {

	tok = token.INVALID
	op = opcode.INVALID
	byteCodePos = s.pos
	data = nil
	number = nil

	if s.pos >= len(s.script) {
		tok = token.ENDOFSCRIPT
	} else {
		bytecode := s.script[s.pos]

		switch {
		case bytecode == 0x00: //OP_0
			tok = token.DATA
			data = _EMPTY_SLICE
		case bytecode <= byte(0x4b): //PUSH CONSTANT
			tok = token.DATA
			dataPos := s.pos + 1
			endOfData := dataPos + int(bytecode)
			if len(s.script) <= endOfData {
				s.raiseError("Script Underflow")
				endOfData = len(s.script)
			}
			data = s.script[dataPos:endOfData]
			s.pos = s.pos + int(bytecode)
		case bytecode == 0x4c: //PUSHDATA1
			tok = token.DATA
			var byteCount = int(s.script[s.pos+1])
			dataPos := s.pos + 2

			endOfData := dataPos + byteCount
			if len(s.script) <= endOfData {
				s.raiseError("Script Underflow")
				endOfData = len(s.script)
			}
			data = s.script[dataPos:endOfData]

			s.pos = s.pos + byteCount + 1
		case bytecode == 0x4d: //PUSHDATA2
			tok = token.DATA
			byteCount := uint16(s.script[s.pos+1]) | uint16(s.script[s.pos+2]<<8)
			dataPos := s.pos + 3

			endOfData := dataPos + int(byteCount)
			if len(s.script) <= endOfData {
				s.raiseError("Script Underflow")
				endOfData = len(s.script)
			}
			data = s.script[dataPos:endOfData]

			s.pos = s.pos + int(byteCount) + 2
		case bytecode == 0x4e: //PUSHDATA4
			tok = token.DATA
			byteCount := uint32(s.script[s.pos+1]) | uint32(s.script[s.pos+2]<<8) | uint32(s.script[s.pos+3]<<16) | uint32(s.script[s.pos+4]<<24)
			dataPos := s.pos + 5
			endOfData := dataPos + int(byteCount)
			if len(s.script) <= endOfData {
				s.raiseError("Script Underflow")
				endOfData = len(s.script)
			}
			data = s.script[dataPos:endOfData]
			s.pos = s.pos + int(byteCount) + 2
		case bytecode == 0x4f: //1Negate
			tok = token.NUMBER
			number = big.NewInt(-1)
		case bytecode <= byte(0x60):
			tok = token.NUMBER
			number = big.NewInt(int64(bytecode - 80))
		case opcode.OpCode(bytecode).IsValid():
			tok = token.OPERATION
			op = opcode.OpCode(bytecode)
		case bytecode == 0xab: //CODESEPARATOR
			tok = token.CODESEPARATOR
		case bytecode == 0x63: //IF
			tok = token.IF
		case bytecode == 0x64: //NOTIF
			tok = token.NOTIF
		case bytecode == 0x67: //ELSE
			tok = token.ELSE
		case bytecode == 0x68: //ENDIF
			tok = token.ENDIF
		default:
			tok = token.INVALID
			op = opcode.INVALID
			data = nil
			number = nil
			s.raiseError("Invalid Token")
		}
	}

	endByteCodePos = s.pos + 1

	if tok != token.ENDOFSCRIPT {
		s.pos++
	}

	return
}
