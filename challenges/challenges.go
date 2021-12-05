package challenges

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Challenge struct {
	Year          int
	Day           int
	Part          int
	InputFileName string
	Input         []string
	Exec          challengeFunc
}

func GetChallenge(year, day, part, exampleNumber int) (Challenge, error) {
	y, ok := allChallenges[year]
	if !ok || y == nil {
		return Challenge{}, fmt.Errorf("unknown year: '%d'", year)
	}

	d, ok := allChallenges[year][day]
	if !ok || d == nil {
		return Challenge{}, fmt.Errorf("unknown day: Year '%d' Day '%d'", year, day)
	}

	challenge, ok := allChallenges[year][day][part]
	if !ok {
		return Challenge{}, fmt.Errorf("unknown challenge: Year '%d' Day '%d' Part '%d'", year, day, part)
	}

	if exampleNumber != -1 {
		prefix := strings.TrimSuffix(challenge.InputFileName, ".txt")
		challenge.InputFileName = fmt.Sprintf("%s_example%d.txt", prefix, exampleNumber)
	}

	return challenge, nil
}

func RegisterChallengeFunc(year int, day int, part int, inputFileName string, execFunc challengeFunc) {
	if allChallenges == nil {
		allChallenges = make(map[int]map[int]map[int]Challenge)
	}
	_, ok := allChallenges[year]
	if !ok {
		allChallenges[year] = make(map[int]map[int]Challenge)
	}
	_, ok = allChallenges[year][day]
	if !ok {
		allChallenges[year][day] = make(map[int]Challenge)
	}
	allChallenges[year][day][part] = Challenge{
		Year:          year,
		Day:           day,
		Part:          part,
		InputFileName: inputFileName,
		Exec:          execFunc,
	}
}

type challengeFunc func(input []string) (string, error)

func (ch Challenge) Run() (string, error) {
	if ch.Exec == nil {
		return "", fmt.Errorf("unknown challenge")
	}
	var input []string
	var err error
	if ch.InputFileName != "" {
		input, err = readInputFile(ch.Year, ch.InputFileName)
		if err != nil {
			return "", err
		}
		ch.Input = input
	}
	CurrentChallenge = &ch
	defer func() {
		CurrentChallenge = nil
	}()
	return ch.Exec(input)
}

func (ch Challenge) Name() string {
	if ch.Year == 0 {
		return "Unknown Year"
	}
	if ch.Day == 0 {
		return "Unknown Challenge"
	}
	return fmt.Sprintf("Year %d - Day %d - Part %d", ch.Year, ch.Day, ch.Part)
}

var CurrentChallenge *Challenge

var allChallenges map[int]map[int]map[int]Challenge

func readInputFile(year int, inputFileName string) ([]string, error) {
	pwd, _ := os.Getwd()
	fp := filepath.Join(pwd, "inputs", fmt.Sprintf("%d", year), inputFileName)
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
