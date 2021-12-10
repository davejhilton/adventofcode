package codegen

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var insertPattern = regexp.MustCompile(`.*CODEGEN: INSERT HERE*`)

func GenerateChallengeTemplate(year int, day int) error {

	dirName := fmt.Sprintf("challenges/%d/day%.2d", year, day)
	goFile := fmt.Sprintf("%s/day%.2d.go", dirName, day)

	if _, err := os.Stat(goFile); err == nil {
		// file already exists - don't overwrite it!
		return fmt.Errorf("file '%s' already exists", goFile)
	} else if !errors.Is(err, os.ErrNotExist) {
		// we expected os.ErrNotExist... but got some other error
		return err
	}

	if err := os.Mkdir(dirName, 0755); err != nil {
		return err
	}

	writer, err := os.Create(goFile)
	if err != nil {
		return err
	}
	defer writer.Close()

	tpl, err := template.New("NewChallenge").ParseFiles("codegen/challenge-template.tpl")
	if err != nil {
		fmt.Printf("ERROR rendering template: %s\n", err)
		return err
	}

	tpl.ExecuteTemplate(writer, "challenge-template.tpl", &map[string]int{
		"Year": year,
		"Day":  day,
	})

	fmt.Printf("Created challenge file:      %s\n", goFile)

	inputFile := fmt.Sprintf("%s/day%.2d.txt", dirName, day)
	f2, err := os.Create(inputFile)
	if err != nil {
		return err
	}
	f2.Close()
	fmt.Printf("Created empty input file:    %s\n", inputFile)

	exampleInputFile := fmt.Sprintf("%s/day%.2d_example1.txt", dirName, day)
	f3, err := os.Create(exampleInputFile)
	if err != nil {
		return err
	}
	f3.Close()
	fmt.Printf("Created empty example file:  %s\n", exampleInputFile)

	loaderPath := "challenges/loader/loader.go"
	loaderFile, err := os.OpenFile(loaderPath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer loaderFile.Close()

	lines := make([]string, 0)

	pkgPath := fmt.Sprintf("github.com/davejhilton/adventofcode/challenges/%d/day%.2d", year, day)
	importStmt := fmt.Sprintf(`	_ "%s"`, pkgPath)
	importAlreadyExists := false

	scanner := bufio.NewScanner(loaderFile)
	for scanner.Scan() {
		text := scanner.Text()
		if text == importStmt {
			importAlreadyExists = true
		}
		if insertPattern.MatchString(text) && !importAlreadyExists {
			fmt.Printf("Modifying loader.go to add import for '%s'.\n", pkgPath)
			lines = append(lines, importStmt)
		}
		lines = append(lines, text)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if importAlreadyExists {
		return nil
	}

	loaderFile.Truncate(0)
	loaderFile.Seek(0, 0)
	loaderFile.WriteString(strings.Join(lines, "\n"))

	return nil
}
