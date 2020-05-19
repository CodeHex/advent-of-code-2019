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
}

// Holds the planets and what they directly orbit
type planetMap map[string][]string

func newOrbitMap(inputLines []string) planetMap {
	m := make(planetMap)
	for _, line := range inputLines {
		parts := strings.Split(line, ")")
		m[parts[0]] = append(m[parts[0]], parts[1])
	}
	return m
}

// orbits are calculated by adding the direct orbits, plus
// all direct and indirect orbits of inner planets
func (p planetMap) calculateOrbits(planet string) int {
	// Start with direct orbits
	total := len(p[planet])

	// Add all other orbits of inner planets
	for _, innerPlanet := range p[planet] {
		total += p.calculateOrbits(innerPlanet)
	}
	return total
}
