package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

	in := make(chan int64)
	out := make(chan int64)
	c, err := newChannelComputer(inputText, in, out)
	if err != nil {
		fmt.Printf("unable to convert input to int code memory, %s\n", err.Error())
		os.Exit(1)
	}
	c.disableLog = true

	// Start computer running
	endCh := make(chan struct{})
	go func() {
		err = c.run()
		if err != nil {
			fmt.Printf("error running program: %s\n", err.Error())
			os.Exit(1)
		}
		close(endCh)
	}()

	const (
		Up    = 1
		Down  = 2
		Left  = 3
		Right = 4
	)
	turnRight := map[int]int{
		Up:    Right,
		Right: Down,
		Down:  Left,
		Left:  Up,
	}
	turnLeft := map[int]int{
		Up:    Left,
		Left:  Down,
		Down:  Right,
		Right: Up,
	}

	moveForward := map[int]func(panel) panel{
		Up:    func(p panel) panel { return panel{p.x, p.y - 1} },
		Down:  func(p panel) panel { return panel{p.x, p.y + 1} },
		Left:  func(p panel) panel { return panel{p.x - 1, p.y} },
		Right: func(p panel) panel { return panel{p.x + 1, p.y} },
	}

	panels := make(map[panel]int64)
	cursor := panel{0, 0}
	direction := Up
	var currentPanel int64 = 1

	terminated := false
	for !terminated {
		select {

		case in <- currentPanel:
			// Feed current panel color
		case color := <-out:
			// If we have an output paint the panel
			panels[cursor] = color
			move := <-out
			if move == 0 {
				direction = turnLeft[direction]
			} else {
				direction = turnRight[direction]
			}
			cursor = moveForward[direction](cursor)
			currentPanel = panels[cursor]
		case <-endCh:
			terminated = true
		}
	}

	print(panels)
}

type panel struct {
	x, y int
}

func print(panels map[panel]int64) {
	minX, minY, maxX, maxY := 0, 0, 0, 0
	for panel := range panels {
		if panel.x < minX {
			minX = panel.x
		}
		if panel.x > maxX {
			maxX = panel.x
		}
		if panel.y < minY {
			minY = panel.y
		}
		if panel.y > maxY {
			maxY = panel.y
		}
	}

	// Make array
	var canvas [][]int64
	height := maxY - minY + 1
	width := maxX - minX + 1

	for j := 0; j < height; j++ {
		canvas = append(canvas, make([]int64, width))
	}

	for panel, val := range panels {
		canvas[panel.y-minY][panel.x-minX] = val
	}

	// Print array
	for j := 0; j < len(canvas); j++ {
		line := fmt.Sprintf("%v", canvas[j])
		line = strings.ReplaceAll(line, "0", " ")
		line = strings.ReplaceAll(line, "1", "*")
		fmt.Println(line)
	}
}
