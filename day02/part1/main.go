package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	p, err := newProgram(scanner.Text())
	if err != nil {
		fmt.Printf("unable to convert input to int code memory, %s\n", err.Error())
		return
	}

	err = p.run()
	if err != nil {
		fmt.Printf("error running program: %s\n", err.Error())
		return
	}

	// Print out the first value of the programs memory
	fmt.Printf("Value at zero index is: %d\n", p.readAddr(0))
}
