package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/kfultz07/go-dataframe"
)

// readCSVtoDataFrame reads a CSV file and returns a DataFrame
// readCSVtoDataFrame reads a CSV file and returns a DataFrame
func readCSVtoDataFrame(filePath string, filename string) dataframe.DataFrame {
	// Read the CSV into a DataFrame
	df := dataframe.CreateDataFrame(filePath, filename)
	return df
}

// GenreData stores the count and sum of ratings for each genre
type GenreData struct {
	Count int
	Sum   float64
}

// GenreStats stores the statistics for all genres
type GenreStats map[string]GenreData

// List of valid genres
var validGenres = map[string]struct{}{
	"Action":             {},
	"Adventure":          {},
	"Animation":          {},
	"Children":           {},
	"Comedy":             {},
	"Crime":              {},
	"Documentary":        {},
	"Drama":              {},
	"Fantasy":            {},
	"Film-Noir":          {},
	"Horror":             {},
	"IMAX":               {},
	"Musical":            {},
	"Mystery":            {},
	"Romance":            {},
	"Sci-Fi":             {},
	"Thriller":           {},
	"War":                {},
	"Western":            {},
	"(no genres listed)": {},
}

// mergeCSVWithDataFrameConcurrently merges a list of CSV files concurrently with the given DataFrame
// It calculates the number of ratings and average rating after each merge and stores that in a list.
func mergeCSVWithDataFrameConcurrently(filePaths []string, existingDF *dataframe.DataFrame) (GenreStats, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	genreStats := make(GenreStats)
	ch := make(chan struct {
		err error
	}, len(filePaths))

	sem := make(chan struct{}, 20)

	// Process each file concurrently
	for _, filePath := range filePaths {
		wg.Add(1)

		sem <- struct{}{}

		go func(filePath string) {
			defer wg.Done()
			defer func() { <-sem }()
			fmt.Println("Processing: ", filePath)
			// Read the CSV into a new DataFrame
			newDF := readCSVtoDataFrame("output", filePath)
			mu.Lock()

			// Merge the new DataFrame with the existing one
			mergedDF, err := existingDF.InnerMerge(&newDF, "movieId")
			mu.Unlock()
			if err != nil {
				panic(err)
			}

			mu.Lock()
			defer mu.Unlock()
			for _, row := range mergedDF.FrameRecords {
				// row.Val() is used to extract the value in a specific column while iterating
				genre := row.Val("genres", mergedDF.Headers)
				if _, exists := validGenres[genre]; exists {
					rating, err := strconv.ParseFloat(row.Val("rating", mergedDF.Headers), 64)
					if err != nil {
						fmt.Println("Error converting string to float:", err)
						return
					}
					genreData := genreStats[genre]
					genreData.Count++
					genreData.Sum += rating
					genreStats[genre] = genreData
				}
			}
			fmt.Println("Finalizando cuentas: ", filePath)

			// Store the results in the channel
			ch <- struct {
				err error
			}{err: nil}
		}(filePath)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(ch)

	// Collect results from the channel
	for result := range ch {
		if result.err != nil {
			log.Printf("Error processing file: %v", result.err)
		}
	}

	return genreStats, nil
}

func calculateTotalAndAverage(genreStats GenreStats) {
	// Calculate and print the total sum and average for each genre
	var totalSum float64
	var totalCount int

	for genre, stats := range genreStats {
		fmt.Printf("Genre: %s, Count: %d, Sum: %.2f, Average: %.2f\n", genre, stats.Count, stats.Sum, stats.Sum/float64(stats.Count))
		totalSum += stats.Sum
		totalCount += stats.Count
	}

	// Calculate the overall total and average
	if totalCount > 0 {
		avg := totalSum / float64(totalCount)
		fmt.Printf("Total Sum: %.2f, Total Count: %d, Total Average: %.2f\n", totalSum, totalCount, avg)
	}
}

// Worker function to process a chunk of rows
func processChunk(rows [][]string, wg *sync.WaitGroup, numberOfChunks int, parentDir string, nameOfArchive string) {
	defer wg.Done()

	// Convert rows to CSV string and read as DataFrame
	chunkName := fmt.Sprintf("%s-%d.csv", nameOfArchive, numberOfChunks)
	chunkName = filepath.Join(parentDir, chunkName)

	// Open a CSV file for writing
	file, err := os.Create(chunkName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)

	// Write each row to the CSV file
	mySlice := []string{"userId", "movieId", "rating", "timestamp"}
	writer.Write(mySlice)
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			fmt.Println("Error writing row to CSV:", err)
			return
		}
	}

	// Flush any buffered data to the file
	writer.Flush()

	// Check for any errors during writing
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}

	fmt.Println("Data written to output.csv successfully.")
}

func generateBatchesFromCsvFile(nameOfCsvFile string, nameOfArchive string, numberOfRowsPerChunk int, batchFileNameConvetion string) {
	parentRoot := nameOfArchive
	err := os.MkdirAll(parentRoot, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	f, err := os.Open(nameOfCsvFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	var wg sync.WaitGroup
	numRowsPerChunk := numberOfRowsPerChunk // Adjust based on memory/size

	// Read and process CSV in chunks
	numberOfChunks := 0
	for {
		rows := make([][]string, 0, numRowsPerChunk)
		for i := 0; i < numRowsPerChunk; i++ {
			row, err := r.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				panic(err)
			}
			rows = append(rows, row)
		}
		if len(rows) == 0 {
			break
		}
		numberOfChunks += 1

		// Launch a goroutine to process each chunk
		wg.Add(1)
		go processChunk(rows, &wg, numberOfChunks, parentRoot, batchFileNameConvetion)
	}
	wg.Wait()
	fmt.Println("Termine de hacer la division")
}

func getNumberTotalNumberOfFilesInDir(batchesDirectory string) int {

	// Open the directory
	files, err := os.ReadDir(batchesDirectory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return -1
	}

	// Count the number of files
	fileCount := 0
	for _, file := range files {
		if !file.IsDir() { // Check if it's a file, not a directory
			fileCount++
		}
	}

	fmt.Printf("Total number of files in directory %s: %d\n", batchesDirectory, fileCount)
	return fileCount - 1
}

func main() {

	nameOfArchive := "batch"
	batchesDir := "output"
	ratingsFileName := "ratings.csv"

	generateBatchesFromCsvFile(ratingsFileName, batchesDir, 1000, nameOfArchive)

	numberOfFiles := getNumberTotalNumberOfFilesInDir(batchesDir)
	// Example of an existing DataFrame (replace with your actual DataFrame)
	existingDF := dataframe.CreateDataFrame("output", "movies.csv")

	// List of CSV files to process concurrently
	n := numberOfFiles // Set the desired number of strings
	files := make([]string, n)

	for i := 0; i <= n-1; i++ {
		chunkName := fmt.Sprintf("%s-%d.csv", nameOfArchive, i+1)
		files[i] = chunkName
	}

	// Merge the CSV files concurrently and get the row counts and average ratings
	genreStats, err := mergeCSVWithDataFrameConcurrently(files, &existingDF)
	if err != nil {
		log.Fatalf("Error merging CSV files: %v", err)
	}

	// Print the number of ratings and average ratings after each merge
	calculateTotalAndAverage(genreStats)
}
