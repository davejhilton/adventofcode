package codegen

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/davejhilton/adventofcode/log"
)

const msgWidth = -35

var insertPattern = regexp.MustCompile(`.*CODEGEN: INSERT HERE*`)

func GenerateChallengeTemplate(year int, day int) error {

	yearPath := fmt.Sprintf("./challenges/%d", year)
	if err := createDirectory(yearPath, true); err != nil {
		return err
	}
	dirPath := fmt.Sprintf("%s/day%.2d", yearPath, day)
	if err := createDirectory(dirPath, false); err != nil {
		return err
	}
	fmt.Printf("%s%s\n", log.Colorize("Created challenge dir:", log.Green, msgWidth), dirPath)

	goFilePath := fmt.Sprintf("%s/day%.2d.go", dirPath, day)
	if err := generateCodeFile(goFilePath, year, day); err != nil {
		return err
	}

	loaderFilePath := "./challenges/loader/loader.go"
	if err := addPackageImport(loaderFilePath, year, day); err != nil {
		return err
	}

	inputFilePath := fmt.Sprintf("%s/day%.2d.txt", dirPath, day)
	if err := generateInputFile(inputFilePath); err != nil {
		return err
	}

	exampleFilePath := fmt.Sprintf("%s/day%.2d_example1.txt", dirPath, day)
	if err := generateInputFile(exampleFilePath); err != nil {
		return err
	}

	return nil
}

func createDirectory(dirPath string, silent bool) error {
	err := os.Mkdir(dirPath, 0755)
	if errors.Is(err, os.ErrExist) {
		if !silent {
			fmt.Printf("%s%s\n", log.Colorize("Dir already exists:", log.Yellow, msgWidth), dirPath)
		}
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func generateCodeFile(filePath string, year int, day int) error {
	if _, err := os.Stat(filePath); err == nil {
		// file already exists - don't overwrite it!
		fmt.Printf("%s%s\n", log.Colorize("Challenge file already exists:", log.Yellow, msgWidth), filePath)
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		// we expected os.ErrNotExist... but got some other error
		return err
	}

	writer, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer writer.Close()

	tpl, err := template.New("NewChallenge").ParseFiles("codegen/challenge-template.tpl")
	if err != nil {
		fmt.Printf("%s %s\n", log.Colorize("ERROR rendering template:", log.Red, 0), err)
		return err
	}

	tpl.ExecuteTemplate(writer, "challenge-template.tpl", &map[string]int{
		"Year": year,
		"Day":  day,
	})

	fmt.Printf("%s%s\n", log.Colorize("Created challenge file:", log.Green, msgWidth), filePath)
	return nil
}

func generateInputFile(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			// we expected os.ErrNotExist... but got some other error
			return err
		}
	} else {
		// file already exists - don't overwrite it!
		fmt.Printf("%s%s\n", log.Colorize("Input file already exists:", log.Yellow, msgWidth), filePath)
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("%s%s\n", log.Colorize("Created empty input file:", log.Green, msgWidth), filePath)
	return nil
}

func addPackageImport(loaderFilePath string, year int, day int) error {
	loaderFile, err := os.OpenFile(loaderFilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer loaderFile.Close()

	lines := make([]string, 0)

	pkgPath := fmt.Sprintf("github.com/davejhilton/adventofcode/challenges/%d/day%.2d", year, day)
	importStmt := fmt.Sprintf(`	_ "%s"`, pkgPath)

	scanner := bufio.NewScanner(loaderFile)
	for scanner.Scan() {
		text := scanner.Text()
		if text == importStmt {
			fmt.Printf("%s%s\n", log.Colorize("Package import already exists:", log.Yellow, msgWidth), loaderFilePath)
			return nil
		}
		if insertPattern.MatchString(text) {
			fmt.Printf("%s%s\n", log.Colorize("Added package import statement:", log.Green, msgWidth), loaderFilePath)
			lines = append(lines, importStmt)
		}
		lines = append(lines, text)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	loaderFile.Truncate(0)
	loaderFile.Seek(0, 0)
	loaderFile.WriteString(strings.Join(lines, "\n"))

	return nil
}
