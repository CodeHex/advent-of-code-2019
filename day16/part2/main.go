package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// Get signal and expand to 10,000 times as much
	signal := parseInput(input)
	expandedSignal := []int{}
	for i := 0; i < 10000; i++ {
		expandedSignal = append(expandedSignal, signal...)
	}
	signal = expandedSignal
	fmt.Println("signal generated")

	// Read offset
	offset := calcOffset(signal)
	fmt.Printf("offset: %d\n", offset)

	// Check that we are past half way through the signal. This will ensure that
	// x[k] = SUM^N_n=k x[n]
	// for the number of digits we want
	if offset < (len(signal)/2)+1 {
		panic("unsupported, offset must be at least halfway through the signal")
	}

	// Truncate signal as previous digits are not required
	signal = signal[offset:]
	fmt.Printf("truncated signal length: %d\n", len(signal))

	// Apply 100 times
	for i := 1; i <= 100; i++ {
		signal = transformSum(signal)
	}

	fmt.Printf("after 100 phases: %s\n", signalString(signal[0:8]))
}

func parseInput(input string) []int {
	digits := make([]int, len(input))
	var err error
	for i, char := range input {
		digits[i], err = strconv.Atoi(string(char))
		if err != nil {
			panic(errors.Wrap(err, "unable to parse input"))
		}
	}
	return digits
}

func calcOffset(signal []int) int {
	offsetStr := ""
	for _, val := range signal[0:7] {
		offsetStr = offsetStr + strconv.Itoa(val)
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		panic(err)
	}
	return offset
}

func transformSum(signal []int) []int {
	newSignal := make([]int, len(signal))

	// Since we are past half way, each signal entry is transformed to the sum of
	// itself and proceeding digits i.e.
	// N = length of the truncated signal
	// k = index of the signal entry
	// x = original signal truncated by the offset
	// X = transformed signal truncated by the offset
	//
	//             N
	//   X[k] =(  SUM  x[n]  )  REM 10
	//            n=k
	//
	// Therefore we start at the end of the signal (k=N) and move backwards so we
	// can keep a running total i.e.
	//
	//  X[N] = x[N]  and  X[k-1] =  (X[k] + x[k]) REM 10
	//
	total := 0
	for i := len(signal) - 1; i >= 0; i-- {
		total += signal[i]
		newSignal[i] = total % 10
	}
	return newSignal
}

func signalString(signal []int) string {
	out := ""
	for _, digit := range signal {
		out += strconv.Itoa(digit)
	}
	return out
}
