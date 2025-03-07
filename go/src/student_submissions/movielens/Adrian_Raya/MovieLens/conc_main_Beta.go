package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

// Definimos la estructura Movie para almacenar los datos de las películas
type Movie struct {
	MovieID string parquet:"name=movieId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"
	Title   string parquet:"name=title, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"
	Genres  string parquet:"name=genres, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"
}

// Definimos la estructura Rating para almacenar los datos de las calificaciones
type Rating struct {
	UserID    string  parquet:"name=userId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"
	MovieID   string  parquet:"name=movieId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"
	Rating    float32 parquet:"name=rating, type=FLOAT"
	Timestamp int64   parquet:"name=timestamp, type=INT64"
}

// Función para contar y sumar las calificaciones de películas, usando concurrencia
func countAndSumMovieRatings(m_path string, r_paths []string) map[string]map[string]float32 {
	var wg sync.WaitGroup                                  // Usamos un WaitGroup para esperar a que todas las goroutines terminen
	genreRatingData := make(map[string]map[string]float32) // Mapa para guardar el conteo y suma de ratings por género
	var mu sync.Mutex                                      // Mutex para controlar el acceso a genreRatingData y evitar problemas de concurrencia

	// Recorremos cada archivo de ratings para procesarlos
	for _, r_path := range r_paths {
		wg.Add(1) // Añadimos una goroutine al WaitGroup
		go func(r_path string) {
			defer wg.Done() // Avisa que terminó al WaitGroup

			// Abrimos el archivo de películas
			fr, err := local.NewLocalFileReader(m_path)
			if err != nil {
				log.Fatalf("No se pudo abrir el archivo de películas: %v", err)
			}
			defer fr.Close()

			// Creamos un lector Parquet para las películas
			pr, err := reader.NewParquetReader(fr, new(Movie), 4)
			if err != nil {
				log.Fatalf("Error al crear el lector Parquet para películas: %v", err)
			}
			defer pr.ReadStop()

			// Leemos todas las películas y creamos un mapa MovieID -> Géneros
			numMovies := int(pr.GetNumRows())
			movies := make([]Movie, numMovies)
			if err = pr.Read(&movies); err != nil {
				log.Fatalf("Error al leer datos de películas: %v", err)
			}

			// Guardamos el ID de la película y sus géneros en un mapa
			movieGenresMap := make(map[string]string)
			for _, movie := range movies {
				movieGenresMap[movie.MovieID] = movie.Genres
			}

			// Abrimos el archivo de calificaciones
			fr2, err := local.NewLocalFileReader(r_path)
			if err != nil {
				log.Fatalf("No se pudo abrir el archivo de calificaciones: %v", err)
			}
			defer fr2.Close()

			// Creamos un lector Parquet para las calificaciones
			pr2, err := reader.NewParquetReader(fr2, new(Rating), 4)
			if err != nil {
				log.Fatalf("Error al crear el lector Parquet para calificaciones: %v", err)
			}
			defer pr2.ReadStop()

			// Leemos todas las calificaciones y las procesamos
			numRatings := int(pr2.GetNumRows())
			ratings := make([]Rating, numRatings)
			if err = pr2.Read(&ratings); err != nil {
				log.Fatalf("Error al leer datos de calificaciones: %v", err)
			}

			// Mapa local para guardar datos de géneros por cada archivo de calificaciones
			localGenreData := make(map[string]map[string]float32)

			// Procesamos cada calificación
			for _, rating := range ratings {
				if genres, exists := movieGenresMap[rating.MovieID]; exists { // Verificamos que el ID de la película exista en el mapa de géneros
					for _, genre := range strings.Split(genres, "|") { // Dividimos los géneros (pueden ser varios)
						if _, ok := localGenreData[genre]; !ok {
							localGenreData[genre] = map[string]float32{"count": 0, "sum": 0.0} // Iniciamos el género si no está en el mapa
						}
						localGenreData[genre]["count"]++              // Aumentamos el conteo
						localGenreData[genre]["sum"] += rating.Rating // Sumamos el rating
					}
				}
			}

			// Pasamos los datos locales al mapa global con Mutex para evitar errores de concurrencia
			mu.Lock()
			for genre, data := range localGenreData {
				if _, ok := genreRatingData[genre]; !ok {
					genreRatingData[genre] = map[string]float32{"count": 0, "sum": 0.0}
				}
				genreRatingData[genre]["count"] += data["count"]
				genreRatingData[genre]["sum"] += data["sum"]
			}
			mu.Unlock()
		}(r_path)
	}

	wg.Wait() // Esperamos a que terminen todas las goroutines

	// Calculamos el promedio de rating para cada género
	for genre, data := range genreRatingData {
		if data["count"] > 0 {
			genreRatingData[genre]["average"] = data["sum"] / data["count"]
		}
	}

	return genreRatingData
}

func main() {
	// Ruta de prueba para archivos de películas y calificaciones
	var sMovieParquetPath = "C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/movie_go.parquet"
	sRatingParquetPaths := []string{"C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/Ratings_1.parquet",
		"C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/Ratings_2.parquet",
		//"C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/Ratings_3.parquet",
		//"C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/Ratings_4.parquet",
		//"C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/Ratings_5.parquet",
		//"C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/Ratings_6.parquet",
		//"C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/Ratings_7.parquet",
		//"C:/Users/PC/Desktop/DiplomadoIA/Go/MovieLens/Ratings_8.parquet",
	}

	// Llamamos a la función para contar y sumar los ratings por género
	genreRatingData := countAndSumMovieRatings(sMovieParquetPath, sRatingParquetPaths)

	// Imprimimos los resultados de conteo, suma y promedio por género
	fmt.Println("Datos de Rating por Género:")
	for genre, data := range genreRatingData {
		fmt.Printf("Género: %s, Conteo: %v, Suma: %v, Promedio: %.2f\n", genre, data["count"], data["sum"], data["average"])
	}
}