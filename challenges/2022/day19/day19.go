package aoc2022_day19

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

type State [8]int

func part1(input []string) (string, error) {
	blueprints := parseInput(input)
	// log.Debugln("Parsed Input:")
	// log.DebugJSON(blueprints, true)

	// map of blueprint id -> number of geodes
	bestResults := runMiningOperation(blueprints, 24)

	log.Debugln("Best Results for each blueprint")
	result := 0
	for id, best := range bestResults {
		log.Debugf("Blueprint %d: %d Geodes (quality: %d)\n", id, best, id*best)
		result += id * best
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	blueprints := parseInput(input)
	// log.Debugln("Parsed Input:")
	// log.DebugJSON(blueprints, true)

	if len(blueprints) > 3 {
		blueprints = blueprints[:3]
	}

	bestResults := runMiningOperation(blueprints, 32)

	log.Debugln("Best Results for each blueprint")
	result := 1
	for id, best := range bestResults {
		log.Debugf("Blueprint %d: %d Geodes\n", id, best)
		result *= best
	}

	return fmt.Sprintf("%d", result), nil
}

func runMiningOperation(blueprints []Blueprint, minutes int) map[int]int {
	// map of blueprint id -> number of geodes
	bestResults := make(map[int]int)

	for _, bp := range blueprints {
		startState := State{
			0, // ORE
			0, // CLAY
			0, // OBSIDIAN
			0, // GEODE
			1, // ORE_ROBOT
			0, // CLAY_ROBOT
			0, // OBSIDIAN_ROBOT
			0, // GEODE_ROBOT
		}
		// map of "state" -> empty struct
		var prevStates = map[State]struct{}{
			startState: {},
		}
		var maxGeodes int
		var maxRobots int
		for minute := 0; minute < minutes; minute++ {
			log.Debugf("Starting blueprint %d minute %d (%d states) (%d max geodes)\n", bp.Id, minute, len(prevStates), maxGeodes)
			var newStates = make(map[State]struct{})
			for state := range prevStates {
				// let the robots mine their things
				newResources := [4]int{
					state[ORE] + state[ORE_ROBOT],
					state[CLAY] + state[CLAY_ROBOT],
					state[OBSIDIAN] + state[OBSIDIAN_ROBOT],
					state[GEODE] + state[GEODE_ROBOT],
				}

				if newResources[GEODE] > maxGeodes {
					maxGeodes = newResources[GEODE]
				} else if newResources[GEODE] < maxGeodes && state[GEODE_ROBOT] < maxRobots-1 {
					continue // if it's more than 1 robot behind, it's never going to catch up
				}

				if minute == minutes-1 {
					continue //
				}

				purchases := 0
				// figure out what to buy
				for rType, costs := range bp.Costs {
					if rType == ORE && state[ORE_ROBOT] > 4 {
						continue // don't waste time buying more ore robots than needed
					} else if rType == CLAY && state[CLAY_ROBOT] > 16 {
						continue // don't waste time buying more clay robots than needed
					} else if rType == OBSIDIAN && state[OBSIDIAN_ROBOT] > 18 {
						continue // don't waste time buying more obsidian robots than needed
					}
					hasEnough := true
					for _, c := range costs {
						if state[c.Type] < c.Amount {
							hasEnough = false
							break
						}
					}
					if hasEnough {
						purchases++
						var newState = State{
							newResources[ORE], newResources[CLAY], newResources[OBSIDIAN], newResources[GEODE],
							state[ORE_ROBOT], state[CLAY_ROBOT], state[OBSIDIAN_ROBOT], state[GEODE_ROBOT],
						}
						for _, c := range costs {
							newState[c.Type] = newState[c.Type] - c.Amount // subtract costs
						}
						newState[rType+4] = newState[rType+4] + 1 // add robot
						if _, ok := prevStates[newState]; !ok {
							newStates[newState] = struct{}{}
						}
						if newState[GEODE_ROBOT] > maxRobots {
							maxRobots = newState[GEODE_ROBOT]
						}
					}
				}
				if purchases < 4 {
					newStates[State{
						newResources[ORE], newResources[CLAY], newResources[OBSIDIAN], newResources[GEODE],
						state[ORE_ROBOT], state[CLAY_ROBOT], state[OBSIDIAN_ROBOT], state[GEODE_ROBOT],
					}] = struct{}{}
				}
			}
			prevStates = newStates
		}
		bestResults[bp.Id] = maxGeodes
	}

	return bestResults
}

func parseInput(input []string) []Blueprint {
	blueprints := make([]Blueprint, 0, len(input))
	for _, s := range input {
		var id, oreOre, clayOre, obsOre, obsClay, geoOre, geoObs int
		fmt.Sscanf(
			s,
			"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&id, &oreOre, &clayOre, &obsOre, &obsClay, &geoOre, &geoObs,
		)
		blueprints = append(blueprints, Blueprint{
			Id: id,
			Costs: map[ResourceType][]Cost{
				ORE:      {{ORE, oreOre}},
				CLAY:     {{ORE, clayOre}},
				OBSIDIAN: {{ORE, obsOre}, {CLAY, obsClay}},
				GEODE:    {{ORE, geoOre}, {OBSIDIAN, geoObs}},
			},
		})
	}
	return blueprints
}

type ResourceType int
type RobotType int

const (
	// these are used as indexes in the 'state' arrays
	ORE            ResourceType = 0
	CLAY           ResourceType = 1
	OBSIDIAN       ResourceType = 2
	GEODE          ResourceType = 3
	ORE_ROBOT      RobotType    = 4
	CLAY_ROBOT     RobotType    = 5
	OBSIDIAN_ROBOT RobotType    = 6
	GEODE_ROBOT    RobotType    = 7
)

type Blueprint struct {
	Id    int
	Costs map[ResourceType][]Cost
}

type Cost struct {
	Type   ResourceType
	Amount int
}

type Robot struct {
	Produces ResourceType
	Costs    []*Cost
}

func init() {
	challenges.RegisterChallengeFunc(2022, 19, 1, "day19.txt", part1)
	challenges.RegisterChallengeFunc(2022, 19, 2, "day19.txt", part2)
}
