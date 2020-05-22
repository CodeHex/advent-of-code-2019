package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
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
	m := newAsteroidMap(inputData)
	m.shootAsteroids()
	m.print()
}

type coord struct {
	x, y int
}

type asteroidMap struct {
	asteroids [][]int
	station   coord
}

func newAsteroidMap(inputData []string) asteroidMap {
	chart := asteroidMap{}
	for i, row := range inputData {
		chart.asteroids = append(chart.asteroids, nil)
		for _, val := range row {
			isAsteroid := val == '#'
			entry := -1
			if isAsteroid {
				entry = 0
			}
			chart.asteroids[i] = append(chart.asteroids[i], entry)
		}
	}
	chart.analyze()
	chart.station = chart.highestAsteroidValue()
	return chart
}

func (a asteroidMap) set(c coord, val int) {
	a.asteroids[c.y][c.x] = val
}

func (a asteroidMap) get(c coord) int {
	return a.asteroids[c.y][c.x]
}

func (a asteroidMap) print() {
	const reset = "\033[0m"
	const green = "\033[32m"

	for y, row := range a.asteroids {
		fmt.Printf("[")
		for x, val := range row {
			switch {
			case val == -1:
				fmt.Printf("   ")
			case x == a.station.x && y == a.station.y:
				fmt.Printf("%s  X%s", green, reset)
			default:
				fmt.Printf("%3d", val)
			}
			if x != len(row)-1 {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("]\n")
	}
	fmt.Printf("\nStation is at (%d,%d)\n", a.station.x, a.station.y)
}

// replaces asteroid values with the number of asteroids that are in LOS and
// calculates where the station should be
func (a asteroidMap) analyze() {
	for y, row := range a.asteroids {
		for x, val := range row {
			if val != -1 {
				c := coord{x: x, y: y}
				a.set(c, a.calcLOS(c))
			}
		}
	}
	a.station = a.highestAsteroidValue()
}

// calculate the coord with the highest current asteroid value
func (a asteroidMap) highestAsteroidValue() coord {
	most := 0
	var mostX, mostY int
	for y, row := range a.asteroids {
		for x, val := range row {
			if most < val {
				mostX, mostY = x, y
				most = val
			}
		}
	}
	return coord{x: mostX, y: mostY}
}

type gradient struct {
	// The gradient of the line
	m float32

	// Is the angle from the vertical line straight up (negative y). When shooting, this angle
	// will rotate from 0 -> 2 pi
	angleRads float64

	// whether you
	dir bool
}

// calculateLOS counts how many target asteroids can be seen by the provided asteroid.
// Each target asteroid is projected on to a line equation (one for forward direction and one for
// backwards. Unique equations are recorded as only one
// asteroid can be seen on that line
func (a asteroidMap) calcLOS(mapCoord coord) int {
	// treat provided asteroid as 0,0, so all line equations are y = mx
	xMod := -mapCoord.x
	yMod := -mapCoord.y

	grads := make(map[gradient]struct{})
	for j, row := range a.asteroids {
		for i, val := range row {
			sourceCoord := coord{x: i + xMod, y: j + yMod}
			// ignore empty space and source asteroid
			if val == -1 || (sourceCoord.x == 0 && sourceCoord.y == 0) {
				continue
			}

			grads[calcGrad(sourceCoord)] = struct{}{}
		}
	}
	return len(grads)
}

type asteroidPos struct {
	mapCoord     coord
	stationCoord coord
}

// contains all asteroids coordinates keyed by line gradient for line equations from the station
type vectors map[gradient][]asteroidPos

// calculates all asteroids in the form of vectors from the statiomn
func (a asteroidMap) newVectors() vectors {
	// treat station as 0,0
	xMod := -a.station.x
	yMod := -a.station.y

	vectors := make(map[gradient][]asteroidPos)

	for j, row := range a.asteroids {
		for i, val := range row {
			mapCoord := coord{x: i, y: j}
			stationCoord := coord{x: i + xMod, y: j + yMod}
			target := asteroidPos{mapCoord: mapCoord, stationCoord: stationCoord}
			// Ignore empty space and station
			if val == -1 || (stationCoord.x == 0 && stationCoord.y == 0) {
				continue
			}
			// Find line gradient for the targe asteroid
			grad := calcGrad(stationCoord)
			vectors[grad] = append(vectors[grad], target)
		}
	}

	// Sort targets closet to fartherest
	for grad := range vectors {
		sort.Slice(vectors[grad], func(i, j int) bool {
			x1, y1 := vectors[grad][i].stationCoord.x, vectors[grad][i].stationCoord.y
			x2, y2 := vectors[grad][j].stationCoord.x, vectors[grad][j].stationCoord.y
			return (x1*x1)+(y1*y1) < (x2*x2)+(y2*y2)
		})
	}
	return vectors
}

// shoot asteroids by rotating the laser continuously until it does a full sweep without shooting anything
func (a asteroidMap) shootAsteroids() {
	vectors := a.newVectors()

	// Order the gradients by angle to "rotate" around all the angles
	var rotateGrads []gradient
	for grad := range vectors {
		rotateGrads = append(rotateGrads, grad)
	}
	sort.Slice(rotateGrads, func(i, j int) bool {
		return rotateGrads[i].angleRads < rotateGrads[j].angleRads
	})

	noShotsFired := false
	shotCounter := 0
	for !noShotsFired {
		noShotsFired = true
		for _, grad := range rotateGrads {
			// Move to next line if no more asteroids left on that line
			if len(vectors[grad]) == 0 {
				continue
			}
			// Shoot asteroid
			noShotsFired = false
			shotCounter++
			a.set(vectors[grad][0].mapCoord, shotCounter)
			fmt.Printf("asteroid shot %d  hits asteroid %v\n", shotCounter, vectors[grad][0].mapCoord)
			vectors[grad] = vectors[grad][1:]
		}
	}
}

// calculates the gradient from the provided coord and 0,0. The line is split so any line with angle 0 > pi
// is considered forward, and pi > 2 pi considered backwards
func calcGrad(c coord) gradient {
	g := gradient{
		m:         float32(c.y) / float32(c.x),
		angleRads: math.Atan(-1 * float64(c.x) / float64(c.y)),
	}

	switch {
	case c.x == 0:
		g.dir = c.y < 0
	default:
		g.dir = c.x > 0
	}

	if g.angleRads < 0 {
		g.angleRads += math.Pi
	}
	if !g.dir {
		g.angleRads += math.Pi
	}
	return g
}
