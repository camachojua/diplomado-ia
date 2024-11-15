package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Función para manejar errores y evitar repetición de código
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Contar líneas en el archivo
func LineCounter(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

func FragCsvtoCsv(inpFn string) {

	inpDir := "./dat/in/"
	outDir := "./dat/out/"

	//Abrir el archivo de entrada para contar las líneas
	filePath := inpDir + inpFn + ".csv"
	fc, err := os.Open(filePath)
	checkError(err)
	defer fc.Close()

	//var rc io.Reader = fc

	totInpLns, err := LineCounter(fc) //total input lines
	checkError(err)
	fmt.Println("Number of lines in file =", totInpLns)
	fmt.Println("continue")

	// open input file again, this time for reading content
	f, err := os.Open(filePath)
	checkError(err)
	defer f.Close()

	//var r io.Reader = f
	csvReader := csv.NewReader(f)

	// Definir el número de archivos de salida
	nf := 10
	nRowsOutFile := totInpLns / nf
	if totInpLns%nf != 0 {
		nRowsOutFile++
	}

	header, err := csvReader.Read()
	checkError(err)

	//Proceso de fragmentación
	outFmt := ".csv"
	inpLnCt := 0 //Contador de lineas de entrada

	for i := 0; i < nf; i++ {
		outFn := outDir + inpFn + "_" + strconv.Itoa(i) + outFmt
		fmt.Println("FileName = ", outFn)

		csvFile, err := os.Create(outFn)
		checkError(err)
		defer csvFile.Close()

		w := csv.NewWriter(csvFile)
		//defer w.Flush()
		w.Write(header)

		for j := 0; j < nRowsOutFile; j++ {
			inpLnCt++
			rc, err := csvReader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Fatal(err)
			}
			w.Write(rc)
		}

		w.Flush()
		checkError(w.Error())
		//defer csvFile.Close() // file will be closed when the function exits
		fmt.Println("inpLnCnt=", inpLnCt)
	}

}

func main() {
	fmt.Println("FileFrag Start")

	FragCsvtoCsv("ratings")

	fmt.Println("FileFrag End")
}
