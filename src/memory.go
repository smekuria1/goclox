package src

import (
	"log"
	"reflect"
)

var logger = log.Default()

func GrowCapacity(oldcap int) int {
	if oldcap < 8 {
		return 8
	}

	return oldcap * 2
}

func GrowArrayChunks(code []uint8, oldcap, newcap int) []uint8 {
	return Reallocate(code, oldcap, newcap).([]uint8)
}

func GrowArrayValueArray(valarray []Value, oldcap, newcap int) []Value {
	return Reallocate(valarray, oldcap, newcap).([]Value)
}

func GrowArrayLines(lines []int, oldcap, newcap int) []int {
	return Reallocate(lines, oldcap, newcap).([]int)
}

func GrowArrayEntries(entries []Entry, oldcap, newcap int) []Entry {
	return Reallocate(entries, oldcap, newcap).([]Entry)
}

func FreeArray(array any, cap int) {
	Reallocate(array, cap, 0)
}

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

func FreeObjects(object *Obj) {
	for object != nil {
		next := object.Next
		object.Next = nil
		object = next
	}
}
