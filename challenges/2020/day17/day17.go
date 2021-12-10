package aoc2020_day17

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

//  Z  Y  X

func part1(input []string) (string, error) {
	zSpace := parseInput(input)

	log.Debugln("Before any cycles\n")
	print3DSpace(zSpace)
	log.Debugln("---------------------\n")

	nActive := 0
	for cycle := 0; cycle <= 5; cycle++ {

		zSize := len(zSpace) + 2
		size := len(zSpace[0]) + 2
		log.Debugf("zSize: %d\n", zSize)
		log.Debugf("size: %d\n", size)
		// FIRST: expand the grid!
		newZSpace := make([][][]value, zSize)
		newZSpace[0] = make([][]value, size)
		newZSpace[zSize-1] = make([][]value, size)
		for i := 0; i < size; i++ {
			newZSpace[0][i] = make([]value, size)
			newZSpace[zSize-1][i] = make([]value, size)
		}
		for oldZ, ySpace := range zSpace {
			z := oldZ + 1
			newZSpace[z] = make([][]value, size)
			newZSpace[z][0] = make([]value, size)
			newZSpace[z][size-1] = make([]value, size)
			for oldY, xSpace := range ySpace {
				y := oldY + 1
				newZSpace[z][y] = make([]value, size)
				log.Debugf("zlen: %d, ylen: %d, xlen: %d\n", len(newZSpace), len(newZSpace[z]), len(newZSpace[z][y]))
				for oldX, active := range xSpace {
					x := oldX + 1
					log.Debugf("x: %d, y: %d, z: %d\n", x, y, z)
					newZSpace[z][y][x] = active
				}
			}
		}
		// NEXT: clone the grid
		zSpace = make([][][]value, zSize)
		for z, ySpace := range newZSpace {
			zSpace[z] = make([][]value, size)
			for y, xSpace := range ySpace {
				zSpace[z][y] = make([]value, size)
				for x, active := range xSpace {
					zSpace[z][y][x] = active
				}
			}
		}

		// log.Debugf("AFTER EXPANDING: Cycle #%d:\n\n", cycle+1)
		// print3DSpace(zSpace)

		// NEXT: evaluate everything

		nActive = 0
		for z, ySpace := range zSpace {
			for y, xSpace := range ySpace {
				for x, active := range xSpace {
					activeNeighbors := 0
					for nz := z - 1; nz <= z+1; nz++ {
						if activeNeighbors > 3 {
							break
						}
						for ny := y - 1; ny <= y+1; ny++ {
							if activeNeighbors > 3 {
								break
							}
							for nx := x - 1; nx <= x+1; nx++ {
								if activeNeighbors > 3 {
									break
								}
								if nz >= 0 && ny >= 0 && nx >= 0 && nz < zSize && ny < size && nx < size && !(nz == z && ny == y && nx == x) {
									if zSpace[nz][ny][nx] {
										activeNeighbors++
									}
								}
							}
						}
					}
					if (active && activeNeighbors == 2) || activeNeighbors == 3 {
						newZSpace[z][y][x] = true
						nActive++
					} else {
						newZSpace[z][y][x] = false
					}
				}
			}
		}

		log.Debugf("After Cycle #%d:\n", cycle+1)
		log.Debugf("(# of active cubes: %d)\n", nActive)
		print3DSpace(newZSpace)
		log.Debugln("---------------------\n")
		// print3DSpace(zSpace)
		// log.Debugln("---------------------\n")
		zSpace = newZSpace
	}

	// log.Debugln("Foobar")

	return fmt.Sprintf("%d", nActive), nil
}

func part2(input []string) (string, error) {
	wSpace := make([][][][]value, 0, 1)
	wSpace = append(wSpace, parseInput(input))

	log.Debugln("Before any cycles\n")
	print4DSpace(wSpace)
	log.Debugln("---------------------\n")

	nActive := 0
	for cycle := 0; cycle <= 5; cycle++ {

		wSize := len(wSpace) + 2
		size := len(wSpace[0][0]) + 2
		// log.Debugf("wSize: %d\n", wSize)
		// log.Debugf("size: %d\n", size)

		// FIRST: expand the grid!
		newWSpace := make([][][][]value, wSize)
		newWSpace[0] = make([][][]value, wSize)
		newWSpace[wSize-1] = make([][][]value, wSize)
		for i := 0; i < wSize; i++ {
			newWSpace[0][i] = make([][]value, size)
			newWSpace[wSize-1][i] = make([][]value, size)
		}
		for oldW, zSpace := range wSpace {
			w := oldW + 1
			newWSpace[w] = make([][][]value, wSize)
			newWSpace[w][0] = make([][]value, size)
			newWSpace[w][wSize-1] = make([][]value, size)
			for i := 0; i < wSize; i++ {
				newWSpace[0][i] = make([][]value, size)
				newWSpace[w][i] = make([][]value, size)
				newWSpace[wSize-1][i] = make([][]value, size)
				for j := 0; j < size; j++ {
					newWSpace[0][i][j] = make([]value, size)
					newWSpace[w][i][j] = make([]value, size)
					newWSpace[wSize-1][i][j] = make([]value, size)
				}
			}
			for oldZ, ySpace := range zSpace {
				z := oldZ + 1
				newWSpace[w][z] = make([][]value, size)
				newWSpace[w][z][0] = make([]value, size)
				newWSpace[w][z][size-1] = make([]value, size)
				for oldY, xSpace := range ySpace {
					y := oldY + 1
					_ = y
					newWSpace[w][z][y] = make([]value, size)
					// log.Debugf("wlen: %d, zlen: %d, ylen: %d, xlen: %d\n", len(newWSpace), len(newWSpace[w]), len(newWSpace[w][z]), len(newWSpace[w][z][y]))
					for oldX, active := range xSpace {
						x := oldX + 1
						// log.Debugf("x: %d, y: %d, z: %d\n", x, y, z)
						newWSpace[w][z][y][x] = active
					}
				}
			}
		}
		// NEXT: clone the grid
		wSpace = make([][][][]value, wSize)
		for w, zSpace := range newWSpace {
			wSpace[w] = make([][][]value, wSize)
			for z, ySpace := range zSpace {
				// log.Debugf("w/wlen/newwlen: %d/%d/%d, z/zlen/newzlen: %d/%d/%d\n", w, len(wSpace), len(newWSpace), z, len(wSpace[w]), len(newWSpace[w]))
				wSpace[w][z] = make([][]value, size)
				for y, xSpace := range ySpace {
					wSpace[w][z][y] = make([]value, size)
					for x, active := range xSpace {
						wSpace[w][z][y][x] = active
					}
				}
			}
		}

		// log.Debugf("AFTER EXPANDING: Cycle #%d:\n\n", cycle+1)
		// print4DSpace(newWSpace)

		// NEXT: evaluate everything
		nActive = 0
		for w, zSpace := range wSpace {
			for z, ySpace := range zSpace {
				for y, xSpace := range ySpace {
					for x, active := range xSpace {
						activeNeighbors := 0
						for nw := w - 1; nw <= w+1; nw++ {
							if activeNeighbors > 3 {
								break
							}
							for nz := z - 1; nz <= z+1; nz++ {
								if activeNeighbors > 3 {
									break
								}
								for ny := y - 1; ny <= y+1; ny++ {
									if activeNeighbors > 3 {
										break
									}
									for nx := x - 1; nx <= x+1; nx++ {
										if activeNeighbors > 3 {
											break
										}
										if nw >= 0 && nz >= 0 && ny >= 0 && nx >= 0 && nw < wSize && nz < wSize && ny < size && nx < size && !(nw == w && nz == z && ny == y && nx == x) {
											if wSpace[nw][nz][ny][nx] {
												activeNeighbors++
											}
										}
									}
								}
							}
						}
						if (active && activeNeighbors == 2) || activeNeighbors == 3 {
							newWSpace[w][z][y][x] = true
							nActive++
						} else {
							// log.Debugf("w/wlen/newwlen: %d/%d/%d, z/zlen/newzlen: %d/%d/%d, y/ylen/newylen: %d/%d/%d, x/xlen/newxlen: %d/%d/%d\n", w, len(wSpace), len(newWSpace), z, len(wSpace[w]), len(newWSpace[w]), y, len(wSpace[w][z]), len(newWSpace[w][z]), x, len(wSpace[w][z][y]), len(newWSpace[w][z][y]))
							newWSpace[w][z][y][x] = false
						}
					}
				}
			}
		}

		log.Debugf("After Cycle #%d:\n", cycle+1)
		log.Debugf("(# of active cubes: %d)\n", nActive)
		print4DSpace(newWSpace)
		log.Debugln("---------------------\n")

		wSpace = newWSpace
	}

	return fmt.Sprintf("%d", nActive), nil
}

type value bool

func (v value) String() string {
	if v {
		return "#"
	} else {
		return "."
	}
}

func print3DSpace(zSpace [][][]value) {
	for z, ySpace := range zSpace {
		log.Debugf("z=%d\n", z)
		for _, xSpace := range ySpace {
			log.Debug("  ")
			for _, active := range xSpace {
				log.Debugf("%s", active)
			}
			log.Debugln()
		}
		log.Debugln()
	}
}

func print4DSpace(wSpace [][][][]value) {
	for w, zSpace := range wSpace {
		for z, ySpace := range zSpace {
			log.Debugf("z=%d, w=%d\n", z, w)
			for _, xSpace := range ySpace {
				log.Debug("  ")
				for _, active := range xSpace {
					log.Debugf("%s", active)
				}
				log.Debugln()
			}
			log.Debugln()
		}
	}
}

func parseInput(input []string) [][][]value {
	y := make([][]value, 0, len(input))
	for _, line := range input {
		x := make([]value, 0, len(line))
		for _, c := range line {
			if c == '#' {
				x = append(x, true)
			} else {
				x = append(x, false)
			}
		}
		y = append(y, x)
	}
	return [][][]value{y}
}

func init() {
	challenges.RegisterChallengeFunc(2020, 17, 1, "day17.txt", part1)
	challenges.RegisterChallengeFunc(2020, 17, 2, "day17.txt", part2)
}
