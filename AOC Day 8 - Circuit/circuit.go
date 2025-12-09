package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type PointNode struct {
	X, Y, Z float64
	ID      int
}

type Edge struct {
	FROM_NODE, TO_NODE PointNode
	FROM_ID, TO_ID     int
	WEIGHT             float64
}

type UnionFind struct {
	parent []int // intially nodes are parents of themselves
	rank   []int // height of the tree at parent rank -> use it to balance the tree
	size   []int // number of items in the circuit
}

func createUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: []int{},
		rank:   []int{},
		size:   []int{},
	}
	for i := 0; i < n; i++ {
		uf.parent = append(uf.parent, i)
		uf.rank = append(uf.rank, 0)
		uf.size = append(uf.size, 1) // intially all nodes are seperate curcuits
	}
	return uf
}

func (uf *UnionFind) find(node int) int {
	// start from any node, and traverse up until you find root node
	if uf.parent[node] != node {
		uf.parent[node] = uf.find(uf.parent[node])
	}
	return uf.parent[node]
}

func (uf *UnionFind) union(nodeA, nodeB int) bool {
	rootA := uf.find(nodeA)
	rootB := uf.find(nodeB)

	// already on the same set
	if rootA == rootB {
		return false
	}
	// always attach the shorter tree under the larger one
	if uf.rank[rootA] < uf.rank[rootB] {
		uf.parent[rootA] = rootB
		uf.size[rootB] += uf.size[rootA]
	} else if uf.rank[rootB] < uf.rank[rootA] {
		uf.parent[rootB] = rootA
		uf.size[rootA] += uf.size[rootB]
	} else { // in case they're equal, arbitrarily choose and increase rank
		uf.parent[rootB] = rootA
		uf.rank[rootA]++
		uf.size[rootA] += uf.size[rootB]
	}
	// union performed
	return true
}

func main() {
	file, err := os.Open("filePath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	points := []PointNode{}
	scanner := bufio.NewScanner(file)
	pos := 0
	for scanner.Scan() {
		line := scanner.Text()
		nodeCoordinates := strings.Split(line, ",")
		x, _ := strconv.ParseFloat(nodeCoordinates[0], 64)
		y, _ := strconv.ParseFloat(nodeCoordinates[1], 64)
		z, _ := strconv.ParseFloat(nodeCoordinates[2], 64)
		points = append(points, PointNode{x, y, z, pos})
		pos++
	}

	edges := []Edge{}
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			edge := Edge{points[i], points[j], points[i].ID, points[j].ID, calculatePointDistance(points[i], points[j])}
			edges = append(edges, edge)
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].WEIGHT < edges[j].WEIGHT
	})

	closestEdges := edges[0:1000]

	uf := createUnionFind(len(points))
	// create unions from all the edges (they contain from - to nodes)
	for _, edge := range closestEdges {
		uf.union(edge.FROM_ID, edge.TO_ID)
	}

	// collect evrything into a map, grouping them by the root -> all elements in the same uniion will return the same root
	// with find, so we can find how many points each set contains
	circuits := map[int][]int{}
	for i := 0; i < len(points); i++ {
		root := uf.find(i)
		circuits[root] = append(circuits[root], i)
	}
	//convert to slice to use slice sorting
	circuitSlices := [][]int{}
	for _, points := range circuits {
		circuitSlices = append(circuitSlices, points)
	}
	sort.Slice(circuitSlices, func(i, j int) bool {
		return len(circuitSlices[i]) > len(circuitSlices[j])
	})

	res := len(circuitSlices[0]) * len(circuitSlices[1]) * len(circuitSlices[2])

	fmt.Println(res)
	ufMST := createUnionFind(len(points))
	mst := []Edge{}
	// for part 2 we need MST so we need all edges
	for _, edge := range edges {
		if ufMST.union(edge.FROM_ID, edge.TO_ID) {
			mst = append(mst, edge)
		}

		if len(mst) == len(points)-1 {
			break
		}
	}
	lastEdge := mst[len(mst)-1]
	cableLength := lastEdge.FROM_NODE.X * lastEdge.TO_NODE.X
	fmt.Println(cableLength)
}

func calculatePointDistance(pointA PointNode, pointB PointNode) float64 {
	dx := pointA.X - pointB.X
	dy := pointA.Y - pointB.Y
	dz := pointA.Z - pointB.Z
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return dist
}
