package ast

import (
	"encoding/hex"
	"github.com/spearson78/guardian/script/lexer"
	"github.com/spearson78/guardian/script/opcode"
	"github.com/spearson78/guardian/script/scanner"
	"math/big"
	"strings"
	"testing"
)

func TestStandardTransactionToBitcoinAddressLexer(t *testing.T) {

	script := `DUP
HASH160
0x89abcdefabbaabbaabbaabbaabbaabbaabbaabba
EQUALVERIFY
CHECKSIG
`

	checkData, _ := hex.DecodeString("89abcdefabbaabbaabbaabbaabbaabbaabbaabba")

	l := new(lexer.Lexer)
	l.Init(strings.NewReader(script), nil)

	tree, err := Parse(l)
	if err != nil {
		t.Errorf("Failed %v", err)
	}

	root := new(SimpleBlock)
	root.NodeList = []Node{
		&Operation{
			ParentBlock: root,
			OpCode:      opcode.DUP,
		},
		&Operation{
			ParentBlock: root,
			OpCode:      opcode.HASH160,
		},
		&Data{
			ParentBlock: root,
			Value:       checkData,
		},
		&Operation{
			ParentBlock: root,
			OpCode:      opcode.EQUALVERIFY,
		},
		&Operation{
			ParentBlock: root,
			OpCode:      opcode.CHECKSIG,
		},
	}

	if !root.Equals(tree) {
		t.Errorf("AST Mismatch")
	}

}

func TestStandardTransactionToBitcoinAddressScanner(t *testing.T) {

	script, _ := hex.DecodeString("76A91489ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA88AC")
	checkData, _ := hex.DecodeString("89ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA")

	s := new(scanner.Scanner)

	errorReported := false

	s.Init(script, func(pos int, msg string) {
		errorReported = true
	})

	tree, err := Parse(s)
	if err != nil {
		t.Errorf("Failed %v", err)
	}

	root := new(SimpleBlock)
	root.NodeList = []Node{
		&Operation{
			ParentBlock: root,
			OpCode:      opcode.DUP,
		},
		&Operation{
			ParentBlock: root,
			OpCode:      opcode.HASH160,
		},
		&Data{
			ParentBlock: root,
			Value:       checkData,
		},
		&Operation{
			ParentBlock: root,
			OpCode:      opcode.EQUALVERIFY,
		},
		&Operation{
			ParentBlock: root,
			OpCode:      opcode.CHECKSIG,
		},
	}

	if !root.Equals(tree) {
		t.Errorf("AST Mismatch")
	}

}

func TestIfElseEndif(t *testing.T) {

	script := `
IF
1
ELSE
2
ENDIF
`

	l := new(lexer.Lexer)
	l.Init(strings.NewReader(script), nil)

	tree, err := Parse(l)
	if err != nil {
		t.Errorf("Failed %v", err)
	}

	root := new(SimpleBlock)
	root.NodeList = []Node{
		&IfStmt{
			Body: &SimpleBlock{
				NodeList: []Node{
					&Number{
						Value: big.NewInt(1),
					},
				},
			},
			Else: &SimpleBlock{
				NodeList: []Node{
					&Number{
						Value: big.NewInt(2),
					},
				},
			},
		},
	}

	if !root.Equals(tree) {
		t.Errorf("AST Mismatch")
	}

}

func BenchmarkStandardTransactionToBitcoinAddress(b *testing.B) {
	script, _ := hex.DecodeString("76A91489ABCDEFABBAABBAABBAABBAABBAABBAABBAABBA88AC")

	s := new(scanner.Scanner)

	for i := 0; i < b.N; i++ {
		s.Init(script, nil)
		Parse(s)
	}
}

func BenchmarkNops(b *testing.B) {
	script, _ := hex.DecodeString("61616161616161616161")

	s := new(scanner.Scanner)

	for i := 0; i < b.N; i++ {
		s.Init(script, nil)
		Parse(s)
	}
}
