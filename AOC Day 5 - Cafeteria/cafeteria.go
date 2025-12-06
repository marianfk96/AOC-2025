package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("filepath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	ranges := []string{}
	ingredients := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// range
		if strings.Contains(line, "-") {
			ranges = append(ranges, line)
		} else {
			ingredients = append(ingredients, line)
		}
	}
	freshIngredients := []string{}
	fresh := 0
	for _, ingr := range ingredients {
		id, _ := strconv.Atoi(ingr)
		if isFresh(id, ranges) {
			freshIngredients = append(freshIngredients, ingr)
			fresh++
		}
	}
	sortedRanges := [][]int{}
	for i := 0; i < len(ranges); i++ {
		r1, r2 := rangeToInts(ranges[i])
		sortedRanges = append(sortedRanges, []int{r1, r2})
	}
	// sort the ranges after creating the 2 d array, to make sure i < i+2 < i+4 etc etc
	sort.Slice(sortedRanges, func(i int, j int) bool {
		return sortedRanges[i][0] < sortedRanges[j][0]
	})
	flattenedRanges := []int{sortedRanges[0][0], sortedRanges[0][1]}
	for i := 1; i < len(sortedRanges); i++ {
		flattenedRanges = flattenAndMergeRanges(sortedRanges[i][0], sortedRanges[i][1], flattenedRanges)
	}
	possibleFresh := countPotentiallyFreshIds(flattenedRanges)
	fmt.Println(possibleFresh)
}

// full on tilt :)
func flattenAndMergeRanges(r1 int, r2 int, flattenedRanges []int) []int {
	for i := 0; i < len(flattenedRanges); i += 2 {
		// is completely covered -> don't add anywhere, break the loop .
		if r1 >= flattenedRanges[i] && r2 <= flattenedRanges[i+1] {
			break
		}
		//completely engulfs current range -> replacr it and check next ranges for potential removals
		if r1 <= flattenedRanges[i] && r2 >= flattenedRanges[i+1] {
			flattenedRanges[i] = r1
			flattenedRanges[i+1] = r2
			// find where the overlap ends, completely removing intermiediate ranges
			overlapEnd := findOverlapEnd(r2, i+2, flattenedRanges)
			flattenedRanges = append(flattenedRanges[:i+2], flattenedRanges[overlapEnd:]...)
			break
		}
		// extends to the righg
		if r1 >= flattenedRanges[i] && r1 <= flattenedRanges[i+1] && r2 > flattenedRanges[i+1] {
			flattenedRanges[i+1] = r2
			overlapEnd := findOverlapEnd(r2, i+2, flattenedRanges)
			flattenedRanges = append(flattenedRanges[:i+2], flattenedRanges[overlapEnd:]...)
			break
		}
		// no overlap
		if r2 < flattenedRanges[i] {
			newEntries := []int{r1, r2}
			flattenedRanges = append(flattenedRanges[:i], append(newEntries, flattenedRanges[i:]...)...)
			break
		}
		if r1 > flattenedRanges[i+1] {
			// keep checking -> it'll overlap in next iterations and will be handled
			if i+2 < len(flattenedRanges) {
				continue
			}
			// if we reached the end, we can append
			flattenedRanges = append(flattenedRanges, r1, r2)
			break
		}
	}
	return flattenedRanges
}

func findOverlapEnd(r2 int, startignIdx int, flattenedRanges []int) int {
	for i := startignIdx; i < len(flattenedRanges); i++ {
		if flattenedRanges[i] > r2 {
			return i
		}
	}
	// r2 is larger than all our ranges
	return len(flattenedRanges)
}

func countPotentiallyFreshIds(sortedRanges []int) int {
	sum := 0
	for i := 0; i < len(sortedRanges)-1; i += 2 {
		// formula is r - l + 1
		gap := sortedRanges[i+1] - sortedRanges[i] + 1
		sum += gap

	}
	return sum
}

func isFresh(id int, ranges []string) bool {
	for _, r := range ranges {
		r1, r2 := rangeToInts(r)
		if r1 <= id && r2 >= id {
			return true
		}
	}
	return false
}

func rangeToInts(r string) (int, int) {
	parts := strings.Split(r, "-")
	r1, _ := strconv.Atoi(parts[0])
	r2, _ := strconv.Atoi(parts[1])
	return r1, r2
}
