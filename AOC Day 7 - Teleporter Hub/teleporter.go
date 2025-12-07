package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("beams.txt")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	startRow := 0
	startCol := 0
	splitBeams := []string{}
	for i := 0; i < len(lines)-1; i += 2 {
		line := lines[i]
		tmp := []rune(lines[i])
		tmp2 := []rune(lines[i+1])
		for j, c := range line {
			if i > 0 {
				prev := []rune(splitBeams[i-1])
				if c == '.' && prev[j] == '|' {
					tmp2[j] = '|'
				}
			}
			if string(c) == "S" {
				tmp2[j] = '|'
				startRow = i
				startCol = j
			} else if string(c) == "^" {
				tmp[j-1] = '|'
				tmp[j+1] = '|'
				tmp2[j-1] = '|'
				tmp2[j+1] = '|'
			}
		}
		splitBeams = append(splitBeams, string(tmp))
		splitBeams = append(splitBeams, string(tmp2))
	}
	splits := 0
	for i, beam := range splitBeams {
		for j, c := range beam {
			if i > 0 && c == '^' && splitBeams[i-1][j] == '|' {
				splits++
			}
		}
	}
	fmt.Println(splits)
	// splitBeams already represents a tree structure
	// S is the root node.
	// | is the edge that connects nodes
	// ^ is represents a node.
	// which means i can run dfs and count every time i reach a leaf node
	// graph is directed and acyclic + multiple "parents" can lead to the same node -> I need memoization

	// conver splitBeams to [][]rune
	beamsRunes := make([][]rune, len(splitBeams))
	for i, val := range splitBeams {
		beamsRunes[i] = []rune(val)
	}

	// decalre a map : key [2]int = an array of 2 integers, value = int
	// essentially key is the 'coordinates' of the node, value is the number of possible paths from that node
	pathsFromNode := make(map[[2]int]int)
	possiblePaths := dfs(beamsRunes, startRow, startCol, pathsFromNode)
	fmt.Println("Possible paths: ", possiblePaths)

}

func dfs(beams [][]rune, startRow int, startCol int, pathsFromNode map[[2]int]int) int {
	// start processing from root
	key := [2]int{startRow, startCol}
	val, exists := pathsFromNode[key]
	// if node has been processed return the memoized result
	if exists {
		return val
	}
	// if we've reached the last row, a complete path has been found -> return 1
	if startRow == len(beams)-1 {
		pathsFromNode[key] = 1
		return 1
	}
	total := 0
	nextRow := startRow + 1
	cell := beams[startRow][startCol]

	if cell == '^' {
		// recurse left
		if startCol > 0 {
			total += dfs(beams, startRow, startCol-1, pathsFromNode)
		}
		// recurse right
		if startCol < len(beams[0])-1 {
			total += dfs(beams, startRow, startCol+1, pathsFromNode)
		}
	} else {
		// did not find a node -> continue to the next line
		total += dfs(beams, nextRow, startCol, pathsFromNode)
	}

	// once out of the loop, store total
	pathsFromNode[key] = total
	return total
}
