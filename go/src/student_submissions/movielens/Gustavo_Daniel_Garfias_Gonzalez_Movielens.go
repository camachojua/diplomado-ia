package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Movie struct {
	MovieID int
	Genres  string
}

func dividir(nombre string, N int) {
	file, err := os.Open(nombre)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	ratings, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error al leer el archivo CSV: %v", err)
	}
	x := len(ratings) - 1
	partes := x / N
	for i := 0; i < N; i++ {
		inicio := (i * partes) + 1
		fin := (i + 1) * partes
		if i == N-1 {
			resto := x - fin
			fin += resto
		}
		newFilename := fmt.Sprintf("rating%d.csv", i+1)
		newfile, err := os.Create(newFilename)
		if err != nil {
			log.Fatalf("Error al crear archivo %s: %v", newFilename, err)
		}

		writer := csv.NewWriter(newfile)
		writer.Write(ratings[0]) // Escribir encabezado

		for j := inicio; j <= fin; j++ {
			err := writer.Write(ratings[j])
			if err != nil {
				log.Fatalf("Error al escribir en archivo %s: %v", newFilename, err)
			}
		}
		writer.Flush()
		newfile.Close()
		fmt.Printf("Archivo Creado: %s\n", newFilename)
	}
}

func leerCSV(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error al abrir archivo %s: %v", fileName, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error al leer archivo CSV %s: %v", fileName, err)
	}
	return records
}

func resultados(n int, movies map[int]string) []string {
	ratingParts := leerCSV(fmt.Sprintf("rating%d.csv", n))
	return contar(ratingParts, movies)
}
func contar(ratingParts [][]string, movies map[int]string) []string {
	var generos []string
	for idx, rating := range ratingParts {

		if idx == 0 {
			continue
		}
		// Verificar que haya al menos 2 columnas

		if len(rating) <= 1 {
			log.Printf("Error: Rating tiene menos columnas de las esperadas")
			continue
		}

		// Convertir MovieID a int
		ratingMovieID, err := strconv.Atoi(rating[1])
		if err != nil {
			log.Printf("Error al convertir MovieID '%s': %v", rating[1], err)
			continue
		}

		// Buscar el género correspondiente
		if genres, found := movies[ratingMovieID]; found {
			generos = append(generos, strings.Split(genres, "|")...)
		} else {
			log.Printf("No se encontró el MovieID %d en movies", ratingMovieID)
		}
	}
	return generos
}

func leerMoviesCSV(fileName string) map[int]string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error al abrir archivo %s: %v", fileName, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error al leer archivo CSV %s: %v", fileName, err)
	}

	movies := make(map[int]string)
	for _, record := range records[1:] { // Ignorar encabezado
		movieID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Error al convertir MovieID en %s: %v", fileName, err)
			continue
		}
		movies[movieID] = record[2] // Suponiendo que `Genres` está en la tercera columna
	}
	return movies
}

func main() {
	start := time.Now()
	N := 10
	nombre := "ratings.csv"
	dividir(nombre, N)
	movies := leerMoviesCSV("movies.csv")
	vDatos := make([][]string, N)
	var wg sync.WaitGroup

	for n := 1; n <= N; n++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Printf("Iniciando procesamiento para archivo rating%d.csv\n", n)
			vDatos[n-1] = resultados(n, movies)
			fmt.Printf("Finalizado procesamiento para archivo rating%d.csv\n", n)
		}(n)
	}
	wg.Wait()

	// Contar géneros
	genreCount := make(map[string]int)
	for _, genres := range vDatos {
		for _, genre := range genres {
			genreCount[genre]++
		}
	}

	// Ordenar géneros
	type GenreRanking struct {
		Genre string
		Count int
	}
	var sortedGenres []GenreRanking
	for genre, count := range genreCount {
		sortedGenres = append(sortedGenres, GenreRanking{Genre: genre, Count: count})
	}
	sort.Slice(sortedGenres, func(i, j int) bool {
		return sortedGenres[i].Count > sortedGenres[j].Count
	})

	// Mostrar resultados
	fmt.Println("Géneros ordenados por cantidad de rankings:")
	for _, genre := range sortedGenres {
		fmt.Printf("%s: %d\n", genre.Genre, genre.Count)
	}
	duration := time.Since(start)
	fmt.Printf("Duración total = %.2f s\n", duration.Seconds())
}
