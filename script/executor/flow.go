package executor

import (
	"errors"
)

func op_NOP(c *Context) error {
	return nil
}

func op_VERIFY(c *Context) error {
	top := c.PopBool()
	if !top {
		return errors.New("OP_VERIFY False Transaction Invalid")
	}

	return nil
}

func op_RETURN(c *Context) error {
	return errors.New("OP_RETURN Transaction Invalid")
}
