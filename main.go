package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davejhilton/adventofcode2020/challenges"
	"github.com/davejhilton/adventofcode2020/log"
)

func main() {

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
	parseArgs(&day, &part)

	log.EnableDebugLogs(verbose)

	challenge, err := challenges.GetChallenge(day, part, exampleNum)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("The specified day number does not exist!\nError: %w", err))
		os.Exit(1)
	}

	startTime := time.Now()

	solution, err := challenge.Run()

	execTime := time.Now().Sub(startTime)

	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("Error running challenge!\nError: %w", err))
		os.Exit(1)
	}

	fmt.Printf("====================\n%s\n", challenge.Name())
	fmt.Printf("%s\n", challenge.InputFileName)
	fmt.Printf("====================\nSolution: %s\n====================\n", solution)
	if printTime {
		fmt.Printf("time: %v\n", execTime)
	}
}

func parseArgs(day *int, part *int) {
	if arg0 := flag.Arg(0); arg0 != "" {
		d, err := strconv.Atoi(arg0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: First argument (day) must be a number, got: '%s'\n", arg0)
			os.Exit(1)
		}
		*day = d
	}
	if arg1 := flag.Arg(1); arg1 != "" {
		p, err := strconv.Atoi(arg1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Second argument (part) must be a number, got: '%s'\n", arg1)
			os.Exit(1)
		}
		*part = p
	}
}

func printUsage() {
	prgParts := strings.Split(os.Args[0], "/")
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [OPTIONS] [day] [part]\n\n", prgParts[len(prgParts)-1])
	// flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "[OPTIONS]")
	fmt.Fprintln(os.Stderr, "    --example, -e      [int] the example input file number to use")
	fmt.Fprintln(os.Stderr, "    --verbose, -v      enable debug logging")
	//
	fmt.Fprintln(os.Stderr, "[day]                  which day's challenge to run (integer)")
	fmt.Fprintln(os.Stderr, "[part]                 which of the specified [day]'s challenges to run (integer)")
}
