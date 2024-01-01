package src

import "fmt"

type Value float64

type ValueArray struct {
	Capacity int
	Count    int
	Values   []Value
}

func InitValueArray(array *ValueArray) {
	array.Values = nil
	array.Capacity = 0
	array.Count = 0
}

func WriteValueArray(array *ValueArray, val Value) {
	if array.Capacity < array.Count+1 {
		oldCap := array.Capacity
		array.Capacity = GrowCapacity(oldCap)
		array.Values = GrowArrayValueArray(array.Values, oldCap, array.Capacity)
	}
	array.Values[array.Count] = val
	array.Count++
}

func FreeValueArray(array *ValueArray) {
	FreeArray(array.Values, array.Capacity)
	InitValueArray(array)
}

func PrintValue(value Value) {
	fmt.Printf("%g", value)
}
