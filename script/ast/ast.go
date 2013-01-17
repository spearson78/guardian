package ast

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/spearson78/guardian/script/opcode"
	"math/big"
)

type Visitor interface {
	Visit(node Node) (cont bool, err error)
}

type Node interface {
	Parent() Block
	Pos() int
	String() string
	Equals(node Node) bool
}

type Block interface {
	Node
	List() []Node
	Append(node Node)
	ForEachNode(v Visitor) error
}

type SimpleBlock struct {
	ParentBlock Block
	NodeList    []Node
}

type CodeSeparator struct {
	CodeSeparatorPos int
	ParentBlock      Block
	NodeList         []Node
}

type Operation struct {
	ParentBlock  Block
	OperationPos int
	OpCode       opcode.OpCode
}

type Data struct {
	ParentBlock Block
	DataPos     int
	Value       []byte
}

type Number struct {
	ParentBlock Block
	NumberPos   int
	Value       *big.Int
}

type IfStmt struct {
	ParentBlock Block
	Not         bool
	IfPos       int
	ElsePos     int
	EndIfPos    int
	Body        Block
	Else        Block
}

func (this *SimpleBlock) Parent() Block {
	return this.ParentBlock
}

func (this *SimpleBlock) List() []Node {
	return this.NodeList
}

func (this *SimpleBlock) Append(node Node) {
	this.NodeList = append(this.NodeList, node)
}

func (this *SimpleBlock) Pos() int {
	return this.NodeList[0].Pos()
}

func (this *SimpleBlock) End() int {
	return this.NodeList[len(this.NodeList)-1].Pos() + 1
}

func (this *SimpleBlock) ForEachNode(v Visitor) error {
	for _, node := range this.NodeList {
		cont, err := v.Visit(node)
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}

	return nil
}

func (this *SimpleBlock) Equals(node Node) bool {

	other, ok := node.(*SimpleBlock)
	if !ok {
		return false
	}

	if len(this.NodeList) != len(other.NodeList) {
		return false
	}

	for i, thisSubNode := range this.NodeList {
		otherSubNode := other.NodeList[i]
		if !thisSubNode.Equals(otherSubNode) {
			return false
		}
	}

	return true
}

func (this *SimpleBlock) String() string {
	var buffer bytes.Buffer
	for _, node := range this.NodeList {
		buffer.WriteString(node.String())
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func (this *CodeSeparator) Parent() Block {
	return this.ParentBlock
}

func (this *CodeSeparator) List() []Node {
	return this.NodeList
}

func (this *CodeSeparator) Append(node Node) {
	this.NodeList = append(this.NodeList, node)
}

func (this *CodeSeparator) Pos() int {
	return this.CodeSeparatorPos
}

func (this *CodeSeparator) End() int {
	return this.NodeList[len(this.NodeList)-1].Pos() + 1
}

func (this *CodeSeparator) ForEachNode(v Visitor) error {
	for _, node := range this.NodeList {
		cont, err := v.Visit(node)
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}

	return nil
}

func (this *CodeSeparator) Equals(node Node) bool {
	_, ok := node.(*CodeSeparator)
	if !ok {
		return false
	}

	return true
}

func (this *CodeSeparator) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("CODESEPARATOR\n")
	for _, node := range this.NodeList {
		buffer.WriteString(node.String())
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func (this *Operation) Parent() Block {
	return this.ParentBlock
}

func (this *Operation) Pos() int {
	return this.OperationPos
}

func (this *Operation) Equals(node Node) bool {
	other, ok := node.(*Operation)
	if !ok {
		return false
	}

	return this.OpCode == other.OpCode
}

func (this *Operation) String() string {
	return this.OpCode.String()
}

func (this *Data) Parent() Block {
	return this.ParentBlock
}

func (this *Data) Pos() int {
	return this.DataPos
}

func (this *Data) Equals(node Node) bool {
	other, ok := node.(*Data)
	if !ok {
		return false
	}

	return bytes.Equal(this.Value, other.Value)
}

func (this *Data) String() string {
	return fmt.Sprint("0x", hex.EncodeToString(this.Value))
}

func (this *Number) Parent() Block {
	return this.ParentBlock
}

func (this *Number) Pos() int {
	return this.NumberPos
}

func (this *Number) Equals(node Node) bool {
	other, ok := node.(*Number)
	if !ok {
		return false
	}

	return this.Value.Cmp(other.Value) == 0
}

func (this *Number) String() string {
	return fmt.Sprint(this.Value)
}

func (this *IfStmt) Parent() Block {
	return this.ParentBlock
}

func (this *IfStmt) Pos() int {
	return this.IfPos
}

func (this *IfStmt) End() int {
	return this.EndIfPos
}

func (this *IfStmt) Equals(node Node) bool {
	other, ok := node.(*IfStmt)
	if !ok {
		return false
	}

	if this.Not != other.Not {
		return false
	}

	if !this.Body.Equals(other.Body) {
		return false
	}

	if this.Else != nil {
		if other.Else != nil {
			if !this.Else.Equals(other.Else) {
				return false
			}
		} else {
			return false
		}
	} else {
		if other.Else != nil {
			return false
		}
	}

	return true
}

func (this *IfStmt) String() string {
	var buffer bytes.Buffer
	if this.Not {
		buffer.WriteString("NOTIF\n")
	} else {
		buffer.WriteString("IF\n")
	}
	buffer.WriteString(this.Body.String())
	if this.Else != nil {
		buffer.WriteString("ELSE\n")
		buffer.WriteString(this.Else.String())
	}
	buffer.WriteString("ENDIF\n")

	return buffer.String()
}
