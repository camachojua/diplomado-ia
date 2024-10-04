package fileprocessing

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// SplitBigFile divide un archivo grande en varios archivos más pequeños.
func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {
	// Abrir el archivo original
	filePath := filepath.Join(directory, file_name)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Crear un escáner para leer el archivo línea por línea
	scanner := bufio.NewScanner(file)

	// Contar cuántas líneas tiene el archivo original
	var totalLines int
	for scanner.Scan() {
		totalLines++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error al leer el archivo: %v", err)
	}

	// Volver al inicio del archivo
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("error al reiniciar el archivo: %v", err)
	}
	scanner = bufio.NewScanner(file)

	// Calcular cuántas líneas debe tener cada archivo
	linesPerChunk := totalLines / number_of_chunks
	remainder := totalLines % number_of_chunks

	// Variables para almacenar los nombres de los archivos generados
	var outputFiles []string

	// Dividir el archivo en múltiples archivos más pequeños
	for i := 0; i < number_of_chunks; i++ {
		outputFileName := fmt.Sprintf("part_%d.csv", i+1)
		outputFilePath := filepath.Join(directory, outputFileName)
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			fmt.Printf("error al crear el archivo de salida: %v", err)
		}
		defer outputFile.Close()

		writer := bufio.NewWriter(outputFile)
		linesToWrite := linesPerChunk
		if i < remainder {
			linesToWrite++ // Distribuir las líneas sobrantes
		}

		// Escribir las líneas correspondientes a este archivo
		for j := 0; j < linesToWrite && scanner.Scan(); j++ {
			_, err := writer.WriteString(scanner.Text() + "\n")
			if err != nil {
				fmt.Printf("error al escribir en el archivo de salida: %v", err)
			}
		}

		// Asegurarse de que todo se haya escrito en el archivo
		writer.Flush()

		// Agregar el archivo generado a la lista de archivos de salida
		outputFiles = append(outputFiles, outputFileName)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error al escanear el archivo: %v", err)
	}

	return outputFiles
}
