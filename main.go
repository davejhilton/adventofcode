package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"

	_ "github.com/davejhilton/adventofcode/challenges/2020"
	_ "github.com/davejhilton/adventofcode/challenges/2021"
)

func main() {
	var year int = 1
	var day int = 1
	var part int = 1
	var verbose bool
	var printTime bool
	var exampleNum int

	flag.Usage = printUsage
	flag.BoolVar(&verbose, "verbose", false, "verbose:     if true, will print debug logs")
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&printTime, "time", true, "time:        if true, print the execution time after the solution")
	flag.BoolVar(&printTime, "t", true, "")
	flag.IntVar(&exampleNum, "example", -1, "example-num: the number of the example input file to run")
	flag.IntVar(&exampleNum, "e", -1, "")
	flag.Parse()
	parseArgs(&year, &day, &part)

	log.EnableDebugLogs(verbose)

	// newChallengeFromTemplate(year, day)

	// os.Exit(0)

	challenge, err := challenges.GetChallenge(year, day, part, exampleNum)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("the specified year/day/challenge does not exist!\nError: %w", err))
		os.Exit(1)
	}

	startTime := time.Now()

	solution, err := challenge.Run()

	execTime := time.Since(startTime)

	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error running challenge!\nError: %w", err))
		os.Exit(1)
	}

	fmt.Printf("====================\n%s\n", challenge.Name())
	fmt.Printf("%s\n", challenge.InputFileName)
	fmt.Printf("====================\nSolution: %s\n====================\n", solution)
	if printTime {
		fmt.Printf("time: %v\n", execTime)
	}
}

func parseArgs(year *int, day *int, part *int) {
	if arg0 := flag.Arg(0); arg0 != "" {
		y, err := strconv.Atoi(arg0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: First argument (year) must be a number, got: '%s'\n", arg0)
			os.Exit(1)
		}
		*year = y
	}
	if arg1 := flag.Arg(1); arg1 != "" {
		d, err := strconv.Atoi(arg1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Second argument (day) must be a number, got: '%s'\n", arg1)
			os.Exit(1)
		}
		*day = d
	}
	if arg2 := flag.Arg(2); arg2 != "" {
		p, err := strconv.Atoi(arg2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Third argument (part) must be a number, got: '%s'\n", arg2)
			os.Exit(1)
		}
		*part = p
	}
}

func printUsage() {
	prgParts := strings.Split(os.Args[0], "/")
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [OPTIONS] [year] [day] [part]\n\n", prgParts[len(prgParts)-1])
	// flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "[OPTIONS]")
	fmt.Fprintln(os.Stderr, "    --example, -e      [int] the example input file number to use")
	fmt.Fprintln(os.Stderr, "    --verbose, -v      enable debug logging")
	//
	fmt.Fprintln(os.Stderr, "[year]                 which year's challenge to run (integer)")
	fmt.Fprintln(os.Stderr, "[day]                  which day's challenge to run (integer)")
	fmt.Fprintln(os.Stderr, "[part]                 which of the specified [day]'s challenges to run (integer)")
}
