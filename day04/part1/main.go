package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := scanner.Text()
	if !strings.Contains(input, "-") {
		// Single value mode, just checks if the password is valid
		pass, err := newPassword(input)
		if err != nil {
			fmt.Printf("unable to parse input '%s': %s\n", input, err.Error())
			return
		}
		fmt.Printf("valid password: %t\n", pass.isValid())
		return
	}

	// Range mode
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		fmt.Println("unable to parse input, requires single range")
		return
	}
	lowerPass, err := newPassword(parts[0])
	if err != nil {
		fmt.Printf("unable to parse lower range '%s': %s\n", parts[0], err.Error())
		return
	}
	upperPass, err := newPassword(parts[1])
	if err != nil {
		fmt.Printf("unable to parse upper range '%s': %s\n", parts[1], err.Error())
		return
	}
	if lowerPass.greaterThan(upperPass) {
		fmt.Println("invalid range lower password is greater than higher password")
		return
	}
	fmt.Printf("from %s to %s\n", lowerPass.string(), upperPass.string())

	// Loop through all passwords in the range checking validity and keeping an invalid and valid count
	pass := lowerPass
	validCount := 0
	invalidCount := 0
	for !pass.greaterThan(upperPass) {
		isValid := pass.isValid()
		fmt.Printf("password %s - valid:%t\n", pass.string(), isValid)
		if isValid {
			validCount++
		} else {
			invalidCount++
		}
		pass = pass.next()
	}
	fmt.Printf("invalid count %d / valid count %d\n", invalidCount, validCount)
}

// ----- password -------

type password [6]int

func newPassword(s string) (password, error) {
	if len(s) != 6 {
		return password{}, errors.New("invalid password, needs to be 6 digits")
	}

	result := password{}
	for i, char := range s {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return password{}, errors.New("invalid password, unable to parse digits")
		}
		result[i] = digit
	}
	return result, nil
}

func (p password) greaterThan(cmp password) bool {
	for i := 0; i < 6; i++ {
		switch {
		case p[i] > cmp[i]:
			return true
		case p[i] < cmp[i]:
			return false
		}
	}
	return false
}

func (p password) string() string {
	return fmt.Sprintf("%d%d%d%d%d%d", p[0], p[1], p[2], p[3], p[4], p[5])
}

func (p password) isValid() bool {
	var isDouble bool
	for i := 0; i < 5; i++ {
		if p[i] == p[i+1] {
			isDouble = true
		}

		// Fail validity if next digit is decresing
		if p[i+1] < p[i] {
			return false
		}
	}
	return isDouble
}

func (p password) next() password {
	for i := 5; i >= 0; i-- {
		if p[i] == 9 {
			p[i] = 0
			continue
		}
		p[i]++
		return p
	}
	return p
}
