package src

import (
	"fmt"
	"strconv"

	"github.com/smekuria1/goclox/globals"
)

var scanner Scanner

type Parser struct {
	Current   Token
	Previous  Token
	HadError  bool
	PanicMode bool
}
type Precedence int

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNAR                  // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

var rules map[globals.TokenType]ParseRule

type ParseRule struct {
	Prefix     Parsefn
	Infix      Parsefn
	Precedence Precedence
}

type Parsefn func(canAssign bool)

var parser Parser

func Compile(source string, chunk *Chunk) bool {

	scanner.InitScanner(source)
	compilingChunk = chunk
	parser.HadError = false
	parser.PanicMode = false
	advance(*scanner.Source)

	// for i := 0; i < scanner.Line; i++ {
	// 	expression()
	// }
	// consume(globals.TOKEN_EOF, "Expect end of expression.")
	for !match(globals.TOKEN_EOF) {
		declaration()
	}
	endCompiler()
	return !parser.HadError

}

func expression() {
	parsePrecendece(PREC_ASSIGNMENT)
}

func declaration() {
	if match(globals.TOKEN_VAR) {
		varDeclaration()
	} else {
		statement()
	}
	if parser.PanicMode {
		synchronize()
	}
}

func varDeclaration() {
	global := parseVariable("Expect variable name. ")
	if match(globals.TOKEN_EQUAL) {
		expression()
	} else {
		emitByte(uint8(globals.OP_NIL))
	}
	consume(globals.TOKEN_SEMICOLON, "Expect ';' after variable declaration.")
	defineVariable(global)
}

func parseVariable(errorMessage string) uint8 {
	consume(globals.TOKEN_IDENTIFIER, errorMessage)
	return identifierConstant(&parser.Previous)
}

func identifierConstant(name *Token) uint8 {
	return makeConstant(ObjStrValue(copyString(name.Start, name.Length, *scanner.Source, ObjStringType)))
}

func defineVariable(global uint8) {
	emityBytes(uint8(globals.OP_DEFINE_GLOBAL), global)
}

func statement() {
	if match(globals.TOKEN_PRINT) {
		printStatement()
	} else {
		expressionStatement()
	}
}

func expressionStatement() {
	expression()
	consume(globals.TOKEN_SEMICOLON, "Expext ';' after expression")
	emitByte(uint8(globals.OP_POP))
}
func match(_type globals.TokenType) bool {
	if !check(_type) {
		return false
	}
	advance(*scanner.Source)
	return true
}

func check(_type globals.TokenType) bool {
	return parser.Current.TOKENType == _type
}

func printStatement() {
	expression()
	consume(globals.TOKEN_SEMICOLON, "Expect ';' after value.")
	emitByte(uint8(globals.OP_PRINT))
}

func synchronize() {
	parser.PanicMode = false
	for parser.Current.TOKENType != globals.TOKEN_EOF {
		if parser.Previous.TOKENType == globals.TOKEN_SEMICOLON {
			return
		}
		switch parser.Current.TOKENType {
		case globals.TOKEN_CLASS:
		case globals.TOKEN_FUN:
		case globals.TOKEN_VAR:
		case globals.TOKEN_FOR:
		case globals.TOKEN_IF:
		case globals.TOKEN_WHILE:
		case globals.TOKEN_PRINT:
		case globals.TOKEN_RETURN:
			return
		default:
			// Do nothing.
		}
		advance(*scanner.Source)
	}
}

func endCompiler() {
	if globals.DEBUG_PRINT_CODE {
		if !parser.HadError {
			DisassembleChunk(currentChunk(), "code")
		}
	}
	emitReturn()

}
func getRule(tokentype globals.TokenType) *ParseRule {
	parserule := rules[tokentype]
	return &parserule
}
func binary(canAssign bool) {
	operatorType := parser.Previous.TOKENType
	rule := getRule(operatorType)
	parsePrecendece(rule.Precedence + 1)

	switch operatorType {
	case globals.TOKEN_BANG_EQUAL:
		emityBytes(uint8(globals.OP_EQUAL), uint8(globals.OP_NOT))
	case globals.TOKEN_EQUAL_EQUAL:
		emitByte(uint8(globals.OP_EQUAL))
	case globals.TOKEN_GREATER:
		emitByte(uint8(globals.OP_GREATER))
	case globals.TOKEN_GREATER_EQUAL:
		emityBytes(uint8(globals.OP_LESS), uint8(globals.OP_NOT))
	case globals.TOKEN_LESS:
		emitByte(uint8(globals.OP_LESS))
	case globals.TOKEN_LESS_EQUAL:
		emityBytes(uint8(globals.OP_GREATER), uint8(globals.OP_NOT))
	case globals.TOKEN_PLUS:
		emitByte(uint8(globals.OP_ADD))
	case globals.TOKEN_MINUS:
		emitByte(uint8(globals.OP_SUBTRACT))
	case globals.TOKEN_STAR:
		emitByte(uint8(globals.OP_MULTIPLY))
	case globals.TOKEN_SLASH:
		emitByte(uint8(globals.OP_DIVIDE))
	}
}

func literal(canAssign bool) {
	switch parser.Previous.TOKENType {
	case globals.TOKEN_FALSE:
		emitByte(uint8(globals.OP_FALSE))
	case globals.TOKEN_NIL:
		emitByte(uint8(globals.OP_NIL))
	case globals.TOKEN_TRUE:
		emitByte(uint8(globals.OP_TRUE))
	default:
		return
	}
}

func grouping(canAssign bool) {
	expression()
	consume(globals.TOKEN_RIGHT_PAREN, "Expect ')' after the expression")
}
func emitReturn() {
	emitByte(uint8(globals.OP_RETURN))
}

func number(canAssign bool) {
	source := *scanner.Source
	value, err := strconv.ParseFloat(source[parser.Previous.Start:parser.Previous.Start+parser.Previous.Length], 64)
	if err != nil {
		error(err.Error())
	}
	emitConstant(NumberValue(value))

}
func stringy(canAssign bool) {
	source := *scanner.Source
	emitConstant(ObjStrValue(copyString(parser.Previous.Start+1, parser.Previous.Length-2, source, ObjStringType)))
}
func variable(canAssign bool) {
	namedVariable(parser.Previous, canAssign)
}

func namedVariable(name Token, canAssign bool) {
	arg := identifierConstant(&name)
	if match(globals.TOKEN_EQUAL) && canAssign {
		expression()
		emityBytes(uint8(globals.OP_SET_GLOBAL), arg)
	} else {
		emityBytes(uint8(globals.OP_GET_GLOBAL), arg)
	}
}
func unary(canAssign bool) {
	opratorType := parser.Previous.TOKENType

	parsePrecendece(PREC_UNAR)

	switch opratorType {
	case globals.TOKEN_MINUS:
		emitByte(uint8(globals.OP_NEGATE))
	case globals.TOKEN_BANG:
		emitByte(uint8(globals.OP_NOT))
	default:
		return
	}
}

func parsePrecendece(precedence Precedence) {

	advance(*scanner.Source)
	prefixRule := getRule(parser.Previous.TOKENType).Prefix
	if prefixRule == nil {
		error("Expect expression")
		return
	}
	canAssign := precedence <= PREC_ASSIGNMENT
	prefixRule(canAssign)

	for precedence <= getRule(parser.Current.TOKENType).Precedence {
		advSource := *scanner.Source
		advance(advSource[parser.Current.Start:])
		infixRule := getRule(parser.Previous.TOKENType).Infix
		infixRule(canAssign)
	}

	if canAssign && match(globals.TOKEN_EQUAL) {
		error("Invalid assignment target")
	}
}

func emitConstant(value Value) {
	emityBytes(uint8(globals.OP_CONSTANT), makeConstant(value))
}
func makeConstant(value Value) uint8 {
	constant := AddConstants(currentChunk(), value)
	if constant > STACK_MAX {
		error("Too many constants in one chunk")
		return 0
	}
	return uint8(constant)
}

func emityBytes(bytecode1, bytecode2 uint8) {
	emitByte(bytecode1)
	emitByte(bytecode2)
}

func consume(tokentype globals.TokenType, message string) {
	if parser.Current.TOKENType == tokentype {
		source := *scanner.Source
		advance(source[parser.Current.Length:])
		return
	}

	errorAtCurrent(message)
}

func emitByte(bytecode uint8) {
	WriteChunk(compilingChunk, bytecode, parser.Previous.Line)
}

var compilingChunk *Chunk

func currentChunk() *Chunk {
	return compilingChunk
}

func advance(source string) {
	parser.Previous = parser.Current

	for {
		parser.Current = scanner.ScanToken(&source)
		if parser.Current.TOKENType != globals.TOKEN_ERROR {
			break
		}

		errorAtCurrent(source[parser.Current.Start:])
	}
}

func errorAtCurrent(message string) {
	errorAt(&parser.Current, message)
}

func error(message string) {
	errorAt(&parser.Previous, message)
}

func errorAt(token *Token, message string) {
	if parser.PanicMode {
		return
	}
	parser.PanicMode = true
	fmt.Printf("Error [line %d],", token.Line)
	source := *scanner.Source
	if token.TOKENType == globals.TOKEN_EOF {
		fmt.Printf(" at end")
	} else if token.TOKENType == globals.TOKEN_ERROR {
		//
	} else {
		fmt.Printf(" at '%s'", string(source[token.Start:token.Start+token.Length]))
	}

	fmt.Printf(": %s\n", message)
	parser.HadError = true
}

func init() {
	rules = map[globals.TokenType]ParseRule{
		globals.TOKEN_LEFT_PAREN:    {grouping, nil, PREC_NONE},
		globals.TOKEN_RIGHT_PAREN:   {nil, nil, PREC_NONE},
		globals.TOKEN_LEFT_BRACE:    {nil, nil, PREC_NONE},
		globals.TOKEN_RIGHT_BRACE:   {nil, nil, PREC_NONE},
		globals.TOKEN_COMMA:         {nil, nil, PREC_NONE},
		globals.TOKEN_DOT:           {nil, nil, PREC_NONE},
		globals.TOKEN_MINUS:         {unary, binary, PREC_TERM},
		globals.TOKEN_PLUS:          {nil, binary, PREC_TERM},
		globals.TOKEN_SEMICOLON:     {nil, nil, PREC_NONE},
		globals.TOKEN_SLASH:         {nil, binary, PREC_FACTOR},
		globals.TOKEN_STAR:          {nil, binary, PREC_FACTOR},
		globals.TOKEN_BANG:          {unary, nil, PREC_NONE},
		globals.TOKEN_BANG_EQUAL:    {nil, binary, PREC_EQUALITY},
		globals.TOKEN_EQUAL:         {nil, nil, PREC_NONE},
		globals.TOKEN_EQUAL_EQUAL:   {nil, binary, PREC_EQUALITY},
		globals.TOKEN_GREATER:       {nil, binary, PREC_COMPARISON},
		globals.TOKEN_GREATER_EQUAL: {nil, binary, PREC_COMPARISON},
		globals.TOKEN_LESS:          {nil, binary, PREC_COMPARISON},
		globals.TOKEN_LESS_EQUAL:    {nil, binary, PREC_COMPARISON},
		globals.TOKEN_IDENTIFIER:    {variable, nil, PREC_NONE},
		globals.TOKEN_STRING:        {stringy, nil, PREC_NONE},
		globals.TOKEN_NUMBER:        {number, nil, PREC_NONE},
		globals.TOKEN_AND:           {nil, nil, PREC_NONE},
		globals.TOKEN_CLASS:         {nil, nil, PREC_NONE},
		globals.TOKEN_ELSE:          {nil, nil, PREC_NONE},
		globals.TOKEN_FALSE:         {literal, nil, PREC_NONE},
		globals.TOKEN_FOR:           {nil, nil, PREC_NONE},
		globals.TOKEN_FUN:           {nil, nil, PREC_NONE},
		globals.TOKEN_IF:            {nil, nil, PREC_NONE},
		globals.TOKEN_NIL:           {literal, nil, PREC_NONE},
		globals.TOKEN_OR:            {nil, nil, PREC_NONE},
		globals.TOKEN_PRINT:         {nil, nil, PREC_NONE},
		globals.TOKEN_RETURN:        {nil, nil, PREC_NONE},
		globals.TOKEN_SUPER:         {nil, nil, PREC_NONE},
		globals.TOKEN_THIS:          {nil, nil, PREC_NONE},
		globals.TOKEN_TRUE:          {literal, nil, PREC_NONE},
		globals.TOKEN_VAR:           {nil, nil, PREC_NONE},
		globals.TOKEN_WHILE:         {nil, nil, PREC_NONE},
		globals.TOKEN_ERROR:         {nil, nil, PREC_NONE},
		globals.TOKEN_EOF:           {nil, nil, PREC_NONE},
	}
}
