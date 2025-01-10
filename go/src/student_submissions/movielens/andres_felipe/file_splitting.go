package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"math"
	"time"
)

//Función para particionar un archivo de datos grande, 
//donde file_name es el nombre del CSV grande; partitions es el número de minis CSV en el que se dividirán 
//los datos; directory es el directorio donde se encuentra el CSV 
func SplitBigFile(file_name string, partitions int, directory string) []string {
	//Se abre el archivo
	filePath := fmt.Sprintf("%s/%s", directory, file_name)
	file, err := os.Open(filePath)
	if err != nil {
    	//Muestra el error y termina el programa si no se puede leer el archivo CSV
    	log.Fatal("No se pudo abrir el archivo CSV: ", err)
	}
	defer file.Close()

	//Se lee el archivo
	CSV_reader, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal("No se pudo leer el archivo CSV: ", err)
	}

	//Número de registos/lineas que contiene el archivo
	n_lines := len(CSV_reader)
	fmt.Printf("El archivo CSV consta de %v registros \n", n_lines)

	//Número de registros por cada partición
	np_lines := math.Floor(float64(n_lines)/float64(partitions))
	fmt.Printf("Por cada partición habrá %v registros \n", np_lines)

	//Concurrencia
	ch := make(chan string)
	for n := 0; n < partitions; n++ {
		go splits(directory, n, partitions, np_lines,  CSV_reader, ch)
	}

	var fileNames []string //Este slice contendrá el nombre de los minis CSV generados
	var fileName string
	for ch_done := 0; ch_done < partitions; ch_done ++ {
		fileName = <- ch
		fileNames = append(fileNames, fileName)
	} 
	return fileNames
}

//Función de escritura de cada partición
func splits(
	directory string,
	n_partition int,
	partitions int,
	np_lines float64,
	CSV_reader [][]string,
	ch chan string,
) {
	//Creación del mini CSV de la n-partición
	mini_filePath := fmt.Sprintf("%s/partition_%v.csv", directory, n_partition)
	mini_fileName := fmt.Sprintf("partition_%v.csv", n_partition)
	mini_file, err := os.Create(mini_filePath)
	if err != nil {
    	log.Fatal("No se pudo crear la partición: ", err)
	}
	defer mini_file.Close()

	// Escribir los datos en la partición/archivo CSV
	start := n_partition * int(np_lines)
	end := start + int(np_lines)
	if n_partition == partitions - 1 {
		end = len(CSV_reader)
	}
	mCSV_writer := csv.NewWriter(mini_file)
	err = mCSV_writer.WriteAll(CSV_reader[start:end])
	if err != nil {
		log.Fatal("No se puede escribir en CSV: ", err)
	}
	mCSV_writer.Flush()

	//Manda el nombre del archivo CSV generado para después agregarlo al slice
	ch <- mini_fileName
}

func main() {
	var startTime time.Time = time.Now()
	slice_names := SplitBigFile("movies.csv", 10, "/home/ws117/diplomado-ia-1/go/MovieLens/movies")
	fmt.Println(slice_names)
	fmt.Printf("La partición tomó %v segundos", time.Since(startTime).Seconds())
}