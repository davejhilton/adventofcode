package aoc2020_day4

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	passports := parseInput(input)
	nValid := countValid(passports, false)
	return fmt.Sprintf("%d", nValid), nil
}

func part2(input []string) (string, error) {
	passports := parseInput(input)
	nValid := countValid(passports, true)
	return fmt.Sprintf("%d", nValid), nil
}

func parseInput(input []string) []passport {
	passports := make([]passport, 0)

	var cur []string
	for _, line := range input {
		// keep reading lines until we reach the end of a passport
		if line != "" {
			if cur == nil {
				cur = make([]string, 0)
			}
			cur = append(cur, line)
			// then parse that passport, store it, and move on
		} else {
			if cur != nil {
				passports = append(passports, parsePassport(cur))
				cur = nil
			}
		}
	}
	// don't forget the last one, which may not have hit an empty line to actually get stored
	if cur != nil {
		passports = append(passports, parsePassport(cur))
	}

	return passports
}

func countValid(passports []passport, strict bool) int {
	nValid := 0
	for i, p := range passports {
		if p.Valid(strict) {
			nValid++
		}
		log.Debugf("%3d: %s\n", i, p.DebugString(strict))
	}
	return nValid
}

func parsePassport(input []string) passport {
	passport := passport{}
	for _, line := range input {
		kvPairs := strings.Split(line, " ")
		for _, pair := range kvPairs {
			kv := strings.Split(pair, ":")
			k, v := kv[0], kv[1]
			var err error
			switch k {
			case "byr":
				passport.BirthYear, err = strconv.Atoi(v)
				if err != nil {
					fmt.Printf("ERROR: expected int for 'byr', got: '%s'\n", v)
				}
			case "iyr":
				passport.IssueYear, err = strconv.Atoi(v)
				if err != nil {
					fmt.Printf("ERROR: expected int for 'iyr', got: '%s'\n", v)
				}
			case "eyr":
				passport.ExpirationYear, err = strconv.Atoi(v)
				if err != nil {
					fmt.Printf("ERROR: expected int for 'eyr', got: '%s'\n", v)
				}
			case "hgt":
				passport.Height = v
			case "hcl":
				passport.HairColor = v
			case "ecl":
				passport.EyeColor = v
			case "pid":
				passport.PassportID = v
			case "cid":
				passport.CountryID = v
			default:
				fmt.Printf("ERROR: unknown key '%s' (value: '%s')\n", k, v)
			}
		}
	}
	return passport
}

type passport struct {
	BirthYear      int
	IssueYear      int
	ExpirationYear int
	Height         string
	HairColor      string
	EyeColor       string
	PassportID     string
	CountryID      string
}

func (p passport) String() string {
	return fmt.Sprintf("[ byr:%-4d   iyr:%-4d   eyr: %-4d   hgt:%-6s   hcl:%-7s   ecl:%-7s   pid:%-10s   cid:%-3s ]",
		p.BirthYear,
		p.IssueYear,
		p.ExpirationYear,
		p.Height,
		p.HairColor,
		p.EyeColor,
		p.PassportID,
		p.CountryID)
}

func (p passport) DebugString(strict bool) string {
	var label string
	if p.Valid(strict) {
		label = log.Colorize(" VALID ", log.Green, -7)
	} else {
		label = log.Colorize("INVALID", log.Red, -7)
	}

	getColor := func(valid bool) log.Color {
		if !valid {
			return log.Red
		}
		return log.Normal
	}
	handleEmpty := func(val interface{}) string {
		result := ""
		switch v := val.(type) {
		case string:
			result = v
		case int:
			if v == 0 {
				result = ""
			} else {
				result = fmt.Sprintf("%d", v)
			}
		default:
			result = fmt.Sprintf("%v", v)
		}

		if result == "" {
			result = "---"
		}
		return result
	}

	return fmt.Sprintf("%s - [ byr:%s   iyr:%s   eyr: %s   hgt:%s   hcl:%s   ecl:%s   pid:%s   cid:%s ]",
		label,
		log.Colorize(handleEmpty(p.BirthYear), getColor(p.validBirthYear(strict)), -4),
		log.Colorize(handleEmpty(p.IssueYear), getColor(p.validIssueYear(strict)), -4),
		log.Colorize(handleEmpty(p.ExpirationYear), getColor(p.validExpirationYear(strict)), -4),
		log.Colorize(handleEmpty(p.Height), getColor(p.validHeight(strict)), -6),
		log.Colorize(handleEmpty(p.HairColor), getColor(p.validHairColor(strict)), -7),
		log.Colorize(handleEmpty(p.EyeColor), getColor(p.validEyeColor(strict)), -7),
		log.Colorize(handleEmpty(p.PassportID), getColor(p.validPassportID(strict)), -10),
		log.Colorize(handleEmpty(p.CountryID), getColor(p.validCountryID(strict)), -5),
	)
}

func (p passport) Valid(strict bool) bool {
	return p.validBirthYear(strict) &&
		p.validIssueYear(strict) &&
		p.validExpirationYear(strict) &&
		p.validHeight(strict) &&
		p.validHairColor(strict) &&
		p.validEyeColor(strict) &&
		p.validPassportID(strict) &&
		p.validCountryID(strict)
}

func (p passport) validBirthYear(strict bool) bool {
	if strict {
		return p.BirthYear >= 1920 && p.BirthYear <= 2002
	} else {
		return p.BirthYear != 0
	}
}

func (p passport) validIssueYear(strict bool) bool {
	if strict {
		return p.IssueYear >= 2010 && p.IssueYear <= 2020
	} else {
		return p.IssueYear != 0
	}
}

func (p passport) validExpirationYear(strict bool) bool {
	if strict {
		return p.ExpirationYear >= 2020 && p.ExpirationYear <= 2030
	} else {
		return p.ExpirationYear != 0
	}
}

func (p passport) validHeight(strict bool) bool {
	if strict {
		if strings.HasSuffix(p.Height, "in") {
			in, _ := strconv.Atoi(strings.Replace(p.Height, "in", "", 1))
			return in >= 59 && in <= 76
		} else if strings.HasSuffix(p.Height, "cm") {
			cm, _ := strconv.Atoi(strings.Replace(p.Height, "cm", "", 1))
			return cm >= 150 && cm <= 193
		}
		return false
	} else {
		return p.Height != ""
	}
}

func (p passport) validHairColor(strict bool) bool {
	if strict {
		match, _ := regexp.MatchString("^#[0-9a-f]{6}$", p.HairColor)
		return match
	} else {
		return p.HairColor != ""
	}
}

func (p passport) validEyeColor(strict bool) bool {
	if strict {
		match, _ := regexp.MatchString("^(amb|blu|brn|gry|grn|hzl|oth)$", p.EyeColor)
		return match
	} else {
		return p.EyeColor != ""
	}
}

func (p passport) validPassportID(strict bool) bool {
	if strict {
		match, _ := regexp.MatchString("^[0-9]{9}$", p.PassportID)
		return match
	} else {
		return p.PassportID != ""
	}
}

func (p passport) validCountryID(strict bool) bool {
	// no validation
	return true
}

func init() {
	challenges.RegisterChallengeFunc(2020, 4, 1, "day04.txt", part1)
	challenges.RegisterChallengeFunc(2020, 4, 2, "day04.txt", part2)
}
