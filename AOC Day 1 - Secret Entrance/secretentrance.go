package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var times_at_zero int = 0
	var current_position int = 50
	var passes_through_zero int = 0
	file, err := os.Open("filepath")
	if err != nil {
		fmt.Println("Error reading the file")
		return
	}
	defer file.Close()
	var positions []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		positions = append(positions, scanner.Text())
	}

	for _, pos := range positions {
		var direction string = string(pos[0])
		moves, _ := strconv.Atoi(pos[1:])
		if direction == "R" {
			distance_to_zero := 100 - current_position
			if distance_to_zero == 0 {
				distance_to_zero = 100 // 100 moves away from next 0
			}
			if moves >= distance_to_zero {
				passes_through_zero += 1 + (moves-distance_to_zero)/100
			}
			current_position = (current_position + moves) % 100
		} else {
			distance_to_zero := current_position
			if distance_to_zero == 0 {
				distance_to_zero = 100 // 100 moves away from next 0
			}
			if moves >= distance_to_zero {
				passes_through_zero += 1 + (moves-distance_to_zero)/100
			}
			current_position = (current_position - moves) % 100
			if current_position < 0 {
				current_position += 100
			}
		}
		if current_position == 0 {
			times_at_zero++
		}
	}
	fmt.Println("Roations around 0: ", passes_through_zero)
}
