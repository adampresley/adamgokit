package datastructures

import "fmt"

/*
LinkedList implements a basic single-linked list, where
each element refers to the next.
*/
type LinkedList[T comparable] struct {
	Head   *SingleListElement[T]
	Length int
}

/*
SingleListElement is a single item in a linked list of type T.
*/
type SingleListElement[T comparable] struct {
	Value T
	Next  *SingleListElement[T]
}

/*
NewLinkedList creates a new linked list of type T.

Example:

	l := NewLinkedList[int]()

	l.Push(2)
	l.Push(4)
	sliceVersion := l.ToSlice() // []int{2, 4}
	second := l.GetAt(1)        // 4

	l.Walk(func(list *datastructures.LinkedList[int], currentEl *datastructures.SingleListElement[int], index int) {
	   // Do stuff
	})
*/
func NewLinkedList[T comparable]() *LinkedList[T] {
	return InitLinkedList([]T{})
}

/*
InitLinkedList creates a new linked list of type T
initialized with a slice of values. The result is a linked
list already populated.
*/
func InitLinkedList[T comparable](values []T) *LinkedList[T] {
	result := &LinkedList[T]{}

	for _, value := range values {
		result.Push(value)
	}

	return result
}

/*
Push adds a new item to the end of the list.
*/
func (sll *LinkedList[T]) Push(value T) {
	var (
		el *SingleListElement[T]
	)

	newEl := &SingleListElement[T]{
		Value: value,
	}

	/*
	 * If this is an empty list, intialize head.
	 */
	if sll.Length == 0 {
		sll.Head = newEl
		sll.Length++
		return
	}

	/*
	 * Otherwise, find the tail and add the new element.
	 */
	el = sll.Head

	for el.Next != nil {
		el = el.Next
	}

	el.Next = newEl
	sll.Length++
}

/*
InsertAt adds a new item at a specified index.
*/
func (sll *LinkedList[T]) InsertAt(value T, index int) error {
	var (
		prev *SingleListElement[T]
		newEl *SingleListElement[T]
	)

	if index < 0 || index >= sll.Length {
		return fmt.Errorf("index out of bounds")
	}

	prev = sll.Head

	for i := 0; i <= index; i++ {
		if index > 0 && i > 0 && i == index - 1 {
			prev = 
		}
	}
}

func (sll *LinkedList[T]) GetAt(index int) *SingleListElement[T] {
	if index < 0 || index > sll.Length-1 {
		return nil
	}

	el := sll.Head

	for i := 0; i < index; i++ {
		el = el.Next
	}

	return el
}

func (sll *LinkedList[T]) Walk(f func(list *LinkedList[T], currentEl *SingleListElement[T], index int)) {
	if sll.Length == 0 {
		return
	}

	el := sll.Head
	index := 0

	for el != nil {
		f(sll, el, index)

		index++
		el = el.Next
	}
}

func (sll *LinkedList[T]) ToSlice() []T {
	result := []T{}

	if sll.Length == 0 {
		return result
	}

	el := sll.Head

	for el != nil {
		result = append(result, el.Value)
		el = el.Next
	}

	return result
}
