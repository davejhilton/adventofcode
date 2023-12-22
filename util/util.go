package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Numeric interface {
	~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type SignedNumeric interface {
	~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type IntLike interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
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

func Abs[T SignedNumeric](n T) T {
	if n < 0 {
		return n * -1
	}
	return n
}

func Atoi(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}

func AtoiSplit(str, sep string) (nums []int) {
	parts := strings.Split(str, sep)
	nums = make([]int, 0, len(parts))
	for _, p := range parts {
		nums = append(nums, Atoi(p))
	}
	return nums
}

var numRegex = regexp.MustCompile(`\d+`)

func ExtractNumbers(s string) []int {
	nums := make([]int, 0)
	for _, n := range numRegex.FindAllString(s, -1) {
		nums = append(nums, Atoi(n))
	}
	return nums
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

func Filter[T any](s []T, f func(T) bool) []T {
	new := make([]T, 0)
	for _, v := range s {
		if f(v) {
			new = append(new, v)
		}
	}
	return new
}

func Contains[T comparable](s []T, v T) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func Sum[T Numeric](s []T) T {
	var sum T
	for _, x := range s {
		sum += x
	}
	return sum
}

func JoinInts(s []int, sep string) string {
	var sb strings.Builder
	for i, n := range s {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(fmt.Sprintf("%d", n))
	}
	return sb.String()
}

func GCD[T IntLike](a, b T) T {
	// Calculate the greatest common divisor (gcd) using the Euclidean algorithm
	// https://en.wikipedia.org/wiki/Euclidean_algorithm
	var temp T
	for b != 0 {
		temp = b
		b = a % b
		a = temp
	}
	return a
}

func LCM[T IntLike](a, b T) T {
	return a * b / GCD(a, b)
}

type Grid [][]string

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		sb.WriteString(fmt.Sprintf("%s\n", strings.Join(row, "")))
	}
	return sb.String()
}

type IntGrid [][]int

func (ig IntGrid) String() string {
	var sb strings.Builder
	for _, row := range ig {
		for _, n := range row {
			sb.WriteString(fmt.Sprintf("%d", n))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type Coord struct {
	Row       int
	Col       int
	IntVal    int
	StringVal string
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.Row, c.Col)
}

func (c Coord) ManhattanDistance(other Coord) int {
	return Abs(c.Row-other.Row) + Abs(c.Col-other.Col)
}

func ShoelaceArea(coords []Coord) int {
	// https://en.wikipedia.org/wiki/Shoelace_formula
	// https://www.mathopenref.com/coordpolygonarea.htmll

	var sum int
	for i := 0; i < len(coords)-1; i++ {
		sum += coords[i].Col*coords[i+1].Row - coords[i+1].Col*coords[i].Row
	}

	return sum / 2
}

func PicksTheorem(innerArea int, perimeterArea int) int {
	// https://en.wikipedia.org/wiki/Pick%27s_theorem
	// https://www.mathopenref.com/polygonirregulararea.html

	// A = i + b/2 - 1
	// i = A - b/2 + 1
	return innerArea - (perimeterArea / 2) + 1
}
