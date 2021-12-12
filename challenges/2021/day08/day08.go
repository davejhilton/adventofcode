package aoc2021_day8

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	inputOutputPairs := parseInput(input)
	result := 0
	for _, io := range inputOutputPairs {
		uniqueLengths := []int{2, 3, 4, 7}
		for _, l := range uniqueLengths {
			if v, ok := io.OutputMap[l]; ok {
				result += len(v)
			}
		}
	}
	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	inputOutputPairs := parseInput(input)
	result := 0

	for _, io := range inputOutputPairs {

		one := sortString(io.InputMap[2][0])   // TR, BR
		four := sortString(io.InputMap[4][0])  // TL, TR, M, BR
		seven := sortString(io.InputMap[3][0]) // T, TR, BR
		eight := sortString(io.InputMap[7][0]) // all segments
		six_nine_zero := io.InputMap[6]

		T := symmetricDiff(seven, one)
		TL_M := symmetricDiff(four, one)
		TR_M_BL := unique(strings.Join([]string{
			symmetricDiff(six_nine_zero[0], six_nine_zero[1]),
			symmetricDiff(six_nine_zero[1], six_nine_zero[2]),
		}, ""))
		TR := intersection(TR_M_BL, one)
		BR := symmetricDiff(seven, strings.Join([]string{T, TR}, ""))
		M := intersection(TL_M, TR_M_BL)
		TL := symmetricDiff(TL_M, M)
		BL := symmetricDiff(TR_M_BL, strings.Join([]string{M, TR}, ""))
		B := symmetricDiff(eight, strings.Join([]string{T, TL, TR, M, BL, BR}, ""))

		two := sortString(strings.Join([]string{T, TR, M, BL, B}, ""))
		three := sortString(strings.Join([]string{T, TR, M, BR, B}, ""))
		five := sortString(strings.Join([]string{T, TL, M, BR, B}, ""))
		six := sortString(strings.Join([]string{T, TL, M, BL, BR, B}, ""))
		nine := sortString(strings.Join([]string{T, TL, TR, M, BR, B}, ""))
		zero := sortString(strings.Join([]string{T, TL, TR, BL, BR, B}, ""))

		numMap := map[string]int{
			zero:  0,
			one:   1,
			two:   2,
			three: 3,
			four:  4,
			five:  5,
			six:   6,
			seven: 7,
			eight: 8,
			nine:  9,
		}
		allSegments := segmentMap{"T": T, "TL": TL, "TR": TR, "M": M, "BL": BL, "BR": BR, "B": B}
		displays := make([]segmentMap, 0)
		numericOutput := 0
		for i, o := range io.Outputs {
			n := numMap[sortString(o)]
			displays = append(displays, NewDigitDisplay(o, allSegments))
			numericOutput += n * int(math.Pow10(len(io.Outputs)-i-1))
		}
		log.Debugf("%s\n%s\n\n", strings.Join(io.Outputs, " "), multiDigitDisplayString(displays...))
		result += numericOutput
	}

	return fmt.Sprintf("%d", result), nil
}

func intersection(str1, str2 string) string {
	var sb strings.Builder
	for _, r := range str1 {
		if strings.ContainsRune(str2, r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func symmetricDiff(str1 string, str2 string) string {
	s1, s2 := str1, str2
	for _, r := range str2 {
		if strings.ContainsRune(str1, r) {
			s1 = strings.Replace(s1, string(r), "", 1)
			s2 = strings.Replace(s2, string(r), "", 1)
		}
	}
	return fmt.Sprintf("%s%s", s1, s2)
}

func unique(str string) string {
	u := make(map[rune]bool)
	for _, r := range str {
		u[r] = true
	}
	r := make([]rune, 0)
	for k := range u {
		r = append(r, k)
	}
	return string(r)
}

func sortString(str string) string {
	split := strings.Split(str, "")
	sort.Strings(split)
	return strings.Join(split, "")
}

func parseInput(input []string) []ioPair {
	pairs := make([]ioPair, 0, len(input))
	for _, s := range input {
		halves := strings.SplitN(s, " | ", 2)
		in := strings.Split(halves[0], " ")
		out := strings.Split(halves[1], " ")
		inMap := make(map[int][]string)
		outMap := make(map[int][]string)
		for i, v := range in {
			l := len(v)
			if _, ok := inMap[l]; !ok {
				inMap[l] = make([]string, 0)
			}
			inMap[l] = append(inMap[l], v)
			in[i] = v
		}
		for i, v := range out {
			l := len(v)
			if _, ok := outMap[l]; !ok {
				outMap[l] = make([]string, 0)
			}
			outMap[l] = append(outMap[l], v)
			out[i] = v
		}
		pairs = append(pairs, ioPair{
			Inputs:    in,
			InputMap:  inMap,
			Outputs:   out,
			OutputMap: outMap,
		})
	}
	return pairs
}

type ioPair struct {
	Inputs    []string
	InputMap  map[int][]string
	Outputs   []string
	OutputMap map[int][]string
}

func (io ioPair) String() string {
	return fmt.Sprintf("%s | %s", strings.Join(io.Inputs, " "), strings.Join(io.Outputs, " "))
}

type segmentMap map[string]string

func EmptySegmentMap() segmentMap {
	return segmentMap{
		"T":  ".",
		"TL": ".",
		"TR": ".",
		"M":  ".",
		"BL": ".",
		"BR": ".",
		"B":  ".",
	}
}

func NewDigitDisplay(segments string, mappings segmentMap) segmentMap {
	display := EmptySegmentMap()
	for k, s := range mappings {
		if strings.Contains(segments, s) {
			display[k] = s
		}
	}
	return display
}

func (sm segmentMap) String() string {
	var sb strings.Builder
	T, TL, TR, M, BL, BR, B := sm["T"], sm["TL"], sm["TR"], sm["M"], sm["BL"], sm["BR"], sm["B"]
	segList := []*string{&T, &TL, &TR, &M, &BL, &BR, &B}
	for _, v := range segList {
		if *v != "." {
			*v = log.Colorize(*v, log.Red, 0)
		}
	}
	sb.WriteString(fmt.Sprintf(" %s%s%s%s \n", T, T, T, T))
	l2_3 := fmt.Sprintf("%s    %s\n", TL, TR)
	sb.WriteString(l2_3)
	sb.WriteString(l2_3)
	sb.WriteString(fmt.Sprintf(" %s%s%s%s \n", M, M, M, M))
	l5_6 := fmt.Sprintf("%s    %s\n", BL, BR)
	sb.WriteString(l5_6)
	sb.WriteString(l5_6)
	sb.WriteString(fmt.Sprintf(" %s%s%s%s ", B, B, B, B))
	return sb.String()
}

func multiDigitDisplayString(digits ...segmentMap) string {
	var sb strings.Builder
	sep := "  "
	displayRows := make([][]string, 0, len(digits))
	for _, d := range digits {
		displayRows = append(displayRows, strings.Split(d.String(), "\n"))
	}
	for rowNum := range displayRows[0] {
		if rowNum != 0 {
			sb.WriteString("\n")
		}
		rowParts := make([]string, 0)
		for _, digit := range displayRows {
			rowParts = append(rowParts, digit[rowNum])
		}
		sb.WriteString(strings.Join(rowParts, sep))
	}
	return sb.String()
}

func init() {
	challenges.RegisterChallengeFunc(2021, 8, 1, "day08.txt", part1)
	challenges.RegisterChallengeFunc(2021, 8, 2, "day08.txt", part2)
}
