package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

	in := make(chan int64)
	out := make(chan int64)
	c, err := newChannelComputer(inputText, in, out)
	if err != nil {
		fmt.Printf("unable to convert input to int code memory, %s\n", err.Error())
		os.Exit(1)
	}
	c.disableLog = true
	c.disableOutLog = true
	c.memory[0] = 2

	// Start computer running
	end := make(chan struct{})
	go func() {
		err = c.run()
		if err != nil {
			fmt.Printf("error running program: %s\n", err.Error())
			os.Exit(1)
		}
		close(end)
	}()

	paddleGame := newGame()
	paddleGame.runGameloop(in, out, end)
}

type JoyPos int64

const (
	JoyLeft   JoyPos = -1
	JoyCenter JoyPos = 0
	JoyRight  JoyPos = 1
)

type game struct {
	board     [][]TileType
	score     int
	ballPos   point
	paddlePos point
}

func newGame() *game {
	return &game{board: [][]TileType{}, score: 0}
}

const scoreOutput = -1

func (g *game) runGameloop(in chan<- int64, out <-chan int64, end <-chan struct{}) {
	for {
		select {
		case in <- int64(g.calcOptimalJoystickPos()):
			g.refreshScreen()
		case startVal := <-out:
			if startVal == scoreOutput {
				<-out //ignore second arg as the we only care about the score
				g.score = int(<-out)
				break
			}
			p := point{x: int(startVal), y: int(<-out)}
			tileType := TileType(<-out)
			g.setTile(p, tileType)
		case <-end:
			g.refreshScreen()
			return
		}
	}
}

func (g *game) setTile(p point, tile TileType) {
	for p.y >= len(g.board) {
		g.board = append(g.board, []TileType{})
	}
	for p.x >= len(g.board[p.y]) {
		g.board[p.y] = append(g.board[p.y], 0)
	}
	g.board[p.y][p.x] = tile

	if tile == TileBall {
		g.ballPos = p
	}
	if tile == TileHorizontal {
		g.paddlePos = p
	}
}

func (g *game) calcOptimalJoystickPos() JoyPos {
	switch {
	case g.ballPos.x > g.paddlePos.x:
		return JoyRight
	case g.ballPos.x < g.paddlePos.x:
		return JoyLeft
	default:
		return JoyCenter
	}
}

func (g *game) refreshScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Printf("SCORE: %d\n", g.score)
	for _, row := range g.board {
		line := ""
		for _, tile := range row {
			line += tile.icon()
		}
		fmt.Println(line)
	}

	// Keep screen visible for at least this time
	time.Sleep(30 * time.Millisecond)
}

type point struct {
	x, y int
}

type TileType int

const (
	TileEmpty      TileType = 0
	TileWall       TileType = 1
	TileBlock      TileType = 2
	TileHorizontal TileType = 3
	TileBall       TileType = 4
)

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

func (t TileType) icon() string {
	switch t {
	case TileEmpty:
		return " "
	case TileWall:
		return "▒"
	case TileBlock:
		return "□"
	case TileHorizontal:
		return "-"
	case TileBall:
		return "●"
	default:
		panic("unknown tile type")
	}
}
