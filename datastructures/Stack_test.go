package datastructures_test

import (
	"testing"

	"github.com/adampresley/adamgokit/datastructures"
	"github.com/stretchr/testify/assert"
)

func TestNewStack(t *testing.T) {
	stack := datastructures.NewStack[int]()
	assert.Equal(t, true, stack.IsEmpty())

	stack.Push(4)
	stack.PushMany([]int{1, 10})

	assert.Equal(t, 3, stack.Size())
	assert.Equal(t, false, stack.IsEmpty())

	pop1 := stack.Pop()
	pop2 := stack.Pop()
	pop3 := stack.Pop()
	pop4 := stack.Pop()

	assert.Equal(t, 10, pop1)
	assert.Equal(t, 1, pop2)
	assert.Equal(t, 4, pop3)
	assert.Equal(t, 0, pop4)
}
