package main

import (
	"fmt"
	"os"
	"text/template"
)

type tplVars struct {
	Year string
	Day  string
}

func newChallengeFromTemplate(year int, day int) error {
	tpl, err := template.New("NewChallenge").ParseFiles("challenges/template/challenge.tpl")
	if err != nil {
		fmt.Printf("ERROR! %s\n", err)
		return err
	}
	tpl.ExecuteTemplate(os.Stdout, "challenge.tpl", &tplVars{
		Year: fmt.Sprintf("%d", year),
		Day:  fmt.Sprintf("%.2d", day),
	})
	fmt.Println()
	return nil
}

/*
func printUsagex() {
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
*/
