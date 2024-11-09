using DataFrames
using Glob
using CSV

function splitFile(filename::String, nbytes::Int)
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

function average(dirPath::String)
    dfMovies = DataFrame(CSV.File("../data/movies.csv"))
    select!(dfMovies, [:movieId, :genres])

    files = filter!(x->contains(x, r"^tmp.*\.csv$"), readdir(dirPath))
    genres = Dict()
    for file in files
        open(file, "r") do f
            dfRatings = DataFrame(CSV.File("../data/ratings.csv"))
            select!(dfRatings, [:movieId, :rating])

            dfMerge = groupby(innerjoin(dfMovies, dfRatings, on = :movieId), :genres)
            dfSum = combine(dfMerge, :rating => sum)

            for row in eachrow(dfSum)
                for key in split(row[1], "|")
                    if !haskey(genres, key)
                        genres[key] = 0
                    end
                    genres[key] += row[2]
                end
            end
        end
    end
    println(genres)
end

splitFile(ARGS[1], parse(Int64, ARGS[2])*1024*1024)
average(ARGS[3])
