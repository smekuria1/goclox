package src

/*Add support for keys of the other primitive types: numbers, Booleans, and nil.
Later, clox will support user-defined classes. If we want to support keys that are
instances of those classes, what kind of complexity does that add
*/

// TableMaxLoad is the maximum load factor for the table.
const TableMaxLoad float32 = 0.75

// Table is a struct representing a table data structure.
type Table struct {
	capacity float32 // The maximum capacity of the table.
	count    float32 // The current count of entries in the table.
	entries  []Entry // The array of entries in the table.
}

// Entry is a struct representing an entry in the table.
type Entry struct {
	key   *ObjectString // The key of the entry, represented as a pointer to an ObjectString
	value Value         // The value of the entry
}

// InitTable initializes the Table struct.
//
// No parameters.
// No return values.
func (table *Table) InitTable() {
	table.capacity = 0
	table.count = 0
	table.entries = nil
}

// Freetable frees the table by releasing its entries and reinitializing it.
//
// No parameters.
// No return types.
func (table *Table) Freetable() {
	FreeArray(table.entries, int(table.capacity))
	table.InitTable()
}

// TableSet sets the value for a given key in the Table.
//
// Parameters:
// - key: a pointer to an ObjectString representing the key.
// - value: the value to be set.
//
// Returns:
// - bool: true if the key is a new key in the table, false otherwise.
func (table *Table) TableSet(key *ObjectString, value Value) bool {
	if table.count+1 > table.capacity*TableMaxLoad {
		oldcap := table.capacity
		capacity := GrowCapacity(int(table.capacity))
		table.adjustTable(int(oldcap), capacity)
	}
	entry, index := table.findEntry(int(table.capacity), key)
	isNewKey := entry.key == nil
	if isNewKey && IsNil(entry.value) {
		table.count++
	}
	table.entries[index].key = key
	table.entries[index].value = value
	return isNewKey
}

// TableDelete deletes an entry from the Table.
//
// It takes a key of type *ObjectString as a parameter and returns a boolean value indicating whether the deletion was successful.
func (table *Table) TableDelete(key *ObjectString) bool {
	if table.count == 0 {
		return false
	}
	entry, _ := table.findEntry(int(table.capacity), key)
	if entry.key == nil {
		return false
	}
	entry.key = nil
	entry.value = BoolValue(true)

	return true
}

// TableAddAll adds all elements from the given table to the current table.
//
// It takes a pointer to the Table struct named "from" as a parameter.
// It does not return anything.
func (table *Table) TableAddAll(from *Table) {
	for i := 0; i < int(from.capacity); i++ {
		entry := from.entries[i]
		if entry.key != nil {
			table.TableSet(entry.key, entry.value)
		}
	}
}

// TableGet retrieves the value associated with the given key in the Table.
//
// The function takes two parameters: key, a pointer to an ObjectString, and value, a pointer to a Value.
// It returns a boolean value indicating whether the key was found in the Table.
func (table Table) TableGet(key *ObjectString, value *Value) bool {
	if table.count == 0 {
		return false
	}

	entry, _ := table.findEntry(int(table.capacity), key)
	if entry.key == nil {
		return false
	}
	*value = entry.value
	return true
}

// findEntry finds the entry in the table with the given capacity and key.
//
// Parameters:
// - capacity: the capacity of the table
// - key: the key to search for
//
// Returns:
// - entry: the entry found
// - index: the index of the entry
func (table *Table) findEntry(capacity int, key *ObjectString) (*Entry, uint32) {
	index := key.Hash % uint32(capacity)
	var tombstone *Entry
	for {
		entry := table.entries[index]
		if entry.key == nil {
			if IsNil(entry.value) {
				if tombstone != nil {
					return tombstone, index
				}
				return &entry, index
			}
			if tombstone == nil {
				tombstone = &entry
			}

		} else if entry.key == key {
			return &entry, index
		}
		index = (index + 1) % uint32(capacity)
	}
}

func (table *Table) adjustTable(oldcap, capacity int) {
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

		dest, _ := table.findEntry(capacity, entry.key)
		dest.key = entry.key
		dest.value = entry.value
		table.count++

	}
	FreeArray(table.entries, int(table.capacity))
	table.entries = entries
	table.capacity = float32(capacity)
}
