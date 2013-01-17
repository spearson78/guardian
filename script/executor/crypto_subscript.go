package executor

import (
	"bytes"
	"errors"
	"github.com/spearson78/guardian/script/scanner"
	"github.com/spearson78/guardian/script/token"
)

func Subscriptify(script []byte, sig []byte) ([]byte, error) {

	s := new(scanner.Scanner)
	s.Init(script, nil)

	subscript := make([]byte, 0, len(script))

loop:
	for {
		tok := s.Scan()
		switch tok {
		case token.DATA:
			//Suppress the signature
			if !bytes.Equal(s.Data(), sig) {
				subscript = append(subscript, s.ByteCode())
				subscript = append(subscript, s.Data()...)
			}
		case token.CODESEPARATOR:
			//Remove code separators
		case token.ENDOFSCRIPT:
			break loop
		case token.INVALID:
			return nil, errors.New("Invalid Token")
		default:
			//All other tokens are passed straight through
			subscript = append(subscript, s.ByteCode())
		}
	}

	if s.ErrorCount() != 0 {
		return nil, errors.New("Scanner Errors Occured")
	}

	return subscript, nil
}
