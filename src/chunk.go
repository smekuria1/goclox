package src

// Chunk represents a chunk of bytecode.
type Chunk struct {
	Code      []uint8    // The bytecode of the chunk.
	Constants ValueArray // The constants of the chunk.
	Lines     []int      // The line numbers of the chunk.
	Count     int        // The number of instructions in the chunk.
	Capacity  int        // The capacity of the chunk.
}

// InitChunk initializes a Chunk.
//
// The function takes a pointer to a Chunk struct as its parameter.
// It sets the Capacity and Count fields of the Chunk to 0.
// It sets the Lines and Code fields of the Chunk to nil.
// It initializes the Constants field of the Chunk using the InitValueArray function.
func InitChunk(chunk *Chunk) {
	chunk.Capacity = 0
	chunk.Count = 0
	chunk.Lines = nil
	chunk.Code = nil
	InitValueArray(&chunk.Constants)
}

// FreeChunk frees the memory allocated for a given Chunk.
//
// It takes a pointer to a Chunk as its parameter.
// The function does not return anything.
func FreeChunk(chunk *Chunk) {
	FreeArray(chunk.Code, chunk.Capacity)
	FreeArray(chunk.Lines, chunk.Capacity)
	FreeValueArray(&chunk.Constants)
	InitChunk(chunk)
}

// WriteChunk writes a bytecode to the given chunk at the specified line.
//
// Parameters:
// - chunk: A pointer to the Chunk struct that represents the chunk.
// - bytecode: The bytecode to be written to the chunk.
// - line: The line number where the bytecode is written.
func WriteChunk(chunk *Chunk, bytecode uint8, line int) {
	if chunk.Capacity < chunk.Count+1 {
		oldcapacity := chunk.Capacity
		chunk.Capacity = GrowCapacity(oldcapacity)
		chunk.Code = GrowArrayChunks(chunk.Code, oldcapacity, chunk.Capacity)
		chunk.Lines = GrowArrayLines(chunk.Lines, oldcapacity, chunk.Capacity)
	}

	chunk.Code[chunk.Count] = bytecode
	chunk.Lines[chunk.Count] = line
	chunk.Count++
}

// AddConstants adds a constant value to the chunk's list of constants.
//
// Parameters:
// - chunk: a pointer to the Chunk struct representing the chunk of bytecode.
// - val: the Value to be added to the constants list.
//
// Returns:
// - int: the index of the added constant in the constants list.
func AddConstants(chunk *Chunk, val Value) int {
	WriteValueArray(&chunk.Constants, val)
	return chunk.Constants.Count - 1
}
