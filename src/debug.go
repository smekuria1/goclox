package src

import (
	"fmt"

	"github.com/smekuria1/goclox/globals"
)

func DisassembleChunk(chunk *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < chunk.Count; {
		offset = DisassembleInstruction(chunk, offset)
	}
}

func DisassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	if offset > 0 && chunk.Lines[offset] == chunk.Lines[offset-1] {
		fmt.Printf(" | ")
	} else {
		fmt.Printf("%4d ", chunk.Lines[offset])
	}

	instruction := chunk.Code[offset]

	switch instruction {
	case uint8(globals.OP_RETURN):
		return simpleInstruction("OP_RETURN", offset)
	case uint8(globals.OP_CONSTANT):
		return constantInstruction("OP_CONSTANT", chunk, offset)
	case uint8(globals.OP_NIL):
		return simpleInstruction("OP_NIL", offset)
	case uint8(globals.OP_TRUE):
		return simpleInstruction("OP_TRUE", offset)
	case uint8(globals.OP_FALSE):
		return simpleInstruction("OP_FALSE", offset)
	case uint8(globals.OP_NEGATE):
		return simpleInstruction("OP_NEGATE", offset)
	case uint8(globals.OP_ADD):
		return simpleInstruction("OP_ADD", offset)
	case uint8(globals.OP_SUBTRACT):
		return simpleInstruction("OP_SUBTRACT", offset)
	case uint8(globals.OP_MULTIPLY):
		return simpleInstruction("OP_MULTIPLY", offset)
	case uint8(globals.OP_DIVIDE):
		return simpleInstruction("OP_DIVIDE", offset)
	case uint8(globals.OP_NOT):
		return simpleInstruction("OP_NOT", offset)
	case uint8(globals.OP_EQUAL):
		return simpleInstruction("OP_EQUAL", offset)
	case uint8(globals.OP_GREATER):
		return simpleInstruction("OP_GREATER", offset)
	case uint8(globals.OP_LESS):
		return simpleInstruction("OP_LESS", offset)
	case uint8(globals.OP_PRINT):
		return simpleInstruction("OP_PRINT", offset)
	case uint8(globals.OP_POP):
		return simpleInstruction("OP_POP", offset)
	case uint8(globals.OP_DEFINE_GLOBAL):
		return constantInstruction("OP_DEFINE_GLOBAL", chunk, offset)
	case uint8(globals.OP_GET_GLOBAL):
		return constantInstruction("OP_GET_GLOBAL", chunk, offset)
	case uint8(globals.OP_SET_GLOBAL):
		return constantInstruction("OP_SET_GLOBAL", chunk, offset)
	default:
		fmt.Println("Unknown opcode ", instruction)
		return offset + 1
	}

}

func simpleInstruction(opcode string, offset int) int {
	fmt.Printf("%s\n", opcode)
	return offset + 1
}

func constantInstruction(opcode string, chunk *Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", opcode, constant)
	PrintValue(chunk.Constants.Values[constant])
	fmt.Printf("'\n")
	return offset + 2
}
