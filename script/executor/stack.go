package executor

import (
	"errors"
	"github.com/spearson78/guardian/encoding/scriptint"
	"math/big"
)

func op_TOALTSTACK(c *Context) error {
	c.AltPush(c.Pop())
	return nil
}

func op_FROMALTSTACK(c *Context) error {
	c.Push(c.AltPop())
	return nil
}

func op_IFDUP(c *Context) error {
	top := c.Peek()
	if top == nil {
		return errors.New("Stack Underflow")
	}
	if scriptint.Decode(top).Sign() == 0 {
		c.Push(top)
	}
	return nil
}

func op_DEPTH(c *Context) error {
	c.PushNumber(new(big.Int).SetInt64(int64(c.Depth())))
	return nil
}

func op_DROP(c *Context) error {
	c.Pop()
	return nil
}

func op_DUP(c *Context) error {
	top := c.Peek()
	if top == nil {
		return errors.New("Stack Underflow")
	}
	c.Push(top)

	return nil
}

func op_NIP(c *Context) error {
	x2 := c.Pop()
	c.Pop()

	c.Push(x2)
	return nil

}

func op_OVER(c *Context) error {
	x1 := c.PeekN(1)
	c.Push(x1)

	return nil
}

func op_PICK(c *Context) error {

	n := int(c.PopNumber().Int64())

	xn := c.PeekN(n)
	c.Push(xn)

	return nil
}

func op_ROLL(c *Context) error {

	n := int(c.PopNumber().Int64())

	xn := c.PopN(n)
	c.Push(xn)

	return nil
}

func op_ROT(c *Context) error {

	x3 := c.Pop()
	x2 := c.Pop()
	x1 := c.Pop()

	c.Push(x2)
	c.Push(x3)
	c.Push(x1)

	return nil
}

func op_SWAP(c *Context) error {

	x2 := c.Pop()
	x1 := c.Pop()

	c.Push(x2)
	c.Push(x1)

	return nil
}

func op_TUCK(c *Context) error {

	x2 := c.Pop()
	x1 := c.Pop()

	c.Push(x2)
	c.Push(x1)
	c.Push(x2)

	return nil
}

func op_TWODROP(c *Context) error {

	c.Pop()
	c.Pop()

	return nil
}

func op_TWODUP(c *Context) error {

	x1 := c.PeekN(1)
	x2 := c.PeekN(0)

	c.Push(x1)
	c.Push(x2)

	return nil
}

func op_THREEDUP(c *Context) error {

	x1 := c.PeekN(2)
	x2 := c.PeekN(1)
	x3 := c.PeekN(0)

	c.Push(x1)
	c.Push(x2)
	c.Push(x3)

	return nil
}

func op_TWOOVER(c *Context) error {

	x1 := c.PeekN(3)
	x2 := c.PeekN(2)

	c.Push(x1)
	c.Push(x2)

	return nil
}

func op_TWOROT(c *Context) error {

	x2 := c.PopN(4)
	x1 := c.PopN(4)

	c.Push(x1)
	c.Push(x2)

	return nil
}

func op_TWOSWAP(c *Context) error {

	x2 := c.PopN(2)
	x1 := c.PopN(2)

	c.Push(x1)
	c.Push(x2)

	return nil
}
