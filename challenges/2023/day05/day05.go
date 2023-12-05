package aoc2023_day5

import (
	"fmt"
	"math"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type SeedList []int

type RangeMap struct {
	SourceStart      int
	SourceEnd        int
	SourceToDestDiff int
	DestStart        int
	DestEnd          int
	Size             int
}

type SeedRange struct {
	Start int
	End   int
}

type TypeMap struct {
	SourceType string
	DestType   string
	RangeMaps  []RangeMap
}

func (tm TypeMap) String() string {
	return fmt.Sprintf("%s-to-%s map: %v\n", tm.SourceType, tm.DestType, tm.RangeMaps)
}

type Almanac struct {
	Seeds    SeedList
	TypeMaps []TypeMap
}

func part1(input []string) (string, error) {
	almanac := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", almanac)

	locations := make(map[int]int, len(almanac.Seeds))
	minLocation := math.MaxInt
	for _, seed := range almanac.Seeds {
		log.Debugf("Seed: %d\n", seed)

		curVal := seed // find the location for this seed
		for _, tm := range almanac.TypeMaps {
			for _, rm := range tm.RangeMaps {
				if curVal >= rm.SourceStart && curVal <= rm.SourceEnd {
					curVal = curVal + rm.SourceToDestDiff
					break
				}
			}
		}
		locations[seed] = curVal
		if curVal < minLocation {
			minLocation = curVal
		}
	}

	result := minLocation
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	almanac := parseInput(input)

	seedRanges := make([]SeedRange, 0, len(almanac.Seeds)/2)
	for i := 0; i < len(almanac.Seeds); i += 2 {
		seedRanges = append(seedRanges, SeedRange{Start: almanac.Seeds[i], End: almanac.Seeds[i] + almanac.Seeds[i+1]})
	}

	cur := -1
	found := false
	for !found {
		cur++
		curVal := cur
		if cur%1000000 == 0 {
			log.Debugf("Checking Location: %d\n", cur)
		}
		for i := len(almanac.TypeMaps) - 1; i >= 0; i-- {
			tm := almanac.TypeMaps[i]
			for _, rm := range tm.RangeMaps {
				if curVal >= rm.DestStart && curVal <= rm.DestEnd {
					curVal = curVal - rm.SourceToDestDiff
					break
				}
			}
		}
		for _, sr := range seedRanges {
			if curVal >= sr.Start && curVal <= sr.End {
				found = true
				break
			}
		}
	}

	result := cur
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) Almanac {
	i := 0
	almanac := Almanac{
		Seeds:    util.AtoiSplit(strings.Split(input[i], ": ")[1], " "),
		TypeMaps: make([]TypeMap, 0),
	}
	i += 2 // skip blank line

	parseTypeMap := func(idx int, input []string) (TypeMap, int) {
		types := strings.Split(strings.Split(input[idx], " ")[0], "-")
		sourceType, destType := types[0], types[2]
		idx++
		typeMap := TypeMap{
			SourceType: sourceType,
			DestType:   destType,
			RangeMaps:  make([]RangeMap, 0),
		}
		for idx < len(input) && input[idx] != "" {
			var sourceStart, destStart, size int
			fmt.Sscanf(input[idx], "%d %d %d", &destStart, &sourceStart, &size)
			idx++
			typeMap.RangeMaps = append(typeMap.RangeMaps, RangeMap{
				SourceStart:      sourceStart,
				SourceEnd:        sourceStart + size - 1,
				SourceToDestDiff: destStart - sourceStart,
				DestStart:        destStart,
				DestEnd:          destStart + size - 1,
				Size:             size,
			})
		}
		idx++
		return typeMap, idx
	}

	for i < len(input) {
		typeMap, idx := parseTypeMap(i, input)
		almanac.TypeMaps = append(almanac.TypeMaps, typeMap)
		i = idx
	}
	return almanac
}

func init() {
	challenges.RegisterChallengeFunc(2023, 5, 1, "day05.txt", part1)
	challenges.RegisterChallengeFunc(2023, 5, 2, "day05.txt", part2)
}
