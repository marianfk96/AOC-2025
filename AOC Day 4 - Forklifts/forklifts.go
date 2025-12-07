package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("filepath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	var paperGrid [][]string
	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lines++
		row := strings.Split(line, "")
		paperGrid = append(paperGrid, row)
	}
	lineLength := len(paperGrid[0])
	accessible := 0
	for i := 0; i < lines; i++ {
		for j := 0; j < lineLength; j++ {
			if paperGrid[i][j] != "@" {
				continue
			}
			adjacent := countAdjacent(i, j, paperGrid, lineLength, lines)
			if adjacent < 4 {
				accessible++
			}
		}
	}
	canBeRemoved := removeAccessible(paperGrid, lines, lineLength)
	fmt.Println(accessible, canBeRemoved)
}

func removeAccessible(paperGrid [][]string, lines int, lineLength int) int {
	accessible := 0
	removedItems := 0
	for i := 0; i < lines; i++ {
		for j := 0; j < lineLength; j++ {
			if paperGrid[i][j] != "@" {
				continue
			}
			adjacent := countAdjacent(i, j, paperGrid, lineLength, lines)
			if adjacent < 4 {
				accessible++
				paperGrid[i][j] = "x"
				removedItems++
			}
		}
	}
	if removedItems == 0 {
		return 0
	}
	return removedItems + removeAccessible(paperGrid, lines, lineLength)
}

func countAdjacent(i int, j int, paperGrid [][]string, cols int, rows int) int {
	adjacent := 0
	//left
	if j-1 >= 0 && paperGrid[i][j-1] == "@" {
		adjacent++
	}
	// diagonal left
	if i-1 >= 0 && j-1 >= 0 && paperGrid[i-1][j-1] == "@" {
		adjacent++
	}
	// above
	if i-1 >= 0 && paperGrid[i-1][j] == "@" {
		adjacent++
	}
	// diagonal right
	if i-1 >= 0 && j+1 < cols && paperGrid[i-1][j+1] == "@" {
		adjacent++
	}
	// right
	if j+1 < cols && paperGrid[i][j+1] == "@" {
		adjacent++
	}
	// low diagonal right
	if i+1 < rows && j+1 < cols && paperGrid[i+1][j+1] == "@" {
		adjacent++
	}
	// below
	if i+1 < rows && paperGrid[i+1][j] == "@" {
		adjacent++
	}
	// low left diagonal
	if i+1 < rows && j-1 >= 0 && paperGrid[i+1][j-1] == "@" {
		adjacent++
	}
	return adjacent
}
