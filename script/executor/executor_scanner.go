package executor

import (
	"errors"
	"github.com/spearson78/guardian/script/scanner"
	"github.com/spearson78/guardian/script/token"
)

func seekEndIf(s *scanner.Scanner) error {
	ifDepth := 0
loop:
	for {

		tok := s.Scan()
		switch tok {
		case token.ENDOFSCRIPT:
			return errors.New("Unclosed IF Statements")
		case token.IF:
			ifDepth++
		case token.ENDIF:
			if ifDepth == 0 {
				break loop
			} else {
				ifDepth--
			}
		}
	}

	return nil

}

func seekElse(s *scanner.Scanner) error {
	ifDepth := 0
loop:
	for {

		tok := s.Scan()
		switch tok {
		case token.ENDOFSCRIPT:
			return errors.New("Unclosed IF Statements")
		case token.IF:
			ifDepth++
		case token.ENDIF:
			if ifDepth == 0 {
				break loop
			} else {
				ifDepth--
			}
		case token.ELSE:
			if ifDepth == 0 {
				break loop
			}
		}
	}

	return nil
}

func (this *Executor) execScanner(script []byte) error {
	s := new(scanner.Scanner)
	s.Init(script, nil)

	ifDepth := 0

loop:
	for {
		tok := s.Scan()
		switch tok {
		case token.DATA:
			this.Push(s.Data())
		case token.OPERATION:
			impl := opCodeImpls[s.Op()]
			if impl == nil {
				return errors.New("Unknown OpCode - " + s.Op().String())
			}

			err := impl(&this.Context)
			if err != nil {
				return err
			}
		case token.ENDOFSCRIPT:
			break loop
		case token.NUMBER:
			this.PushNumber(s.Number())

		case token.CODESEPARATOR:
			this.codeSeparatorPos = s.Pos()
		case token.IF:
			ifDepth++
			top := this.PopBool()
			if !top {
				seekElse(s)
			}
		case token.NOTIF:
			ifDepth++
			top := this.PopBool()
			if top {
				seekElse(s)
			}
		case token.ELSE:
			if ifDepth == 0 {
				return errors.New("Unexpected EndIf")
			}
			ifDepth--
			seekEndIf(s)
		case token.ENDIF:
			if ifDepth == 0 {
				return errors.New("Unexpected Else")
			}
			ifDepth--
		case token.INVALID:
			return errors.New("Invalid Token")

		}
	}

	if ifDepth != 0 {
		return errors.New("Unclosed If")
	}

	if s.ErrorCount() != 0 {
		return errors.New("TokenSource reported errors")
	}

	return nil
}
