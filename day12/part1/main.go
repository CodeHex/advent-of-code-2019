package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var inputData []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		inputData = append(inputData, line)
	}

	zeroSystem := newSystem(inputData)
	timeline := generateTimeline(zeroSystem, 1000)
	timeline.print(1000)
}

type systemTimeline []system

type system struct {
	moons    []*moon
	timeStep int
}

type vector struct{ x, y, z int }

type moon struct {
	p vector
	v vector
}

func (m moon) pe() int          { return mod(m.p.x) + mod(m.p.y) + mod(m.p.z) }
func (m moon) ke() int          { return mod(m.v.x) + mod(m.v.y) + mod(m.v.z) }
func (m moon) totalEnergy() int { return m.pe() * m.ke() }

func (s system) totalEnergy() int { return s.sum(func(m *moon) int { return m.totalEnergy() }) }

func newSystem(inputData []string) system {
	moons := make([]*moon, len(inputData))
	for i, moonData := range inputData {
		moons[i] = newMoon(moonData)
	}
	return system{moons: moons}
}

func newMoon(inputData string) *moon {
	inputData = strings.Trim(inputData, "<>")
	inputData = strings.ReplaceAll(inputData, " ", "")
	coords := strings.Split(inputData, ",")
	points := make(map[string]int)
	for _, entry := range coords {
		s := strings.Split(entry, "=")
		points[s[0]], _ = strconv.Atoi(s[1])
	}
	return &moon{p: vector{x: points["x"], y: points["y"], z: points["z"]}}
}

func (v vector) string() string {
	return fmt.Sprintf("<x=%3d, y=%3d, z=%3d>", v.x, v.y, v.z)
}

func (m moon) print() {
	fmt.Printf("pos=%s, vel=%s\n", m.p.string(), m.v.string())
}

func (s system) print() {
	fmt.Printf("After %d steps:\n", s.timeStep)
	for _, moon := range s.moons {
		moon.print()
	}
	fmt.Printf("Sum of total energy: %d\n", s.totalEnergy())
}

func (s systemTimeline) print(steps ...int) {
	if len(steps) == 0 {
		for _, sys := range s {
			sys.print()
			fmt.Println()
		}
		return
	}

	for _, step := range steps {
		s[step].print()
		fmt.Println()
	}
}

func (s system) advance() system {
	nextMoons := make([]*moon, len(s.moons))

	for i, currentMoon := range s.moons {
		newVelocity := gravity(currentMoon, s.moons)
		nextMoon := &moon{v: newVelocity, p: addVector(currentMoon.p, newVelocity)}
		nextMoons[i] = nextMoon
	}

	return system{
		moons:    nextMoons,
		timeStep: s.timeStep + 1,
	}
}

func gravity(m *moon, allMoons []*moon) vector {
	newVector := m.v
	for _, moonPair := range allMoons {
		if moonPair == m {
			continue
		}
		velAdjustVector := vector{
			x: velocityAdjustment(m.p.x, moonPair.p.x),
			y: velocityAdjustment(m.p.y, moonPair.p.y),
			z: velocityAdjustment(m.p.z, moonPair.p.z),
		}

		newVector = addVector(newVector, velAdjustVector)
	}
	return newVector
}

func velocityAdjustment(moonCoord int, pairMoonCoord int) int {
	switch {
	case moonCoord > pairMoonCoord:
		return -1
	case moonCoord < pairMoonCoord:
		return 1
	default:
		return 0
	}
}

func generateTimeline(zeroPoint system, steps int) systemTimeline {
	result := make([]system, steps+1)
	currentSystem := zeroPoint
	result[0] = currentSystem
	for i := 1; i < steps+1; i++ {
		currentSystem = currentSystem.advance()
		result[i] = currentSystem
	}
	return result
}

func addVector(v1 vector, v2 vector) vector {
	return vector{
		x: v1.x + v2.x,
		y: v1.y + v2.y,
		z: v1.z + v2.z,
	}
}

func (s system) sum(quantity func(*moon) int) int {
	total := 0
	for _, moon := range s.moons {
		total += quantity(moon)
	}
	return total
}

func mod(i int) int {
	if i > 0 {
		return i
	}
	return -1 * i
}
