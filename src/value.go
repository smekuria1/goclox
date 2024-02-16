package src

import (
	"bytes"
	"fmt"
)

type ValueType int

// Value types
const (
	ValBool ValueType = iota
	ValNil
	ValObjStr
	ValObj
	ValNumber
)

// Value represents a value in the language
type Value struct {
	Type ValueType   // The type of the value
	As   interface{} // The value
}

// ValueArray represents an array of values.
type ValueArray struct {
	Capacity int     // The maximum number of elements that can be stored in the array.
	Count    int     // The current number of elements in the array.
	Values   []Value // The array of values.
}

// BoolValue returns a Value with Type ValBool and As value.
func BoolValue(value bool) Value {
	return Value{Type: ValBool, As: value}
}

// ObjStrValue returns a Value with Type ValObjStr and As value.
//
// value: a pointer to an ObjectString.
// Returns: a Value.
func ObjStrValue(value *ObjectString) Value {
	return Value{Type: ValObjStr, As: value}
}

// ObjFunctionValue returns the value of the ObjFunction.
//
// value *ObjFunction - the ObjFunction parameter
// Value - the return type
func ObjFunctionValue(value *ObjFunction) Value {
	return Value{Type: ValObj, As: value}
}

// OBJ_VAL description of the Go function.
//
// It takes a parameter object of type *Obj and returns a Value type.
func ObjVal(object *ObjFunction) Value {
	return Value{Type: ValueType(ValObj), As: object}
}

// NilValue returns a Value with Type ValNil and a nil As field.
//
// NilValue does not take any parameters.
// It returns a Value.
func NilValue() Value {
	return Value{Type: ValNil, As: nil}
}

// NumberValue creates a Value struct with the given float64 value.
//
// Parameters:
// - value: The float64 value to be assigned to the Value struct.
//
// Returns:
// The created Value struct.
func NumberValue(value float64) Value {
	return Value{Type: ValNumber, As: value}
}

// AsBool returns the boolean value of the given Value.
//
// It takes a single parameter:
// - value: the Value to convert to a boolean.
//
// It returns a boolean value.
func AsBool(value Value) bool {
	return value.As.(bool)
}

// AsObj returns the *Obj value from the given Value.
//
// It takes a single parameter:
// - value: the Value to extract the *Obj from.
//
// It returns a *Obj.
func AsObj(value Value) *Obj {
	return value.As.(*Obj)
}

// AsNumber returns the value of the input parameter as a float64.
//
// value: The value to be converted.
// Returns: The value as a float64.
func AsNumber(value Value) float64 {
	return value.As.(float64)
}

// IsBool checks if the given value is a boolean.
//
// value: the value to be checked.
// bool: true if the value is a boolean, false otherwise.
func IsBool(value Value) bool {
	return value.Type == ValBool
}

// IsNil checks if the given value is nil.
//
// It takes a parameter of type Value and returns a boolean value.
func IsNil(value Value) bool {
	return value.Type == ValNil
}

// IsValObj checks if the given value is of type ValObj
//
// value: the value to be checked.
// Returns: true if the value is of type ValObjStr, false otherwise.
func IsValObj(value Value) bool {
	return value.Type == ValObj
}

// IsNumber checks if the given value is of type number.
//
// value: the value to be checked.
// bool: true if the value is of type number, false otherwise.
func IsNumber(value Value) bool {
	return value.Type == ValNumber
}

// InitValueArray initializes the given ValueArray.
//
// Takes a pointer to a ValueArray as a parameter.
// Does not return anything.
func InitValueArray(array *ValueArray) {
	array.Values = nil
	array.Capacity = 0
	array.Count = 0
}

// WriteValueArray writes a value to the given ValueArray.
//
// It takes a pointer to a ValueArray and a Value as parameters.
// The ValueArray is dynamically resized if its capacity is not sufficient to accommodate the new value.
// After writing the value, the count of the ValueArray is incremented by 1.
func WriteValueArray(array *ValueArray, val Value) {
	if array.Capacity < array.Count+1 {
		oldCap := array.Capacity
		array.Capacity = GrowCapacity(oldCap)
		array.Values = GrowArrayValueArray(array.Values, oldCap, array.Capacity)
	}
	array.Values[array.Count] = val
	array.Count++
}

// FreeValueArray frees the memory allocated for the values in the given ValueArray.
//
// It takes a pointer to a ValueArray as a parameter.
// There is no return value.
func FreeValueArray(array *ValueArray) {
	FreeArray(array.Values, array.Capacity)
	InitValueArray(array)
}

// PrintValue prints the value of a given Value object.
//
// It takes a Value object as a parameter and prints its value based on its type:
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
	case ValObj:
		printFunction(AsFunction(value))
	}
}

// printObjectStr prints the string representation of an object.
//
// It takes a Value object as a parameter.
// It does not return anything.
func printObjectStr(object Value) {
	fmt.Printf("%s", AsCString(object))
}

// printFunction prints the function object.
func printFunction(function *ObjFunction) {

	if function.name == nil {
		fmt.Printf("<script>")
	} else {
		fmt.Printf("%s", string(function.name.Chars))
	}
}

// valuesEqual checks if two values are equal.
//
// It takes two parameters, a and b, of type Value.
// It returns a boolean value indicating whether the values are equal.
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

// removeNullBytes removes null bytes from the given input byte array.
//
// It takes an input byte array as a parameter and returns a modified byte array
// with all null bytes removed.
func removeNullBytes(input []byte) []byte {
	return bytes.Replace(input, []byte{0}, []byte{}, -1)
}
