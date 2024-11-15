package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func split_csv() {
	// Abrir el archivo CSV
	largeCSVFile, err := os.Open("ratings.csv")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer largeCSVFile.Close()

	reader := csv.NewReader(largeCSVFile)

	// Crear canales y espera para sincronización
	const numFiles = 10
	recordChans := make([]chan []string, numFiles)
	var wg sync.WaitGroup

	// Inicializar gorutinas
	for i := 0; i < numFiles; i++ {
		recordChans[i] = make(chan []string, 1000)
		wg.Add(1)
		go func(i int, recordChan chan []string) {
			defer wg.Done()
			// Crear el archivo CSV
			fileName := fmt.Sprintf("ratings_part_%d.csv", i+1)
			smallCSVFile, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error al crear el archivo:", err)
				return
			}
			defer smallCSVFile.Close()

			writer := csv.NewWriter(smallCSVFile)

			// Escribir registros recibidos en el canal
			for record := range recordChan {
				if err := writer.Write(record); err != nil {
					fmt.Println("Error al escribir en el archivo:", err)
					return
				}
			}

			// Asegurarse de que todos los datos se escriban en el archivo
			writer.Flush()
			if err := writer.Error(); err != nil {
				fmt.Println("Error al finalizar la escritura:", err)
				return
			}

			fmt.Printf("Archivo %s creado con éxito.\n", fileName)
		}(i, recordChans[i])
	}

	// Leer el archivo CSV línea por línea y distribuir los registros
	recordIndex := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error al leer el archivo CSV:", err)
			break
		}
		// Enviar el registro al canal correspondiente
		recordChans[recordIndex%numFiles] <- record
		recordIndex++
	}

	// Cerrar todos los canales
	for i := 0; i < numFiles; i++ {
		close(recordChans[i])
	}

	// Esperar a que todas las gorutinas terminen
	wg.Wait()
}

// Función para leer los archivos CSV divididos
func ReadRatingsCsvFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error al abrir el archivo %s: %v", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error al leer el archivo %s: %v", filename, err)
	}

	// Devuelve las filas, excepto la cabecera
	return records[1:], nil
}

// Función para procesar archivo de ratings
func Mt_FindRatingsWorker(w int, ci chan int, kg []string, ca *[][]int, va *[][]float64, movies map[string]string) {
	aFileName := fmt.Sprintf("ratings_part_%d.csv", w)
	println("Trabajador ", w, " está procesando el archivo ", aFileName)

	// Leer el archivo de ratings
	ratings, err := ReadRatingsCsvFile(aFileName)
	if err != nil {
		fmt.Println("Error:", err)
		ci <- 1
		return
	}

	ng := len(kg) // Número de géneros a procesar
	start := time.Now()

	// Procesar cada fila del archivo ratings
	for _, row := range ratings {
		movieID := row[1]
		ratingStr := row[2]
		genres, exists := movies[movieID]
		if !exists {
			continue
		}

		// Convertir la calificación a float64
		rating, _ := strconv.ParseFloat(ratingStr, 64)

		// Contar calificaciones por cada género correspondiente
		for ig := 0; ig < ng; ig++ {
			if strings.Contains(genres, kg[ig]) {
				(*ca)[ig][w-1] += 1
				(*va)[ig][w-1] += rating
			}
		}
	}

	duration := time.Since(start)
	fmt.Println("Duración = ", duration)
	fmt.Println("Trabajador ", w, " completó su tarea")

	// Notificar al master que este trabajador ha completado su trabajo
	ci <- 1
}

// Función para cargar la información de movies en un mapa de movieId a géneros
func loadMovies(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error al abrir el archivo %s: %v", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error al leer el archivo %s: %v", filename, err)
	}

	movies := make(map[string]string)
	for _, row := range records[1:] {
		movieID := row[0]
		genres := row[2]
		movies[movieID] = genres
	}

	return movies, nil
}

// Función maestro que coordina el procesamiento y genera la respuesta final
func Mt_FindRatingsMaster() {
	fmt.Println("Iniciando Mt_FindRatingsMaster")
	start := time.Now()
	nf := 10 // Número de archivos con ratings

	// Lista de géneros conocidos
	kg := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	ng := len(kg) // Número de géneros conocidos

	// Inicializar las matrices para los resultados parciales
	ra := make([][]float64, ng)
	ca := make([][]int, ng)

	for i := 0; i < ng; i++ {
		ra[i] = make([]float64, nf)
		ca[i] = make([]int, nf)
	}

	// Crear el canal para sincronizar
	ci := make(chan int)

	// Cargar la información de películas desde movies
	movies, err := loadMovies("movies.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Lanzar
	for i := 0; i < nf; i++ {
		go Mt_FindRatingsWorker(i+1, ci, kg, &ca, &ra, movies)
	}

	// Esperar a que todos completen su tarea
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

	locCount := make([]int, ng)
	locVals := make([]float64, ng)
	locAvg := make([]float64, ng)

	for i := 0; i < ng; i++ {
		for j := 0; j < nf; j++ {
			locCount[i] += ca[i][j]
			locVals[i] += ra[i][j]
		}

		// Calcular el promedio
		if locCount[i] > 0 {
			locAvg[i] = locVals[i] / float64(locCount[i])
		} else {
			locAvg[i] = 0.0
		}
	}

	// Imprimir resultados finales
	fmt.Println("\nResultados finales:")
	for i := 0; i < ng; i++ {
		fmt.Printf("%2d  %-20s  %8d  %12.2f\n", i, kg[i], locCount[i], locAvg[i])
	}

	duration := time.Since(start)
	fmt.Printf("Duración total = %.9fs\n", duration.Seconds())
	fmt.Println("Mt_FindRatingsMaster ha finalizado.")
}

func main() {
	split_csv()
	Mt_FindRatingsMaster()
}
