package challenges

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/davejhilton/adventofcode2020/log"
)

func day10_part1(input []string, isExample bool) (string, error) {
	nums := day10_getSortedJoltages(input)

	ones := 0
	threes := 0

	for i := 0; i < len(nums); i++ {
		if i != 0 {
			diff := nums[i] - nums[i-1]
			if diff == 1 {
				ones++
			} else if diff == 3 {
				threes++
			}
			log.Debugf("step size from %d to %d: %d\n", nums[i-1], nums[i], diff)
		}
	}

	log.Debugf("number of +1 steps: %d\n", ones)
	log.Debugf("number of +3 steps: %d\n", threes)

	return fmt.Sprintf("%d", ones*threes), nil
}

func day10_part2(input []string, isExample bool) (string, error) {
	nums := day10_getSortedJoltages(input)

	log.Debugf("%v\n", nums)

	optionsFromVal := make(map[int]int)

	for i := len(nums) - 1; i >= 0; i-- {
		val := nums[i]
		if i == len(nums)-1 {
			optionsFromVal[val] = 1 // last value has only 1 possible combination/option, by definition
			log.Debugf("FIRST: %d - %d\n", i, val)
		} else {
			sum := 0
			log.Debugf("i = %d, val = %d\n", i, val)
			for diff := 1; diff <= 3; diff++ { // for each valid increment...
				log.Debugf("  -- checking %d... ", val+diff)
				if nextValOptions, ok := optionsFromVal[val+diff]; ok { // if an adapter exists for {current joltage} + {this increment}
					sum += nextValOptions // taking that path next gives us {that many} ways to the end goal from here
					log.Debugf(" + %d options to the target\n", nextValOptions)
				} else {
					log.Debugln("(nope)")
				}
			}
			optionsFromVal[val] = sum // memo for how many options there are to get to the end from this joltage
			log.Debugf("MEMO: value %d -> %d options to the target\n\n", val, sum)
		}
	}
	result := optionsFromVal[0] // total # of options from the first joltage (0) to the end

	return fmt.Sprintf("%d", result), nil
}

func day10_getSortedJoltages(input []string) []int {
	nums := make([]int, 0, len(input)+2)
	nums = append(nums, 0) // add the "starting" joltage of 0
	for _, s := range input {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	sort.Ints(nums)                          // sort em
	nums = append(nums, nums[len(nums)-1]+3) // add the "device joltage" at the end
	return nums
}

func init() {
	registerChallengeFunc(10, 1, "day10.txt", day10_part1)
	registerChallengeFunc(10, 2, "day10.txt", day10_part2)
}
