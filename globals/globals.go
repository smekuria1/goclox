package globals

type OpCode uint8

const (
	OpReturn OpCode = iota
	OpNegate
	OpPrint
	OpPop
	OpDefineGlobal
	OpGetGlobal
	OpSetGlobal
	OpSetLocal
	OpGetLocal
	OpNil
	OpTrue
	OpFalse
	OpEqual
	OpGreater
	OpLess
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNot
	OpConstant
)

type TokenType int

const (
	// Single-character tokens.
	TokenLeftParen TokenType = iota
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenCOMMA
	TokenDOT
	TokenMINUS
	TokenPLUS
	TokenSEMICOLON
	TokenSLASH
	TokenSTAR

	TokenBANG
	TokenBANG_EQUAL
	TokenEQUAL
	TokenEQUAL_EQUAL
	TokenGREATER
	TokenGREATER_EQUAL
	TokenLESS
	TokenLESS_EQUAL

	TokenIDENTIFIER
	TokenSTRING
	TokenNUMBER

	TokenAND
	TokenCLASS
	TokenELSE
	TokenFALSE
	TokenFOR
	TokenFUN
	TokenIF
	TokenNIL
	TokenOR
	TokenPRINT
	TokenRETURN
	TokenSUPER
	TokenTHIS
	TokenTRUE
	TokenVAR
	TokenWHILE

	TokenERROR
	TokenEOF
)

var DEBUG_TRACE_EXECUTION = false
var DEBUG_PRINT_CODE = false
