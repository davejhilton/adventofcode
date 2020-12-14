package challenges

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Challenge struct {
	Day           int
	Part          int
	InputFileName string
	Input         []string
	Exec          challengeFunc
}

func GetChallenge(day, part, exampleNumber int) (Challenge, error) {
	d, ok := allChallenges[day]
	if !ok || d == nil {
		return Challenge{}, fmt.Errorf("Unknown Day: '%d'", day)
	}

	challenge, ok := allChallenges[day][part]
	if !ok {
		return Challenge{}, fmt.Errorf("Unknown Challenge: Day '%d' Part '%d'", day, part)
	}

	if exampleNumber != -1 {
		prefix := strings.TrimSuffix(challenge.InputFileName, ".txt")
		challenge.InputFileName = fmt.Sprintf("%s_example%d.txt", prefix, exampleNumber)
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
		Day:           day,
		Part:          part,
		InputFileName: inputFileName,
		Exec:          execFunc,
	}
}

type challengeFunc func(input []string) (string, error)

func (ch Challenge) Run() (string, error) {
	if ch.Exec == nil {
		return "", fmt.Errorf("Unknown Challenge")
	}
	var input []string
	var err error
	if ch.InputFileName != "" {
		input, err = readInputFile(ch.InputFileName)
		if err != nil {
			return "", err
		}
		ch.Input = input
	}
	currentChallenge = &ch
	defer func() {
		currentChallenge = nil
	}()
	return ch.Exec(input)
}

func (ch Challenge) Name() string {
	if ch.Day == 0 {
		return "Unknown Challenge"
	}
	return fmt.Sprintf("Day %d - Part %d", ch.Day, ch.Part)
}

var currentChallenge *Challenge

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
