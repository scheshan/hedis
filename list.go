package hedis

type List struct {
	head *listNode
	tail *listNode
	size int
}

type listNode struct {
	prev *listNode
	next *listNode
	v    interface{}
}

func NewList() *List {
	return new(List)
}

func (t *List) AddFirst(v interface{}) {
	node := new(listNode)
	node.v = v

	if t.head != nil {
		t.head.prev = node
	}
	node.next = t.head.prev
	t.head = node

	if t.tail == nil {
		t.tail = node
	}

	t.size++
}

func (t *List) AddLast(v interface{}) {
	node := new(listNode)
	node.v = v

	if t.tail != nil {
		t.tail.next = node
	}
	node.prev = t.tail
	t.tail = node

	if t.head == nil {
		t.head = node
	}

	t.size++
}

func (t *List) RemoveFirst() (v interface{}, find bool) {
	v = nil
	find = false

	if t.size == 0 {
		return
	}

	node := t.head
	v = node.v
	find = true

	if node.next != nil {
		node.next.prev = nil
	}
	t.head = node.next

	if node == t.tail {
		t.tail = nil
	}

	t.size--
	return
}

func (t *List) RemoveLast() (v interface{}, find bool) {
	v = nil
	find = false

	if t.size == 0 {
		return
	}

	node := t.tail
	v = node.v
	find = true

	if node.prev != nil {
		node.prev.next = nil
	}
	t.tail = node.prev

	if node == t.head {
		t.head = nil
	}

	t.size--
	return
}

func (t *List) Remove(v interface{}) {
	if t.size == 0 {
		return
	}

	n := t.head
	for n != nil {
		if n.v == v {
			if n.prev != nil {
				n.prev.next = n.next
			}
			if n.next != nil {
				n.next.prev = n.prev
			}
			if n == t.head {
				t.head = n.next
			}
			if n == t.tail {
				t.tail = n.prev
			}

			t.size--
		}

		n = n.next
	}
}
