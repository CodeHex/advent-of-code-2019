package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
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

	clearScreen()
	comp, _ := newComputer(inputText, nil)
	comp.disableLog = true
	comp.disableOutLog = true

	comp.run()

	state := newMapState(comp.outputs)
	state.print()
	alignment := state.calcAlignment()
	fmt.Printf("Alignment: %d\n", alignment)

	path := state.calcPath()
	fmt.Printf("Path: %s\n", path)

	// Output is

	// L,6,R,12,L,6,L,8,L,8,                 A
	// L,6,R,12,L,6,L,8,L,8,                 A
	// L,6,R,12,R,8,L,8,                     B
	// L,4,L,4,L,6,                          C
	// L,6,R,12,R,8,L,8,                     B
	// L,6,R,12,L,6,L,8,L,8,                 A
	// L,4,L,4,L,6,                          C
	// L,6,R,12,R,8,L,8,                     B
	// L,4,L,4,L,6,                          C
	// L,6,R,12,L,6,L,8,L,8                  A

	newInput := "A,A,B,C,B,A,C,B,C,A\nL,6,R,12,L,6,L,8,L,8\nL,6,R,12,R,8,L,8\nL,4,L,4,L,6\nn\n"

	newInputMemory := "2" + inputText[1:]
	inComp, _ := newInOutComputer(newInputMemory)

	go func() {
		for _, v := range newInput {
			inComp.in <- int64(v)
		}
	}()

	var lastval int64
	for a := range inComp.out {
		lastval = a
	}
	fmt.Printf("DUST COLLECTED: %d\n", lastval)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

type mapState struct {
	data           [][]byte
	robotX, robotY int
	depth          int
	length         int
}

func newMapState(out []int64) *mapState {
	byteLine := make([]byte, len(out))
	for i, v := range out {
		byteLine[i] = byte(v)
	}
	outString := string(byteLine)

	lines := strings.Split(outString, "\n")
	lines = lines[:len(lines)-2] // Remove last empty line
	depth := len(lines)
	length := len(lines[0])

	result := &mapState{data: make([][]byte, depth), length: length, depth: depth}

	for j := 0; j < depth; j++ {
		result.data[j] = make([]byte, length)
		for i := 0; i < length; i++ {
			result.data[j][i] = lines[j][i]
			switch lines[j][i] {
			case '^':
				result.robotX, result.robotY = i, j
			case 'v':
				result.robotX, result.robotY = i, j
			case '>':
				result.robotX, result.robotY = i, j
			case '<':
				result.robotX, result.robotY = i, j
			}
		}

	}
	return result
}

func (m *mapState) print() {
	lines := make([]string, m.depth)
	for k, newline := range m.data {
		lines[k] = string(newline)
	}
	screen := strings.Join(lines, "\n") + "\n"
	fmt.Println(screen)
	fmt.Printf("Current position: [%d,%d], direction: %s\n", m.robotX, m.robotY, string(m.data[m.robotY][m.robotX]))
}

func (m *mapState) calcAlignment() int {
	sum := 0
	for j := 1; j < m.depth-1; j++ {
		for i := 1; i < m.length-1; i++ {
			if m.data[j][i] == '#' &&
				m.data[j-1][i] == '#' &&
				m.data[j+1][i] == '#' &&
				m.data[j][i+1] == '#' &&
				m.data[j][i-1] == '#' {
				sum += (i * j)
			}
		}
	}
	return sum
}

func (m *mapState) calcPath() string {
	result := []byte{}

	for {
		// Work out direction
		var turnOp byte
		switch m.data[m.robotY][m.robotX] {
		case '^':
			if m.robotX != 0 && m.data[m.robotY][m.robotX-1] == '#' {
				m.data[m.robotY][m.robotX] = '<'
				turnOp = 'L'
			} else if m.robotX != m.length-1 && m.data[m.robotY][m.robotX+1] == '#' {
				m.data[m.robotY][m.robotX] = '>'
				turnOp = 'R'
			}
		case 'v':
			if m.robotX != 0 && m.data[m.robotY][m.robotX-1] == '#' {
				m.data[m.robotY][m.robotX] = '<'
				turnOp = 'R'
			} else if m.robotX != m.length-1 && m.data[m.robotY][m.robotX+1] == '#' {
				m.data[m.robotY][m.robotX] = '>'
				turnOp = 'L'
			}
		case '<':
			if m.robotY != 0 && m.data[m.robotY-1][m.robotX] == '#' {
				m.data[m.robotY][m.robotX] = '^'
				turnOp = 'R'
			} else if m.robotY != m.depth-1 && m.data[m.robotY+1][m.robotX] == '#' {
				m.data[m.robotY][m.robotX] = 'v'
				turnOp = 'L'
			}
		case '>':
			if m.robotY != 0 && m.data[m.robotY-1][m.robotX] == '#' {
				m.data[m.robotY][m.robotX] = '^'
				turnOp = 'L'
			} else if m.robotY != m.depth-1 && m.data[m.robotY+1][m.robotX] == '#' {
				m.data[m.robotY][m.robotX] = 'v'
				turnOp = 'R'
			}
		}

		// No new direction, we have reached the end
		if turnOp == 0 {
			break
		}

		result = append(result, turnOp, byte(','))

		// Work out Steps forward
		steps := 0
		for {
			finished := false
			switch m.data[m.robotY][m.robotX] {
			case '^':
				if m.robotY != 0 && m.data[m.robotY-1][m.robotX] == '#' {
					m.data[m.robotY][m.robotX] = '#'
					m.data[m.robotY-1][m.robotX] = '^'
					m.robotY--
				} else {
					finished = true
				}

			case 'v':
				if m.robotY != m.depth-1 && m.data[m.robotY+1][m.robotX] == '#' {
					m.data[m.robotY][m.robotX] = '#'
					m.data[m.robotY+1][m.robotX] = 'v'
					m.robotY++
				} else {
					finished = true
				}

			case '>':
				if m.robotX != m.length-1 && m.data[m.robotY][m.robotX+1] == '#' {
					m.data[m.robotY][m.robotX] = '#'
					m.data[m.robotY][m.robotX+1] = '>'
					m.robotX++
				} else {
					finished = true
				}

			case '<':
				if m.robotX != 0 && m.data[m.robotY][m.robotX-1] == '#' {
					m.data[m.robotY][m.robotX] = '#'
					m.data[m.robotY][m.robotX-1] = '<'
					m.robotX--
				} else {
					finished = true
				}
			}
			if finished {
				break
			}
			steps++
		}

		for _, v := range strconv.Itoa(steps) {
			result = append(result, byte(v))
		}
		result = append(result, byte(','))

	}
	return string(result[:len(result)-1])
}
