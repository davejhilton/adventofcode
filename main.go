package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/davejhilton/adventofcode2020/challenges"
)

func main() {
	day := 0
	part := 1
	if len(os.Args) > 1 {
		d, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("ERROR: Unknown Day: '%s'\n", os.Args[1])
			os.Exit(1)
		}
		day = d
	}
	if len(os.Args) > 2 {
		part, _ = strconv.Atoi(os.Args[2])
	}
	if len(os.Args) > 2 {
		p, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("ERROR: Unknown Challenge: Day '%d' Part '%s'\n", day, os.Args[2])
			os.Exit(1)
		}
		part = p
	}

	challenge, err := challenges.GetChallenge(day, part)
	if err != nil {
		fmt.Println("Error finding challenge!", err)
		os.Exit(1)
	}
	solution, err := challenge.Run()
	if err != nil {
		fmt.Println("Error running challenge!", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n====================\nSolution: %s\n====================\n", challenge.Name(), solution)
}
