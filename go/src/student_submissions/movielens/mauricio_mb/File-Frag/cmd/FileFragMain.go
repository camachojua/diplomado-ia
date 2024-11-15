package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const inpDir string = "../dat/in/"
const outDir string = "../dat/out/"

func FragCsvToCsv(inpFn string) {
	// open input file for couting number of lines
	fc, err := os.Open(inpDir + inpFn + ".csv")
	if err != nil {
		log.Fatal(err)
	}

	defer fc.Close()
	var rc io.Reader = fc

	totInpLns, err := LineCounter(rc) // total input lines
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Number of lines in file = ", totInpLns)
	fmt.Println("continue")
	totInpLns -= 1
	// open input file again, this time for reading content
	f, err := os.Open(inpDir + inpFn + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var r io.Reader = f
	csvReader := csv.NewReader(r)

	// read the header
	header, err := csvReader.Read()
	if err != nil {
		log.Fatal("Error reading header:", err)
	}

	// define the number of files and how many rows for the file
	nf := 10
	nRowsOutFile := totInpLns / nf
	outFmt := ".csv"
	inpLnCt := 0 // input line count
	start := time.Now()

	for i := 0; i < nf; i++ {
		outFn := outDir + inpFn + "_" + fmt.Sprintf("%02d", (i+1)) + outFmt
		fmt.Println("FileName =", outFn)
		csvFile, err := os.Create(outFn)
		if err != nil {
			log.Fatal(err)
		}

		w := csv.NewWriter(csvFile)
		//defer w.Flush()
		if err := w.Write(header); err != nil {
			log.Fatal("Error writing header:", err)
		}

		//w.Flush()
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
		csvFile.Close() // file will be closed when the function finish
		fmt.Println("inpLnCnt=", inpLnCt)
	}
	duration := time.Since(start)
	fmt.Println("Duration =", duration)
}

func LineCounter(r io.Reader) (int, error) {
	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}
	return count, nil
}

func testLineCounter(inpFn string) {
	fc, err := os.Open(inpDir + inpFn + ".csv")
	if err != nil {
		log.Fatal(err)
	}

	defer fc.Close()
	var rc io.Reader = fc
	totInpLns, err := LineCounter(rc) // total input lines
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Number of lines in file =", totInpLns)
}

func main() {
	fmt.Println("FileFrag Start")
	//testLineCounter("ratings")
	//FragCsvToCsv("ratings")
	FragCsvToCsv("ratings")
	fmt.Println("FileFrag End")

}
