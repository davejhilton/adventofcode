package aoc2023_day20

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

const (
	FLIP_FLOP   = "Flip-Flop"
	CONJUNCTION = "Conjunction"
	BROADCAST   = "Broadcast"
	EMPTY       = "Empty"
)

type Pulse struct {
	From  string
	Value int
	To    string
}

type CPU struct {
	Modules          map[string]*Module
	ButtonPushes     int
	PulseCounts      map[int]int
	bus              []Pulse
	ModuleSendCounts map[string]map[int]int
}

func NewCPU(modules map[string]*Module) *CPU {
	// initialize all Conjunction modules to have a memory of 0 for each input

	modulePulseCounts := make(map[string]map[int]int)
	for _, m := range modules {
		for _, dest := range m.Destinations {
			if dest.Type == CONJUNCTION {
				dest.Memory[m.Name] = 0
			}
		}
		modulePulseCounts[m.Name] = make(map[int]int)
	}

	for _, m := range modules {
		if m.Type == CONJUNCTION {
			log.Debugf("Conjunction %s initialized with memory:\n", m.Name)
			for input, value := range m.Memory {
				log.Debugf("    %s: %d\n", input, value)
			}
		}
	}

	return &CPU{
		Modules:          modules,
		ButtonPushes:     0,
		PulseCounts:      make(map[int]int),
		bus:              make([]Pulse, 0),
		ModuleSendCounts: make(map[string]map[int]int),
	}
}

func (c *CPU) Emit(pulse Pulse) {
	c.bus = append(c.bus, pulse)
	c.PulseCounts[pulse.Value]++
	if c.ModuleSendCounts[pulse.From] == nil {
		c.ModuleSendCounts[pulse.From] = make(map[int]int)
	}
	c.ModuleSendCounts[pulse.From][pulse.Value]++
}

func (c *CPU) PushButton() {
	log.Debugf("\n=====================\nPushButton #%d\n=====================\n", c.ButtonPushes+1)
	c.ButtonPushes++
	c.Emit(Pulse{"button", 0, "broadcaster"})
	for len(c.bus) > 0 {
		pulse := c.bus[0]
		c.bus = c.bus[1:]
		c.Modules[pulse.To].Receive(pulse, c)
	}
}

func (c *CPU) String() string {
	var sb strings.Builder
	for _, m := range c.Modules {
		sb.WriteString(fmt.Sprintf("%s\n", m))
	}
	return sb.String()
}

type Module struct {
	Name         string
	Type         string
	Destinations []*Module
	CurrentValue int
	Memory       map[string]int
}

func (m *Module) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s (%s) -> ", m.Name, m.Type))
	for _, dest := range m.Destinations {
		sb.WriteString(fmt.Sprintf("%s, ", dest.Name))
	}
	return sb.String()
}

func (m *Module) Receive(pulse Pulse, c *CPU) {
	log.Debugf("%s (%s) received %d\n", m.Name, m.Type, pulse.Value)
	switch m.Type {
	case FLIP_FLOP:
		if pulse.Value == 0 {
			if m.CurrentValue == 1 {
				m.CurrentValue = 0
			} else {
				m.CurrentValue = 1
			}
			for _, dest := range m.Destinations {
				log.Debugf("    sending %d to %s\n", m.CurrentValue, dest.Name)
				c.Emit(Pulse{m.Name, m.CurrentValue, dest.Name})
			}
		}
	case CONJUNCTION:
		var out int
		m.Memory[pulse.From] = pulse.Value
		for _, v := range m.Memory {
			if v == 0 {
				out = 1
				break
			}
		}
		for _, dest := range m.Destinations {
			log.Debugf("    sending %d to %s\n", out, dest.Name)
			c.Emit(Pulse{m.Name, out, dest.Name})
		}
	case BROADCAST:
		for _, dest := range m.Destinations {
			log.Debugf("    sending %d to %s\n", pulse.Value, dest.Name)
			c.Emit(Pulse{m.Name, pulse.Value, dest.Name})
		}
	case EMPTY:
		// do nothing
	}

}

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	cpu := NewCPU(parsed)

	for i := 0; i < 1000; i++ {
		cpu.PushButton()
	}

	result := cpu.PulseCounts[1] * cpu.PulseCounts[0]
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	cpu := NewCPU(parsed)

	rxInput := ""
	for _, m := range cpu.Modules {
		for _, dest := range m.Destinations {
			if dest.Name == "rx" {
				rxInput = m.Name
				break
			}
		}
	}
	rxParentInputs := make(map[string][]int)
	for _, m := range cpu.Modules {
		for _, dest := range m.Destinations {
			if dest.Name == rxInput {
				rxParentInputs[m.Name] = make([]int, 0)
				log.Debugf("Found rx input: %s\n", m.Name)
			}
		}
	}
	i := 0
	for cpu.ModuleSendCounts["rx"][0] < 1 {
		i++
		if i%100000 == 0 {
			log.Debugf("PushButton #%d\n", cpu.ButtonPushes+1)
		}
		cpu.PushButton()
		pending := len(rxParentInputs)
		for name, val := range rxParentInputs {
			if len(val) < 2 {
				n1 := cpu.ModuleSendCounts[name][1]
				if n1 == 1 && len(val) == 0 {
					val = append(val, i)
					log.Debugf("Input %s got a 1 after %d cycles\n", name, i)
					rxParentInputs[name] = val
				} else if n1 > 1 && len(val) == 1 {
					val = append(val, i)
					log.Debugf("Found cycle for %s: %d\n", name, val[1]-val[0])
					rxParentInputs[name] = val
					pending--
				}
			} else {
				pending--
			}
		}
		if pending == 0 {
			break
		}
	}

	log.Debugf("module counts: %v\n", cpu.ModuleSendCounts)
	result := 1
	for _, val := range rxParentInputs {
		result = util.LCM(result, val[1]-val[0])
	}
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) map[string]*Module {
	modules := make(map[string]*Module)
	destMap := make(map[string][]string)
	for _, s := range input {
		parts := strings.Split(s, " -> ")
		name := parts[0]
		mType := EMPTY
		if strings.HasPrefix(name, "%") {
			mType = FLIP_FLOP
			name = name[1:]
		} else if strings.HasPrefix(name, "&") {
			mType = CONJUNCTION
			name = name[1:]
		} else if name == "broadcaster" {
			mType = BROADCAST
		}

		if len(parts) > 1 {
			destMap[name] = strings.Split(parts[1], ", ")
		}
		modules[name] = &Module{
			Name:         name,
			Type:         mType,
			Destinations: make([]*Module, 0),
			CurrentValue: 0,
			Memory:       make(map[string]int),
		}
	}
	for name, dests := range destMap {
		for _, destStr := range dests {
			if _, ok := modules[destStr]; !ok {
				modules[destStr] = &Module{
					Name:         destStr,
					Type:         EMPTY,
					Destinations: make([]*Module, 0),
					CurrentValue: 0,
					Memory:       make(map[string]int),
				}
			}
			modules[name].Destinations = append(modules[name].Destinations, modules[destStr])
		}
	}
	return modules
}

func init() {
	challenges.RegisterChallengeFunc(2023, 20, 1, "day20.txt", part1)
	challenges.RegisterChallengeFunc(2023, 20, 2, "day20.txt", part2)
}
