package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"sync"
)

func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {
	var filesCreated []string

	//Inicia variable tiempo
	fmt.Println("Inicia la variable del tiempo ")
	t1 := time.Now()
	fmt.Println("Tiempo Inicial es :",t1)
	data := readCSVFile(directory + file_name)
	rowsPerChunk := getRowsPerChunk(len(data), number_of_chunks)
	fmt.Println("Rows in file: ", len(data))
	fmt.Println("Rows per chunk: ", rowsPerChunk)
	i := 0
	var wg sync.WaitGroup
	for i < number_of_chunks-1 {
		chunkName := "archivo_salida_" + strconv.Itoa(i) + ".csv"
		chunkPath := directory + chunkName
		chunkData := data[(rowsPerChunk * i):(rowsPerChunk * (i + 1))]
		wg.Add(1)

		go func() {
			defer wg.Done()
			csvFile, err := os.Create(chunkPath)
			if err != nil {
				log.Fatalf("Falla / error : %s", err)
			}
			defer csvFile.Close()
			writer := csv.NewWriter(csvFile)
			defer writer.Flush()
			err = writer.WriteAll(chunkData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(len(chunkData), " Registros_Totales ", chunkPath)

	    }()

		filesCreated = append(filesCreated, chunkName)
		i += 1
	}
	chunkName := "archivo_salida_" + strconv.Itoa(i) + ".csv"
	chunkPath := directory + chunkName
	chunkData := data[rowsPerChunk*i:]
    wg.Add(1)
	go func ()  {
		defer wg.Done()
		csvFile, err := os.Create(chunkPath)
		if err != nil {
			log.Fatalf("Falla / error : %s", err)
		}
		defer csvFile.Close()
		writer := csv.NewWriter(csvFile)
		defer writer.Flush()
		err = writer.WriteAll(chunkData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(chunkData), " Registros_Totales ", chunkPath)

	}()

	filesCreated = append(filesCreated, chunkName)

	// Finaliza el conteo del tiempo
	fmt.Println("Finaliza la variable del tiempo ")
	now := time.Now()
	diff:= now.Sub(t1)
	fmt.Println("Tiempo Final en segundos es: ",diff.Seconds())
	wg.Wait()	
	return filesCreated
}

func getRowsPerChunk(numberOfRows int, numberOfChunks int) int {
	var rowsPerChunk int
	rowsPerChunk = numberOfRows / numberOfChunks
	return rowsPerChunk
}

func readCSVFile(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	return data
}

func main() {
	SplitBigFile("ratings.csv", 100, "./test/")
}
