package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
)

func lineCounter(r io.Reader) (int, error) {
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

// TestSplitBigFile calls fileprocessing.SplitBigFile with a test file,
// the function splits a big, bulky file into smaller chunks suitable for
// concurrent processing
func TestSplitBigFile(t *testing.T) {
	directory := "./test/"
	file_name := "ratings.csv"
	number_of_chunks := 10

	t.Log("Start the file partitioning process.")

	// We ask for splitting the file in 10 chunks
	// expecting for the files to be written in the specified directory
	files_generated := SplitBigFile(file_name, number_of_chunks, directory)

	t.Log("File partitioning ended.")
	t.Log("Validating the number of files and each file length.")

	// We try to read the "test" directory
	for file := range files_generated {
		current_file, file_error := os.Open(directory + file_name)
		t.Log(current_file)
		if file_error != nil {
			t.Log(file_error)
			t.Fail()
		}

		lines_number, count_error := lineCounter(current_file)
		if count_error != nil {
			t.Log(count_error)
			t.Fail()
		}
		if lines_number < 1000 || lines_number > 1001 {
			t.Log("The file", file, " has ", lines_number, " number of lines. We're expecting 1000 or 1001 lines at most")
			t.Fail()
		}

		current_file.Close()
	}

	if number_of_chunks != len(files_generated) {
		t.Log("got ", files_generated, " new files, wanted ", number_of_chunks)
		t.Fail()
	}
}
