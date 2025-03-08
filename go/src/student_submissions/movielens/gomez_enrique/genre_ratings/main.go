package main

import (
	"fmt"
	"movielens/split"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	dataframe "github.com/kfultz07/go-dataframe"
)

type Stats struct {
	Rating       float64
	Observations int
}

var lock = &sync.RWMutex{}

func main() {
	// Arg 1: file to split
	filename := os.Args[1]

	// Arg 2: total chunks
	totalChunks, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	// Arg3: out-path
	outPath := os.Args[3]

	// Arg4: file to inner join with
	filenameJoin := os.Args[4]

	// Split file
	start := time.Now()
	split.Split(filename, totalChunks, outPath)
	elapsed := time.Since(start)
	fmt.Printf("Split: %dms\n", elapsed.Milliseconds())

	// TODO: use generic filepaths
	dfMovies := dataframe.CreateDataFrame("", filenameJoin)

	// Create a map to keep track of stats
	stats := make(map[string]*Stats)

	// Process file
	matches, _ := filepath.Glob(filepath.Join(outPath, "tmp_ratings*.csv"))

	// Parallel procesing
	start = time.Now()
	var wg sync.WaitGroup
	for _, p := range matches {
		wg.Add(1)
		go func() {
			processFile(p, dfMovies, stats)
			wg.Done()
		}()
	}
	wg.Wait()
	elapsed = time.Since(start)
	fmt.Printf("Merge & Count: %dms\n\n", elapsed.Milliseconds())

	// Sequential processing
	//for _, p := range matches {
	//	processFile(p, dfMovies, stats)
	//}

	// Sort stats
	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Print stats
	for _, k := range keys {
		fmt.Println(k, stats[k].Rating/float64(stats[k].Observations))
	}
}

func processFile(filename string, dfMovies dataframe.DataFrame, stats map[string]*Stats) {
	dfRatings := dataframe.CreateDataFrame("./", filename)
	// TODO: check if I'm not loosing any row with values after replacing headers
	dfRatings.Headers = map[string]int{"userId": 0, "movieId": 1, "rating": 2, "timestamp": 4}
	err := dfRatings.Merge(&dfMovies, "movieId", "genres")
	if err != nil {
		panic(err)
	}

	for _, row := range dfRatings.KeepColumns([]string{"genres", "rating"}).FrameRecords {
		// Get rating data
		rating, err := strconv.ParseFloat(row.Data[1], 64)
		if err != nil {
			panic(err)
		}

		// Split into individual genres
		genres := strings.Split(row.Data[0], "|")

		// Count per genre
		lock.Lock()
		for _, key := range genres {
			_, keyExists := stats[key]
			if !keyExists {
				stats[key] = new(Stats)
			}
			stats[key].Rating += rating
			stats[key].Observations += 1
		}
		lock.Unlock()
	}
}
