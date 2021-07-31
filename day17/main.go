package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	clearScreen()
	outputScreenBytes := make([]byte, len(comp.outputs))
	for i, v := range comp.outputs {
		outputScreenBytes[i] = byte(v)
	}
	outputScreen := string(outputScreenBytes)
	alignment, newScreen := calcAlignment(outputScreen)
	fmt.Println(newScreen)
	fmt.Printf("alignment = %d\n", alignment)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func calcAlignment(out string) (int, string) {
	lines := strings.Split(out, "\n")
	lines = lines[:len(lines)-2] // Remove last empty line
	depth := len(lines)
	length := len(lines[0])

	fmt.Println(length, depth)

	newScreenPixels := make([][]byte, depth)
	sum := 0
	for j := 0; j < depth; j++ {
		newScreenPixels[j] = make([]byte, length)
		for i := 0; i < length; i++ {
			newScreenPixels[j][i] = lines[j][i]

			if i == 0 || i == length-1 || j == 0 || j == depth-1 {
				continue
			}

			if lines[j][i] == '#' &&
				lines[j-1][i] == '#' &&
				lines[j+1][i] == '#' &&
				lines[j][i+1] == '#' &&
				lines[j][i-1] == '#' {
				newScreenPixels[j][i] = 'O'
				sum += (i * j)
			}
		}
	}
	newlines := make([]string, depth)
	for k, newline := range newScreenPixels {
		newlines[k] = string(newline)
	}
	newScreen := strings.Join(newlines, "\n") + "\n"
	return sum, newScreen
}
