function split!(filename::String, nbytes::Int)
    # Open file to split
    file = open(filename, "r")

    # Get file size
    fileSize = stat(filename).size

    # Init counters to keep track of bytes written
    byteCounter = 0
    fileCounter = 0
    while byteCounter < fileSize
        # Init byte array to store bytes read
        readBytes = UInt8[]

        # Init string that will contain the EOL
        charsUntilEol = ""

        # Read 'nbytes' and until EOL
        readbytes!(file, readBytes, nbytes)
        charsUntilEol = readuntil(file, '\n', keep = true)

        # Write read bytes
        fileSuffix = lpad(string(fileCounter), 2, "0") 
        outFilename = "tmp" * fileSuffix * ".csv"
        open(outFilename, "w") do f
            byteCounter += write(f, readBytes)
            byteCounter += write(f, charsUntilEol)
        end

        # Increment file suffix
        fileCounter += 1
    end

    close(file)
end

split!(ARGS[1], 50*1024*1024)
