package executor

import (
	"errors"
	"github.com/spearson78/guardian/script/ast"
	"github.com/spearson78/guardian/script/scanner"
)

func (this *Executor) Visit(node ast.Node) (bool, error) {

	switch n := node.(type) {
	case *ast.CodeSeparator:
		this.Context.codeSeparatorPos = n.CodeSeparatorPos
		return true, n.ForEachNode(this)
	case ast.Block:
		return true, n.ForEachNode(this)
	case *ast.Operation:

		impl := opCodeImpls[n.OpCode]
		if impl == nil {
			return false, errors.New("Unknown OpCode - " + n.OpCode.String())
		}

		err := impl(&this.Context)
		if err != nil {
			return false, err
		}

	case *ast.Data:
		this.Push(n.Value)
	case *ast.Number:
		this.PushNumber(n.Value)
	case *ast.IfStmt:
		top := this.PopBool()

		if top != n.Not {
			return true, n.Body.ForEachNode(this)
		} else {
			return true, n.Else.ForEachNode(this)
		}
	default:
		return false, errors.New("Unknown Node Type - " + n.String())
	}

	return true, nil
}

func (this *Executor) execAst(script []byte) error {
	s := new(scanner.Scanner)
	s.Init(script, nil)
	block, err := ast.Parse(s)
	if err != nil {
		return err
	}

	err = block.ForEachNode(this)
	if err != nil {
		return err
	}

	if s.ErrorCount() != 0 {
		return errors.New("TokenSource reported errors")
	}

	return nil
}
