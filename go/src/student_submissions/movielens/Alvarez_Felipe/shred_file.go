//Fragmenta el csv y lo convierte en parquets

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

//El archivo ratings.csv cuenta con 4 columnas [userId movieId rating timestamp]
//que contienen datos de tipo [int, int, float, int]

type Rating_Obj struct {
	UserId    int64   `parquet:"name=userId, type=INT64"`
	MovieId   int64   `parquet:"name=movieId, type=INT64"`
	Rating    float64 `parquet:"name=rating, type=DOUBLE"`
	Timestamp int64   `parquet:"name=timestamp, type=INT64"`
}

//Contamos el número de renglones que tiene el csv y lo repartimos entre el
//número de rutinas que correran simultaneamente

func getFileLength(sPath string) (int, error) {
	//Obtiene el número de registros del archivo
	file, err := os.Open(sPath) //Abre el archivo o regresa un error

	if err != nil { //En caso de error, te lo hace saber
		fmt.Println("Error al leer el archivo")
		return 0, err
	}

	defer file.Close() //Cerramos el archivo

	scanner := bufio.NewScanner(file) //Crea un scaner que lee el contenido del archivo linea por linea
	lineCount := 0                    //Iniciamos el conteo de lineas
	for scanner.Scan() {              //Recorremos el archivo por cada línea
		lineCount++ //Aumentamos la cuenta de líneas
	}

	return lineCount, scanner.Err() //Retornamos el conteo de líneas totales y si hubo algún error
}

func getSubFilesInfo(iFileRows int) (int, []int64) {
	//Dado el número de filas del archivo, encuentra el número de cpu disponibles y
	//regresa el número óptimo de archivos a dividir y un slice conteniendo los archivos
	iNumCPU := runtime.NumCPU()

	fmt.Println("Tu archivo se fragmentará en", iNumCPU, "archivos")

	iSubFileRows := iFileRows / iNumCPU //Calculamos el número de filas por archivo

	switch {
	case (iSubFileRows * iNumCPU) != iFileRows: //En caso de no ser divisible exactamente en archivos de igual tamaño

		fmt.Println("No se puede dividir exactamente en", iNumCPU, "archivos de", iSubFileRows, "filas")
		fmt.Println("Se requiere que el último archivo contenga", iSubFileRows+(iFileRows-(iSubFileRows*iNumCPU)), "registros")
		sliceSubFileSize := make([]int64, iNumCPU) //Declaramos un slice con tantas entradas como numero de archivos

		for i := 0; i <= (iNumCPU - 1); i++ { //Llenamos el array con las lineas que debe contener cada archivo
			sliceSubFileSize[i] = int64(iSubFileRows)
		}

		sliceSubFileSize[iNumCPU-1] = int64(iSubFileRows + (iFileRows - (iSubFileRows * iNumCPU))) //Le damos al último archivo la longitud diferente

		return iNumCPU, sliceSubFileSize

	case ((iFileRows / iNumCPU) * iNumCPU) == iFileRows:
		fmt.Println("Se puede dividir exactamente en", iFileRows/iNumCPU, "archivos")
		sliceSubFileSize := make([]int64, iNumCPU) //Declaramos un slice con tantas entradas como numero de archivos

		for i := 0; i < iNumCPU; i++ { //Llenamos el array con las lineas que debe contener cada archivo
			sliceSubFileSize[i] = int64(iSubFileRows)
		}

		return iNumCPU, sliceSubFileSize
	}

	return iNumCPU, nil
}

func getFileContenToObj(sPath string, iStartRow int64, iEndRow int64) ([]Rating_Obj, error) {

	arrContent := make([]Rating_Obj, 0) //Creamos un array vacío

	file, err := os.Open(sPath) //Abre el archivo o regresa un error

	if err != nil { //En caso de error, te lo hace saber
		fmt.Println("Error al leer el archivo")
		return arrContent, err
	}

	defer file.Close() //Cerramos el archivo

	scanner := bufio.NewScanner(file) //Crea un scaner que lee el contenido del archivo linea por linea

	var iCurrentLine int64 = -1

	for scanner.Scan() { //Recorremos el archivo por cada línea

		iCurrentLine++ //Pasamos a lla siguiente línea

		if ((iCurrentLine >= iStartRow) && (iCurrentLine <= iEndRow)) || (iCurrentLine == 0) { //Mientras este en el rango que nos interesa ó si es el primero

			sLine := scanner.Text() //Tomamos la info de la línea

			sFields := strings.Split(sLine, ",") //La separamos

			//Guardamos todos los valores con su correcto tipo de dato
			iUserId, _ := strconv.ParseInt(sFields[0], 10, 64)
			iMovieId, _ := strconv.ParseInt(sFields[1], 10, 64)
			fRating, _ := strconv.ParseFloat(sFields[2], 64)
			iTimestamp, _ := strconv.ParseInt(sFields[3], 10, 64)

			//Creamos un nuevo objeto para el llenado de los datos
			ratobjNewRow := Rating_Obj{UserId: iUserId, MovieId: iMovieId, Rating: fRating, Timestamp: iTimestamp}

			//Guardamos el contenido
			arrContent = append(arrContent, ratobjNewRow)

		}
	}

	return arrContent, scanner.Err()
}

func getParquetFromObj(sParquetPath string, arrContent []Rating_Obj) {
	//Funcion que crea un arhivo parquet a partir de
	//uno de los objetos creados previamente
	// Crear el archivo Parquet
	parquetFile, err := local.NewLocalFileWriter(sParquetPath)

	if err != nil {
		print(err)
		log.Fatalf("No se pudo crear el archivo Parquet: %s", err)
	}
	defer parquetFile.Close()

	// Escribir los datos en el archivo Parquet
	pw, err := writer.NewParquetWriter(parquetFile, new(Rating_Obj), 4)
	if err != nil {
		log.Fatalf("No se pudo crear el escritor de Parquet: %s", err)
	}
	defer pw.WriteStop()

	for _, rating := range arrContent {

		if err := pw.Write(rating); err != nil {
			log.Fatalf("No se pudo escribir en el archivo Parquet: %s", err)
		}
	}

	fmt.Println("Archivo Parquet creado exitosamente. ", sParquetPath)
}

func getFirstnLastIndex(sliceSubFileSize []int64) [][]int64 {

	var arrIndexArr [][]int64

	start := int64(1)
	end := int64(sliceSubFileSize[0])

	for i := 0; i < len(sliceSubFileSize); i++ {

		if i > 0 {
			start = end + int64(1)
			end = start + sliceSubFileSize[i] - int64(1)
		}

		// Crear un nuevo array con el rango y añadirlo a la lista de resultados
		rangeArray := []int64{start, end}
		arrIndexArr = append(arrIndexArr, rangeArray)

	}

	return arrIndexArr
}

type MovieObj struct {
	MovieId int64  `parquet:"name=movieId, type=int64"`
	Title   string `parquet:"name=title, type=STRING"`
	Genres  string `parquet:"name=genres, type=UTF8"`
}

func csvToParquet(sPath, sParquetPath string) {
	//Funcion que convierte un archivo csv en parquet

	arrContent := make([]MovieObj, 0) //Creamos un array vacío

	file, err := os.Open(sPath) //Abre el archivo o regresa un error

	if err != nil { //En caso de error, te lo hace saber
		fmt.Println("Error al leer el archivo")
	}

	defer file.Close() //Cerramos el archivo

	scanner := bufio.NewScanner(file) //Crea un scaner que lee el contenido del archivo linea por linea

	for scanner.Scan() { //Recorremos el archivo por cada línea

		sLine := scanner.Text() //Tomamos la info de la línea

		sFields := strings.Split(sLine, ",") //La separamos

		//Guardamos todos los valores con su correcto tipo de dato
		iMovieId, _ := strconv.ParseInt(sFields[1], 10, 64)
		sTitle := sFields[2]
		sGenre := sFields[3]

		//Creamos un nuevo objeto para el llenado de los datos
		movobjNewRow := MovieObj{MovieId: iMovieId, Title: sTitle, Genres: sGenre}

		//Guardamos el contenido
		arrContent = append(arrContent, movobjNewRow)

	}

	parquetFile, err := local.NewLocalFileWriter(sParquetPath)

	if err != nil {
		print(err)
		log.Fatalf("No se pudo crear el archivo Parquet: %s", err)
	}

	defer parquetFile.Close()

	// Escribir los datos en el archivo Parquet
	pw, err := writer.NewParquetWriter(parquetFile, new(MovieObj), 4)
	if err != nil {
		log.Fatalf("No se pudo crear el escritor de Parquet: %s", err)
	}
	defer pw.WriteStop()

	for _, rating := range arrContent {

		if err := pw.Write(rating); err != nil {
			log.Fatalf("No se pudo escribir en el archivo Parquet: %s", err)
		}
	}

	fmt.Println("Archivo Parquet creado exitosamente. ", sParquetPath)

}
