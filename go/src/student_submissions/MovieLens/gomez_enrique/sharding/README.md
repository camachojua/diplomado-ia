# How to run

```console
$ go build
$ ./sharding file_to_split total_number_of_output_files
```

# sharding module vs. fileprocessing module

Both, `sharding` (this module) and `fileprocessing` modules hold the same logic, except for the fact that `sharding` doesn't fix the **total number of lines** per output file. Instead, what remains constant is the **total number of bytes** per output file, except possibly for the output file that contains the EOF character.

# Incompatibility with unit test API

This module is incompatible with the unit test API of the `fileprocessing` module since it doesn't fix the total number of lines per output file.

A compatible unit test might have to compare the total number of lines of the input file with the sum of each output file's total number of lines.
