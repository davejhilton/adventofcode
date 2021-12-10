package aoc2020_day20

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	topLeftCorner := arrangeTiles(input)

	// Now find all the corners, and multiply their IDs

	curTile := topLeftCorner
	log.Debugf("\n- FOUND 1st CORNER: %d\n", curTile.Id)
	result := curTile.Id // store the first corner

	// next move to the top right corner
	for curTile.Neighbors[RIGHT] != nil {
		curTile = curTile.Neighbors[RIGHT]
	}
	log.Debugf("- FOUND 2nd CORNER: %d\n", curTile.Id)
	result *= curTile.Id // store the second corner

	// next move to the bottom right corner
	for curTile.Neighbors[BOTTOM] != nil {
		curTile = curTile.Neighbors[BOTTOM]
	}
	log.Debugf("- FOUND 3rd CORNER: %d\n", curTile.Id)
	result *= curTile.Id // store the third corner

	// lastly, move to the bottom left corner
	for curTile.Neighbors[LEFT] != nil {
		curTile = curTile.Neighbors[LEFT]
	}
	log.Debugf("- FOUND 4th CORNER: %d\n", curTile.Id)
	result *= curTile.Id // store the last corner

	return fmt.Sprintf("%d", result), nil
}

func arrangeTiles(input []string) *tile {
	tiles := parse(input)

	// FIRST, find all the "Neighbor" relations for every tile
	changed := true
	var firstTile *tile
	for changed {
		changed = false
		for _, tile := range tiles {
			if firstTile == nil {
				px := tile.Orientations[0]
				tile.Pixels = &px
				firstTile = tile
			}
			if tile.Pixels != nil {
				for _, dir := range DIRECTIONS {
					if tile.Neighbors[dir] == nil {
						log.Debugf("Finding Tile %d's %-6s NEIGHBOR... - ", tile.Id, DIR_NAMES[dir])
						neighbor, o := tile.FindNeighbor(dir, tiles)
						if neighbor != nil {
							log.Debugf("found %4d\n", neighbor.Id)
							if neighbor.Pixels == nil && o >= 0 {
								px := neighbor.Orientations[o]
								neighbor.Pixels = &px
							}
							tile.Neighbors[dir] = neighbor
							neighbor.Neighbors[(dir+2)%4] = tile
							changed = true
						} else {
							log.Debugln("EMPTY")
						}
					}
				}
			}
		}
	}

	// find the top left corner...
	curTile := firstTile
	for curTile.Neighbors[LEFT] != nil {
		curTile = curTile.Neighbors[LEFT]
	}
	for curTile.Neighbors[TOP] != nil {
		curTile = curTile.Neighbors[TOP]
	}

	topLeftCorner := curTile
	log.Debugf("\nFINAL LAYOUT:\n-------------------------\n")
	leftTile := topLeftCorner
	for leftTile != nil {
		curTile = leftTile
		for curTile != nil {
			log.Debugf("  %4d", curTile.Id)
			curTile = curTile.Neighbors[RIGHT]
		}
		log.Debugln()
		leftTile = leftTile.Neighbors[BOTTOM]
	}
	log.Debugln()

	return topLeftCorner
}

func part2(input []string) (string, error) {
	topLeftCorner := arrangeTiles(input)

	log.Debug("\n\n\n\n")
	log.Debug("*****************************************\n")
	log.Debug("*                PART 2                 *\n")
	log.Debug("*****************************************\n")
	log.Debug("\n\n\n\n")

	// First, let's chop off the borders and create our raw image

	tileSize := len(*(topLeftCorner.Pixels))
	pixels := make(pixels, 0)
	leftTile := topLeftCorner
	for leftTile != nil {
		pixelRow := 1
		for ; pixelRow < tileSize-1; pixelRow++ {
			row := make([]rune, 0, tileSize-2)
			curTile := leftTile
			for curTile != nil {
				pix := *curTile.Pixels
				row = append(row, pix[pixelRow][1:tileSize-1]...)
				curTile = curTile.Neighbors[RIGHT]
			}
			pixels = append(pixels, row)
		}
		leftTile = leftTile.Neighbors[BOTTOM]
	}

	log.Debug("Pixels without borders:\n\n")
	log.Debugln(pixels)
	log.Debug("\n\n\n\n")

	var monsters [][]int
	for flipped := 0; flipped <= 1; flipped++ {
		for rotated := 0; rotated <= 3; rotated++ {

			log.Debug("\n\n")
			log.Debugln(pixels)
			log.Debug("\n\n")
			monsters = findSeamonsters(pixels)

			if len(monsters) > 0 {
				break
			}

			pixels = rotatePixels90Degrees(pixels)
		}

		if len(monsters) > 0 {
			break
		}
		pixels = reflectPixelsAroundHorizontalAxis(pixels)
	}

	log.Debugf("Found %d Seamonsters!\n", len(monsters))

	result := replaceSeamonstersAndCount(pixels, monsters)

	return fmt.Sprintf("%d", result), nil
}

type pixels [][]rune

func (p pixels) String() string {
	var b strings.Builder
	for i, row := range p {
		b.WriteString(string(row))
		if i != len(p)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (p pixels) GetBorder(direction int) string {
	switch direction {

	case TOP:
		return string(p[0])

	case RIGHT:
		px := make([]rune, 0, len(p))
		for i := 0; i < len(p); i++ {
			px = append(px, p[i][len(p[i])-1])
		}
		return string(px)

	case LEFT:
		px := make([]rune, 0, len(p))
		for i := 0; i < len(p); i++ {
			px = append(px, p[i][0])
		}
		return string(px)

	case BOTTOM:
		return string(p[len(p)-1])
	}
	return ""
}

type tile struct {
	Id           int
	Pixels       *pixels
	Orientations []pixels
	Neighbors    []*tile /* 0 = TOP, 1 = RIGHT, 2 = BOTTOM, 3 = LEFT */
}

func (t tile) FindNeighbor(direction int, tiles map[int]*tile) (*tile, int) {
	if t.Pixels == nil {
		panic("Can't find neighbor for a tile without a pinned orientation!")
	}
	border := t.Pixels.GetBorder(direction)

	oppositeDirection := (direction + 2) % 4

	for _, tile := range tiles {
		if tile.Pixels != nil {
			if tile.Pixels.GetBorder(oppositeDirection) == border {
				return tile, -1
			}
		} else {
			for o, pixels := range tile.Orientations {
				if pixels.GetBorder(oppositeDirection) == border {
					return tile, o
				}
			}
		}
	}
	return nil, -1
}

func NewTile(id int, pxs pixels) *tile {
	tile := tile{
		Id:           id,
		Orientations: make([]pixels, 0, 8),
		Neighbors:    make([]*tile, 4),
	}

	tile.Orientations = append(tile.Orientations, pxs)

	for i := 0; i < 7; i++ {
		var newPixels pixels
		if i == 3 {
			newPixels = reflectPixelsAroundHorizontalAxis(pxs)
		} else {
			newPixels = rotatePixels90Degrees(pxs)
		}
		tile.Orientations = append(tile.Orientations, newPixels)
		pxs = newPixels
	}
	return &tile
}

func rotatePixels90Degrees(pxs pixels) pixels {
	newPixels := make(pixels, len(pxs))
	for i := 0; i < len(pxs); i++ {
		newPixels[i] = make([]rune, len(pxs))
		for j := 0; j < len(pxs); j++ {
			newPixels[i][j] = pxs[len(pxs)-j-1][i]
		}
	}
	return newPixels
}

func reflectPixelsAroundHorizontalAxis(pxs pixels) pixels {
	newPixels := make(pixels, len(pxs))
	copy(newPixels, pxs)
	for j := 0; j < len(newPixels)/2; j++ {
		newPixels[j], newPixels[len(newPixels)-j-1] = newPixels[len(newPixels)-j-1], newPixels[j]
	}
	return newPixels
}

func findSeamonsters(pxs pixels) [][]int {
	monsters := make([][]int, 0)
	for i := 0; i < len(pxs)-2; i++ {
		left := 0
		for left < len(pxs[i])-20 {
			if m := seamonsterRegex_row2.FindAllStringIndex(string(pxs[i+1])[left:], 1); m != nil {
				if n := seamonsterRegex_row1.FindAllStringIndex(string(pxs[i][left+m[0][0]:left+m[0][1]]), 1); n != nil {
					if o := seamonsterRegex_row3.FindAllStringIndex(string(pxs[i+2][left+m[0][0]:left+m[0][1]]), 1); o != nil {
						log.Debugf("Found Seamonster at %d,%d\n", i, left+m[0][0])
						monsters = append(monsters, []int{i, left + m[0][0], left + m[0][1]})
					}
				}
				left = left + m[0][0] + 1
			} else {
				break
			}
		}
	}

	return monsters
}

func replaceSeamonstersAndCount(pxs pixels, monsters [][]int) int {
	for i := 0; i < len(monsters); i++ {
		row, left := monsters[i][0], monsters[i][1]
		for j := 0; j < len(seamonsterPieces); j++ {
			for _, k := range seamonsterPieces[j] {
				pxs[row+j][left+k] = 'O'
			}
		}
	}

	log.Debug("\n\n")
	log.Debugln(pxs)
	log.Debug("\n\n")

	var count int
	for i := 0; i < len(pxs); i++ {
		for j := 0; j < len(pxs[i]); j++ {
			if pxs[i][j] == '#' {
				count++
			}
		}
	}

	return count
}

const (
	TOP    = 0
	RIGHT  = 1
	BOTTOM = 2
	LEFT   = 3
)

var (
	DIRECTIONS = []int{TOP, RIGHT, BOTTOM, LEFT}
	DIR_NAMES  = []string{"TOP", "RIGHT", "BOTTOM", "LEFT"}
)

var (
	seamonsterRegex_row1 = regexp.MustCompile(`..................#.`)
	seamonsterRegex_row2 = regexp.MustCompile(`#....##....##....###`)
	seamonsterRegex_row3 = regexp.MustCompile(`.#..#..#..#..#..#...`)
	seamonsterPieces     = [][]int{
		{18},
		{0, 5, 6, 11, 12, 17, 18, 19},
		{1, 4, 7, 10, 13, 16},
	}
)

func parse(input []string) map[int]*tile {
	tiles := make(map[int]*tile)

	var i int
	for i < len(input) {
		var id int
		fmt.Sscanf(input[i], "Tile %d:", &id)
		i++
		pixels := make(pixels, 0)
		for ; i < len(input) && input[i] != ""; i++ {
			row := make([]rune, 0)
			for _, p := range input[i] {
				row = append(row, rune(p))
			}
			pixels = append(pixels, row)
		}
		i++
		tiles[id] = NewTile(id, pixels)
	}
	return tiles
}

func init() {
	challenges.RegisterChallengeFunc(2020, 20, 1, "day20.txt", part1)
	challenges.RegisterChallengeFunc(2020, 20, 2, "day20.txt", part2)
}
