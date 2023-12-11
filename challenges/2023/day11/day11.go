package aoc2023_day11

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Image [][]string

func (i Image) String() string {
	var sb strings.Builder
	for _, row := range i {
		for _, c := range row {
			sb.WriteString(c)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type SparseImage struct {
	Rows   map[int][]int
	Height int
	Width  int
}

func (si SparseImage) String() string {
	var sb strings.Builder
	for r := 0; r < si.Height; r++ {
		if row, ok := si.Rows[r]; ok {
			for i := 0; i < si.Width; i++ {
				if util.Contains(row, i) {
					sb.WriteString("#")
				} else {
					sb.WriteString(".")
				}
			}
		} else {
			for i := 0; i < si.Width; i++ {
				sb.WriteString(".")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func sparseImageFromImage(image Image) SparseImage {
	si := SparseImage{
		Rows:   make(map[int][]int),
		Height: len(image),
		Width:  len(image[0]),
	}
	for r, row := range image {
		for c, col := range row {
			if col == "#" {
				si.Rows[r] = append(si.Rows[r], c)
			}
		}
	}
	return si
}

type Coord struct {
	Row int
	Col int
}

func (c Coord) Equals(other Coord) bool {
	return c.Row == other.Row && c.Col == other.Col
}

type CoordPair struct {
	A Coord
	B Coord
}

func (cp CoordPair) MinDistance() int {
	return util.Abs(cp.A.Row-cp.B.Row) + util.Abs(cp.A.Col-cp.B.Col)
}

func expandImage(image Image, replaceWith int) SparseImage {
	expandBy := replaceWith - 1
	si := sparseImageFromImage(image)
	expanded := SparseImage{
		Rows:   make(map[int][]int),
		Height: len(image),
		Width:  len(image[0]),
	}

	indexMap := make(map[int]int)

	modIndex := 0
	for i := 0; i < len(image); i++ {
		if _, ok := si.Rows[i]; ok {
			indexMap[i] = modIndex
			expanded.Rows[modIndex] = make([]int, len(si.Rows[i]))
			copy(expanded.Rows[modIndex], si.Rows[i])
		} else {
			modIndex += expandBy
			expanded.Height += expandBy
		}
		modIndex++
	}

	log.Debugf("Index Map: %v\n", indexMap)

	for c := len(image[0]) - 1; c >= 0; c-- {
		hasGalaxy := false
		for _, row := range image {
			if row[c] == "#" {
				hasGalaxy = true
				break
			}
		}
		if !hasGalaxy {
			log.Debugf("Expanding column %d by %d\n", c, expandBy)
			expanded.Width += expandBy
			for i := range image {
				if modIndex, ok := indexMap[i]; ok {
					eRow := expanded.Rows[modIndex]
					for j := len(eRow) - 1; j >= 0; j-- {
						if eRow[j] > c {
							log.Debugf("Moving row %d (%d)'s column %d up by %d\n", i, modIndex, eRow[j], expandBy)
							eRow[j] += expandBy
						}
					}
				}
			}
		}
	}

	return expanded
}

func part1(input []string) (string, error) {
	image := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", image)

	sparse := expandImage(image, 2)
	log.Debugf("Expanded Input:\n%s\n", sparse)

	coords := make([]Coord, 0)
	for r, row := range sparse.Rows {
		for _, c := range row {
			coords = append(coords, Coord{Row: r, Col: c})
		}
	}

	// find all combinations of coords
	coordPairs := make([]CoordPair, 0)
	for i, c1 := range coords {
		for j := i + 1; j < len(coords); j++ {
			c2 := coords[j]
			coordPairs = append(coordPairs, CoordPair{A: c1, B: c2})
		}
	}

	log.Debugf("Coord Pairs:\n%d\n", len(coordPairs))
	var result int

	for _, cp := range coordPairs {
		result += cp.MinDistance()
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	image := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", image)

	si := expandImage(image, 1000000)

	log.Debugf("Sparse Image:\n%s\n", si)

	coords := make([]Coord, 0)
	for r, row := range si.Rows {
		for _, c := range row {
			coords = append(coords, Coord{Row: r, Col: c})
		}
	}

	// find all combinations of coords
	coordPairs := make([]CoordPair, 0)
	for i, c1 := range coords {
		for j := i + 1; j < len(coords); j++ {
			c2 := coords[j]
			coordPairs = append(coordPairs, CoordPair{A: c1, B: c2})
		}
	}

	var result int
	for _, cp := range coordPairs {
		result += cp.MinDistance()
	}

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) Image {
	image := make(Image, 0, len(input))
	for _, s := range input {
		image = append(image, strings.Split(s, ""))
	}
	return image
}

func init() {
	challenges.RegisterChallengeFunc(2023, 11, 1, "day11.txt", part1)
	challenges.RegisterChallengeFunc(2023, 11, 2, "day11.txt", part2)
}
