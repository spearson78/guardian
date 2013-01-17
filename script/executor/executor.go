package executor

import (
	"errors"
	"github.com/spearson78/guardian/encoding/scriptint"
	"github.com/spearson78/guardian/script/opcode"
	"math/big"
)

type SignatureCheck interface {
	CheckSig(pk []byte, hashType uint32, sig []byte, subScript []byte) error
}

type Context struct {
	stack    [][]byte
	altstack [][]byte

	script           []byte
	codeSeparatorPos int
	signatureCheck   SignatureCheck
}

type OpCodeImplementation func(context *Context) error

func op_DISABLED(c *Context) error {
	return errors.New("DISABLED OPCODE")
}

func compositeOp(c *Context, ops ...OpCodeImplementation) error {
	for _, op := range ops {
		err := op(c)
		if err != nil {
			return err
		}
	}

	return nil
}

var opCodeImpls = [...]OpCodeImplementation{
	//Flow
	opcode.NOP:    op_NOP,
	opcode.VERIFY: op_VERIFY,
	opcode.RETURN: op_RETURN,

	//Stack
	opcode.TOALTSTACK:   op_TOALTSTACK,
	opcode.FROMALTSTACK: op_FROMALTSTACK,
	opcode.IFDUP:        op_IFDUP,
	opcode.DEPTH:        op_DEPTH,
	opcode.DROP:         op_DROP,
	opcode.DUP:          op_DUP,
	opcode.NIP:          op_NIP,
	opcode.OVER:         op_OVER,
	opcode.PICK:         op_PICK,
	opcode.ROLL:         op_ROLL,
	opcode.ROT:          op_ROT,
	opcode.SWAP:         op_SWAP,
	opcode.TUCK:         op_TUCK,
	opcode.TWODROP:      op_TWODROP,
	opcode.TWODUP:       op_TWODUP,
	opcode.THREEDUP:     op_THREEDUP,
	opcode.TWOOVER:      op_TWOOVER,
	opcode.TWOROT:       op_TWOROT,
	opcode.TWOSWAP:      op_TWOSWAP,

	//Splice
	opcode.CAT:    op_DISABLED,
	opcode.SUBSTR: op_DISABLED,
	opcode.LEFT:   op_DISABLED,
	opcode.RIGHT:  op_DISABLED,
	opcode.SIZE:   op_DISABLED,

	//Bitwise
	opcode.INVERT:      op_DISABLED,
	opcode.AND:         op_DISABLED,
	opcode.OR:          op_DISABLED,
	opcode.XOR:         op_DISABLED,
	opcode.EQUAL:       op_EQUAL,
	opcode.EQUALVERIFY: op_EQUALVERIFY,

	//Arithmetic
	opcode.ONEADD:             op_ONEADD,
	opcode.ONESUB:             op_ONESUB,
	opcode.TWOMUL:             op_DISABLED,
	opcode.TWODIV:             op_DISABLED,
	opcode.NEGATE:             op_NEGATE,
	opcode.ABS:                op_ABS,
	opcode.NOT:                op_NOT,
	opcode.ZERONOTEQUAL:       op_ZERONOTEQUAL,
	opcode.ADD:                op_ADD,
	opcode.SUB:                op_SUB,
	opcode.MUL:                op_DISABLED,
	opcode.DIV:                op_DISABLED,
	opcode.MOD:                op_DISABLED,
	opcode.LSHIFT:             op_DISABLED,
	opcode.RSHIFT:             op_DISABLED,
	opcode.BOOLAND:            op_BOOLAND,
	opcode.BOOLOR:             op_BOOLOR,
	opcode.NUMEQUAL:           op_NUMEQUAL,
	opcode.NUMEQUALVERIFY:     op_NUMEQUALVERIFY,
	opcode.NUMNOTEQUAL:        op_NUMNOTEQUAL,
	opcode.LESSTHAN:           op_LESSTHAN,
	opcode.GREATERTHAN:        op_GREATERTHAN,
	opcode.LESSTHANOREQUAL:    op_LESSTHANOREQUAL,
	opcode.GREATERTHANOREQUAL: op_GREATERTHANOREQUAL,
	opcode.MIN:                op_MIN,
	opcode.MAX:                op_MAX,
	opcode.WITHIN:             op_WITHIN,

	//Crypto
	opcode.RIPEMD160:           op_RIPEMD160,
	opcode.SHA1:                op_SHA1,
	opcode.HASH160:             op_HASH160,
	opcode.HASH256:             op_HASH256,
	opcode.CHECKSIG:            op_CHECKSIG,
	opcode.CHECKSIGVERIFY:      op_CHECKSIGVERIFY,
	opcode.CHECKMULTISIG:       op_CHECKMULTISIG,
	opcode.CHECKMULTISIGVERIFY: op_CHECKMULTISIGVERIFY,

	//Reserved NOPs
	opcode.NOP1:  op_NOP,
	opcode.NOP2:  op_NOP,
	opcode.NOP3:  op_NOP,
	opcode.NOP4:  op_NOP,
	opcode.NOP5:  op_NOP,
	opcode.NOP6:  op_NOP,
	opcode.NOP7:  op_NOP,
	opcode.NOP8:  op_NOP,
	opcode.NOP9:  op_NOP,
	opcode.NOP10: op_NOP,
}

type Executor struct {
	Context
}

func (this *Executor) Init(signatureCheck SignatureCheck) {
	this.Context.signatureCheck = signatureCheck
}

func (this *Context) Push(data []byte) {
	this.stack = append(this.stack, data)
}

func (this *Context) Depth() int {
	return len(this.stack)
}

func (this *Context) AltPush(data []byte) {
	this.altstack = append(this.altstack, data)
}

func (this *Context) AltPop() []byte {
	if len(this.altstack) == 0 {
		panic("Alt Stack Underflow")
	}

	top := this.altstack[len(this.altstack)-1]
	if len(this.altstack) == 1 {
		this.altstack = nil
	} else {
		this.altstack = this.altstack[:len(this.altstack)-1]
	}

	return top
}

func (this *Context) Peek() []byte {
	if len(this.stack) == 0 {
		return nil
	}

	return this.stack[len(this.stack)-1]
}

func (this *Context) PeekN(n int) []byte {
	if len(this.stack) <= n {
		return nil
	}

	return this.stack[len(this.stack)-1-n]
}

func (this *Context) PopN(n int) []byte {
	if len(this.stack) <= n {
		return nil
	}

	xn := this.stack[len(this.stack)-1-n]

	poppedStack := this.stack[:len(this.stack)-1-n]
	remainingStack := this.stack[len(this.stack)-n:]

	this.stack = append(poppedStack, remainingStack...)

	return xn
}

func (this *Context) PopNumber() *big.Int {
	return scriptint.Decode(this.Pop())
}

func (this *Context) PopBool() bool {
	return scriptint.Decode(this.Pop()).Sign() != 0
}

func (this *Context) PushNumber(n *big.Int) {
	this.Push(scriptint.Encode(n))
}

func (this *Context) PushBool(b bool) {
	if b {
		this.Push([]byte{1})
	} else {
		this.Push([]byte{})
	}
}

func (this *Context) Pop() []byte {
	if len(this.stack) == 0 {
		panic("Stack Underflow")
	}

	top := this.stack[len(this.stack)-1]
	if len(this.stack) == 1 {
		this.stack = nil
	} else {
		this.stack = this.stack[:len(this.stack)-1]
	}

	return top
}

//Must use []byte as checksig needs access to the script
func (this *Executor) Execute(script []byte) error {

	this.codeSeparatorPos = 0
	this.script = script

	//return this.execScanner(script)
	return this.execAst(script)
}
