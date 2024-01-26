package src

import (
	"bytes"

	"github.com/smekuria1/goclox/globals"
)

// Scanner represents a lexical scanner for the source code.
type Scanner struct {
	Start   int     // Start represents the start position of the scanner.
	Current int     // Current represents the current position of the scanner.
	Line    int     // Line represents the current line number.
	Source  *string // Source is a pointer to the source code being scanned.
}

// Token represents a lexical token in the code.
type Token struct {
	TOKENType globals.TokenType // Represents the type of the token.
	Start     int               // Represents the starting position of the token.
	Length    int               // Represents the length of the token.
	Line      int               // Represents the line number where the token is found.
}

// InitScanner initializes the Scanner struct with the given source.
//
// Parameters:
// - source: a string representing the source code to be scanned.
//
// Return type: none.
func (scanner *Scanner) InitScanner(source string) {
	scanner.Start = 0
	scanner.Current = 0
	scanner.Source = &source
	scanner.Line = 1
}

// ScanToken scans the source string and returns a Token.
//
// It takes a pointer to a Scanner and a pointer to a string as parameters.
// The Scanner is used to keep track of the current position in the source string.
// The source string contains the code to be scanned.
//
// It returns a Token based on the current character in the source string.
// The Token represents the type of the current character.
func (scanner *Scanner) ScanToken(source *string) Token {
	scanner.skipWhitespace()
	scanner.Start = scanner.Current
	if scanner.isAtEnd() {
		return makeToken(globals.TokenEOF, scanner)
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
		return makeToken(globals.TokenLeftParen, scanner)
	case ')':
		return makeToken(globals.TokenRightParen, scanner)
	case '{':
		return makeToken(globals.TokenLeftBrace, scanner)
	case '}':
		return makeToken(globals.TokenRightBrace, scanner)
	case ';':
		return makeToken(globals.TokenSEMICOLON, scanner)
	case ',':
		return makeToken(globals.TokenCOMMA, scanner)
	case '.':
		return makeToken(globals.TokenDOT, scanner)
	case '-':
		return makeToken(globals.TokenMINUS, scanner)
	case '+':
		return makeToken(globals.TokenPLUS, scanner)
	case '/':
		return makeToken(globals.TokenSLASH, scanner)
	case '*':
		return makeToken(globals.TokenSTAR, scanner)
	case '!':
		return makeToken(
			func() globals.TokenType {
				if scanner.match('=') {
					return globals.TokenBANG_EQUAL
				}
				return globals.TokenBANG
			}(), scanner)
	case '=':
		return makeToken(
			func() globals.TokenType {
				if scanner.match('=') {
					return globals.TokenEQUAL_EQUAL
				}
				return globals.TokenEQUAL
			}(), scanner)
	case '<':
		return makeToken(
			func() globals.TokenType {
				if scanner.match('=') {
					return globals.TokenLESS_EQUAL
				}
				return globals.TokenLESS
			}(), scanner)
	case '>':
		return makeToken(
			func() globals.TokenType {
				if scanner.match('=') {
					return globals.TokenGREATER_EQUAL
				}
				return globals.TokenGREATER
			}(), scanner)
	case '"':
		return scanner.checkString()
	default:
		// Handle other cases or return an Error token
		return makeErrorToken("Unexpected character.", scanner)
	}

}

// identifier scans and returns an identifier Token.
//
// The identifier function scans the input string and checks if the current character is an alphabetic character or a digit. It continues scanning until it reaches a non-alphanumeric character. It then returns a Token representing the identifier.
//
// It takes no parameters.
// It returns a Token.
func (scanner *Scanner) identifier() Token {
	for scanner.isAlpha(scanner.peek()) || scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}
	return makeToken(scanner.identifierType(), scanner)
}

// identifierType returns the token type of an identifier in the Scanner struct.
//
// No parameters.
// Returns the TokenType of the identifier.
func (scanner *Scanner) identifierType() globals.TokenType {
	source := *scanner.Source
	switch source[scanner.Start] {
	case 'a':
		return scanner.checkKeyword(1, 2, "nd", globals.TokenAND)
	case 'c':
		return scanner.checkKeyword(1, 4, "lass", globals.TokenCLASS)
	case 'e':
		return scanner.checkKeyword(1, 3, "lse", globals.TokenELSE)
	case 'i':
		return scanner.checkKeyword(1, 1, "f", globals.TokenIF)
	case 'n':
		return scanner.checkKeyword(1, 2, "il", globals.TokenNIL)
	case 'o':
		return scanner.checkKeyword(1, 1, "r", globals.TokenOR)
	case 'p':
		return scanner.checkKeyword(1, 4, "rint", globals.TokenPRINT)
	case 'r':
		return scanner.checkKeyword(1, 5, "eturn", globals.TokenRETURN)
	case 's':
		return scanner.checkKeyword(1, 4, "uper", globals.TokenSUPER)
	case 'v':
		return scanner.checkKeyword(1, 2, "ar", globals.TokenVAR)
	case 'w':
		return scanner.checkKeyword(1, 4, "hile", globals.TokenWHILE)
	case 'f':
		firstCheck := scanner.checkKeyword(1, 4, "alse", globals.TokenFALSE)
		if firstCheck == globals.TokenIDENTIFIER {
			firstCheck = scanner.checkKeyword(1, 2, "or", globals.TokenFOR)
		}
		return firstCheck
	default:
		return globals.TokenIDENTIFIER

	}
}

// checkKeyword checks if a given substring matches a keyword and returns the corresponding token type.
//
// Parameters:
//
// - start: the starting index of the substring to check.
//
// - length: the length of the substring to check.
//
// - rest: the keyword to match against the substring.
//
// - tokenType: the token type to return if the substring matches the keyword.
//
// Returns the token type corresponding to the keyword if the substring matches the keyword, otherwise returns globals.TokenIDENTIFIER.
func (scanner *Scanner) checkKeyword(start, length int, rest string, tokenType globals.TokenType) globals.TokenType {
	if scanner.Current-scanner.Start == start+length &&
		bytes.Equal([]byte((*scanner.Source)[scanner.Start+start:scanner.Start+start+length]), []byte(rest)) {
		return tokenType
	}
	return globals.TokenIDENTIFIER
}

// isAlpha checks if the given rune is an alphabetic character or an underscore.
//
// Parameters:
// - c: the rune to be checked.
//
// Returns:
// - bool: true if the rune is an alphabetic character or an underscore, false otherwise.
func (scanner *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isDigit checks if the given character is a digit.
//
// c - the character to be checked.
// Returns true if the character is a digit, false otherwise.
func (scanner *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// number scans and returns a Token representing a number.
//
// It scans digits until a non-digit character is encountered.
// If a dot is encountered, it scans more digits to handle floating-point numbers.
// The function returns a Token of type TokenNUMBER.
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
	return makeToken(globals.TokenNUMBER, scanner)
}

// checkString scans the input and checks if it is a valid string.
//
// It scans the input until it finds a closing double quote (") or reaches the end of the input.
// It increments the line count if a newline character is encountered.
// If the input ends without finding a closing double quote, it returns an Error token.
// Otherwise, it creates a token of type TokenSTRING and returns it.
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
	return makeToken(globals.TokenSTRING, scanner)
}

// advance advances the scanner to the next rune and returns it.
//
// It increments the scanner's current position and returns the rune at that position.
// If the scanner has reached the end, it returns 0 or any appropriate value to indicate the end.
func (scanner *Scanner) advance() rune {
	if !scanner.isAtEnd() {
		scanner.Current++
		ret := *scanner.Source
		return rune(ret[scanner.Current-1])
	}
	return 0 // or any appropriate value to indicate the end
}

// skipWhitespace skips over any whitespace characters in the input string.
// It advances the scanner's position until a non-whitespace character is encountered.
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

// peekNext returns the next rune in the scanner's source code without consuming it.
//
// It checks if the scanner is at the end or if the current position is at the end of the source code,
// and returns 0 or any appropriate value to indicate the end.
// Otherwise, it returns the rune at the next position.
func (scanner *Scanner) peekNext() rune {
	if scanner.isAtEnd() || scanner.Current+1 >= len(*scanner.Source) {
		return 0 // or any appropriate value to indicate the end
	}
	check := *scanner.Source
	return rune(check[scanner.Current+1])
}

// peek returns the next rune in the input without consuming it.
//
// It returns 0 if the scanner is at the end of the input.
// It takes no parameters.
// It returns a rune.
func (scanner *Scanner) peek() rune {
	if scanner.isAtEnd() {
		return 0 // or any appropriate value to indicate the end
	}
	check := *scanner.Source
	return rune(check[scanner.Current])
}

// match checks if the next character in the source matches the expected rune.
//
// It takes the expected rune as a parameter.
// It returns a boolean value indicating whether the match was successful.
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

// isAtEnd checks if the scanner has reached the end of the source string.
//
// It returns a boolean value indicating whether the scanner's current position is greater than or equal to
// the length of the source string.
func (scanner *Scanner) isAtEnd() bool {
	return scanner.Current >= len(*scanner.Source)
}

// makeToken creates a token of the given type using the provided scanner.
//
// Parameters:
// - tokentype: the type of the token to create.
//
// - scanner: a pointer to the Scanner object for tokenization.
//
// Returns:
// - Token: the created token.
func makeToken(tokentype globals.TokenType, scanner *Scanner) Token {
	var token Token
	token.TOKENType = tokentype
	token.Start = scanner.Start
	token.Length = scanner.Current - scanner.Start
	token.Line = scanner.Line
	return token
}

// makeErrorToken creates an Error token with the given message and scanner.
//
// Parameters:
//
// - message: the Error message for the token.
//
// - scanner: the scanner used to create the token.
//
// Return:
// - Token: the Error token.
func makeErrorToken(message string, scanner *Scanner) Token {
	var token Token
	token.TOKENType = globals.TokenERROR
	token.Start = scanner.Start
	token.Length = len(message)
	token.Line = scanner.Line
	return token
}
