package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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

	out := make(chan int64)
	c, err := newChannelComputer(inputText, nil, out)
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

	// Keep a list of tile types with points
	records := make(map[TileType][]point)

	completed := false
	for !completed {
		select {
		case x := <-out:
			y := <-out
			tileType := TileType(<-out)
			records[tileType] = append(records[tileType], point{x: x, y: y})
		case <-endCh:
			completed = true
		}
	}

	for t, points := range records {
		fmt.Printf("%s tile (occurs %d times)\n", t.string(), len(points))
	}
}

type point struct {
	x, y     int64
	tileType TileType
}

type TileType int

func (t TileType) string() string {
	switch t {
	case TileEmpty:
		return "Empty"
	case TileWall:
		return "Wall"
	case TileBlock:
		return "Block"
	case TileHorizontal:
		return "Horizontal"
	case TileBall:
		return "Ball"
	default:
		panic("unknown tile type")
	}
}

const (
	TileEmpty      TileType = 0
	TileWall       TileType = 1
	TileBlock      TileType = 2
	TileHorizontal TileType = 3
	TileBall       TileType = 4
)
