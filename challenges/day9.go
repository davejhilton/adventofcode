package challenges

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode2020/log"
)

func day9_part1(input []string) (string, error) {
	nums := day9_parseInts(input)
	log.Debugf("%v\n", nums)

	r := -1
	window := 25

	for i := window; i < len(nums); i++ {
		if !day9_checkSum(nums[i], nums[i-window:i]) {
			r = nums[i]
			break
		}
	}

	return fmt.Sprintf("%d", r), nil
}

func day9_checkSum(val int, nums []int) bool {
	log.Debugf("Checking '%d' against: %v\n", val, nums)
	for i, n := range nums {
		if n >= val {
			continue
		}
		for _, m := range nums[i+1:] {
			if m+n == val {
				return true
			}
		}
	}
	return false
}

func day9_part2(input []string) (string, error) {
	nums := day9_parseInts(input)

	target := 2089807806
	f, l := 0, 0

	sum := nums[0]
	for l < len(nums) {
		if sum < target {
			l++
			sum += nums[l]
		} else if sum > target {
			sum -= nums[f]
			f++
		} else {
			break
		}
	}
	log.Debugf("RANGE: %d - %d\n", f, l)
	log.Debugf("%v\n", nums[f:l+1])
	sm, lg := int((^uint(0))>>1), 0
	for i := f; i <= l; i++ {
		if nums[i] < sm {
			sm = nums[i]
		}
		if nums[i] > lg {
			lg = nums[i]
		}
	}
	r := sm + lg
	log.Debugf("SMALLEST: %d  +  LARGEST: %d  =  %d\n", sm, lg, r)

	return fmt.Sprintf("%d", r), nil
}

func day9_parseInts(input []string) []int {
	nums := make([]int, 0, len(input))
	for _, s := range input {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	return nums
}

func init() {
	registerChallengeFunc(9, 1, "day9.txt", day9_part1)
	registerChallengeFunc(9, 2, "day9.txt", day9_part2)
}
