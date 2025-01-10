package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	dataframe "github.com/kfultz07/go-dataframe"
)

func countRows(csvFile string) int {
	// Cuenta el numero de renglones que contiene el archivo CSV
	file, err := os.Open(csvFile)

	if err != nil {
		log.Fatal("Error while reading the CSV file", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	/*
		encabezados, err := reader.Read()
		if err != nil {
			log.Fatal("Error while reading the CSV file", err)
		}

		fmt.Println("Encabezado del archivo CSV contiene:", encabezados)*/

	rowCount := 0
	for {
		_, err := reader.Read()
		if err != nil {
			break
		}
		rowCount++
	}
	return rowCount
}

func processPart(csvFile string, start int, end int, part int, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal("Error while reading the CSV file", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Leer encabezados
	headers, err := reader.Read()
	if err != nil {
		log.Fatal("Error while reading CSV headers", err)
	}

	// Saltar filas hasta `start`
	for i := 0; i < start; i++ {
		reader.Read()
	}

	// Crear archivo de salida para esta parte
	aFileName := fmt.Sprintf("ratings_%02d.csv", part)
	outFile, err := os.Create(aFileName)
	if err != nil {
		log.Fatal("Error while creating output CSV file", err)
	}
	defer outFile.Close()
	writer := csv.NewWriter(outFile)

	// Escribir encabezado
	writer.Write(headers)

	for i := start; i < end; i++ {
		record, err := reader.Read()
		if err != nil {
			break
		}
		writer.Write(record)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatal("Error while writing to CSV file", err)
	}
	fmt.Printf("Finished part %02d\n", part)
}

func csvToDataFrame(dirPath string, archivo string) dataframe.DataFrame {
	// Llamada a CreateDataFrame que devuelve el dataframe solicitado
	df := dataframe.CreateDataFrame(dirPath, archivo)
	//fmt.Println("Columnas del DataFrame:", df.Headers)
	return df
}

func Mt_FindRatingsMaster(nf int) {
	fmt.Println("In MtFindRatingsMaster")

	//nf := 12 // number of files with ratings is also number of threads for multi-threading

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
	//var ci = make(chan int) // create the channel to sync all workers

	movies := csvToDataFrame("C:/Users/Gustavo/Documents/ml-25m/", "movies.csv")

	var wg sync.WaitGroup

	// run FindRatings in 10 workers
	for i := 0; i < nf; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			Mt_FindRatingsWorker(i+1, kg, &ca, &ra, movies)
		}(i)
	}

	wg.Wait()
	// wait for the workers
	/*
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
		}*/
	// all workers completed their work. Collect results and produce report
	locCount := make([]int, ng)
	locVals := make([]float64, ng)
	for i := 0; i < ng; i++ {
		for j := 0; j < nf; j++ {
			locCount[i] += ca[i][j]
			locVals[i] += ra[i][j]
		}
	}
	// Calcular los promedios
	for i := 0; i < ng; i++ {
		avg := 0.0
		if locCount[i] != 0 {
			avg = locVals[i] / float64(locCount[i])
		}
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", kg[i]), "  ", fmt.Sprintf("%8d", locCount[i]), "  ", fmt.Sprintf("promedio = %.5f", avg))
	}

	println("Mt_FindRatingsMaster is Done")
}

func Mt_FindRatingsWorker(w int, kg []string, ca *[][]int, va *[][]float64, movies dataframe.DataFrame) {
	aFileName := "ratings_" + fmt.Sprintf("%02d", w) + ".csv"
	println("Worker  ", fmt.Sprintf("%02d", w), "  is processing file ", aFileName)
	ratings := csvToDataFrame("C:/Users/Gustavo/go/src/movielens/", aFileName)
	ng := len(kg)

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

	fmt.Println("Worker ", w, " completed")

	// notify master that this worker has completed its job
	//ci <- 1
}

func main() {
	start_time := time.Now()

	ratings := "C:/Users/Gustavo/Documents/ml-25m/ratings.csv"
	//movies := "C:/Users/Gustavo/Documents/ml-25m/movies.csv"
	rowsRatings := countRows(ratings) // Calculamos el numero de renglones que contiene el archivo ratings
	//rowsMovies := countRows(movies)
	numCPU := runtime.NumCPU()         // Calculamos el nivel de concurrencia (El numero de archivos que se partira el original)
	rowsPerCPU := rowsRatings / numCPU // Calculamos el número de lineas que contiene cada partición

	fmt.Printf("N° Renglones de Ratings: %v\n", rowsRatings)
	fmt.Printf("El archivo ratings se dividira en: %v partes y seran %v renglones por partición\n", numCPU, rowsPerCPU)

	var wg sync.WaitGroup
	// Asegurarse de que todas las gorutinas terminen antes de que la ejecución continue
	for i := 0; i < numCPU; i++ {
		inicio := i * rowsPerCPU
		fin := inicio + rowsPerCPU
		if i == numCPU-1 {
			fin = rowsRatings
		}
		wg.Add(1) // Cada vez que inicia una nueva gorutina, aumenta el contador de gorutinas activas

		go processPart(ratings, inicio, fin, i+1, &wg)
		//Iniciar la gorutina que ejecutara processPart sin bloquear el programa principal
	}
	wg.Wait() /// El programa principal espera hasta que termina la gortuina
	//Termina division de archivos

	Mt_FindRatingsMaster(numCPU)

	end_time := time.Now()
	fmt.Println("Tiempo transcurrido:", end_time.Sub(start_time).Seconds())
}
