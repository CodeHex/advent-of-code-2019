package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputData := scanner.Text()

	layers, err := parseLayers(inputData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	leastZeros := pixelHeight * pixelWidth
	var leastLayer *layer
	for i, l := range layers {
		fmt.Printf("Layer %d\n", i)
		l.print()
		fmt.Printf("Stats: %v\n\n", l.digitCount)

		if l.digitCount[0] < leastZeros {
			leastZeros = l.digitCount[0]
			leastLayer = l
		}
	}

	fmt.Println()
	leastLayer.print()
	fmt.Printf("Number of 1 digits * number of 2 digits = %d\n", leastLayer.digitCount[1]*leastLayer.digitCount[2])
}

const pixelWidth = 25
const pixelHeight = 6

type layer struct {
	pixels     [][]int
	digitCount map[int]int
}

func (l *layer) print() {
	for _, r := range l.pixels {
		fmt.Println(r)
	}
}

func makeLayer() *layer {
	p := make([][]int, pixelHeight)
	for i := range p {
		p[i] = make([]int, pixelWidth)
	}
	return &layer{pixels: p, digitCount: make(map[int]int)}
}

func parseLayers(inputData string) ([]*layer, error) {
	// Split into layers
	var result []*layer
	layerPixels := pixelWidth * pixelHeight
	for i := 0; i < len(inputData); i = i + (layerPixels) {
		l := makeLayer()
		layerData := inputData[i : i+layerPixels]

		// Split layer into rows
		for y := 0; y < pixelHeight; y++ {
			rowData := layerData[y*pixelWidth : (y+1)*pixelWidth]

			// Store each pixel
			for x := 0; x < pixelWidth; x++ {
				val, err := strconv.Atoi(string(rowData[x]))
				if err != nil {
					return nil, errors.Wrap(err, "unable to read int from input")
				}
				l.pixels[y][x] = val
				l.digitCount[val]++
			}
		}
		result = append(result, l)
	}
	return result, nil
}
