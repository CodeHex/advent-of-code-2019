package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	template, err := newComputer(scanner.Text())
	if err != nil {
		fmt.Printf("unable to convert input to int code memory, %s\n", err.Error())
		return
	}

	const requiredResult = 19690720

	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			c := template.clone()
			c.storeAtAddr(1, noun)
			c.storeAtAddr(2, verb)

			err = c.run()
			if err != nil {
				fmt.Printf("error running program: %s\n", err.Error())
				return
			}

			if c.readAddr(0) == requiredResult {
				fmt.Printf("Found result at noun: %d and verb %d (val at zero is %d)\n", noun, verb, c.readAddr(0))
				fmt.Printf("Calculated value 100 * noun + verb is: %d\n", (noun*100)+verb)
				return
			}
		}
	}
	fmt.Printf("No valid noun verb combination found for %d\n", requiredResult)
}
