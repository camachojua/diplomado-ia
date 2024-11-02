package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

/* Este codigo toma el archivo rating.csv y lo divide en N partes
deonde N es un numero que suministra el usuario en el Joiner.go
no se usan dataframes ya que no pude lograr que corriera la libreria
de df en mi computadora, en vez de eso se manekan los datos mediante
estrcuturas*/

// Estructura Rating, mediante la cual se guardan los datos del csv
type Rating struct {
	UserID    int64
	MovieID   int64
	Rating    float64
	Timestamp int64
}

/*
	Función writeCSVPart que escribe cada chunk en un archivo CSV

independiente. Recibe el nombre del archivo y el conjunto de datos
como un slice de slices de strings.
*/
func writeCSVPart(fileName string, data [][]string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

/*
Función Splitter: divide el archivo ratings.csv en N partes
sin cargar todo el archivo en memoria, en su lugar lo lee línea
por línea para mayor eficiencia.
*/
func Splitter(csvFile string, numParts int) error {
	file, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return err
	}

	// Leer todas las líneas y cuenta el total
	var records [][]string
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		records = append(records, record)
	}

	// Definimos las variables de conteo
	totalLines := len(records)
	linesPerPart := totalLines / numParts
	residuo := totalLines % numParts

	// Inicializamos un WaitGroup para manejar las gorutinas
	var wg sync.WaitGroup

	// Dividimos las líneas en partes para procesarlas con gorutinas
	startIndex := 0
	for part := 1; part <= numParts; part++ {
		currentPartLines := linesPerPart
		// Distribuimos de manera equitativa los datos en los archivos
		if part <= residuo {
			currentPartLines++
		}

		// Definimos el rango de líneas para la parte actual
		endIndex := startIndex + currentPartLines
		partRecords := append([][]string{header}, records[startIndex:endIndex]...)
		startIndex = endIndex

		// Llamada a la gorutina para escribir el archivo x parte
		wg.Add(1)
		go func(part int, partRecords [][]string) {
			defer wg.Done()
			fileName := fmt.Sprintf("ratings_%d.csv", part)
			if err := writeCSVPart(fileName, partRecords); err != nil {
				log.Printf("error al escribir %s: %v", fileName, err)
			}
		}(part, partRecords)
	}

	// Esperar a que todas las gorutinas terminen
	wg.Wait()
	return nil

}

// Aqui solo llamamos a la funcion Splitter y la usamos para llamarla en el codigo de Joiner
func Splitter_main(num_parts int) (tiempo float64) {
	comienzo := time.Now()
	if err := Splitter("ratings.csv", num_parts); err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		fmt.Println("Archivo CSV partido correctamente.")
	}

	return float64(time.Since(comienzo))
}
