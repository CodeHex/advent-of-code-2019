package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Continually scan in mass values from stdin and print out a running total of the calculated fuel
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	totalFuel := 0

	for scanner.Scan() {
		massString := scanner.Text()
		mass, err := strconv.Atoi(massString)
		if err != nil {
			fmt.Printf("unable to convert line to int:'%s', %s\n", massString, err.Error())
			return
		}

		fuel, err := calculateFuel(mass)
		if err != nil {
			fmt.Printf("unable to calculate fuel for mass:%d, %s\n", mass, err.Error())
			return
		}

		totalFuel += fuel
		fmt.Printf("mass %d needs %d fuel (total fuel %d)\n", mass, fuel, totalFuel)
	}

	fmt.Printf("\nTotal fuel: %d\n", totalFuel)
}

func calculateFuel(mass int) (int, error) {
	// integer division is floor division
	result := mass / 3
	result = result - 2

	if result < 1 {
		return 0, errors.Errorf("invalid amount of fuel require (%d)", result)
	}
	return result, nil
}
