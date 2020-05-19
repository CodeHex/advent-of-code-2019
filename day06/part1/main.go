package main

import (
	"bufio"
	"fmt"
	"os"
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

	// Total orbits over all planets in the map
	totalOrbits := 0
	orbitMap := newOrbitMap(inputData)
	for planet := range orbitMap {
		totalOrbits += orbitMap.calculateOrbits(planet)
	}

	fmt.Printf("Total orbits: %d\n", totalOrbits)
	fmt.Printf("Shortest path between YOU and SAN is %d hops\n", orbitMap.shortestHops())
}

// Holds the planets and what they directly orbit
type planetMap map[string]string

func newOrbitMap(inputLines []string) planetMap {
	m := make(planetMap)
	for _, line := range inputLines {
		parts := strings.Split(line, ")")
		m[parts[1]] = parts[0]
	}
	return m
}

// orbits are calculated by adding the direct orbits, plus
// all direct and indirect orbits of inner planets
func (p planetMap) calculateOrbits(planet string) int {
	inner, ok := p[planet]
	if !ok {
		return 0
	}
	return p.calculateOrbits(inner) + 1
}

func (p planetMap) generatePath(entry string) []string {
	var path []string
	ok := true
	for ok {
		path = append(path, entry)
		entry, ok = p[entry]
	}

	reversedPath := make([]string, len(path))
	for i, entry := range path {
		reversedPath[len(path)-1-i] = entry
	}
	return reversedPath
}

func (p planetMap) shortestHops() int {
	pathToSAN := p.generatePath("SAN")
	pathToYOU := p.generatePath("YOU")

	// Find first index where prev index differs (common ancestor)
	for i := 0; i < len(pathToSAN); i++ {
		if pathToSAN[i] != pathToYOU[i] {
			// The i -1 must be the common node, so total orbit hop is
			// hops from SANS -> COMMON -> YOU
			return (len(pathToSAN) - i - 1) + (len(pathToYOU) - i - 1)
		}
	}
	return 0
}
