/* Autor Martha Rico Diener
Equipo Los cantantes
Este programa procesa archivos CSV que contienen calificaciones de películas y
estadísticas sobre géneros. Utiliza concurrencia para manejar múltiples archivos
de calificaciones y calcula el promedio de calificaciones por género.
*/

package main

/*Se importan varios paquetes necesarios para la funcionalidad del programa:
encoding/csv para manejar archivos CSV
fmt para imprimir en consola
log para registrar errores
os para operaciones de archivos
sort para ordenar datos,
strings para manipulación de cadenas
sync para gestionar concurrencia
time para medir el tiempo de ejecución.
*/

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
	"strconv"
)

func dividir() { //Esta función divide el archivo ratings en 10 partes iguales

	startTime := time.Now() //Para medir el tiempo de ejecución

	// Cargar el archivo CSV
	df, headers, err := loadCSV("ratings.csv")
	if err != nil {
		panic(err)
	}

	// Definir el número de partes en las que quieres dividir el archivo
	N := 10

	// Calcular el tamaño de cada parte
	numRows := len(df)
	rowsPerPart := numRows / N

	// Crear un loop para dividir el DataFrame y guardar cada parte
	for i := 0; i < N; i++ {
		startRow := i * rowsPerPart
		endRow := startRow + rowsPerPart

		// Asegurarse de no exceder el número de filas
		if i == N-1 {
			endRow = numRows
		}

		// Crear y guardar la parte como un nuevo archivo CSV
		partDF := df[startRow:endRow]
		err := saveCSV(fmt.Sprintf("ratings_part_%d.csv", i+1), headers, partDF)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("El archivo ha sido dividido en %d partes.\n", N)
	elapsedTime := time.Since(startTime)
	fmt.Printf("Tiempo transcurrido hasta este moemnto: %s\n", elapsedTime)
}

// loadCSV carga un archivo CSV y devuelve los datos como una matriz de cadenas y los encabezados
func loadCSV(filename string) ([][]string, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	// Separar los encabezados de los datos
	headers := records[0]
	data := records[1:]

	return data, headers, nil
}

// saveCSV guarda una parte de datos en un nuevo archivo CSV con encabezados
func saveCSV(filename string, headers []string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir los encabezados en el archivo CSV
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Escribir las filas en el archivo CSV
	return writer.WriteAll(data)
}

/*
 **Función `processRatings`**:
   - Toma el nombre de un archivo de calificaciones,
   un mapa para almacenar estadísticas de géneros,
   un mutex para sincronización y un `WaitGroup`
   para esperar a que se completen las goroutines.
*/

func processRatings(ratingsFile string, genreData map[string]*GenreStats, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	//Abre el archivo CSV
	file, err := os.Open(ratingsFile)
	if err != nil {
		log.Printf("Error al abrir %s: %v", ratingsFile, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	ratingsData, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error al leer %s: %v", ratingsFile, err)
		return
	}

	for _, rating := range ratingsData[1:] { // Ignorar encabezados
		movieId := rating[1]
		var ratingValue float64
		if len(rating) > 2 {
			_, err := fmt.Sscanf(rating[2], "%f", &ratingValue) // Leer la calificación como float
			if err != nil {
				log.Printf("Error al leer la calificación: %v", err)
			}
		}

		// Buscar el título y género correspondiente en moviesData
		for _, movie := range moviesData[1:] { // Ignorar encabezados
			if movie[0] == movieId {
				genres := strings.Split(movie[2], "|")
				mu.Lock()
				for _, genre := range genres {
					if _, exists := genreData[genre]; !exists {
						genreData[genre] = &GenreStats{}
					}
					genreData[genre].Count++
					genreData[genre].Sum += ratingValue
				}
				mu.Unlock()
				break
			}
		}
	}
}

var moviesData [][]string

// Estructura para almacenar estadísticas de cada género
type GenreStats struct {
	Count int
	Sum   float64
}


// Función para formatear números con comas
func formatNumber(n int) string {
	s := strconv.Itoa(n)
	nLen := len(s)
	if nLen <= 3 {
		return s
	}

	var result strings.Builder
	for i, digit := range s {
		if (nLen-i)%3 == 0 && i != 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}
	return result.String()
}

func main() {

	startTime := time.Now() //Para medir el tiempo de ejecución

	dividir() //dividir el archivo ratings.csv en 10

	// Abrir y leer el archivo movies.csv
	moviesFile, err := os.Open("movies.csv")
	if err != nil {
		log.Fatalf("Error al abrir movies.csv: %v", err)
	}
	defer moviesFile.Close()

	moviesReader := csv.NewReader(moviesFile)
	moviesData, err = moviesReader.ReadAll()
	if err != nil {
		log.Fatalf("Error al leer movies.csv: %v", err)
	}

	// Mapa para almacenar estadísticas de calificaciones por género
	genreData := make(map[string]*GenreStats)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Procesar cada archivo de ratings concurrentemente
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go processRatings(fmt.Sprintf("ratings_part_%d.csv", i), genreData, &mu, &wg)
	}

	wg.Wait()

	// Crear un slice para almacenar los géneros y sus estadísticas
	type genreStat struct {
		genre   string
		count   int
		average float64
	}

	var sortedGenres []genreStat
	for genre, stats := range genreData {
		average := 0.0
		if stats.Count > 0 {
			average = float64(stats.Sum) / float64(stats.Count)
		}
		sortedGenres = append(sortedGenres, genreStat{genre, stats.Count, average})
	}

	// Ordenar el slice por el nombre del género
	sort.Slice(sortedGenres, func(i, j int) bool {
		return sortedGenres[i].genre < sortedGenres[j].genre
	})

/*
// Imprimir los resultados
fmt.Println("Conteo de calificaciones por género:")
for _, gs := range sortedGenres {
	fmt.Printf("%s: \t \t %s  \t \t %.2f\n", gs.genre, formatNumber(gs.count), gs.average)
}
*/

// Imprimir los resultados
fmt.Println()
fmt.Println("Conteo de calificaciones por género:")
fmt.Printf("%-20s|%20s\t|\t%20s\n", "Género", "Conteo", "Promedio de calificaciones")
fmt.Println("----------------------------------------------------------")
for _, gs := range sortedGenres {
	// Ajusta los anchos según tus necesidades
	fmt.Printf("%-20s| %20s\t|\t%6.2f\n", gs.genre, formatNumber(gs.count), gs.average)
}

	// Calcular y mostrar el tiempo total de ejecución
	elapsedTime := time.Since(startTime)
	fmt.Printf("Tiempo total de ejecución: %s\n", elapsedTime)
}
