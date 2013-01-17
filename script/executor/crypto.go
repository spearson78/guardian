package executor

import (
	_ "code.google.com/p/go.crypto/ripemd160"
	"github.com/spearson78/guardian/crypto/sha256d"

	"crypto"
	_ "crypto/sha1"
	_ "crypto/sha256"
	"errors"
	"fmt"
)

func op_RIPEMD160(c *Context) error {

	top := c.Pop()

	h := crypto.RIPEMD160.New()
	_, err := h.Write(top)
	if err != nil {
		return errors.New("Hash Failed")
	}

	res := h.Sum(nil)

	c.Push(res)

	return nil
}

func op_SHA1(c *Context) error {

	top := c.Pop()

	h := crypto.SHA1.New()
	_, err := h.Write(top)
	if err != nil {
		return errors.New("Hash Failed")
	}

	res := h.Sum(nil)

	c.Push(res)

	return nil
}

func op_HASH160(c *Context) error {

	top := c.Pop()

	h := crypto.SHA256.New()
	_, err := h.Write(top)
	if err != nil {
		return errors.New("Hash Failed")
	}

	res := h.Sum(nil)

	h = crypto.RIPEMD160.New()
	_, err = h.Write(res)
	if err != nil {
		return errors.New("Hash Failed")
	}

	res = h.Sum(nil)

	c.Push(res)

	return nil
}

func op_HASH256(c *Context) error {

	top := c.Pop()

	h := sha256d.New()
	_, err := h.Write(top)
	if err != nil {
		return errors.New("Hash Failed")
	}

	res := h.Sum(nil)

	c.Push(res)

	return nil
}

func op_CHECKSIG(c *Context) error {
	if c.signatureCheck == nil {
		return errors.New("No SignatureCheck Implementation")
	}

	pk := c.Pop()
	sig := c.Pop()

	subScript, err := Subscriptify(c.script[c.codeSeparatorPos:], sig)
	if err != nil {
		return err
	}

	hashType := uint32(sig[len(sig)-1])
	sigVal := sig[:len(sig)-1]

	err = c.signatureCheck.CheckSig(pk, hashType, sigVal, subScript)
	if err != nil {
		fmt.Println(err)
	}
	c.PushBool(err == nil)

	return nil
}

func op_CHECKSIGVERIFY(c *Context) error {
	return compositeOp(c, op_CHECKSIG, op_VERIFY)
}

func op_CHECKMULTISIG(c *Context) error {

	pkcount := int(c.PopNumber().Int64())
	for i := 0; i < pkcount; i++ {
		_ = c.Pop()
	}
	sigcount := int(c.PopNumber().Int64())
	for i := 0; i < sigcount; i++ {
		_ = c.Pop()
	}

	//Due to a bug in the reference client, one extra unused value is removed from the stack. 
	_ = c.Pop()

	//TODO: Implement the necessary checkmultisig callback
	c.PushBool(false)

	return nil
}

func op_CHECKMULTISIGVERIFY(c *Context) error {
	return compositeOp(c, op_CHECKMULTISIG, op_VERIFY)
}
