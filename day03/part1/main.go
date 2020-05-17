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
		fmt.Printf("unable to parse first wire path, %s\n", err.Error())
		return
	}

	// To check for intersections we will use cartesian coords with 0,0 being the central port (x increase to the right and y increasing up)
	// The Manhatten distance from the central port to any point should be |x| + |y|
	//
	// To calculate the minimum, first create a hash map of all the points crossed by the first wire
	// Then traverse the points crossed by the second wire. If any of them exist for the first wire this is considered an intersection (cross point) and recorde
	//
	// Once all points have been tranversed on the secomd wire, calculate the Manhatten distance for all cross points and report the smallest

	wire1Map := generatePointMap(wire1)
	wire2Map := generatePointMap(wire2)
	crossPoints, closestPoint := compareAndSortWireMaps(wire1Map, wire2Map)

	if len(crossPoints) == 0 {
		fmt.Println("unable to determine distance, no crossing point found")
		return
	}

	for _, crossPoint := range crossPoints {
		fmt.Printf("cross point at %v with dist %d\n", crossPoint, crossPoint.ManhattenDist())
	}
	fmt.Printf("\nclosest cross point at %v with dist %d\n", closestPoint, closestPoint.ManhattenDist())
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

type WireMap map[WirePoint]struct{}

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
	points := map[WirePoint]struct{}{
		currentPoint: {},
	}

	var addToX int
	var addToY int
	for _, entry := range path {
		addToX = 0
		addToY = 0
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
			currentPoint = WirePoint{x: currentPoint.x + addToX, y: currentPoint.y + addToY}
			points[currentPoint] = struct{}{}
		}
	}
	return points
}

func compareAndSortWireMaps(wire1Map WireMap, wire2Map WireMap) ([]WirePoint, WirePoint) {
	var crossPoints []WirePoint

	// Build up collection of cross points
	for point := range wire2Map {
		if _, ok := wire1Map[point]; ok && !point.IsOrigin() {
			crossPoints = append(crossPoints, point)
		}
	}

	sort.Slice(crossPoints, func(i, j int) bool {
		return crossPoints[i].ManhattenDist() < crossPoints[j].ManhattenDist()
	})

	closest := WirePoint{}
	if len(crossPoints) > 0 {
		closest = crossPoints[0]
	}
	return crossPoints, closest
}
