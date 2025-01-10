package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/kfultz07/go-dataframe"
)

func fileSharding() {

	// Read the file
	scanner := bufio.NewScanner(os.Stdin)
	inputFile := "ratings.csv"
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Print("File opened successfully\n")

	// Count the number of lines
	scanner = bufio.NewScanner(file)
	numberLines := 0
	for scanner.Scan() {
		numberLines++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Calculate the number of lines per file
	fmt.Printf("Number of lines in the file: %d\n", numberLines)
	numberFiles := 10
	linesPerFile := numberLines / numberFiles
	//remainder := numberLines % numberFiles

	// Extract the base name of the original file
	baseName := inputFile[:len(inputFile)-4]

	//remove the folder sharded if it exists even if it is not empty
	if _, err := os.Stat("sharded_" + baseName); !os.IsNotExist(err) {
		os.RemoveAll("sharded_" + baseName)
	}

	// Create the files into a folder called (sharded + baseName) if does not exist create it making it os agnostic
	//os.Stat returns an error if the file does not exist

	if _, err := os.Stat("sharded_" + baseName); os.IsNotExist(err) {
		os.Mkdir("sharded_"+baseName, 0755)
	}

	files := make([]*os.File, numberFiles)
	for i := 0; i < numberFiles; i++ {
		// Create the file with the name baseName + i

		fileName := fmt.Sprintf("%s_%02d.csv", baseName, i+1)
		file, err := os.Create("sharded_" + baseName + "/" + fileName)

		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		files[i] = file
	}

	// Write the lines to the files
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	fileIndex := 0
	lineCount := 0
	var header string
	for scanner.Scan() {
		//Maintain the header for each file
		if lineCount == 0 && fileIndex == 0 {
			header = scanner.Text() + "\n"
		}
		if lineCount == 0 {
			files[fileIndex].WriteString(header)
		}
		line := scanner.Text()
		files[fileIndex].WriteString(line + "\n")
		lineCount++
		if lineCount == linesPerFile {
			fileIndex++
			lineCount = 0
			// Fileindex is the index of the file to write to
			if fileIndex >= len(files) {
				// If the fileIndex is greater than the number of files, write the remainder to the last file
				fileIndex = len(files) - 1

				fmt.Print("Done sharding the file\n")
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the files
	for _, file := range files {
		file.Close()
	}

}

func ReadMoviesCsvFile(fileName string) dataframe.DataFrame {
	currentDir := path.Dir(fileName)
	movies := dataframe.CreateDataFrame(currentDir, fileName)
	return movies
}
func ReadRatingsCsvFile(aFileName string) dataframe.DataFrame {
	currentDir := path.Dir(aFileName)
	ratings := dataframe.CreateDataFrame(currentDir+"/sharded_ratings", aFileName)
	return ratings
}

func Mt_FindRatingsMaster() {
	fmt.Println("In MtFindRatingsMaster")
	start := time.Now()
	nf := 10 // number of files with ratings is also number of threads for multi-threading

	// kg is a 1D array that contains the Known Genres
	kg := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	ng := len(kg) // number of known genres
	// ra is a 2D array where the ratings values for each genre are maintained.
	// The columns signal/maintain the core number where a worker is running.
	// The rows in that column maintain the rating values for that core and that genre
	ra := make([][]float64, ng)
	// ca is a 2D array where the count of Ratings for each genre is maintained
	// The columns signal the core number where the worker is running
	// The rows in that column maintain the counts for that that genre
	ca := make([][]int, ng)
	// populate the ng rows of ra and ca with nf columns
	for i := 0; i < ng; i++ {
		ra[i] = make([]float64, nf)
		ca[i] = make([]int, nf)
	}
	var ci = make(chan int) // create the channel to sync all workers
	movies := ReadMoviesCsvFile("movies.csv")
	// run FindRatings in 10 workers
	for i := 0; i < nf; i++ {
		go Mt_FindRatingsWorker(i+1, ci, kg, &ca, &ra, movies)
	}
	// wait for the workers
	iMsg := 0
	go func() {
		for {
			i := <-ci
			iMsg += i
		}
	}()
	for {
		if iMsg == 10 {
			break
		}
	}
	// all workers completed their work. Collect results and produce report
	locCount := make([]int, ng)
	locVals := make([]float64, ng)
	for i := 0; i < ng; i++ {
		for j := 0; j < nf; j++ {
			locCount[i] += ca[i][j]
			locVals[i] += ra[i][j]
		}
	}
	for i := 0; i < ng; i++ {
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", kg[i]), "  ", fmt.Sprintf("%8d", locCount[i]))
	}
	//Print the average ratings for each genre
	//print space line
	fmt.Println()
	fmt.Println("Average Ratings for each genre")
	for i := 0; i < ng; i++ {
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", kg[i]), "  ", fmt.Sprintf("%8.2f", locVals[i]/float64(locCount[i])))
	}

	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	println("Mt_FindRatingsMaster is Done")
}

func Mt_FindRatingsWorker(w int, ci chan int, kg []string, ca *[][]int, va *[][]float64, movies dataframe.DataFrame) {
	aFileName := "ratings_" + fmt.Sprintf("%02d", w) + ".csv"
	println("Worker  ", fmt.Sprintf("%02d", w), "  is processing file ", aFileName, "\n")

	ratings := ReadRatingsCsvFile(aFileName)
	ng := len(kg)
	start := time.Now()

	// import all records from the movies DF into the ratings DF, keeping genres column from movies
	//df.Merge is the equivalent of an inner-join in the DF lib I am using here
	ratings.Merge(&movies, "movieId", "genres")

	// We only need "genres" and "ratings" to find Count(Ratings | Genres), so keep only those columns
	grcs := [2]string{"genres", "rating"} // grcs => Genres Ratings Columns
	grDF := ratings.KeepColumns(grcs[:])  // grDF => Genres Ratings DF
	for ig := 0; ig < ng; ig++ {
		for _, row := range grDF.FrameRecords {
			if strings.Contains(row.Data[0], kg[ig]) {
				(*ca)[ig][w-1] += 1
				v, _ := strconv.ParseFloat((row.Data[1]), 32) // do not check for error
				(*va)[ig][w-1] += v
			}
		}
	}
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	fmt.Println("Worker ", w, " completed")

	// notify master that this worker has completed its job
	ci <- 1
}

func main() {
	fileSharding()
	Mt_FindRatingsMaster()
}
