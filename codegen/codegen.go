package codegen

import (
	"errors"
	"fmt"
	"os"
	"text/template"
)

func GenerateChallengeTemplate(year int, day int) error {
	goFile := fmt.Sprintf("challenges/%d/day%.2d.go", year, day)

	if _, err := os.Stat(goFile); err == nil {
		// file already exists - don't overwrite it!
		return fmt.Errorf("file '%s' already exists", goFile)
	} else if !errors.Is(err, os.ErrNotExist) {
		// we expected os.ErrNotExist... but got some other error
		return err
	}

	writer, err := os.Create(goFile)
	if err != nil {
		return err
	}
	defer writer.Close()

	tpl, err := template.New("NewChallenge").ParseFiles("challenges/template/challenge.tpl")
	if err != nil {
		fmt.Printf("ERROR rendering template: %s\n", err)
		return err
	}

	tpl.ExecuteTemplate(writer, "challenge.tpl", &map[string]string{
		"Year": fmt.Sprintf("%d", year),
		"Day":  fmt.Sprintf("%.2d", day),
	})

	fmt.Printf("Created challenge file:      %s\n", goFile)

	inputFile := fmt.Sprintf("inputs/%d/day%.2d.txt", year, day)
	f2, err := os.Create(inputFile)
	if err != nil {
		return err
	}
	f2.Close()
	fmt.Printf("Created empty input file:    %s\n", inputFile)

	exampleInputFile := fmt.Sprintf("inputs/%d/day%.2d_example1.txt", year, day)
	f3, err := os.Create(exampleInputFile)
	if err != nil {
		return err
	}
	f3.Close()
	fmt.Printf("Created empty example file:  %s\n", exampleInputFile)

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