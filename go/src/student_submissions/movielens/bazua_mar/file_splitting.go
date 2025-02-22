/// Movielens en go
/// Mar Bazúa

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kfultz07/go-dataframe"
)

type RatingObj struct {
	UserId    int64
	movieId   int64
	rating    float64
	Timestamp int64
}

// Procesa los registros y los guarda en archivos de salida
func processRecords(records [][]string, outputIndex int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Crear el archivo de salida
	outputFileName := fmt.Sprintf("./output/ratings_%02d.csv", outputIndex)
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("Error al crear el archivo %s: %v", outputFileName, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezados
	writer.Write([]string{"UserId", "movieId", "rating", "Timestamp"})

	// Procesar y escribir cada registro en el archivo
	for _, record := range records {
		userId, _ := strconv.ParseInt(record[0], 10, 64)
		movieId, _ := strconv.ParseInt(record[1], 10, 64)
		rating, _ := strconv.ParseFloat(record[2], 64)
		timestamp, _ := strconv.ParseInt(record[3], 10, 64)

		data := RatingObj{
			UserId:    userId,
			movieId:   movieId,
			rating:    rating,
			Timestamp: timestamp,
		}

		// Convertir el objeto Rating a un registro CSV
		recordToWrite := []string{
			strconv.FormatInt(data.UserId, 10),
			strconv.FormatInt(data.movieId, 10),
			strconv.FormatFloat(data.rating, 'f', 1, 64),
			strconv.FormatInt(data.Timestamp, 10),
		}
		writer.Write(recordToWrite)
	}
}

func main() {
	start := time.Now()

	executeSplit := true // Cambia a false si no deseas ejecutar el Split de Ratings

	if executeSplit {
		fmt.Println("Split Ratings will be executed:")
		splitStart := time.Now()

		// Leer el archivo CSV
		file, err := os.Open("./ml-25m/ratings.csv")
		if err != nil {
			log.Fatalf("Error al abrir el archivo CSV: %v", err)
		}
		defer file.Close()

		records, err := csv.NewReader(file).ReadAll()
		if err != nil {
			log.Fatalf("Error al leer el archivo CSV: %v", err)
		}

		// Remover encabezados si existen
		if len(records) > 0 && records[0][0] == "UserId" {
			records = records[1:]
		}

		// Determinar el tamaño de cada chunk para los N archivos de salida
		totalJobs := 10
		sizeRange := (len(records)) / totalJobs

		// Crear WaitGroup para sincronizar las goroutines
		var wg sync.WaitGroup

		// Procesar los registros en paralelo, dividiéndolos en chunks
		for i := 0; i < totalJobs; i++ {
			start := i * sizeRange
			end := start + sizeRange
			if end > len(records) {
				end = len(records)
			}

			wg.Add(1)
			go processRecords(records[start:end], i+1, &wg)
		}

		// Esperar a que todas las goroutines terminen
		wg.Wait()

		fmt.Println("Split Ratings execution duration:", time.Since(splitStart))
	}

	fmt.Println("Join Ratings with Movielens using go rutines will be executed:")
	Mt_FindRatingsMaster()

	fmt.Println("Total execution duration:", time.Since(start))
}

// Segundo Código (modificado para trabajar en conjunto con el primero)
func Mt_FindRatingsMaster() {
	fmt.Println("In MtFindRatingsMaster")
	start := time.Now()
	nf := 10 // number of files with ratings is also number of threads for multi-threading

	kg := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	ng := len(kg)

	ra := make([][]float64, ng)
	ca := make([][]int, ng)

	for i := 0; i < ng; i++ {
		ra[i] = make([]float64, nf)
		ca[i] = make([]int, nf)
	}
	var ci = make(chan int)
	movies := ReadMoviesCsvFile("./ml-25m/", "movies.csv")

	for i := 0; i < nf; i++ {
		go Mt_FindRatingsWorker(i+1, ci, kg, &ca, &ra, movies)
	}

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

	locCount := make([]int, ng)
	locVals := make([]float64, ng)
	for i := 0; i < ng; i++ {
		for j := 0; j < nf; j++ {
			locCount[i] += ca[i][j]
			locVals[i] += ra[i][j]
		}
	}

	// Calcular y mostrar el promedio de ratings por género
	for i := 0; i < ng; i++ {
		avgRating := 0.0
		if locCount[i] > 0 {
			avgRating = locVals[i] / float64(locCount[i])
		}
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", kg[i]), "  ", fmt.Sprintf("%8d", locCount[i]), "  ", fmt.Sprintf("Avg Rating: %.2f", avgRating))
	}

	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	println("Mt_FindRatingsMaster is Done")
}

func Mt_FindRatingsWorker(w int, ci chan int, kg []string, ca *[][]int, va *[][]float64, movies dataframe.DataFrame) {
	aFileName := "./output/ratings_" + fmt.Sprintf("%02d", w) + ".csv"
	println("Worker  ", fmt.Sprintf("%02d", w), "  is processing file ", aFileName, "\n")

	bObName := fmt.Sprintf("ratings_%02d.csv", w)
	ratings := ReadRatingsCsvFile("./output/", bObName)
	ng := len(kg)
	start := time.Now()

	ratings.Merge(&movies, "movieId", "genres")

	grcs := [2]string{"genres", "rating"}
	grDF := ratings.KeepColumns(grcs[:])
	for ig := 0; ig < ng; ig++ {
		for _, row := range grDF.FrameRecords {
			if strings.Contains(row.Data[0], kg[ig]) {
				(*ca)[ig][w-1] += 1
				v, _ := strconv.ParseFloat((row.Data[1]), 32)
				(*va)[ig][w-1] += v
			}
		}
	}
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	fmt.Println("Worker ", w, " completed")

	ci <- 1
}

// Función para leer el archivo de películas
func ReadMoviesCsvFile(filePath string, fileOb string) dataframe.DataFrame {
	// Abrir el archivo
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error al abrir el archivo %s: %v", filePath, err)
	}
	defer file.Close()

	// Leer el CSV y devolver el DataFrame
	return dataframe.CreateDataFrame(filePath, fileOb)
}

// Función para leer el archivo de ratings
func ReadRatingsCsvFile(filePath string, fileOb string) dataframe.DataFrame {
	// Abrir el archivo
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error al abrir el archivo %s: %v", filePath, err)
	}
	defer file.Close()

	// Leer el CSV y devolver el DataFrame
	return dataframe.CreateDataFrame(filePath, fileOb)
}
