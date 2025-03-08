package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

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

// Intento de mejora que no funcionó
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

// Definimos la estructura Movie para almacenar los datos de las películas
type Movie struct {
	MovieID string `parquet:"name=movieId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	Title   string `parquet:"name=title, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	Genres  string `parquet:"name=genres, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
}

// Definimos la estructura Rating para almacenar los datos de las calificaciones
type Rating struct {
	UserID    string  `parquet:"name=userId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	MovieID   string  `parquet:"name=movieId, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	Rating    float32 `parquet:"name=rating, type=FLOAT"`
	Timestamp int64   `parquet:"name=timestamp, type=INT64"`
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

func countAndSumMovieRatings_limited(m_path string, r_paths []string) map[string]map[string]float32 {
	var wg sync.WaitGroup
	genreRatingData := make(map[string]map[string]float32) // Mapa global para los resultados
	var mu sync.Mutex                                      // Mutex para proteger el acceso a genreRatingData

	// Dividimos los archivos en dos grupos
	mid := len(r_paths) / 2
	groups := [][]string{r_paths[:mid], r_paths[mid:]}

	// Función para procesar un grupo de archivos
	processGroup := func(files []string) {
		defer wg.Done() // Señal al WaitGroup de que esta goroutine terminó

		// Mapa local para acumular datos del grupo
		localGenreData := make(map[string]map[string]float32)

		for _, r_path := range files {
			// Carga de películas en memoria para asociar géneros
			fr, err := local.NewLocalFileReader(m_path)
			if err != nil {
				log.Fatalf("No se pudo abrir el archivo de películas: %v", err)
			}
			defer fr.Close()

			pr, err := reader.NewParquetReader(fr, new(Movie), 4)
			if err != nil {
				log.Fatalf("Error al crear el lector Parquet: %v", err)
			}
			defer pr.ReadStop()

			numMovies := int(pr.GetNumRows())
			movies := make([]Movie, numMovies)
			if err = pr.Read(&movies); err != nil {
				log.Fatalf("Error al leer datos de películas: %v", err)
			}

			movieGenresMap := make(map[string]string)
			for _, movie := range movies {
				movieGenresMap[movie.MovieID] = movie.Genres
			}

			// Carga y procesamiento de ratings
			fr2, err := local.NewLocalFileReader(r_path)
			if err != nil {
				log.Fatalf("No se pudo abrir el archivo de calificaciones: %v", err)
			}
			defer fr2.Close()

			pr2, err := reader.NewParquetReader(fr2, new(Rating), 4)
			if err != nil {
				log.Fatalf("Error al crear el lector Parquet para calificaciones: %v", err)
			}
			defer pr2.ReadStop()

			numRatings := int(pr2.GetNumRows())
			ratings := make([]Rating, numRatings)
			if err = pr2.Read(&ratings); err != nil {
				log.Fatalf("Error al leer datos de calificaciones: %v", err)
			}

			// Contar y sumar ratings por género
			for _, rating := range ratings {
				if genres, exists := movieGenresMap[rating.MovieID]; exists {
					for _, genre := range strings.Split(genres, "|") {
						if _, ok := localGenreData[genre]; !ok {
							localGenreData[genre] = map[string]float32{"count": 0, "sum": 0.0}
						}
						localGenreData[genre]["count"]++
						localGenreData[genre]["sum"] += rating.Rating
					}
				}
			}
		}

		// Bloqueamos el acceso al mapa global para agregar resultados del grupo
		mu.Lock()
		for genre, data := range localGenreData {
			if _, ok := genreRatingData[genre]; !ok {
				genreRatingData[genre] = map[string]float32{"count": 0, "sum": 0.0}
			}
			genreRatingData[genre]["count"] += data["count"]
			genreRatingData[genre]["sum"] += data["sum"]
		}
		mu.Unlock()
	}

	// Lanzamos solo 2 goroutines para los dos grupos de archivos
	for _, group := range groups {
		wg.Add(1)
		go processGroup(group)
	}

	wg.Wait() // Esperamos a que ambas goroutines terminen

	// Calculamos el promedio de ratings para cada género
	for genre, data := range genreRatingData {
		if data["count"] > 0 {
			genreRatingData[genre]["average"] = data["sum"] / data["count"]
		}
	}

	return genreRatingData
}
