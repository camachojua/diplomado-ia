package main

import (
    "fmt"
    "time"
	"strings"
	"strconv"
	"github.com/kfultz07/go-dataframe"
	"os"
	"log"
)

func Mt_FindRatingsMaster() {
	fmt.Println("In MtFindRatingsMaster")
	start := time.Now()
	nf := 2 // number of files with ratings is also number of threads for multi-threading
	
	// generos is a 1D array that contains the Known Genres
	generos := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	ng := len(generos) // number of known genres
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
	var ci = make(chan int)		// create the channel to sync all workers
	
	movies := ReadMoviesCsvFile("./ADCC_project/csv_files/")
	// run FindRatings in 10 workers
	for i := 0; i < nf; i++ {
		go Mt_FindRatingsWorker(i+1, ci, generos, &ca, &ra, movies)
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
	locCount := make([]int, ng)
	locVals := make([]float64, ng)
	for i := 0; i < ng; i++ {
		for j := 0; j < nf; j++ {
			locCount[i] += ca[i][j]
			locVals[i] += ra[i][j]
		}
	}
	for i := 0; i < ng; i++ {
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", generos[i]), "  ", fmt.Sprintf("%8d", locCount[i]))
	}
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	println("Mt_FindRatingsMaster is Done")
}

func Mt_FindRatingsWorker(w int, ci chan int, generos []string, ca *[][]int, va *[][]float64, movies dataframe.DataFrame) {
	//aFileName := "./ADCC_project/splited_files/ratings_" + fmt.Sprintf("%02d", w) + ".csv"
	aFileName := "ratings_" + fmt.Sprintf("%d", w-1) + ".csv"
	path:="./ADCC_project/splited_files/"

	println("Worker  ", fmt.Sprintf("%02d", w), "  is processing file ", aFileName, "\n")

	ratings := ReadRatingsCsvFile(path,aFileName)
	ng := len(generos)
	start := time.Now()

	// import all records from the movies DF into the ratings DF, keeping genres column from movies
       //df.Merge is the equivalent of an inner-join in the DF lib I am using here

	//fmt.Printf("Columnas de movies: %v (para worker %d)\n", movies.Columns(), w)
	//fmt.Printf("Columnas de ratings: %v (para worker %d)\n", ratings.Columns(), w)

	//fmt.Printf("Worker %d: Realizando merge\n", w)
	ratings.Merge(&movies, "movieId", "genres")
	//fmt.Printf("Worker %d: Merge realizado\n", w)

	// We only need "genres" and "ratings" to find Count(Ratings | Genres), so keep only those columns
	grcs := [2]string{"genres", "rating"} // grcs => Genres Ratings Columns
	grDF := ratings.KeepColumns(grcs[:])  // grDF => Genres Ratings DF
	for ig := 0; ig < ng; ig++ {
		for _, row := range grDF.FrameRecords {
			if strings.Contains(row.Data[0], generos[ig]) {
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

func ReadMoviesCsvFile(filePath string) dataframe.DataFrame {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer f.Close()

	df := dataframe.CreateDataFrame(filePath, "movies.csv")
	return df
}

func ReadRatingsCsvFile(filePath string, fileName string) dataframe.DataFrame {

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer f.Close()

	df := dataframe.CreateDataFrame(filePath, fileName)
	return df
}

func main() {
	Mt_FindRatingsMaster()
}