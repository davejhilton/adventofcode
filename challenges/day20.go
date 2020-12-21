package challenges

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

var tilesPerRow int
var pixelsPerRow int

var tiles map[int]*day20_tile
var positions [][]int

func day20_part1(input []string) (string, error) {
	pixelsPerRow = len(input[1])
	tiles = day20_parse(input)

	tilesPerRow = int(math.Sqrt(float64(len(tiles))))

	// for id, tile := range tiles {
	// 	log.Debugf("Tile %d:\n%s\n", id, tile)
	// }

	// 1. start with whichever tile...
	// 2. match left and move left until nothing else matches left
	// 3. match up and move up until nothing else matches up
	// 		now we have our top left corner
	// 4. match right and move right until nothing else matches right
	//		now we have our top right corner (and top row)
	// 5. match down and move down until nothing else matches
	//		now we have our bottom right corner
	// 6. match left and move left until nothing else matches
	//		now we have our bottom left corner
	// 7. match up and move up and verify that we hit our top left corner when we expect to

	var curTile *day20_tile
	// curTile = tiles[3221]
	// grab a single tile from our map, as our arbitrary starting point
	for id := range tiles {
		curTile = tiles[id]
		log.Debugf("Starting with tile: %d\n", id)
		break
	}

	// curTile = tiles[3079]
	curTile.OrientationLocked = true

	lOrder := make([]int, 0)
	lOrder = append(lOrder, curTile.Id)
	// match tiles to the left until we find the left edge of our "image"
	for {
		if id, rotation, flip, ok := day20_findNeighbor(day20_LEFT, curTile.Id); ok {
			log.Debugf("FOUND MATCH: id=%d\n\n", id)
			day20_addNeighbor(curTile.Id, id, day20_LEFT, rotation, flip)

			lOrder = append(lOrder, id)
			curTile = tiles[id]
		} else {
			break //we found the right edge
		}
	}

	log.Debugln("HIT THE LEFT EDGE!")
	for i := len(lOrder) - 1; i >= 0; i-- {
		if i != len(lOrder)-1 {
			log.Debug(" <-- ")
		}
		log.Debugf("%d", lOrder[i])
	}
	log.Debugln()

	var topLeftTile *day20_tile

	tOrder := make([]int, 0)

	for {
		if id, rotation, flip, ok := day20_findNeighbor(day20_TOP, curTile.Id); ok {
			log.Debugf("FOUND MATCH: id=%d\n\n", id)
			log.Debugln(tiles[id])
			log.Debug("\n\n")
			day20_addNeighbor(curTile.Id, id, day20_TOP, rotation, flip)

			tOrder = append(tOrder, id)
			curTile = tiles[id]
		} else {
			break //we found the TOP edge
		}
	}

	topLeftTile = curTile

	log.Debugln("HIT THE TOP LEFT CORNER!")
	for i := len(tOrder) - 1; i >= 0; i-- {
		if i != len(tOrder)-1 {
			log.Debug(" ^\n |\n |\n")
		}
		log.Debugf("%d\n", tOrder[i])
	}
	if len(tOrder) != 0 {
		log.Debug(" ^\n |\n |\n")
	}
	for i := len(lOrder) - 1; i >= 0; i-- {
		if i != len(lOrder)-1 {
			log.Debug(" <-- ")
		}
		log.Debugf("%d", lOrder[i])
	}
	log.Debugln()

	// NOW: start figuring crap out top left to bottom right

	positions = make([][]int, tilesPerRow)
	for i := range positions {
		positions[i] = make([]int, tilesPerRow)
	}

	positions[0][0] = topLeftTile.Id

	log.Debugf("\n\n*********************\n")
	for _, row := range positions {
		for _, id := range row {
			log.Debugf("%4d\t", id)
		}
		log.Debugln()
	}
	log.Debugf("*********************\n\n")

	for i := 1; i < tilesPerRow; i++ {
		if id, rotation, flip, ok := day20_findNeighbor(day20_BOTTOM, curTile.Id); ok {
			log.Debugf("FOUND MATCH: id=%d\n\n", id)
			log.Debugln(tiles[id])
			log.Debug("\n\n")
			day20_addNeighbor(curTile.Id, id, day20_BOTTOM, rotation, flip)

			curTile = tiles[id]
			positions[i][0] = id
		} else {
			log.Debugf("BREAKING FROM top-down loop at index: %d\n", i)
			log.Debugf("NOTHING BELOW Tile %d\n", curTile.Id)
			break //we found the right edge
		}
	}

	for i, row := range positions {
		if row[0] == 0 {
			log.Debugf("EMPTY ID: BREAKING FROM left-right loop at index: %d\n", i)
			break
		}
		curTile = tiles[row[0]]
		for j := 1; j < tilesPerRow; j++ {
			if id, rotation, flip, ok := day20_findNeighbor(day20_RIGHT, curTile.Id); ok {
				log.Debugf("FOUND MATCH: id=%d\n\n", id)
				log.Debugln(tiles[id])
				log.Debug("\n\n")
				day20_addNeighbor(curTile.Id, id, day20_RIGHT, rotation, flip)

				lOrder = append(lOrder, id)
				curTile = tiles[id]
				positions[i][j] = id
			} else {
				log.Debugf("BREAKING FROM left-right loop at index: %d,%d\n", i, j)
				log.Debugf("NOTHING RIGHT of Tile %d\n", curTile.Id)
				break //we found the right edge
			}
		}
	}

	product := 1
	log.Debugf("\n\n*********************\n")
	for i, row := range positions {
		for j, id := range row {
			log.Debugf("%4d\t", id)
			if i == 0 || i == tilesPerRow-1 {
				if j == 0 || j == tilesPerRow-1 {
					product *= positions[i][j]
				}
			}
		}
		log.Debugln()
	}

	return fmt.Sprintf("%d", product), nil
}

func day20_part2(input []string) (string, error) {

	day20_part1(input)

	log.Debug("\n\n\n\n")
	log.Debug("*****************************************\n")
	log.Debug("*                PART 2                 *\n")
	log.Debug("*****************************************\n")
	log.Debug("\n\n\n\n")

	//(pixelsPerRow-2)*tilesPerRow

	pixels := make([][]rune, 0, pixelsPerRow-2)
	for tilesRow := 0; tilesRow < len(positions); tilesRow++ {
		pixelRow := 1
		for ; pixelRow < pixelsPerRow-1; pixelRow++ {
			row := make([]rune, 0, pixelsPerRow-2)
			for tilesCol := 0; tilesCol < tilesPerRow; tilesCol++ {
				pix := tiles[positions[tilesRow][tilesCol]].Pixels
				for _, p := range pix[pixelRow][1 : pixelsPerRow-1] {
					row = append(row, p)
				}
			}
			pixels = append(pixels, row)
		}
	}

	log.Debugln(day20_pixelsToString(pixels))
	log.Debug("\n\n\n\n")

	var monsters [][]int
	for flipped := 0; flipped <= 1; flipped++ {
		for rotated := 0; rotated <= 3; rotated++ {

			log.Debug("\n\n")
			log.Debugln(day20_pixelsToString(pixels))
			log.Debug("\n\n")
			monsters = day20_findSeamonsters(pixels)

			if len(monsters) > 0 {
				break
			}

			pixels = day20_rotatePixels(pixels, 90)
		}

		if len(monsters) > 0 {
			break
		}
		pixels = day20_flipPixels(pixels, false)
	}

	log.Debugf("Found %d Seamonsters!\n", len(monsters))

	result := day20_replaceSeamonstersAndCount(pixels, monsters)

	return fmt.Sprintf("%d", result), nil
}

func day20_pixelsToString(pixels [][]rune) string {
	var b strings.Builder
	for i, row := range pixels {
		b.WriteString(string(row))
		if i != len(pixels)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func day20_findSeamonsters(pixels [][]rune) [][]int {

	monsters := make([][]int, 0)
	for i := 0; i < len(pixels)-2; i++ {
		s := 0
		for s < len(pixels[i])-20 {
			if m := day20_seamonsterRegex_row2.FindAllStringIndex(string(pixels[i+1])[s:], 1); m != nil {
				if n := day20_seamonsterRegex_row1.FindAllStringIndex(string(pixels[i][s+m[0][0]:s+m[0][1]]), 1); n != nil {
					if o := day20_seamonsterRegex_row3.FindAllStringIndex(string(pixels[i+2][s+m[0][0]:s+m[0][1]]), 1); o != nil {
						log.Debugf("found row3 match at %d,%d\n", i, s+m[0][0])
						monsters = append(monsters, []int{i, s + m[0][0], s + m[0][1]})
					}
				}
				s = s + m[0][0] + 1
			} else {
				break
			}
		}
	}

	return monsters
}

func day20_replaceSeamonstersAndCount(pixels [][]rune, monsters [][]int) int {

	for i := 0; i < len(monsters); i++ {
		r, s := monsters[i][0], monsters[i][1]
		for j := 0; j < len(day20_seamonsterReplacements); j++ {
			for _, n := range day20_seamonsterReplacements[j] {
				pixels[r+j][s+n] = 'O'
			}
		}
	}

	// for i := 0; i < len(pixels); i++ {
	// 	for j := 0; j < len(pixels); j++ {
	// 		if pixels[i][j] == '#' {
	// 			pixels[i][j] = '.'
	// 		} else if pixels[i][j] == '.' {
	// 			pixels[i][j] = ' '
	// 		}
	// 	}
	// }

	log.Debug("\n\n")
	log.Debugln(day20_pixelsToString(pixels))
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

/*
                  #
#    ##    ##    ###
 #  #  #  #  #  #
*/

var day20_seamonsterRegex_row1 = regexp.MustCompile(`..................#.`)
var day20_seamonsterRegex_row2 = regexp.MustCompile(`#....##....##....###`)
var day20_seamonsterRegex_row3 = regexp.MustCompile(`.#..#..#..#..#..#...`)
var day20_seamonsterReplacements = [][]int{
	[]int{18},
	[]int{0, 5, 6, 11, 12, 17, 18, 19},
	[]int{1, 4, 7, 10, 13, 16},
}

func day20_addNeighbor(id1 int, id2 int, id1Border int, id2Rotation int, id2Flip bool) {
	tile1 := tiles[id1]
	tile2 := tiles[id2]

	if tile1.Neighbors[id1Border] == tile2 {
		return
	} else if tile1.Neighbors[id1Border] != nil {
		panic(fmt.Sprintf("ALREADY USED THIS BORDER! id1=%d, id2=%d, id1Border=%d, id2Rotation=%d, id2Flip=%v, neighbor=%d\n", id1, id2, id1Border, id2Rotation, id2Flip, tile1.Neighbors[id1Border].Id))
	}

	if id2Rotation != 0 {
		log.Debugf("ROTATING %d DEGREES:\n\n", id2Rotation)
		tile2.Rotate(id2Rotation)
		log.Debugf("%s\n\n", tile2)
		log.Debugf("%s\n\n", day20_bordersToString(tile2.Borders))
	}
	if id2Flip {
		log.Debug("FLIPPING:\n\n")
		tile2.Flip(id1Border%2 == 0)
		log.Debugf("%s\n\n", tile2)
		log.Debugf("%s\n\n", day20_bordersToString(tile2.Borders))
	}
	tile2.OrientationLocked = true
	tile1.Neighbors[id1Border] = tile2
	tile2.Neighbors[(id1Border+2)%4] = tile1
}

func day20_findNeighbor(dir int, tileId int) (id int, rotation int, flip bool, ok bool) {
	log.Debugf("**********\nFINDING %s NEIGHBOR FOR %d...\n", dirToName(dir), tileId)
	tile := tiles[tileId]
	if tile.Neighbors[dir] != nil {
		log.Debugf("- FOUND CACHED NEIGHBOR: %d\n", tile.Neighbors[dir].Id)
		return tile.Neighbors[dir].Id, 0, false, true
	}

	border := tile.Borders[dir]

	var borderPos int
	var reversed bool
	var matched bool
	for tId, t := range tiles {
		if tId == tileId {
			continue
		}

		if t.OrientationLocked {
			oppositeBorder := (dir + 2) % 4
			log.Debugf("********** ORIENTATION IS LOCKED! ID=%d, dir=%d, opp=%d\n", tId, dir, oppositeBorder)
			if t.Borders[oppositeBorder] == border {
				if t.Neighbors[oppositeBorder] != nil {
					panic(fmt.Sprintf("ALREADY USED THIS BORDER! id=%d, tileId=%d, dir=%d, bdr=%d, neighborId=%d\n", id, tileId, dir, oppositeBorder, t.Neighbors[oppositeBorder].Id))
				}
				return tId, 0, false, true
				// log.Debugf("FOUND MATCH for %d(%s): %d(%s)\n", tileId, dirToName(dir), tId, dirToName(pos))
			} else {
				continue
			}
		}

		for pos, b := range t.Borders {
			if b == border {
				log.Debugf("FOUND MATCH for %d(%s): %d(%s)\n", tileId, dirToName(dir), tId, dirToName(pos))
				matched = true
				id = tId
				borderPos = pos
				reversed = false
				break
			} else if day20_strRev(b) == border {
				log.Debugf("FOUND MATCH for %d(%s): %d(%s)[FLIPPED]\n", tileId, dirToName(dir), tId, dirToName(pos))
				matched = true
				id = tId
				borderPos = pos
				reversed = true
				break
			} else {
				log.Debugf("  >%d.%6s  does not match %d.%6s - %s != %s\n", tId, dirToName(pos), tileId, dirToName(dir), b, border)
				log.Debugf("  >%d.%6s^ does not match %d.%6s - %s != %s\n", tId, dirToName(pos), tileId, dirToName(dir), day20_strRev(b), border)
			}
		}
		if matched {
			break
		}
	}

	if matched {
		if dir == borderPos {
			rotation = 180
		} else if dir == (borderPos+1)%4 {
			rotation = 270
		} else if dir == (borderPos+2)%4 {
			rotation = 0
		} else if dir == (borderPos+3)%4 {
			rotation = 90
		}

		log.Debugf("DIR: %d, BORDER_POS: %d, REVERSED: %v", dir, borderPos, reversed)
		flip = reversed
		if (dir < 2 && borderPos < 2) || (dir >= 2 && borderPos >= 2) {
			flip = !flip
			log.Debugf(", FLIP: %v", flip)
		}
		log.Debugln()
		return id, rotation, flip, true
	}

	return 0, 0, false, false
}

func day20_parse(input []string) map[int]*day20_tile {
	tiles := make(map[int]*day20_tile)

	var top strings.Builder
	var right strings.Builder
	var bottom strings.Builder
	var left strings.Builder
	var i int
	for i < len(input) {
		tile := day20_tile{
			Pixels:    make([][]rune, 0, pixelsPerRow),
			Borders:   make([]string, 4),
			Neighbors: make([]*day20_tile, 4),
		}
		fmt.Sscanf(input[i], "Tile %d:", &tile.Id)
		i++
		for ; i < len(input) && input[i] != ""; i++ {
			row := make([]rune, 0, pixelsPerRow)
			for _, p := range input[i] {
				row = append(row, rune(p))
			}
			tile.Pixels = append(tile.Pixels, row)
		}
		i++

		for j := 0; j < pixelsPerRow; j++ {
			top.WriteRune(tile.Pixels[0][j])
			bottom.WriteRune(tile.Pixels[pixelsPerRow-1][j])
			left.WriteRune(tile.Pixels[j][0])
			right.WriteRune(tile.Pixels[j][pixelsPerRow-1])
		}
		tile.Borders[day20_TOP] = top.String()
		tile.Borders[day20_RIGHT] = right.String()
		tile.Borders[day20_BOTTOM] = bottom.String()
		tile.Borders[day20_LEFT] = left.String()

		tiles[tile.Id] = &tile

		top.Reset()
		right.Reset()
		bottom.Reset()
		left.Reset()
	}
	return tiles
}

type day20_tile struct {
	Id                int
	Pixels            [][]rune
	Borders           []string      /* 0 = TOP, 1 = RIGHT, 2 = BOTTOM, 3 = LEFT */
	Neighbors         []*day20_tile /* 0 = TOP, 1 = RIGHT, 2 = BOTTOM, 3 = LEFT */
	OrientationLocked bool
}

func (t day20_tile) String() string {
	var b strings.Builder
	for _, row := range t.Pixels {
		for _, p := range row {
			b.WriteRune(p)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (t *day20_tile) Flip(horizontal bool) {
	t.Pixels = day20_flipPixels(t.Pixels, horizontal)
	if horizontal {
		t.Borders = []string{
			day20_strRev(t.Borders[day20_TOP]),    // TOP
			t.Borders[day20_LEFT],                 // RIGHT
			day20_strRev(t.Borders[day20_BOTTOM]), // BOTTOM
			t.Borders[day20_RIGHT],                // LEFT
		}
	} else {
		t.Borders = []string{
			t.Borders[day20_BOTTOM],              // TOP
			day20_strRev(t.Borders[day20_RIGHT]), // RIGHT
			t.Borders[day20_TOP],                 // BOTTOM
			day20_strRev(t.Borders[day20_LEFT]),  // LEFT
		}
	}
}

func (t *day20_tile) Rotate(degrees int) {
	t.Pixels = day20_rotatePixels(t.Pixels, degrees)
	switch degrees {
	case 90:
		t.Borders = []string{
			day20_strRev(t.Borders[day20_LEFT]),  // TOP
			t.Borders[day20_TOP],                 // RIGHT
			day20_strRev(t.Borders[day20_RIGHT]), // BOTTOM
			t.Borders[day20_BOTTOM],              // LEFT
		}
	case 180:
		t.Borders = []string{
			day20_strRev(t.Borders[day20_BOTTOM]), // TOP
			day20_strRev(t.Borders[day20_LEFT]),   // RIGHT
			day20_strRev(t.Borders[day20_TOP]),    // BOTTOM
			day20_strRev(t.Borders[day20_RIGHT]),  // BOTTOM
		}
	case 270:
		t.Borders = []string{
			t.Borders[day20_RIGHT],                // TOP
			day20_strRev(t.Borders[day20_BOTTOM]), // RIGHT
			t.Borders[day20_LEFT],                 // BOTTOM
			day20_strRev(t.Borders[day20_TOP]),    // LEFT
		}
	}
}

func day20_rotatePixels(pixels [][]rune, degrees int) [][]rune {
	switch degrees {
	case 90:
		newPixels := make([][]rune, len(pixels))
		for i := 0; i < len(pixels); i++ {
			newPixels[i] = make([]rune, len(pixels))
			for j := 0; j < len(pixels); j++ {
				newPixels[i][j] = pixels[len(pixels)-j-1][i]
			}
		}
		pixels = newPixels
	case 180:
		for i := 0; i < len(pixels); i++ {
			for j := 0; j < len(pixels)/2; j++ {
				pixels[i][j], pixels[i][len(pixels)-j-1] = pixels[i][len(pixels)-j-1], pixels[i][j]
			}
		}
		for i := 0; i < len(pixels)/2; i++ {
			pixels[i], pixels[len(pixels)-i-1] = pixels[len(pixels)-i-1], pixels[i]
		}
	case 270:
		newPixels := make([][]rune, len(pixels))
		for i := 0; i < len(pixels); i++ {
			newPixels[i] = make([]rune, len(pixels))
			for j := 0; j < len(pixels); j++ {
				newPixels[i][j] = pixels[j][len(pixels)-i-1]
			}
		}
		pixels = newPixels
	}
	return pixels
}

func day20_flipPixels(pixels [][]rune, horizontal bool) [][]rune {
	if horizontal {
		for i := 0; i < len(pixels); i++ {
			for j := 0; j < len(pixels)/2; j++ {
				pixels[i][j], pixels[i][len(pixels)-j-1] = pixels[i][len(pixels)-j-1], pixels[i][j]
			}
		}
	} else {
		for i := 0; i < len(pixels)/2; i++ {
			pixels[i], pixels[len(pixels)-i-1] = pixels[len(pixels)-i-1], pixels[i]
		}
	}
	return pixels
}

func day20_strRev(str string) string {
	r := []rune(str)
	for i := 0; i < len(r)/2; i++ {
		r[i], r[len(r)-i-1] = r[len(r)-i-1], r[i]
	}
	return string(r)
}

func day20_bordersToString(b []string) string {
	var s strings.Builder
	s.WriteString(b[day20_TOP])
	s.WriteString("\n")
	for i := 1; i < pixelsPerRow-1; i++ {
		s.WriteByte(b[day20_LEFT][i])
		s.WriteString(fmt.Sprintf("%*s", pixelsPerRow-2, ""))
		s.WriteByte(b[day20_RIGHT][i])
		s.WriteString("\n")
	}
	for i := 0; i < pixelsPerRow; i++ {
		s.WriteByte(b[day20_BOTTOM][i])
	}
	return s.String()
}

func init() {
	registerChallengeFunc(20, 1, "day20.txt", day20_part1)
	registerChallengeFunc(20, 2, "day20.txt", day20_part2)
}

func dirToName(dir int) string {
	return ([]string{"TOP", "RIGHT", "BOTTOM", "LEFT"})[dir]
}

const day20_TOP = 0
const day20_RIGHT = 1
const day20_BOTTOM = 2
const day20_LEFT = 3
