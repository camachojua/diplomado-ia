package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-gota/gota/dataframe"
)

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
	mySlice := []string{"USERID", "movieId", "rating", "TIMESTAMP"}
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

func elGfazo(batchesDirectory string, numberOfWorkers int, moviesDirectory string, nameOfArchive string) {
	var wg sync.WaitGroup
	//Read the total number of files in the directory
	numberOfFiles := getNumberTotalNumberOfFilesInDir(batchesDirectory)

	//Based of that calculate with the #numberOfFiles/#numberOfWorkers to know how many files are going to be taken by each worker
	filenames := generateFileNames(nameOfArchive, numberOfFiles)
	var batchesPerWorker int = numberOfFiles / numberOfWorkers
	remainingFilesThatDontComplete := numberOfFiles % numberOfWorkers
	ratingsChanel := make(chan map[string][]float64)

	for i := 0; i < numberOfWorkers; i++ {
		start := i * batchesPerWorker   // Start index for the current worker
		end := start + batchesPerWorker // End index for the current worker

		// If this is the last worker, include any remaining files
		if i == numberOfWorkers-1 && remainingFilesThatDontComplete != 0 {
			end = start + remainingFilesThatDontComplete
		}
		filenamesDesignated := filenames[start:end]
		fmt.Println(filenamesDesignated)
		wg.Add(1)
		go losGodines(filenamesDesignated, ratingsChanel, &wg)

	}

	categoryCountRating := map[string][]float64{
		"Action":             {0, 0},
		"Adventure":          {0, 0},
		"Animation":          {0, 0},
		"Children":           {0, 0},
		"Comedy":             {0, 0},
		"Crime":              {0, 0},
		"Documentary":        {0, 0},
		"Drama":              {0, 0},
		"Fantasy":            {0, 0},
		"Film-Noir":          {0, 0},
		"Horror":             {0, 0},
		"IMAX":               {0, 0},
		"Musical":            {0, 0},
		"Mystery":            {0, 0},
		"Romance":            {0, 0},
		"Sci-Fi":             {0, 0},
		"Thriller":           {0, 0},
		"War":                {0, 0},
		"Western":            {0, 0},
		"(no genres listed)": {0, 0},
	}

	// Read from the channel and aggregate values (consumes each message as it's read)
	for data := range ratingsChanel {
		fmt.Println("fina", data)
		for category, values := range data {
			if _, exists := categoryCountRating[category]; exists {
				categoryCountRating[category][0] += values[0] // Aggregate the first value
				categoryCountRating[category][1] += values[1] // Aggregate the second value
			} else {
				// Handle any new categories that weren't initialized
				categoryCountRating[category] = values
			}
		}
	}

	fmt.Println("TERMINO MASTER")
	go func() {
		wg.Wait()
		close(ratingsChanel)
	}()

	// Print the aggregated results
	for category, values := range categoryCountRating {
		fmt.Printf("%s: Count=%.1f, Rating=%.1f\n", category, values[0], values[1])
	}

}

func generateFileNames(nameOfArchive string, numberOfFiles int) []string {

	var fileNames []string
	for i := 0; i < numberOfFiles; i++ {
		fileNames = append(fileNames, fmt.Sprintf("%s-%d.csv", nameOfArchive, i+1))
	}

	return fileNames
}

func losGodines(filesDesignated []string, ch chan<- map[string][]float64, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Worker started")
	category_count_rating := map[string][]float64{
		"Action":             {0, 0},
		"Adventure":          {0, 0},
		"Animation":          {0, 0},
		"Children":           {0, 0},
		"Comedy":             {0, 0},
		"Crime":              {0, 0},
		"Documentary":        {0, 0},
		"Drama":              {0, 0},
		"Fantasy":            {0, 0},
		"Film-Noir":          {0, 0},
		"Horror":             {0, 0},
		"IMAX":               {0, 0},
		"Musical":            {0, 0},
		"Mystery":            {0, 0},
		"Romance":            {0, 0},
		"Sci-Fi":             {0, 0},
		"Thriller":           {0, 0},
		"War":                {0, 0},
		"Western":            {0, 0},
		"(no genres listed)": {0, 0},
	}

	var partialReviewsDataframe dataframe.DataFrame = getDataframeFromCsv(filesDesignated, "Pse para aca 1")
	var moviesDataframe dataframe.DataFrame = getDataframeFromCsvMovies("Pse para aca 2")
	joinedDataframe := partialReviewsDataframe.LeftJoin(moviesDataframe, "movieId")
	fmt.Println("dddddd")
	fmt.Println(joinedDataframe)
	fmt.Println("xx")
	for i := 0; i < joinedDataframe.Nrow(); i++ {
		row := joinedDataframe.Subset(i) // Get the row as a DataFrame
		generes_value := row.Col("genres").Elem(0).String()
		ratings_value := row.Col("rating").Elem(0).Float()
		fmt.Println("ddddd")
		fmt.Println(generes_value)
		fmt.Println(ratings_value)
		fmt.Println("xxxx")
		if generes_value == "" {
			fmt.Println("no category for row", row)
			new_count_value := category_count_rating["(no genres listed)"][0] + 1
			new_review_value := category_count_rating["(no genres listed)"][1] + ratings_value
			category_count_rating["(no genres listed)"][0] = new_count_value
			category_count_rating["(no genres listed)"][1] = new_review_value
		} else {
			//We have to split string
			// Split the string by "|"
			categoryArray := strings.Split(generes_value, "|")
			for _, element := range categoryArray {
				new_count_value := category_count_rating[element][0] + 1
				new_review_value := category_count_rating[element][1] + ratings_value
				category_count_rating[element][0] = new_count_value
				category_count_rating[element][1] = new_review_value
			}
			// Print the result
			fmt.Println(categoryArray)
		}
	}
	fmt.Println("Termine worker")
	fmt.Println("category_count_rating")
	ch <- category_count_rating
}

func getDataframeFromCsv(filesDesignated []string, errsor string) dataframe.DataFrame {
	// Create an empty DataFrame to hold the concatenated result

	var combinedDF dataframe.DataFrame

	for _, path := range filesDesignated {
		// Open each CSV file
		complete_path := "batchArchive20/" + path
		file, err := os.Open(complete_path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Read the CSV file into a Gota DataFrame
		df := dataframe.ReadCSV(file)

		// Append each DataFrame to the combined DataFrame
		if combinedDF.Nrow() == 0 {
			combinedDF = df
		} else {
			combinedDF = combinedDF.Concat(df)
		}
	}
	fmt.Println(errsor)
	return combinedDF

}

func getDataframeFromCsvMovies(errsor string) dataframe.DataFrame {
	// Create an empty DataFrame to hold the concatenated result

	fmt.Println(errsor)

	// Open each CSV file
	file, err := os.Open("movies.csv")
	if err != nil {
		panic("Algo trono boss,  es movies")
	}
	defer file.Close()

	// Read the CSV file into a Gota DataFrame
	df := dataframe.ReadCSV(file)

	return df

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
	return fileCount
}
func main() {
	nameOfArchive := "batch"
	batchesDir := "batchArchive20"
	ratingsFileName := "ratings.csv"

	generateBatchesFromCsvFile(ratingsFileName, batchesDir, 10000, nameOfArchive)
	elGfazo(batchesDir, 100, "movies.csv", nameOfArchive)

}
