package fileprocessing

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func SplitBigFile(fileName string, numberOfChunks int, directory string) []string {
	filePath := fmt.Sprintf("%s/%s", directory, fileName)
	file, err := os.Open(filePath)
	if err != nil{
		log.Fatal(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	partitions := make(map[string][][]string)
	recordNum := 0
	headers, err := csvReader.Read()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			log.Fatal(err)
		}
		process(numberOfChunks, recordNum, record, partitions)
		recordNum++
	}

	var partitionsChannel = make(chan int)
	for i := 0; i < len(partitions); i++ {
		go savePartitions(filePath, i, partitionsChannel, partitions, headers)
	}
	iMsg := 0
	go func(){
		for {
			i := <-partitionsChannel
			iMsg += i
		}
	}()
	for {
		if iMsg == len(partitions) {
			break
		}
	}

	var partitionFiles []string
	for key := range partitions{
		partitionFileName := fmt.Sprintf("%s_partition_%s.csv", fileName, key)
		partitionFiles = append(partitionFiles, partitionFileName)
	}
	return partitionFiles
}

func process(numberOfChunks int, workerNum int, record []string, partitions map[string][][]string){
	partition := strconv.Itoa(workerNum % numberOfChunks + 1)
	if value, exists := partitions[partition]; exists {
		partitions[partition] = append(value, record)
	} else {
		partitions[partition] = [][]string{record}
	}
}

func savePartitions(filePath string, workerNum int, channel chan int, partitions map[string][][]string, headers []string){
	partitionFileName := fmt.Sprintf("%s_partition_%d.csv", filePath, workerNum + 1)
	partition := strconv.Itoa(workerNum + 1)
	partitionRows := partitions[partition]
	partitionRowsWithHeaders := append([][]string{headers}, partitionRows...)
	csvFile, err := os.OpenFile(partitionFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.WriteAll(partitionRowsWithHeaders)
	csvWriter.Flush()
	channel <- 1
}