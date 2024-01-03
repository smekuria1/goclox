package src

import (
	"fmt"

	"github.com/smekuria1/goclox/globals"
)

var scanner Scanner

func Compile(source string) {

	scanner.InitScanner(source)
	line := -1

	for {
		token := scanner.ScanToken(&source)
		if token.Line != line {
			fmt.Printf("%4d ", token.Line)
			line = token.Line
		} else {
			fmt.Printf("    | ")
		}

		fmt.Printf("%2d '%.*s'\n", token.TOKENType, token.Length, string(source[token.Start:]))

		if token.TOKENType == globals.TOKEN_EOF {
			break
		}
	}
}
