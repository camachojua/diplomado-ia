package fileprocessing


import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

// Definimos la estructura de Ratings

type Rating struct {
	UserID    int64   `parquet:"name=userid, type=INT64"`
	MovieID   int64   `parquet:"name=movieid, type=INT64"`
	Rating    float64 `parquet:"name=rating, type=DOUBLE"`
	Timestamp int64   `parquet:"name=timestamp, type=INT64"`
}

func SplitBigFile(file_name string, number_of_chunks int, directory string) ([]string, error) {

	file, err := os.Open(file_name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Preparamos el archivo para leerlo
	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for i := range headers {
		headers[i] = strings.ToLower(headers[i])
	}

	var wg sync.WaitGroup
	rowsPerChunk := make([][][]string, number_of_chunks)
	chunkSize := 0

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}

		// Asignamos el número de chunk en el que estamos trabajando
		chunkIndex := chunkSize % number_of_chunks
		rowsPerChunk[chunkIndex] = append(rowsPerChunk[chunkIndex], record)
		chunkSize++
	}

	// Creamos un slice para almacenar los nombres de los archivos generados
	nombres := make([]string, number_of_chunks)

	// Definimos el ciclo donde vamos a trabajar
	for j := 0; j < number_of_chunks; j++ {
		wg.Add(1)

		go func(chunkNum int, rows [][]string) {
			defer wg.Done()

			// Creamos la ruta, el nombre del archivo y lo mandos a la lista de nombres

			fileName := fmt.Sprintf("Rating_%d.parquet", chunkNum)
			outputFilePath := filepath.Join(directory, fileName)
			nombres[chunkNum] = fileName

			fw, err := local.NewLocalFileWriter(outputFilePath)
			if err != nil {
				log.Printf("Error creando archivo parquet %d: %v", chunkNum, err)
				return
			}
			defer fw.Close()

			pw, err := writer.NewParquetWriter(fw, new(Rating), 32)
			if err != nil {
				log.Printf("Error creando writer para el chunk %d: %v", chunkNum, err)
				return
			}
			defer pw.WriteStop()

			for _, record := range rows {
				var rating Rating
				rating.UserID, _ = strconv.ParseInt(record[0], 10, 64)
				rating.MovieID, _ = strconv.ParseInt(record[1], 10, 64)
				rating.Rating, _ = strconv.ParseFloat(record[2], 64)
				rating.Timestamp, _ = strconv.ParseInt(record[3], 10, 64)
				if err := pw.Write(rating); err != nil {
					log.Printf("Error escribiendo Rating en chunk %d: %v", chunkNum, err)
				}
			}

		}(j, rowsPerChunk[j])
	}

	wg.Wait()
	return nombres, nil // Devuelve la lista de nombres

}

func fileprocessing() {

	nombres, err := SplitBigFile("ratings.csv", 50, `test/`)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	for _, nombre := range nombres {
		fmt.Println(nombre)
	}
}

