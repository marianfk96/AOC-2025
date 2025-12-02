package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	file, err := os.ReadFile("filepath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	input := string(file)
	inputArray := strings.Split(input, ",")
	invalidCodes := []int{}
	invalidSequances := []int{}
	invalidCodeSum := 0
	invalidSequncesSum := 0

	for _, val := range inputArray {
		strA, strB, _ := strings.Cut(val, "-")
		minDigits := len(strA)
		maxDigits := len(strB)
		// avoid looping a to b -> get a range and generate invalid ids within that range
		res, sum := generateRepeatedNumbers(strA, strB, minDigits, maxDigits)
		invalidCodes = append(invalidCodes, res...)
		invalidCodeSum += sum
		resp, summ := repeatedSeaquences(strA, strB, minDigits, maxDigits)
		invalidSequances = append(invalidSequances, resp...)
		invalidSequncesSum += summ
	}
	fmt.Println(invalidSequncesSum)
}

func generateRepeatedNumbers(a string, b string, minDigits int, maxDigits int) ([]int, int) {
	numA, _ := strconv.Atoi(a)
	numB, _ := strconv.Atoi(b)
	if minDigits%2 == 1 {
		minDigits++
	}
	sum := 0

	results := []int{}
	for i := minDigits; i <= maxDigits; i += 2 {
		halfLenght := i / 2         // we know i always even
		smallest := toIntPow(i - 1) // smallest number with i digits
		largest := toIntPow(i) - 1  // largest number with i digits
		// clamp restraints of the loop between numA and numB
		if numA > smallest {
			smallest = numA
		}
		if numB < largest {
			largest = numB
		}
		for n := smallest; n <= largest; n++ {
			strN := strconv.Itoa(n)
			firstHalf := strN[0:halfLenght]
			secondHalf := strN[halfLenght:]
			if firstHalf == secondHalf {
				results = append(results, n)
				sum += n
			}
		}
	}
	return results, sum
}

func repeatedSeaquences(a string, b string, minDigits int, maxDigits int) ([]int, int) {
	numA, _ := strconv.Atoi(a)
	numB, _ := strconv.Atoi(b)
	sum := 0
	results := []int{}
	for i := minDigits; i <= maxDigits; i++ {
		smallest := toIntPow(i - 1) // smallest number with i digits
		largest := toIntPow(i) - 1  // largest number with i digits
		// clamp restraints of the loop between numA and numB
		if numA > smallest {
			smallest = numA
		}
		if numB < largest {
			largest = numB
		}
		for n := smallest; n <= largest; n++ {
			if isRepeatingPattern(n) {
				results = append(results, n)
				sum += n
			}
		}
	}
	return results, sum
}

func isRepeatingPattern(num int) bool {
	numS := strconv.Itoa(num)
	digits := len(numS)
	// if i larger than half the numebr digits, we cannot have repeating sequences -> stop there
	for i := 1; i <= digits/2; i++ {
		// we want perfect division -> prime numbers only have 1-digit sequences, everything else
		// divisor-digit sequences -> continue to next iteration
		if digits%i != 0 {
			continue
		}
		//get first substring of size i
		substring := numS[0:i]
		repetiotions := digits / i
		if strings.Repeat(substring, repetiotions) == numS {
			return true
		}
	}
	return false
}

// math pow returns float which may cause issues with very large numbers -> custom func
func toIntPow(pow int) int {
	res := 1
	for i := 0; i < pow; i++ {
		res *= 10
	}
	return res
}
