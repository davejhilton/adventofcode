package aoc2021_day16

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	mainPacket := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", mainPacket)

	result := getVersionSum(mainPacket)
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	mainPacket := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", mainPacket)

	result := eval(mainPacket)
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) *Packet {
	hexDigits := strings.Split(input[0], "")
	var sb strings.Builder
	for _, hexDigit := range hexDigits {
		num, _ := strconv.ParseUint(hexDigit, 16, 8)
		fmt.Fprintf(&sb, "%04b", num)
	}
	str := sb.String()
	// log.Debugf("PARSED BITS: %s (%d)\n", str, len(str))
	packet, bits := parsePacket(str)
	log.Debugln(packet)
	log.Debugf("Remaining bits: %d\n", len(bits))

	return packet
}

func getVersionSum(p *Packet) int {
	// log.Debugf("Version: %#v\n", p.Version)
	sum := p.Version
	if p.TypeID != LITERAL {
		for _, sp := range p.SubPackets {
			sum += getVersionSum(sp)
		}
	}
	return sum
}

func eval(p *Packet) (result int) {
	switch p.TypeID {
	case 0:
		// sum
		sum := 0
		for _, sp := range p.SubPackets {
			sum += eval(sp)
		}
		result = sum
	case 1:
		// product
		product := 1
		for _, sp := range p.SubPackets {
			product *= eval(sp)
		}
		result = product
	case 2:
		// minimum
		min := math.MaxInt
		for _, sp := range p.SubPackets {
			min = util.Min(min, eval(sp))
		}
		result = min
	case 3:
		// maximum
		max := math.MinInt
		for _, sp := range p.SubPackets {
			max = util.Max(max, eval(sp))
		}
		result = max
	case 4:
		// literal value
		result = p.LiteralValue
	case 5:
		// greater than
		if eval(p.SubPackets[0]) > eval(p.SubPackets[1]) {
			result = 1
		} else {
			result = 0
		}
	case 6:
		// less than
		if eval(p.SubPackets[0]) < eval(p.SubPackets[1]) {
			result = 1
		} else {
			result = 0
		}
	case 7:
		// equal to
		if eval(p.SubPackets[0]) == eval(p.SubPackets[1]) {
			result = 1
		} else {
			result = 0
		}
	}
	return result
}

func parsePacket(bits string) (*Packet, string) {
	version, typeId := parsePacketHeaders(bits[:6])
	log.Debugf("Version: %d, typeID: %d\n", version, typeId)

	bits = bits[6:]
	if typeId == LITERAL {
		return parseLiteralValuePacket(version, typeId, bits)
	} else {
		return parseOperatorPacket(version, typeId, bits)
	}
}

func parseLiteralValuePacket(version int, typeId int, bits string) (*Packet, string) {
	// dbg := fmt.Sprintf("RAW PACKET SEGMENTS: %3b %3b", version, typeId)
	packetLen := 6 // header
	value := ""
	chunkNum := 1
	log.Debugln("parsing LITERAL packet")
	for {
		more, chunk := bits[:1], bits[1:5]
		bits = bits[5:]
		// dbg = fmt.Sprintf("%s %s%s", dbg, more, chunk)
		packetLen += 5
		value = fmt.Sprintf("%s%s", value, chunk)
		chunkNum++
		if more == "0" {
			break
		}
	}

	// trailingBitCount := 4 - (packetLen % 4)
	// if trailingBitCount != 4 {
	// 	// tb := bits[:trailingBitCount]
	// 	// dbg = fmt.Sprintf("%s %s", dbg, tb)
	// 	bits = bits[trailingBitCount:]
	// }

	// log.Debugln(dbg)
	packet := &Packet{
		Version:      version,
		TypeID:       typeId,
		LiteralValue: binToInt(value),
	}
	return packet, bits
}

func parseOperatorPacket(version int, typeId int, bits string) (*Packet, string) {

	lengthType := bits[:1]
	bits = bits[1:]
	subPackets := make([]*Packet, 0)
	if lengthType == "1" {
		// numPackets, 11 bits
		numPackets := binToInt(bits[:11])
		bits = bits[11:]
		log.Debugf("parsing OPERATOR packet with %d sub packets\n", numPackets)
		var p *Packet
		for i := 0; i < numPackets; i++ {
			p, bits = parsePacket(bits)
			subPackets = append(subPackets, p)
		}
	} else {
		// totalLength, 15 bits
		numBits := binToInt(bits[:15])
		bits = bits[15:]
		innerBits := bits[:numBits]
		bits = bits[numBits:]
		log.Debugf("parsing OPERATOR packet with %d inner bits\n", numBits)
		var p *Packet
		for len(innerBits) > 0 {
			p, innerBits = parsePacket(innerBits)
			subPackets = append(subPackets, p)
		}
	}
	packet := &Packet{
		Version:    version,
		TypeID:     typeId,
		SubPackets: subPackets,
	}
	return packet, bits
}

func parsePacketHeaders(sixBits string) (version int, typeId int) {
	return binToInt(sixBits[:3]), binToInt(sixBits[3:6])
}

func binToInt(bits string) int {
	v, _ := strconv.ParseUint(bits, 2, 64)
	return int(v)
}

const LITERAL = 4

type Packet struct {
	Version      int
	TypeID       int
	LiteralValue int
	SubPackets   []*Packet
}

func (p *Packet) String() string {
	if p.TypeID == LITERAL {
		return fmt.Sprintf("PACKET[ Version: %d - Type: %d (LITERAL)  - Value: %d ]\n", p.Version, p.TypeID, p.LiteralValue)
	}
	return fmt.Sprintf("PACKET[ Version: %d - Type: %d (OPERATOR) - SubPackets: %d ]\n", p.Version, p.TypeID, len(p.SubPackets))
}

func init() {
	challenges.RegisterChallengeFunc(2021, 16, 1, "day16.txt", part1)
	challenges.RegisterChallengeFunc(2021, 16, 2, "day16.txt", part2)
}
