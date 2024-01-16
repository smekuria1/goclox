package src

import (
	"log"
	"reflect"
)

var logger = log.Default()

// GrowCapacity returns the new capacity after growing the old capacity.
//
// oldcap - the old capacity (int)
// Returns the new capacity (int)
func GrowCapacity(oldcap int) int {
	if oldcap < 8 {
		return 8
	}

	return oldcap * 2
}

// GrowArrayChunks is a Go function that takes in a code slice of uint8, along with the oldcap and newcap as integers.
//
// It returns a new code slice of the same type ([]uint8).
func GrowArrayChunks(code []uint8, oldcap, newcap int) []uint8 {
	return Reallocate(code, oldcap, newcap).([]uint8)
}

// GrowArrayValueArray returns a new array of Values with a larger capacity.
//
// Parameters:
// - valarray: The original array of Values.
// - oldcap: The old capacity of the array.
// - newcap: The new capacity of the array.
//
// Return type:
// - []Value: The new array of Values.
func GrowArrayValueArray(valarray []Value, oldcap, newcap int) []Value {
	return Reallocate(valarray, oldcap, newcap).([]Value)
}

// GrowArrayLines returns a new slice with the same elements as the original slice, but with a larger capacity.
//
// It takes in three parameters:
// - lines: a slice of integers representing the original array
// - oldcap: an integer representing the old capacity of the array
// - newcap: an integer representing the new capacity of the array
//
// It returns a new slice of integers with the same elements as the original slice, but with a larger capacity.
func GrowArrayLines(lines []int, oldcap, newcap int) []int {
	return Reallocate(lines, oldcap, newcap).([]int)
}

// GrowArrayEntries returns a new slice of Entry with a larger capacity.
//
// The function takes in the following parameters:
// - entries: a slice of Entry that contains the current entries.
// - oldcap: an integer that represents the old capacity of the slice.
// - newcap: an integer that represents the new capacity of the slice.
//
// The function returns a new slice of Entry with the updated capacity.
func GrowArrayEntries(entries []Entry, oldcap, newcap int) []Entry {
	return Reallocate(entries, oldcap, newcap).([]Entry)
}

// FreeArray releases the memory occupied by the given array.
//
// The function takes two parameters:
// - `array` of type `any`, which is the array to be freed.
// - `cap` of type `int`, which is the capacity of the array.
//
// The function does not return anything.
func FreeArray(array any, cap int) {
	Reallocate(array, cap, 0)
}

// Reallocate reallocates the memory of a pointer to a new size.
//
// It takes in three parameters:
// - pointer: the pointer to reallocate the memory for.
// - oldSize: the current size of the memory block.
// - newSize: the new size of the memory block.
//
// It returns an interface{} which is the reallocated pointer.
func Reallocate(pointer interface{}, oldSize, newSize int) interface{} {
	oldptrvalue := reflect.ValueOf(pointer)
	if newSize == 0 {
		return reflect.Zero(reflect.TypeOf(oldptrvalue)).Interface()
	}
	newptrvalue := reflect.MakeSlice(oldptrvalue.Type(), newSize, newSize)
	if newptrvalue.IsNil() {
		logger.Fatalln("Realloc eror")
	}
	//fmt.Println("Growing Array")
	reflect.Copy(newptrvalue, oldptrvalue)

	return newptrvalue.Interface()

}

// FreeObjects frees all objects in the linked list starting from the given object.
//
// object: a pointer to the first object in the linked list.
func FreeObjects(object *Obj) {
	for object != nil {
		next := object.Next
		object.Next = nil
		object = next
	}
}
