package main

import (
	"fmt"
	"os"
	"strings"
)

// calculateSum adds all numbers in a slice
func calculateSum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num // Assignment
	}
	return sum
}

// processData handles various data processing operations
func processData(data []string, options map[string]bool) ([]string, error) {
	if len(data) == 0 { // Condition
		return nil, fmt.Errorf("empty data provided") // Branch (fmt.Errorf)
	}

	results := make([]string, 0, len(data)) // Assignment

	for _, item := range data { // Condition
		if options["uppercase"] { // Condition
			item = strings.ToUpper(item) // Assignment + Branch (strings.ToUpper)
		} else if options["lowercase"] { // Condition
			item = strings.ToLower(item) // Assignment + Branch (strings.ToLower)
		}

		switch {
		case options["trim"]: // Condition
			item = strings.TrimSpace(item) // Assignment + Branch (strings.TrimSpace)
		case options["prefix"] && len(item) > 0: // Condition
			item = "PREFIX_" + item // Assignment
		}

		results = append(results, item) // Assignment + Branch (append)
	}

	return results, nil
}

// main function with complex decision logic
func main() {
	args := os.Args[1:] // Assignment + Branch (os.Args)

	options := map[string]bool{ // Assignment
		"uppercase": false,
		"lowercase": false,
		"trim":      true,
		"prefix":    false,
	}

	numbers := []int{1, 2, 3, 4, 5} // Assignment
	sum := calculateSum(numbers)    // Assignment + Branch (calculateSum)

	if sum > 10 && len(args) > 0 { // Condition
		options["uppercase"] = true // Assignment
	} else if sum <= 10 || len(args) == 0 { // Condition
		options["lowercase"] = true // Assignment
	}

	data := []string{"  Hello  ", "World  ", "  ABC Metrics  "} // Assignment

	results, err := processData(data, options) // Assignment + Branch (processData)
	if err != nil {                            // Condition
		fmt.Println("Error:", err) // Branch (fmt.Println)
		os.Exit(1)                 // Branch (os.Exit)
	}

	for i, result := range results { // Condition
		fmt.Printf("Result %d: %s\n", i+1, result) // Branch (fmt.Printf)
	}
}
