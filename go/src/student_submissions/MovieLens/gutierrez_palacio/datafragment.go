package main

//Autores: Mariana Gutiérrez y Rafael Palacio

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	//"strconv"
	//"strings"
	"time"
	//"github.com/rocketlaunchr/dataframe-go"
)

func main() {
	start := time.Now()

	SegmFiles("ratings", 10)
	//Mt_FindRatingsMaster()

	elapsed := time.Since(start)
	fmt.Printf("The code took %s to execute\n", elapsed)
}

func SegmFiles(filename string, numfiles int) {
	//filename := "ratings"
	//numfiles := 10
	fmt.Printf("INPUT: %s.csv\n", filename)

	file, err := os.Open(filename + ".csv")
	if err != nil {
		log.Fatal("Error while opening the file:", err)
	}
	defer file.Close()

	numlines, err := LineCounter(file)
	if err != nil {
		fmt.Println("Error counting lines:", err)
		return
	}

	//fmt.Println("Number of lines in the file:", numlines)
	linesPerFile := numlines / numfiles
	fmt.Printf("The %d new files will each contain approximately %d lines\n", numfiles, linesPerFile)

	file.Seek(0, io.SeekStart)
	lineCount := 0
	for i := 0; i < numfiles; i++ {
		namearchivofrag := filename + fmt.Sprintf("%d.csv", i)
		archivofrag, err := os.Create(namearchivofrag)
		if err != nil {
			fmt.Println("Error creating the file:", err)
			return
		}
		defer archivofrag.Close()

		scanner := bufio.NewScanner(file)
		writer := bufio.NewWriter(archivofrag)

		for scanner.Scan() {
			lineCount++
			if lineCount > linesPerFile*i && lineCount <= linesPerFile*(i+1) {
				line := scanner.Text()
				_, err := writer.WriteString(line + "\n")
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			} else if lineCount > linesPerFile*(i+1) {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading the file:", err)
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("Error flushing to file:", err)
		}
		fmt.Printf("OUTPUT: %s\n", namearchivofrag)
	}
}

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}
