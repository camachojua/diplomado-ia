package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kfultz07/go-dataframe"
)

func main() {

	fila_peli := "movies_large.csv"
	procesaArchivoMultiHilo(fila_peli)

}

// Esta función lee un archivo y carga los datos de dicho archivo en memoria.
// Para compartir los datos que están en memoria la función regresa un
// "apuntador" a la memoria ocupada por el archivo leído.  Pueden ver el
// apuntador como una flecha que le dice a Go "aquí están los datos".  Dicho
// apuntador es del tipo de dato llamado "os.File", el "*" indica que es
// apuntador.

// Por cuestiones técnicas también debemos de regresar una función "func()" ¿Y
// eso de qué me sirve? En Go cuando abrimos un recurso hay que "cerrarlo"
// cuando ya no lo utilizamos, en particular los archivos se abren, se leen o
// escriben y por último se cierran esa función nos ayudará a cerrar el archivo

// Tarea moral: ¿Qué tipos de validaciones debería de tener esta función?
func leeArchivo(archivo string) (*os.File, func()) {
	datos, err := os.Open(archivo) // Leemos el archivo
	if err != nil {
		panic("Error al abrir el archivo.")
	}

	// Esto es una variable que representa una función vean que del lado
	// derecho diche "func"
	cerrar_archivo := func() {
		err := datos.Close()

		if err != nil {
			panic("No pudimos cerrar el archivo :(")
		}
	}

	return datos, cerrar_archivo
}

// Esta función lee los datos binarios de un archivo, los interpreta como un csv
// Y regresa una matriz que representa al archivo CSV
func leeArchivoCsv(archivo *os.File) [][]string {
	// Acá tenemos los datos "crudos" el archivo es un csv por lo que
	// debemos decirle a Go que los lea Para ello utilizamos csv.NewReader
	parser := csv.NewReader(archivo)
	parser.Read()

	// Leemos todos los registros del archivo, como todo en la vida puede
	// haber un error entonces revisamos si hubo un error antes de regresar
	// la información
	records, err := parser.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Printf("El archivo CSV tiene %d registros\n", len(records))

	return records
}

// Leemos el archivo CSV para convertirlo en un DataFrame
func csvToDataFrame(archivo string) dataframe.DataFrame {
	// La firma de la función necesita una ruta de directorio usamos el
	// alias "./" para referirnos al directorio actual
	df := dataframe.CreateDataFrame("./", archivo)
	return df
}

// Esta función orquesta la tarea de procesar de manera concurrente la tarea de
// realizar un inner join entre el archivo "movies.csv" y el archivo
// "ratings.csv" con el find de encontrar el número de calificaciones por archivo.
//
// Para realizar esto necesitamos varias cosas:
// - Crear canales de comunicación entre el orquestador (esta función) y los workers
// - Esperar a que todos los workers terminen de trabajar
// - Leer la respuesta de cada uno de los workers
// - Consolidar la respuesta en un único resultado

// - El orquestador crea el canal de comunicación y crea los workers utilizando go-rutinas
// - Los workers sincronizan su operación a través del canal de sincronización, para evitar conflictos utilizamos el "id" del worker como identificador de canal
// - Los workers reciven un apuntador al arreglo de resultados, el arreglo de resultados tiene la forma [numero_filas][numero_columnas]
// - Cada worker puede escribir información en las entradas tipo [numero_filas][j]
// - En este caso el código de sincronización es nulo (porque no tenemos condiciones de carrera)
func procesaArchivoMultiHilo(archivo string) {
	fmt.Println("El orquestador del proceso ha iniciado su ejecución.")

	// Definimos el número de archivos a generar, también lo usamos como el
	// número de hilos que usaremos
	nivel_multiprogramacion := 10
	// numero_generos := 20

	generos_conocidos := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	numero_generos := len(generos_conocidos)

	// Matriz de 2 dimensiones que guardará los resultados (género y número
	// de calificaciones) Las columnas hacen referencia al "id" del worker
	// que está trabajando su parte del progrema Las filas contienen las
	// calificaciones entcontradas por un worker especifico dado un género
	arreglo_resultados := make([][]float64, numero_generos)

	// Esta matriz de 2 dimensiones almacena la cuenta total que llevamos
	// hasta el momento Las columnas hacen referencia al "id" del worker que
	// está trabajando su parte del problema Las filas mantienen el conteo
	// de las caliificaciones encontradas por un worker específico dado un
	// género
	arreglo_conteo := make([][]int, numero_generos)
	// Llenamos las "numero_generos" filas de los arreglos
	// "arreglo_resultados" y "arreglo_conteo", les ponemos
	// "nivel_multiprogramacion" columnas
	for i := 0; i < numero_generos; i++ {
		arreglo_resultados[i] = make([]float64, nivel_multiprogramacion)
		arreglo_conteo[i] = make([]int, nivel_multiprogramacion)
	}
	// Creamos el canal de comunicación
	var ci = make(chan int)

	movies := csvToDataFrame("movies_large.csv") // Sacamos los dataframes del archivo csv

	// Creamos 10 workers, cada worker se encarga de leer su archivo correspondiente
	for i := 0; i < nivel_multiprogramacion; i++ {
		go encuentraCalificaciones(i+1, ci, generos_conocidos, &arreglo_conteo, &arreglo_resultados, movies)
	}

	// Esperamos que los workers terminen de trabajar
	iMsg := 0
	// Consultamos el valor del canal de comunicación cuando un worker
	// termina de trabajar grita ¡Ya acabé!, ese grito puede ser escuchado a
	// través del canal de comunicación
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

	// Acá "consolidamos la información"
	locCount := make([]int, numero_generos)
	locVals := make([]float64, numero_generos)
	locAVG := make([]float64, numero_generos)
	for i := 0; i < numero_generos; i++ {
		for j := 0; j < nivel_multiprogramacion; j++ {
			locCount[i] += arreglo_conteo[i][j]
			locVals[i] += arreglo_resultados[i][j]
		}
		locAVG[i] = float64(locVals[i]) / float64(locCount[i])
	}

	// Acá imprimimos los resultados
	for i := 0; i < numero_generos; i++ {

		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", generos_conocidos[i]), "  ", fmt.Sprintf("%8d", locCount[i]), "  ", fmt.Sprintf("%.2f", locAVG[i]), "  ", fmt.Sprintf("%.2f", locVals[i]))
	}

	println("Fin del orquestador.")
}

// Este es el trabajo que realizará cada worker.
// La forma de sincronizar la información es a través del canal de comunicación (llamado "ci").
// Cada worker recibe un apuntador al arreglo de resultados (llamado arreglo_resultados)
// dicho arreglo tiene la forma [numero_filas][numero_columnas].
// Cada worker puede escribir información en [numero_filas][j], j es la j-ésima columna del arreglo de resultados.
func encuentraCalificaciones(worker_id int, ci chan int, generos_conocidos []string, arreglo_conteo *[][]int, arreglo_valor *[][]float64, movies dataframe.DataFrame) {
	// Acá hago trampa, ya que supongo que en el directorio actual se encuentra el archivo "ratings.csv" partido
	// Tu tarea es crear la función que parte dicho archivo
	// ratings_1.csv, ..., raitngs_10.csv <= Ustedes deben generar esto con código
	ratings_chiquito := "ratings_parte_" + fmt.Sprintf("%02d", worker_id) + ".csv"
	println("Worker  ", fmt.Sprintf("%02d", worker_id), " procesará el archivo ", ratings_chiquito, "\n")
	ratings := csvToDataFrame(ratings_chiquito)

	numero_generos := len(generos_conocidos)

	tiempo_inicial := time.Now()

	// Importamos todos los registros del dataframe de películas en el dataframe de ratings
	// La operación data_frame.Merge es equivalente a un inner-join en esta biblioteca específica
	ratings.Merge(&movies, "movieId", "genres")

	// De todo el dataset sólo nos interesan los generos y los ratings,
	// entonces mantenemos únicamente esas columnas del dataframe

	columnas_que_nos_interesan := [2]string{"genres", "rating"}
	genero_df := ratings.KeepColumns(columnas_que_nos_interesan[:])

	// ====================================================================
	// Acá tenemos que pensar cómo hacer el conteo y cómo reportar el resultado

	movies.CountRecords()

	for i := 0; i < numero_generos; i++ {
		for _, row := range genero_df.FrameRecords {
			if strings.Contains(row.Data[0], generos_conocidos[i]) {
				(*arreglo_conteo)[i][worker_id-1] += 1
				v, _ := strconv.ParseFloat((row.Data[1]), 32) // do not check for error
				(*arreglo_valor)[i][worker_id-1] += v
			}

		}
	}

	// ====================================================================

	tiempo_final := time.Since(tiempo_inicial)
	fmt.Println("Tiempo en procesar = ", tiempo_final)
	fmt.Println("Worker ", worker_id, " ha terminado 🟢")

	// Le decimos al orquestador que hemos terminado
	// Al "prender" el canal de comunicación establecemos indirectamente un protocolo de sincronización
	ci <- 1
}
