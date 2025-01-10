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
)

// Rating representa una calificación de usuario por película
type Rating struct {
	UserId  int
	MovieId int
	Rating  float64
}

// Movie representa la info de cada película
type Movie struct {
	MovieId int
	Title   string
	Genres  string
}

// Esta función lee cualquier archivo CSV y usa una función para convertir cada fila en el tipo adecuado
func readCSV[T any](filePath string, mapFunc func([]string) T) ([]T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []T
	for _, record := range records[1:] { // Ignoramos la primera fila, que son los encabezados
		data = append(data, mapFunc(record))
	}
	return data, nil
}

// Convierte cada fila del CSV a Rating
func mapRating(record []string) Rating {
	userId, _ := strconv.Atoi(record[0])
	movieId, _ := strconv.Atoi(record[1])
	rating, _ := strconv.ParseFloat(record[2], 64)
	return Rating{
		UserId:  userId,
		MovieId: movieId,
		Rating:  rating,
	}
}

// Convierte cada fila del CSV a Movie
func mapMovie(record []string) Movie {
	movieId, _ := strconv.Atoi(record[0])
	title := record[1]
	genres := record[2]
	return Movie{
		MovieId: movieId,
		Title:   title,
		Genres:  genres,
	}
}

// Divide el archivo de ratings en chunks para facilitar el procesamiento
func fileSplitter[T any](n int, data []T, path string, fileNamePrefix string, writeFunc func(fileName string, chunk []T)) {
	numRows := len(data)
	rowsPerFile := numRows / n
	remainingRows := numRows % n

	startTime := time.Now()

	for i := 0; i < n; i++ {
		startIndex := i * rowsPerFile
		var endIndex int
		if i == n-1 {
			endIndex = startIndex + rowsPerFile + remainingRows
		} else {
			endIndex = startIndex + rowsPerFile
		}

		fileName := fmt.Sprintf("%s/%s_%d.csv", path, fileNamePrefix, i+1)
		writeFunc(fileName, data[startIndex:endIndex])

		fmt.Printf("Chunk %d: Inicio: %d, Final: %d\n", i+1, startIndex+1, endIndex)
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Tiempo total de ejecución: %v\n", elapsedTime)
}

// Escribe un chunk de Ratings en un CSV
func writeCSVRatings(fileName string, chunk []Rating) {
	output, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	defer writer.Flush()

	headers := []string{"User Id", "MovieId", "Rating"}
	if err := writer.Write(headers); err != nil {
		log.Fatal(err)
	}

	for _, record := range chunk {
		row := []string{
			strconv.Itoa(record.UserId),
			strconv.Itoa(record.MovieId),
			fmt.Sprintf("%.4f", record.Rating),
		}
		if err := writer.Write(row); err != nil {
			log.Fatal(err)
		}
	}
}

// Lee todos los archivos CSV que dividimos antes y los recombina en un mapa
func readSeparatedFiles(baseName string, n int, path string) (map[string][]Rating, error) {
	dataframes := make(map[string][]Rating)

	for i := 1; i <= n; i++ {
		fileName := fmt.Sprintf("%s/%s_%d.csv", path, baseName, i)
		varName := fmt.Sprintf("%s_%d", baseName, i)

		data, err := readCSV(fileName, mapRating)
		if err != nil {
			return nil, fmt.Errorf("error leyendo %s: %v", fileName, err)
		}

		dataframes[varName] = data
	}

	return dataframes, nil
}

// Cuenta los géneros a partir de ratings y películas
func calcularGeneros(ratings []Rating, movies []Movie) map[string]int {
	contador := make(map[string]int)
	movieGenres := make(map[int]string)
	for _, movie := range movies {
		movieGenres[movie.MovieId] = movie.Genres
	}

	for _, rating := range ratings {
		if genres, exists := movieGenres[rating.MovieId]; exists {
			generos := strings.Split(genres, "|")
			for _, genero := range generos {
				contador[strings.TrimSpace(genero)]++
			}
		}
	}

	return contador
}

func main() {
	ratings, err := readCSV("/Users/rodrigozeferino/Documents/Diplomado/Go/Movielens/ml-25m/ratings.csv", mapRating)
	if err != nil {
		log.Fatal(err)
	}

	movies, err := readCSV("/Users/rodrigozeferino/Documents/Diplomado/Go/Movielens/ml-25m/movies.csv", mapMovie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Filas ratings: %d, Filas movies: %d\n", len(ratings), len(movies))

	n := 10
	path := "/Users/rodrigozeferino/Documents/Diplomado/Go/Movielens/ml-25m/"

	fileSplitter(n, ratings, path, "ratings", writeCSVRatings)

	dataframes, err := readSeparatedFiles("ratings", n, path)
	if err != nil {
		log.Fatal(err)
	}

	totalRows := 0
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("ratings_%d", i+1)
		if data, exists := dataframes[key]; exists {
			totalRows += len(data)
		} else {
			fmt.Printf("Key %s no existe.\n", key)
		}
	}
	fmt.Printf("Total de filas: %d\n", totalRows)

	kg := []string{
		"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)",
	}

	ra := make([][]float64, len(kg))
	for i := range ra {
		ra[i] = make([]float64, n)
	}

	ca := make([][]int, len(kg))
	for i := range ca {
		ca[i] = make([]int, n)
	}

	numFiles := 10
	results := make([]map[string]int, numFiles)
	var wg sync.WaitGroup

	for i := 1; i <= numFiles; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			filePath := fmt.Sprintf("%sratings_%d.csv", path, i)
			ratings, err := readCSV(filePath, mapRating)
			if err != nil {
				log.Printf("Error leyendo archivo %s: %v\n", filePath, err)
				results[i-1] = make(map[string]int)
				return
			}
			results[i-1] = calcularGeneros(ratings, movies)
		}(i)
	}

	wg.Wait()

	for i, contador := range results {
		fmt.Printf("Géneros en ratings_%d.csv:\n", i+1)
		for genero, count := range contador {
			fmt.Printf("%s: %d\n", genero, count)
		}
	}

	sumaTotal := make(map[string]int)
	for _, diccionario := range results {
		if diccionario != nil {
			for genero, conteo := range diccionario {
				sumaTotal[genero] += conteo
			}
		}
	}

	fmt.Println("Total de géneros:", sumaTotal)
}
