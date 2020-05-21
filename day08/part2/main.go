package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	compressed := compressLayers(layers)
	compressed.print()
}

const pixelWidth = 25
const pixelHeight = 6

const transparentPixel = 2

type layer struct {
	pixels     [][]int
	digitCount map[int]int
}

func (l *layer) print() {
	for _, r := range l.pixels {
		lineStr := fmt.Sprintf("%v", r)
		lineStr = strings.ReplaceAll(lineStr, "0", " ")
		lineStr = strings.ReplaceAll(lineStr, "1", "*")
		fmt.Println(lineStr)
	}
}

func makeLayer() *layer {
	p := make([][]int, pixelHeight)
	for i := range p {
		p[i] = make([]int, pixelWidth)
	}
	return &layer{pixels: p, digitCount: make(map[int]int)}
}

func makeTransparentLayer() *layer {
	p := makeLayer()
	for i, row := range p.pixels {
		for j := range row {
			p.pixels[i][j] = transparentPixel
		}
	}
	return p
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

func compressLayers(layers []*layer) *layer {
	result := makeTransparentLayer()
	for _, l := range layers {
		combineLayers(result, l)
	}
	return result
}

func combineLayers(result *layer, l *layer) {
	for i, row := range l.pixels {
		for j, pixVal := range row {
			// If the top layer is transparent, compress to lower layer, other keep original pixel
			if result.pixels[i][j] == transparentPixel {
				result.pixels[i][j] = pixVal
			}
		}
	}
}
