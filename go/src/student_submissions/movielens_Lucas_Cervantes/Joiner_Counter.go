package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
Este codigo toma los archivos divididos de ratings en csv
les hace join con movies.csv, cuenta el numero de ratings
por genero y suma el valor de los ratings por genero.

No usamos dataframes en go (porque la verdad no pude hacer
que corriera en mi computadora) y en ves de eso se manejan
los datos mediante estructuras definidas para archivo y
para el "Join"

*/

/* Definimos 3 estructuras, una para Movies, una para Ratings
y una para JoinesRow */

type Movie struct {
	MovieID int64
	Title   string
	Genres  string
}

type Ratings struct {
	UserID    int64
	MovieID   int64
	Rating    float64
	Timestamp int64
}

type JoinedRow struct {
	Rating float64
	Genres string
}

// Función para leer las peliculas desde csv
func readMovies(filePath string) ([]Movie, error) {
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

	movies := make([]Movie, 0, len(records))

	// Convertimos los tipos del csv segun se necesiten en la estrcutura con strconv
	for _, record := range records {
		movieID, _ := strconv.ParseInt(record[0], 10, 64)

		// Llenado de la estrcutura
		movies = append(movies, Movie{
			MovieID: movieID,
			Title:   record[1],
			Genres:  record[2],
		})
	}

	return movies, nil
}

// Función para leer los ratings desde csv
func readRatings(filePath string) ([]Ratings, error) {
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

	ratings := make([]Ratings, 0, len(records))

	// Convertimos los tipos del csv segun se necesiten en la estrcutura con strconv
	for _, record := range records {
		UserID, _ := strconv.ParseInt(record[0], 10, 64)
		MovieID, _ := strconv.ParseInt(record[1], 10, 64)
		Rating, _ := strconv.ParseFloat(record[2], 64)
		Timestamp, _ := strconv.ParseInt(record[3], 10, 64)

		// Llenado de la estrcutura
		ratings = append(ratings, Ratings{
			UserID:    UserID,
			MovieID:   MovieID,
			Rating:    Rating,
			Timestamp: Timestamp,
		})
	}

	return ratings, nil
}

// Función para realizar el join entre películas y ratings
func joinMoviesAndRatings(movies []Movie, ratings []Ratings) []JoinedRow {

	// Mapa para guardar Ratings | Generos
	genresMap := make(map[int64]string)

	/* Para cada pelicula de la estrcutura Movie
	genresMap relaciona el MovieID con el genero
	*/
	for _, movie := range movies {
		genresMap[movie.MovieID] = movie.Genres
	}

	/* Despues para cada fila de la estrcutura rating
	genresMap busca un MovieID de ratings, si este existe
	dentro de gnresMap, entonces se hace un append sobre
	JoinedRows del Rating y el genero correspondiente al
	MovieID
	*/
	var joinedRows []JoinedRow
	for _, rating := range ratings {
		if genres, exists := genresMap[rating.MovieID]; exists {
			joinedRows = append(joinedRows, JoinedRow{
				Rating: rating.Rating,
				Genres: genres,
			})
		}
	}

	return joinedRows
}

func main() {

	// Comenzamos un timer
	start := time.Now()
	// Dividimos el archivo rating.csv mediante Splitter_main
	tiempo_split := float64(Splitter_main(3))

	// Esta funcion busca el numero de Nucleos disponibles en la PC
	numCPUs := runtime.NumCPU()
	// Aqui establesco el no. max de nucleos del CPU
	runtime.GOMAXPROCS(numCPUs)

	// Sobre la misma carpeta donde esta el codigo, se buscan coincidencias con ratings_n
	matches, err := filepath.Glob("ratings_*")
	if err != nil {
		log.Fatal(err)
	}

	if len(matches) == 0 {
		log.Fatal("No se encontraron archivos coincidentes")
	}

	// Estrcutura definida para guardar la cuenta de ratings y sumarlos
	type GenreStats struct {
		Count     int
		RatingSum float64
	}

	/* comenzamos un wg y mu, el wg es un wait group que permite esperar a
	que un grupo de gorutinas termine antes de continuar con el siguiente
	bloque del codigo

	mu es una estrcutura de exclucion mutua para evitar condiciones de carrera
	cuando un grupo de gorutinas quieren acceder a una sola variable. Podria
	evitarse el uso de mu enviendo los resultados de localcounts sobre un canal
	*/
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Mapa donde se almacenan los datos datos finales
	genreCounts := make(map[string]GenreStats)

	// Inicializacion de la variable movies_total con la funcion readMovies
	movies_total, _ := readMovies("Movies.csv")

	for _, file := range matches {
		wg.Add(1)

		// Copiamos files sobre local_file, para evitar errores dentro de las gorutinas
		local_file := file
		go func(file string) {
			defer wg.Done()

			// Iinicializamos el rating como local_rating para cada archivo fraccionado
			// del csv dentro del for de matches
			local_rating, _ := readRatings(local_file)
			joined := joinMoviesAndRatings(movies_total, local_rating)

			// localcount es un mapa con un string y una estrcutura GenreStats
			localcount := make(map[string]GenreStats)
			for _, row := range joined {
				/* Para cada fila de joined, definimos genres como el split de los Generos por |
				y el rating como el rating de la fila*/
				genres := strings.Split(row.Genres, "|")
				rating := row.Rating
				for _, genre := range genres {
					/* Usamos stats para manejar la estrcutura GenreStats,
					sumamos un 1 sobre stats.count por genero, y sumamos
					el valor del rating sobre stats.RatingSum
					al final solo guadarmos el valor de stats sobre
					localcount[genre]
					*/
					stats := localcount[genre]
					stats.Count++
					stats.RatingSum += rating
					localcount[genre] = stats

				}
			}

			/* Con mu.Lock() restringimos el acceso a la variable genreCounts
			a una sola gorutina para evitar condiciones de carrera
			*/
			mu.Lock()
			for genre, localStats := range localcount {
				globalStats := genreCounts[genre]
				globalStats.Count += localStats.Count
				globalStats.RatingSum += localStats.RatingSum
				genreCounts[genre] = globalStats
			}
			// liberamos el acceso a genreCounts para la proxima gorituna
			mu.Unlock()
		}(file)
	}

	wg.Wait()

	fmt.Print("Conteo de ratings y promedio por genero:\n")
	// Imprimimos resultados para cada genero de genreCounts
	for genre, stats := range genreCounts {
		fmt.Printf("%-20s: %-10d, %-10.3f\n", genre, stats.Count, stats.RatingSum/float64(stats.Count))

	}

	fmt.Printf("\nTiempo de Split: %v", time.Duration(tiempo_split))
	fmt.Printf("\nTiempo de procesamiento: %v", time.Since(start)-time.Duration(tiempo_split))
	fmt.Printf("\nTiempo de ejecucion total: %v", time.Since(start))

}
