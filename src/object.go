package src

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
	return AsObj(value).Obj.Type
}

func IsObjType(value Value, objType ObjType) bool {
	return IsValObj(value) && AsObj(value).Obj.Type == objType
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
	copy(heapChars, source[start:start+length])
	return allocateString(heapChars, length, _type, hash)
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
