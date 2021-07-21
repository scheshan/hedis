package core

const maxHashTableSize = 1 << 30

type Hash struct {
	items []*hashItem
}

type hashItem struct {
	hash  int
	key   *String
	value interface{}
	pre   *hashItem
	next  *hashItem
}

func (t *Hash) tableSize(size int) int {
	if size > maxHashTableSize {
		return maxHashTableSize
	}

	n := size - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16

	if n < 0 {
		return 1
	}

	if n >= maxHashTableSize {
		return maxHashTableSize
	}

	return n + 1
}

func (t *Hash) find(key *String) *hashItem {
	hc := key.HashCode()
	ind := hc & (len(t.items) - 1)

	cur := t.items[ind]
	for cur != nil {
		if cur.hash == hc && cur.key.Equal(key) {
			return cur
		}

		cur = cur.next
	}

	return nil
}

func (t *Hash) Contains(key *String) bool {
	return t.find(key) != nil
}

func (t *Hash) Get(key *String) interface{} {
	item := t.find(key)
	if item == nil {
		return nil
	}

	return item.value
}

func NewHashSize(size int) *Hash {
	h := &Hash{}

	size = h.tableSize(size)
	h.items = make([]*hashItem, size, size)

	return h
}
