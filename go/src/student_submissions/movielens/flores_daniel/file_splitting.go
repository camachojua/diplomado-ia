package fileprocessing

func SplitBigFile(file_name string, number_of_chunks int, directory string) []string {

	return mySplitFile(file_name, number_of_chunks, directory)
}
