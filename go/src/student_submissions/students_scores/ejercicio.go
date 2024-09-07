package main

import (
	"fmt"
)

type student struct {
	Name string
	Score float32
	Grade string
}

type students [] student 

func (ss students) scores (score float32) int {
	for i, v := range ss {
		fmt.Printf("scoring student %s...\n\n", v.Name)
		ss[i].Score = score
	}
	return 0
}

func (ss students) grades() int {
	for i, v := range ss {
		fmt.Printf("grading student %s...\n", v.Name)
		var grade string
		switch {
		case v.Score > 0 && v.Score < 1:
			grade = "X"
		case v.Score >= 1:
			grade = "Y"
		}
		ss[i].Grade = grade
		fmt.Printf("Score=%f, grade=%s\n\n", v.Score, grade)
	}
	return 0
}

func (ss students) printStudents() int {
	for _, v := range ss {
		fmt.Printf("Name: %s, Score: %f, Grade: %s\n",
			v.Name,
			v.Score,
			v.Grade)
	}
	fmt.Printf("\n")
	return 0
}

func main() {
	var s = student{"Sonia Santiago", 0.1, "A"}
	var ss = students{
		s,
	 	{"Alberto Almanza", 0.1, "A"},
	 	{"Bto Almanza", 9.1, "B"},
	}
	
	ss.printStudents()
	ss.scores(1.0)
	ss.printStudents()
	ss.grades()
	ss.printStudents()
}
