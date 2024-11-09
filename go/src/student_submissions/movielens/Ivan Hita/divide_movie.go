package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

func SplitCSV(fileName string, numPartitions int, directory string) ([]string, error) {
	filePath := fmt.Sprintf("%s/%s", directory, fileName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", directory)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	partitions := make(map[string][][]string)

	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty file: %s", fileName)
		}
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	recordCount := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read record: %w", err)
		}
		assignPartition(numPartitions, recordCount, record, partitions)
		recordCount++
	}

	var wg sync.WaitGroup
	for i := 0; i < numPartitions; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := savePartition(fileName, i, partitions, headers); err != nil {
				log.Printf("failed to save partition %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	var partitionFiles []string
	for key := range partitions {
		partitionFileName := fmt.Sprintf("%s_partition_%s.csv", fileName, key)
		partitionFiles = append(partitionFiles, partitionFileName)
	}
	return partitionFiles, nil
}

func assignPartition(numPartitions, recordIndex int, record []string, partitions map[string][][]string) {
	partitionKey := strconv.Itoa(recordIndex%numPartitions + 1)
	partitions[partitionKey] = append(partitions[partitionKey], record)
}

func savePartition(fileName string, partitionIndex int, partitions map[string][][]string, headers []string) error {
	partitionKey := strconv.Itoa(partitionIndex + 1)
	partitionFileName := fmt.Sprintf("%s_partition_%d.csv", fileName, partitionIndex+1)
	partitionRows := partitions[partitionKey]

	rowsWithHeaders := append([][]string{headers}, partitionRows...)
	file, err := os.Create(partitionFileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", partitionFileName, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(rowsWithHeaders); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", partitionFileName, err)
	}
	writer.Flush()
	return nil
}

func main() {
	// Example usage
	partitionFiles, err := SplitCSV("example.csv", 5, "./data")
	if err != nil {
		log.Fatalf("Error splitting CSV: %v", err)
	}

	fmt.Println("Generated files:", partitionFiles)
}
