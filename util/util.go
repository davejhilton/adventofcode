package util

import (
	"math"
	"strconv"
)

func Min(ints ...int) int {
	min := math.MaxInt
	for _, n := range ints {
		if n < min {
			min = n
		}
	}
	return min
}

func Max(nums ...int) int {
	max := math.MinInt
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	return max
}

func Abs(n int) int {
	if n < 0 {
		return n * -1
	}
	return n
}

func Atoi(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}
