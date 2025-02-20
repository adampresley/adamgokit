package datastructures

/*
Stack implements a traditional LIFO structure.
*/
type Stack[T any] struct {
	values []T
}

/*
NewStack creates a new unbounded stack of type T.

Example:

	s := datastructures.NewStack[int]()
	s.Push(3)
	s.PushMany([]int{4, 5, 6})
	value := s.Pop()
	l := s.Size()
	isEmpty := s.IsEmpty()

	// value == 6
	// l == 3
	// isempty == false
*/
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

/*
IsEmpty returns true if there are no values on the stack.
*/
func (s *Stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

/*
Pop pulls the top item off of the stack and returns its value.
*/
func (s *Stack[T]) Pop() T {
	var value T

	if s.IsEmpty() {
		return value
	}

	value = s.Top()
	s.values = s.values[:len(s.values)-1]

	return value
}

/*
Push adds a new item to the top of the stack.
*/
func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

/*
PushMany adds a set of values to the top of the stack.
The last item in the array is the new top of the stack.
*/
func (s *Stack[T]) PushMany(values []T) {
	s.values = append(s.values, values...)
}

/*
Size returns the number of items in the stack.
*/
func (s *Stack[T]) Size() int {
	return len(s.values)
}

/*
Top returns the value of the top item in the stack
without removing it.
*/
func (s *Stack[T]) Top() T {
	return s.values[len(s.values)-1]
}
