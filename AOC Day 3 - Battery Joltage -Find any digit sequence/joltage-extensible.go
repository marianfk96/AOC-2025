package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	file, err := os.Open("filePath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	var banks [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var bank []int
		for _, ch := range line {
			bank = append(bank, int(ch-'0'))
		}
		banks = append(banks, bank)
	}
	joltageSum := 0
	turboJoltageSum := 0
	for _, bank := range banks {
		joltageSum += maxJoltage(bank, 2)
		turboJoltageSum += maxJoltage(bank, 12)
	}
	fmt.Println(joltageSum)
	fmt.Println(turboJoltageSum)
}

func maxJoltage(bank []int, maxDigits int) int {

	start := 0
	maxSequence := 0 // ! not -1, we're multiplying it !
	// decrease the remaining digits we need to select
	for remainingDigits := maxDigits; remainingDigits > 0; remainingDigits-- {
		// make sure we have enough entries to check afterwards
		end := len(bank) - remainingDigits
		max := -1
		posMax := -1
		for i := start; i <= end; i++ {
			if bank[i] > max {
				max = bank[i]
				posMax = i
			}
		}
		maxSequence = maxSequence*10 + max
		// we now start from the enxt element after max
		start = posMax + 1
	}
	return maxSequence
}
