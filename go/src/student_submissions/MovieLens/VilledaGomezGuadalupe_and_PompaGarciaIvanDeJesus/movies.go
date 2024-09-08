/* *****************************************************
 *
 * Diplomado Inteligencia Artificial y Ciencia de Datos
 * 		23 de agosto de 2024
 *
 *  Code Challenge:
 *		Tomar un archivo CSV y escribir n archivos
 *    de tipo parquet.
 *
 *	Nuestro programa asume que existe el archivo
 *  data/ratings.csv relativo al ejecutable y que tiene
 *  el mismo formato que el encontrado en:
 *  https://grouplens.org/datasets/movielens/25m/
 *
 * 	Compilado:
 *			go build
 *  Ejecución:
 *			./movielens
 *  Ejecución directa:
 *			go run movies.go
 *
 *	Equipo:
 * 		Guadalupe Villeda Gómez
 * 			<lupis_act@ciencias.unam.mx>
 * 		Ivan Pompa-García
 *			<ivanjpg@ekbt.nl>
 *
 * *****************************************************
 */

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

// Definimos la estructura Rating que sirve de base para
// la escritura hacia parquet.
type Rating struct {
	UserID    int32   `parquet:"name=user_id, type=INT32"`
	MovieID   int32   `parquet:"name=movie_id, type=INT32"`
	Rating    float32 `parquet:"name=rating, type=FLOAT"`
	Timestamp int64   `parquet:"name=timestamp, type=INT64"`
}

// La función perror revisa si el
// argumento que recibe es nulo, es decir, si no hay
// error. En caso de que exista, terminamos la ejecución
// de todo el programa.
// El manejo de errores puede ser más fino, pero se
// deja así por simplicidad.
func perror(err error) {
	if err != nil {
		panic(err)
	}
}

// La escritura de los archivos parquet se extrae del
// cuerpo principal del programa con la finalidad de
// utilizar concurrencia para acelerar la escritura
// de los diferentes archivos.
func writeParquetFile(ichannel chan int, i int, nFiles int, rowsRemaining int, rowsPerFile int, rows *[][]string) {
	// Debido al modo en que calculamos el número de
	// registros adecuados para cada archivo parquet (rowsPerFile),
	// en el archivo final debemos tener algunas filas
	// extras (rowsRemaining).

	// Para todos los archivos [0,nFiles-2] no tenemos
	// registros extras.
	var extraFinalRows int = 0

	// Para el archivo final, debemos tener filas extras.
	if i == nFiles-1 {
		extraFinalRows = rowsRemaining
	}

	// Abrimos el archivo donde se escribirá el parquet.
	fileName := fmt.Sprintf("data/ratings-%d.parquet", i)
	pqFile, err := local.NewLocalFileWriter(fileName)
	perror(err)

	// Creamos el escritor de parquet, donde especificamos
	// el tipo de estructura que recibirá (Rating)
	pqWriter, err := writer.NewParquetWriter(pqFile, new(Rating), 2)
	perror(err)

	// Recorremos cada uno de los registros leídos desde
	// el CSV, dependiendo de en qué archivo estemos.
	for j := i * rowsPerFile; j < (i+1)*rowsPerFile+extraFinalRows; j++ {
		row := (*rows)[j]

		// Realizamos las conversiones adecuadas para almacenar
		// la información en la estructura Rating. Nótese
		// que no realizamos manejo de errores por simplicidad.
		userid, _ := strconv.ParseInt(row[0], 10, 32)
		movieid, _ := strconv.ParseInt(row[1], 10, 32)
		rating, _ := strconv.ParseFloat(row[2], 32)
		timestamp, _ := strconv.ParseInt(row[3], 10, 64)

		// Creamos una nueva estructura Rating, con algunos
		// ajustes en las conversiones por los tamaños de las
		// variables.
		ratingObj := Rating{
			UserID:    int32(userid),
			MovieID:   int32(movieid),
			Rating:    float32(rating),
			Timestamp: timestamp,
		}

		// Escribimos el registro dentro del parquet.
		err = pqWriter.Write(ratingObj)
		perror(err)
	}

	// En este punto hemos terminado de procesar los
	// registros correspondientes a uno de los archivos
	// parquet.

	// Indicamos que hemos terminado de escribir en el
	// parquet.
	pqWriter.WriteStop()
	// Cerramos el archivo que contiene el formato parquet.
	pqFile.Close()

	// A través de un canal, enviamos la notificación de
	// que esta gofunction ha concluido.
	ichannel <- 1
}

func main() {
	// Definimos cuántos archivos parquet queremos.
	var nFiles int = 10
	// Cuántos registros tendremos por archivo.
	var rowsPerFile int
	// Los registros "sobrantes" irán al último archivo.
	var rowsRemaining int
	// Contador de control para saber cuántas gofunctions
	// han terminado de escribir su archivo.
	var finishCounter int = 0

	// Canal de comunicación entre las gofunctions
	// y main()
	intChannel := make(chan int)

	fmt.Printf("Leyendo el archivo CSV...")
	// Registramos el tiempo de inicio de lectura del CSV.
	var startTime time.Time = time.Now()
	// Abrimos el CSV.
	csvFile, err := os.Open("data/ratings.csv")
	perror(err)

	// Generamos un nuevo lector de CSV a partir del
	// archivo abierto.
	reader := csv.NewReader(csvFile)

	// Evitamos la verificación del número de columnas
	// de cada registro del CSV. Cada columna puede tener
	// el número de registros que sea.
	reader.FieldsPerRecord = -1

	// Leemos todo el contenido del CSV.
	rows, err := reader.ReadAll()
	perror(err)

	// Ya hemos leído los datos, no es necesario dejar
	// abierto el archivo.
	csvFile.Close()

	// Registramos en pantalla el tiempo en segundos
	// que tardó el programa en leer el CSV.
	fmt.Println("La lectura del CSV tardó", time.Since(startTime).Seconds(), "segundos")
	// Informamos cuántos registros leímos.
	fmt.Println("El archivo tiene", len(rows), "registros.")

	// Calculamos cuántos registros irán en cada archivo
	// parquet, con excepción del último.
	rowsPerFile = int(len(rows) / nFiles)
	// Cuántos registros "extras" tendrá el último archivo.
	rowsRemaining = len(rows) % nFiles

	fmt.Println("Escribiendo", nFiles, "archivos parquet...")
	// Registramos el inicio de la escritura de los parquet.
	startTime = time.Now()
	// Iniciamos una gofunction por cada archivo
	// parquet a generar. Esto puede ser contraproducente
	// y debe moderarse dependiendo del número de núcleos
	// del procesador.
	for i := 0; i < nFiles; i++ {
		go writeParquetFile(intChannel, i, nFiles, rowsRemaining, rowsPerFile, &rows)
	}

	// Cada gofunction regresa un 1 cuando termina de
	// escribir el archivo parquet. Entramos en un ciclo
	// infinito, leyendo el mensaje de término de cada
	// gofunction. Cuando todas terminan, salimos del ciclo.
	for {
		finishCounter += <-intChannel

		if finishCounter == nFiles {
			break
		}
	}

	// Informamos cuánto tardó la escritura de todos
	// los archivos parquet.
	fmt.Println("Parquet files writting took", time.Since(startTime).Seconds(), "seconds")
}
