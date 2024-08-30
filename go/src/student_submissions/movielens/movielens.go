package fileprocessing

import (
	"dataframe"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Mt_FindRatingsMaster() {
	fmt.Println("In MtFindRatingsMaster")
	start := time.Now()
	number_of_files := 10 // number of files with ratings is also number of threads for multi-threading

	// know_genres is a 1D array that contains the Known Genres
	know_genres := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	number_of_genres := len(kg) // number of know genres

	// "ratings" is a 2D array where the ratings values for each genre are maintained.
	// The columns signal/maintain the core number where a worker is running.
	// Tho rows in that column maintain the rating values for that core and that genre
	ratings := make([][]float64, number_of_genres)
	// "count_all" is a 2D array where the count of Ratings for each genre is maintained.
	// The columns signal the core number where the worker is running.
	// The rows in that column maintain the counts of that genre
	count_all := make([][]int, number_of_genres)
	// populate de "number_of_genres" rows of "ratings" and "count_all" with "number_of_files" columns
	for i := 0; i < number_of_genres; i++ {
		ratings[i] = make([]float64, number_of_files)
		count_all[i] = make([]int, number_of_files)
	}
	var ci = make(chan int)                                    // create the channel to sync all workers
	movies = SplitBigFile("movies.csv", number_of_files, "./") // THIS IS THE CODE YOU NEED TO DEVELOP in the "file_splitting.go" file
	// run FindRatings in 10 workers
	for i := 0; i < number_of_files; i++ {
		go Mt_FindRatingsWorker(i+1, ci, know_genres, &count_all, &ratings, movies)
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
	local_count := make([]int, number_of_genres)
	local_values := make([]float64, number_of_genres)
	for i := 0; i < number_of_genres; i++ {
		for j := 0; j < number_of_files; j++ {
			local_count[i] += count_all[i][j]
			local_values[i] += ratings[i][j]
		}
	}
	for i := 0; i < number_of_genres; i++ {
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", know_genres[i]), "  ", fmt.Sprintf("%8d", local_count[i]))
	}
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	println("Mt_FindRatingsMaster is Done")
}

func Mt_FindRatingsWorker(worker int, ci chan int, know_genres []string, count_all *[][]int, value *[][]float64, movies dataframe.DataFrame) {
	aFileName := "ratings_" + fmt.Sprintf("%02d", worker) + ".csv"
	println("Worker ", fmt.Sprintf("%02d", worker), " is processing file ", aFileName, "\n")

	ratings := ReadRatingsCsvFile(aFileName) // THIS IS ANOTHER PIECE OF THE PUZZLE YOU NEED TO DEVELOP
	number_of_genres := len(know_genres)
	start := time.Now()

	// Import all records from the movies DF (DataFrame) into the ratings DF, keeping genres column from movies
	// data_frame.Merge is the equivalent of an inner-join in the DF lib I am using here
	ratings.Merge(&movies, "movieId", "genres")

	// We only need "genres" and "ratings" to find Count(Ratings | Genres), so keep only those columns
	grcs := [2]string("genres", "rating") // grcs => Genres Ratings Columns
	grDF := ratings.KeepColumns(grcs[:])  // grDF => Genres Ratings D

	for ig := 0; ig < number_of_genres; ig++ {
		for _, row := range grDF.FrameRecords {
			if strings.Contains(row.Data[0], know_genres[ig]) {
				(*count_all)[ig][worker-1] += 1
				v, _ := strconv.ParseFloat((row.Data[1]), 32) // do not check for error
				(*value)[ig][worker-1] += v
			}
		}
	}

	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	fmt.Println("Worker ", worker, " completed")

	// notify to the master that this worker has completed it's job
	ci <- 1
}

func main() {
	Mt_FindRatingsMaster()
}
