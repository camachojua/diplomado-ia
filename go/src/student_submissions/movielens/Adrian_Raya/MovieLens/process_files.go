package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

// Definición de las estructuras para las películas y calificaciones
// Definir la estructura Movie

type Movie struct {
	MovieID string `parquet:"name=movieId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	Title   string `parquet:"name=title, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	Genres  string `parquet:"name=genres, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
}

type Rating struct {
	UserID    string  `parquet:"name=userId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	MovieID   string  `parquet:"name=movieId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	Rating    float32 `parquet:"name=rating, type=FLOAT"`
	Timestamp int64   `parquet:"name=timestamp, type=INT64"`
}

/*
var knownGenres = []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
	"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
	"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}
*/

func count_MovieRating(m_path string, r_path string) {

	fr, err := local.NewLocalFileReader(m_path)
	if err != nil {
		log.Fatalf("No se puede abrir el archivo: %v", err)
	}
	defer fr.Close()

	// Crear un nuevo reader de Parquet
	pr, err := reader.NewParquetReader(fr, new(Movie), 4)
	if err != nil {
		log.Fatalf("Error al crear el ParquetReader: %v", err)
	}
	defer pr.ReadStop()

	// Leer las películas
	num := int(pr.GetNumRows())
	movies := make([]Movie, num)
	if err = pr.Read(&movies); err != nil {
		log.Fatalf("Error al leer los datos: %v", err)
	}

	// Mostrar las primeras películas
	for i := 0; i < 5 && i < len(movies); i++ {
		fmt.Printf("MovieID: %s, Title: %s, Genres: %s\n", movies[i].MovieID, movies[i].Title, movies[i].Genres)
	}

	//Lee movies bien

	// Leer el archivo de ratings
	fr2, err2 := local.NewLocalFileReader(r_path)
	if err2 != nil {
		log.Fatalf("No se puede abrir el archivo: %v", err2)
	}
	defer fr2.Close()

	// Crear un nuevo reader de Parquet
	pr2, err2 := reader.NewParquetReader(fr2, new(Rating), 4) // Aquí usamos fr2 en lugar de fr
	if err2 != nil {
		log.Fatalf("Error al crear el ParquetReader: %v", err2)
	}
	defer pr2.ReadStop()

	// Leer las calificaciones
	num2 := int(pr2.GetNumRows())
	ratings := make([]Rating, num2)
	if err2 = pr2.Read(&ratings); err2 != nil {
		log.Fatalf("Error al leer los datos: %v", err2)
	}

	// Mostrar las primeras calificaciones
	for i := 0; i < 5 && i < len(ratings); i++ {
		fmt.Printf("UserID: %s, MovieID: %s, Rating: %f, Timestamp: %d\n", ratings[i].UserID, ratings[i].MovieID, ratings[i].Rating, ratings[i].Timestamp)
	}
}

func count_MovieRating2(m_path string, r_path string) {
	// Leer el archivo de movies
	movieFile := m_path
	fr, err := local.NewLocalFileReader(movieFile)
	if err != nil {
		log.Fatalf("No se puede abrir el archivo: %v", err)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(Movie), 4)
	if err != nil {
		log.Fatalf("Error al crear el ParquetReader: %v", err)
	}
	defer pr.ReadStop()

	// Leer todas las películas
	numMovies := int(pr.GetNumRows())
	movies := make([]Movie, numMovies)
	if err = pr.Read(&movies); err != nil {
		log.Fatalf("Error al leer los datos de movies: %v", err)
	}

	// Crear un mapa de MovieID para búsqueda rápida
	movieMap := make(map[string]bool)
	for _, movie := range movies {
		movieMap[movie.MovieID] = true
	}

	// Leer el archivo de ratings
	ratingFile := r_path
	fr2, err2 := local.NewLocalFileReader(ratingFile)
	if err2 != nil {
		log.Fatalf("No se puede abrir el archivo: %v", err2)
	}
	defer fr2.Close()

	pr2, err2 := reader.NewParquetReader(fr2, new(Rating), 4)
	if err2 != nil {
		log.Fatalf("Error al crear el ParquetReader: %v", err2)
	}
	defer pr2.ReadStop()

	// Leer todas las calificaciones
	numRatings := int(pr2.GetNumRows())
	ratings := make([]Rating, numRatings)
	if err2 = pr2.Read(&ratings); err2 != nil {
		log.Fatalf("Error al leer los datos de ratings: %v", err2)
	}

	// Filtrar calificaciones por MovieID
	filteredRatings := []Rating{}
	for _, rating := range ratings {
		if _, exists := movieMap[rating.MovieID]; exists {
			filteredRatings = append(filteredRatings, rating)
		}
	}

	// Mostrar las primeras calificaciones filtradas
	for i := 0; i < 5 && i < len(filteredRatings); i++ {
		fmt.Printf("UserID: %s, MovieID: %s, Rating: %f, Timestamp: %d\n", filteredRatings[i].UserID, filteredRatings[i].MovieID, filteredRatings[i].Rating, filteredRatings[i].Timestamp)
	}

	fmt.Printf("Total calificaciones filtradas: %d\n", len(filteredRatings))
}

func count_MovieRating3(m_path string, r_path string) map[string]int {
	// Lista de géneros conocidos
	knownGenres := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	// Crear un mapa para contar los géneros
	genreCount := make(map[string]int)

	// Leer el archivo de movies
	movieFile := m_path
	fr, err := local.NewLocalFileReader(movieFile)
	if err != nil {
		log.Fatalf("No se puede abrir el archivo: %v", err)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(Movie), 4)
	if err != nil {
		log.Fatalf("Error al crear el ParquetReader: %v", err)
	}
	defer pr.ReadStop()

	// Leer todas las películas
	numMovies := int(pr.GetNumRows())
	movies := make([]Movie, numMovies)
	if err = pr.Read(&movies); err != nil {
		log.Fatalf("Error al leer los datos de movies: %v", err)
	}

	// Crear un mapa de MovieID -> Géneros
	movieGenresMap := make(map[string]string)
	for _, movie := range movies {
		movieGenresMap[movie.MovieID] = movie.Genres
	}

	// Leer el archivo de ratings
	ratingFile := r_path
	fr2, err2 := local.NewLocalFileReader(ratingFile)
	if err2 != nil {
		log.Fatalf("No se puede abrir el archivo: %v", err2)
	}
	defer fr2.Close()

	pr2, err2 := reader.NewParquetReader(fr2, new(Rating), 4)
	if err2 != nil {
		log.Fatalf("Error al crear el ParquetReader: %v", err2)
	}
	defer pr2.ReadStop()

	// Leer todas las calificaciones
	numRatings := int(pr2.GetNumRows())
	ratings := make([]Rating, numRatings)
	if err2 = pr2.Read(&ratings); err2 != nil {
		log.Fatalf("Error al leer los datos de ratings: %v", err2)
	}

	// Filtrar calificaciones por MovieID y contar géneros
	for _, rating := range ratings {
		if genres, exists := movieGenresMap[rating.MovieID]; exists {
			// Dividir los géneros por el carácter '|'
			genreList := strings.Split(genres, "|")
			for _, genre := range genreList {
				if contains(knownGenres, genre) {
					genreCount[genre]++
				} else {
					genreCount["Unknown"]++
				}
			}
		}
	}

	// Mostrar los resultados del conteo de géneros
	fmt.Println("Conteo de géneros:")
	for genre, count := range genreCount {
		fmt.Printf("%s: %d\n", genre, count)
	}

	return genreCount
}

// Función para verificar si un género es conocido
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func count_MovieRating4(m_path string, r_path string) map[string]int {
	// Lista de géneros conocidos
	knownGenres := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	// Crear un mapa para contar los géneros
	genreCount := make(map[string]int)

	// Leer el archivo de movies
	movieFile := m_path
	fr, err := local.NewLocalFileReader(movieFile)
	if err != nil {
		log.Fatalf("No se puede abrir el archivo: %v", err)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(Movie), 4)
	if err != nil {
		log.Fatalf("Error al crear el ParquetReader: %v", err)
	}
	defer pr.ReadStop()

	// Leer todas las películas
	numMovies := int(pr.GetNumRows())
	movies := make([]Movie, numMovies)
	if err = pr.Read(&movies); err != nil {
		log.Fatalf("Error al leer los datos de movies: %v", err)
	}

	// Crear un mapa de MovieID -> Géneros
	movieGenresMap := make(map[string]string)
	for _, movie := range movies {
		movieGenresMap[movie.MovieID] = movie.Genres
	}

	// Leer el archivo de ratings
	ratingFile := r_path
	fr2, err2 := local.NewLocalFileReader(ratingFile)
	if err2 != nil {
		log.Fatalf("No se puede abrir el archivo: %v", err2)
	}
	defer fr2.Close()

	pr2, err2 := reader.NewParquetReader(fr2, new(Rating), 4)
	if err2 != nil {
		log.Fatalf("Error al crear el ParquetReader: %v", err2)
	}
	defer pr2.ReadStop()

	// Leer todas las calificaciones
	numRatings := int(pr2.GetNumRows())
	ratings := make([]Rating, numRatings)
	if err2 = pr2.Read(&ratings); err2 != nil {
		log.Fatalf("Error al leer los datos de ratings: %v", err2)
	}

	// Filtrar calificaciones por MovieID y contar géneros
	for _, rating := range ratings {
		if genres, exists := movieGenresMap[rating.MovieID]; exists {
			// Dividir los géneros por el carácter '|'
			genreList := strings.Split(genres, "|")
			for _, genre := range genreList {
				if contains(knownGenres, genre) {
					genreCount[genre]++
				} else {
					genreCount["Unknown"]++
				}
			}
		}
	}

	// Devolver el conteo de géneros
	return genreCount
}
