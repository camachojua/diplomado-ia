package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type RatingObj struct {
	UserId    int64
	MovieId   int64
	Rating    float64
	Timestamp int64
}

func processRecords(records [][]string, outputIndex int, wg *sync.WaitGroup, outputDirectory string) {
	defer wg.Done()

	outputFileName := fmt.Sprintf(outputDirectory+"%02d.csv", outputIndex)
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("Error al crear el archivo %s: %v", outputFileName, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Encabezados
	writer.Write([]string{"UserId", "MovieId", "Rating", "Timestamp"})

	for _, record := range records {
		userId, _ := strconv.ParseInt(record[0], 10, 64)
		movieId, _ := strconv.ParseInt(record[1], 10, 64)
		rating, _ := strconv.ParseFloat(record[2], 64)
		timestamp, _ := strconv.ParseInt(record[3], 10, 64)

		data := RatingObj{
			UserId:    userId,
			MovieId:   movieId,
			Rating:    rating,
			Timestamp: timestamp,
		}

		recordToWrite := []string{
			strconv.FormatInt(data.UserId, 10),
			strconv.FormatInt(data.MovieId, 10),
			strconv.FormatFloat(data.Rating, 'f', 1, 64),
			strconv.FormatInt(data.Timestamp, 10),
		}
		writer.Write(recordToWrite)
	}
}

func SplitBigFile(fileName string, numberOfChunks int, inputDirectory string) []string {
	file, err := os.Open(inputDirectory + fileName)
	if err != nil {
		log.Fatalf("Error al abrir el archivo CSV: %v", err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatalf("Error al leer el archivo CSV: %v", err)
	}

	// Remover encabezados
	if len(records) > 0 && records[0][0] == "UserId" {
		records = records[1:]
	}

	sizeRange := len(records) / numberOfChunks
	var wg sync.WaitGroup
	fileList := make([]string, 0, numberOfChunks)

	for i := 0; i < numberOfChunks; i++ {
		start := i * sizeRange
		end := start + sizeRange
		if end > len(records) {
			end = len(records)
		}

		wg.Add(1)
		go func(index int) {
			processRecords(records[start:end], index, &wg, "./ratings")
			fileList = append(fileList, fmt.Sprintf("ratings_%02d.csv", index))
		}(i + 1)
	}

	wg.Wait()

	return fileList
}

func main() {
	start := time.Now()
	aux := SplitBigFile("ratings.csv", 10, "./ml-25m/")
	fmt.Println("The file names are", aux)
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
}
