package challenges2021

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func day02_part1(input []string) (string, error) {
	vectors := day02_parse(input)

	depth, xPos := 0, 0

	for _, v := range vectors {
		switch v.direction {
		case "forward":
			xPos += v.magnitude
		case "down":
			depth += v.magnitude
		case "up":
			depth -= v.magnitude
		}
	}

	result := depth * xPos

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day02_part2(input []string) (string, error) {
	vectors := day02_parse(input)
	depth, xPos, aim := 0, 0, 0

	for _, v := range vectors {
		switch v.direction {
		case "forward":
			xPos += v.magnitude
			depth += aim * v.magnitude
		case "down":
			aim += v.magnitude
		case "up":
			aim -= v.magnitude
		}
	}

	result := depth * xPos

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day02_parse(input []string) []*day02_vector {
	vectors := make([]*day02_vector, 0)
	for _, line := range input {
		parts := strings.SplitN(line, " ", 2)
		magnitude, _ := strconv.Atoi(parts[1])
		vectors = append(vectors, &day02_vector{
			direction: parts[0],
			magnitude: magnitude,
		})
	}
	return vectors
}

type day02_vector struct {
	direction string
	magnitude int
}

func init() {
	challenges.RegisterChallengeFunc(2021, 2, 1, "day02.txt", day02_part1)
	challenges.RegisterChallengeFunc(2021, 2, 2, "day02.txt", day02_part2)
}
