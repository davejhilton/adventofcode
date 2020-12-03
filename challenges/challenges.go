package challenges

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type Challenge struct {
	day           int
	part          int
	inputFileName string
	exec          challengeFunc
}

func GetChallenge(day, part int) (Challenge, error) {
	d, ok := allChallenges[day]
	if !ok || d == nil {
		return Challenge{}, fmt.Errorf("Unknown Day: '%d'", day)
	}

	challenge, ok := allChallenges[day][part]
	if !ok {
		return Challenge{}, fmt.Errorf("Unknown Challenge: Day '%d' Part '%d'", day, part)
	}

	return challenge, nil
}

func registerChallengeFunc(day int, part int, inputFileName string, execFunc challengeFunc) {
	if allChallenges == nil {
		allChallenges = make(map[int]map[int]Challenge)
	}
	_, ok := allChallenges[day]
	if !ok {
		allChallenges[day] = make(map[int]Challenge)
	}
	allChallenges[day][part] = Challenge{
		day:           day,
		part:          part,
		inputFileName: inputFileName,
		exec:          execFunc,
	}
}

type challengeFunc func(input []string) (string, error)

func (c Challenge) Run() (string, error) {
	if c.exec == nil {
		return "", fmt.Errorf("Unknown Challenge")
	}
	var input []string
	var err error
	if c.inputFileName != "" {
		input, err = readInputFile(c.inputFileName)
		if err != nil {
			return "", err
		}
	}
	return c.exec(input)
}

func (c Challenge) Name() string {
	if c.day == 0 {
		return "Unknown Challenge"
	}
	return fmt.Sprintf("Day %d - Part %d", c.day, c.part)
}

var allChallenges map[int]map[int]Challenge

func readInputFile(inputFileName string) ([]string, error) {
	pwd, _ := os.Getwd()
	fp := filepath.Join(pwd, "inputs", inputFileName)
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	input := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return input, nil
}
