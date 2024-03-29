package src

import (
	"fmt"
	"math"
	"strconv"

	"github.com/smekuria1/goclox/globals"
)

var scanner Scanner

// Parser is responsible for parsing the source code and generating the abstract syntax tree (AST).
type Parser struct {
	// Current represents the current token being processed.
	Current Token

	// Previous represents the previously processed token.
	Previous Token

	// HadError indicates whether an Error occurred during parsing.
	HadError bool

	// PanicMode indicates whether the parser is in panic mode.
	PanicMode bool
}
type Precedence int

// Precedence represents the precedence level of a token.
const (
	PrecNONE Precedence = iota
	PrecASSIGNMENT
	PrecOR
	PrecAND
	PrecEQUALITY
	PrecCOMPARISON
	PrecTERM
	PrecFACTOR
	PrecUNAR
	PrecCALL
	PrecPRIMARY
)

// Uint8Count represents the maximum number of local variables.
const Uint8Count = StackMax

var rules map[globals.TokenType]ParseRule

// ParseRule represents the parsing rule for a specific token type.
type ParseRule struct {
	// Prefix is the parsing function for an expression with the token as a prefix.
	Prefix Parsefn

	// Infix is the parsing function for an expression with the token as an infix.
	Infix Parsefn

	// Precedence represents the precedence level of the token.
	Precedence Precedence
}

// FunctionType represents the type of a function.
type FunctionType int

const (
	// FunctionTypeNative represents a native function.
	TypeFunction FunctionType = iota
	TypeScript
)

// Compiler represents a compiler object.
type Compiler struct {
	locals     [Uint8Count]Local // An array of `Local` objects with a length of `Uint8Count`.
	localCount int               // Keeps track of the number of local variables.
	scopeDepth int               // Represents the depth of the current scope.
	function   *ObjFunction      // Represents the current function being compiled.
	funcType   FunctionType      // Represents the type of the current function being compiled.
	encolsing  *Compiler         // Represents the compiler that encloses the current compiler.
}

// Local represents a local variable in the compiler.
type Local struct {
	name  Token // The name of the local variable.
	depth int   // The depth of the local variable within the scope.
}

// Parsefn represents the parsing function for a specific token type.
type Parsefn func(canAssign bool)

var parser Parser

// var compilingChunk *Chunk
var current *Compiler = nil

// InitCompiler initializes the compiler.
//
// It takes a pointer to a Compiler struct as a parameter.
func InitCompiler(compiler *Compiler, _type FunctionType) {
	compiler.encolsing = current
	compiler.function = nil
	compiler.funcType = _type
	compiler.localCount = 0
	compiler.scopeDepth = 0
	compiler.function = NewFunction()
	current = compiler
	if _type != TypeScript {
		current.function.name = copyString(parser.Previous.Start, parser.Previous.Length, ObjStringType)
	}
	local := &current.locals[current.localCount]
	local.depth = 0
	local.name.Start = 0
	local.name.Length = 0
}

func currentChunk() *Chunk {
	return &current.function.chunk
}

// Compile compiles the source code into an ObjFunction.
//
// It takes a source string and a pointer to a Chunk as parameters and returns a pointer to an ObjFunction.
func Compile(source string, chunk *Chunk) *ObjFunction {

	scanner.InitScanner(source)
	var compiler Compiler
	InitCompiler(&compiler, TypeScript)
	// compilingChunk = chunk
	parser.HadError = false
	parser.PanicMode = false
	advance(*scanner.Source)

	// for i := 0; i < scanner.Line; i+=1 {
	// 	expression()
	// }
	// consume(globals.TokenEOF, "Expect end of expression.")
	for !match(globals.TokenEOF) {
		declaration()
	}
	function := endCompiler()

	if parser.HadError {
		return nil
	}
	return function

}

// expression is a Go function that parses the precedence of an assignment.
//
// It does not take any parameters.
// It does not return any values.
func expression() {
	parsePrecendece(PrecASSIGNMENT)
}

// declaration represents a Go function that handles the declaration of a variable or a statement.
//
// It does not take any parameters and does not return any values.
func declaration() {
	if match(globals.TokenFUN) {
		functionDeclaration()
	} else if match(globals.TokenVAR) {
		varDeclaration()
	} else {
		statement()
	}
	if parser.PanicMode {
		synchronize()
	}
}

func functionDeclaration() {
	global := parseVariable("Expect function name.")
	markInitialized()
	function(TypeFunction)
	defineVariable(global)
}

func markInitialized() {
	if current.scopeDepth == 0 {
		return
	}
	current.locals[current.localCount-1].depth = current.scopeDepth
}

func function(_type FunctionType) {
	var compiler Compiler
	InitCompiler(&compiler, _type)
	beginScope()

	consume(globals.TokenLeftParen, "Expect '(' after function name.")
	if !check(globals.TokenRightParen) {
		for {
			current.function.arity++
			if current.function.arity > 255 {
				errorAtCurrent("Can't have more than 255 parameters.")
			}
			paramConstant := parseVariable("Expect parameter name.")
			defineVariable(paramConstant)
			if !match(globals.TokenCOMMA) {
				break
			}
		}
	}
	consume(globals.TokenRightParen, "Expect ')' after parameters.")

	consume(globals.TokenLeftBrace, "Expect '{' before function body.")
	block()

	function := endCompiler()
	emityBytes(uint8(globals.OpConstant), makeConstant(ObjVal(function)))
}

// varDeclaration is a function that performs variable declaration.
//
// It takes no parameters and does not return anything.
func varDeclaration() {
	global := parseVariable("Expect variable name. ")
	if match(globals.TokenEQUAL) {
		expression()
	} else {
		emitByte(uint8(globals.OpNil))
	}
	consume(globals.TokenSEMICOLON, "Expect ';' after variable declaration.")
	defineVariable(global)
}

// parseVariable parses the variable and returns a uint8 value.
//
// It takes an `errorMessage` string as a parameter.
// The function consumes the `globals.TokenIDENTIFIER` and `errorMessage`.
// It then declares a variable and checks the `current.scopeDepth`.
// If the `current.scopeDepth` is greater than 0, it returns 0.
// Otherwise, it returns the identifier constant of `parser.Previous`.
func parseVariable(errorMessage string) uint8 {
	consume(globals.TokenIDENTIFIER, errorMessage)
	declareVariable()
	if current.scopeDepth > 0 {
		return 0
	}
	return identifierConstant(&parser.Previous)
}

// declareVariable is a function that declares a variable.
//
// This function checks if the current scope depth is 0 and returns if it is.
// It then assigns the value of the previous parser to the 'name' variable.
// The function iterates over the local variables in reverse order, checking
// if there is already a variable with the same name in the current scope.
// If there is, it raises an Error. Otherwise, it adds the 'name' to the
// list of local variables.
func declareVariable() {
	if current.scopeDepth == 0 {
		return
	}
	name := &parser.Previous
	for i := current.localCount - 1; i >= 0; i-- {
		local := &current.locals[i]
		if local.depth != -1 && local.depth < current.scopeDepth {
			break
		}
		if identfierEqual(name, &local.name) {
			Error("Already variable with this name in this scope")
		}
	}
	addLocal(name)
}

// identfierEqual checks if two tokens have the same identifier.
//
// Parameters:
// - a: pointer to the first token
// - b: pointer to the second token
//
// Returns:
// - bool: true if the tokens have the same identifier, false otherwise
func identfierEqual(a, b *Token) bool {
	if a.Length != b.Length {
		return false
	}
	source := *scanner.Source
	aChar := source[a.Start : a.Start+a.Length]
	bChar := source[b.Start : b.Start+b.Length]

	return memcmp([]byte(aChar), []byte(bChar), a.Length) == 0
}

// addLocal adds a local variable to the current function.
//
// Parameters:
// - name: a pointer to a Token representing the name of the variable.
//
// Returns: None.
func addLocal(name *Token) {
	if current.localCount == Uint8Count {
		Error("Too many local variables in function")
		return
	}

	current.locals[current.localCount].name = *name
	current.locals[current.localCount].depth = current.scopeDepth
	current.localCount++
}

// identifierConstant generates a constant identifier.
//
// Parameters:
// - name: a pointer to a Token representing the name of the identifier.
//
// Returns:
// - uint8: the generated constant identifier.
func identifierConstant(name *Token) uint8 {
	return makeConstant(ObjStrValue(copyString(name.Start, name.Length, ObjStringType)))
}

// defineVariable defines a global variable.
//
// The function takes a single parameter, `global`, which is of type `uint8`.
// It does not return any values.
func defineVariable(global uint8) {
	if current.scopeDepth > 0 {
		markInitialized()
		return
	}
	emityBytes(uint8(globals.OpDefineGlobal), global)
}

// statement is a Go function that performs a specific task based on the current token.
//
// It checks if the current token matches the TokenPRINT and calls the printStatement function if it does.
// If the current token matches the TokenLeftBrace, it calls the beginScope, block, and endScope functions to handle a block of code.
// If the current token does not match any of the above, it calls the expressionStatement function.
func statement() {
	if match(globals.TokenPRINT) {
		printStatement()
	} else if match(globals.TokenLeftBrace) {
		beginScope()
		block()
		endScope()
	} else if match(globals.TokenIF) {
		ifStatement()
	} else if match(globals.TokenFOR) {
		forStatement()
	} else if match(globals.TokenWHILE) {
		whileStatement()
	} else {
		expressionStatement()
	}
}

// forStatement is a function that processes the for loop
func forStatement() {
	beginScope()
	consume(globals.TokenLeftParen, "Expect '(' after 'for'.")
	if match(globals.TokenSEMICOLON) {
		//No initialier
	} else if match(globals.TokenVAR) {
		varDeclaration()
	} else {
		expressionStatement()
	}

	loopStart := currentChunk().Count
	exitJump := -1
	if !match(globals.TokenSEMICOLON) {
		expression()
		consume(globals.TokenSEMICOLON, "Expect ';' after loop condition")

		// jump out of loop
		exitJump = int(emitJump(uint8(globals.OpJumpFalse)))
		emitByte(uint8(globals.OpPop))
	}

	if !match(globals.TokenRightParen) {
		bodyJump := emitJump(uint8(globals.OpJump))

		incrementStart := currentChunk().Count
		expression()
		emitByte(uint8(globals.OpPop))
		consume(globals.TokenRightParen, "Expect ')' after for clauses.")

		emitLoop(loopStart)
		loopStart = incrementStart
		patchJump(bodyJump)
	}

	statement()
	emitLoop(loopStart)
	if exitJump != -1 {
		patchJump(uint32(exitJump))
		emitByte(uint8(globals.OpPop))
	}
	endScope()
}

// whileStatement is a function that processes the while loop
func whileStatement() {
	loopStart := currentChunk().Count
	consume(globals.TokenLeftParen, "Expect '(' after 'while'")
	expression()
	consume(globals.TokenRightParen, "Expect ')' after condition")

	exitJump := emitJump(uint8(globals.OpJumpFalse))

	emitByte(uint8(globals.OpPop))
	statement()

	emitLoop(loopStart)

	patchJump(exitJump)
	emitByte(uint8(globals.OpPop))

}

// ifStatement is a function that processes an if statement.
func ifStatement() {
	consume(globals.TokenLeftParen, "Expect '(' after 'if'")
	expression()
	consume(globals.TokenRightParen, "Expect ')' after condition")

	thenJump := emitJump(uint8(globals.OpJumpFalse))
	emitByte(uint8(globals.OpPop))
	statement()

	elseJump := emitJump(uint8(globals.OpJump))
	patchJump(thenJump)

	emitByte(uint8(globals.OpPop))
	if match(globals.TokenELSE) {
		statement()
	}
	patchJump(elseJump)

}

func emitLoop(loopStart int) {
	emitByte(uint8(globals.OpLoop))

	offset := currentChunk().Count - loopStart + 2
	if offset > math.MaxUint16 {
		Error("Loop body too large")
	}
	emitByte(uint8((offset >> 8) & 0xff))
	emitByte(uint8(offset & 0xff))
}

// emitJump is a function that emits a jump instruction.
func emitJump(opCode uint8) uint32 {
	emitByte(opCode)
	emitByte(0xff)
	emitByte(0xff)
	return uint32(currentChunk().Count - 2)
}

// patchJump is a function that patches a jump instruction.
func patchJump(offset uint32) {
	jump := currentChunk().Count - int(offset) - 2

	if jump > math.MaxUint16 {
		Error("Too much code to jump over")
	}
	currentChunk().Code[offset] = uint8((jump >> 8) & 0xff)
	currentChunk().Code[offset+1] = uint8(jump & 0xff)
}

// block is a function that processes a block of code.
//
// It iterates through the code until either a right brace token or
// an end of file token is encountered. During each iteration, it
// calls the declaration function to process the code.
//
// There are no parameters.
//
// The function does not return anything.
func block() {
	for !check(globals.TokenRightBrace) && !check(globals.TokenEOF) {
		declaration()
	}

	consume(globals.TokenRightBrace, "Expect '}' after block")
}

// beginScope increments the scope depth.
//
// No parameters.
// No return types.
func beginScope() {
	current.scopeDepth++
}

// endScope decrements the scope depth and pops local variables until the
// current scope depth is reached.
//
// No parameters.
// No return type.
func endScope() {
	current.scopeDepth--
	for current.localCount > 0 && current.locals[current.localCount-1].depth > current.scopeDepth {
		emitByte(uint8(globals.OpPop))
		current.localCount--
	}
}

// expressionStatement description of the Go function.
//
// This function does not have any parameters.
// It does not return anything.
func expressionStatement() {
	expression()
	consume(globals.TokenSEMICOLON, "Expext ';' after expression")
	emitByte(uint8(globals.OpPop))
}

// match checks if the given token type matches the current token and advances the scanner.
//
// _type: the token type to match
// Returns: true if the token type matches and the scanner has been advanced, false otherwise
func match(_type globals.TokenType) bool {
	if !check(_type) {
		return false
	}
	advance(*scanner.Source)
	return true
}

// check checks if the current token type matches the given token type.
//
// _type: the token type to check against.
// bool: true if the current token type matches the given token type, false otherwise.
func check(_type globals.TokenType) bool {
	return parser.Current.TOKENType == _type
}

// printStatement prints a statement.
//
// This function calls the expression function, then the consume function, and finally the emitByte function.
// It doesn't take any parameters and doesn't return any values.
func printStatement() {
	expression()
	consume(globals.TokenSEMICOLON, "Expect ';' after value.")
	emitByte(uint8(globals.OpPrint))
}

// synchronize is a Go function that synchronizes the parser state.
//
// It iterates through the tokens in the parser until it reaches the end of the file or encounters a semicolon.
// During the iteration, it checks the type of each token and performs specific actions based on the type.
//
// No parameters are required for this function.
// This function does not return any values.
func synchronize() {
	parser.PanicMode = false
	for parser.Current.TOKENType != globals.TokenEOF {
		if parser.Previous.TOKENType == globals.TokenSEMICOLON {
			return
		}
		switch parser.Current.TOKENType {
		case globals.TokenCLASS:
		case globals.TokenFUN:
		case globals.TokenVAR:
		case globals.TokenFOR:
		case globals.TokenIF:
		case globals.TokenWHILE:
		case globals.TokenPRINT:
		case globals.TokenRETURN:
			return
		default:
			// Do nothing.
		}
		advance(*scanner.Source)
	}
}

// endCompiler is a Go function that is responsible for ending the compiler process.
//
// This function does not take any parameters.
// It does not return any values.
func endCompiler() *ObjFunction {
	emitReturn()
	function := current.function
	if globals.DEBUG_PRINT_CODE {
		if !parser.HadError {
			if function.name != nil {
				DisassembleChunk(currentChunk(), string(function.name.Chars))
			} else {
				DisassembleChunk(currentChunk(), "script")
			}
		}
	}
	current = current.encolsing
	return function

}

// getRule returns the ParseRule associated with the given TokenType.
//
// It takes a TokenType as input parameter and returns a pointer to a ParseRule.
func getRule(tokentype globals.TokenType) *ParseRule {
	parserule := rules[tokentype]
	return &parserule
}

func call(canAssign bool) {
	argcount := argumentList()
	emityBytes(uint8(globals.OpCall), argcount)
}

func argumentList() uint8 {
	argcount := uint8(0)
	if !check(globals.TokenRightParen) {
		for {
			expression()
			argcount++
			if argcount == 255 {
				errorAtCurrent("Can't have more than 255 arguments.")
			}
			if !match(globals.TokenCOMMA) {
				break
			}
		}

	}
	consume(globals.TokenRightParen, "Expect ')' after arguments.")

	return argcount
}

// binary represents a function that performs a binary operation based on the given operator type.
//
// It takes a boolean parameter canAssign, which indicates whether the operation can be assigned to a variable.
// This function does not return any value.
func binary(canAssign bool) {
	operatorType := parser.Previous.TOKENType
	rule := getRule(operatorType)
	parsePrecendece(rule.Precedence + 1)

	switch operatorType {
	case globals.TokenBANG_EQUAL:
		emityBytes(uint8(globals.OpEqual), uint8(globals.OpNot))
	case globals.TokenEQUAL_EQUAL:
		emitByte(uint8(globals.OpEqual))
	case globals.TokenGREATER:
		emitByte(uint8(globals.OpGreater))
	case globals.TokenGREATER_EQUAL:
		emityBytes(uint8(globals.OpLess), uint8(globals.OpNot))
	case globals.TokenLESS:
		emitByte(uint8(globals.OpLess))
	case globals.TokenLESS_EQUAL:
		emityBytes(uint8(globals.OpGreater), uint8(globals.OpNot))
	case globals.TokenPLUS:
		emitByte(uint8(globals.OpAdd))
	case globals.TokenMINUS:
		emitByte(uint8(globals.OpSubtract))
	case globals.TokenSTAR:
		emitByte(uint8(globals.OpMultiply))
	case globals.TokenSLASH:
		emitByte(uint8(globals.OpDivide))
	}
}

// literal generates bytecode for literal values.
//
// The function takes a boolean parameter `canAssign` which determines if the
// literal value can be assigned.
// It does not return any value.
func literal(canAssign bool) {
	switch parser.Previous.TOKENType {
	case globals.TokenFALSE:
		emitByte(uint8(globals.OpFalse))
	case globals.TokenNIL:
		emitByte(uint8(globals.OpNil))
	case globals.TokenTRUE:
		emitByte(uint8(globals.OpTrue))
	default:
		return
	}
}

// grouping is a Go function that performs some operation.
//
// It takes a boolean parameter, canAssign, which determines whether the function can perform an assignment operation.
//
// The function does not return any value.
func grouping(canAssign bool) {
	expression()
	consume(globals.TokenRightParen, "Expect ')' after the expression")
}

// emitReturn emits the return opcode.
//
// This function does not take any parameters.
// It does not return anything.
func emitReturn() {
	emitByte(uint8(globals.OpNil))
	emitByte(uint8(globals.OpReturn))
}

// number is a function that performs some operation on a given input.
//
// It takes a boolean argument canAssign, which determines whether the function can assign a value.
// The function does not return anything.
func number(canAssign bool) {
	source := *scanner.Source
	value, err := strconv.ParseFloat(source[parser.Previous.Start:parser.Previous.Start+parser.Previous.Length], 64)
	if err != nil {
		Error(err.Error())
	}
	emitConstant(NumberValue(value))

}

// stringy generates a string based on the given value of canAssign.
//
// It takes a boolean parameter canAssign which determines whether the generated string can be assigned or not.
// The function does not return any value.
func stringy(canAssign bool) {
	emitConstant(ObjStrValue(copyString(parser.Previous.Start+1, parser.Previous.Length-2, ObjStringType)))
}

// variable is a Go function that takes a boolean parameter canAssign.
// The function calls the namedVariable function passing parser.Previous and canAssign as arguments.
func variable(canAssign bool) {
	namedVariable(parser.Previous, canAssign)
}

// namedVariable is a function that takes a name Token and a canAssign boolean as parameters.
//
// The function resolves the local variable using the current scope and the name Token. If the
// variable is found in the current scope, it uses the OpGetLocal and OpSetLocal opcodes to get
// and set the variable value. If the variable is not found in the current scope, it uses the
// OpGetGlobal and OpSetGlobal opcodes to get and set the variable value. If the canAssign
// parameter is true and there is an EQUAL token, the function calls the expression() function and
// emits the set opcode and the argument. Otherwise, it emits the get opcode and the argument.
func namedVariable(name Token, canAssign bool) {
	var (
		getOp globals.OpCode
		setOp globals.OpCode
	)

	arg := resolveLocal(current, &name)
	if arg != -1 {
		getOp = globals.OpGetLocal
		setOp = globals.OpSetLocal
	} else {
		arg = int(identifierConstant(&name))
		getOp = globals.OpGetGlobal
		setOp = globals.OpSetGlobal
	}
	if match(globals.TokenEQUAL) && canAssign {
		expression()
		emityBytes(uint8(setOp), uint8(arg))
	} else {
		emityBytes(uint8(getOp), uint8(arg))
	}
}

// resolveLocal finds the index of a local variable in the compiler's local
// array with a given name.
//
// Parameters:
// - compiler: a pointer to the Compiler struct
// - name: a pointer to the Token struct representing the name of the variable
//
// Return:
// - int: the index of the local variable in the compiler's local array, or -1 if not found
func resolveLocal(compiler *Compiler, name *Token) int {
	for i := compiler.localCount - 1; i >= 0; i-- {
		local := &compiler.locals[i]
		if identfierEqual(name, &local.name) {
			return i
		}

	}
	return -1
}

// unary performs a unary operation based on the operator type.
//
// It takes a boolean parameter, canAssign, to determine if the unary operation can be assigned.
// The function does not return any values.
func unary(canAssign bool) {
	opratorType := parser.Previous.TOKENType

	parsePrecendece(PrecUNAR)

	switch opratorType {
	case globals.TokenMINUS:
		emitByte(uint8(globals.OpNegate))
	case globals.TokenBANG:
		emitByte(uint8(globals.OpNot))
	default:
		return
	}
}

func and(canAssign bool) {
	endJump := emitJump(uint8(globals.OpJumpFalse))
	emitByte(uint8(globals.OpPop))

	parsePrecendece(PrecAND)

	patchJump(endJump)
}

func or(canAssign bool) {
	elseJump := emitJump(uint8(globals.OpJumpFalse))
	endJump := emitJump(uint8(globals.OpJump))

	patchJump(elseJump)
	emitByte(uint8(globals.OpPop))

	parsePrecendece(PrecOR)
	patchJump(endJump)
}

// parsePrecendece parses the precedence of a given Precedence.
//
// It advances the scanner source and retrieves the prefix rule for the
// previous token type. If the prefix rule is nil, it throws an Error and
// returns. It determines if the precedence is less than or equal to
// PrecASSIGNMENT and assigns the result to canAssign. It then calls the
// prefix rule with the canAssign parameter.
//
// It enters a loop where it checks if the precedence is less than or equal
// to the precedence of the current token type. If true, it advances the
// source, retrieves the infix rule for the previous token type, and calls
// the infix rule with the canAssign parameter.
//
// After the loop, it checks if canAssign is true and if the current token
// type matches TokenEQUAL. If true, it throws an Error for invalid
// assignment target.
func parsePrecendece(precedence Precedence) {

	advance(*scanner.Source)
	prefixRule := getRule(parser.Previous.TOKENType).Prefix
	if prefixRule == nil {
		Error("Expect expression")
		return
	}
	canAssign := precedence <= PrecASSIGNMENT
	prefixRule(canAssign)

	for precedence <= getRule(parser.Current.TOKENType).Precedence {
		advSource := *scanner.Source
		advance(advSource[parser.Current.Start:])
		infixRule := getRule(parser.Previous.TOKENType).Infix
		infixRule(canAssign)
	}

	if canAssign && match(globals.TokenEQUAL) {
		Error("Invalid assignment target")
	}
}

// emitConstant generates a constant value for the Go function.
//
// It takes a value of type Value as a parameter.
// It does not return anything.
func emitConstant(value Value) {
	emityBytes(uint8(globals.OpConstant), makeConstant(value))
}

// makeConstant generates a new constant value in the current chunk.
//
// value: the value to be added as a constant.
// Returns: the index of the constant in the chunk.
func makeConstant(value Value) uint8 {
	constant := AddConstants(currentChunk(), value)
	if constant > StackMax {
		Error("Too many constants in one chunk")
		return 0
	}
	return uint8(constant)
}

// emityBytes emits two bytes of bytecode.
//
// The function takes two parameters, "bytecode1" and "bytecode2", both of type uint8.
// It does not return anything.
func emityBytes(bytecode1, bytecode2 uint8) {
	emitByte(bytecode1)
	emitByte(bytecode2)
}

// consume consumes a token of the given type and advances the parser.
//
// Parameters:
// - tokentype: the type of token to consume.
// - message: the Error message to display if the token type does not match.
//
// Return type: None.
func consume(tokentype globals.TokenType, message string) {
	if parser.Current.TOKENType == tokentype {
		source := *scanner.Source
		advance(source[parser.Current.Length:])
		return
	}

	Error(message)
}

// emitByte writes a bytecode to the compiling chunk.
//
// bytecode: the bytecode to be written.
// Returns: nothing.
func emitByte(bytecode uint8) {
	WriteChunk(currentChunk(), bytecode, parser.Previous.Line)
}

// advance advances the parser to the next token in the source string.
//
// Parameters:
// - source: a string representing the source code to be parsed.
//
// Return type: None.
func advance(source string) {
	parser.Previous = parser.Current

	for {
		parser.Current = scanner.ScanToken(&source)
		if parser.Current.TOKENType != globals.TokenERROR {
			break
		}

		errorAtCurrent(source[parser.Current.Start:])
	}
}

// errorAtCurrent is a function that takes a message as a parameter and calls the errorAt function with the parser.Current variable and the message. It does not return any value.
//
// Parameters:
// - message: a string representing the Error message.
func errorAtCurrent(message string) {
	errorAt(&parser.Current, message)
}

// Error is a function that handles errors and logs them.
//
// It takes a message string as a parameter and calls the errorAt function
// passing the address of the parser.Previous variable and the Error message.
func Error(message string) {
	errorAt(&parser.Previous, message)
}

// errorAt prints an Error message and sets the parser in panic mode.
//
// It takes a token pointer and a message string as parameters.
// It does not return anything.
func errorAt(token *Token, message string) {
	if parser.PanicMode {
		return
	}
	parser.PanicMode = true
	fmt.Printf("Error [line %d],", token.Line)
	source := *scanner.Source
	if token.TOKENType == globals.TokenEOF {
		fmt.Printf(" at end")
	} else if token.TOKENType == globals.TokenERROR {
		//
	} else {
		fmt.Printf(" at '%s'", string(source[token.Start:token.Start+token.Length]))
	}

	fmt.Printf(": %s\n", message)
	parser.HadError = true
}

// init initializes the rules map for the TokenType and corresponding ParseRule.
//
// No parameters.
// No return type.
func init() {
	rules = map[globals.TokenType]ParseRule{
		globals.TokenLeftParen:     {grouping, call, PrecNONE},
		globals.TokenRightParen:    {nil, nil, PrecNONE},
		globals.TokenLeftBrace:     {nil, nil, PrecNONE},
		globals.TokenRightBrace:    {nil, nil, PrecNONE},
		globals.TokenCOMMA:         {nil, nil, PrecNONE},
		globals.TokenDOT:           {nil, nil, PrecNONE},
		globals.TokenMINUS:         {unary, binary, PrecTERM},
		globals.TokenPLUS:          {nil, binary, PrecTERM},
		globals.TokenSEMICOLON:     {nil, nil, PrecNONE},
		globals.TokenSLASH:         {nil, binary, PrecFACTOR},
		globals.TokenSTAR:          {nil, binary, PrecFACTOR},
		globals.TokenBANG:          {unary, nil, PrecNONE},
		globals.TokenBANG_EQUAL:    {nil, binary, PrecEQUALITY},
		globals.TokenEQUAL:         {nil, nil, PrecNONE},
		globals.TokenEQUAL_EQUAL:   {nil, binary, PrecEQUALITY},
		globals.TokenGREATER:       {nil, binary, PrecCOMPARISON},
		globals.TokenGREATER_EQUAL: {nil, binary, PrecCOMPARISON},
		globals.TokenLESS:          {nil, binary, PrecCOMPARISON},
		globals.TokenLESS_EQUAL:    {nil, binary, PrecCOMPARISON},
		globals.TokenIDENTIFIER:    {variable, nil, PrecNONE},
		globals.TokenSTRING:        {stringy, nil, PrecNONE},
		globals.TokenNUMBER:        {number, nil, PrecNONE},
		globals.TokenAND:           {nil, and, PrecNONE},
		globals.TokenCLASS:         {nil, nil, PrecNONE},
		globals.TokenELSE:          {nil, nil, PrecNONE},
		globals.TokenFALSE:         {literal, nil, PrecNONE},
		globals.TokenFOR:           {nil, nil, PrecNONE},
		globals.TokenFUN:           {nil, nil, PrecNONE},
		globals.TokenIF:            {nil, nil, PrecNONE},
		globals.TokenNIL:           {literal, nil, PrecNONE},
		globals.TokenOR:            {nil, or, PrecNONE},
		globals.TokenPRINT:         {nil, nil, PrecNONE},
		globals.TokenRETURN:        {nil, nil, PrecNONE},
		globals.TokenSUPER:         {nil, nil, PrecNONE},
		globals.TokenTHIS:          {nil, nil, PrecNONE},
		globals.TokenTRUE:          {literal, nil, PrecNONE},
		globals.TokenVAR:           {nil, nil, PrecNONE},
		globals.TokenWHILE:         {nil, nil, PrecNONE},
		globals.TokenERROR:         {nil, nil, PrecNONE},
		globals.TokenEOF:           {nil, nil, PrecNONE},
	}
}
