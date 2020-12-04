package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/davejhilton/adventofcode2020/challenges"
	"github.com/davejhilton/adventofcode2020/log"
)

func main() {

	var day int = 1
	var part int = 1
	var verbose bool
	var noColor bool

	flag.Usage = printUsage
	flag.BoolVar(&verbose, "v", false, "verbose: if enabled, will print debug logs")
	flag.BoolVar(&noColor, "nocolor", false, "nocolor: if true, debug logs will not have coloring")
	flag.Parse()
	parseArgs(&day, &part)

	log.EnableDebugLogs(verbose)
	log.EnableColors(!noColor)

	challenge, err := challenges.GetChallenge(day, part)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("The specified day number does not exist!\nError: %w", err))
		os.Exit(1)
	}

	solution, err := challenge.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("Error running challenge!\nError: %w", err))
		os.Exit(1)
	}

	fmt.Printf("====================\n%s\n====================\nSolution: %s\n====================\n", challenge.Name(), solution)
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
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [day] [part]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n [day]\twhich day's challenge to run (integer)\n")
	fmt.Fprintf(os.Stderr, "[part]\twhich of the specified [day]'s challenges to run (integer)\n")
}
