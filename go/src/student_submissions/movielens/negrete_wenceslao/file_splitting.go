// Package fileprocessing provides functionality to split large CSV files into smaller chunks
// while preserving the header row in each partition.
package fileprocessing

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Partition represents a collection of CSV records
type Partition struct {
	Records [][]string
}

// AddRecord adds a new record to the partition
func (p *Partition) AddRecord(record []string) {
	p.Records = append(p.Records, record)
}

// Len returns the number of records in the partition
func (p *Partition) Len() int {
	return len(p.Records)
}

// SplitBigFile splits a CSV file into a specified number of chunks.
// It takes the filename, desired number of chunks, and target directory as parameters.
// Returns a slice of partition filenames and any error that occurred.
func SplitBigFile(fileName string, numberOfChunks int, directory string) ([]string, error) {
	if numberOfChunks <= 0 {
		return nil, fmt.Errorf("number of chunks must be positive, got %d", numberOfChunks)
	}

	// Validate directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", directory)
	}

	filePath := filepath.Join(directory, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	partitions := make(map[int]*Partition)
	recordNum := 0

	headers, err := csvReader.Read()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("reading headers: %w", err)
	}

	// Read and distribute records
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading record: %w", err)
		}
		
		if err := distributeRecord(numberOfChunks, recordNum, record, partitions); err != nil {
			return nil, fmt.Errorf("distributing record: %w", err)
		}
		recordNum++
	}

	// Create buffered channels for error handling and synchronization
	done := make(chan error, len(partitions))
	partitionFiles := make([]string, len(partitions))
	
	// Launch goroutines to save each partition
	for partitionNum := range partitions {
		partitionFiles[partitionNum] = fmt.Sprintf("%s_partition_%d.csv", fileName, partitionNum+1)
		go savePartition(filePath, partitionNum, partitions[partitionNum], headers, done)
	}

	// Wait for all partitions to complete and collect any errors
	for i := 0; i < len(partitions); i++ {
		if err := <-done; err != nil {
			return nil, fmt.Errorf("saving partition %d: %w", i, err)
		}
	}

	return partitionFiles, nil
}

// distributeRecord assigns a record to its corresponding partition based on the record number
// and desired number of chunks.
func distributeRecord(numberOfChunks, recordNum int, record []string, partitions map[int]*Partition) error {
	// Calculate partition number (0 to numberOfChunks-1)
	partitionNum := recordNum % numberOfChunks
	
	if partitions[partitionNum] == nil {
		partitions[partitionNum] = &Partition{
			Records: make([][]string, 0, 1000), // Pre-allocate space for better performance
		}
	}
	
	partitions[partitionNum].AddRecord(record)
	return nil
}

// savePartition writes a partition of records to a new CSV file, including the header row.
func savePartition(filePath string, partitionNum int, partition *Partition, headers []string, done chan<- error) {
	var err error
	defer func() {
		done <- err
	}()

	partitionFileName := fmt.Sprintf("%s_partition_%d.csv", filePath, partitionNum+1)
	
	csvFile, err := os.OpenFile(partitionFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		err = fmt.Errorf("creating partition file: %w", err)
		return
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write headers
	if err = csvWriter.Write(headers); err != nil {
		err = fmt.Errorf("writing headers: %w", err)
		return
	}

	// Pre-allocate buffer based on partition size for better performance
	if partition.Len() > 1000 {
		csvWriter.Buffer(make([]byte, 0, partition.Len()*100)) // Estimate 100 bytes per record
	}

	// Write records in batches for better performance
	for _, record := range partition.Records {
		if err = csvWriter.Write(record); err != nil {
			err = fmt.Errorf("writing record: %w", err)
			return
		}
	}
}