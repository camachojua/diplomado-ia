package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Preference matrices
	proposerPreferences := initProposerPreferences(os.Args[1])
	acceptorPreferences := initAcceptorPreferences(os.Args[2])

	if len(proposerPreferences) != len(acceptorPreferences) {
		panic(errors.New("Matrices don't match."))
	}

	// Helper data structures
	pending := make([]int, len(proposerPreferences))
	next := make([]int, len(proposerPreferences))
	current := make([]int, len(proposerPreferences))
	for i := 0; i < len(proposerPreferences); i++ {
		pending[i] = i  // not-yet-accepted proposers
		next[i] = 0     // next-in-line acceptor (indexed by proposers)
		current[i] = -1 // current match (indexed by acceptors)
	}

	for len(pending) > 0 {
		// Get next in line pending proposer
		p := pending[0]

		// Get p's next to propose to
		var a = proposerPreferences[p][next[p]]

		fmt.Printf("p%d proposing to a%d\n", p, a)

		if current[a] == -1 || acceptorPreferences[a][p] < acceptorPreferences[a][current[a]] {
			// Add current[a] to pending list
			if current[a] > -1 {
				pending = append(pending, current[a])
			}

			// Accept proposal
			current[a] = p
			fmt.Printf("Proposal accepted: (p%d, a%d)\n\n", p, a)

			// Remove p from pending list
			pending = pending[1:]
		}
		next[p]++
	}

	// Print result
	fmt.Println("Final matches:")
	for i := 0; i < len(current); i++ {
		fmt.Printf("(a%d, p%d)\n", i, current[i])
	}
}

func readPreferences(filepath string) [][]int {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// File's first row defines matrix size
	var firstRow []int
	if scanner.Scan() {
		firstRow = line2arr(scanner.Text())
		if len(firstRow) == 0 {
			panic(errors.New("Empty line. Aborting."))
		}
	} else {
		panic(errors.New("Empty file. Aborting."))
	}

	// Create matrix
	matrix := make([][]int, len(firstRow))
	matrix[0] = firstRow

	// Finish reading file
	currRow := 1
	for scanner.Scan() {
		if currRow > len(firstRow) {
			panic(errors.New("Not a squared matrix. Aborting."))
		}

		row := line2arr(scanner.Text())
		if len(row) != len(firstRow) {
			panic(errors.New("Not a squared matrix. Aborting."))
		}

		matrix[currRow] = row
		currRow++
	}
	return matrix
}

func line2arr(line string) []int {
	splitLine := strings.Split(line, " ")
	arr := make([]int, len(splitLine))
	for i, token := range splitLine {
		value, err := strconv.Atoi(token)
		if err != nil {
			panic(err)
		}
		arr[i] = value
	}
	return arr
}

func initProposerPreferences(filepath string) [][]int {
	return readPreferences(filepath)
}

func initAcceptorPreferences(filepath string) [][]int {
	x := readPreferences(filepath)

	xp := make([][]int, len(x))
	for i := range x {
		xp[i] = make([]int, len(x))
		for j, jval := range x[i] {
			xp[i][jval] = j
		}
	}

	return xp
}
