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

func (table *Table) InitTable() {
	table.capacity = 0
	table.count = 0
	table.entries = nil
}

func (table *Table) Freetable() {
	FreeArray(table.entries, int(table.capacity))
	table.InitTable()
}

func (table *Table) TableSet(key *ObjectString, value Value) bool {
	if table.count+1 > table.capacity*TABLE_MAX_LOAD {
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

func (to *Table) TableAddAll(from *Table) {
	for i := 0; i < int(from.capacity); i++ {
		entry := from.entries[i]
		if entry.key != nil {
			to.TableSet(entry.key, entry.value)
		}
	}
}

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

func (table *Table) findEntry(capacity int, key *ObjectString) (*Entry, uint32) {
	index := key.Hash % uint32(capacity)
	var tombstone *Entry = nil
	for {
		entry := table.entries[index]
		if entry.key == nil {
			if IsNil(entry.value) {
				if tombstone != nil {
					return tombstone, index
				}
				return &entry, index
			} else {
				if tombstone == nil {
					tombstone = &entry
				}
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
