package src

type Chunk struct {
	Code      []uint8
	Constants ValueArray
	Lines     []int
	Count     int
	Capacity  int
}

func InitChunk(chunk *Chunk) {
	chunk.Capacity = 0
	chunk.Count = 0
	chunk.Lines = nil
	chunk.Code = nil
	InitValueArray(&chunk.Constants)
}
func FreeChunk(chunk *Chunk) {
	FreeArray(chunk.Code, chunk.Capacity)
	FreeArray(chunk.Lines, chunk.Capacity)
	FreeValueArray(&chunk.Constants)
	InitChunk(chunk)
}
func WriteChunk(chunk *Chunk, bytecode uint8, line int) {
	if chunk.Capacity < chunk.Count+1 {
		oldcapacity := chunk.Capacity
		chunk.Capacity = GrowCapacity(oldcapacity)
		chunk.Code = GrowArrayChunks(chunk.Code, oldcapacity, chunk.Capacity)
		chunk.Lines = GrowArrayLines(chunk.Lines, oldcapacity, chunk.Capacity)
	}

	chunk.Code[chunk.Count] = bytecode
	chunk.Lines[chunk.Count] = line
	chunk.Count += 1
}

func AddConstants(chunk *Chunk, val Value) int {
	WriteValueArray(&chunk.Constants, val)
	return chunk.Constants.Count - 1
}
