package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

func Split_csv() {
	// Abrir el archivo CSV
	largeCSVFile, err := os.Open("ratings.csv")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer largeCSVFile.Close()

	reader := csv.NewReader(largeCSVFile)

	// Crear canales y espera para sincronización
	const numFiles = 10
	recordChans := make([]chan []string, numFiles)
	var wg sync.WaitGroup

	// Inicializar gorutinas
	for i := 0; i < numFiles; i++ {
		recordChans[i] = make(chan []string, 1000)
		wg.Add(1)
		go func(i int, recordChan chan []string) {
			defer wg.Done()
			// Crear el archivo CSV
			fileName := fmt.Sprintf("ratings_part_%d.csv", i+1)
			smallCSVFile, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error al crear el archivo:", err)
				return
			}
			defer smallCSVFile.Close()

			writer := csv.NewWriter(smallCSVFile)

			// Escribir registros recibidos en el canal
			for record := range recordChan {
				if err := writer.Write(record); err != nil {
					fmt.Println("Error al escribir en el archivo:", err)
					return
				}
			}

			// Asegurarse de que todos los datos se escriban en el archivo
			writer.Flush()
			if err := writer.Error(); err != nil {
				fmt.Println("Error al finalizar la escritura:", err)
				return
			}

			fmt.Printf("Archivo %s creado con éxito.\n", fileName)
		}(i, recordChans[i])
	}

	// Leer el archivo CSV línea por línea y distribuir los registros
	recordIndex := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error al leer el archivo CSV:", err)
			break
		}
		// Enviar el registro al canal correspondiente
		recordChans[recordIndex%numFiles] <- record
		recordIndex++
	}

	// Cerrar todos los canales
	for i := 0; i < numFiles; i++ {
		close(recordChans[i])
	}

	// Esperar a que todas las gorutinas terminen
	wg.Wait()
}
