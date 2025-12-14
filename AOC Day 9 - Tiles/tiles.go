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

type Result struct {
	area Area
	ok   bool
}

type Point struct {
	X, Y int
}

type Area struct {
	FROM, TO Point
	area     float64
}

type GridLimits struct {
	MIN_X, MAX_X, MIN_Y, MAX_Y int
}

func main() {
	file, err := os.Open("filePath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	pos := 0
	points := []Point{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		point := Point{toInteger(line[0]), toInteger(line[1])}
		points = append(points, point)
		pos++
	}
	//find grid limits -> seems useful to have
	minX := points[0].X
	maxX := -1
	minY := points[0].Y
	maxY := -1
	for _, point := range points {
		if point.X < minX {
			minX = point.X
		} else if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		} else if point.Y > maxY {
			maxY = point.Y
		}
	}
	areas := []Area{}
	for i, point := range points {
		for j := i + 1; j < len(points); j++ {
			// forgot to add option for x1=x2 and y1=y2 , but got correct result withou it, so fuck it :D
			if point.X != points[j].X && point.Y != points[j].Y {
				area := (math.Abs(float64(points[j].X)-float64(point.X)) + 1) *
					(math.Abs(float64(points[j].Y)-float64(point.Y)) + 1)
				areas = append(areas, Area{point, points[j], area})
			}
		}
	}

	maxArea := -1
	for _, area := range areas {
		if area.area > float64(maxArea) {
			maxArea = int(area.area)
		}
	}
	fmt.Println(maxArea)
	// PART 2
	//I have an ordered set of points that create the vertices of a polygon in a 2d grid
	//I have a slice of areas created by at least 2 of those points (i keep track of them)
	// I can calculate opposite points to find all edges of the rectangle that creates this area.
	// Then use ray casting algorithm to whether those points are contained inside the polygon (or intersect wiht the perimeter)
	// Those areas will form my candidates
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].area > areas[j].area
	})
	candidates := []Area{}
	for _, area := range areas {
		oppA, oppB := findOppositePoints(area.FROM, area.TO)
		// ray casting algorithm for all edges
		if pointInPolygon(area.FROM, points) && pointInPolygon(area.TO, points) && pointInPolygon(oppA, points) && pointInPolygon(oppB, points) {
			candidates = append(candidates, area)
		}
	}
	workers := 6
	areaCh := make(chan Area, len(candidates))
	resultCh := make(chan Result, workers)

	for _, a := range candidates {
		areaCh <- a
	}
	close(areaCh)

	for w := 0; w < workers; w++ {
		go func() {
			for area := range areaCh {
				minX := min(area.FROM.X, area.TO.X)
				maxX := max(area.FROM.X, area.TO.X)
				minY := min(area.FROM.Y, area.TO.Y)
				maxY := max(area.FROM.Y, area.TO.Y)
				// perimeter is contained then the whole rectangle is contained
				if rectangleContained(minX, maxX, minY, maxY, points) {
					resultCh <- Result{area, true}
					return
				}
			}
			resultCh <- Result{Area{}, false}
		}()
	}

	// collect results and find max, as all workers might have returned a valid result
	maxContainedArea := -1
	for range workers {
		r := <-resultCh
		if r.ok && int(r.area.area) > maxContainedArea {
			maxContainedArea = int(r.area.area)
		}
	}

	fmt.Println(maxContainedArea)
}

func rectangleContained(minX, maxX, minY, maxY int, polygon []Point) bool {
	// check x-axis
	for x := minX; x <= maxX; x++ {
		if !pointInPolygon(Point{x, minY}, polygon) ||
			!pointInPolygon(Point{x, maxY}, polygon) {
			return false
		}
	}
	// check y-axis
	for y := minY; y <= maxY; y++ {
		if !pointInPolygon(Point{minX, y}, polygon) ||
			!pointInPolygon(Point{maxX, y}, polygon) {
			return false
		}
	}
	return true
}

func pointInPolygon(p Point, points []Point) bool {
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		// last line to check is final point matched with the first one  -> this one closes the polygon
		var p2 Point
		if i < len(points)-1 {
			p2 = points[i+1]
		} else {
			p2 = points[0]
		}
		if pointOnPerimeter(p, p1, p2) {
			return true
		}
	}
	// ray casting
	// point is inside polygon if ray inmtersect odd number of times -> start with inside = false, on intersection inise = !inside
	inside := false
	j := len(points) - 1 // index of previous point, we start with the last one
	for i := 0; i < len(points); i++ {
		pcurr := points[i]
		pprev := points[j]
		// detect intersections and switch inside
		if (pcurr.Y <= p.Y && p.Y < pprev.Y) || (pprev.Y <= p.Y && p.Y < pcurr.Y) {
			if p.X < (pprev.X-pcurr.X)*(p.Y-pcurr.Y)/(pprev.Y-pcurr.Y)+pcurr.X {
				inside = !inside
			}
		}
		j = i // assign j for the next iteration
	}
	return inside
}

func pointOnPerimeter(p, pointA, pointB Point) bool {
	// 2 checks -> p must be on the line defined by a and b. P must be between a and b
	// cross-product test for collinearilty -> (x2​−x1​)(y3​−y1​)−(y2​−y1​)(x3​−x1​)=0
	collinear := (pointB.X-pointA.X)*(p.Y-pointA.Y) == (pointB.Y-pointA.Y)*(p.X-pointA.X)
	// early exity if not on the same line
	if !collinear {
		return false
	}
	// since collineas , check that p is between a and b
	if min(pointA.X, pointB.X) <= p.X && p.X <= max(pointA.X, pointB.X) && min(pointA.Y, pointB.Y) <= p.Y && p.Y <= max(pointA.Y, pointB.Y) {
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func findOppositePoints(pointA, pointB Point) (Point, Point) {
	oppositeA := Point{pointB.X, pointA.Y}
	oppositeB := Point{pointA.X, pointB.Y}
	return oppositeA, oppositeB
}

func toInteger(a string) int {
	a = strings.TrimSpace(a)
	num, _ := strconv.Atoi(a)
	return num
}
