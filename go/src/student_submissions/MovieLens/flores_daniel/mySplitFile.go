package fileprocessing

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func getBatch(reader *csv.Reader, number_of_chunks int) (int, int) {
	rowCount := 0
	for {
		_, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal("Error reading record: ", err)
		}
		rowCount++
	}

	return int(math.Ceil(float64(rowCount) / float64(number_of_chunks))), rowCount
}

func mySplitFile(file_name_ string, number_of_chunks int, directory string) []string {

	// Read files

	file, err := os.Open(directory + file_name_)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Getting the size per file
	batch, len := getBatch(reader, number_of_chunks)
	fmt.Println("Registros por archivo " + strconv.Itoa(batch))
	fmt.Println("Registros del archivo archivo original " + strconv.Itoa(len))
	file.Close()

	// Read de file again to recover the pointer
	file, err = os.Open(directory + file_name_)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()
	reader = csv.NewReader(file)

	// Setting varibales for the loop
	count := 0
	batch_count := 0
	file_number := 0
	var files []string
	file_name := "ratings_" + strconv.Itoa(file_number) + ".csv"

	//File for first loop
	files = append(files, file_name)
	file_name = directory + file_name
	csvFile, err := os.Create(file_name)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	for count < len {

		// Read de old file and write the new
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal("Error reading record: ", err)
		}

		if err := csvwriter.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}

		count++
		batch_count++

		// If the file is full, close it and create a new one
		if batch_count == batch {
			csvwriter.Flush()
			csvFile.Close()
			fmt.Println("El archivo " + file_name + " ha sido creado con éxito")

			file_number++

			// This avoid create a new wmpty file in the end of the process
			if file_number != number_of_chunks {
				file_name = "ratings_" + strconv.Itoa(file_number) + ".csv"
				files = append(files, file_name)
				file_name = directory + file_name
				csvFile, err = os.Create(file_name)
				if err != nil {
					log.Fatalf("failed creating file: %s", err)
				}

				csvwriter = csv.NewWriter(csvFile)
			}
			batch_count = 0
		}
	}

	// To ensure that the las file is closed
	if batch_count > 0 {
		csvwriter.Flush()
		csvFile.Close()
		fmt.Println("El archivo " + file_name + " ha sido creado con éxito (último archivo)")
	}

	return files
}
