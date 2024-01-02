package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/smekuria1/goclox/globals"
	"github.com/smekuria1/goclox/src"
)

var repl = flag.Bool("repl", false, "Start with REPL mode")
var filename = flag.String("file", "", "Path to gocloxfile")
var help = flag.Bool("help", false, "Display help")

func main() {
	flag.IntVar(&globals.DEBUG_TRACE_EXECUTION, "debug", 0, "Turn on debug mode")
	flag.Parse()

	if *help {
		fmt.Println("-debug bool")
		fmt.Println("    Turn on debug mode")
		fmt.Println("-file string")
		fmt.Println("    Path to gocloxfile")
		fmt.Println("-repl bool")
		fmt.Println("    Start with REPL mode (default \"0\")")
		fmt.Println("-help")
		fmt.Println("    Print this help message")
		return
	}
	src.InitVM()
	if *repl {
		fmt.Println("Running in REPL mode")
		replFunc()
	} else {
		if *filename == "" {
			fmt.Print("Usage: goclox -file path\n")
			return
		}

		fmt.Println("Running file", *filename)
		source := readFile(*filename)
		fmt.Printf("source: %v\n", source)

	}

	src.FreeVM()

}

func readFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		panic("File not found")
	}

	return string(file)

}

func replFunc() {
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		scanner.Scan()
		line = scanner.Text()
		if len(line) != 0 {
			fmt.Println(line)
		} else {
			break
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error : ", scanner.Err())
	}
}
