package challenges2021

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func day08_part1(input []string) (string, error) {
	inputOutputPairs := day08_parse(input)
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

func day08_part2(input []string) (string, error) {
	inputOutputPairs := day08_parse(input)
	result := 0

	for _, io := range inputOutputPairs {

		one := io.InputMap[2][0]      // TR, BR
		four := io.InputMap[4][0]     // TL, TR, M, BR
		seven := io.InputMap[3][0]    // T, TR, BR - the segment that's unique from 'one' is T
		eight := io.InputMap[7][0]    // all segments
		sixNineZero := io.InputMap[6] // all segments

		T := withoutRunes(seven, []rune(one))
		tl_m := withoutRunes(four, []rune(one))
		tr_bl_m := dedup(strings.Join([]string{
			xor(sixNineZero[0], sixNineZero[1]),
			xor(sixNineZero[1], sixNineZero[2]),
		}, ""))
		TR := and(tr_bl_m, one)
		BR := xor(seven, strings.Join([]string{T, TR}, ""))
		M := and(tl_m, tr_bl_m)
		TL := xor(tl_m, M)
		BL := xor(tr_bl_m, strings.Join([]string{M, TR}, ""))
		B := xor(eight, strings.Join([]string{T, TL, TR, M, BL, BR}, ""))

		segMap := *NewSegMap()
		segMap["T"] = T
		segMap["TL"] = TL
		segMap["TR"] = TR
		segMap["M"] = M
		segMap["BL"] = BL
		segMap["BR"] = BR
		segMap["B"] = B

		two := sortStr(strings.Join([]string{T, TR, M, BL, B}, ""))
		three := sortStr(strings.Join([]string{T, TR, M, BR, B}, ""))
		five := sortStr(strings.Join([]string{T, TL, M, BR, B}, ""))
		six := sortStr(strings.Join([]string{T, TL, M, BL, BR, B}, ""))
		nine := sortStr(strings.Join([]string{T, TL, TR, M, BR, B}, ""))
		zero := sortStr(strings.Join([]string{T, TL, TR, BL, BR, B}, ""))

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
		nOutputs := 4
		num := 0
		displays := make([]digitDisplay, 0)
		for i, o := range io.Outputs {
			n := numMap[o]
			displays = append(displays, NewDigitDisplay(o, segMap))
			num += n * int(math.Pow10(nOutputs-i-1))
		}
		log.Debugf("%s\n%s\n\n", strings.Join(io.Outputs, " "), multiDigitDisplayString(displays...))
		result += num
	}

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func withoutRunes(str string, runes []rune) string {
	for _, r := range runes {
		if strings.ContainsRune(str, r) {
			str = strings.Replace(str, string(r), "", 1)
		}
	}
	return str
}

func and(str1, str2 string) string {
	var sb strings.Builder
	for _, r := range str1 {
		if strings.ContainsRune(str2, r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func xor(strs ...string) string {
	cp := strs[0:]
	for i := 0; i < len(strs)-1; i += 1 {
		str := strs[i]
		for _, r := range str {
			for j := 0; j < len(strs); j += 1 {
				if i != j {
					if strings.ContainsRune(strs[j], r) {
						for k := range cp {
							cp[k] = strings.Replace(cp[k], string(r), "", 1)
						}
					}
				}
			}
		}
	}
	return strings.Join(cp, "")
}

func dedup(str string) string {
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

func sortStr(str string) string {
	split := strings.Split(str, "")
	sort.Strings(split)
	return strings.Join(split, "")
}

func day08_parse(input []string) []*ioPair {
	pairs := make([]*ioPair, 0, len(input))
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
			sorted := sortStr(v)
			inMap[l] = append(inMap[l], sorted)
			in[i] = sorted
		}
		for i, v := range out {
			l := len(v)
			if _, ok := outMap[l]; !ok {
				outMap[l] = make([]string, 0)
			}
			sorted := sortStr(v)
			outMap[l] = append(outMap[l], sorted)
			out[i] = sorted
		}
		pairs = append(pairs, &ioPair{
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

type segMap map[string]string

type digitDisplay struct {
	Segments string
	Mappings segMap
	display  segMap
}

func NewDigitDisplay(segments string, mappings segMap) digitDisplay {
	display := *NewSegMap()
	for k, s := range mappings {
		if strings.Contains(segments, s) {
			display[k] = s
		}
	}
	return digitDisplay{
		Segments: segments,
		Mappings: mappings,
		display:  display,
	}
}

func (dd digitDisplay) String() string {
	return dd.display.String()
}

func (d segMap) String() string {
	var sb strings.Builder
	T, TL, TR, M, BL, BR, B := d["T"], d["TL"], d["TR"], d["M"], d["BL"], d["BR"], d["B"]
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

func multiDigitDisplayString(digits ...digitDisplay) string {
	var sb strings.Builder
	sep := "  "
	digitLines := make([][]string, 0, len(digits))
	for _, d := range digits {
		digitLines = append(digitLines, strings.Split(d.String(), "\n"))
	}
	for rowNum := range digitLines[0] {
		if rowNum != 0 {
			sb.WriteString("\n")
		}
		rowParts := make([]string, 0)
		for _, digit := range digitLines {
			rowParts = append(rowParts, digit[rowNum])
		}
		sb.WriteString(strings.Join(rowParts, sep))
	}
	return sb.String()
}

func NewSegMap() *segMap {
	return &segMap{
		"T":  ".",
		"TL": ".",
		"TR": ".",
		"M":  ".",
		"BL": ".",
		"BR": ".",
		"B":  ".",
	}
}

func init() {
	challenges.RegisterChallengeFunc(2021, 8, 1, "day08.txt", day08_part1)
	challenges.RegisterChallengeFunc(2021, 8, 2, "day08.txt", day08_part2)
}
