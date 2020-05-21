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

	c, err := newComputer(inputText, scanner)
	if err != nil {
		fmt.Printf("unable to convert input to int code memory, %s\n", err.Error())
		os.Exit(1)
	}
	c.disableLog = true

	err = c.run()
	if err != nil {
		fmt.Printf("error running program: %s\n", err.Error())
		os.Exit(1)
	}
}
