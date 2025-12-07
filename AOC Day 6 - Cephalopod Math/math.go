package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	file, err := os.Open("filePath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	operands := [][]string{}
	operators := [][]string{}
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		parts := strings.Fields(line)
		if isOperator(parts[0]) {
			operators = append(operators, parts)
		} else {
			operands = append(operands, parts)
		}
	}
	lines = lines[:len(lines)-1]
	// part 1
	sum := 0
	for j := 0; j < len(operators[0]); j++ { // operators is not 2d after all , but leaving it here :)
		for i := 0; i < len(operators); i++ {
			if operators[i][j] == "+" {
				sum += sumByColumn(j, operands)
			} else {
				sum += multiplyByColumn(j, operands)
			}
		}
	}
	fmt.Println(sum)
	// part 2
	boundaries := []int{}
	firstLine := lines[0]
	// find the boundaries -> empty for EVERY line
	for col, ch := range firstLine {
		if ch == ' ' && isBoundary(col, lines) {
			boundaries = append(boundaries, col)
		}
	}
	formattedInput := [][]string{}
	for _, line := range lines {
		splits := []string{}
		start := 0
		// extract the substrings, preserving spaces, but skiping boundaries
		for _, bound := range boundaries {
			splits = append(splits, line[start:bound])
			start = bound + 1
		}
		// add the last line
		splits = append(splits, line[start:])
		formattedInput = append(formattedInput, splits)
	}
	transformedNums := [][]string{}
	for col := 0; col < len(formattedInput[0]); col++ {
		colArr := convertColToArray(col, formattedInput)
		trans := transformColumn(colArr)
		transformedNums = append(transformedNums, trans)
	}
	//ttransformed nums are now my operands and i can rerun logic for part 1
	sum2 := 0
	for j := 0; j < len(operators[0]); j++ { // operators is not 2d after all , but leaving it here :)
		for i := 0; i < len(operators); i++ {
			if operators[i][j] == "+" {
				sum2 += sumByLine(transformedNums[j])
			} else {
				sum2 += multiplyByLine(transformedNums[j])
			}
		}
	}
	fmt.Println(sum2)
}

func sumByLine(arr []string) int {
	sum := 0
	for _, val := range arr {
		intVal := toInteger(val)
		sum += intVal
	}
	return sum
}

func multiplyByLine(arr []string) int {
	sum := 1
	for _, val := range arr {
		intVal := toInteger(val)
		sum *= intVal
	}
	return sum
}

func transformColumn(arr []string) []string {
	trans := []string{}
	lenght := len(arr[0])
	for col := 0; col < lenght; col++ {
		s := ""
		for row := 0; row < len(arr); row++ {
			line := arr[row]
			s += string(line[col])
		}
		trans = append(trans, s)
	}
	return trans
}

func isBoundary(col int, lines []string) bool {
	for r := 1; r < len(lines); r++ {
		line := lines[r]
		if line[col] != ' ' {
			return false
		}
	}
	return true
}

func convertColToArray(col int, operands [][]string) []string {
	arr := []string{}
	for i := 0; i < len(operands); i++ {
		arr = append(arr, operands[i][col])
	}
	return arr
}

func sumByColumn(col int, operands [][]string) int {
	sum := 0
	for i := 0; i < len(operands); i++ {
		sum += toInteger(operands[i][col])
	}
	return sum
}

func multiplyByColumn(col int, operands [][]string) int {
	sum := 1
	for i := 0; i < len(operands); i++ {
		sum *= toInteger(operands[i][col])
	}
	return sum
}

func toInteger(a string) int {
	a = strings.TrimSpace(a)
	num, _ := strconv.Atoi(a)
	return num
}

func isOperator(c string) bool {
	return c == "+" || c == "*" || c == "/" || c == "%" || c == "-"
}
