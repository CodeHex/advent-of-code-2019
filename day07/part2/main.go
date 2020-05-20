package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"gonum.org/v1/gonum/stat/combin"
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

	phasesSettings := combin.Permutations(5, 5)
	// shift all permuations to range 5 ot 9
	for _, perm := range phasesSettings {
		for i, val := range perm {
			perm[i] = val + 5
		}
	}

	var maxOutput int
	var maxPerm []int
	for _, perm := range phasesSettings {
		c, err := newFeedbackComputer(inputText, "AMP A", "AMP B", "AMP C", "AMP D", "AMP E")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		c.runAsync()

		err = c.loadPhases(perm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		c.input(0)
		c.waitForCompletion()

		output := c.output()

		if output > maxOutput {
			maxOutput = output
			maxPerm = perm
		}

		fmt.Printf("perm %v gives output value %d\n", perm, output)
	}

	fmt.Println()
	fmt.Printf("max perm %v gives output value %d\n", maxPerm, maxOutput)
}
