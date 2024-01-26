package src

import (
	"reflect"
	"testing"

	"github.com/smekuria1/goclox/globals"
)

func TestScanner_InitScanner(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		scanner *Scanner
		args    args
	}{
		{
			"TestScanner_InitScanner", &Scanner{}, args{"var x=smekuria1;"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.scanner.InitScanner(tt.args.source)
		})
	}
}

func TestScanner_ScanToken(t *testing.T) {
	type args struct {
		source *string
	}
	tests := []struct {
		name    string
		scanner *Scanner
		args    args
		want    Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.ScanToken(tt.args.source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.ScanToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_identifier(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.identifier(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.identifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_identifierType(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    globals.TokenType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.identifierType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.identifierType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_checkKeyword(t *testing.T) {
	type args struct {
		start     int
		length    int
		rest      string
		tokenType globals.TokenType
	}
	tests := []struct {
		name    string
		scanner *Scanner
		args    args
		want    globals.TokenType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.checkKeyword(tt.args.start, tt.args.length, tt.args.rest, tt.args.tokenType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.checkKeyword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_isAlpha(t *testing.T) {
	type args struct {
		c rune
	}
	tests := []struct {
		name    string
		scanner *Scanner
		args    args
		want    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.isAlpha(tt.args.c); got != tt.want {
				t.Errorf("Scanner.isAlpha() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_isDigit(t *testing.T) {
	type args struct {
		c rune
	}
	tests := []struct {
		name    string
		scanner *Scanner
		args    args
		want    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.isDigit(tt.args.c); got != tt.want {
				t.Errorf("Scanner.isDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_number(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.number(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.number() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_checkString(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.checkString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.checkString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_advance(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    rune
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.advance(); got != tt.want {
				t.Errorf("Scanner.advance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_skipWhitespace(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.scanner.skipWhitespace()
		})
	}
}

func TestScanner_peekNext(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    rune
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.peekNext(); got != tt.want {
				t.Errorf("Scanner.peekNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_peek(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    rune
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.peek(); got != tt.want {
				t.Errorf("Scanner.peek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_match(t *testing.T) {
	type args struct {
		expected rune
	}
	tests := []struct {
		name    string
		scanner *Scanner
		args    args
		want    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.match(tt.args.expected); got != tt.want {
				t.Errorf("Scanner.match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_isAtEnd(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scanner.isAtEnd(); got != tt.want {
				t.Errorf("Scanner.isAtEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeToken(t *testing.T) {
	type args struct {
		tokentype globals.TokenType
		scanner   *Scanner
	}
	tests := []struct {
		name string
		args args
		want Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeToken(tt.args.tokentype, tt.args.scanner); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeErrorToken(t *testing.T) {
	type args struct {
		message string
		scanner *Scanner
	}
	tests := []struct {
		name string
		args args
		want Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeErrorToken(tt.args.message, tt.args.scanner); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeErrorToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
