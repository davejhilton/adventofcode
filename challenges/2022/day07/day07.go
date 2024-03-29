package aoc2022_day7

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

var (
	cd_regex   = regexp.MustCompile(`^\$\scd\s(.*)$`)
	dir_regex  = regexp.MustCompile(`^dir\s(.*)$`)
	file_regex = regexp.MustCompile(`^(\d+)\s(.*)$`)
)

func part1(input []string) (string, error) {
	root := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", root.TreeString(0))

	var sumSizes func(*dir, int) int

	sumSizes = func(d *dir, threshold int) int {
		total := 0
		if d.Size() <= threshold {
			total += d.Size()
		}
		for _, n := range d.contents {
			if d, ok := n.(*dir); ok {
				total += sumSizes(d, threshold)
			}
		}
		return total
	}

	result := sumSizes(root, 100000)

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	root := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", root.TreeString(0))

	var findDirToDelete func(*dir, int, int) int

	findDirToDelete = func(d *dir, min int, smallest int) int {
		size := d.Size()
		if size >= min {
			if size < smallest {
				smallest = util.Min(smallest, size)
			}
			for _, n := range d.contents {
				if d, ok := n.(*dir); ok {
					smallest = findDirToDelete(d, min, smallest)
				}
			}
		}
		return smallest
	}

	diskSize := 70000000
	freeSpaceRequirement := 30000000
	minToDelete := freeSpaceRequirement - (diskSize - root.Size())
	result := findDirToDelete(root, minToDelete, diskSize)

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) *dir {
	root := &dir{
		name:     "/",
		contents: make(map[string]node, 0),
		parent:   nil,
	}

	var cwd *dir = root
	for _, s := range input {
		switch s[0:4] {
		case "$ ls":
			// nothing to do
		case "dir ":
			matches := dir_regex.FindStringSubmatch(s)
			name := matches[1]
			cwd.AddChild(&dir{
				name:     name,
				contents: make(map[string]node),
				parent:   cwd,
			})
		case "$ cd":
			matches := cd_regex.FindStringSubmatch(s)
			name := matches[1]
			if name == "/" {
				cwd = root
			} else if name == ".." {
				cwd = cwd.Parent()
			} else {
				cwd = cwd.contents[name].(*dir)
			}
		default:
			matches := file_regex.FindStringSubmatch(s)
			cwd.AddChild(&file{
				name:   matches[2],
				size:   util.Atoi(matches[1]),
				parent: cwd,
			})
		}
	}
	return root
}

func init() {
	challenges.RegisterChallengeFunc(2022, 7, 1, "day07.txt", part1)
	challenges.RegisterChallengeFunc(2022, 7, 2, "day07.txt", part2)
}

// ===========================
//      DIRECTORY STRUCT
// ===========================

type dir struct {
	name     string
	contents map[string]node
	parent   *dir
}

func (d dir) Name() string {
	return d.name
}

func (d dir) Size() int {
	var size int
	for _, n := range d.contents {
		size += n.Size()
	}
	return size
}

func (d dir) Parent() *dir {
	return d.parent
}

func (d dir) String() string {
	return fmt.Sprintf("%s (dir)", d.name)
}

func (d dir) TreeString(depth int) string {
	var b strings.Builder
	pad := " "
	childPad := " └─── "
	if depth != 0 {
		pad = fmt.Sprintf(" %s└─── ", strings.Repeat("     ", depth-1))
		childPad = fmt.Sprintf("     %s", pad)
	}
	b.WriteString(fmt.Sprintf("%s%s", pad, d))
	for _, n := range d.contents {
		if d, ok := n.(*dir); ok {
			b.WriteString(fmt.Sprintf("\n%s", d.TreeString(depth+1)))
		} else {
			b.WriteString(fmt.Sprintf("\n%s%s", childPad, n))
		}
	}
	return b.String()
}

func (d *dir) AddChild(n node) {
	if n != nil {
		d.contents[n.Name()] = n
	}
}

// ===========================
//         FILE STRUCT
// ===========================

type file struct {
	name   string
	size   int
	parent *dir
}

func (f file) Name() string {
	return f.name
}

func (f file) Size() int {
	return f.size
}

func (f file) Parent() *dir {
	return f.parent
}

func (f file) String() string {
	return fmt.Sprintf("%s (file, size=%d)", f.name, f.size)
}

// =====================================
//      FILESYSTEM NODE INTERFACE
// =====================================

type node interface {
	Name() string
	Size() int
	Parent() *dir
	String() string
}
