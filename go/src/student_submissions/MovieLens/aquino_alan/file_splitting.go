package fileprocessing

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {

	const BACH_FILE_PREFIX = "bachesito"

	// 1. Split big CSV File into Smaller Chucks and save it concurretly for fun

	// Handle the reading of the file with dataframe
	// 1.1 Open File
	file, err := os.Open(file_name)
	if err != nil {
		panic(err)
	} //Found an error
	defer file.Close() //Close at end of function excecution

	total_lines := countTotalLines(file)
	print(total_lines)

	file.Seek(0, io.SeekStart) // Reset file pointer to the beginning

	csvReader := csv.NewReader(file)

	var batch_size int = (total_lines / number_of_chunks)
	fmt.Println(batch_size)

	write_file_dirs := []string{}

	header, err := csvReader.Read()
	if err != nil {
		fmt.Println("Error reading header")
	}

	for i := 0; ; i++ {
		var chunk [][]string
		chunk = append(chunk, header)
		for len(chunk) < batch_size+1 {

			row, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error reading row:", err)
				continue
			}
			chunk = append(chunk, row)
		}
		if len(chunk) < 2 {
			break
		}

		fmt.Println("Processing #", i, "chunk with", len(chunk)-1, "records")
		writeToCsv(i, chunk, directory+BACH_FILE_PREFIX+strconv.Itoa(i))

		write_file_dirs = append(write_file_dirs, BACH_FILE_PREFIX+strconv.Itoa(i)+".csv")

		if err == io.EOF {
			break
		}
	}
	return write_file_dirs
}

func writeToCsv(batch_num int, chunk [][]string, output_dir string) {

	file_dir := output_dir + ".csv"

	outFile, err := os.OpenFile(file_dir, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	csvWriter := csv.NewWriter(outFile)
	defer csvWriter.Flush()

	for i, row := range chunk {
		err = csvWriter.Write(row)
		if err != nil {
			fmt.Println("Error writing row #", i, "of batch", batch_num)
		}
	}

}

func countTotalLines(file *os.File) int {
	scanner := bufio.NewScanner(file)
	totalLines := 0
	for scanner.Scan() {
		totalLines++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error counting lines:", err)
	}
	return totalLines
}
