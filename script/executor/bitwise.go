package executor

import (
	"bytes"
)

func op_EQUAL(c *Context) error {
	a := c.Pop()
	b := c.Pop()

	c.PushBool(bytes.Equal(a, b))

	return nil
}

func op_EQUALVERIFY(c *Context) error {
	return compositeOp(c, op_EQUAL, op_VERIFY)
}
