package hedis

type Dict struct {
	dt        []*dictTable
	rehashIdx int32
}

type dictTable struct {
	table    []*dictEntry
	used     uint32
	size     uint32
	sizeMask uint32
}

type dictEntry struct {
	key   string
	value interface{}
	next  *dictEntry
}

func NewDict() *Dict {
	d := new(Dict)
	d.dt = make([]*dictTable, 2, 2)
	d.rehashIdx = -1

	return d
}

func (t *Dict) Put(k string, v interface{}) {
	t.expand()

	if t.isRehashing() {
		t.rehash(1)
	}

	de := t.find(k)
	if de != nil {
		de.value = v
		return
	}

	dt := t.dt[0]
	if t.isRehashing() {
		dt = t.dt[1]
	}

	idx := t.hashCode(k) & dt.sizeMask

	head := dt.table[idx]

	newHead := new(dictEntry)
	newHead.key = k
	newHead.value = v
	newHead.next = head

	dt.table[idx] = newHead
	dt.used++
}

func (t *Dict) Remove(k string) bool {
	if t.isRehashing() {
		t.rehash(1)
	}

	hc := t.hashCode(k)
	find := false

	for i := 0; i <= 1; i++ {
		dt := t.dt[i]
		if dt == nil {
			continue
		}

		idx := hc & dt.sizeMask
		node := dt.table[idx]
		var pre *dictEntry = nil
		for node != nil {
			if node.key == k {
				if pre != nil {
					pre.next = node.next
				} else {
					dt.table[idx] = node.next
				}
				find = true
				break
			}

			node = node.next
		}
	}

	return find
}

func (t *Dict) Get(k string) (v interface{}, find bool) {
	de := t.find(k)
	if de == nil {
		return nil, false
	}

	return de.value, true
}

func (t *Dict) Exists(k string) bool {
	de := t.find(k)

	return de != nil
}

func (t *Dict) hashCode(key string) uint32 {
	var h uint32 = 5381

	for _, str := range key {
		h = ((h << 5) + h) + uint32(str)
	}

	return h
}

func (t *Dict) isRehashing() bool {
	return t.rehashIdx >= 0
}

func (t *Dict) find(k string) *dictEntry {
	hc := t.hashCode(k)
	for i := 0; i <= 1; i++ {
		dt := t.dt[i]
		if dt == nil {
			continue
		}

		idx := hc & dt.sizeMask

		head := dt.table[idx]
		for head != nil {
			if head.key == k {
				return head
			}
			head = head.next
		}
	}

	return nil
}

func (t *Dict) expand() {
	if t.isRehashing() {
		return
	}

	dt := t.dt[0]

	size := DictInitSize
	if dt != nil && dt.size > 0 {
		if dt.used < dt.size || dt.size >= DictMaxSize {
			return
		}

		size = dt.size << 1
		if size > DictMaxSize {
			size = DictMaxSize
		}
	}

	dt = new(dictTable)
	dt.size = size
	dt.sizeMask = size - 1
	dt.table = make([]*dictEntry, size)

	if t.dt[0] == nil {
		t.dt[0] = dt
	} else {
		t.dt[1] = dt
		t.rehashIdx = 0
	}
}

func (t *Dict) rehash(n int) {
	if !t.isRehashing() {
		return
	}

	for n > 0 {
		if t.dt[0].used == 0 {
			t.dt[0] = t.dt[1]
			t.dt[1] = nil
			t.rehashIdx = -1
			return
		}

		for t.dt[0].table[t.rehashIdx] == nil {
			t.rehashIdx++
			if uint32(t.rehashIdx) >= t.dt[0].size {
				return
			}
		}

		head := t.dt[0].table[t.rehashIdx]
		for head != nil {
			node := head
			head = node.next

			idx := t.hashCode(node.key) & t.dt[1].sizeMask
			node.next = t.dt[1].table[idx]
			t.dt[1].table[idx] = node
			t.dt[1].used++
			t.dt[0].used--
			t.dt[0].table[t.rehashIdx] = nil
		}

		t.rehashIdx++
		n--
	}
}
