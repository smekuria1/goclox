package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/profile"
	"github.com/smekuria1/goclox/globals"
	"github.com/smekuria1/goclox/src"
)

var repl = flag.Bool("repl", false, "Start with REPL mode")
var filename = flag.String("file", "", "Path to gocloxfile")
var help = flag.Bool("help", false, "Display help")
var cpuprof = flag.Bool("cpuprof", false, "write cpu profile to file")

//var memprof = flag.Bool("memprof", false, "write memory profile to `file`")

func main() {

	flag.BoolVar(&globals.DEBUG_TRACE_EXECUTION, "debugT", false, "Turn on debug trace execution mode")
	flag.BoolVar(&globals.DEBUG_PRINT_CODE, "debugC", false, "Turn on debug print code mode")
	flag.Parse()
	if *cpuprof {
		defer profile.Start(profile.ProfilePath(".")).Stop()
	}
	if *help {
		fmt.Println("-debugT bool")
		fmt.Println("    Turn on debug trace execution mode")
		fmt.Println("-debugC bool")
		fmt.Println("    Turn on debug print code mode")
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
		src.Interpret(source)

	}

	src.FreeVM()

}

func readFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic("File not found")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		// if scanner.Text() == "\n" {
		// 	continue
		// }
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")

}

func replFunc() {
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Goclox v0.0.1, type q or press Enter to quit REPL")
	for {
		fmt.Printf("> ")
		scanner.Scan()
		line = scanner.Text()
		if line == "q" {
			break
		}
		if len(line) != 0 {
			src.Interpret(line)
		} else {
			break
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error : ", scanner.Err())
	}
}
