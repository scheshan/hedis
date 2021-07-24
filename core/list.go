package core

type ListFilter func(value interface{}) bool

type List struct {
	size int
	head *ListItem
	tail *ListItem
}

type ListItem struct {
	value interface{}
	pre   *ListItem
	next  *ListItem
}

func (t *List) AddHead(value interface{}) {
	li := &ListItem{
		value: value,
	}

	if t.head == nil {
		t.head = li
		t.tail = li
	} else {
		li.next = t.head
		t.head.pre = li
		t.head = li
	}

	t.size++
}

func (t *List) AddTail(value interface{}) {
	li := &ListItem{
		value: value,
	}

	if t.tail == nil {
		t.head = li
		t.tail = li
	} else {
		li.pre = t.tail
		t.tail.next = li
		t.tail = li
	}

	t.size++
}

func (t *List) GetHead() (interface{}, *ListItem, bool) {
	if t.head == nil {
		return nil, nil, false
	}

	return t.head.value, t.head, true
}

func (t *List) GetTail() (interface{}, *ListItem, bool) {
	if t.tail == nil {
		return nil, nil, false
	}

	return t.tail.value, t.tail, true
}

func (t *List) Remove(item *ListItem) {
	if item.pre != nil {
		item.pre.next = item.next
	}
	if item.next != nil {
		item.next.pre = item.pre
	}
	if item == t.head {
		t.head = t.head.next
	}
	if item == t.tail {
		t.tail = t.tail.pre
	}

	item.pre = nil
	item.next = nil
	t.size--
}

func (t *List) Len() int {
	return t.size
}

func (t *List) Filter(filter ListFilter) (interface{}, *ListItem, bool) {
	cur := t.head
	for cur != nil {
		if filter(cur.value) {
			return cur.value, cur, true
		}
	}

	return nil, nil, false
}

func NewList() *List {
	return &List{}
}
