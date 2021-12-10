package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/codegen"
	"github.com/davejhilton/adventofcode/log"

	_ "github.com/davejhilton/adventofcode/challenges/loader"
)

func main() {
	var create bool
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
	parseArgs(&create, &year, &day, &part)

	log.EnableDebugLogs(verbose)

	if create {
		err := codegen.GenerateChallengeTemplate(year, day)
		if err != nil {
			fmt.Printf("error generating code file: %s\n", err)
			fmt.Println("Aborting.")
			os.Exit(1)
		}
		fmt.Println("\nDone.")
		os.Exit(0)
	}

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

func parseArgs(create *bool, year *int, day *int, part *int) {
	if arg0 := flag.Arg(0); arg0 != "" {
		if arg0 == "create" {
			*create = true
		} else {
			n, err := strconv.Atoi(arg0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: First argument (year) must be a number, got: '%s'\n", arg0)
				os.Exit(1)
			}
			*year = n
		}
	}
	if arg1 := flag.Arg(1); arg1 != "" {
		argName := "day"
		if *create {
			argName = "year"
		}
		n, err := strconv.Atoi(arg1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Second argument (%s) must be a number, got: '%s'\n", argName, arg1)
			os.Exit(1)
		}
		if *create {
			*year = n
		} else {
			*day = n
		}
	}
	if arg2 := flag.Arg(2); arg2 != "" {
		argName := "part"
		if *create {
			argName = "day"
		}
		n, err := strconv.Atoi(arg2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Third argument (%s) must be a number, got: '%s'\n", argName, arg2)
			os.Exit(1)
		}
		if *create {
			*day = n
		} else {
			*part = n
		}
	}
}

func printUsage() {
	prgParts := strings.Split(os.Args[0], "/")
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [OPTIONS] [year] [day] [part]\n\n", prgParts[len(prgParts)-1])
	// flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "year   which year's challenge to run (integer)")
	fmt.Fprintln(os.Stderr, "day    which day's challenge to run (integer)")
	fmt.Fprintln(os.Stderr, "part   which of the specified [day]'s challenges to run (integer)")
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	fmt.Fprintln(os.Stderr, "     --example, -e      the example input file number to use (integer)")
	fmt.Fprintln(os.Stderr, "     --verbose, -v      enable debug logging")
	//
	fmt.Fprintln(os.Stderr)
	//
	fmt.Fprintf(os.Stderr, "\n\nCodeGen Usage:\n  %s create [year] [day]\n\n", prgParts[len(prgParts)-1])
	fmt.Fprintln(os.Stderr, "year   which year folder to create within (integer)")
	fmt.Fprintln(os.Stderr, "day    which day to create starter files for (integer)")
}
