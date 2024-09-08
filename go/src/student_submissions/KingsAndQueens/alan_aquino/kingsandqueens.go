package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	queens := loadPlayers("./queens.txt")
	kings := loadPlayers("./kings.txt")

	fmt.Println("Initial Queens Preferences")
	fmt.Println(queens, "\n")

	fmt.Println("Initial Kings Preferences")
	fmt.Println(kings, "\n")

	fmt.Println("Starting Sim")
	simulate(kings, queens)
}

func simulate(kings_preferences map[string][]string, queens_preferences map[string][]string) {

	//Map with Queen and King Pairings
	kings_to_queens, queens_to_kings := inizilize_states(kings_preferences, queens_preferences)

	for !(is_all_queens_and_kings_are_matched(kings_to_queens, queens_to_kings)) {
		//Repeat util all kings and kings are paired

		for king, queens := range kings_preferences {
			if kings_to_queens[king] != "UNMATCHED" {
				//If a king is already matched, do not match again, skip ...
				break
			}
			for _, queen := range queens {
				if queen_is_unmatched(queens_to_kings, queen) {
					queens_to_kings[queen] = king
					kings_to_queens[king] = queen
					break
				} else if is_new_partner_hotter(queen, king, queens_to_kings, queens_preferences) {
					previous_partner := queens_to_kings[queen]
					kings_to_queens[previous_partner] = "UNMATCHED"
					queens_to_kings[queen] = king
					kings_to_queens[king] = queen
					break
				}
			}
		}

	}
	fmt.Println("Final Pairings:")
	fmt.Println(queens_to_kings)
	fmt.Println("FINISHED")
}

func is_new_partner_hotter(actual_queen string, candidate_king string, queens_to_kings map[string]string, queens_preferences map[string][]string) bool {

	//Check position of current king in queen king pairing
	actual_queen_preferences := queens_preferences[actual_queen]

	//Who is currently the fortunate

	current_king := queens_to_kings[actual_queen]

	//Set initial values

	candidate_preference_level := -10

	current_king_preference_level := -10

	//retrive levels

	for preference_level, king := range actual_queen_preferences {
		if candidate_king == king {
			candidate_preference_level = preference_level
		}

		if current_king == king {
			current_king_preference_level = preference_level
		}
	}

	if candidate_preference_level < current_king_preference_level {
		//he is hotter, now the current king is the ex
		return true
	}
	return false

}

func is_all_queens_and_kings_are_matched(k_to_q map[string]string, q_to_k map[string]string) bool {
	var is_matched bool
	is_matched = true

	for _, k := range q_to_k {
		if k == "UNMATCHED" {
			is_matched = false
		}
	}

	for _, q := range k_to_q {
		if q == "UNMATCHED" {
			is_matched = false
		}
	}

	return is_matched
}

func queen_is_unmatched(q_to_k map[string]string, queen string) bool {
	if q_to_k[queen] == "UNMATCHED" {
		return true
	}
	return false
}

func inizilize_states(kings map[string][]string, queens map[string][]string) (map[string]string, map[string]string) {
	kings_to_queens := make(map[string]string)
	queens_to_kings := make(map[string]string)

	for king, _ := range kings {
		kings_to_queens[king] = "UNMATCHED"
	}

	for queen, _ := range queens {
		queens_to_kings[queen] = "UNMATCHED"
	}

	return kings_to_queens, queens_to_kings
}

func loadPlayers(fileDirectory string) map[string][]string {

	result := make(map[string][]string)

	file, err := os.Open(fileDirectory)

	if err != nil {
		fmt.Println("Error open file", err)
	}

	defer file.Close() //Ensure the the file os closed after reading

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		elements := strings.Split(line, ",")
		result[elements[0]] = elements[1:]
	}

	return result
}
