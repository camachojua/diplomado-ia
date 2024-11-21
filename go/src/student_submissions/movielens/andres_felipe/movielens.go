package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/kfultz07/go-dataframe"
)

// ====================================================================

// Esta función se encarga de leer un archivo CSV
func readCSV(directory, file_name string) [][]string {
	//Se abre el archivo
	filePath := fmt.Sprintf("%s/%s", directory, file_name)
	file, err := os.Open(filePath)
	if err != nil {
		//Muestra el error y termina el programa si no se puede leer el archivo CSV
		log.Fatal("No se pudo abrir el archivo CSV: ", err)
	}
	defer file.Close()
	//Se lee el archivo linea por linea para mejorar el tiempo de ejecución
	csvReader := csv.NewReader(file)
	var records [][]string
	for {
		line, err := csvReader.Read()
		if err != nil {
			break
		}
		records = append(records, line)
	}
	return records
}

// Esta función se encarga de crear un archivo CSV a partir de un conjunto de datos
func writeCSV(directory, file_name string, data [][]string) {
	//Creación del CSV
	filePath := fmt.Sprintf("%s/%s", directory, file_name)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("No se pudo crear el archivo CSV: ", err)
	}
	defer file.Close()

	// Escribir los datos en el archivo CSV
	fileWriter := csv.NewWriter(file)
	err = fileWriter.WriteAll(data)
	if err != nil {
		log.Fatal("No se puede escribir en CSV: ", err)
	}
	fileWriter.Flush()
}

// Función para particionar un archivo de datos grande,
// donde file_name es el nombre del CSV grande; partitions es el número de minis CSV en el que se dividirán
// los datos; directory es el directorio donde se encuentra el CSV
func SplitBigFile(directory string, file_name string, partitions int) []string {

	startTime := time.Now()

	println("Inicio de lectura y partición del archivo CSV")

	//Lectura del archivo
	records := readCSV(directory, file_name)

	//Número de registos/lineas que contiene el archivo
	n_lines := len(records)
	fmt.Printf("El archivo CSV consta de %v registros \n", n_lines)

	//Número de registros por cada partición
	np_lines := int(math.Floor(float64(n_lines) / float64(partitions)))
	fmt.Printf("Por cada partición habrá %v registros \n", np_lines)

	//Concurrencia
	ch := make(chan string, partitions)
	var wg sync.WaitGroup
	for n := 0; n < partitions; n++ {
		wg.Add(1)
		go splits(directory, n, partitions, np_lines, records, ch, &wg)
	}

	// Espera a que todas las goroutines terminen
	wg.Wait()
	close(ch)

	// Este slice contendrá el nombre de los minis CSV generados
	var fileNames []string
	for fileName := range ch {
		fileNames = append(fileNames, fileName)
	}

	fmt.Println("La partición tomó", time.Since(startTime).Seconds(), "segundos")
	return fileNames
}

// Función de escritura de cada partición
func splits(
	directory string,
	n_partition int,
	partitions int,
	np_lines int,
	data [][]string,
	ch chan string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	//Nombre del archivo
	mini_fileName := fmt.Sprintf("partition_%v.csv", n_partition)
	//Limites inicial y final de los datos
	start := n_partition*np_lines + 1
	end := start + np_lines
	if n_partition == partitions-1 {
		end = len(data)
	}
	//Cada partición debe contener los nombres de cada columna
	data = append([][]string{data[0]}, data[start:end]...)
	//Creación del mini CSV de la n-partición
	writeCSV(directory, mini_fileName, data)

	//Manda el nombre del archivo CSV generado para después agregarlo al slice
	ch <- mini_fileName
}

// ====================================================================

func procesaArchivoMultiHilo(directory1, directory2 string, threads int) {
	println("El orquestador del proceso ha iniciado su ejecución")

	generos := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	numero_generos := len(generos)

	// Matriz de 2 dimensiones (Número de generos x Número de hilos/workers/particiones)
	// Guarda la suma de las calificaciones para un género,
	// dado los registros en una partición de los datos.
	arreglo_calificacion := make([][]float64, numero_generos)

	// Matriz de 2 dimensiones (Número de generos x Número de particiones)
	// Guarda el número de veces que aparece un género en los registros de una partición
	arreglo_conteo := make([][]int, numero_generos)

	for i := 0; i < numero_generos; i++ {
		arreglo_calificacion[i] = make([]float64, threads)
		arreglo_conteo[i] = make([]int, threads)
	}

	// Creamos el canal de comunicación
	var ch = make(chan int)

	// Se crea un DataFrame del archivo CSV de movies
	movies := dataframe.CreateDataFrame(directory2, "movies.csv")

	// Creamos 10 workers, cada worker se encarga de leer su archivo/partición correspondiente
	for i := 0; i < threads; i++ {
		go encuentraCalificaciones(i, ch, generos, &arreglo_conteo, &arreglo_calificacion, movies, directory1)
	}

	// Esperamos que los workers terminen de trabajar
	iMsg := 0
	// Consultamos el valor del canal de comunicación cuando un worker
	// termina de trabajar grita ¡Ya acabé!, ese grito puede ser escuchado a
	// través del canal de comunicación
	go func() {
		for {
			i := <-ch
			iMsg += i
		}
	}()
	for {
		if iMsg == threads {
			break
		}
	}

	// Acá "consolidamos la información"
	locCount := make([]int, numero_generos)    // Conteo total por género
	locVals := make([]float64, numero_generos) // Suma total de ratings
	locMean := make([]float64, numero_generos) // Calificación promedio por género
	for i := 0; i < numero_generos; i++ {
		for j := 0; j < threads; j++ {
			locCount[i] += arreglo_conteo[i][j]
			locVals[i] += arreglo_calificacion[i][j]
		}
		locMean[i] = locVals[i] / float64(locCount[i])
	}

	// Acá imprimimos los resultados
	fmt.Println(" n   ", fmt.Sprintf("%20s", "GÉNERO"), "  ", fmt.Sprintf("%8s", "CONTEO"), "  ", fmt.Sprintf("%3s", "CAL. PROM."))
	for i := 0; i < numero_generos; i++ {
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", generos[i]), "  ", fmt.Sprintf("%8d", locCount[i]), "  ", fmt.Sprintf("%.2f", locMean[i]))
	}

	println("Fin del orquestador.")
}

func encuentraCalificaciones(
	worker_id int,
	ch chan int,
	generos_conocidos []string,
	arreglo_conteo *[][]int,
	arreglo_calificacion *[][]float64,
	movies dataframe.DataFrame,
	directory string,
) {
	partitionCSV := "partition_" + fmt.Sprintf("%v", worker_id) + ".csv"
	fmt.Println("Worker  ", fmt.Sprintf("%v", worker_id), " procesará el archivo ", partitionCSV)

	// Se crea un DataFrame a aprtir del archivo CSV de la partición
	partitionDF := dataframe.CreateDataFrame(directory, partitionCSV)

	tiempo_inicial := time.Now()

	// Inner-join entre los DataFrames de movies y (partición de) ratings
	// Nos interesa los datos de genres del datagrame movies
	partitionDF.Merge(&movies, "movieId", "genres")

	// De todo el dataset sólo nos interesan los generos y los ratings,
	// entonces mantenemos únicamente esas columnas del dataframe
	columnas_que_nos_interesan := [2]string{"genres", "rating"}
	ratings_genre := partitionDF.KeepColumns(columnas_que_nos_interesan[:])

	// Iteración por género
	for g, genero := range generos_conocidos {
		for _, row := range ratings_genre.FrameRecords {
			//row.Data[0] son los géneros de la pelicula, row.Data[1] es el rating
			if strings.Contains(row.Data[0], genero) == true {
				r, _ := strconv.ParseFloat(row.Data[1], 64)
				(*arreglo_calificacion)[g][worker_id] += r
				(*arreglo_conteo)[g][worker_id] += 1
			}
		}
	}

	tiempo_final := time.Since(tiempo_inicial)
	fmt.Println("Worker ", worker_id, " ha terminado 🟢")
	fmt.Println("Tiempo en procesar = ", tiempo_final)
	ch <- 1
}

func main() {
	runtime.GOMAXPROCS(12)

	//división
	directory1 := "/home/ws117/diplomado-ia/go/src/my_code/movielens/ratings"
	directory2 := "/home/ws117/diplomado-ia/go/src/my_code/movielens/movies"
	file_name := "ratings_p.csv"
	partitions := 10

	SplitBigFile(directory1, file_name, partitions)
	procesaArchivoMultiHilo(directory1, directory2, partitions)
}