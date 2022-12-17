package util

import (
	"encoding/json"
	"sort"
	"strconv"
)

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func Min[T Numeric](nums ...T) T {
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	if len(nums) > 0 {
		return nums[0]
	}
	return 0
}

func Max[T Numeric](nums ...T) T {
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	if len(nums) > 0 {
		return nums[len(nums)-1]
	}
	return 0
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

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func ToJSON(v any, pretty bool) string {
	var j []byte
	if pretty {
		j, _ = json.MarshalIndent(v, "", "  ")
	} else {
		j, _ = json.Marshal(v)
	}
	return string(j)
}
