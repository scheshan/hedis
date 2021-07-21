package session

type State int

const (
	StateNonAuth State = iota
	StateOpen
	StateError
	StateClosed
)
