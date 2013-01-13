package token

type Token byte

const (
	//Constants
	INVALID Token = iota
	ENDOFSCRIPT
	DATA
	NUMBER
	OPERATION
	CODESEPARATOR
	IF
	NOTIF
	ELSE
	ENDIF
)

var tokens = [...]string{
	INVALID:       "INVALID",
	ENDOFSCRIPT:   "ENDOFSCRIPT",
	DATA:          "DATA",
	NUMBER:        "NUMBER",
	OPERATION:     "OPERATION",
	CODESEPARATOR: "CODESEPARATOR",
	IF:            "IF",
	NOTIF:         "NOTIF",
	ELSE:          "ELSE",
	ENDIF:         "ENDIF",
}

func (this Token) String() string {
	return tokens[this]
}
