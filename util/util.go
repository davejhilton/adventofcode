package util

import (
	"strconv"
)

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func Min[T Numeric](nums ...T) T {
	var min T
	for _, n := range nums {
		if n < min {
			min = n
		}
	}
	return min
}

func Max[T Numeric](nums ...T) T {
	var max T
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	return max
}

func Abs[T Numeric](n T) T {
	if n < 0 {
		return n * -1
	}
	return n
}

func Atoi(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}
