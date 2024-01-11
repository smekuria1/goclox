package src

import (
	"bytes"
)

type ObjType int

const (
	ObjStringType ObjType = iota
)

type Obj struct {
	Type ObjType
	Next *Obj
}

type ObjectString struct {
	Obj    Obj
	Length int
	Chars  []byte
	Hash   uint32
}

func OBJStrType(value Value) ObjType {
	return AsObj(value).Type
}

func IsObjType(value Value, objType ObjType) bool {
	return IsValObj(value) && AsObjString(value).Obj.Type == objType
}

func IsString(value Value) bool {
	return IsObjType(value, ObjStringType)
}

// Method to cast a Value to ObjectString
func AsObjString(value Value) *ObjectString {
	return value.As.(*ObjectString)
}

// Method to get the Chars field of an ObjectString
func AsCString(value Value) string {
	objString := AsObjString(value).Chars
	return string(objString[:len(objString)-1])
}

func copyString(start, length int, source string, _type ObjType) *ObjectString {
	heapChars := make([]byte, length+1)
	chars := source[start : start+length]
	hash := hashString([]byte(chars), length)
	interned := tableFindString(&vm.strings, []byte(chars), length, hash)
	if interned != nil {
		return interned
	}
	copy(heapChars, []byte(chars))
	return allocateString(heapChars, length, _type, hash)
}

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

func hashString(key []byte, length int) uint32 {
	hash := uint32(2166136261)
	const prime = 16777619

	for i := 0; i < length; i++ {
		hash ^= uint32(key[i])
		hash *= prime
	}
	return hash
}

func allocateString(chars []byte, length int, _type ObjType, hash uint32) *ObjectString {
	str := &ObjectString{
		Length: length,
		Chars:  chars,
		Hash:   hash,
		Obj:    allocateObject(_type),
	}
	TableSet(&vm.strings, str, NilValue())
	return str
}

func allocateObject(_type ObjType) Obj {
	obj := Obj{
		Type: ObjType(_type),
	}
	obj.Next = vm.objects
	vm.objects = &obj
	return obj
}
