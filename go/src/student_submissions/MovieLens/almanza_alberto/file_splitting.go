package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"sync"
)

func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {
	var chunks []string
	
	file, _ := os.Open(directory + file_name)
	data, _ := csv.NewReader(file).ReadAll()
	file.Close()
	fmt.Printf("\nFile with %d rows\n", len(data))
	
	chunk_size := len(data) / number_of_chunks 
	fmt.Printf("Chunk size: %d\n\n", chunk_size)
	
	start_time:= time.Now()
	fmt.Println("Tiempo inicial:", start_time)

	var wg sync.WaitGroup
	for i := 0; i <= number_of_chunks; i++ {
		chunk_name := fmt.Sprintf("CHUNK_%d.csv", i)
		chunk_path := fmt.Sprintf("%s%s",directory, chunk_name)
		end_slice := chunk_size * (i + 1)
		if end_slice > len(data) { end_slice = len(data) }
		chunk_data := data[chunk_size * i:end_slice]
		len_chunk := len(chunk_data)
		if len_chunk == 0 { break }
		wg.Add(1)
		go func() {
			defer wg.Done()
			chunk_file, _ := os.Create(chunk_path)
			writer := csv.NewWriter(chunk_file)
			writer.WriteAll(chunk_data)
			writer.Flush()
			chunk_file.Sync()
			chunk_file.Close()
			fmt.Printf("file %s (size: %d)\n", chunk_path, len_chunk)
			
		}()
		chunks = append(chunks, chunk_name)
	}

	end_time := time.Now()
	fmt.Println("Tiempo final:", end_time)
	fmt.Println("Tiempo transcurrido:", end_time.Sub(start_time).Seconds())

	wg.Wait()
	
	return chunks
}

func main() {
	fmt.Println(SplitBigFile("ratings.csv", 10, "./test/"))
}
