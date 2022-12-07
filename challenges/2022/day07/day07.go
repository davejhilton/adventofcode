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
			if n.Type() == "dir" {
				total += sumSizes(n.(*dir), threshold)
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
				if n.Type() == "dir" {
					smallest = findDirToDelete(n.(*dir), min, smallest)
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
		path:     "/",
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
				path:     fmt.Sprintf("%s/%s", cwd.Path(), name),
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
				if d, ok := cwd.contents[name]; ok && d.Type() == "dir" {
					cwd = cwd.contents[name].(*dir)
				} else {
					fmt.Println(log.Colorize(fmt.Sprintf("ERROR: Cannot cd into directory %s! Does not exist or is not a directory! (cwd: %s)", name, cwd.Path()), log.Red, 0))
				}
			}
		default:
			matches := file_regex.FindStringSubmatch(s)
			cwd.AddChild(&file{
				name:   matches[2],
				path:   fmt.Sprintf("%s/%s", cwd.Path(), matches[2]),
				size:   util.Atoi(matches[1]),
				parent: cwd,
			})
		}
	}
	return root
}

type node interface {
	Type() string
	Name() string
	Size() int
	Path() string
	Parent() *dir
	String() string
}

type dir struct {
	name     string
	path     string
	contents map[string]node
	parent   *dir
}

func (d dir) Type() string {
	return "dir"
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

func (d dir) Path() string {
	return d.path
}

func (d dir) Parent() *dir {
	return d.parent
}

func (d dir) String() string {
	return fmt.Sprintf("%s (dir)", d.name)
}

func (d dir) TreeString(depth int) string {
	var b strings.Builder
	pad := strings.Repeat("  ", depth)
	b.WriteString(fmt.Sprintf("%s- %s", pad, d))
	for _, n := range d.contents {
		if n.Type() == "file" {
			b.WriteString(fmt.Sprintf("\n  %s- %s", pad, n))
		} else {
			b.WriteString(fmt.Sprintf("\n%s", n.(*dir).TreeString(depth+1)))
		}
	}
	return b.String()
}

func (d *dir) AddChild(n node) {
	if n != nil {
		d.contents[n.Name()] = n
	}
}

type file struct {
	name   string
	path   string
	size   int
	parent *dir
}

func (f file) Type() string {
	return "file"
}

func (f file) Name() string {
	return f.name
}

func (f file) Size() int {
	return f.size
}

func (f file) Path() string {
	return f.path
}

func (f file) Parent() *dir {
	return f.parent
}

func (f file) String() string {
	return fmt.Sprintf("%s (file, size=%d)", f.name, f.size)
}

func init() {
	challenges.RegisterChallengeFunc(2022, 7, 1, "day07.txt", part1)
	challenges.RegisterChallengeFunc(2022, 7, 2, "day07.txt", part2)
}
