package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// memoization caches a result -> only nodes that can lead to 'out' are left in the graph after pruning.
// so all we need to know when reaching a node is whethwe we've passed fft and dac -> if we have, valid path exists
type MemoState struct {
	node      string
	passedFFT bool
	passedDAC bool
}

func main() {
	file, err := os.Open("filePath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	// read the file as a map, for O(1) lookup times -> fingers crossed this'll help for part 2
	paths := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		node := parts[0]
		//split a string on whitespace
		neighbors := strings.Fields(parts[1])
		paths[string(node)] = neighbors
	}
	// save results for later ?
	results := [][]string{}
	visited := make(map[string]bool)
	// start dfs from node "you" and with an initially empty path
	dfs(paths, "you", []string{}, visited, &results)

	fmt.Println(len(results))

	//PART 2
	// graph is huge -> prune nodes that cannot reach out, they're useless to our problem
	// instead of checking whether a node can reach out, we can reverse tha graph and check which nodes 'out' can reach
	reversePaths := make(map[string][]string)

	for fromNode, neighbors := range paths {
		for _, toNode := range neighbors {
			// flip the stored path for each node
			reversePaths[toNode] = append(reversePaths[toNode], fromNode)
		}
	}
	// prune the graph
	prunedPaths := pruneGraph(paths, reversePaths, "out")
	memo := make(map[MemoState]int)
	validPaths := countValidPaths("svr", prunedPaths, false, false, memo)
	fmt.Println("no. of valid paths: ", validPaths)
}

func countValidPaths(node string, prunedPaths map[string][]string, passedFFT, passedDAC bool, memo map[MemoState]int) int {
	state := MemoState{
		node:      node,
		passedFFT: passedFFT,
		passedDAC: passedDAC,
	}
	// if we've already seen this state, return its value (mnemoized result)
	val, exists := memo[state]
	if exists {
		return val
	}
	if node == "out" {
		// if we've passed fft and dac when we reach out -> valid path
		if passedFFT && passedDAC {
			return 1
		} else {
			return 0
		}
	}
	//update the flags
	if node == "fft" {
		passedFFT = true
	} else if node == "dac" {
		passedDAC = true
	}
	totalPaths := 0
	//recurse to neighboring nodes
	for _, next := range prunedPaths[node] {
		totalPaths += countValidPaths(next, prunedPaths, passedFFT, passedDAC, memo)
	}

	// store the result and return the number of paths
	memo[state] = totalPaths
	return totalPaths
}

func pruneGraph(paths, reversePaths map[string][]string, target string) map[string][]string {
	prunedPaths := make(map[string][]string)
	canReachOut := make(map[string]bool)
	markUsefulNodes(target, reversePaths, canReachOut)
	for node, neighbors := range paths {
		// skip nodes that can't reach out
		if !canReachOut[node] {
			continue
		}
		for _, next := range neighbors {
			if canReachOut[next] {
				prunedPaths[node] = append(prunedPaths[node], next)
			}
		}
	}
	return prunedPaths
}

func markUsefulNodes(node string, reversePaths map[string][]string, canReach map[string]bool) {
	if canReach[node] {
		return
	}
	canReach[node] = true
	for _, prevNode := range reversePaths[node] {
		markUsefulNodes(prevNode, reversePaths, canReach)
	}
}

func dfs(
	paths map[string][]string,
	startingPoint string,
	currentPath []string,
	visited map[string]bool,
	results *[][]string,
) {
	// add node we're working with to the current path
	currentPath = append(currentPath, startingPoint)
	// when we reach out, we're done exploring this path (branch)
	if startingPoint == "out" {
		// change object references before appending
		temp := make([]string, len(currentPath))
		copy(temp, currentPath)
		*results = append(*results, temp)
		return
	}
	//mark node as visited
	visited[startingPoint] = true
	for _, next := range paths[startingPoint] {
		// run dfs for next unvisited node
		if !visited[next] {
			dfs(paths, next, currentPath, visited, results)
		}
	}
	// when we're done exploring a specific 'branch' , remove starting point (node currently being examined) from visited
	// so it can be reused
	visited[startingPoint] = false
}
