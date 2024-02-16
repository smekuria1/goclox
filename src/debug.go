package src

import (
	"fmt"

	"github.com/smekuria1/goclox/globals"
)

// DisassembleChunk prints the disassembled instructions of a given chunk.
//
// The function takes a pointer to a Chunk struct and a name string as parameters.
// It does not return any value.
func DisassembleChunk(chunk *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < chunk.Count; {
		offset = DisassembleInstruction(chunk, offset)
	}
}

// DisassembleInstruction disassembles an instruction in the given chunk at the specified offset.
//
// Parameters:
// - chunk: A pointer to the Chunk struct representing the chunk of code.
// - offset: An integer representing the offset of the instruction in the chunk.
//
// Return:
// - An integer representing the new offset after processing the instruction.
func DisassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	if offset > 0 && chunk.Lines[offset] == chunk.Lines[offset-1] {
		fmt.Printf(" | ")
	} else {
		fmt.Printf("%4d ", chunk.Lines[offset])
	}

	instruction := chunk.Code[offset]

	switch instruction {
	case uint8(globals.OpReturn):
		return simpleInstruction("OpReturn", offset)
	case uint8(globals.OpConstant):
		return constantInstruction("OpConstant", chunk, offset)
	case uint8(globals.OpNil):
		return simpleInstruction("OpNil", offset)
	case uint8(globals.OpTrue):
		return simpleInstruction("OpTrue", offset)
	case uint8(globals.OpFalse):
		return simpleInstruction("OpFalse", offset)
	case uint8(globals.OpNegate):
		return simpleInstruction("OpNegate", offset)
	case uint8(globals.OpAdd):
		return simpleInstruction("OpAdd", offset)
	case uint8(globals.OpSubtract):
		return simpleInstruction("OpSubtract", offset)
	case uint8(globals.OpMultiply):
		return simpleInstruction("OpMultiply", offset)
	case uint8(globals.OpDivide):
		return simpleInstruction("OpDivide", offset)
	case uint8(globals.OpNot):
		return simpleInstruction("OpNot", offset)
	case uint8(globals.OpEqual):
		return simpleInstruction("OpEqual", offset)
	case uint8(globals.OpGreater):
		return simpleInstruction("OpGreater", offset)
	case uint8(globals.OpLess):
		return simpleInstruction("OpLess", offset)
	case uint8(globals.OpPrint):
		return simpleInstruction("OpPrint", offset)
	case uint8(globals.OpPop):
		return simpleInstruction("OpPop", offset)
	case uint8(globals.OpDefineGlobal):
		return constantInstruction("OpDefineGlobal", chunk, offset)
	case uint8(globals.OpGetGlobal):
		return constantInstruction("OpGetGlobal", chunk, offset)
	case uint8(globals.OpSetGlobal):
		return constantInstruction("OpSetGlobal", chunk, offset)
	case uint8(globals.OpGetLocal):
		return byteInstruction("OpGetLocal", chunk, offset)
	case uint8(globals.OpSetLocal):
		return byteInstruction("OpSetLocal", chunk, offset)
	case uint8(globals.OpJump):
		return jumpInstruction("OpJump", 1, chunk, offset)
	case uint8(globals.OpJumpFalse):
		return jumpInstruction("OpJumpElse", 1, chunk, offset)
	case uint8(globals.OpLoop):
		return jumpInstruction("OpLoop", -1, chunk, offset)
	case uint8(globals.OpCall):
		return byteInstruction("OpCall", chunk, offset)
	default:
		fmt.Println("Unknown opcode ", instruction)
		return offset + 1
	}

}

// simpleInstruction prints the given opcode and returns the offset incremented by 1.
//
// Parameters:
// - opcode: a string representing the opcode to be printed.
// - offset: an integer representing the current offset.
//
// Returns:
// - an integer representing the new offset after incrementing it by 1.
func simpleInstruction(opcode string, offset int) int {
	fmt.Printf("%s\n", opcode)
	return offset + 1
}

// constantInstruction prints an opcode and its corresponding constant value.
//
// It takes in the opcode string, the chunk pointer, and the offset integer as parameters.
// It returns an integer representing the updated offset.
func constantInstruction(opcode string, chunk *Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", opcode, constant)
	PrintValue(chunk.Constants.Values[constant])
	fmt.Printf("'\n")
	return offset + 2
}

// byteInstruction prints the opcode and slot of a byte instruction.
//
// It takes the following parameter(s):
// - opcode: a string representing the opcode of the byte instruction.
// - chunk: a pointer to the Chunk struct representing the chunk of code being executed.
// - offset: an integer representing the offset of the current byte instruction in the chunk.
//
// It returns an integer representing the updated offset after processing the byte instruction.
func byteInstruction(opcode string, chunk *Chunk, offset int) int {
	slot := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d\n", opcode, slot)
	return offset + 2
}

// jumpInstruction
func jumpInstruction(name string, sign int, chunk *Chunk, offset int) int {
	jump := uint16(chunk.Code[offset+1])<<8 | uint16(chunk.Code[offset+2])
	fmt.Printf("%-16s %4d -> %d\n", name, offset, offset+3+sign*int(jump))
	return offset + 3
}
