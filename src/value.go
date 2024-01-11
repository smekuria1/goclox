package src

import (
	"bytes"
	"fmt"
)

type ValueType int

const (
	ValBool ValueType = iota
	ValNil
	ValObjStr
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

func ObjStrValue(value *ObjectString) Value {
	return Value{Type: ValObjStr, As: value}
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
func AsObj(value Value) *Obj {
	return value.As.(*Obj)
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

func IsValObj(value Value) bool {
	return value.Type == ValObjStr
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
	case ValObjStr:
		printObjectStr(value)
	}
}

func printObjectStr(object Value) {
	switch object.As.(*ObjectString).Obj.Type {
	case ObjStringType:
		fmt.Printf("%s", AsCString(object))
	}
}

func valuesEqual(a, b Value) bool {
	if a.Type != b.Type {
		return false
	}

	switch a.Type {
	case ValBool:
		return AsBool(a) == AsBool(b)
	case ValNil:
		return true
	case ValNumber:
		return AsNumber(a) == AsNumber(b)
	case ValObjStr:
		aString := removeNullBytes(AsObjString(a).Chars)
		bString := removeNullBytes(AsObjString(b).Chars)

		return bytes.Equal(aString, bString)
	}

	return false
}

func removeNullBytes(input []byte) []byte {
	return bytes.Replace(input, []byte{0}, []byte{}, -1)
}
