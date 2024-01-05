package src

import "fmt"

type ValueType int

const (
	ValBool ValueType = iota
	ValNil
	ValNumber
)

type Value struct {
	Type ValueType
	As   interface{}
}

type ValueArray struct {
	Capacity int
	Count    int
	Values   []Value
}

// Functions to create specific value instances
func BoolValue(value bool) Value {
	return Value{Type: ValBool, As: value}
}

func NilValue() Value {
	return Value{Type: ValNil, As: nil}
}

func NumberValue(value float64) Value {
	return Value{Type: ValNumber, As: value}
}

// Functions to access values and check types
func AsBool(value Value) bool {
	return value.As.(bool)
}

func AsNumber(value Value) float64 {
	return value.As.(float64)
}

func IsBool(value Value) bool {
	return value.Type == ValBool
}

func IsNil(value Value) bool {
	return value.Type == ValNil
}

func IsNumber(value Value) bool {
	return value.Type == ValNumber
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
	switch value.Type {
	case ValBool:
		fmt.Print(AsBool(value))
	case ValNil:
		fmt.Print("nil")
	case ValNumber:
		fmt.Print(AsNumber(value))
	}
}
