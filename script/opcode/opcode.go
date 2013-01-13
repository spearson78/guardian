package opcode

type OpCode byte

const (
	//Flow Control
	NOP    OpCode = 0x61
	VERIFY OpCode = 0x69
	RETURN OpCode = 0x6a

	//Stack
	TOALTSTACK   OpCode = 0x6b
	FROMALTSTACK OpCode = 0x6c
	IFDUP        OpCode = 0x73
	DEPTH        OpCode = 0x74
	DROP         OpCode = 0x75
	DUP          OpCode = 0x76
	NIP          OpCode = 0x77
	OVER         OpCode = 0x78
	PICK         OpCode = 0x79
	ROLL         OpCode = 0x7a
	ROT          OpCode = 0x7b
	SWAP         OpCode = 0x7c
	TUCK         OpCode = 0x7d
	TWODROP      OpCode = 0x6d
	TWODUP       OpCode = 0x6e
	THREEDUP     OpCode = 0x6f
	TWOOVER      OpCode = 0x70
	TWOROT       OpCode = 0x71
	TWOSWAP      OpCode = 0x72

	//Splice
	CAT    OpCode = 0x7e
	SUBSTR OpCode = 0x7f
	LEFT   OpCode = 0x80
	RIGHT  OpCode = 0x81
	SIZE   OpCode = 0x82

	//Bitwise
	INVERT      OpCode = 0x83
	AND         OpCode = 0x84
	OR          OpCode = 0x85
	XOR         OpCode = 0x86
	EQUAL       OpCode = 0x87
	EQUALVERIFY OpCode = 0x88

	//Arithmetic
	ONEADD            OpCode = 0x8b
	ONESUB            OpCode = 0x8c
	TWOMUL            OpCode = 0x8d
	TWODIV            OpCode = 0x8e
	NEGATE            OpCode = 0x8f
	ABS               OpCode = 0x90
	NOT               OpCode = 0x91
	ZERONOTEQUAL      OpCode = 0x92
	ADD               OpCode = 0x93
	SUB               OpCode = 0x94
	MUL               OpCode = 0x95
	DIV               OpCode = 0x96
	MOD               OpCode = 0x97
	LSHIFT            OpCode = 0x98
	RSHIFT            OpCode = 0x99
	BOOLAND           OpCode = 0x9a
	BOOLOR            OpCode = 0x9b
	NUMEQUAL          OpCode = 0x9c
	NUMEQUALVERIFY    OpCode = 0x9d
	NUMNOTEQUAL       OpCode = 0x9e
	LESSTHAN          OpCode = 0x9f
	GREATERTHAN       OpCode = 0xa0
	LESSTHANOREQUAL   OpCode = 0xa1
	GREATETHANOREQUAL OpCode = 0xa2
	MIN               OpCode = 0xa3
	MAX               OpCode = 0xa4
	WITHIN            OpCode = 0xa5

	//Crypto
	RIPEMD160           OpCode = 0xa6
	SHA1                OpCode = 0xa7
	SHA256              OpCode = 0xa8
	HASH160             OpCode = 0xa9
	HASH256             OpCode = 0xaa
	CODESEPARATOR       OpCode = 0xab
	CHECKSIG            OpCode = 0xac
	CHECKSIGVERIFY      OpCode = 0xad
	CHECKMULTISIG       OpCode = 0xae
	CHECKMULTISIGVERIFY OpCode = 0xaf

	//Pseudo
	PUBKEYHASH OpCode = 0xfd
	PUBKEY     OpCode = 0xfe
	INVALID    OpCode = 0xff

	//Reserved
	RESERVED  OpCode = 0x50
	VER       OpCode = 0x62
	VERIF     OpCode = 0x65
	VERNOTIF  OpCode = 0x66
	RESERVED1 OpCode = 0x89
	RESERVED2 OpCode = 0x8a
	NOP1      OpCode = 0xb0
	NOP2      OpCode = 0xb1
	NOP3      OpCode = 0xb2
	NOP4      OpCode = 0xb3
	NOP5      OpCode = 0xb4
	NOP6      OpCode = 0xb5
	NOP7      OpCode = 0xb6
	NOP8      OpCode = 0xb7
	NOP9      OpCode = 0xb8
	NOP10     OpCode = 0xb9
)

var opcodes = [...]string{
	//Flow Control
	NOP:    "NOP",
	VERIFY: "VERIFY",
	RETURN: "RETURN",

	//Stack
	TOALTSTACK:   "TOALTSTACK",
	FROMALTSTACK: "FROMALTSTACK",
	IFDUP:        "IFDUP",
	DEPTH:        "DEPTH",
	DROP:         "DROP",
	DUP:          "DUP",
	NIP:          "NIP",
	OVER:         "OVER",
	PICK:         "PICK",
	ROLL:         "ROLL",
	ROT:          "ROT",
	SWAP:         "SWAP",
	TUCK:         "TUCK",
	TWODROP:      "TWODROP",
	TWODUP:       "TWODUP",
	THREEDUP:     "THREEDUP",
	TWOOVER:      "TWOOVER",
	TWOROT:       "TWOROT",
	TWOSWAP:      "TWOSWAP",

	//Splice
	CAT:    "CAT",
	SUBSTR: "SUBSTR",
	LEFT:   "LEFT",
	RIGHT:  "RIGHT",
	SIZE:   "SIZE",

	//Bitwise
	INVERT:      "INVERT",
	AND:         "AND",
	OR:          "OR",
	XOR:         "XOR",
	EQUAL:       "EQUAL",
	EQUALVERIFY: "EQUALVERIFY",

	//Arithmetic
	ONEADD:            "ONEADD",
	ONESUB:            "ONESUB",
	TWOMUL:            "TWOMUL",
	TWODIV:            "TWODIV",
	NEGATE:            "NEGATE",
	ABS:               "ABS",
	NOT:               "NOT",
	ZERONOTEQUAL:      "ZERONOTEQUAL",
	ADD:               "ADD",
	SUB:               "SUB",
	MUL:               "MUL",
	DIV:               "DIV",
	MOD:               "MOD",
	LSHIFT:            "LSHIFT",
	RSHIFT:            "RSHIFT",
	BOOLAND:           "BOOLAND",
	BOOLOR:            "BOOLOR",
	NUMEQUAL:          "NUMEQUAL",
	NUMEQUALVERIFY:    "NUMEQUALVERIFY",
	NUMNOTEQUAL:       "NUMNOTEQUAL",
	LESSTHAN:          "LESSTHAN",
	GREATERTHAN:       "GREATERTHAN",
	LESSTHANOREQUAL:   "LESSTHANOREQUAL",
	GREATETHANOREQUAL: "GREATETHANOREQUAL",
	MIN:               "MIN",
	MAX:               "MAX",
	WITHIN:            "WITHIN",

	//Crypto
	RIPEMD160:           "RIPEMD160",
	SHA1:                "SHA1",
	SHA256:              "SHA256",
	HASH160:             "HASH160",
	HASH256:             "HASH256",
	CHECKSIG:            "CHECKSIG",
	CHECKSIGVERIFY:      "CHECKSIGVERIFY",
	CHECKMULTISIG:       "CHECKMULTISIG",
	CHECKMULTISIGVERIFY: "CHECKMULTISIGVERIFY",

	//Pseudo
	PUBKEYHASH: "PUBKEYHASH",
	PUBKEY:     "PUBKEY",
	INVALID:    "INVALID",

	//Reserved
	RESERVED:  "RESERVED",
	VER:       "VER",
	VERIF:     "VERIF",
	VERNOTIF:  "VERNOTIF",
	RESERVED1: "RESERVED1",
	RESERVED2: "RESERVED2",
	NOP1:      "NOP1",
	NOP2:      "NOP2",
	NOP3:      "NOP3",
	NOP4:      "NOP4",
	NOP5:      "NOP5",
	NOP6:      "NOP6",
	NOP7:      "NOP7",
	NOP8:      "NOP8",
	NOP9:      "NOP9",
	NOP10:     "NOP10",
}

func (this OpCode) IsFlow() bool {
	return this == NOP ||
		this == VERIFY ||
		this == RETURN
}

func (this OpCode) IsStack() bool {
	return this == FROMALTSTACK ||
		this == IFDUP ||
		this == DEPTH ||
		this == DROP ||
		this == DUP ||
		this == NIP ||
		this == OVER ||
		this == PICK ||
		this == ROLL ||
		this == ROT ||
		this == SWAP ||
		this == TUCK ||
		this == TWODROP ||
		this == TWODUP ||
		this == THREEDUP ||
		this == TWOOVER ||
		this == TWOROT ||
		this == TWOSWAP
}

func (this OpCode) IsSplice() bool {

	return this == CAT ||
		this == SUBSTR ||
		this == LEFT ||
		this == RIGHT ||
		this == SIZE
}

func (this OpCode) IsBitwise() bool {

	return this == INVERT ||
		this == AND ||
		this == OR ||
		this == XOR ||
		this == EQUAL ||
		this == EQUALVERIFY
}

func (this OpCode) IsArithmetic() bool {

	return this == ONEADD ||
		this == ONESUB ||
		this == TWOMUL ||
		this == TWODIV ||
		this == NEGATE ||
		this == ABS ||
		this == ZERONOTEQUAL ||
		this == ADD ||
		this == SUB ||
		this == MUL ||
		this == DIV ||
		this == MOD ||
		this == LSHIFT ||
		this == RSHIFT ||
		this == BOOLAND ||
		this == BOOLOR ||
		this == NUMEQUAL ||
		this == NUMEQUALVERIFY ||
		this == NUMNOTEQUAL ||
		this == LESSTHAN ||
		this == GREATERTHAN ||
		this == LESSTHANOREQUAL ||
		this == GREATETHANOREQUAL ||
		this == MIN ||
		this == MAX ||
		this == WITHIN
}

func (this OpCode) IsCrypto() bool {
	return this == RIPEMD160 ||
		this == SHA1 ||
		this == SHA256 ||
		this == HASH160 ||
		this == HASH256 ||
		this == CODESEPARATOR ||
		this == CHECKSIG ||
		this == CHECKSIGVERIFY ||
		this == CHECKMULTISIG ||
		this == CHECKMULTISIGVERIFY
}

func (this OpCode) IsOpCode() bool {
	return this.IsFlow() ||
		this.IsStack() ||
		this.IsSplice() ||
		this.IsBitwise() ||
		this.IsArithmetic() ||
		this.IsCrypto()
}

func (this OpCode) IsPsuedo() bool {
	return this == PUBKEYHASH ||
		this == PUBKEY ||
		this == INVALID
}

func (this OpCode) IsReserved() bool {
	return this == RESERVED ||
		this == VER ||
		this == VERIF ||
		this == VERNOTIF ||
		this == RESERVED1 ||
		this == RESERVED2
}

func (this OpCode) IsNop() bool {
	return this == NOP ||
		this == NOP1 ||
		this == NOP2 ||
		this == NOP3 ||
		this == NOP4 ||
		this == NOP5 ||
		this == NOP6 ||
		this == NOP7 ||
		this == NOP8 ||
		this == NOP9 ||
		this == NOP10
}

func (this OpCode) IsValid() bool {
	//Psuedo MUST not be in this list
	return this.IsOpCode() || this.IsReserved() || this.IsNop()
}

func (this OpCode) String() string {
	if this.IsValid() || this.IsPsuedo() {
		return opcodes[this]
	}
	return "UNKNOWN"
}
