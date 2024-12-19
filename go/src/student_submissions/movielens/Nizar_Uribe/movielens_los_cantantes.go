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

// Mapa de géneros por movieId
var moviesMap = make(map[string][]string)
var mu sync.Mutex

// Estructura para almacenar el conteo y suma de calificaciones por género
type GenreStats struct {
	Count int
	Sum   float64
}

// Función para cargar `movies.csv` y construir el mapa de géneros
func loadMovies() error {
	file, err := os.Open("movies.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read() // Saltar encabezado

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		movieId := record[0]
		genres := strings.Split(record[2], "|")
		moviesMap[movieId] = genres
	}

	return nil
}

// Función para procesar cada archivo `ratings_part_*` y acumular los resultados en `genreStats`
func processFile(filename string, genreStats map[string]*GenreStats, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error al abrir %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read() // Saltar encabezado

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		movieId := record[1]
		rating, _ := strconv.ParseFloat(record[2], 64)

		mu.Lock()
		genres, exists := moviesMap[movieId]
		mu.Unlock()

		if !exists {
			continue
		}

		mu.Lock()
		for _, genre := range genres {
			if _, ok := genreStats[genre]; !ok {
				genreStats[genre] = &GenreStats{}
			}
			genreStats[genre].Count++
			genreStats[genre].Sum += rating
		}
		mu.Unlock()
	}
}

func main() {
	// Guardar el tiempo de inicio
	start := time.Now()

	// Cargar los géneros de `movies.csv`
	if err := loadMovies(); err != nil {
		fmt.Println("Error cargando movies.csv:", err)
		return
	}

	// Mapa concurrente para almacenar conteo y suma de calificaciones por género
	genreStats := make(map[string]*GenreStats)

	// Canal para sincronizar el procesamiento de archivos en paralelo
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		filename := fmt.Sprintf("ratings_part_%d.csv", i)
		wg.Add(1)
		go processFile(filename, genreStats, &wg)
	}

	wg.Wait()

	// Crear archivo CSV de resultados
	outputFile, err := os.Create("genre_ratings_summary_go1.csv")
	if err != nil {
		fmt.Println("Error creando genre_ratings_summary_go1.csv:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Escribir encabezado
	writer.Write([]string{"Genre", "Count", "AverageRating"})

	// Calcular promedio y escribir los resultados en el archivo
	for genre, stats := range genreStats {
		avgRating := stats.Sum / float64(stats.Count)
		writer.Write([]string{genre, strconv.Itoa(stats.Count), fmt.Sprintf("%.2f", avgRating)})
	}

	// Calcular y mostrar el tiempo de ejecución
	duration := time.Since(start)
	fmt.Printf("Archivo 'genre_ratings_summary_go1.csv' creado con éxito.\n")
	fmt.Printf("Tiempo de ejecución: %v\n", duration)
}
