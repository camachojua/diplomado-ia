package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Function to split a csv file into small files.
// You provided a name for the file, the number of chuncks wich it will be divided and the directory
// where the file is located and the new ones will be saved.
func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {

	data := ReadCsv(file_name, directory)

	fmt.Printf("%v rows in file %s\n", len(data), file_name)
	rowsPerFile := len(data) / number_of_chunks
	var filesCreated []string

	for i := 0; i < number_of_chunks; i++ {
		tempName := file_name + "_part_" + strconv.Itoa(i)
		path := directory + tempName
		tempData := data[i*rowsPerFile : (i+1)*rowsPerFile]
		WriteCsv(tempData, tempName, path)
		filesCreated = append(filesCreated, tempName)
	}

	return filesCreated
}

// Open and read a csv file and returns the content.
func ReadCsv(fileName string, directory string) [][]string {
	file, err := os.Open(directory + fileName + ".csv")

	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	data, err := csvReader.ReadAll()

	if err != nil {
		log.Fatalf("Error extracting data from file %v: %s", fileName, err)
	}
	return data
}

// Create a csv file with the name and data provided in the path defined.
func WriteCsv(data [][]string, name string, path string) {
	csvFile, err := os.Create(path + name + ".csv")
	if err != nil {
		log.Fatalf("Error creating new csv file %v: %s", name, err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	err = writer.WriteAll(data)
	if err != nil {
		log.Fatalf("Error writing new csv file %v: %s", name, err)
	}

	fmt.Printf("File %s has been created with %v rows\n", name, len(data))
}

func main() {
	files := SplitBigFile("ratings", 10, "/mnt/c/Users/omarjh/Documents/Diplomado_IA/ejercicios/")
	fmt.Println(files)
}
