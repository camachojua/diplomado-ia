package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {

	file, err := os.Open(directory + file_name)

	if err != nil {
		fmt.Println("ora")
	}

	reader := csv.NewReader(file)

	dataset, err := reader.ReadAll()

	if err != nil {
		println("ora2")
	}

	datalength := len(dataset)

	step := int(math.Round(float64(datalength / number_of_chunks)))

	// a := dataset[:100][:]

	// println(a, "khe")

	var sliceToPrint [][]string

	var end int

	var filenameToWrite string

	var inResult []string

	for i := 0; i < number_of_chunks; i++ {

		if i == number_of_chunks-1 {
			end = int(math.Max(float64(datalength), float64(step*(i+1))))
		} else {
			end = step * (i + 1)
		}
		sliceToPrint = dataset[i*step : end][:]

		filenameToWrite = file_name + strconv.Itoa(i) + ".csv"

		fileToWrite, err := os.Create(directory + filenameToWrite)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		csv.NewWriter(fileToWrite).WriteAll(sliceToPrint)

		fileToWrite.Close()

		inResult = append(inResult, filenameToWrite)
	}

	return inResult
}

func main() {
	SplitBigFile("ratings.csv", 10, "./")
}
