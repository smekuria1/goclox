package src

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/smekuria1/goclox/globals"
)

// StackMax represents the maximum size of the stack.
const StackMax = 256

// FrameMax represents the maximum number of call frames.
const FrameMax = 64

// VM represents a virtual machine.
type VM struct {
	chunk *Chunk // Stores the bytecode of the program being executed.
	//ip             []uint8 // Keeps track of the current instruction pointer.
	//instructionPtr int
	frame      [FrameMax]CallFrame // Stores the call frames of the virtual machine.
	frameCount int                 // Keeps track of the number of call frames.
	stack      []Value             // Stores the values of the virtual machine's stack.
	stackTop   int                 // Keeps track of the top of the stack.
	objects    *Obj                // Stores a linked list of all dynamically allocated objects.
	strings    *Table              // Stores a table of string objects.
	globals    *Table              // Stores a table of global variables.

}

// InterpretResult represents the result of an interpretation.
type InterpretResult int

const (
	// InterpretOk indicates successful interpretation
	InterpretOk InterpretResult = iota

	// InterpretCompileError indicates a compilation Error during interpretation
	InterpretCompileError

	// InterpretRuntimeError indicates a runtime Error during interpretation
	InterpretRuntimeError
)

type CallFrame struct {
	function *ObjFunction // Stores the function object of the function being called.
	slots    []Value      // Stores the slots of the call frame.
	// slotTop  int          // Keeps track of the top of the slots.
	fp    []uint8 // Stores the frame pointers of the call frame.
	fpPtr int     //  tracks the current frame pointer
}

var vm VM

// InitVM initializes the virtual machine.
//
// It resets the stack, clears the objects, and initializes the strings and globals tables.
func InitVM() {
	vm.ResetStack()
	// vm.instructionPtr = 0
	vm.objects = nil
	vm.strings = &Table{}
	vm.globals = &Table{}
	vm.stack = make([]Value, StackMax)
	vm.globals.InitTable()
	vm.strings.InitTable()
}

// ResetStack resets the stack of the VM.
//
// No parameters.
// No return type.
func (vm *VM) ResetStack() {
	vm.stackTop = 0
	vm.frameCount = 0
}

// FreeVM frees the virtual machine by calling the Freetable method on the vm.strings and vm.globals variables,
// and the FreeObjects function on the vm.objects variable.
//
// No parameters.
// No return value.
func FreeVM() {

	vm.strings.Freetable()
	vm.globals.Freetable()
	FreeObjects(vm.objects)
}

// Push pushes a value onto the stack.
//
// value: the value to be pushed onto the stack.
func (vm *VM) Push(value Value) {
	// //DEBATING HOW TO HANDLE OVERFLOW
	// if vm.stackTop+1 >= StackMax {

	// }
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

// Pop removes and returns the top element from the stack.
//
// No parameters.
// Returns a Value.
func (vm *VM) Pop() Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

// Peek returns the value at the top of the stack without removing it.
//
// It does not take any parameters.
// It returns a Value.
func (vm *VM) Peek(index ...int) Value {
	if vm.stackTop != 0 {
		return vm.stack[vm.stackTop-1]
	}
	return vm.stack[vm.stackTop]
}

// Interpret interprets the given source code and returns the interpretation result.
//
// Parameters:
// - source: The source code to be interpreted as a string.
//
// Return type:
// - InterpretResult: The result of the interpretation.
func Interpret(source string) InterpretResult {
	var chunk Chunk
	InitChunk(&chunk)
	function := Compile(source, &chunk)
	if function == nil {
		FreeChunk(&chunk)
		return InterpretCompileError
	}
	vm.Push(ObjVal(function))
	callValue(ObjVal(function), 0)
	result := vm.run()
	FreeChunk(&chunk)
	return result

}

// BinaryOp performs a binary operation on the top two values in the VM's stack
// using the provided operation function.
//
// Parameters:
// - op: The operation function that takes two values and returns a value.
//
// Return type: error
func (vm *VM) BinaryOp(op func(Value, Value) Value, offset ...int) error {
	b := vm.Pop()
	a := vm.Pop()
	if a.Type != ValNumber || b.Type != ValNumber {
		return errors.New("cannot compare non numeric")
	}
	vm.Push(op(a, b))
	return nil
}

// ReadByteVM reads a single byte from the VM's instruction pointer.
//
// No parameters.
// Returns a uint8 value.
func (frame *CallFrame) ReadByteVM() uint8 {
	// Dereference the slice pointer and take the address of the first element.
	//result := (*uint8)(unsafe.Pointer(&(vm.ip)[0]))
	result := frame.fp[frame.fpPtr]
	// Increment the slice pointer to point to the next element.
	//vm.ip = (vm.ip)[1:]

	frame.fpPtr++

	return result
}

// ReadConstant retrieves a constant value from the chunk's constant pool.
//
// It reads a byte from the VM's instruction stream and uses it as an index
// into the constant pool. The constant value at that index is then returned.
// The constant pool is stored in the `Values` field of the `Constants`
// field of the `chunk` field of the `VM` struct.
//
// Returns the constant value retrieved from the constant pool.
func (frame *CallFrame) ReadConstant() Value {
	result := frame.function.chunk.Constants.Values[frame.ReadByteVM()]
	return result
}

// runtimeError handles runtime errors in the VM.
//
// It takes the offset and runoffset integers as parameters.
// It does not return anything.
func (vm *VM) runtimeError(offset int, runoffset int, message ...string) {

	// frame := &vm.frame[vm.frameCount-1]
	// instruction := frame.fp[frame.fpPtr]
	// line := frame.function.chunk.Lines[int(instruction)]
	// fmt.Printf("%s line[%d]\n", message[0], line)
	// if len(message) > 1 {
	// 	fmt.Printf("%s", message[1])
	// }
	for i := vm.frameCount - 1; i >= 0; i-- {
		frame := &vm.frame[i]
		function := frame.function
		instruction := frame.fp[frame.fpPtr]

		fmt.Printf("[line %d] in %s\n", function.chunk.Lines[int(instruction)], string(function.name.Chars))
		if function.name == nil {
			fmt.Printf("script\n")
		} else {
			fmt.Printf("function %s\n", string(function.name.Chars))
		}
	}

	vm.ResetStack()
}

// ReadShort reads a 16-bit value from the VM's instruction stream.
func (frame *CallFrame) ReadShort() uint16 {
	value := uint16(frame.fp[frame.fpPtr])<<8 | uint16(frame.fp[frame.fpPtr+1])

	frame.fpPtr += 2

	return value
}

func callValue(calle Value, argcount int) bool {
	if IsValObj(calle) {
		switch AsFunction(calle).obj.Type {
		case ObjFunctionType:
			return fcall(AsFunction(calle), argcount)
		default:
			break
		}
	}
	vm.runtimeError(0, 0, "Can only call functions and classes.")
	return false
}

func fcall(function *ObjFunction, argcount int) bool {
	if argcount != function.arity {
		vm.runtimeError(0, 0, "Expected", strconv.Itoa(function.arity), "arguments but got", strconv.Itoa(argcount))
		return false
	}
	if vm.frameCount == FrameMax {
		vm.runtimeError(0, 0, "Stack overflow.")
		return false
	}

	frame := &vm.frame[vm.frameCount]
	frame.function = function
	frame.fp = function.chunk.Code

	frame.slots = vm.stack[vm.stackTop-argcount-1:]
	vm.frameCount++
	return true
}

/*
run executes the bytecode in the VM's chunk until an Error occurs or the program completes.
During execution, the function interprets each bytecode instruction, performing the
corresponding operations such as pushing constants onto the stack, performing binary
operations, and handling control flow instructions. If debugging is enabled, it prints
the stack and disassembled instructions at each step.

Parameters:
- vm: A pointer to the Virtual Machine executing the bytecode.

Returns:
- InterpretResult: Indicates the result of the interpretation, such as success, Error, or runtime Error.
*/
func (vm *VM) run() InterpretResult {

	frame := &vm.frame[vm.frameCount-1]

	offset := 0
	runoffset := 0
	for {
		if globals.DEBUG_TRACE_EXECUTION {
			fmt.Printf("     ")
			for slot := 0; slot < int(vm.stackTop); slot++ {
				fmt.Print("[")
				PrintValue(vm.stack[slot])
				fmt.Print("]")

			}
			fmt.Print("\n")
			offset = DisassembleInstruction(&frame.function.chunk, frame.fpPtr)
		}

		instruction := frame.ReadByteVM()
		//fmt.Printf("instruction: %v\n", instruction)
		switch instruction {
		case uint8(globals.OpConstant):
			constant := frame.ReadConstant()
			vm.Push(constant)
			runoffset += 2
			//break
		case uint8(globals.OpNil):
			vm.Push(NilValue())
			runoffset++
		case uint8(globals.OpTrue):
			vm.Push(BoolValue(true))
			runoffset++
		case uint8(globals.OpFalse):
			vm.Push(BoolValue(false))
			runoffset++
		case uint8(globals.OpEqual):
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(BoolValue(valuesEqual(a, b)))
			runoffset++
		case uint8(globals.OpPrint):
			PrintValue(vm.Pop())
			fmt.Printf("\n")
			runoffset++
		case uint8(globals.OpPop):
			vm.Pop()
			runoffset++
		case uint8(globals.OpCall):
			argcount := frame.ReadByteVM()
			if !callValue(vm.Peek(int(argcount)), int(argcount)) {
				return InterpretRuntimeError
			}
			runoffset += 1
		case uint8(globals.OpGetGlobal):
			name := frame.readString()
			var value Value
			runoffset += 2
			if !vm.globals.TableGet(name, &value) {
				vm.runtimeError(offset, runoffset, "Undefined Variable")
				return InterpretRuntimeError
			}
			vm.Push(value)
		case uint8(globals.OpSetGlobal):
			runoffset += 2
			name := frame.readString()
			if vm.globals.TableSet(name, vm.Peek()) {
				vm.globals.TableDelete(name)
				vm.runtimeError(offset, runoffset, "Undefined variable", string(name.Chars))
				return InterpretRuntimeError
			}
		case uint8(globals.OpDefineGlobal):
			runoffset += 2
			name := frame.readString()
			peeked := vm.Peek()
			vm.globals.TableSet(name, peeked)
			vm.Pop()

		case uint8(globals.OpGetLocal):
			runoffset += 2
			slot := frame.ReadByteVM()
			vm.Push(frame.slots[slot])
		case uint8(globals.OpSetLocal):
			runoffset += 2
			slot := frame.ReadByteVM()
			frame.slots[slot] = vm.Peek()
		case uint8(globals.OpReturn):
			result := vm.Pop()
			vm.frameCount--
			if vm.frameCount == 0 {
				return InterpretOk
			}
			vm.stack = frame.slots
			vm.Push(result)
			frame = &vm.frame[vm.frameCount-1]
			runoffset = 0
		case uint8(globals.OpJumpFalse):
			runoffset += 3
			offsetJumpFalse := frame.ReadShort()
			if isFalsey(vm.Peek()) {
				frame.fpPtr += int(offsetJumpFalse)
			}
		case uint8(globals.OpJump):
			runoffset += 3
			offsetJump := frame.ReadShort()
			frame.fpPtr += int(offsetJump)
		case uint8(globals.OpLoop):
			offsetLoop := int(frame.ReadShort())
			frame.fpPtr -= int(offsetLoop)
		case uint8(globals.OpGreater):
			runoffset++
			err := vm.BinaryOp(func(v1, v2 Value) Value { return BoolValue(v1.As.(float64) > v2.As.(float64)) })
			if err != nil {
				vm.runtimeError(offset, runoffset, err.Error())
				return InterpretRuntimeError
			}
		case uint8(globals.OpLess):
			runoffset++
			vm.BinaryOp(func(v1, v2 Value) Value { return BoolValue(v1.As.(float64) < v2.As.(float64)) })
		case uint8(globals.OpNegate):
			runoffset++
			vm.Push(Value{Type: ValNumber, As: -vm.Pop().As.(float64)})
		case uint8(globals.OpAdd):
			runoffset++
			b := vm.Pop()
			a := vm.Pop()
			if IsObjType(b, ObjStringType) && IsObjType(a, ObjStringType) {
				aString := AsObjString(a)
				bString := AsObjString(b)
				resultString := append(aString.Chars, bString.Chars...)
				hash := hashString(resultString, len(resultString))
				resultObj := allocateString(resultString, len(resultString), ObjStringType, hash)
				vm.Push(ObjStrValue(resultObj))
			} else if IsNumber(a) && IsNumber(b) {
				b := AsNumber(b)
				a := AsNumber(a)
				vm.Push(NumberValue(a + b))
			} else {
				vm.runtimeError(offset, runoffset, "Operands must be two numbers or two strings.")
				return InterpretRuntimeError
			}
		case uint8(globals.OpSubtract):
			runoffset++
			vm.BinaryOp(func(v1, v2 Value) Value { return Value{Type: ValNumber, As: v1.As.(float64) - v2.As.(float64)} })
		case uint8(globals.OpMultiply):
			runoffset++
			vm.BinaryOp(func(v1, v2 Value) Value { return Value{Type: ValNumber, As: v1.As.(float64) * v2.As.(float64)} })
		case uint8(globals.OpDivide):
			runoffset++
			vm.BinaryOp(func(v1, v2 Value) Value { return Value{Type: ValNumber, As: v1.As.(float64) / v2.As.(float64)} })
		case uint8(globals.OpNot):
			runoffset++
			vm.Push(BoolValue(isFalsey(vm.Pop())))
		default:
			fmt.Println("Runtime Error at", vm.chunk.Lines[offset])
			return InterpretRuntimeError
		}

	}

}

// isFalsey checks if a value is falsey.
//
// It takes a parameter `val` of type `Value`.
// It returns a boolean value indicating whether the value is falsey.
func isFalsey(val Value) bool {
	return IsNil(val) || (IsBool(val) && !AsBool(val))
}

// readString returns an ObjectString.
//
// No parameters.
// Returns *ObjectString.
func (frame *CallFrame) readString() *ObjectString {
	return AsObjString(frame.ReadConstant())
}
