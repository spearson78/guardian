package executor

import (
	"math/big"
)

var c0 = big.NewInt(0)
var c1 = big.NewInt(1)

//TODO:Arithmetic is only allowed on 4 byte integers consider a PopInt to validate the value on the stack is compatible

func op_ONEADD(c *Context) error {
	n := c.PopNumber()
	n.Add(n, c1)
	c.PushNumber(n)
	return nil
}

func op_ONESUB(c *Context) error {
	n := c.PopNumber()
	n.Sub(n, c1)
	c.PushNumber(n)
	return nil
}

func op_NEGATE(c *Context) error {
	n := c.PopNumber()
	n.Neg(n)
	c.PushNumber(n)
	return nil
}

func op_ABS(c *Context) error {
	n := c.PopNumber()
	n.Abs(n)
	c.PushNumber(n)
	return nil
}

func op_NOT(c *Context) error {
	c.PushBool(!c.PopBool())
	return nil
}

func op_ZERONOTEQUAL(c *Context) error {
	c.PushBool(c.PopBool())
	return nil
}

func op_ADD(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	a.Add(a, b)

	c.PushNumber(a)
	return nil
}

func op_SUB(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	a.Sub(a, b)

	c.PushNumber(a)
	return nil
}

func op_BOOLAND(c *Context) error {
	b := c.PopBool()
	a := c.PopBool()

	c.PushBool(a && b)
	return nil
}

func op_BOOLOR(c *Context) error {
	b := c.PopBool()
	a := c.PopBool()

	c.PushBool(a || b)
	return nil
}

func op_NUMEQUAL(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	c.PushBool(a.Cmp(b) == 0)

	return nil
}

func op_NUMEQUALVERIFY(c *Context) error {
	return compositeOp(c, op_NUMEQUAL, op_VERIFY)
}

func op_NUMNOTEQUAL(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	c.PushBool(a.Cmp(b) != 0)

	return nil
}

func op_LESSTHAN(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	c.PushBool(a.Cmp(b) < 0)

	return nil
}

func op_GREATERTHAN(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	c.PushBool(a.Cmp(b) > 0)

	return nil
}

func op_LESSTHANOREQUAL(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	c.PushBool(a.Cmp(b) <= 0)

	return nil
}

func op_GREATERTHANOREQUAL(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	c.PushBool(a.Cmp(b) >= 0)

	return nil
}

func op_MIN(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	if a.Cmp(b) > 0 {
		c.PushNumber(b)
	} else {
		c.PushNumber(a)
	}

	return nil
}

func op_MAX(c *Context) error {
	b := c.PopNumber()
	a := c.PopNumber()

	if a.Cmp(b) > 0 {
		c.PushNumber(a)
	} else {
		c.PushNumber(b)
	}

	return nil
}

func op_WITHIN(c *Context) error {
	max := c.PopNumber()
	min := c.PopNumber()
	x := c.PopNumber()

	if x.Cmp(max) > 0 {
		c.PushBool(false)
	} else if x.Cmp(min) < 0 {
		c.PushBool(false)
	} else {
		c.PushBool(true)
	}

	return nil
}
