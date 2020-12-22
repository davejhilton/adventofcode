package challenges

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

func day20_part1(input []string) (string, error) {

	topLeftCorner := day20_arrangeTiles(input)

	// Now find all the corners, and multiply their IDs

	curTile := topLeftCorner
	log.Debugf("\n- FOUND 1st CORNER: %d\n", curTile.Id)
	result := curTile.Id // store the first corner

	// next move to the top right corner
	for curTile.Neighbors[day20_RIGHT] != nil {
		curTile = curTile.Neighbors[day20_RIGHT]
	}
	log.Debugf("- FOUND 2nd CORNER: %d\n", curTile.Id)
	result *= curTile.Id // store the second corner

	// next move to the bottom right corner
	for curTile.Neighbors[day20_BOTTOM] != nil {
		curTile = curTile.Neighbors[day20_BOTTOM]
	}
	log.Debugf("- FOUND 3rd CORNER: %d\n", curTile.Id)
	result *= curTile.Id // store the third corner

	// lastly, move to the bottom left corner
	for curTile.Neighbors[day20_LEFT] != nil {
		curTile = curTile.Neighbors[day20_LEFT]
	}
	log.Debugf("- FOUND 4th CORNER: %d\n", curTile.Id)
	result *= curTile.Id // store the last corner

	return fmt.Sprintf("%d", result), nil
}

func day20_arrangeTiles(input []string) *day20_tile {
	tiles := day20_parse(input)

	// FIRST, find all the "Neighbor" relations for every tile
	var changed = true
	var firstTile *day20_tile
	for changed {
		changed = false
		for _, tile := range tiles {
			if firstTile == nil {
				px := tile.Orientations[0]
				tile.Pixels = &px
				firstTile = tile
			}
			if tile.Pixels != nil {
				for _, dir := range day20_DIRECTIONS {
					if tile.Neighbors[dir] == nil {
						log.Debugf("Finding Tile %d's %-6s NEIGHBOR... - ", tile.Id, day20_DIR_NAMES[dir])
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
	for curTile.Neighbors[day20_LEFT] != nil {
		curTile = curTile.Neighbors[day20_LEFT]
	}
	for curTile.Neighbors[day20_TOP] != nil {
		curTile = curTile.Neighbors[day20_TOP]
	}

	topLeftCorner := curTile
	log.Debugf("\nFINAL LAYOUT:\n-------------------------\n")
	leftTile := topLeftCorner
	for leftTile != nil {
		curTile = leftTile
		for curTile != nil {
			log.Debugf("  %4d", curTile.Id)
			curTile = curTile.Neighbors[day20_RIGHT]
		}
		log.Debugln()
		leftTile = leftTile.Neighbors[day20_BOTTOM]
	}
	log.Debugln()

	return topLeftCorner
}

func day20_part2(input []string) (string, error) {
	topLeftCorner := day20_arrangeTiles(input)

	log.Debug("\n\n\n\n")
	log.Debug("*****************************************\n")
	log.Debug("*                PART 2                 *\n")
	log.Debug("*****************************************\n")
	log.Debug("\n\n\n\n")

	// First, let's chop off the borders and create our raw image

	tileSize := len(*(topLeftCorner.Pixels))
	pixels := make(day20_pixels, 0)
	leftTile := topLeftCorner
	for leftTile != nil {
		pixelRow := 1
		for ; pixelRow < tileSize-1; pixelRow++ {
			row := make([]rune, 0, tileSize-2)
			curTile := leftTile
			for curTile != nil {
				pix := *curTile.Pixels
				for _, p := range pix[pixelRow][1 : tileSize-1] {
					row = append(row, p)
				}
				curTile = curTile.Neighbors[day20_RIGHT]
			}
			pixels = append(pixels, row)
		}
		leftTile = leftTile.Neighbors[day20_BOTTOM]
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
			monsters = day20_findSeamonsters(pixels)

			if len(monsters) > 0 {
				break
			}

			pixels = day20_rotatePixels90Degrees(pixels)
		}

		if len(monsters) > 0 {
			break
		}
		pixels = day20_reflectPixelsAroundHorizontalAxis(pixels)
	}

	log.Debugf("Found %d Seamonsters!\n", len(monsters))

	result := day20_replaceSeamonstersAndCount(pixels, monsters)

	return fmt.Sprintf("%d", result), nil
}

type day20_pixels [][]rune

func (p day20_pixels) String() string {
	var b strings.Builder
	for i, row := range p {
		b.WriteString(string(row))
		if i != len(p)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (p day20_pixels) GetBorder(direction int) string {
	switch direction {

	case day20_TOP:
		return string(p[0])

	case day20_RIGHT:
		px := make([]rune, 0, len(p))
		for i := 0; i < len(p); i++ {
			px = append(px, p[i][len(p[i])-1])
		}
		return string(px)

	case day20_LEFT:
		px := make([]rune, 0, len(p))
		for i := 0; i < len(p); i++ {
			px = append(px, p[i][0])
		}
		return string(px)

	case day20_BOTTOM:
		return string(p[len(p)-1])
	}
	return ""
}

type day20_tile struct {
	Id           int
	Pixels       *day20_pixels
	Orientations []day20_pixels
	Neighbors    []*day20_tile /* 0 = TOP, 1 = RIGHT, 2 = BOTTOM, 3 = LEFT */
}

func (t day20_tile) FindNeighbor(direction int, tiles map[int]*day20_tile) (*day20_tile, int) {
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

func NewTile(id int, pixels day20_pixels) *day20_tile {
	tile := day20_tile{
		Id:           id,
		Orientations: make([]day20_pixels, 0, 8),
		Neighbors:    make([]*day20_tile, 4),
	}

	tile.Orientations = append(tile.Orientations, pixels)

	for i := 0; i < 7; i++ {
		var newPixels day20_pixels
		if i == 3 {
			newPixels = day20_reflectPixelsAroundHorizontalAxis(pixels)
		} else {
			newPixels = day20_rotatePixels90Degrees(pixels)
		}
		tile.Orientations = append(tile.Orientations, newPixels)
		pixels = newPixels
	}
	return &tile
}

func day20_rotatePixels90Degrees(pixels day20_pixels) day20_pixels {
	newPixels := make(day20_pixels, len(pixels))
	for i := 0; i < len(pixels); i++ {
		newPixels[i] = make([]rune, len(pixels))
		for j := 0; j < len(pixels); j++ {
			newPixels[i][j] = pixels[len(pixels)-j-1][i]
		}
	}
	return newPixels
}

func day20_reflectPixelsAroundHorizontalAxis(pixels day20_pixels) day20_pixels {
	newPixels := make(day20_pixels, len(pixels))
	copy(newPixels, pixels)
	for j := 0; j < len(newPixels)/2; j++ {
		newPixels[j], newPixels[len(newPixels)-j-1] = newPixels[len(newPixels)-j-1], newPixels[j]
	}
	return newPixels
}

func day20_findSeamonsters(pixels day20_pixels) [][]int {

	monsters := make([][]int, 0)
	for i := 0; i < len(pixels)-2; i++ {
		left := 0
		for left < len(pixels[i])-20 {
			if m := day20_seamonsterRegex_row2.FindAllStringIndex(string(pixels[i+1])[left:], 1); m != nil {
				if n := day20_seamonsterRegex_row1.FindAllStringIndex(string(pixels[i][left+m[0][0]:left+m[0][1]]), 1); n != nil {
					if o := day20_seamonsterRegex_row3.FindAllStringIndex(string(pixels[i+2][left+m[0][0]:left+m[0][1]]), 1); o != nil {
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

func day20_replaceSeamonstersAndCount(pixels day20_pixels, monsters [][]int) int {

	for i := 0; i < len(monsters); i++ {
		row, left := monsters[i][0], monsters[i][1]
		for j := 0; j < len(day20_seamonsterPieces); j++ {
			for _, k := range day20_seamonsterPieces[j] {
				pixels[row+j][left+k] = 'O'
			}
		}
	}

	log.Debug("\n\n")
	log.Debugln(pixels)
	log.Debug("\n\n")

	var count int
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			if pixels[i][j] == '#' {
				count++
			}
		}
	}

	return count
}

const day20_TOP = 0
const day20_RIGHT = 1
const day20_BOTTOM = 2
const day20_LEFT = 3

var day20_DIRECTIONS = []int{day20_TOP, day20_RIGHT, day20_BOTTOM, day20_LEFT}
var day20_DIR_NAMES = []string{"TOP", "RIGHT", "BOTTOM", "LEFT"}

var day20_seamonsterRegex_row1 = regexp.MustCompile(`..................#.`)
var day20_seamonsterRegex_row2 = regexp.MustCompile(`#....##....##....###`)
var day20_seamonsterRegex_row3 = regexp.MustCompile(`.#..#..#..#..#..#...`)
var day20_seamonsterPieces = [][]int{
	[]int{18},
	[]int{0, 5, 6, 11, 12, 17, 18, 19},
	[]int{1, 4, 7, 10, 13, 16},
}

func day20_parse(input []string) map[int]*day20_tile {
	tiles := make(map[int]*day20_tile)

	var i int
	for i < len(input) {
		var id int
		fmt.Sscanf(input[i], "Tile %d:", &id)
		i++
		pixels := make(day20_pixels, 0)
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
	registerChallengeFunc(20, 1, "day20.txt", day20_part1)
	registerChallengeFunc(20, 2, "day20.txt", day20_part2)
}
