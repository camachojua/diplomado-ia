package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func SplitBigFile(fileName string, numberOfChunks int, directory string) []string {
	// Intentar abrir el archivo CSV
	fileName = filepath.Join(directory, fileName)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("cannot find the path to read the files, sorry :(")
		fmt.Println(fileName)
		return nil
	}
	defer file.Close()

	// Intentar leer el archivo CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil
	}

	// Calcular el tamaño de cada partición
	totalRecords := len(records)
	partitionSize := totalRecords / numberOfChunks
	if totalRecords%numberOfChunks != 0 {
		partitionSize++
	}

	// Slice para almacenar los nombres de los archivos generados
	var fileNames []string

	// Particionar y escribir los nuevos archivos CSV
	for i := 0; i < numberOfChunks; i++ {
		start := i * partitionSize
		end := start + partitionSize
		if end > totalRecords {
			end = totalRecords
		}

		// Crear un nuevo archivo CSV para la partición
		outputFileName := filepath.Join(directory, "output_"+strconv.Itoa(i+1)+".csv")
		justName := "output_" + strconv.Itoa(i+1) + ".csv"
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Println("cannot find the path to save the files, sorry :(")
			return nil
		}
		defer outputFile.Close()

		// Escribir la partición en el archivo CSV
		writer := csv.NewWriter(outputFile)
		err = writer.WriteAll(records[start:end])
		if err != nil {
			fmt.Println("cannot find  :(")
			return nil
		}
		writer.Flush()

		// Agregar el nombre del archivo al slice
		fileNames = append(fileNames, justName)

		fmt.Println("Archivo", outputFileName, "creado con éxito.")
	}

	fmt.Println("the length of the array is: ", len(fileNames))
	fmt.Println(fileNames)

	return fileNames
}
