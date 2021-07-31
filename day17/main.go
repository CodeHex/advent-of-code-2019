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

	comp.run()
	clearScreen()
	outputScreenBytes := make([]byte, len(comp.outputs))
	for i, v := range comp.outputs {
		outputScreenBytes[i] = byte(v)
	}
	outputScreen := string(outputScreenBytes)
	fmt.Println(outputScreen)
	fmt.Printf("alignment = %d\n", calcAlignment(outputScreen))
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func calcAlignment(out string) int {
	lines := strings.Split(out, "\n")
	lines = lines[:len(lines)-1] // Remove last empty line
	n := len(lines)
	sum := 0
	for j := 1; j < n-2; j++ {
		for i := 1; i < len(lines[j])-2; i++ {
			if lines[j][i] == '#' &&
				lines[j-1][i] == '#' &&
				lines[j+1][i] == '#' &&
				lines[j][i+1] == '#' &&
				lines[j][i-1] == '#' {
				sum += (i * j)
			}
		}
	}
	return sum
}
