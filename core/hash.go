package core

import "math/rand"

const maxHashTableSize = 1 << 30
const maxHashSize = 1 << 31

type HashFunc func(key *String, value interface{})

var HashDefaultValue = true

//Hash	哈希数据结构
//
//内部维护了2个数组，t1是主要存储元素的数组，t2用来做扩容迁移。当容量达到阈值后，Hash触发rehash。
//此时会递进的将t1元素迁移到t2元素。tIndex表示了当前迁移的数组下标，当tIndex达到len(t1)后，迁移
//结束。
type Hash struct {
	size   int
	tIndex int
	t1     []*hashItem
	t2     []*hashItem
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

func (t *Hash) findInTable(key *String, table []*hashItem) (int, *hashItem) {
	h := key.HashCode()
	ind := h & (len(table) - 1)

	cur := table[ind]
	for cur != nil {
		if cur.hash == h && cur.key.Equal(key) {
			return ind, cur
		}

		cur = cur.next
	}

	return -1, nil
}

func (t *Hash) isTransferring() bool {
	return t.tIndex > -1
}

func (t *Hash) find(key *String) *hashItem {
	_, hi := t.findInTable(key, t.t1)

	if hi == nil && t.isTransferring() {
		_, hi = t.findInTable(key, t.t2)
	}

	return hi
}

func (t *Hash) transfer() {
	if !t.isTransferring() {
		return
	}

	for t.tIndex < len(t.t1) && t.t1[t.tIndex] == nil {
		t.tIndex++
	}

	if t.tIndex < len(t.t1) {
		cur := t.t1[t.tIndex]
		for cur != nil {
			ind := cur.hash&len(t.t2) - 1

			next := cur.next

			cur.pre = nil
			cur.next = t.t2[ind]
			if t.t2[ind] != nil {
				t.t2[ind].pre = cur
			}
			t.t2[ind] = cur

			cur = next
		}

		t.tIndex++
	}

	if t.tIndex == len(t.t1) {
		t.tIndex = -1
		t.t1 = t.t2
		t.t2 = nil
	}
}

func (t *Hash) ensureSize() {
	if t.isTransferring() {
		return
	}

	if len(t.t1) == maxHashTableSize {
		return
	}

	if t.size >= len(t.t1)*3/4 {
		newSize := len(t.t1) << 1

		t.t2 = make([]*hashItem, newSize, newSize)
		t.tIndex = 0
	}
}

func (t *Hash) removeInTable(key *String, table []*hashItem) bool {
	ind, hi := t.findInTable(key, table)
	if hi == nil {
		return false
	}

	if hi.pre != nil {
		hi.pre.next = hi.next
	}
	if hi.next != nil {
		hi.next.pre = hi.pre
	}

	if table[ind] == hi {
		table[ind] = hi.next
	}

	hi.pre = nil
	hi.next = nil

	t.size--

	return true
}

func (t *Hash) iterate(ind int, table []*hashItem, hashFunc HashFunc) {
	head := table[ind]
	for head != nil {
		hashFunc(head.key, head.value)
		head = head.next
	}
}

func (t *Hash) iterateTable(table []*hashItem, hashFunc HashFunc) {
	for i := 0; i < len(table); i++ {
		t.iterate(i, table, hashFunc)
	}
}

func (t *Hash) randomItem(head *hashItem) *hashItem {
	num := 0
	cur := head
	for cur != nil {
		num++
		cur = cur.next
	}

	num = rand.Intn(num)
	cur = head
	for num > 0 {
		cur = cur.next
		num--
	}

	return cur
}

func (t *Hash) Contains(key *String) bool {
	t.transfer()

	return t.find(key) != nil
}

func (t *Hash) Get(key *String) (interface{}, bool) {
	t.transfer()

	item := t.find(key)
	if item == nil {
		return nil, false
	}

	return item.value, true
}

func (t *Hash) Put(key *String, value interface{}) {
	t.transfer()

	if t.size >= maxHashSize {
		panic("hash memory overflow")
	}

	hi := t.find(key)
	if hi != nil {
		hi.value = value
		return
	}

	hi = &hashItem{}
	hi.hash = key.HashCode()
	hi.key = key
	hi.value = value

	tb := t.t1
	if t.isTransferring() {
		tb = t.t2
	}

	ind := hi.hash & (len(tb) - 1)
	hi.next = tb[ind]
	if tb[ind] != nil {
		tb[ind].pre = hi
	}
	tb[ind] = hi

	t.size++

	t.ensureSize()

	return
}

func (t *Hash) Remove(key *String) bool {
	t.transfer()

	if t.removeInTable(key, t.t1) {
		return true
	}

	if t.isTransferring() && t.removeInTable(key, t.t2) {
		return true
	}

	return false
}

func (t *Hash) Iterate(hashFunc HashFunc) {
	t.iterateTable(t.t1, hashFunc)

	if t.isTransferring() {
		t.iterateTable(t.t2, hashFunc)
	}
}

func (t *Hash) Size() int {
	return t.size
}

func (t *Hash) Random() (key *String, value interface{}, find bool) {
	if t.size == 0 {
		return nil, nil, false
	}

	var head *hashItem
	if t.isTransferring() {
		for head == nil {
			ind := t.tIndex + rand.Intn(len(t.t1)+len(t.t2)-t.tIndex)
			if ind < len(t.t1) {
				head = t.t1[ind]
			} else {
				head = t.t2[ind-len(t.t1)]
			}
		}
	} else {
		for head == nil {
			ind := rand.Intn(len(t.t1))
			head = t.t1[ind]
		}
	}
	hi := t.randomItem(head)

	return hi.key, hi.value, true
}

func NewHashSize(size int) *Hash {
	h := &Hash{}

	size = h.tableSize(size)
	h.t1 = make([]*hashItem, size, size)

	h.tIndex = -1

	return h
}

func NewHash() *Hash {
	return NewHashSize(16)
}
