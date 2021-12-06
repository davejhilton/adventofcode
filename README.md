# Advent of Code 2020 - 2021

This is my attempt at following along with the [Advent of Code](https://adventofcode.com) challenges for 2020 and 2021, using Golang.


## To run the challenges:

- install golang 1.17 or later
- clone this repo
- `cd` into the root of the repo
- run `go run main.go <year> <day> <part>`, where `<year>` is 2020 or 2021, `<day>` is the day number, and `<part>` is which challenge number for that day.
  - e.g., to run challenge #2 for Day 1 in 2021's event, you'd do `go run main.go 2021 1 2`

## To generate starter files for new challenges:

- run `go run main.go create <year> <day>`
  - e.g., to generate files for Day 7 in 2021's event, you'd do `go run main.go create 2021 7`
  - this would create:
    - `challenges/2021/day07.go` - prepopulated with the basic code structure needed for running a challenge
    - `inputs/2021/day07.txt` - an empty input file
    - `inputs/2021/day07_example1.txt` - an empty example input file
