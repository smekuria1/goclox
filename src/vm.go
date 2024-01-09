package src

import (
	"fmt"
	"unsafe"

	"github.com/smekuria1/goclox/globals"
)

const STACK_MAX = 30000

type VM struct {
	chunk    *Chunk
	ip       []uint8
	stack    [STACK_MAX]Value
	stackTop int32
	objects  *Obj
}
type InterpretResult int

const (
	// INTERPRET_OK indicates successful interpretation
	INTERPRET_OK InterpretResult = iota

	// INTERPRET_COMPILE_ERROR indicates a compilation error during interpretation
	INTERPRET_COMPILE_ERROR

	// INTERPRET_RUNTIME_ERROR indicates a runtime error during interpretation
	INTERPRET_RUNTIME_ERROR
)

var vm VM

func InitVM() {
	vm.ResetStack()
	vm.objects = nil
}

func (vm *VM) ResetStack() {
	vm.stackTop = 0
}
func FreeVM() {
	FreeObjects(vm.objects)
}

func (vm *VM) Push(value Value) {
	// //DEBATING HOW TO HANDLE OVERFLOW
	// if vm.stackTop+1 >= STACK_MAX {

	// }
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) Pop() Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) Peek() Value {
	return vm.stack[vm.stackTop]
}

func Interpret(source string) InterpretResult {
	var chunk Chunk
	InitChunk(&chunk)

	if !Compile(source, &chunk) {
		FreeChunk(&chunk)
		return INTERPRET_COMPILE_ERROR
	}

	vm.chunk = &chunk
	vm.ip = vm.chunk.Code

	result := vm.run()
	FreeChunk(&chunk)
	return result

}
func (vm *VM) BinaryOp(op func(Value, Value) Value) {
	b := vm.Pop()
	a := vm.Pop()
	vm.Push(op(a, b))
}

func (vm *VM) READ_BYTE() uint8 {
	// Dereference the slice pointer and take the address of the first element.
	result := (*uint8)(unsafe.Pointer(&(vm.ip)[0]))

	// Increment the slice pointer to point to the next element.
	vm.ip = (vm.ip)[1:]

	return *result
}

func (vm *VM) READ_CONSTANT() Value {
	result := vm.chunk.Constants.Values[vm.READ_BYTE()]
	return result
}

func (vm *VM) runtimeError(offset int, message ...string) {
	line := vm.chunk.Lines[offset]
	fmt.Printf("%s , [line %d] in script\n", message[0], line)
	vm.ResetStack()
}

/*
run executes the bytecode in the VM's chunk until an error occurs or the program completes.
During execution, the function interprets each bytecode instruction, performing the
corresponding operations such as pushing constants onto the stack, performing binary
operations, and handling control flow instructions. If debugging is enabled, it prints
the stack and disassembled instructions at each step.

Parameters:
- vm: A pointer to the Virtual Machine executing the bytecode.

Returns:
- InterpretResult: Indicates the result of the interpretation, such as success, error, or runtime error.
*/
func (vm *VM) run() InterpretResult {
	offset := 0
	for {
		if globals.DEBUG_TRACE_EXECUTION {
			fmt.Printf("     ")
			for slot := 0; slot < int(vm.stackTop); slot++ {
				fmt.Print("[")
				PrintValue(vm.stack[slot])
				fmt.Print("]")

			}
			fmt.Print("\n")
			offset = DisassembleInstruction(vm.chunk, offset)
		}

		instruction := vm.READ_BYTE()
		//fmt.Printf("instruction: %v\n", instruction)
		switch instruction {
		case uint8(globals.OP_CONSTANT):
			constant := vm.READ_CONSTANT()
			vm.Push(constant)
			//break
		case uint8(globals.OP_NIL):
			vm.Push(NilValue())
		case uint8(globals.OP_TRUE):
			vm.Push(BoolValue(true))
		case uint8(globals.OP_FALSE):
			vm.Push(BoolValue(false))
		case uint8(globals.OP_EQUAL):
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(BoolValue(valuesEqual(a, b)))
		case uint8(globals.OP_RETURN):
			//TODO: Just for debugging remove when adding actual print functionality
			for i := 0; i < scanner.Line; i++ {
				PrintValue(vm.Pop())
				fmt.Print("\n")
			}
			return INTERPRET_OK
		case uint8(globals.OP_GREATER):
			vm.BinaryOp(func(v1, v2 Value) Value { return BoolValue(v1.As.(float64) > v2.As.(float64)) })
		case uint8(globals.OP_LESS):
			vm.BinaryOp(func(v1, v2 Value) Value { return BoolValue(v1.As.(float64) < v2.As.(float64)) })
		case uint8(globals.OP_NEGATE):
			vm.Push(Value{Type: ValNumber, As: -vm.Pop().As.(float64)})
		case uint8(globals.OP_ADD):
			b := vm.Pop()
			a := vm.Pop()
			if IsObjType(b, ObjStringType) && IsObjType(a, ObjStringType) {
				aString := AsObjString(a)
				bString := AsObjString(b)
				resultString := append(aString.Chars, bString.Chars...)
				resultObj := allocateString(resultString, len(resultString), ObjStringType)
				vm.Push(ObjStrValue(resultObj))
			} else if IsNumber(a) && IsNumber(b) {
				b := AsNumber(b)
				a := AsNumber(a)
				vm.Push(NumberValue(a + b))
			} else {
				vm.runtimeError(offset, "Operands must be two numbers or two strings.")
				return INTERPRET_RUNTIME_ERROR
			}
		case uint8(globals.OP_SUBTRACT):
			vm.BinaryOp(func(v1, v2 Value) Value { return Value{Type: ValNumber, As: v1.As.(float64) - v2.As.(float64)} })
		case uint8(globals.OP_MULTIPLY):
			vm.BinaryOp(func(v1, v2 Value) Value { return Value{Type: ValNumber, As: v1.As.(float64) * v2.As.(float64)} })
		case uint8(globals.OP_DIVIDE):
			vm.BinaryOp(func(v1, v2 Value) Value { return Value{Type: ValNumber, As: v1.As.(float64) / v2.As.(float64)} })
		case uint8(globals.OP_NOT):
			vm.Push(BoolValue(isFalsey(vm.Pop())))
		default:
			fmt.Println("Runtime Error at", vm.chunk.Lines[offset])
			return INTERPRET_RUNTIME_ERROR
		}

	}

}

func isFalsey(val Value) bool {
	return IsNil(val) || (IsBool(val) && !AsBool(val))
}
