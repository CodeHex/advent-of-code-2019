package main

import (
	"bufio"
	"fmt"
	"os"
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
	m.analyze()
	m.print()
}

type asteroidMap [][]int

func newAsteroidMap(inputData []string) asteroidMap {
	var chart asteroidMap

	for i, row := range inputData {
		chart = append(chart, nil)
		for _, val := range row {
			isAsteroid := val == '#'
			entry := -1
			if isAsteroid {
				entry = 0
			}
			chart[i] = append(chart[i], entry)
		}
	}
	return chart
}

func (a asteroidMap) print() {
	const reset = "\033[0m"
	const green = "\033[32m"

	mostX, mostY := a.mostAsteroids()
	for y, row := range a {
		fmt.Printf("[")
		for x, val := range row {
			switch {
			case val == -1:
				fmt.Printf("     ")
			case x == mostX && y == mostY:
				fmt.Printf("%s%3d%s, ", green, val, reset)
			default:
				fmt.Printf("%3d, ", val)
			}
		}
		fmt.Printf("]\n")
	}
}

func (a asteroidMap) analyze() {
	for y, row := range a {
		for x, val := range row {
			if val != -1 {
				a[y][x] = a.calculateLOS(x, y)
			}
		}
	}
}

func (a asteroidMap) mostAsteroids() (x, y int) {
	most := 0
	var mostX, mostY int
	for y, row := range a {
		for x, val := range row {
			if most < val {
				mostX = x
				mostY = y
				most = val
			}
		}
	}
	return mostX, mostY
}

type grad struct {
	mTop    int
	mBottom int
	dir     bool
}

func (a asteroidMap) calculateLOS(x, y int) int {
	// treat asteriod as 0,0
	xMod := -x
	yMod := -y

	grads := make(map[grad]struct{})
	for j, row := range a {
		for i, val := range row {
			if val == -1 || (i == x && j == y) {
				continue
			}

			// Add it to an equation list of lines
			mTop := j + yMod
			mBottom := i + xMod
			factor := GCD(mTop, mBottom)
			if factor == 0 {
				factor = 1
			}
			if factor < 0 {
				factor = factor * -1
			}
			m := grad{mTop / factor, mBottom / factor, mTop > 0 || mBottom > 0}
			//fmt.Printf("--- i %d, j %d, mTop %d, mBottom %d, factor %d\n", i, j, mTop, mBottom, factor)
			grads[m] = struct{}{}
		}
	}
	//fmt.Printf("asteroid [%d,%d]\n", x, y)
	//fmt.Printf("%v\n", grads)
	//fmt.Println()
	return len(grads)
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
