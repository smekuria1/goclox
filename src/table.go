package src

/*Add support for keys of the other primitive types: numbers, Booleans, and nil.
Later, clox will support user-defined classes. If we want to support keys that are
instances of those classes, what kind of complexity does that add
*/
const TABLE_MAX_LOAD float32 = 0.75

type Table struct {
	capacity float32
	count    float32
	entries  []Entry
}

type Entry struct {
	key   *ObjectString
	value Value
}

func InitTable(table *Table) {
	table.capacity = 0
	table.count = 0
	table.entries = nil
}

func Freetable(table *Table) {
	FreeArray(table.entries, int(table.capacity))
	InitTable(table)
}

func TableSet(table *Table, key *ObjectString, value Value) bool {
	if table.count+1 > table.capacity*TABLE_MAX_LOAD {
		oldcap := table.capacity
		capacity := GrowCapacity(int(table.capacity))
		adjustTable(table, int(oldcap), capacity)
	}
	entry := findEntry(table.entries, int(table.capacity), key)
	isNewKey := entry.key == nil
	if isNewKey && IsNil(entry.value) {
		table.count++
	}
	entry.key = key
	entry.value = value
	return isNewKey
}

func TableDelete(table *Table, key *ObjectString) bool {
	if table.count == 0 {
		return false
	}
	entry := findEntry(table.entries, int(table.capacity), key)
	if entry.key == nil {
		return false
	}
	entry.key = nil
	entry.value = BoolValue(true)

	return true
}

func TableAddAll(from, to *Table) {
	for i := 0; i < int(from.capacity); i++ {
		entry := from.entries[i]
		if entry.key != nil {
			TableSet(to, entry.key, entry.value)
		}
	}
}

func TableGet(table *Table, key *ObjectString, value *Value) bool {
	if table.count == 0 {
		return false
	}

	entry := findEntry(table.entries, int(table.capacity), key)
	if entry.key == nil {
		return false
	}
	*value = entry.value
	return true
}

func findEntry(entries []Entry, capacity int, key *ObjectString) *Entry {
	index := key.Hash % uint32(capacity)
	var tombstone *Entry = nil
	for {
		entry := entries[index]
		if entry.key == nil {
			if IsNil(entry.value) {
				if tombstone != nil {
					return tombstone
				}
				return &entry
			} else {
				if tombstone == nil {
					tombstone = &entry
				}
			}
		} else if entry.key == key {
			return &entry
		}
		index = (index + 1) % uint32(capacity)
	}
}

func adjustTable(table *Table, oldcap, capacity int) {
	entries := GrowArrayEntries(table.entries, oldcap, capacity)

	for i := 0; i < capacity; i++ {
		entries[i].key = nil
		entries[i].value = NilValue()
	}
	table.count = 0
	for i := 0; i < int(table.capacity); i++ {
		entry := table.entries[i]
		if entry.key == nil {
			continue
		}

		dest := findEntry(entries, capacity, entry.key)
		dest.key = entry.key
		dest.value = entry.value
		table.count++

	}
	FreeArray(table.entries, int(table.capacity))
	table.entries = entries
	table.capacity = float32(capacity)
}
