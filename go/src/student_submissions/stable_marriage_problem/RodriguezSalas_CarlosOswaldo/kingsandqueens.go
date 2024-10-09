package main

import (
	"fmt"
)

type Queen struct {
	Match       string
	Preferences []string
	Proposals   map[string]bool
}

type King struct {
	Match       string
	Preferences []string
	Proposals   map[string]bool
}

func kingsandqueens(kings map[string]*King, kings_names []string, queens map[string]*Queen, queens_names []string) (map[string]*King, map[string]*Queen) {

	s := 1

	for s > 0 {

		s = 0

		for _, n := range kings_names {
			if (*kings[n]).Match == "" {

				for _, q := range (*kings[n]).Preferences {
					if !(*kings[n]).Proposals[q] {
						(*kings[n]).Proposals[q] = true
						(*queens[q]).Proposals[n] = true
						break
					}
				}
			}
		}

		for _, m := range queens_names {
			for _, k := range (*queens[m]).Preferences {
				if (*queens[m]).Proposals[k] {

					if _, ok := kings[(*queens[m]).Match]; ok {
						(*kings[(*queens[m]).Match]).Match = ""
					}

					(*queens[m]).Match = k
					(*kings[k]).Match = m
					break
				}
			}

		}

		for _, k := range kings {
			if (*k).Match == "" {
				s++
			}
		}

		for _, q := range queens {
			if (*q).Match == "" {
				s++
			}
		}

	}

	return kings, queens
}

func main() {

	queens_names := []string{"r01", "r02", "r03", "r04", "r05", "r06", "r07", "r08", "r09", "r10"}
	kings_names := []string{"s01", "s02", "s03", "s04", "s05", "s06", "s07", "s08", "s09", "s10"}

	q_preferences := [][]string{
		{"s02", "s06", "s10", "s07", "s09", "s01", "s04", "s05", "s03", "s08"},
		{"s02", "s01", "s03", "s06", "s07", "s04", "s09", "s05", "s10", "s08"},
		{"s06", "s02", "s05", "s07", "s08", "s03", "s09", "s01", "s04", "s10"},
		{"s06", "s10", "s03", "s01", "s09", "s08", "s07", "s04", "s02", "s05"},
		{"s10", "s08", "s06", "s04", "s01", "s07", "s03", "s05", "s09", "s02"},
		{"s02", "s01", "s05", "s09", "s10", "s04", "s06", "s07", "s03", "s08"},
		{"s10", "s07", "s08", "s06", "s02", "s01", "s03", "s05", "s04", "s09"},
		{"s07", "s10", "s02", "s01", "s09", "s04", "s08", "s05", "s03", "s06"},
		{"s09", "s03", "s08", "s07", "s06", "s02", "s01", "s05", "s10", "s04"},
		{"s05", "s08", "s07", "s01", "s02", "s10", "s03", "s09", "s06", "s04"}}

	k_preferences := [][]string{
		{"r01", "r05", "r03", "r09", "r10", "r04", "r06", "r02", "r08", "r07"},
		{"r03", "r08", "r01", "r04", "r05", "r06", "r02", "r10", "r09", "r07"},
		{"r08", "r05", "r01", "r04", "r02", "r06", "r09", "r07", "r03", "r10"},
		{"r09", "r06", "r04", "r07", "r08", "r05", "r10", "r02", "r03", "r01"},
		{"r10", "r04", "r02", "r03", "r06", "r05", "r01", "r09", "r08", "r07"},
		{"r02", "r01", "r04", "r07", "r05", "r09", "r03", "r10", "r08", "r06"},
		{"r07", "r05", "r09", "r02", "r03", "r01", "r04", "r08", "r10", "r06"},
		{"r01", "r05", "r08", "r06", "r09", "r03", "r10", "r02", "r07", "r04"},
		{"r08", "r03", "r04", "r07", "r02", "r01", "r06", "r09", "r10", "r05"},
		{"r01", "r06", "r10", "r07", "r05", "r02", "r04", "r03", "r09", "r08"}}

	kings := make(map[string]*King, len(kings_names))
	queens := make(map[string]*Queen, len(queens_names))

	for i, name := range kings_names {
		var k King
		k.Proposals = make(map[string]bool)
		k.Preferences = k_preferences[i]
		kings[name] = &k
	}
	for i, name := range queens_names {
		var q Queen
		q.Proposals = make(map[string]bool)
		q.Preferences = q_preferences[i]
		queens[name] = &q
	}

	kings, queens = kingsandqueens(kings, kings_names, queens, queens_names)

	for i, k := range kings {
		fmt.Println(i + "-" + k.Match)
	}

	for i, k := range queens {
		fmt.Println(i + "-" + k.Match)
	}

}
