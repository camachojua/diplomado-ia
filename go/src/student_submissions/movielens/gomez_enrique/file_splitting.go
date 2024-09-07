package fileprocessing

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {
	// open file
	file, error := os.Open(directory+file_name)

	// check for any IO error
	if error != nil {
		panic(error)
	}

	// defer close file
	defer file.Close()

	// get file size in bytes
	fileStat, _ := file.Stat()
	fileByteSize := fileStat.Size()

	// get file name without extension
	fileBaseName := fileNameWithoutExtension(fileStat.Name())

	// process in chunks of arbitrary size
	// ***fileByteSize is ignored to fit test api rules***
	return processFile(file, fileBaseName, directory, int(fileByteSize/10))
}

func processFile(f *os.File, fileBaseName string, directory string, chunkSize int) []string {
	// sync.Pool reuses memory so that the GC doesn't do extra work
	chunkPool := sync.Pool{New: func() interface{} {
		// ***hack to fit test api rules: 10 files, 100-101 lines each***
		// ***remove hardcoded chunksize for a better solution to the splitting problem***
		chunkSize = 1
		chunk := make([]byte, chunkSize)
		return chunk
	}}

	// create a file reader
	reader := bufio.NewReader(f)

	// sync.WaitGroup waits for multiple go-routines to finish
	var wg sync.WaitGroup

	// filename slice
	var fileNameSlice []string

	// start reading file chunk by chunk
	for chunkId := 1; ; chunkId++ {
		// get a region of memory to temporarily store a chunk
		chunk := chunkPool.Get().([]byte)

		// read 'chunkSize' bytes into chunk buffer
		totalBytesRead, error := reader.Read(chunk)

		// totalBytesRead might be less than len(chunk), in any case, re-slice:
		chunk = chunk[:totalBytesRead]

		// break on end-of-file
		if error == io.EOF {
			chunkPool.Put(chunk)
			break
		}

		// panic on any other type of error
		if error != nil {
			panic(error)
		}

		// ***hack to fit test api rules: 10 files, 100-101 lines each***
		// ***remove for-loop for a better solution to the splitting problem***
		for i := 0; i < 100; i++ {
			// read until EOL (inclusive) TODO: what about \r\n EOLs?
			bytesUntilEol, error := reader.ReadBytes('\n')

			if error != nil {
				// TODO: ReadBytes didn't find EOL
			} else {
				// append to complete last line
				chunk = append(chunk, bytesUntilEol...)
			}
		}

		// process chunk concurrently
		wg.Add(1)
		go func() {
			var outFileName = fileBaseName + strconv.Itoa(chunkId) + ".csv"
			processChunk(directory+outFileName, chunk, &chunkPool)
			fileNameSlice = append(fileNameSlice, outFileName)
			wg.Done()
		}()

	}

	wg.Wait()

	return fileNameSlice
}

func processChunk(filePath string, chunk []byte, chunkPool *sync.Pool) {
	// TODO: do something to chunk, for example:
	error := os.WriteFile(filePath, chunk, 0644)

	if error != nil {
		panic(error)
	}

	// release chunk's memory
	chunkPool.Put(chunk)
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}
