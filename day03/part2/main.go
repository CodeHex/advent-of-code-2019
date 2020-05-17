package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	// Read in wire data
	wire1, err := parseWirePath(scanner.Text())
	if err != nil {
		fmt.Printf("unable to parse first wire path, %s\n", err.Error())
		return
	}

	scanner.Scan()
	wire2, err := parseWirePath(scanner.Text())
	if err != nil {
		fmt.Printf("unable to parse second wire path, %s\n", err.Error())
		return
	}

	// To check for intersections we will use cartesian coords with 0,0 being the central port (x increase to the right and y increasing up)
	// The Manhatten distance from the central port to any point should be |x| + |y|
	//
	// To calculate the minimum, first create a hash map of all the points crossed by the first wire
	// ADDITION - Each point will record the minimum steps to get there when the map is being created
	// Then traverse the points crossed by the second wire. If any of them exist for the first wire this is considered an intersection (cross point) and recorde
	//
	// Once all points have been tranversed on the secomd wire, calculate the Manhatten distance for all cross points and report the smallest

	wire1Map := generatePointMap(wire1)
	wire2Map := generatePointMap(wire2)
	crossWires := crossPoints(wire1Map, wire2Map)

	if len(crossWires) == 0 {
		fmt.Println("unable to determine distance, no crossing point found")
		return
	}

	for crossPoint, steps := range crossWires {
		fmt.Printf("cross point at %v with dist %d and steps %d\n", crossPoint, crossPoint.ManhattenDist(), steps)
	}

	fmt.Println()
	manDistPoint := closestManhattenPoint(crossWires)
	fmt.Printf("closest Manhatten dist cross point at    %v with dist    %d\n", manDistPoint, manDistPoint.ManhattenDist())
	stepsPoint, steps := closestStepsPoint(crossWires)
	fmt.Printf("closest steps cross point at             %v with steps   %d\n", stepsPoint, steps)
}

type Direction string

const (
	Right Direction = "R"
	Left  Direction = "L"
	Up    Direction = "U"
	Down  Direction = "D"
)

var validDirections = []Direction{Up, Down, Right, Left}

type PathEntry struct {
	direction Direction
	distance  int
}

func parseWirePath(wireData string) ([]PathEntry, error) {
	parts := strings.Split(wireData, ",")
	if len(parts) == 0 {
		return nil, errors.New("invalid path, no data provided")
	}

	result := make([]PathEntry, len(parts))
	for i, part := range parts {
		result[i] = PathEntry{}
		for _, dir := range validDirections {
			if strings.HasPrefix(part, string(dir)) {
				result[i].direction = dir
				break
			}
		}
		if result[i].direction == "" {
			return nil, errors.Errorf("invalid path, unrecognized direction for entry '%s'", part)
		}

		dist, err := strconv.Atoi(strings.TrimPrefix(part, string(result[i].direction)))
		if err != nil {
			return nil, errors.Wrapf(err, "invalid path, unable to parse distance for entry '%s'", part)
		}
		result[i].distance = dist
	}
	return result, nil
}

type WirePoint struct {
	x int
	y int
}

type WireMap map[WirePoint]int

func (w WirePoint) IsOrigin() bool {
	return w.x == 0 && w.y == 0
}

func (w WirePoint) ManhattenDist() int {
	xMag := w.x
	if xMag < 0 {
		xMag = -w.x
	}

	yMag := w.y
	if yMag < 0 {
		yMag = -w.y
	}
	return xMag + yMag
}

func generatePointMap(path []PathEntry) WireMap {
	currentPoint := WirePoint{0, 0}
	points := map[WirePoint]int{
		currentPoint: 0,
	}

	var steps int
	for _, entry := range path {
		addToX := 0
		addToY := 0
		switch entry.direction {
		case Up:
			addToY = 1
		case Down:
			addToY = -1
		case Right:
			addToX = 1
		case Left:
			addToX = -1
		}

		for i := 0; i < entry.distance; i++ {
			steps += 1
			currentPoint = WirePoint{x: currentPoint.x + addToX, y: currentPoint.y + addToY}
			if _, ok := points[currentPoint]; !ok {
				points[currentPoint] = steps
			}
		}
	}
	return points
}

func crossPoints(wire1Map WireMap, wire2Map WireMap) map[WirePoint]int {
	crossPoints := make(map[WirePoint]int)

	// Build up collection of cross points with combined steps to get there
	for point, wire2Steps := range wire2Map {
		if wire1Steps, ok := wire1Map[point]; ok && !point.IsOrigin() {
			crossPoints[point] = wire1Steps + wire2Steps
		}
	}
	return crossPoints
}

func closestManhattenPoint(crossPoints map[WirePoint]int) WirePoint {
	points := make([]WirePoint, 0, len(crossPoints))
	for point := range crossPoints {
		points = append(points, point)
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].ManhattenDist() < points[j].ManhattenDist()
	})

	if len(crossPoints) > 0 {
		return points[0]
	}
	return WirePoint{}
}

func closestStepsPoint(crossPoints map[WirePoint]int) (WirePoint, int) {
	point := WirePoint{}
	steps := -1
	for crossPoint, crossSteps := range crossPoints {
		if steps == -1 || crossSteps < steps {
			point = crossPoint
			steps = crossSteps
		}
	}
	return point, steps
}
