package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var inputText string

	if len(os.Args[1:]) == 1 {
		inputBytes, err := ioutil.ReadFile(os.Args[1:][0])
		if err != nil {
			fmt.Printf("unable to read input file, %s\n", err.Error())
			os.Exit(1)
		}
		inputText = string(inputBytes)
	} else {
		fmt.Println("ENTER INT CODE")
		scanner.Scan()
		inputText = scanner.Text()
	}

	clearScreen()
	fmt.Println("mapping out corridors...")
	field := generateField(inputText)
	fmt.Println("releasing oxygen in 5 seconds...")
	time.Sleep(5 * time.Second)
	field.releaseOxygen()
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

type Move int

const (
	MoveNorth Move = 1
	MoveSouth Move = 2
	MoveWest  Move = 3
	MoveEast  Move = 4
)

func (m Move) String() string {
	switch m {
	case MoveNorth:
		return "N"
	case MoveSouth:
		return "S"
	case MoveEast:
		return "E"
	case MoveWest:
		return "W"
	default:
		return strconv.Itoa(int(m))
	}
}

var reverseMove = map[Move]Move{
	MoveNorth: MoveSouth,
	MoveSouth: MoveNorth,
	MoveWest:  MoveEast,
	MoveEast:  MoveWest,
}

type Out int

const (
	OutWall   Out = 0
	OutOK     Out = 1
	OutOxygen Out = 2
)

type point struct {
	x, y int
}

type field struct {
	visited            map[point]struct{}
	unvisited          []point
	distance           map[point]int
	prevNode           map[point]*point
	prevNodeReverseDir map[point]Move
	walls              map[point]struct{}
	oxygenTank         *point
	vacuum             map[point]struct{}
	oxygenated         map[point]struct{}
}

func newField() *field {
	f := &field{
		visited:            make(map[point]struct{}),
		unvisited:          []point{{0, 0}},
		distance:           make(map[point]int),
		prevNode:           make(map[point]*point),
		prevNodeReverseDir: make(map[point]Move),
		walls:              make(map[point]struct{}),
		vacuum:             make(map[point]struct{}),
		oxygenated:         make(map[point]struct{}),
	}
	f.distance[point{0, 0}] = 0
	return f
}

func nearestPoints(p point) map[Move]point {
	return map[Move]point{
		MoveNorth: point{p.x, p.y - 1},
		MoveSouth: point{p.x, p.y + 1},
		MoveEast:  point{p.x + 1, p.y},
		MoveWest:  point{p.x - 1, p.y},
	}
}

func generateField(inputProgram string) *field {
	result := newField()

	// Keep looping until all nodes have been visited
	counter := 0
	for len(result.unvisited) != 0 {
		counter++
		// Get next node to process
		node := result.unvisited[0]
		result.unvisited = result.unvisited[1:]
		result.visited[node] = struct{}{}

		// Work out directions to node
		path := result.calculatePath(node)
		//clearScreen()
		//fmt.Printf("visiting %v\n", node)
		//result.print()

		// Get a computer and move to the point
		c, _ := newInOutComputer(inputProgram)
		for _, m := range path {
			output := Out(c.Input(int(m)))
			if output == OutWall {
				panic("unexpected wall")
			}
		}

		// Work out the nearest points (and the directions to them)
		nearest := nearestPoints(node)
		for dir, np := range nearest {
			output := Out(c.Input(int(dir)))
			switch output {
			case OutWall:
				result.walls[np] = struct{}{}
			case OutOK, OutOxygen:
				// Reverse back
				revOut := Out(c.Input(int(reverseMove[dir])))
				if revOut == OutWall {
					panic("unexpected wall when reversing")
				}
				if _, ok := result.visited[np]; !ok {
					toAdd := true
					for _, u := range result.unvisited {
						if u == np {
							toAdd = false
							break
						}
					}
					if toAdd {
						result.unvisited = append(result.unvisited, np)
					}
				}

				dist := result.distance[node] + 1
				val, ok := result.distance[np]
				if !ok || dist < val {
					nodecpy := node
					result.distance[np] = dist
					result.prevNode[np] = &nodecpy
					result.prevNodeReverseDir[np] = dir
				}
				if output == OutOxygen {
					npcpy := np
					result.oxygenTank = &npcpy
				}
			}
		}
		if counter%10 == 0 || len(result.unvisited) == 0 {
			clearScreen()
			fmt.Printf("analysed %d nodes\n\n", counter)
			result.print()
			time.Sleep(200 * time.Millisecond)
		}

	}
	return result
}

func (f *field) calculatePath(node point) []Move {
	var backwardsPath []Move
	prevNodeDir := f.prevNodeReverseDir[node]
	prevNode := f.prevNode[node]
	for prevNode != nil {
		backwardsPath = append(backwardsPath, prevNodeDir)
		prevNodeDir = f.prevNodeReverseDir[*prevNode]
		prevNode = f.prevNode[*prevNode]
	}
	var path []Move
	for i := len(backwardsPath) - 1; i >= 0; i-- {
		path = append(path, backwardsPath[i])
	}
	return path
}

func (f *field) print() {
	f.printDetailed(true, false)
}

func (f *field) printDetailed(showPath bool, showOxygen bool) {
	allPoints := make(map[point]struct{})
	for visitedPoint := range f.visited {
		allPoints[visitedPoint] = struct{}{}
	}
	for _, unvisitedPoint := range f.unvisited {
		allPoints[unvisitedPoint] = struct{}{}
	}
	for wall := range f.walls {
		allPoints[wall] = struct{}{}
	}
	if f.oxygenTank != nil {
		allPoints[*f.oxygenTank] = struct{}{}
	}

	minX, minY, maxX, maxY := 0, 0, 0, 0
	for p := range allPoints {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	// Make array
	var canvas [][]string
	height := maxY - minY + 1
	width := maxX - minX + 1

	for j := 0; j < height; j++ {
		canvas = append(canvas, make([]string, width))
	}

	for visitedPoint := range f.visited {
		canvas[visitedPoint.y-minY][visitedPoint.x-minX] = "."
	}
	for _, unvisitedPoint := range f.unvisited {
		canvas[unvisitedPoint.y-minY][unvisitedPoint.x-minX] = "_"
	}
	for wall := range f.walls {
		canvas[wall.y-minY][wall.x-minX] = "▒"
	}
	if showOxygen {
		for p := range f.oxygenated {
			canvas[p.y-minY][p.x-minX] = "0"
		}
	}
	if f.oxygenTank != nil {
		canvas[f.oxygenTank.y-minY][f.oxygenTank.x-minX] = "O"
		if showPath {
			node := f.prevNode[*f.oxygenTank]
			for node != nil {
				canvas[node.y-minY][node.x-minX] = "●"
				node = f.prevNode[*node]
			}
		}
	}
	canvas[-minY][-minX] = "X"

	for j := 0; j < len(canvas); j++ {
		line := fmt.Sprintf("%1v", canvas[j])
		fmt.Println(line[1 : len(line)-1])
	}
	if f.oxygenTank != nil && showPath {
		fmt.Printf("\ncurrent shortest distance is %d\n", f.distance[*f.oxygenTank])
	}
	fmt.Println()
}

func (f *field) releaseOxygen() {
	mins := 0
	for p := range f.visited {
		f.vacuum[p] = struct{}{}
	}
	f.oxygenated[*f.oxygenTank] = struct{}{}

	counter := 0
	for len(f.vacuum) != 0 {
		counter++
		var oxyToAdd []point
		for o := range f.oxygenated {
			nps := nearestPoints(o)
			for _, np := range nps {
				if _, ok := f.walls[np]; ok {
					continue
				}
				if _, ok := f.vacuum[np]; ok {
					oxyToAdd = append(oxyToAdd, np)
				}
			}
		}

		for _, oxy := range oxyToAdd {
			f.oxygenated[oxy] = struct{}{}
			delete(f.vacuum, oxy)
		}
		mins++

		if counter%10 == 0 || len(f.vacuum) == 0 {
			clearScreen()
			fmt.Printf("after %d mins\n\n", mins)
			f.printDetailed(false, true)
			time.Sleep(200 * time.Millisecond)
		}
	}
}
