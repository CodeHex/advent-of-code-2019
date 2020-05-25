package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var inputData []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		inputData = append(inputData, line)
	}

	formulaBook := newBook(inputData)
	numOfOre := formulaBook.decomposeToOre(newComponent("1 FUEL"))
	fmt.Printf("number of ore required to produce FUEL is %d\n", numOfOre)

	startOre := int64(1000000000000)
	numOfFuel := formulaBook.createFuelFromOre(startOre, numOfOre)
	fmt.Printf("number of fuel created from %d ORE is %d\n", startOre, numOfFuel)
}

type book map[string]*formula

type component struct {
	chemical string
	amount   int64
}

type formula struct {
	product  component
	reagents []component
}

func newBook(inputData []string) book {
	recipes := make(map[string]*formula, len(inputData))
	for _, entry := range inputData {
		formula := newFormula(entry)
		recipes[formula.product.chemical] = formula
	}
	return recipes
}

func newComponent(entry string) component {
	parts := strings.Split(entry, " ")
	amount, _ := strconv.Atoi(parts[0])
	chemical := parts[1]
	return component{chemical: chemical, amount: int64(amount)}
}

func newFormula(entry string) *formula {
	parts := strings.Split(entry, "=>")

	reagentParts := strings.Split(parts[0], ",")
	reagents := make([]component, len(reagentParts))
	for i, reagentInput := range reagentParts {
		reagents[i] = newComponent(strings.TrimSpace(reagentInput))
	}

	product := newComponent(strings.TrimSpace(parts[1]))
	return &formula{reagents: reagents, product: product}
}

func (c component) String() string {
	return fmt.Sprintf("%d %s", c.amount, c.chemical)
}

func (f *formula) String() string {
	reagentStr := ""
	for i, reagent := range f.reagents {
		reagentStr += reagent.String()
		if i != len(f.reagents)-1 {
			reagentStr += " + "
		}
	}
	return fmt.Sprintf("%s = %s", reagentStr, f.product.String())
}

func (b book) String() string {
	result := ""
	count := 0
	for _, formula := range b {
		result += formula.String()
		if count != len(b)-1 {
			result += "\n"
		}
		count++
	}
	return result
}

type balancer map[string]int64

func newBalancer(comp component) balancer {
	balancer := make(map[string]int64)
	balancer[comp.chemical] = comp.amount
	return balancer
}

// primary formulas  chooses the first formula that will break down
func (b balancer) nextPrimaryFormula(formulaBook book) *formula {
	for chemical := range b {
		// If there is no chemical to convert or we have no formula move on
		if b[chemical] <= 0 || formulaBook[chemical] == nil {
			continue
		}
		for _, r := range formulaBook[chemical].reagents {
			if b[r.chemical] > 0 {
				return formulaBook[chemical]
			}
		}
	}
	return nil
}

// secondary formulas are all currently stored
// and will add to existing store ingredients
func (b balancer) nextFormula(formulaBook book) *formula {
	for chemical := range b {
		// If there is no chemical to convert or we have no formula move on
		if b[chemical] <= 0 || formulaBook[chemical] == nil || chemical == "ORE" {
			continue
		}
		return formulaBook[chemical]
	}
	return nil
}

func (b balancer) applyFormula(f *formula) balancer {
	chemical := f.product.chemical
	numOfDecomps := int64(0)
	for numOfDecomps*f.product.amount < b[chemical] {
		numOfDecomps++
	}
	b[chemical] -= numOfDecomps * f.product.amount
	for _, reagent := range f.reagents {
		b[reagent.chemical] += numOfDecomps * reagent.amount
	}
	for k, v := range b {
		if v == 0 {
			delete(b, k)
		}
	}
	return b
}

func (b book) decomposeToOre(comp component) int64 {
	// keep a list of chemicals we need to compose the component
	// negative numbers indicate we have a surplus after the reactions
	store := newBalancer(comp)
	store[comp.chemical] = comp.amount

	// only stop when there are no formulas to apply that can decompose the products
	for {
		if primaryFormula := store.nextPrimaryFormula(b); primaryFormula != nil {
			store = store.applyFormula(primaryFormula)
			continue
		}
		nextFormula := store.nextFormula(b)
		if nextFormula == nil {
			break
		}
		store = store.applyFormula(nextFormula)
	}
	return store["ORE"]
}

func (b balancer) createChemical(formulaBook book, comp component) (balancer, bool) {
	formula := formulaBook[comp.chemical]
	numOfReactions := int64(1)
	for numOfReactions*formula.product.amount < comp.amount {
		numOfReactions++
	}

	// Create the minimum number of reagents
	allReagentsFound := false
	for !allReagentsFound {
		allReagentsFound = true
		for _, reagent := range formula.reagents {
			reqAmount := reagent.amount * numOfReactions
			if b[reagent.chemical] >= reqAmount {
				continue
			}
			// If we don't have enough ore, fail the creation
			if reagent.chemical == "ORE" {
				return b, false
			}

			compToCreate := component{
				chemical: reagent.chemical,
				amount:   reqAmount - b[reagent.chemical],
			}
			var ok bool
			b, ok = b.createChemical(formulaBook, compToCreate)
			if !ok {
				return b, ok
			}
			allReagentsFound = false
			break
		}
	}

	// Perform the reaction
	b[comp.chemical] += formula.product.amount * numOfReactions
	for _, r := range formula.reagents {
		b[r.chemical] -= r.amount * numOfReactions
	}
	return b, true
}

func (b book) createFuelFromOre(ore int64, orePerFuel int64) int64 {
	store := newBalancer(component{"ORE", ore})
	ok := true

	for ok {
		// Jump ahead by the approx number of fuel we can definitely create avoiding extra reagents creates
		// Once we run out of ore to create a single one, try to create one at a time
		fuelJump := store["ORE"] / orePerFuel
		if fuelJump == 0 {
			fuelJump = 1
		}
		store, ok = store.createChemical(b, component{"FUEL", fuelJump})
	}
	return store["FUEL"]
}
