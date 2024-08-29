package fileprocessing
//Autor: Luis Eduardo Sanchez Perez

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"path/filepath"
)

func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {
	
	// Se abre el archivo CSV original
	ruta := filepath.Join(directory, file_name)
	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Println("No se puede abrir el archivo: %w", err)
		return []string{}

	}
	defer archivo.Close()

	// Se crea un lector del archivo CSV
	reader := csv.NewReader(archivo)

	// Se leen todas las líneas del archivo CSV
	filas, err := reader.ReadAll()
	if err != nil {
		fmt.Println("No se puede leer el archivo: %w", err)
		return []string{}
	}

	// Se debe calcular cuántas filas habrá en cada archivo partido
	numElementos := len(filas)
	numRenglonesArchivo := numElementos / number_of_chunks

	//Aquí se van a recolectar los nombres de los archivos
	var nombresDeArchivo []string

	// Dividir y guardar las partes
	for i := 0; i < number_of_chunks; i++ {
		inicio := i * numRenglonesArchivo //Desde que renglón se va a escribir en el subarchivo
		fin := inicio + numRenglonesArchivo //Hasta qué renglón se va escribir en el subarchivo

		// La última iteración tendrá solo los renglones faltantes
		if i == number_of_chunks-1 {
			fin = numElementos
		}

		// Se crea el archivo CSV para cada parte
		nArchivoSalida := "ratings_" + strconv.Itoa(i+1) + ".csv"
		
		//Aquí se van recolectando los nombres de los archivos:
		nombresDeArchivo = append(nombresDeArchivo, nArchivoSalida)
		
		//Se determina el directorio y nombre donde se guardará el archivo partido
		rutaChunk := filepath.Join(directory, nArchivoSalida)
		archivoSalida, err := os.Create(rutaChunk)
		if err != nil {
			fmt.Println("No se pudo crear el archivo: %w", err)
			return []string{}
		}
		defer archivoSalida.Close()

		// Una vez creado el archivo partido, se procede a llamar al writer para escribir en el mismo
		writer := csv.NewWriter(archivoSalida)

		// Con este ciclo se escriben las líneas en el archivo con el rango indicado
		for _, filas := range filas[inicio:fin] {
			if err := writer.Write(filas); err != nil {
				fmt.Println("No se pudo escribir en el archivo: %w", err)
				return []string{}
			}
		}

		// Se colocan en el archivo los datos del buffer
		writer.Flush()
		
		if err := writer.Error(); err != nil {
			fmt.Println("No se pudo escribir en el archivo: %w", err)
			return []string{}
		}
	}

	//Se regresan los nombres de los archivos como se indica en el requerimiento
	return nombresDeArchivo
}
