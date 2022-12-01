package aoc2021_day19

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func printScannerList(msg string, scanners []*ScannerField) {
	log.Debugf("%s [", msg)
	for ii, o := range scanners {
		if ii > 0 {
			log.Debug(", ")
		}
		log.Debugf("%d", o.ScannerNumber)
	}
	log.Debugln("]")
}

func part1(input []string) (string, error) {
	allScanners := parseInput(input)
	mappedScanners := BuildMap(allScanners)

	uniqueBeaconCoords := make(map[string]bool)
	for _, s := range mappedScanners {
		for _, b := range s.Beacons {
			coordStr := fmt.Sprintf("%d,%d,%d", b["x"], b["y"], b["z"])
			uniqueBeaconCoords[coordStr] = true
		}
	}

	allBeacons := make([]string, 0)
	for c := range uniqueBeaconCoords {
		allBeacons = append(allBeacons, c)
	}

	sort.Slice(allBeacons, func(i, j int) bool {
		iParts := strings.Split(allBeacons[i], ",")
		jParts := strings.Split(allBeacons[j], ",")
		i0, j0 := util.Atoi(iParts[0]), util.Atoi(jParts[0])
		if i0 < j0 {
			return true
		} else if i0 > j0 {
			return false
		} else {
			i1, j1 := util.Atoi(iParts[1]), util.Atoi(jParts[1])
			if i1 < j1 {
				return true
			} else if i1 > j1 {
				return false
			} else {
				i2, j2 := util.Atoi(iParts[2]), util.Atoi(jParts[2])
				if i2 < j2 {
					return true
				} else if i2 > j2 {
					return false
				} else {
					return i < j
				}
			}
		}
	})
	// sort.Strings(allBeacons)
	log.Debugln("Final Results\n-------------------------")
	for _, c := range allBeacons {
		log.Debugln(c)
	}

	result := len(uniqueBeaconCoords)
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	allScanners := parseInput(input)
	mappedScanners := BuildMap(allScanners)

	maxDist := 0
	for i, s1 := range mappedScanners {
		for j := i + 1; j < len(mappedScanners); j++ {
			s2 := mappedScanners[j]
			maxDist = util.Max(maxDist, calculateManhattanDistance(s1, s2))
		}
	}

	return fmt.Sprintf("%d", maxDist), nil
}

func parseInput(input []string) []*ScannerField {
	var scannerRegex = regexp.MustCompile(`--- scanner (\d+) ---`)
	var scanners = make([]*ScannerField, 0)

	var scan *ScannerField
	for _, s := range input {
		if s == "" {
			scan.Distances = CalculateDistances(scan.Beacons)
			scanners = append(scanners, scan)
			continue
		}
		if scannerRegex.MatchString(s) {
			var num int
			fmt.Sscanf(s, "--- scanner %d ---", &num)
			log.Debugf("Parsing scanner line: %s -> %d\n", s, num)
			scan = &ScannerField{
				ScannerNumber: num,
				Beacons:       make([]Coordinate, 0),
			}
			continue
		}
		var x, y, z int
		fmt.Sscanf(s, "%d,%d,%d", &x, &y, &z)
		c := Coordinate{"x": x, "y": y, "z": z}
		scan.Beacons = append(scan.Beacons, c)
	}
	scan.Distances = CalculateDistances(scan.Beacons)
	scanners = append(scanners, scan)
	return scanners
}

func BuildMap(allScanners []*ScannerField) []*ScannerField {
	for _, sc := range allScanners {
		log.Debugf("Scanner %d (%d beacons)", sc.ScannerNumber, len(sc.Beacons))
		for j, b := range sc.Beacons {
			log.Debugf("\n\t%2d - (%4d, %4d, %4d)", j, b["x"], b["y"], b["z"])
		}
		log.Debug("\n\n")
	}

	knownScanners := []*ScannerField{allScanners[0]}
	unknownScanners := allScanners[1:]

	printScannerList("Known scanners:  ", knownScanners)
	printScannerList("Unknown scanners:", unknownScanners)

	for len(unknownScanners) > 0 {
		fmt.Printf("%d more 'Unknown' Scanners...\n", len(unknownScanners))
		for i := 0; i < len(unknownScanners); i++ {
			unknown := unknownScanners[i]
			match, rIdx, matchedScanner, matchedDistances := FindRotationMatch(unknown, knownScanners)
			fmt.Printf("Scanner %d match: %v\n", unknown.ScannerNumber, match)
			if match {
				log.Debugf("Scanner %d matches scanner %d at rotation: %d\n", unknown.ScannerNumber, matchedScanner.ScannerNumber, rIdx)
				// fix the rotation,
				unknown.ApplyRotation(rIdx)
				// update the relative positions
				unknown.UpdateToAbsolutePosition(matchedScanner, matchedDistances)

				log.Debugf("Moving scanner %d to the known list:\n", unknown.ScannerNumber)
				printScannerList("Known scanners before:", knownScanners)
				// add this scanner to "known" list
				knownScanners = append(knownScanners[0:], unknown)
				printScannerList("Known scanners after: ", knownScanners)
				// delete this scanner from "unknown" list
				printScannerList("Unknown scanners before:", unknownScanners)
				unknownScanners = append(unknownScanners[:i], unknownScanners[i+1:]...)
				printScannerList("Unknown scanners after: ", unknownScanners)
				break
			}
		}
	}

	return knownScanners
}

func calculateManhattanDistance(s1, s2 *ScannerField) int {
	xDist := util.Abs(s1.ScannerCoords["x"] - s2.ScannerCoords["x"])
	yDist := util.Abs(s1.ScannerCoords["y"] - s2.ScannerCoords["y"])
	zDist := util.Abs(s1.ScannerCoords["z"] - s2.ScannerCoords["z"])
	return xDist + yDist + zDist
}

type Coordinate map[string]int

type ScannerField struct {
	ScannerNumber int
	Beacons       []Coordinate
	Distances     map[string]DistanceSet
	ScannerCoords Coordinate
}

func (s *ScannerField) GetRotatedCoords(rIdx int) []Coordinate {
	r := rotations[rIdx]
	coords := make([]Coordinate, 0, len(s.Beacons))
	for _, b := range s.Beacons {
		log.Debugf("ROTATING(%d: %d): (%4d, %4d, %4d)", s.ScannerNumber, rIdx, b["x"], b["y"], b["z"])
		c := Coordinate{
			"x": b[r.Props[0]] * r.Mult[0],
			"y": b[r.Props[1]] * r.Mult[1],
			"z": b[r.Props[2]] * r.Mult[2],
		}
		log.Debugf(" -> (%4d, %4d, %4d)\n\n", c["x"], c["y"], c["z"])
		coords = append(coords, c)
	}
	return coords
}

func (s *ScannerField) ApplyRotation(rIdx int) {
	if rIdx != 0 {
		if cache, ok := rotationCache[s.ScannerNumber][rIdx]; ok {
			s.Beacons = cache.C
			s.Distances = cache.D
		} else {
			s.Beacons = s.GetRotatedCoords(rIdx)
			s.Distances = CalculateDistances(s.Beacons)
		}
	}
}

func CalculateDistances(coords []Coordinate) map[string]DistanceSet {
	dist := make(map[string]DistanceSet)
	log.Debugf("Calculating distances - length: %d\n", len(coords))
	for i, c1 := range coords {
		for j := i + 1; j < len(coords); j++ {
			c2 := coords[j]
			key := fmt.Sprintf("%d,%d", i, j)
			dist[key] = NewDistanceSet(c1, c2)
		}
	}

	return dist
}

func (s *ScannerField) UpdateToAbsolutePosition(relativeTo *ScannerField, distMap map[string]string) {

	xAdjust, yAdjust, zAdjust := 0, 0, 0

	i := 0
	for myKey, theirKey := range distMap {
		log.Debugf("MAP: %s -> %s\n", myKey, theirKey)
		var i1m, i2m, i1t, i2t int
		fmt.Sscanf(myKey, "%d,%d", &i1m, &i2m)
		fmt.Sscanf(theirKey, "%d,%d", &i1t, &i2t)

		my1 := s.Beacons[i1m]
		my2 := s.Beacons[i2m]
		their1 := relativeTo.Beacons[i1t]
		their2 := relativeTo.Beacons[i2t]

		log.Debugf("\tMINE:   (%4d, %4d, %4d) -> (%4d, %4d, %4d)\n",
			my1["x"], my1["y"], my1["z"], my2["x"], my2["y"], my2["z"],
		)
		log.Debugf("\tTHEIRS: (%4d, %4d, %4d) -> (%4d, %4d, %4d)\n\n\n",
			their1["x"], their1["y"], their1["z"], their2["x"], their2["y"], their2["z"],
		)

		theirXDiff := their1["x"] - their2["x"]
		theirYDiff := their1["y"] - their2["y"]
		theirZDiff := their1["z"] - their2["z"]
		myXDiff := my1["x"] - my2["x"]
		myYDiff := my1["y"] - my2["y"]
		myZDiff := my1["z"] - my2["z"]

		if myXDiff != theirXDiff || myYDiff != theirYDiff || myZDiff != theirZDiff {
			fmt.Printf("\tx: %d vs %d\n\ty: %d vs %d\n\tz: %d vs %d\n\n",
				myXDiff, theirXDiff,
				myYDiff, theirYDiff,
				myZDiff, theirZDiff,
			)
			my1, my2 = my2, my1
		}

		myXDiff = my1["x"] - my2["x"]
		myYDiff = my1["y"] - my2["y"]
		myZDiff = my1["z"] - my2["z"]
		if myXDiff != theirXDiff || myYDiff != theirYDiff || myZDiff != theirZDiff {
			fmt.Printf("WELL CRAP. THEY AREN'T RIGHT...\nScanner %d doesn't map to Scanner %d\n", s.ScannerNumber, relativeTo.ScannerNumber)
			fmt.Printf("\tx: %d vs %d\n\ty: %d vs %d\n\tz: %d vs %d\n\n",
				myXDiff, theirXDiff,
				myYDiff, theirYDiff,
				myZDiff, theirZDiff,
			)
			os.Exit(1)
		}

		xAdj := their1["x"] - my1["x"]
		yAdj := their1["y"] - my1["y"]
		zAdj := their1["z"] - my1["z"]

		if i != 0 {
			if xAdj != xAdjust || yAdj != yAdjust || zAdj != zAdjust {
				fmt.Printf("WELL CRAP. GETTING DIFFERENT ADJUSTMENT VALUES FOR Scanner %d\n", s.ScannerNumber)
				fmt.Printf("BEFORE - x: %4d, y: %4d, z: %4d\n", xAdjust, yAdjust, zAdjust)
				fmt.Printf("NOW    - x: %4d, y: %4d, z: %4d\n", xAdj, yAdj, zAdj)
				os.Exit(1)
			}
		}
		xAdjust = xAdj
		yAdjust = yAdj
		zAdjust = zAdj
		i++
	}

	log.Debugf("UPDATING BEACONS MAP for Scanner %d\n", s.ScannerNumber)
	for i, b := range s.Beacons {
		log.Debugf(" - beacon %d: (%d,%d,%d)\n", i, b["x"], b["y"], b["z"])
		b["x"] += xAdjust
		b["y"] += yAdjust
		b["z"] += zAdjust
	}
	s.ScannerCoords = Coordinate{
		"x": xAdjust * -1,
		"y": yAdjust * -1,
		"z": zAdjust * -1,
	}
}

func FindRotationMatch(s *ScannerField, known []*ScannerField) (bool, int, *ScannerField, map[string]string) {
	dist := s.Distances
	for rIdx := 0; rIdx <= 23; rIdx++ {
		if rIdx != 0 {
			if cache, ok := rotationCache[s.ScannerNumber][rIdx]; ok {
				dist = cache.D
			} else {
				coords := s.GetRotatedCoords(rIdx)
				dist = CalculateDistances(coords)
				if _, ok2 := rotationCache[s.ScannerNumber]; !ok2 {
					rotationCache[s.ScannerNumber] = make(map[int]struct {
						C []Coordinate
						D map[string]DistanceSet
					})
				}
				rotationCache[s.ScannerNumber][rIdx] = struct {
					C []Coordinate
					D map[string]DistanceSet
				}{
					C: coords,
					D: dist,
				}
			}
		}
		for _, other := range known {
			if other.ScannerNumber == s.ScannerNumber {
				log.Debugf("SOMEHOW, unknown scanner %d is in the known list... [", s.ScannerNumber)
				for ii, o := range known {
					if ii > 0 {
						log.Debug(", ")
					}
					log.Debugf("%d", o.ScannerNumber)
				}
				log.Debugln("]")
				continue
			}
			knownDists := other.Distances

			if matches := FindDistanceOverlaps(dist, knownDists); len(matches) >= 12 {
				return true, rIdx, other, matches
			}
		}
	}

	return false, -1, nil, map[string]string{}
}

func FindDistanceOverlaps(d1, d2 map[string]DistanceSet) map[string]string {
	matches := make(map[string]string)
	for k1, v1 := range d1 {
		for k2, v2 := range d2 {
			if v1 == v2 {
				matches[k1] = k2
			}
		}
	}
	return matches
}

type Rotation struct {
	Props []string
	Mult  []int
}

var rotations = []Rotation{
	{Props: []string{"x", "y", "z"}, Mult: []int{1, 1, 1}},
	{Props: []string{"x", "y", "z"}, Mult: []int{-1, -1, 1}},
	{Props: []string{"x", "y", "z"}, Mult: []int{-1, 1, -1}},
	{Props: []string{"x", "y", "z"}, Mult: []int{1, -1, -1}},
	{Props: []string{"x", "z", "y"}, Mult: []int{-1, -1, -1}},
	{Props: []string{"x", "z", "y"}, Mult: []int{-1, 1, 1}},
	{Props: []string{"x", "z", "y"}, Mult: []int{1, -1, 1}},
	{Props: []string{"x", "z", "y"}, Mult: []int{1, 1, -1}},
	{Props: []string{"y", "x", "z"}, Mult: []int{-1, -1, -1}},
	{Props: []string{"y", "x", "z"}, Mult: []int{-1, 1, 1}},
	{Props: []string{"y", "x", "z"}, Mult: []int{1, -1, 1}},
	{Props: []string{"y", "x", "z"}, Mult: []int{1, 1, -1}},
	{Props: []string{"y", "z", "x"}, Mult: []int{-1, -1, 1}},
	{Props: []string{"y", "z", "x"}, Mult: []int{-1, 1, -1}},
	{Props: []string{"y", "z", "x"}, Mult: []int{1, -1, -1}},
	{Props: []string{"y", "z", "x"}, Mult: []int{1, 1, 1}},
	{Props: []string{"z", "x", "y"}, Mult: []int{-1, -1, 1}},
	{Props: []string{"z", "x", "y"}, Mult: []int{-1, 1, -1}},
	{Props: []string{"z", "x", "y"}, Mult: []int{1, -1, -1}},
	{Props: []string{"z", "x", "y"}, Mult: []int{1, 1, 1}},
	{Props: []string{"z", "y", "x"}, Mult: []int{-1, -1, -1}},
	{Props: []string{"z", "y", "x"}, Mult: []int{-1, 1, 1}},
	{Props: []string{"z", "y", "x"}, Mult: []int{1, -1, 1}},
	{Props: []string{"z", "y", "x"}, Mult: []int{1, 1, -1}},
}

var rotationCache = make(map[int]map[int]struct {
	C []Coordinate
	D map[string]DistanceSet
})

type DistanceSet string

func NewDistanceSet(c1, c2 Coordinate) DistanceSet {
	return DistanceSet(fmt.Sprintf(
		"%d,%d,%d",
		// util.Abs(c1["x"]-c2["x"]),
		// util.Abs(c1["y"]-c2["y"]),
		// util.Abs(c1["z"]-c2["z"]),
		c1["x"]-c2["x"],
		c1["y"]-c2["y"],
		c1["z"]-c2["z"],
	))
}

func init() {
	challenges.RegisterChallengeFunc(2021, 19, 1, "day19.txt", part1)
	challenges.RegisterChallengeFunc(2021, 19, 2, "day19.txt", part2)
}
