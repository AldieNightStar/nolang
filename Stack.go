package nolang

type Stack[T any] struct {
	Data []T
	Ptr  int
}

func NewStack[T any](size int) *Stack[T] {
	return &Stack[T]{Data: make([]T, size), Ptr: 0}
}

func (s *Stack[T]) Push(dat T) bool {
	if s.Ptr >= len(s.Data) {
		return false
	}
	s.Data[s.Ptr] = dat
	s.Ptr += 1
	return true
}

func (s *Stack[T]) Pop(def T) (T, bool) {
	if s.Ptr <= 0 {
		return def, false
	}
	s.Ptr -= 1
	return s.Data[s.Ptr], true
}
