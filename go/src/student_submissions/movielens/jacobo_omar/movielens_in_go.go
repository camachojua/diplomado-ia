package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/kfultz07/go-dataframe"
)

// Read and split the big file.
// ###########################################################

// Open and read a csv file and returns the content in an array.
func ReadCsvToArray(fileName string, directory string) [][]string {
	file, err := os.Open(directory + fileName + ".csv")

	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	data, err := csvReader.ReadAll()

	if err != nil {
		log.Fatalf("Error extracting data from file %v: %s", fileName, err)
	}
	return data
}

// Create a csv file with the name and data provided in the path defined.
func WriteCsv1(data [][]string, name string, path string) {
	csvFile, err := os.Create(path + name + ".csv")
	if err != nil {
		log.Fatalf("Error creating new csv file %v: %s", name, err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	err = writer.WriteAll(data)
	if err != nil {
		log.Fatalf("Error writing new csv file %v: %s", name, err)
	}

	fmt.Printf("File %s has been created with %v rows\n", name, len(data))
}

// Function to split a csv file into small files.
// You provided a name for the file, the number of chuncks wich it will be divided and the directory
// where the file is located and the new ones will be saved.
func SplitFile(file_name string, number_of_chunks int, directory string) []string {
	t1 := time.Now()
	data := ReadCsvToArray(file_name, directory)

	//Extrae el encabezado para cada CSV
	header := data[0]
	//Quita el encabezado antes de dividir
	data = data[1:]

	fmt.Printf("%v rows in file %s\n", len(data), file_name)
	rowsPerFile := len(data) / number_of_chunks
	var filesCreated []string

	for i := 0; i < number_of_chunks; i++ {
		tempName := file_name + "_" + fmt.Sprintf("%02d", i+1)
		fmt.Printf("%s\n", tempName)
		path := directory
		fmt.Printf("%s\n", path)
		tempData := append([][]string{header}, data[i*rowsPerFile:(i+1)*rowsPerFile]...)
		WriteCsv1(tempData, tempName, path)
		filesCreated = append(filesCreated, tempName)
	}
	tf := time.Since(t1)
	fmt.Println("Executed in:", tf, "seconds")
	return filesCreated
}


// Define the number of workers.
// ###########################################################

// Get the number of CPUs.
func number_of_workers() int {
	//return 10
	return runtime.GOMAXPROCS(0)
}

// Definir la funcion que ejecutara el worker

// Open and read a csv file and returns a dataframe.
func ReadCsvToDataframe(filePath string) dataframe.DataFrame {
	path := "/mnt/c/Users/omarm/Downloads/ml-25m/"
	df := dataframe.CreateDataFrame(path, filePath)
	return df
}

// Find ratings of ratings files using multithreating, w: The number of the worker that will
// execute the function ci: Channel, knowGenres: List with the genres names, ca: Counter array,
// va: Ratings values array, movies: Dataframe with the tittles and genres of the movies.
func Mt_FindRatingsWorker(w int, ci chan int, knowGenres []string, ca *[][]int, va *[][]float64, movies dataframe.DataFrame) {
	aFileName := "ratings_" + fmt.Sprintf("%02d", w) + ".csv"
	println("Worker  ", fmt.Sprintf("%02d", w), "  is processing file ", aFileName, "\n")

	start := time.Now()
	ratings := ReadCsvToDataframe(aFileName)
	nGenres := len(knowGenres)

	// import all records from the movies DF into the ratings DF, keeping genres column from movies
	//df.Merge is the equivalent of an inner-join in the DF lib I am using here
	ratings.Merge(&movies, "movieId", "genres")

	// We only need "genres" and "ratings" to find Count(Ratings | Genres), so keep only those columns
	grcs := [2]string{"genres", "rating"} // grcs => Genres Ratings Columns
	grDF := ratings.KeepColumns(grcs[:])  // grDF => Genres Ratings DF
	for ig := 0; ig < nGenres; ig++ {
		for _, row := range grDF.FrameRecords {
			if strings.Contains(row.Data[0], knowGenres[ig]) {
				(*ca)[ig][w-1] += 1
				v, _ := strconv.ParseFloat((row.Data[1]), 32) // do not check for error
				(*va)[ig][w-1] += v
			}
		}
	}
	duration := time.Since(start)
	fmt.Println("Worker ", w, " completed")
	fmt.Println("Duration = ", duration)

	// notify master that this worker has completed its job
	ci <- 1
}

// Master function that execute the find ratings problem.
func Mt_FindRatingsMaster() {
	fmt.Println("In MtFindRatingsMaster\n")
	start := time.Now()
	nf := number_of_workers()

	//SplitFile("ratings", nf, "C:\\Users\\configurar\\Documents\\2024\\CursoGo\\peliculasEQUIPO\\")

	// knowGenres is a 1D array that contains the Known Genres
	knowGenres := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	nGenres := len(knowGenres) // number of known genres
	// ra is a 2D array where the ratings values for each genre are maintained.
	// The columns signal/maintain the core number where a worker is running.
	// The rows in that column maintain the rating values for that core and that genre
	ra := make([][]float64, nGenres)
	// ca is a 2D array where the count of Ratings for each genre is maintained
	// The columns signal the core number where the worker is running
	// The rows in that column maintain the counts for that that genre
	ca := make([][]int, nGenres)
	// populate the nGenres rows of ra and ca with nf columns
	for i := 0; i < nGenres; i++ {
		ra[i] = make([]float64, nf)
		ca[i] = make([]int, nf)
	}
	var ci = make(chan int) // create the channel to sync all workers
	movies := ReadCsvToDataframe("movies.csv")
	println("Reading of movies.csv completed\n")
	// run FindRatings in 10 workers
	for i := 0; i < nf; i++ {
		go Mt_FindRatingsWorker(i+1, ci, knowGenres, &ca, &ra, movies)
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
		if iMsg == nf {
			break
		}
	}
	// all workers completed their work. Collect results and produce report
	locCount := make([]int, nGenres)
	locVals := make([]float64, nGenres)
	locPromedio := make([]float64, nGenres)
	for i := 0; i < nGenres; i++ {
		for j := 0; j < nf; j++ {
			locCount[i] += ca[i][j]
			locVals[i] += ra[i][j]
		}
		locPromedio[i] = locVals[i] / float64(locCount[i])
	}
	for i := 0; i < nGenres; i++ {
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", knowGenres[i]), "  ", fmt.Sprintf("%8d", locCount[i]), " ", fmt.Sprintf("%.2f", locPromedio[i]))
	}
	duration := time.Since(start)
	fmt.Println("\nMt_FindRatingsMaster duration = ", duration)
	println("Mt_FindRatingsMaster is Done")
}

func main() {
	println("Starting the movielens problem...\n")
	t0:=time.Now()
	SplitFile("ratings", number_of_workers(), "/mnt/c/Users/omarm/Downloads/ml-25m/")
	println("Files splited, starting the Mt_FindRatingsMaster...\n")
	Mt_FindRatingsMaster()
	duration := time.Since(t0)
	fmt.Println("Total duration for split files and find ratings = ", duration)

}