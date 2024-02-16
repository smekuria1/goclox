package src

import (
	"bytes"
)

type ObjType int

const (
	ObjStringType ObjType = iota // The type of the string object.
	ObjFunctionType
)

// Obj represents an object in the code.
type Obj struct {
	Type ObjType // The type of the object.
	Next *Obj    // The next object in the list.
}

// ObjFunction represents a function object in the code.
type ObjFunction struct {
	obj   Obj
	arity int
	chunk Chunk
	name  *ObjectString
}

// ObjectString represents a string object in the code.
type ObjectString struct {
	Obj    Obj    // The object representing the string.
	Length int    // The length of the string.
	Chars  []byte // The byte array of characters composing the string.
	Hash   uint32 // The hash value of the string.
}

// NewFunction initializes and returns a new ObjFunction.
//
// No parameters.
// Returns a pointer to ObjFunction.
func NewFunction() *ObjFunction {
	function := &ObjFunction{}

	function.arity = 0
	function.obj = allocateObject(ObjFunctionType)
	function.chunk = Chunk{}
	InitChunk(&function.chunk)
	function.name = nil

	return function
}

// AsFunction returns the ObjFunction from the given Value.
//
// value Value
// *ObjFunction
func AsFunction(value Value) *ObjFunction {
	return value.As.(*ObjFunction)
}

// IsFunction checks if the given value is a function.
//
// value Value
// bool
func IsFunction(value Value) bool {
	return IsObjType(value, ObjFunctionType)
}

// OBJStrType returns the ObjType of the given Value.
//
// It takes a single parameter:
// - value: the Value to determine the ObjType for.
//
// It returns the ObjType of the given Value.
func OBJType(value Value) ObjType {
	return AsObj(value).Type
}

// IsObjType checks if the value is of a specific object type.
//
// Parameters:
// - value: The value to check.
// - objType: The object type to compare against.
//
// Returns:
// - bool: True if the value is of the specified object type, false otherwise.
func IsObjType(value Value, objType ObjType) bool {
	return IsValObj(value) && AsObjString(value).Obj.Type == objType
}

// IsString checks if the given value is a string.
//
// value: The value to check.
// Returns: A boolean indicating if the value is a string.
func IsString(value Value) bool {
	return IsObjType(value, ObjStringType)
}

// AsObjString returns the ObjectString representation of a Value.
//
// value: The Value to convert to an ObjectString.
// Returns: A pointer to the ObjectString representation of the Value.
func AsObjString(value Value) *ObjectString {
	return value.As.(*ObjectString)
}

// AsCString returns a string representation of a given Value as a C string.
//
// It takes a Value as a parameter.
// It returns a string.
func AsCString(value Value) string {
	objString := AsObjString(value).Chars
	return string(objString[:len(objString)-1])
}

// copyString is a function that creates a new ObjectString by copying a substring of a source string.
//
// It takes the starting index of the substring, the length of the substring, the source string, and the type of object as parameters.
// It returns a pointer to the newly created ObjectString.
func copyString(start, length int, _type ObjType) *ObjectString {
	source := *scanner.Source
	heapChars := make([]byte, length+1)
	chars := source[start : start+length]
	hash := hashString([]byte(chars), length)
	interned := tableFindString(vm.strings, []byte(chars), length, hash)
	if interned != nil {
		return interned
	}
	copy(heapChars, []byte(chars))
	return allocateString(heapChars, length, _type, hash)
}

// tableFindString finds a string in a table.
//
// Parameters:
// - table: a pointer to a Table object.
// - chars: a byte slice representing the characters to find.
// - length: an integer representing the length of the characters.
// - hash: an unsigned 32-bit integer representing the hash value.
//
// Returns:
// - a pointer to an ObjectString if the string is found, otherwise nil.
func tableFindString(table *Table, chars []byte, length int, hash uint32) *ObjectString {
	if table.count == 0 {
		return nil
	}
	index := hash % uint32(table.capacity)
	for {
		entry := table.entries[index]
		if entry.key == nil {
			if IsNil(entry.value) {
				return nil
			}
		} else if entry.key.Length == length && entry.key.Hash == hash && memcmp(entry.key.Chars, chars, length) == 0 {
			return entry.key
		}
		index = (index + 1) % uint32(table.capacity)
	}

}

// memcmp compares two byte slices up to a specified length and returns an integer indicating their order.
//
// The function takes three parameters:
//   - s1: the first byte slice to compare.
//   - s2: the second byte slice to compare.
//   - length: the number of bytes to compare.
//
// It returns an integer:
//   - 0 if the byte slices are equal up to the specified length.
//   - 1 if the byte slices are not equal up to the specified length, or if the specified length is invalid.
func memcmp(s1, s2 []byte, length int) int {
	// This example uses the built-in bytes.Equal function
	if length <= 0 {
		return 0
	}
	if length > len(s1) || length > len(s2) {
		return 1
	}
	if bytesEqual := bytes.Equal(s1[:length], s2[:length]); bytesEqual {
		return 0
	}
	return 1
}

// hashString calculates the hash value of a given byte array and returns the result as a uint32.
//
// Parameters:
// - key: the byte array to be hashed.
// - length: the length of the byte array.
//
// Returns:
// - uint32: the hash value of the byte array.
func hashString(key []byte, length int) uint32 {
	hash := uint32(2166136261)
	const prime = 16777619

	for i := 0; i < length; i++ {
		hash ^= uint32(key[i])
		hash *= prime
	}
	return hash
}

// allocateString creates a new ObjectString and initializes its properties.
//
// chars is a byte array representing the characters of the string.
// length is the number of characters in the string.
// _type is the type of the object.
// hash is the hash value of the string.
// Returns a pointer to the newly created ObjectString.
func allocateString(chars []byte, length int, _type ObjType, hash uint32) *ObjectString {
	str := &ObjectString{
		Length: length,
		Chars:  chars,
		Hash:   hash,
		Obj:    allocateObject(_type),
	}
	vm.strings.TableSet(str, NilValue())
	return str
}

// allocateObject creates a new object of the given type and adds it to the virtual machine's object list.
//
// _type: The type of the object to be created.
// Returns: The newly created object.
func allocateObject(_type ObjType) Obj {
	obj := Obj{
		Type: ObjType(_type),
	}
	obj.Next = vm.objects
	vm.objects = &obj
	return obj
}
