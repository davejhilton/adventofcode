package challenges2020

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func day24_part1(input []string) (string, error) {
	paths := day24_parse(input)
	blackTiles := make(map[string]bool)

	for _, p := range paths {
		tile := day24_coord{0, 0, 0}
		// log.Debugf("\tTILE %-8s\n", tile)
		for _, d := range p {
			tile.ApplyDelta(d)
			// log.Debugf("\t --> %-8s\n", tile)
		}
		if _, ok := blackTiles[tile.String()]; ok {
			log.Debugf("TILE: %-8s --> WHITE\n", tile)
			delete(blackTiles, tile.String())
		} else {
			log.Debugf("TILE: %-8s --> BLACK\n", tile)
			blackTiles[tile.String()] = true
		}
	}

	return fmt.Sprintf("%d", len(blackTiles)), nil
}

func day24_part2(input []string) (string, error) {
	paths := day24_parse(input)

	// var minX, minY, minZ int = day24_MAX_INT, day24_MAX_INT, day24_MAX_INT
	// var maxX, maxY, maxZ int = day24_MIN_INT, day24_MIN_INT, day24_MIN_INT

	blackTiles := make(map[string]bool)

	for _, p := range paths {
		tile := day24_coord{0, 0, 0}
		for _, d := range p {
			tile.ApplyDelta(d)
		}
		if _, ok := blackTiles[tile.String()]; ok {
			delete(blackTiles, tile.String())
		} else {
			blackTiles[tile.String()] = true
		}
		// if tile[0] < minX {
		// 	minX = tile[0]
		// }
		// if tile[0] > maxX {
		// 	maxX = tile[0]
		// }
		// if tile[1] < minY {
		// 	minY = tile[1]
		// }
		// if tile[1] > maxY {
		// 	maxY = tile[1]
		// }
		// if tile[2] < minZ {
		// 	minZ = tile[2]
		// }
		// if tile[2] > maxZ {
		// 	maxZ = tile[2]
		// }
	}

	var x, y, z int
	c := day24_coord{x, y, z}
	for day := 0; day < 100; day++ {
		newBlackTiles := make(map[string]bool)

		for k := range blackTiles {
			fmt.Sscanf(k, "%3d,%3d,%3d", &x, &y, &z)
			c[0], c[1], c[2] = x, y, z
			// fmt.Printf("\t%s\n", c)
			day24_handleNeighborFlips(c, blackTiles, &newBlackTiles)
		}
		log.Debugf("Day %d: %d\n", day+1, len(newBlackTiles))
		blackTiles = newBlackTiles
	}

	return fmt.Sprintf("%d", len(blackTiles)), nil
}

func day24_parse(input []string) []day24_path {
	paths := make([]day24_path, 0, len(input))
	for _, s := range input {
		p := make(day24_path, 0)
		var i int
		for i < len(s) {
			c := s[i]
			switch c {
			case 'e':
				p = append(p, day24_coord{1, -1, 0})
			case 'w':
				p = append(p, day24_coord{-1, 1, 0})
			case 's':
				i++
				c = s[i]
				if c == 'e' {
					p = append(p, day24_coord{0, -1, 1})
				} else if c == 'w' {
					p = append(p, day24_coord{-1, 0, 1})
				}
			case 'n':
				i++
				c = s[i]
				if c == 'e' {
					p = append(p, day24_coord{1, 0, -1})
				} else if c == 'w' {
					p = append(p, day24_coord{0, 1, -1})
				}
			}
			i++
		}
		paths = append(paths, p)
	}
	return paths
}

func day24_countBlackNeighborTiles(c day24_coord, blackTiles map[string]bool) int {
	deltas := []day24_coord{
		{1, -1, 0}, // e
		{0, -1, 1}, // se
		{-1, 0, 1}, // sw
		{-1, 1, 0}, // w
		{0, 1, -1}, // nw
		{1, 0, -1}, // ne
	}

	count := 0
	n := day24_coord{0, 0, 0}
	for _, d := range deltas {
		n[0] = c[0] + d[0]
		n[1] = c[1] + d[1]
		n[2] = c[2] + d[2]
		if _, ok := blackTiles[n.String()]; ok {
			count++
		}
	}
	return count
}

func day24_handleNeighborFlips(c day24_coord, cur map[string]bool, next *map[string]bool) {
	tiles := []day24_coord{
		c,
		{c[0] + 1, c[1] - 1, c[2] + 0}, // e
		{c[0] + 0, c[1] - 1, c[2] + 1}, // se
		{c[0] - 1, c[1] + 0, c[2] + 1}, // sw
		{c[0] - 1, c[1] + 1, c[2] + 0}, // w
		{c[0] + 0, c[1] + 1, c[2] - 1}, // nw
		{c[0] + 1, c[1] + 0, c[2] - 1}, // ne
	}
	var n int
	var key string
	for _, t := range tiles {
		n = day24_countBlackNeighborTiles(t, cur)
		key = t.String()
		if _, ok := cur[key]; ok {
			if n == 0 || n > 2 {
				delete(*next, key)
			} else {
				(*next)[key] = true
			}
		} else {
			if n == 2 {
				(*next)[key] = true
			}
		}
	}
}

type day24_coord [3]int

func (c *day24_coord) ApplyDelta(delta day24_coord) {
	c[0] += delta[0]
	c[1] += delta[1]
	c[2] += delta[2]
}

func (c day24_coord) String() string {
	return fmt.Sprintf("%3d,%3d,%3d", c[0], c[1], c[2])
}

type day24_path []day24_coord

const (
	day24_MAX_INT = int((^uint(0)) >> 1)
	day24_MIN_INT = -1*day24_MAX_INT - 1
)

/*

		 x,  y,  z

e	->	 1, -1,  0
se	->	 0, -1,  1
sw	->	-1,  0,  1
w	->	-1,  1,  0
nw	->	 0,  1, -1
ne	->	 1,  0, -1

*/

func init() {
	challenges.RegisterChallengeFunc(2020, 24, 1, "day24.txt", day24_part1)
	challenges.RegisterChallengeFunc(2020, 24, 2, "day24.txt", day24_part2)
}
