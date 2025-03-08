/*
Este programa realiza la
*/
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/kfultz07/go-dataframe"
)

/*
Definicion de constantes

	const mvLDir = "/home/luna/Downloads/ml-latest/"
*/
type output struct {
	idGenre    int32
	genre      string
	numRating  float64
	meanRating float64
}

const mvLDir = "/home/luna/Downloads/ml-25m/" //Directorio origen MovieLens

func readFile(archivo string) (*os.File, func()) { //La funcion devuelve un apuntador os.File, esto es recomendable por cuestiones de eficiencia.
	datos, err := os.Open(archivo)
	if err != nil {
		panic("Error al abrir el archivo.")
	}

	cerrar_archivo := func() {
		err := datos.Close()

		if err != nil {
			panic("Error al cerrar el archivo.")
		}
	}

	return datos, cerrar_archivo
}

/*
Esta función lee los datos binarios de un archivo, los interpreta como un csv
Y los regresa en una matriz [i][j] de 2 donde:
[i] = # Renglon del CSV
[j] = # Columna del CSV
*/
func readCsvFile(archivo *os.File) [][]string { //Esta funcion devuelve el conteo de lineas del CSV
	//defer archivo.Close() // Cierra el archivo al final de la función solo si no ha sido cerrado antes

	parser := csv.NewReader(archivo) // Se lee el archivo CSV
	parser.Read()
	records, err := parser.ReadAll() // Si exite algun error este se mostrara

	if err != nil {
		panic(err)
	}

	fmt.Printf("El archivo tiene %d registros\n", len(records)) //Se imprime el total de lineas del CSV
	return records
}

// Leemos el archivo CSV para convertirlo en un DataFrame
func csvToDataFrame(archivo *string) dataframe.DataFrame {
	df := dataframe.CreateDataFrame(mvLDir, *archivo)
	return df
}

/*
	Función para procesar concurrentemente el inner join entre "movies.csv" y cada archivo

"ratings.csv", esta funcion coordina la tarea de contar el numero de ratings por genero.
contando el número de calificaciones por archivo.

		Funciones principales:
		Crea canales de comunicación entre ella y los workers y envia las ordenes a los workers mediante
		go rutinas identificandose por su workerId.
	    Esperar a que todos los workers terminen.
		Lee y consolida las respuestas recibidas de los workers.
*/
func findRatingsMaster() {
	fmt.Println("El orquestador del proceso ha iniciado su ejecución.")

	// Definimos el número de archivos a generar, también lo usamos como el
	// número de hilos que usaremos
	nWorkers := 10
	genNumber := 20

	knowngenres := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	/* Se crea un arreglo de matriz Mij para almacenar los resultados de género y el número de calificaciones
	[i]: Los renglones hacen referencia a cada genero listado.
	[j]: Las columnas contienen las calificaciones entcontradas por un worker especifico por cada uno de los géneros listados
	*/
	ratingSum := make([][]float64, genNumber)

	/* Se crea un arreglo de matriz Mij para almacenar el conteo totall que llevamos hasta el momento
	[i] del (0,19): Los renglones hacen referencia a cada genero listado.
	[j] del (0,9):Las columnas contienen el conteo total de las calificaciones encontradas por un worker específico dado un género
	*/

	ratingCount := make([][]int, genNumber)
	// Llenamos las "genNumber" filas de los arreglos
	// "ratingSum" y "ratingCount", les ponemos
	// "nWorkers" columnas
	for i := 0; i < genNumber; i++ {
		ratingSum[i] = make([]float64, nWorkers)
		ratingCount[i] = make([]int, nWorkers)
	}
	// Creamos el canal de comunicación
	var ci = make(chan int)
	moviesData := "movies.csv"
	movies := csvToDataFrame(&moviesData) // Sacamos los dataframes del archivo csv

	// Se crean 10 workers, cada worker leera su archivo correspondiente
	for i := 0; i < nWorkers; i++ {
		go findRatingsWorker(i+1, ci, knowngenres, &ratingCount, &ratingSum, movies)
	}

	// Espera a que los workers terminen de trabajar
	iMsg := 0
	// Se consulta el valor del canal de comunicación cuando un worker termina de trabajar
	go func() {
		for {
			i := <-ci
			iMsg += i
		}
	}()
	for {
		if iMsg == 10 {
			break
		}
	}

	// Consolidacion del conteo de ratings por genero
	locCount := make([]int, genNumber)
	locVals := make([]float64, genNumber)
	localMeanValues := make([]float64, genNumber)
	var outputs []output
	for i := 0; i < genNumber; i++ {
		for j := 0; j < nWorkers; j++ {
			locCount[i] += ratingCount[i][j]
			locVals[i] += ratingSum[i][j]
		}
		// Valor promedio por fenero.
		localMeanValues[i] = locVals[i] / float64(locCount[i])

		// Estructura de salida

		out := output{
			idGenre:    int32(i),
			genre:      knowngenres[i],
			numRating:  float64(locCount[i]),
			meanRating: localMeanValues[i],
		}

		outputs = append(outputs, out)

	}

	// Mostrar el resultado final
	w1 := tabwriter.NewWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w1, "Género\tidGenero\tNúmero de Valoraciones\tPromedio de Valoraciones\t")
	for _, o := range outputs {
		fmt.Fprintf(w1, "%s\t%d\t%.0f\t%.2f\t\n", o.genre, o.idGenre, o.numRating, o.meanRating)

		//fmt.Printf("Género: %s, ID: %d, Número de Valoraciones: %.0f, Promedio de Valoraciones: %.2f\n",
		//	o.genre, o.idGenre, o.numRating, o.meanRating)
	}
	w1.Flush()
	println("Fin")
}

/* Esta funcion determina los procesos que realizara cada worker:
	Cada worker procesara su respectivo archivo "ratings.csv" y realizara un inner join con el DataFrame de Movies.
	Contabilizara el número de calificaciones por género.

 La información se sincroniza a través del canal de comunicación "ci".
 Cada worker recibira un apuntadir al ratingSum (ratingSum),
  que tiene la forma [número_filas][número_columnas].
 Los workers escriben en las posiciones [número_filas][j], donde j es el id del worker.
*/

func findRatingsWorker(worker_id int, ci chan int, knowngenres []string, ratingCount *[][]int, ratingSum *[][]float64, movies dataframe.DataFrame) {
	ratingSplit := "ratings_" + fmt.Sprintf("%02d", worker_id) + ".csv"
	println("Worker  ", fmt.Sprintf("%02d", worker_id), " procesará el archivo ", ratingSplit, "\n")
	ratings := csvToDataFrame(&ratingSplit)

	genNumber := len(knowngenres)

	start := time.Now()

	/* Se importan todos los registros del dataframe de Movies en el dataframe de ratings y se realiza un "inner join" mediante
	La operación .Merge
	*/
	ratings.Merge(&movies, "movieId", "genres")

	// Solo se necesitan las columnas "genres" y "ratings" por ello unicamente conservaremos estas del Dataframe
	grcs := [2]string{"genres", "rating"} // grcs => Se infican las columnas "genres" y "Ratings"
	grDF := ratings.KeepColumns(grcs[:])  // grDF => Se indica que grupo de columans del DataFrame se conservaran
	for ig := 0; ig < genNumber; ig++ {
		for _, row := range grDF.FrameRecords {
			if strings.Contains(row.Data[0], knowngenres[ig]) {
				(*ratingCount)[ig][worker_id-1] += 1
				v, _ := strconv.ParseFloat((row.Data[1]), 32)
				(*ratingSum)[ig][worker_id-1] += v
			}
		}
	}

	finish := time.Since(start)
	fmt.Println("Worker ", worker_id, "has finished in", finish)

	// Le decimos al orquestador que hemos terminado
	// Al "prender" el canal de comunicación establecemos indirectamente un protocolo de sincronización
	ci <- 1
}

/*
 * Esta es la función que se ejecutará al correr el programa desde la terminal
 *
 * El trabajo de este programa es:
 * - Abrir un archivo .csv
 * - Dividir el archivo más grande en N archivos pequeños
 * - Saber cómo escribir un archivo .csv
 * - Leer datos desde un archivo .csv. Esto implica parsear el archivo
 * - Procesar los datos leídos desde el .csv
 * - Medir el tiempo que tarda el procesamiento.
 * - Imprimir los resultados
 */
func main() {
	// Estos son los archivos que vámos a leer
	// Hay que considerar que dichos archivos están en el directorio actual
	movies := mvLDir + "movies.csv"
	movieCsv := "movies.csv"
	//ratings := mvLDir + "ratings.csv"

	fmt.Println("Comienza la lectura del archivo 'movies.csv'")

	/* Se realiza la lectura del archivo Movies.csv. */
	movies_data, cerrar_archivo_movies := readFile(movies)
	defer cerrar_archivo_movies() // Cerramos el archivo cuando ya no lo usemos

	/* Se interpreta la estructura de los datos, aqui utilizando la de un CSV
	   esta funcion indica la cantidad de registros del archivo, en este caso Movies.csv */

	csv_movies := readCsvFile(movies_data)
	//defer movies_data.Close()
	fmt.Println("La lectura del archivo 'movies.csv' ha concluido")
	fmt.Println("Primer registro del archivo: \n ", csv_movies[0])

	// Para facilitar el uso de los datos se puede recurrir a los DataFrames
	// Aqui se muestran algunos detalles del contenido del archivo
	df_movies := csvToDataFrame(&movieCsv)
	//df_movies.ViewColumns()
	fmt.Println("Cabecera del archivo: ", df_movies.Headers)

	/* Se cronometra el tiempo de ejecucion, en esta parte se procesaran las N divisiones del archivo Ratings.csv
	   y se realizara el conteo de Ratings por genero rating | genre  a traves del campo llave movieId
	*/
	start := time.Now()

	findRatingsMaster() // Inicia el proceso principal del programa

	finish := time.Since(start)

	fmt.Println("Este programa tardó ", finish)
}
