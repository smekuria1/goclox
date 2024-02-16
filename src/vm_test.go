package src

// import "testing"

// func TestVM_run(t *testing.T) {
// 	type fields struct {
// 		chunk    *Chunk
// 		ip       []uint8
// 		stack    [StackMax]Value
// 		stackTop int32
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   InterpretResult
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			vm := &VM{
// 				chunk:    tt.fields.chunk,
// 				ip:       tt.fields.ip,
// 				stack:    tt.fields.stack,
// 				stackTop: tt.fields.stackTop,
// 			}
// 			if got := vm.run(); got != tt.want {
// 				t.Errorf("VM.run() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
