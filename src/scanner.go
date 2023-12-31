package src

import (
	"bytes"

	"github.com/smekuria1/goclox/globals"
)

type Scanner struct {
	Start   int
	Current int
	Line    int
	Source  *string // I know this is not ideal stuck too close to the book
}

type Token struct {
	TOKENType globals.TokenType
	Start     int
	Length    int
	Line      int
}

func (scanner *Scanner) InitScanner(source string) {
	scanner.Start = 0
	scanner.Current = 0
	scanner.Source = &source
	scanner.Line = 1
}

func (scanner *Scanner) ScanToken(source *string) Token {
	scanner.skipWhitespace()
	scanner.Start = scanner.Current
	if scanner.isAtEnd() {
		return makeToken(globals.TOKEN_EOF, scanner)
	}
	c := scanner.advance()
	if scanner.isAlpha(c) {
		return scanner.identifier()
	}
	if scanner.isDigit(c) {
		return scanner.number()
	}
	switch c {
	case '(':
		return makeToken(globals.TOKEN_LEFT_PAREN, scanner)
	case ')':
		return makeToken(globals.TOKEN_RIGHT_PAREN, scanner)
	case '{':
		return makeToken(globals.TOKEN_LEFT_BRACE, scanner)
	case '}':
		return makeToken(globals.TOKEN_RIGHT_BRACE, scanner)
	case ';':
		return makeToken(globals.TOKEN_SEMICOLON, scanner)
	case ',':
		return makeToken(globals.TOKEN_COMMA, scanner)
	case '.':
		return makeToken(globals.TOKEN_DOT, scanner)
	case '-':
		return makeToken(globals.TOKEN_MINUS, scanner)
	case '+':
		return makeToken(globals.TOKEN_PLUS, scanner)
	case '/':
		return makeToken(globals.TOKEN_SLASH, scanner)
	case '*':
		return makeToken(globals.TOKEN_STAR, scanner)
	case '!':
		return makeToken(
			func() globals.TokenType {
				if scanner.match('=') {
					return globals.TOKEN_BANG_EQUAL
				}
				return globals.TOKEN_BANG
			}(), scanner)
	case '=':
		return makeToken(
			func() globals.TokenType {
				if scanner.match('=') {
					return globals.TOKEN_EQUAL_EQUAL
				}
				return globals.TOKEN_EQUAL
			}(), scanner)
	case '<':
		return makeToken(
			func() globals.TokenType {
				if scanner.match('=') {
					return globals.TOKEN_LESS_EQUAL
				}
				return globals.TOKEN_LESS
			}(), scanner)
	case '>':
		return makeToken(
			func() globals.TokenType {
				if scanner.match('=') {
					return globals.TOKEN_GREATER_EQUAL
				}
				return globals.TOKEN_GREATER
			}(), scanner)
	case '"':
		return scanner.checkString()
	default:
		// Handle other cases or return an error token
		return makeErrorToken("Unexpected character.", scanner)
	}

}
func (scanner *Scanner) identifier() Token {
	for scanner.isAlpha(scanner.peek()) || scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}
	return makeToken(scanner.identifierType(), scanner)
}
func (scanner *Scanner) identifierType() globals.TokenType {
	source := *scanner.Source
	switch source[scanner.Start] {
	case 'a':
		return scanner.checkKeyword(1, 2, "nd", globals.TOKEN_AND)
	case 'c':
		return scanner.checkKeyword(1, 4, "lass", globals.TOKEN_CLASS)
	case 'e':
		return scanner.checkKeyword(1, 3, "lse", globals.TOKEN_ELSE)
	case 'i':
		return scanner.checkKeyword(1, 1, "f", globals.TOKEN_IF)
	case 'n':
		return scanner.checkKeyword(1, 2, "il", globals.TOKEN_NIL)
	case 'o':
		return scanner.checkKeyword(1, 1, "r", globals.TOKEN_OR)
	case 'p':
		return scanner.checkKeyword(1, 4, "rint", globals.TOKEN_PRINT)
	case 'r':
		return scanner.checkKeyword(1, 5, "eturn", globals.TOKEN_RETURN)
	case 's':
		return scanner.checkKeyword(1, 4, "uper", globals.TOKEN_SUPER)
	case 'v':
		return scanner.checkKeyword(1, 2, "ar", globals.TOKEN_VAR)
	case 'w':
		return scanner.checkKeyword(1, 4, "hile", globals.TOKEN_WHILE)
	default:
		return globals.TOKEN_IDENTIFIER

	}
}

func (scanner *Scanner) checkKeyword(start, length int, rest string, tokenType globals.TokenType) globals.TokenType {
	if scanner.Current-scanner.Start == start+length &&
		bytes.Equal([]byte((*scanner.Source)[scanner.Start+start:scanner.Start+start+length]), []byte(rest)) {
		return tokenType
	}
	return globals.TOKEN_IDENTIFIER
}
func (scanner *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (scanner *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (scanner *Scanner) number() Token {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && scanner.isDigit(scanner.peekNext()) {
		scanner.advance()
	}
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}
	return makeToken(globals.TOKEN_NUMBER, scanner)
}

func (scanner *Scanner) checkString() Token {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.Line++
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		return makeErrorToken("Unterminated String.", scanner)
	}
	scanner.advance()
	return makeToken(globals.TOKEN_STRING, scanner)
}

func (scanner *Scanner) advance() rune {
	if !scanner.isAtEnd() {
		scanner.Current++
		ret := *scanner.Source
		return rune(ret[scanner.Current-1])
	}
	return 0 // or any appropriate value to indicate the end
}

func (scanner *Scanner) skipWhitespace() {
	for {
		c := scanner.peek()
		switch c {
		case ' ', '\r', '\t':
			scanner.advance()
		case '/':
			if scanner.peekNext() == '/' {
				// A comment goes until the end of the line.
				for scanner.peek() != '\n' && !scanner.isAtEnd() {
					scanner.advance()
				}
			} else {
				return
			}

		case '\n':
			scanner.Line++
			scanner.advance()
		default:
			return
		}
	}
}

func (scanner *Scanner) peekNext() rune {
	if scanner.isAtEnd() || scanner.Current+1 >= len(*scanner.Source) {
		return 0 // or any appropriate value to indicate the end
	}
	check := *scanner.Source
	return rune(check[scanner.Current+1])
}

func (scanner *Scanner) peek() rune {
	if scanner.isAtEnd() {
		return 0 // or any appropriate value to indicate the end
	}
	check := *scanner.Source
	return rune(check[scanner.Current])
}

func (scanner *Scanner) match(expected rune) bool {
	if scanner.isAtEnd() {
		return false
	}
	check := *scanner.Source
	if rune(check[scanner.Current]) != expected {
		return false
	}

	scanner.Current++
	return true
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.Current >= len(*scanner.Source)
}

func makeToken(tokentype globals.TokenType, scanner *Scanner) Token {
	var token Token
	token.TOKENType = tokentype
	token.Start = scanner.Start
	token.Length = scanner.Current - scanner.Start
	token.Line = scanner.Line
	return token
}

func makeErrorToken(message string, scanner *Scanner) Token {
	var token Token
	token.TOKENType = globals.TOKEN_ERROR
	token.Start = scanner.Start
	token.Length = len(message)
	token.Line = scanner.Line
	return token
}
