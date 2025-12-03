package math

import (
	"math"

	"golang.org/x/exp/constraints"
)

func Ceil[T constraints.Float | constraints.Integer](value T) T {
	fvalue := float64(value)
	c := T(math.Ceil(fvalue))
	return c
}

func Floor[T constraints.Float | constraints.Integer](value T) T {
	fvalue := float64(value)
	c := T(math.Floor(fvalue))
	return c
}
