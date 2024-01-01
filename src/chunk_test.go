package src

import "testing"

func TestWriteChunk(t *testing.T) {
	type args struct {
		chunk    *Chunk
		bytecode uint8
		line     int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteChunk(tt.args.chunk, tt.args.bytecode, tt.args.line)
		})
	}
}
