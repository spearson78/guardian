package ast

import (
	"errors"
	"github.com/spearson78/guardian/script/opcode"
	"github.com/spearson78/guardian/script/token"
	"math/big"
)

type TokenSource interface {
	Scan() (tok token.Token)

	Op() opcode.OpCode
	Data() []byte
	Number() *big.Int

	Pos() int

	ErrorCount() int
}

type _IfStack struct {
	stack []*IfStmt
}

func (this *_IfStack) Push(x *IfStmt) {
	this.stack = append(this.stack, x)
}

func (this *_IfStack) Peek() *IfStmt {
	return this.stack[len(this.stack)-1]
}

func (this *_IfStack) Pop() *IfStmt {
	pop := this.stack[len(this.stack)-1]
	this.stack = this.stack[:len(this.stack)-1]
	return pop
}

func (this *_IfStack) Len() int {
	return len(this.stack)
}

func Parse(s TokenSource) (Block, error) {

	root := new(SimpleBlock)

	var currentBlock Block = root
	var ifStack _IfStack

	for {
		tok := s.Scan()
		if tok == token.ENDOFSCRIPT {
			if ifStack.Len() != 0 {
				return nil, errors.New("Unclosed IF Statements")
			}
			break
		}

		switch tok {
		case token.INVALID:
			return nil, errors.New("Invalid Token")
		case token.DATA:
			currentBlock.Append(&Data{
				ParentBlock: currentBlock,
				DataPos:     s.Pos(),
				Value:       s.Data(),
			})
		case token.NUMBER:
			currentBlock.Append(&Number{
				ParentBlock: currentBlock,
				NumberPos:   s.Pos(),
				Value:       s.Number(),
			})
		case token.OPERATION:
			currentBlock.Append(&Operation{
				ParentBlock:  currentBlock,
				OperationPos: s.Pos(),
				OpCode:       s.Op(),
			})
		case token.CODESEPARATOR:
			codeSeparator := &CodeSeparator{
				ParentBlock:      currentBlock,
				CodeSeparatorPos: s.Pos(),
			}
			currentBlock.Append(codeSeparator)
			currentBlock = codeSeparator
		case token.IF:
			ifStmt := &IfStmt{
				ParentBlock: currentBlock,
				Not:         false,
				IfPos:       s.Pos(),
				Body:        new(SimpleBlock),
			}
			currentBlock.Append(ifStmt)
			currentBlock = ifStmt.Body
			ifStack.Push(ifStmt)
		case token.NOTIF:
			ifStmt := &IfStmt{
				ParentBlock: currentBlock,
				Not:         true,
				IfPos:       s.Pos(),
				Body:        new(SimpleBlock),
			}
			currentBlock.Append(ifStmt)
			currentBlock = ifStmt.Body
			ifStack.Push(ifStmt)
		case token.ELSE:
			currentIf := ifStack.Peek()
			currentIf.ElsePos = s.Pos()
			currentIf.Else = new(SimpleBlock)
			currentBlock = currentIf.Else
		case token.ENDIF:
			currentIf := ifStack.Pop()
			currentIf.EndIfPos = s.Pos()
			currentBlock = currentIf.ParentBlock
		}
	}

	if s.ErrorCount() != 0 {
		return nil, errors.New("TokenSource reported errors")
	}

	return root, nil

}
