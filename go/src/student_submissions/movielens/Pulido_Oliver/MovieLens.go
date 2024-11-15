package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// ----------------------------------------------- Función que cuenta las Filas de un archivo --------------------------------------------------

func NumFilas(arch string) int { //Vamos a recibir el nombre del archivo y devolvemos un numero (la cantidad de filas)
	archivo, err := os.Open(arch) //Abrimos el archivo
	if err != nil {
		log.Fatal("Error al abrir el archivo: ", err)
	} //Manejamos el error
	defer archivo.Close() //Cerramos el archivo

	lec := csv.NewReader(archivo) //Necesitamos crear un lector que lea el archivo
	contador := 0                 //Inicializamos un contador para contar las filas

	for { //En éste ciclo for leemos todo el archivo fila por fila
		_, err := lec.Read()
		if err != nil {
			break
		}
		contador++ //Después de cada fila sumamos un 1 al contador, de esta manera contamos
	} //las filas del archivo

	return contador //Se devuelve la cantidad de filas contadas
}

// ----------------------------------- Función que lee una cantidad arbitraria de líneas de un archivo -----------------------------------------

//Recibimos el nombre del archivo, desde qué fila queremos leer y hasta cuál terminaremos de leer

func LeeCSV(archivo string, inicio, final int) [][]string { //Abrimos el archivo
	arch, err := os.Open(archivo)
	if err != nil {
		log.Fatal(err)
	}
	defer arch.Close() //Manejamos los errores y cerramos el archivo

	lec := csv.NewReader(arch) //Creamos un lector

	for i := 0; i < inicio; i++ { //Éste for es para saltar las filas que estén antes de la fila que deseamos leer
		_, err := lec.Read() //Leemos las líneas que queremos saltarnos
		if err != nil {
			return nil //Al retornar nil, (Nada o vacío) aseguramos que se omitirán las filas que no queremos leer
		}
	}

	//Declaramos un slice 2D para almacenar las filas que sí se van a leer
	lineas := [][]string{}             //Es necesaria ésta estructura 2D para poder guardar su contenido en un CSV
	for i := inicio; i <= final; i++ { //Una vez ignoradas las filas leemos el fragmento deseado
		record, err := lec.Read() //Aquí se leen las filas deseadas
		if err != nil {
			break //Si tenemos un error nos salimos del bucle
		}
		lineas = append(lineas, record) //Agregamos lo que leímos al slice
	}

	return lineas //Devolvemos las líneas leídas en el formato [[a11, a12, ..., a1n], [a21, a22, ..., a2n], ..., [an1, an2, ... , ann]]
	//El cuál es el formato requerido para convertirlo a un archivo CSV
}

// --------------------------------------------------- Función que crea un archivo CSV --------------------------------------------------------

//Recibimos las líneas que deseamos agregar al CSV y el nombre del archivo, devolveremos sólo un error

func CreaArch(lineas [][]string, nombre_archivo string) error {

	archivo, err := os.Create(nombre_archivo) //Abrimos el archivo con el nombre que deseamos darle
	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close() //Manejamos los errores y cerramos el archivo

	escr := csv.NewWriter(archivo) //Necesitamos crear un escritor encargado de redactar el CSV
	defer escr.Flush()             //Con ésto se asegura que se escriba todo y se vacíe el buffer

	for _, line := range lineas { //Aquí revisamos cada registro ignorando su índice
		if err := escr.Write(line); err != nil { //El método Write devuelve un error si hay error y nil si no lo hay
			return err //Es por ésto que asignamos err al método, si hay un error al escribir el archivo se mostrará
		}
	}
	return nil //Retornamos nada porque el archivo se genera por sí solo en caso de no haber ningún error
}

// ---------------------------------- Aquí generamos la función de Partición ----------------------------------------------------------------------

// Recibimos el nombre del archivo que se va a particionar y el número de archivos en que lo queremos partir

func Particion(csv string, num_archivos int) {
	archivo := csv
	part := num_archivos

	n := NumFilas(archivo) //Obtenemos el número de filas del archivo a particionar
	div := n / part        //Dividimos el número de filas entre la cantidad de archivos que queremos crear
	//Ésto nos da la cantidad de filas que contendrá cada nuevo archivo

	var inicio int
	var fin int

	for i := 0; i < part; i++ { //En éste for creamos la cantidad de archivos requeridos
		inicio = i * div      //Comenzamos desde el número de filas múltiplo de la división
		fin = ((i + 1) * div) //Terminamos una línea antes de dos veces un múltiplo de la división
		//Ésto asegura que comencemos en la línea siguiente de la última que escribimos en el archivo anterior

		if i == part-1 { //Aquí aseguramos que las filas restantes se agreguen al último archivo
			inicio = i * div //Comenzamos donde debemos comenzar
			fin = n          //Terminamos donde debe ser el final del archivo
		}
		//Para escribir el nuevo archivo debemos leer las filas necesarias y luego escribirlas en dicho archivo
		lines := LeeCSV("ratings.csv", inicio, fin) //Llamamos a la fución que lee para obtener las líneas del archivo original

		nombre := "ratings_" + fmt.Sprintf("%02d", i+1) + ".csv" //El nombre que tendrá cada archivo

		if err := CreaArch(lines, nombre); err != nil { //Aquí creamos el archivo csv, recordemos que le asignamos la variable err
			log.Fatal(err) //Porque la función usa el método que devuelve un error
		}
	}

}

// -------------------------------------------- Creación de los tipos sugeridos --------------------------------------------------------------------

// Vamos a crear las estructuras sugeridas para usar las funciones Worker y Master proporcionadas por el doctor
type MovieObj struct {
	MovieId int64
	Title   string
	Genres  string
}

// Estructura para representar una calificación
type RatingObj struct {
	UserId    int64
	MovieId   int64
	Rating    float64
	Timestamp int64
}

// -------------------------------------------- Función ReadMoviesCsvFile ---------------------------------------------------------------------

// Creamos la función para leer el archivo de ratings CSV que se necesita para echar a andar la función worker del doctor
// Recibimos el nombre del archivo a leer y devolvemos 2 cosas, un error, si lo hay, y el objeto que contendrá las columnas que
// vienen dentro de los archivos de ratings
func ReadRatingsCsvFile(archivo string) ([]RatingObj, error) {
	var ratings []RatingObj //Creamos una variable tipo Rating

	file, err := os.Open(archivo) //Abrimos el archivo que vamos a leer
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer file.Close() //Manejamos los errores y cerramos el archivo

	reader := csv.NewReader(file) //Creamos un lector

	for {
		record, err := reader.Read() //Tomamos las filas leídas
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error al leer el archivo: %w", err)
		}

		//Vamos a parsear cada campo de cada record y con cada uno, creamos objeto RatingObj
		userId, _ := strconv.ParseInt(record[0], 10, 64)
		movieId, _ := strconv.ParseInt(record[1], 10, 64)
		rating, _ := strconv.ParseFloat(record[2], 64)
		timestamp, _ := strconv.ParseInt(record[3], 10, 64)

		ratings = append(ratings, RatingObj{
			UserId:    userId,
			MovieId:   movieId,
			Rating:    rating,
			Timestamp: timestamp,
		})
	}

	return ratings, nil
}

// -------------------------------------------- Función ReadMoviesCsvFile ---------------------------------------------------------------------

// Creamos la función para leer el archivo de movies CSV que se necesita para echar a andar la función master del doctor
// Recibimos el nombre del archivo a leer y devolvemos el objeto que contendrá las columnas que vienen dentro del archivo de movies.csv

func ReadMoviesCsvFile(filename string) []MovieObj {
	var movies []MovieObj //Creamos un objeto tipo Movie

	file, err := os.Open(filename) //Abrimos el archivo
	if err != nil {
		log.Fatalf("No se pudo abrir el archivo: %s", err)
	}
	defer file.Close() //Manejamos los errores y cerramos el archivo

	reader := csv.NewReader(file) //Creamos un lector
	reader.Read()                 // Con ésto, estamos leyendo la primer fila y luego no hacemos nada, ésto garantiza que ignoramos el encabezado

	for {
		record, err := reader.Read() //Con ésto leemos el archivo fila por fila
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error al leer el archivo: %s", err)
		}

		//Sólo es necesario parsear el primer campo de cada record y con todos un objeto MovieObj
		movieId, _ := strconv.ParseInt(record[0], 10, 64)
		title := record[1]
		genres := record[2]

		movies = append(movies, MovieObj{
			MovieId: movieId,
			Title:   title,
			Genres:  genres,
		})
	}

	return movies
}

// -------------------------------------------- Función Worker ----------------------------------------------------------------------------
func Mt_FindRatingsWorker(w int, ci chan int, kg []string, ca *[][]int, va *[][]float64, movies []MovieObj) {
	aFileName := "ratings_" + fmt.Sprintf("%02d", w) + ".csv"
	println("El Worker ", fmt.Sprintf("%02d", w), " está procesando el archivo ", aFileName, "")

	ratings, err := ReadRatingsCsvFile(aFileName)
	if err != nil {
		log.Printf("No se pudo abrir el archivo %s: %v\n", aFileName, err)
		ci <- 1 // Notificar que el worker terminó, aunque sin procesar nada
		return
	}

	//Vamos a crear un mapa para buscar los géneros por movieId rápidamente para no usar el método Merge
	movieGenres := make(map[int64]string)
	for _, movie := range movies {
		movieGenres[movie.MovieId] = movie.Genres
	}

	ng := len(kg)
	start := time.Now()

	//Procesamos los ratings
	for ig := 0; ig < ng; ig++ {
		for _, rating := range ratings {
			//Buscamoslos géneros del movieId actual en el mapa
			genres, exists := movieGenres[rating.MovieId]
			if !exists {
				continue //Si no se encontramos el movieId en el mapa de películas nos saltamos a la fila que sigue
			}

			//Verificamos si el género coincide con la lista conocida
			if strings.Contains(genres, kg[ig]) {
				(*ca)[ig][w-1] += 1
				(*va)[ig][w-1] += rating.Rating
			}
		}
	}

	duration := time.Since(start)

	fmt.Println("\nWorker ", w, " completado")
	fmt.Println("Duración = ", duration)

	// Notificamos al master que este worker ha completado su trabajo
	ci <- 1
}

// ---------------------------------------------------- Función Master ------------------------------------------------------------
func Mt_FindRatingsMaster() {
	fmt.Println("En MtFindRatingsMaster:")

	start := time.Now()
	nf := 10 // Cantidad de archivos de calificaciones y también de subprocesos

	// Lista de géneros conocidos
	kg := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}
	ng := len(kg) // Número de géneros conocidos

	// Matrices 2D para acumular valores de calificaciones y conteo por género y worker
	// ra es una matriz 2D donde se mantienen los valores de calificación para cada género.
	// Las columnas indican/mantienen el número de núcleo en el que se está ejecutando un worker.
	// Las filas de esa columna mantienen los valores de calificación para ese núcleo y ese género

	// ca es una matriz 2D donde se mantiene el recuento de calificaciones para cada género
	// Las columnas indican el número de núcleo donde se está ejecutando el trabajador
	// Las filas en esa columna mantienen los recuentos para ese género

	ra := make([][]float64, ng)
	ca := make([][]int, ng)
	for i := 0; i < ng; i++ {
		ra[i] = make([]float64, nf)
		ca[i] = make([]int, nf)
	}

	// Canal para sincronizar a todos los workers
	ci := make(chan int)
	movies := ReadMoviesCsvFile("movies.csv")

	// Ejecutar FindRatings en 10 workers
	for i := 0; i < nf; i++ {
		go Mt_FindRatingsWorker(i+1, ci, kg, &ca, &ra, movies)
	}

	// Esperar a los workers
	iMsg := 0
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
	// Todos los workers completaron su trabajo. Recopilar los resultados.
	locCount := make([]int, ng)
	locVals := make([]float64, ng)
	for i := 0; i < ng; i++ {
		for j := 0; j < nf; j++ {
			locCount[i] += ca[i][j]
			locVals[i] += ra[i][j]
		}
	}

	// Imprimir resultados: conteo, suma y promedio de calificaciones por género
	fmt.Println("\n ---------------- Resultados: -----------------")
	fmt.Printf("\n%-2s %20s %10s %15s\n", "ID", "Género", "Conteo Raitings", "Promedio Raitings")
	for i := 0; i < ng; i++ {
		var promedio_Rating float64
		if locCount[i] > 0 {
			promedio_Rating = locVals[i] / float64(locCount[i]) // Aquí calculamos el promedio solo si hay calificaciones
		}
		fmt.Printf("%-2d %20s %10d %15.5f\n", i, kg[i], locCount[i], promedio_Rating)
	}

	duration := time.Since(start)

	println("\nMt_FindRatingsMaster completado")
	fmt.Println("Duración = ", duration)

}

// ------------------------------------------------ Función Main ---------------------------------------------------------------------------
func main() {
	//Particion("ratings.csv", 10)  //Descomentar para generar la partición
	Mt_FindRatingsMaster()

}
