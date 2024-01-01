package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/smekuria1/goclox/globals"
	"github.com/smekuria1/goclox/src"
)

var cpuprof = flag.String("cpuprof", "", "write cpu profile to file")
var memprof = flag.String("memprof", "", "write memory profile to `file`")

func main() {
	flag.IntVar(&globals.DEBUG_TRACE_EXECUTION, "debug", 0, "Turn on debug mode")
	flag.Parse()

	if *cpuprof != "" {
		f, err := os.Create(*cpuprof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var chunk src.Chunk
	src.InitVM()
	src.InitChunk(&chunk)
	constant := src.AddConstants(&chunk, 1.2)
	src.WriteChunk(&chunk, uint8(globals.OP_CONSTANT), 122)
	src.WriteChunk(&chunk, uint8(constant), 122)
	constant2 := src.AddConstants(&chunk, 3.4)
	src.WriteChunk(&chunk, uint8(globals.OP_CONSTANT), 124)
	src.WriteChunk(&chunk, uint8(constant2), 122)

	constant3 := src.AddConstants(&chunk, 5.6)
	src.WriteChunk(&chunk, uint8(globals.OP_CONSTANT), 123)
	src.WriteChunk(&chunk, uint8(constant3), 123)
	src.WriteChunk(&chunk, uint8(globals.OP_DIVIDE), 123)
	src.WriteChunk(&chunk, uint8(globals.OP_NEGATE), 123)
	src.WriteChunk(&chunk, uint8(globals.OP_RETURN), 123)

	src.Interpret(&chunk)
	//src.DisassembleChunk(&chunk, "Test chunk")

	if *memprof != "" {
		f, err := os.Create(*memprof)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
	src.FreeChunk(&chunk)
	src.FreeVM()
}
