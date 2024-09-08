package fileprocessing

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func check_error(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func read_csv(directory string) (datos [][]string) {
	f, err := os.Open(directory)

	check_error(err)

	defer f.Close()

	dataReader := csv.NewReader(f)

	data, err := dataReader.ReadAll()

	check_error(err)

	return data
}

func write_csv(datos [][]string, directory string, filename string) {
	f, err := os.Create(directory + filename)

	check_error(err)

	defer f.Close()

	dataWriter := csv.NewWriter(f)

	dataWriter.WriteAll(datos)

	dataWriter.Flush()

}

func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {

	datos := read_csv(directory + file_name)

	lines_number := len(datos) / number_of_chunks

	var names []string

	var chunk_name string

	for i := 1; i <= number_of_chunks; i++ {

		chunk_name = file_name[:len(file_name)-4] + "_p" + strconv.Itoa(i)+".csv"

		fmt.Print(chunk_name)

		if i == number_of_chunks {
			write_csv(datos[lines_number*(i-1):], directory, chunk_name)

		} else {
			write_csv(datos[lines_number*(i-1):lines_number*(i)], directory, chunk_name)
		}

		names = append(names, chunk_name)

	}
	return names
}


