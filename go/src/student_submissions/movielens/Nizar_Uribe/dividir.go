/*package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "time"
)

func main() {
    startTime := time.Now()

    // Cargar el archivo CSV
    df, err := loadCSV("/Users/rico/Documents/Diplomado_IA_Git/diplomado-ia/julia/src/movielens/ratings.csv")
    if err != nil {
        panic(err)
    }

    // Definir el número de partes en las que quieres dividir el archivo
    N := 10

    // Calcular el tamaño de cada parte
    numRows := len(df)
    rowsPerPart := numRows / N

    // Crear un loop para dividir el DataFrame y guardar cada parte
    for i := 0; i < N; i++ {
        startRow := i * rowsPerPart
        endRow := startRow + rowsPerPart

        // Asegurarse de no exceder el número de filas
        if i == N-1 {
            endRow = numRows
        }

        // Crear y guardar la parte como un nuevo archivo CSV
        partDF := df[startRow:endRow]
        err := saveCSV(fmt.Sprintf("ratings_part_%d.csv", i+1), partDF)
        if err != nil {
            panic(err)
        }
    }

    fmt.Printf("El archivo ha sido dividido en %d partes.\n", N)
    elapsedTime := time.Since(startTime)
    fmt.Printf("Elapsed time: %s\n", elapsedTime)
}

// loadCSV carga un archivo CSV y devuelve los datos como una matriz de cadenas
func loadCSV(filename string) ([][]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    return records, nil
}

// saveCSV guarda una parte de datos en un nuevo archivo CSV
func saveCSV(filename string, data [][]string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Escribir las filas en el archivo CSV
    return writer.WriteAll(data)
}
*/