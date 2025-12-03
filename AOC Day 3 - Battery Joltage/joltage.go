package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	file, err := os.Open("filepath")
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
		max := -1
		for i, b := range bank {
			for j := len(bank) - 1; j > i; j-- {
				num := b*10 + bank[j]
				if num > max {
					max = num
				}
			}
		}
		joltageSum += max
		turboJoltageSum += max12DigitJolatge(bank)
	}
	fmt.Println(joltageSum)
	fmt.Println(turboJoltageSum)
}

func max12DigitJolatge(bank []int) int {

	start := 0
	maxSequence := 0 // ! not -1, we're multiplying it !
	// decrease the remaining digits we need to select
	for remainingDigits := 12; remainingDigits > 0; remainingDigits-- {
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
