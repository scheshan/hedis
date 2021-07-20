package session

type SessionList struct {
	size int
	head *Session
	tail *Session
}

func NewSessionList() *SessionList {
	list := &SessionList{}

	return list
}

func (t *SessionList) AddLast(s *Session) {
	if t.tail == nil {
		t.head = s
	} else {
		t.tail.next = s
		s.pre = t.tail
	}
	t.tail = s
	t.size++
}

func (t *SessionList) Remove(s *Session) {
	if s.pre != nil {
		s.pre.next = s.next
	}
	if s.next != nil {
		s.next.pre = s.pre
	}

	if t.head == s {
		t.head = s.next
	}
	if t.tail == s {
		t.tail = s.pre
	}
	t.size--
}
