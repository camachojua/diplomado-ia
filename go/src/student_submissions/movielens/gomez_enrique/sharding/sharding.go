package sharding

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	// get cli arguments
	filename := os.Args[1]
	totalChunks, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	// open file
	file, err := os.Open(filename)

	// check for any IO err
	if err != nil {
		panic(err)
	}

	// get file size in bytes
	fileStat, _ := file.Stat()
	fileByteSize := fileStat.Size()

	// get file name without extension
	fileBaseName := "tmp-" + fileNameWithoutExtension(fileStat.Name())

	// process in chunks of arbitrary size
	processFile(file, fileBaseName, fileByteSize/int64(totalChunks))

	// close file
	file.Close()
}

func processFile(f *os.File, fileBaseName string, chunkSize int64) {
	// sync.Pool reuses memory so that the GC doesn't do extra work
	chunkPool := sync.Pool{New: func() interface{} {
		chunk := make([]byte, chunkSize)
		return chunk
	}}

	// create a file reader
	reader := bufio.NewReader(f)

	// sync.WaitGroup waits for multiple go-routines to finish
	var wg sync.WaitGroup

	// start reading file chunk by chunk
	for chunkId := 1; ; chunkId++ {
		// get a region of memory to temporarily store a chunk
		chunk := chunkPool.Get().([]byte)

		// read 'chunkSize' bytes into chunk buffer
		totalBytesRead, err := io.ReadFull(reader, chunk)

		// totalBytesRead might be less than len(chunk), in any case, re-slice:
		chunk = chunk[:totalBytesRead]

		// break on end-of-file (no bytes were read)
		if err == io.EOF {
			// release chunk's memory space
			chunkPool.Put(chunk)
			break
		}

		bytesUntilEol, err := reader.ReadBytes('\n')

		if err != nil {
			// TODO?: ReadBytes didn't find EOL
		} else {
			// append to complete last line
			chunk = append(chunk, bytesUntilEol...)
		}

		// process chunk concurrently
		wg.Add(1)
		go func() {
			processChunk(fileBaseName+strconv.Itoa(chunkId)+".csv", chunk, &chunkPool)
			wg.Done()
		}()

	}

	wg.Wait()
}

func processChunk(outFileName string, chunk []byte, chunkPool *sync.Pool) {
	// TODO: do something to chunk, for example:
	err := os.WriteFile(outFileName, chunk, 0644)

	if err != nil {
		panic(err)
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
