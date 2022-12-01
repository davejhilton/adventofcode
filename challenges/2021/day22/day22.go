package aoc2021_day22

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	initSteps := parseInput(input)
	const MIN, MAX = -50, 50

	grid := make(Grid)

	cubesOn := 0
	for _, s := range initSteps {
		value := s.TurnOn
		for x := s.X.Start; x <= s.X.Stop; x++ {
			_, ok := grid[x]
			if !ok {
				grid[x] = make(map[int]map[int]bool)
			}
			for y := s.Y.Start; y <= s.Y.Stop; y++ {
				_, ok := grid[x][y]
				if !ok {
					grid[x][y] = make(map[int]bool)
				}
				for z := s.Z.Start; z <= s.Z.Stop; z++ {
					if x < MIN || y < MIN || z < MIN || x > MAX || y > MAX || z > MAX {
						return fmt.Sprintf("%d", cubesOn), nil
					}
					wasOn := grid[x][y][z]
					grid[x][y][z] = value
					if !wasOn && value {
						cubesOn++
					} else if wasOn && !value {
						cubesOn--
					}
				}
			}
		}
	}

	return fmt.Sprintf("%d", cubesOn), nil
}

/*

"LEFT MISS"             A, B ==> A, B
    | |
 | |

"RIGHT MISS"            A, B ==> B, A
 | |
    | |

"LEFT STAGGER"          A, B ==> A, AB, B
    |  |
  |  |

"RIGHT STAGGER"         A, B ==> B, AB, A
  |  |
    |  |

"LEFT-MATCH INNER"      A, B ==> AB, B
  |  |
  | |

"RIGHT-MATCH INNER"     A, B ==> B, AB
  |  |
   | |

"FULL INNER"            A, B ==> B, AB, B
  |   |
   | |

"EXACT MATCH"           A, B ==> AB
  |  |
  |  |

"LEFT-MATCH OUTER"      A, B ==> AB, A
    |  |
    |    |

"RIGHT_MATCH OUTER"     A, B ==> A, AB
    |  |
  |    |

"FULL OUTER"            A, B ==> A, AB, A
    |  |
   |    |

*/

func addRange(rN *Range, ranges *[]*Range) {

	i := 0
	end := len(*ranges)
	for ; i < end; i++ {
		r2 := (*ranges)[i]
		if rN.Dimension == "x" {
			log.Debugf("rN: %d..%d\n", rN.Start, rN.Stop)
		}

		if rN.Start > r2.Stop { // RIGHT MISS
			// r2 | |
			// rN    | |

			// if this is the last r2 in the list, add rN to the end of the list
			if i == len(*ranges)-1 {
				if rN.Dimension == "x" {
					log.Debugln("APPENDING TO END OF LIST (RIGHT MISS)")
				}
				*ranges = append(*ranges, rN)
				if rN.Dimension == "x" {
					PrintList(*ranges)
				}
				return
			}
			continue // otherwise, it might overlap with a different range
		}
		if rN.Stop < r2.Start { // LEFT MISS
			// r2    | |
			// rN | |

			// add rN before r2 in the list. Then return (nothing else will overlap)
			if rN.Dimension == "x" {
				log.Debugf("INSERTING INTO LIST BEFORE %d (LEFT MISS)\n", i)
			}
			insertAtIndex(ranges, i, rN)
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}
			return
		}
		if rN.Start < r2.Start && rN.Stop < r2.Stop { // LEFT STAGGER
			// r2   |   |
			// rN  |  |

			// 1. split rN into rN, rN_joint
			rN_joint := cloneRange(rN)
			rN_joint.Start = r2.Start
			rN.Stop = r2.Start - 1

			// 2. split r2 into r2_joint, r2
			r2_joint := cloneRange(r2)
			r2_joint.Stop = rN_joint.Stop
			r2.Start = r2_joint.Stop + 1

			// 3. merge rN_joint, r2_joint into 'joint'
			// log.Debugln("MERGING CHILD RANGES: (LEFT STAGGER)")
			joint := mergeChildRanges(rN_joint, r2_joint)

			// 4. list should have: [... , rN, joint, r2, ...]
			if rN.Dimension == "x" {
				log.Debugf("INSERTING INTO LIST AT %d (and %d) (LEFT STAGGER)\n", i, i+1)
			}
			insertAtIndex(ranges, i, rN)
			insertAtIndex(ranges, i+1, joint)
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// 5. Then return - nothing else later in the list can overlap
			return
		}
		if rN.Start > r2.Start && rN.Stop > r2.Stop { // RIGHT STAGGER
			// r2 |  |
			// rN  |   |

			// 1. split r2 into r2, r2_joint
			r2_joint := cloneRange(r2)
			r2_joint.Start = rN.Start
			r2.Stop = r2_joint.Start - 1

			// 2. split rN into rN_joint, rN_sep
			rN_joint := cloneRange(rN)
			rN_joint.Stop = r2_joint.Stop
			rN.Start = rN_joint.Stop + 1

			// 3. merge r2_joint, rN_joint into 'joint'
			// log.Debugln("MERGING CHILD RANGES: (RIGHT STAGGER)")
			joint := mergeChildRanges(rN_joint, r2_joint)

			// 4. list should have: [... , r2, joint, ...]
			if rN.Dimension == "x" {
				log.Debugf("INSERTING INTO LIST AT %d (RIGHT STAGGER)\n", i+1)
			}
			insertAtIndex(ranges, i+1, joint)
			if i == end-1 {
				insertAtIndex(ranges, i+2, rN)
				i += 1
				end += 1
			}
			i += 1
			end += 1
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// continue iterating. Other overlaps are still possible
			continue
		}
		if rN.Start == r2.Start && rN.Stop < r2.Stop { // LEFT-MATCH INNER
			// r2 |  |
			// rN | |

			// 1. split r2 into r2_joint, r2
			r2_joint := cloneRange(r2)
			r2_joint.Stop = rN.Stop
			r2.Start = r2_joint.Stop + 1

			// 2. merge rN, r2_joint into 'joint'
			// log.Debugln("MERGING CHILD RANGES: (LEFT-MATCH INNER)")
			joint := mergeChildRanges(rN, r2_joint)

			// 3. list should have: [... , joint, r2, ...]
			if rN.Dimension == "x" {
				log.Debugf("INSERTING INTO LIST AT %d (LEFT-MATCH INNER)\n", i)
			}
			insertAtIndex(ranges, i, joint)
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// 4. return - nothing else later in the list can overlap
			return
		}
		if rN.Start > r2.Start && rN.Stop == r2.Stop { // RIGHT-MATCH INNER
			// r2 |  |
			// rN  | |

			// 1. split r2 into r2, r2_joint
			r2_joint := cloneRange(r2)
			r2_joint.Start = rN.Start
			r2.Stop = r2_joint.Start - 1

			// 2. merge r2_joint, rN into 'joint'
			// log.Debugln("MERGING CHILD RANGES: (RIGHT-MATCH INNER)")
			joint := mergeChildRanges(rN, r2_joint)

			// 3. list should have: [... , r2, joint, ...]
			if rN.Dimension == "x" {
				log.Debugf("INSERTING INTO LIST AT %d (RIGHT-MATCH INNER)\n", i+1)
			}
			insertAtIndex(ranges, i+1, joint)
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// 4. return - nothing else later in the list can overlap
			return
		}
		if rN.Start > r2.Start && rN.Stop < r2.Stop { // FULL INNER
			// r2 |   |
			// rN  | |

			// 1. split r2 into r2, r2_joint, r2_after
			r2_joint := cloneRange(r2)
			r2_joint.Start = rN.Start
			r2_joint.Stop = rN.Stop
			r2_after := cloneRange(r2)
			r2_after.Start = rN.Stop + 1
			r2.Stop = rN.Start - 1

			// 2. merge r2_joint, rN into 'joint'
			// log.Debugln("MERGING CHILD RANGES: (FULL INNER)")
			joint := mergeChildRanges(rN, r2_joint)

			// 3. list should have: [... , r2, joint, r2_after ...]
			if rN.Dimension == "x" {
				log.Debugf("INSERTING INTO LIST AT %d (AND %D) (FULL INNER)\n", i+1, i+2)
			}
			insertAtIndex(ranges, i+1, joint)
			insertAtIndex(ranges, i+2, r2_after)
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// 4. return - nothing else later in the list can overlap
			return
		}

		if rN.Start == r2.Start && rN.Stop == r2.Stop { // EXACT MATCH
			// r2 |  |
			// rN |  |

			// 1. merge r2, rN into 'joint'
			// log.Debugln("MERGING CHILD RANGES: (EXACT MATCH)")
			joint := mergeChildRanges(rN, r2)

			// 2. list should have: [... , joint, ...]
			if rN.Dimension == "x" {
				log.Debugf("REPLACING LIST AT %d (EXACT MATCH)\n", i)
			}
			(*ranges)[i] = joint
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// 3. return - nothing else later in the list can overlap
			return
		}

		if rN.Start == r2.Start && rN.Stop > r2.Stop { // LEFT-MATCH OUTER
			// r2 |  |
			// rN |    |

			// 1. split rN into rN_joint, rN_sep
			rN_joint := cloneRange(rN)
			rN_joint.Stop = r2.Stop
			rN.Start = r2.Stop + 1

			// 2. merge rN_joint, r2
			// log.Debugln("MERGING CHILD RANGES: (LEFT-MATCH OUTER)")
			joint := mergeChildRanges(rN_joint, r2)

			// 3. list should have: [... , joint, ...]
			if rN.Dimension == "x" {
				log.Debugf("REPLACING LIST AT %d (LEFT-MATCH OUTER)\n", i)
			}
			(*ranges)[i] = joint
			if i == end-1 {
				insertAtIndex(ranges, i+1, rN)
				i += 1
				end += 1
			}
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// 4. continue - other ranges might overlap, too
			continue
		}

		if rN.Start < r2.Start && rN.Stop == r2.Stop { // RIGHT-MATCH OUTER
			// r2     |  |
			// rN   |    |

			// 1. split rN into rN_sep, rN_joint
			rN_joint := cloneRange(rN)
			rN_joint.Start = r2.Start
			rN.Stop = r2.Start - 1

			// 2. merge rN_joint, r2 into 'joint'
			// log.Debugln("MERGING CHILD RANGES: (RIGHT-MATCH OUTER)")
			joint := mergeChildRanges(rN_joint, r2)

			// 3. list should have: [... , rN_sep, joint, ...]
			if rN.Dimension == "x" {
				log.Debugf("REPLACING LIST AT %d (AND INSERTING AT %d) (RIGHT-MATCH OUTER)\n", i, i+1)
			}
			insertAtIndex(ranges, i, rN)
			(*ranges)[i+1] = joint
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// 4. return - nothing else can overlap
			return
		}

		if rN.Start < r2.Start && rN.Stop > r2.Stop { // FULL OUTER
			// r2    |  |
			// rN   |    |

			// 1. split rN into rN, rN_joint, rN_after
			rN_joint := cloneRange(rN)
			rN_joint.Start = r2.Start
			rN_joint.Stop = r2.Stop
			rN_before := cloneRange(rN)
			rN_before.Stop = r2.Start - 1
			rN.Start = r2.Stop + 1

			// 2. merge rN_joint, r2 into 'joint'
			// log.Debugln("MERGING CHILD RANGES: (FULL OUTER)")
			joint := mergeChildRanges(rN_joint, r2)

			// 3. list should have: [... , rN_before, joint, ...]
			if rN.Dimension == "x" {
				log.Debugf("INSERTING, REPLACING AT %d, %d (FULL OUTER)\n", i, i+1)
			}
			insertAtIndex(ranges, i, rN_before)
			(*ranges)[i+1] = joint
			if i == end-1 {
				insertAtIndex(ranges, i+2, rN)
				i += 1
				end += 1
			}
			i += 1
			end += 1
			if rN.Dimension == "x" {
				PrintList(*ranges)
			}

			// 4. continue - other things might overlap, too
			continue
		}
	}
}

func subRange(rN *Range, ranges *[]*Range) {

	i := 0
	end := len(*ranges)
	for ; i < end; i++ {
		r2 := (*ranges)[i]

		if rN.Start > r2.Stop { // RIGHT MISS //
			// r2 | |
			// rN    | |

			// nothing to do.
			continue // it might overlap with a different range though, so keep going
		}
		if rN.Stop < r2.Start { // LEFT MISS //
			// r2    | |
			// rN | |

			// just return (nothing will overlap)
			return
		}
		if rN.Start < r2.Start && rN.Stop < r2.Stop { // LEFT STAGGER
			// r2   |   |
			// rN  |  |

			// 1. split r2 into r2_overlap, r2
			r2_overlap := cloneRange(r2)
			r2_overlap.Stop = rN.Stop
			r2.Start = rN.Stop + 1

			// 2. negateRanges rN on r2_overlap
			negateChildRanges(rN, r2_overlap)

			// 3a. if anything left of r2_overlap, list should have: [... , r2_overlap, r2, ...]
			if len(r2_overlap.Ranges) > 0 {
				insertAtIndex(ranges, i, r2_overlap)
			}
			// 3b. otherwise, list should have: [... , r2, ...]
			// (do nothing)

			// 4. Then return - nothing else later in the list can overlap
			return
		}
		if rN.Start > r2.Start && rN.Stop > r2.Stop { // RIGHT STAGGER
			// r2 |  |
			// rN  |   |

			// 1. split r2 into r2, r2_overlap
			r2_overlap := cloneRange(r2)
			r2_overlap.Start = rN.Start
			r2.Stop = rN.Start - 1

			// 2. negateRanges rN on r2_overlap
			negateChildRanges(rN, r2_overlap)

			// 3a. if anything left of r2_overlap, list should have: [... , r2, r2_overlap, ...]
			if len(r2_overlap.Ranges) > 0 {
				insertAtIndex(ranges, i+1, r2_overlap)
				i += 1
				end += 1
			}
			// 3b. otherwise, list should have: [... , r2, ...]
			// (do nothing)

			// 4. continue, since other overlaps may happen
			continue
		}
		if rN.Start == r2.Start && rN.Stop < r2.Stop { // LEFT-MATCH INNER
			// r2 |  |
			// rN | |

			// 1. split r2 into r2_overlap, r2
			r2_overlap := cloneRange(r2)
			r2_overlap.Stop = rN.Stop
			r2.Start = rN.Stop + 1

			// 2. negateRanges rN on r2_overlap
			negateChildRanges(rN, r2_overlap)

			// 3a. if anything left of r2_overlap, list should have: [... , r2_overlap, r2, ...]
			if len(r2_overlap.Ranges) > 0 {
				insertAtIndex(ranges, i, r2_overlap)
			}
			// 3b. otherwise, list should have: [... , r2, ...]
			// (do nothing)

			// 4. return, since nothing else can match
			return
		}
		if rN.Start > r2.Start && rN.Stop == r2.Stop { // RIGHT-MATCH INNER
			// r2 |  |
			// rN  | |

			// 1. split r2 into r2, r2_overlap
			r2_overlap := cloneRange(r2)
			r2_overlap.Start = rN.Start
			r2.Stop = rN.Start - 1

			// 2. negateRanges rN on r2_overlap
			negateChildRanges(rN, r2_overlap)

			// 3a. if anything left of r2_overlap, list should have: [... , r2, r2_overlap, ...]
			if len(r2_overlap.Ranges) > 0 {
				insertAtIndex(ranges, i+1, r2_overlap)
			}
			// 3b. otherwise, list should have: [... , r2, ...]
			// (do nothing)

			// 4. return, since nothing else will overlap
			return
		}
		if rN.Start > r2.Start && rN.Stop < r2.Stop { // FULL INNER
			// r2 |   |
			// rN  | |

			// 1. split r2 into r2, r2_overlap, r2_after
			r2_overlap := cloneRange(r2)
			r2_overlap.Start = rN.Start
			r2_overlap.Stop = rN.Stop
			r2_after := cloneRange(r2)
			r2_after.Start = rN.Stop + 1
			r2.Stop = rN.Start - 1

			// 2. negateRanges rN on r2_overlap
			negateChildRanges(rN, r2_overlap)

			// 3a. if anything left of r2_overlap, list should have: [... , r2, r2_overlap, r2_after, ...]
			if len(r2_overlap.Ranges) > 0 {
				insertAtIndex(ranges, i+1, r2_overlap)
				insertAtIndex(ranges, i+2, r2_after)

				// 3b. otherwise, list should have: [... , r2, r2_after ...]
			} else {
				insertAtIndex(ranges, i+1, r2_after)
			}
			// 4. return, since nothing else can overlap
			return
		}

		if rN.Start == r2.Start && rN.Stop == r2.Stop { // EXACT MATCH
			// r2 |  |
			// rN |  |

			// 1. negateRanges rN on r2
			negateChildRanges(rN, r2)

			// 2a. if anything left of r2, list should have: [... , r2, ...]
			// (nothing to do)

			// 2b. otherwise, list should delete r2!
			if len(r2.Ranges) == 0 {
				deleteAtIndex(ranges, i)
			}

			// 3. return, since nothing else can overlap
			return
		}

		if rN.Start == r2.Start && rN.Stop > r2.Stop { // LEFT-MATCH OUTER
			// r2 |  |
			// rN |    |

			// 1. negateRanges rN on r2
			negateChildRanges(rN, r2)

			// 2a. if anything left of r2, list should have: [... , r2, ...]
			// (nothing to do)

			// 2b. otherwise, list should delete r2!
			if len(r2.Ranges) == 0 {
				deleteAtIndex(ranges, i)
				i -= 1
				end -= 1
			}

			// 3. continue, since other overlaps may happen
			continue
		}

		if rN.Start < r2.Start && rN.Stop == r2.Stop { // RIGHT-MATCH OUTER
			// r2     |  |
			// rN   |    |

			// 1. negateRanges rN on r2
			negateChildRanges(rN, r2)

			// 2a. if anything left of r2, list should have: [... , r2, ...]
			// (nothing to do)

			// 2b. otherwise, list should delete r2!
			if len(r2.Ranges) == 0 {
				deleteAtIndex(ranges, i)
			}

			// 3. return, since nothing else can overlap
			return
		}

		if rN.Start < r2.Start && rN.Stop > r2.Stop { // FULL OUTER
			// r2    |  |
			// rN   |    |

			// 1. negateRanges rN on r2
			negateChildRanges(rN, r2)

			// 2a. if anything left of r2, list should have: [... , r2, ...]
			// (nothing to do)

			// 2b. otherwise, list should delete r2!
			if len(r2.Ranges) == 0 {
				deleteAtIndex(ranges, i)
				i -= 1
				end -= 1
			}

			// 3. continue, since other overlaps may happen
			continue
		}
	}
}

func cloneRange(r *Range) *Range {
	cloned := &Range{
		Start:     r.Start,
		Stop:      r.Stop,
		Ranges:    make([]*Range, 0),
		Dimension: r.Dimension,
		Value:     r.Value,
	}

	for _, rC := range r.Ranges {
		cloned.Ranges = append(cloned.Ranges, cloneRange(rC))
	}

	return cloned
}

func mergeChildRanges(rN *Range, r2 *Range) *Range {
	if rN.Dimension != r2.Dimension {
		log.Debugf("\n\nSOMETHING IS BROKEN. Trying to merge Dimensions that don't match! rN: %s, r2: %s\n\n\n", rN.Dimension, r2.Dimension)
		return nil
	}

	if rN.Start != r2.Start || rN.Stop != r2.Stop {
		log.Debugf("\n\nSOMETHING IS BROKEN. Trying to merge Ranges that haven't been sliced yet! rN: %d..%d, r2: %d..%d\n\n\n", rN.Start, rN.Stop, r2.Start, r2.Stop)
		return nil
	}

	for _, cRange := range rN.Ranges {
		addRange(cRange, &r2.Ranges)
	}
	return r2
}

func negateChildRanges(rN *Range, r2 *Range) {
	if rN.Dimension != r2.Dimension {
		log.Debugf("\n\nSOMETHING IS BROKEN. Trying to negate Dimensions that don't match! rN: %s, r2: %s\n\n\n", rN.Dimension, r2.Dimension)
		return
	}

	for _, cRange := range rN.Ranges {
		subRange(cRange, &r2.Ranges)
	}
}

func insertAtIndex(list *[]*Range, index int, r *Range) {
	if index == len(*list) { // nil, empty slice, or after last element
		*list = append(*list, r)
	} else {
		*list = append((*list)[:index+1], (*list)[index:]...) // index < len(a)
		(*list)[index] = r
	}
}

func deleteAtIndex(list *[]*Range, index int) {
	*list = append((*list)[:index], (*list)[index+1:]...)
}

func part2(input []string) (string, error) {
	rangeList := parseInput2(input)

	// const MIN, MAX = -50, 50

	active := []*Range{rangeList[0]}

	PrintList(active)

	for _, r := range rangeList[1:] {
		if r.Value {
			addRange(r, &active)
		} else {
			subRange(r, &active)
		}
	}

	cubesOn := 0
	for i, r := range active {
		log.Debugf("%d) %s\n", i, r)

		cubesOn += r.Count(1)
	}

	return fmt.Sprintf("%d", cubesOn), nil
}

func PrintList(list []*Range) {
	for i, r := range list {
		log.Debugf("%d) %s\n", i, r)
	}
	log.Debugln("--------------")
}

func parseInput(input []string) []InitStep {
	steps := make([]InitStep, 0, len(input))
	for _, s := range input {
		parts := strings.SplitN(s, " ", 2)
		iStep := InitStep{
			TurnOn: parts[0] == "on",
		}
		ranges := strings.Split(parts[1], ",")
		parseRange := func(in string) Range {
			nums := strings.SplitN(in[2:], "..", 2)
			start, stop := util.Atoi(nums[0]), util.Atoi(nums[1])
			return Range{
				Start: start,
				Stop:  stop,
			}
		}
		iStep.X = parseRange(ranges[0])
		iStep.Y = parseRange(ranges[1])
		iStep.Z = parseRange(ranges[2])
		steps = append(steps, iStep)
	}
	return steps
}

func parseInput2(input []string) []*Range {
	allRanges := make([]*Range, 0, len(input))
	for _, s := range input {
		parts := strings.SplitN(s, " ", 2)
		ranges := strings.Split(parts[1], ",")
		parseRange := func(in string, val bool) Range {
			dim := in[0:1]
			nums := strings.SplitN(in[2:], "..", 2)
			start, stop := util.Atoi(nums[0]), util.Atoi(nums[1])
			return Range{
				Dimension: dim,
				Start:     start,
				Stop:      stop,
				Value:     val,
				Ranges:    make([]*Range, 0),
			}
		}

		on := parts[0] == "on"
		zRange := parseRange(ranges[2], on)
		yRange := parseRange(ranges[1], on)
		yRange.Ranges = append(yRange.Ranges, &zRange)
		xRange := parseRange(ranges[0], on)
		xRange.Ranges = append(xRange.Ranges, &yRange)
		allRanges = append(allRanges, &xRange)
	}
	return allRanges
}

type Grid map[int]map[int]map[int]bool

type Range struct {
	Start     int
	Stop      int
	Ranges    []*Range
	Dimension string
	Value     bool
}

func (r Range) String() string {
	cIndent := "\t"
	cIndent2 := ""
	if r.Dimension == "y" {
		cIndent = "\t"
		cIndent2 = "\t\t"
	}
	mine := fmt.Sprintf("%s: %2d..%2d", r.Dimension, r.Start, r.Stop)
	str := mine
	for i, rC := range r.Ranges {
		if i == 0 {
			str = fmt.Sprintf("%s%s%s", str, cIndent, rC)
		} else {
			str = fmt.Sprintf("%s\n   %s%s%s%s", str, cIndent2, mine, cIndent, rC)
		}
	}
	return str
}

func (r Range) Count(multiplier int) int {
	mine := (r.Stop - r.Start) + 1
	result := 0
	if r.Dimension == "z" {
		result = mine * multiplier
	} else {
		for _, rC := range r.Ranges {
			result += rC.Count(mine * multiplier)
		}
	}
	return result
}

type InitStep struct {
	TurnOn bool
	X      Range
	Y      Range
	Z      Range
}

func init() {
	challenges.RegisterChallengeFunc(2021, 22, 1, "day22.txt", part1)
	challenges.RegisterChallengeFunc(2021, 22, 2, "day22.txt", part2)
}
