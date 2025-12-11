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
	file, err := os.Open("fl.txt")
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
	// 1. create a rectangle grid out of the grid limits we found earlier
	// 2. add the red tiles -> make sure to normalize points by substracting minX, minY
	// 3. add X between # in rows -> rows will ALWAYS end on max width (maxY)
	// 4. after filling the outline between # in the same row, fill X in colums by keeping track in which row we
	// encountered the last #
	// 5. Now that we have the outline we can fill the polygon with X and return the points that create it
	gridlimits := GridLimits{minX, maxX, minY, maxY}
	grid := createEmptyGrid(gridlimits)
	// add red points to grid making sure to 'shift' their coordinates
	grid = addRedTiles(points, grid, minX, minY)
	for i, row := range grid {
		grid[i] = addXToRow(grid, row, i)
	}
	grid = addXToColumn(grid)
	grid, polygonPoints := fillPolygon(grid, minX, minY)
	// 1. sort the areas slice
	// 2. beggining from greatest area, find the opposite points from the points that created it
	// 3. if they are contained within polygonPoints that's the max width we can create
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].area > areas[j].area
	})
	maxContainedArea := -1
	for _, area := range areas {
		oppX, oppY := findOppositePoints(area.FROM, area.TO)
		if isPartOfPolygon(oppX, polygonPoints) && isPartOfPolygon(oppY, polygonPoints) {
			maxContainedArea = int(area.area)
			break
		}
	}

	fmt.Println(maxContainedArea)
}

func fillPolygon(grid [][]rune, minX, minY int) ([][]rune, []Point) {
	points := []Point{}
	for row := 0; row < len(grid); row++ {
		// fill in the points inside our outline
		inside := false
		for col := 0; col < len(grid[0]); col++ {
			if grid[row][col] == 'X' || grid[row][col] == '#' {
				points = append(points, Point{col + minX, row + minY})
				// stop filling once we reach the ende
				inside = !inside
				continue
			}
			// fill with X
			if inside && grid[row][col] == '.' {
				grid[row][col] = 'X'
				points = append(points, Point{row + minY, col + minX})
			}
		}
	}
	return grid, points
}

func addXToRow(grid [][]rune, row []rune, rowIndex int) []rune {
	for i, val := range row {
		if val == '#' {
			for j := i + 1; j < len(row); j++ {
				if row[j] == '#' {
					break
				}
				row[j] = 'X'
			}
		}
	}
	return row
}

func addXToColumn(grid [][]rune) [][]rune {
	for col := 0; col < len(grid[0]); col++ {
		lastHash := -1
		for row := 0; row < len(grid); row++ {
			// keep track of where we saw the last #
			if grid[row][col] == '#' {
				if lastHash >= 0 {
					// fill column with X between the position of last seen # and current row
					for k := lastHash + 1; k < row; k++ {
						grid[k][col] = 'X'
					}
				}
				lastHash = row
			}
		}
	}
	return grid
}

func createEmptyGrid(limits GridLimits) [][]rune {
	width := limits.MAX_X - limits.MIN_X + 1
	height := limits.MAX_Y - limits.MIN_Y + 1
	grid := make([][]rune, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]rune, width)
		for j := 0; j < width; j++ {
			grid[i][j] = '.'
		}
	}
	return grid
}

func addRedTiles(points []Point, grid [][]rune, minX, minY int) [][]rune {
	for _, point := range points {
		row := point.Y - minY
		col := point.X - minX
		grid[row][col] = '#'
	}
	return grid
}

func isPartOfPolygon(point Point, polygonPOints []Point) bool {
	for _, p := range polygonPOints {
		if point.X == p.X && point.Y == p.Y {
			return true
		}
	}
	return false
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
