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

	sys := newSystem(inputData)
	cmpSys := newSystem(inputData)

	nextSystemRepeat(sys, cmpSys)
}

type system struct {
	moons    []*moon
	timeStep int64
}

type vector struct{ x, y, z int }

type moon struct {
	p vector
	v vector
}

func newSystem(inputData []string) *system {
	moons := make([]*moon, len(inputData))
	for i, moonData := range inputData {
		moons[i] = newMoon(moonData)
	}
	return &system{moons: moons}
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

func (v vector) equals(cmp vector) bool {
	return v.x == cmp.x && v.y == cmp.y && v.z == cmp.z
}

func (m *moon) print() {
	fmt.Printf("pos=%s, vel=%s\n", m.p.string(), m.v.string())
}

func (s *system) print() {
	fmt.Printf("After %d steps:\n", s.timeStep)
	for _, moon := range s.moons {
		moon.print()
	}
}

func (s *system) equals(cmp *system) bool {
	for i, moon := range s.moons {
		if !cmp.moons[i].v.equals(moon.v) || !cmp.moons[i].p.equals(moon.p) {
			return false
		}
	}
	return true
}

func (s *system) advance() {
	s.timeStep++
	applyGravity(s.moons[0], s.moons[1], s.moons[2], s.moons[3])
	for i, currentMoon := range s.moons {
		s.moons[i].p = addVector(currentMoon.p, s.moons[i].v)
	}
}

func applyGravity(m1, m2, m3, m4 *moon) {
	updatePair(m1, m2)
	updatePair(m1, m3)
	updatePair(m1, m4)
	updatePair(m2, m3)
	updatePair(m2, m4)
	updatePair(m3, m4)
}

func updatePair(m1 *moon, m2 *moon) {
	switch {
	case m1.p.x > m2.p.x:
		m1.v.x--
		m2.v.x++
	case m1.p.x < m2.p.x:
		m1.v.x++
		m2.v.x--
	default:
	}

	switch {
	case m1.p.y > m2.p.y:
		m1.v.y--
		m2.v.y++
	case m1.p.y < m2.p.y:
		m1.v.y++
		m2.v.y--
	default:
	}

	switch {
	case m1.p.z > m2.p.z:
		m1.v.z--
		m2.v.z++
	case m1.p.z < m2.p.z:
		m1.v.z++
		m2.v.z--
	default:
	}
}

// work out when each dimension repeats and then find the least common multiple to work out when
// then will all be zero
func nextSystemRepeat(s *system, cmp *system) {
	s.advance()
	xRepeat := int64(0)
	yRepeat := int64(0)
	zRepeat := int64(0)
	for xRepeat == 0 || yRepeat == 0 || zRepeat == 0 {
		s.advance()
		if xRepeat == 0 && checkSystemDimension(s, cmp, func(m *moon) int { return m.p.x }, func(m *moon) int { return m.v.x }) {
			xRepeat = s.timeStep
		}
		if yRepeat == 0 && checkSystemDimension(s, cmp, func(m *moon) int { return m.p.y }, func(m *moon) int { return m.v.y }) {
			yRepeat = s.timeStep
		}
		if zRepeat == 0 && checkSystemDimension(s, cmp, func(m *moon) int { return m.p.z }, func(m *moon) int { return m.v.z }) {
			zRepeat = s.timeStep
		}
	}
	fmt.Printf("Repeats at x: %d, y: %d, z: %d (lcm: %d)\n", xRepeat, yRepeat, zRepeat, lcm(xRepeat, yRepeat, zRepeat))

}

func checkSystemDimension(s *system, cmp *system, getP func(*moon) int, getV func(*moon) int) bool {
	for i, moon := range s.moons {
		if getV(moon) != getV(cmp.moons[i]) || getP(moon) != getP(cmp.moons[i]) {
			return false
		}
	}
	return true
}

func addVector(v1 vector, v2 vector) vector {
	return vector{
		x: v1.x + v2.x,
		y: v1.y + v2.y,
		z: v1.z + v2.z,
	}
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int64, integers ...int64) int64 {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
