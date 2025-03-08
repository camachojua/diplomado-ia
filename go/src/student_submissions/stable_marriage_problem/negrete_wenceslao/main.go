// Package main implements the Gale-Shapley stable marriage algorithm
package main

import (
	"fmt"
	"log"
)

// Preferences represents a mapping of individuals to their ranked preferences
type Preferences map[string][]string

// StableMarriage implements the Gale-Shapley algorithm to find stable marriages
type StableMarriage struct {
	kings         Preferences
	queens        Preferences
	kingStacks    map[string]*stack
	queenStacks   map[string]*stack
	matches       map[string]string
	freeKings     *stack
	debugEnabled  bool
}

// NewStableMarriage creates a new instance of the stable marriage solver
func NewStableMarriage(kings, queens Preferences, debug bool) *StableMarriage {
	sm := &StableMarriage{
		kings:         kings,
		queens:        queens,
		kingStacks:    make(map[string]*stack),
		queenStacks:   make(map[string]*stack),
		matches:       make(map[string]string),
		freeKings:     NewStack(),
		debugEnabled:  debug,
	}

	// Initialize preference stacks
	sm.initializeStacks()
	return sm
}

func (sm *StableMarriage) debugPrint(format string, args ...interface{}) {
	if sm.debugEnabled {
		fmt.Printf(format, args...)
	}
}

func (sm *StableMarriage) initializeStacks() {
	// Initialize king preference stacks
	for king, preferences := range sm.kings {
		sm.kingStacks[king] = NewStack()
		for i := len(preferences) - 1; i >= 0; i-- {
			sm.kingStacks[king].push(preferences[i])
		}
	}

	// Initialize queen preference stacks
	for queen, preferences := range sm.queens {
		sm.queenStacks[queen] = NewStack()
		for i := len(preferences) - 1; i >= 0; i-- {
			sm.queenStacks[queen].push(preferences[i])
		}
	}

	// Initialize free kings
	for king := range sm.kings {
		sm.freeKings.push(king)
	}
}

func (sm *StableMarriage) queenPrefersCurrentKing(matchedKing, currentKing string, queenPreferences *stack) bool {
	// Create a clone to avoid modifying the original stack
	preferences := queenPreferences.Clone()
	
	for !preferences.isEmpty() {
		preferredKing, err := preferences.pop()
		if err != nil {
			log.Fatal(err)
		}
		if preferredKing == matchedKing {
			return false
		}
		if preferredKing == currentKing {
			return true
		}
	}
	return false
}

// Solve runs the Gale-Shapley algorithm and returns the stable matches
func (sm *StableMarriage) Solve() map[string]string {
	for !sm.freeKings.isEmpty() {
		sm.debugPrint("===> Free Kings: %v\n", sm.freeKings.print())
		
		currentKing, err := sm.freeKings.pop()
		if err != nil {
			log.Fatal(err)
		}
		sm.debugPrint("===> Current King: %v\n", currentKing)

		for !sm.kingStacks[currentKing].isEmpty() {
			preferredQueen, err := sm.kingStacks[currentKing].pop()
			if err != nil {
				log.Fatal(err)
			}
			sm.debugPrint("===> Preferred Queen: %v\n", preferredQueen)

			matchedKing, matched := sm.matches[preferredQueen]
			if !matched {
				sm.debugPrint("Queen is free, accept the proposal\n")
				sm.matches[preferredQueen] = currentKing
				sm.debugPrint("Matches: %v\n", sm.matches)
				break
			} else {
				sm.debugPrint("Queen has been matched\n")
				if sm.queenPrefersCurrentKing(matchedKing, currentKing, sm.queenStacks[preferredQueen]) {
					sm.debugPrint("Queen prefers the current king, previously matched king %v returns to free kings\n", matchedKing)
					sm.matches[preferredQueen] = currentKing
					sm.freeKings.push(matchedKing)
					break
				} else {
					sm.debugPrint("Queen prefers the previously matched king, current king keeps proposing\n")
				}
			}
		}
	}
	return sm.matches
}

func main() {
	kings := Preferences{
		"r01": []string{"s02", "s06", "s10", "s07", "s09", "s01", "s04", "s05", "s03", "s08"},
		"r02": []string{"s02", "s01", "s03", "s06", "s07", "s04", "s09", "s05", "s10", "s08"},
		"r03": []string{"s06", "s02", "s05", "s07", "s08", "s03", "s09", "s01", "s04", "s10"},
		"r04": []string{"s06", "s10", "s03", "s01", "s09", "s08", "s07", "s04", "s02", "s05"},
		"r05": []string{"s10", "s08", "s06", "s04", "s01", "s07", "s03", "s05", "s09", "s02"},
		"r06": []string{"s02", "s01", "s05", "s09", "s10", "s04", "s06", "s07", "s03", "s08"},
		"r07": []string{"s10", "s07", "s08", "s06", "s02", "s01", "s03", "s05", "s04", "s09"},
		"r08": []string{"s07", "s10", "s02", "s01", "s09", "s04", "s08", "s05", "s03", "s06"},
		"r09": []string{"s09", "s03", "s08", "s07", "s06", "s02", "s01", "s05", "s10", "s04"},
		"r10": []string{"s05", "s08", "s07", "s01", "s02", "s10", "s03", "s09", "s06", "s04"},
	}

	queens := Preferences{
		"s01": []string{"r01", "r05", "r03", "r09", "r10", "r04", "r06", "r02", "r08", "r07"},
		"s02": []string{"r03", "r08", "r01", "r04", "r05", "r06", "r02", "r10", "r09", "r07"},
		"s03": []string{"r08", "r05", "r01", "r04", "r02", "r06", "r09", "r07", "r03", "r10"},
		"s04": []string{"r09", "r06", "r04", "r07", "r08", "r05", "r10", "r02", "r03", "r01"},
		"s05": []string{"r10", "r04", "r02", "r03", "r06", "r05", "r01", "r09", "r08", "r07"},
		"s06": []string{"r02", "r01", "r04", "r07", "r05", "r09", "r03", "r10", "r08", "r06"},
		"s07": []string{"r07", "r05", "r09", "r02", "r03", "r01", "r04", "r08", "r10", "r06"},
		"s08": []string{"r01", "r05", "r08", "r06", "r09", "r03", "r10", "r02", "r07", "r04"},
		"s09": []string{"r08", "r03", "r04", "r07", "r02", "r01", "r06", "r09", "r10", "r05"},
		"s10": []string{"r01", "r06", "r10", "r07", "r05", "r02", "r04", "r03", "r09", "r08"},
	}

	sm := NewStableMarriage(kings, queens, true)
	matches := sm.Solve()
	
	fmt.Println("\nFinal Matches:")
	for queen, king := range matches {
		fmt.Printf("%s ♔ %s\n", king, queen)
	}
}
