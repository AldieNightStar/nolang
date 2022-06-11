package nolang

import "testing"

func TestStack(t *testing.T) {
	s := NewStack[int](4)
	s.Push(32)
	s.Push(12)
	s.Push(14)
	s.Push(44)
	if s.Push(100) {
		t.Fatal("Stack should be overflown")
	}
	for _, n := range []int{44, 14, 12, 32} {
		val, ok := s.Pop(0)
		if !ok {
			t.Fatal("Can't pop values from Stack")
		}
		if n != val {
			t.Fatalf("%d should be %d", val, n)
		}
	}
	if _, ok := s.Pop(0); ok {
		t.Fatal("Empty stack should not return anything")
	}
}
