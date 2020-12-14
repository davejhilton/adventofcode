package challenges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

func day8_part1(input []string) (string, error) {
	program := day8_parseProgram(input)
	instrNum := 1
	acc := 0
	pc := 0
	execHistory := make([]int, len(program))

	log.Debugln("RUNNING:")
	for {

		if pc >= len(program) {
			log.Debugln("EOP")
			break
		}
		log.Debugf("%3d. Executing Instruction: (%d) %s %3d\n", instrNum, pc, program[pc].Op, program[pc].Arg)
		if execHistory[pc] != 0 {
			log.Debugln(log.Colorize("INFINITE LOOP", log.Red, 0))
			break
		}
		execHistory[pc] = instrNum
		instrNum += 1
		pc = day8_execInstruction(program[pc], pc, &acc)
	}
	return fmt.Sprintf("%d", acc), nil
}

func day8_execInstruction(instr day8_instruction, pc int, acc *int) int {
	switch instr.Op {
	case "nop":
		return pc + 1
	case "acc":
		*acc += instr.Arg
		log.Debugf("\tacc: %d\n", *acc)
		return pc + 1
	case "jmp":
		return pc + instr.Arg
	}
	return pc + 1
}

func day8_part2(input []string) (string, error) {
	program := day8_parseProgram(input)
	acc, _ := day8_tryProgram(program, 0, 0, make([]int, len(program)), false, "")
	return fmt.Sprintf("%d", acc), nil
}

type day8_instruction struct {
	Op  string
	Arg int
}

func day8_tryProgram(program []day8_instruction, pc int, acc int, execHistory []int, altered bool, indent string) (int, bool) {

	instrNum := 1

	var success bool

	log.Debugf("%sTRYING:\n", indent)
	for {
		if pc >= len(program) {
			log.Debugf("%s%s\n", indent, log.Colorize("EOP!", log.Green, 0))
			success = true
			break
		}
		log.Debugf("%s%3d. Executing Instruction: (%d) %s %3d\n", indent, instrNum, pc, program[pc].Op, program[pc].Arg)
		if execHistory[pc] != 0 {
			log.Debugf("%s%s\n", indent, log.Colorize("INFINITE LOOP", log.Red, 0))
			success = false
			break
		}

		if !altered && program[pc].Op == "nop" {
			program[pc].Op = "jmp"
			log.Debugf("%sTRYING ALTERNATE PATH...\n", indent)
			altExecHistory := append(make([]int, len(execHistory)), execHistory...)
			altAcc, altSuccess := day8_tryProgram(program, pc, acc, altExecHistory, true, fmt.Sprintf("\t%s", indent))
			if altSuccess {
				acc, success = altAcc, altSuccess
				break
			} else {
				program[pc].Op = "nop"
			}
		} else if !altered && program[pc].Op == "jmp" {
			program[pc].Op = "nop"
			log.Debugf("%sTRYING ALTERNATE PATH...\n", indent)
			altExecHistory := append(make([]int, len(execHistory)), execHistory...)
			altAcc, altSuccess := day8_tryProgram(program, pc, acc, altExecHistory, true, fmt.Sprintf("\t%s", indent))
			if altSuccess {
				acc, success = altAcc, altSuccess
				break
			} else {
				program[pc].Op = "jmp"
			}
		}
		execHistory[pc] = instrNum
		instrNum += 1
		pc = day8_execInstruction(program[pc], pc, &acc)
	}
	return acc, success
}

func day8_parseProgram(input []string) []day8_instruction {
	program := make([]day8_instruction, 0, len(input))
	log.Debugln("PROGRAM:")
	for i, line := range input {
		parts := strings.Split(line, " ")
		program = append(program, day8_instruction{
			Op: parts[0],
		})
		program[i].Arg, _ = strconv.Atoi(parts[1])
		log.Debugf("%3d:\t%s %3d\n", i, program[i].Op, program[i].Arg)
	}
	log.Debugln()
	return program
}

func init() {
	registerChallengeFunc(8, 1, "day08.txt", day8_part1)
	registerChallengeFunc(8, 2, "day08.txt", day8_part2)
}
