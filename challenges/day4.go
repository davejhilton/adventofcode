package challenges

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func day4_part1(input []string) (string, error) {
	passports := day4_parsePassports(input)
	nValid := day4_countValid(passports, false)
	return fmt.Sprintf("%d", nValid), nil
}

func day4_part2(input []string) (string, error) {
	passports := day4_parsePassports(input)
	nValid := day4_countValid(passports, true)
	return fmt.Sprintf("%d", nValid), nil
}

func day4_countValid(passports []day4_passport, strict bool) int {
	nValid := 0
	var ok bool
	for _, p := range passports {
		ok = p.Valid(strict)
		if ok {
			nValid++
		}
		fmt.Printf("valid: %-5v --- %s\n", ok, p)
	}
	return nValid
}

func day4_parsePassports(input []string) []day4_passport {
	passports := make([]day4_passport, 0)
	cur := day4_passport{}
	for _, line := range input {
		if line == "" {
			// empty line. last passport is populated. add it and create a new one to populate next
			passports = append(passports, cur)
			cur = day4_passport{}
			continue
		}

		kvPairs := strings.Split(line, " ")
		for _, pair := range kvPairs {
			kv := strings.Split(pair, ":")
			k, v := kv[0], kv[1]
			var err error
			switch k {
			case "byr":
				cur.BirthYear, err = strconv.Atoi(v)
				if err != nil {
					fmt.Printf("ERROR: expected int for 'byr', got: '%s'\n", v)
				}
			case "iyr":
				cur.IssueYear, err = strconv.Atoi(v)
				if err != nil {
					fmt.Printf("ERROR: expected int for 'iyr', got: '%s'\n", v)
				}
			case "eyr":
				cur.ExpirationYear, err = strconv.Atoi(v)
				if err != nil {
					fmt.Printf("ERROR: expected int for 'eyr', got: '%s'\n", v)
				}
			case "hgt":
				cur.Height = v
			case "hcl":
				cur.HairColor = v
			case "ecl":
				cur.EyeColor = v
			case "pid":
				cur.PassportID = v
			case "cid":
				cur.CountryID = v
			default:
				fmt.Printf("ERROR: unknown key '%s' (value: '%s')\n", k, v)
			}
		}
	}
	return passports
}

func init() {
	registerChallengeFunc(4, 1, "day4.txt", day4_part1)
	registerChallengeFunc(4, 2, "day4.txt", day4_part2)
}

type day4_passport struct {
	BirthYear      int
	IssueYear      int
	ExpirationYear int
	Height         string
	HairColor      string
	EyeColor       string
	PassportID     string
	CountryID      string
}

func (p day4_passport) String() string {
	return fmt.Sprintf("[[   byr:%-4d   iyr:%-4d   eyr: %-4d   hgt:%-6s   hcl:%-7s   ecl:%-7s   pid:%-10s   cid:%-3s   ]]",
		p.BirthYear,
		p.IssueYear,
		p.ExpirationYear,
		p.Height,
		p.HairColor,
		p.EyeColor,
		p.PassportID,
		p.CountryID)
}

func (p day4_passport) Valid(strict bool) bool {

	return p.validBirthYear(strict) &&
		p.validIssueYear(strict) &&
		p.validExpirationYear(strict) &&
		p.validHeight(strict) &&
		p.validHairColor(strict) &&
		p.validEyeColor(strict) &&
		p.validPassportID(strict) &&
		p.validCountryID(strict)
}

func (p day4_passport) validBirthYear(strict bool) bool {
	if strict {
		return p.BirthYear >= 1920 && p.BirthYear <= 2002
	} else {
		return p.BirthYear != 0
	}
}

func (p day4_passport) validIssueYear(strict bool) bool {
	if strict {
		return p.IssueYear >= 2010 && p.IssueYear <= 2020
	} else {
		return p.IssueYear != 0
	}
}

func (p day4_passport) validExpirationYear(strict bool) bool {
	if strict {
		return p.ExpirationYear >= 2020 && p.ExpirationYear <= 2030
	} else {
		return p.ExpirationYear != 0
	}
}

func (p day4_passport) validHeight(strict bool) bool {
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

func (p day4_passport) validHairColor(strict bool) bool {
	if strict {
		match, _ := regexp.MatchString("^#[0-9a-f]{6}$", p.HairColor)
		return match
	} else {
		return p.HairColor != ""
	}
}

func (p day4_passport) validEyeColor(strict bool) bool {
	if strict {
		match, _ := regexp.MatchString("^(amb|blu|brn|gry|grn|hzl|oth)$", p.EyeColor)
		return match
	} else {
		return p.EyeColor != ""
	}
}

func (p day4_passport) validPassportID(strict bool) bool {
	if strict {
		match, _ := regexp.MatchString("^[0-9]{9}$", p.PassportID)
		return match
	} else {
		return p.PassportID != ""
	}
}

func (p day4_passport) validCountryID(strict bool) bool {
	// no validation
	return true
}
