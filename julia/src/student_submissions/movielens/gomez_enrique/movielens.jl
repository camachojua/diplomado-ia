using DataFrames
using Glob
using CSV
using Printf

mutable struct Stats
    rating::Float64
    observations::Int64
end

function splitBigFile(filename::String, nbytes::Int)
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

function processFile!(filename::String, dfMovies::DataFrame, stats::Dict)
    dfRatings = DataFrame(CSV.File("../data/ratings.csv"))
    select!(dfRatings, [:movieId, :rating])

    dfMerge = groupby(innerjoin(dfRatings, dfMovies, on = :movieId), :genres)
    dfCombine = combine(dfMerge, [:rating] .=> [sum length])

    for row in eachrow(dfCombine)
        # Rating
        rating = row[2]
        observations = row[3]

        # Split genres
        genres = split(row[1], "|")

        # Count per genre
        @sync for key in genres
            if !haskey(stats, key)
                stats[key] = Stats(0.0, 0)
            end
            stats[key].rating += rating
            stats[key].observations += observations
        end
    end
end

splitBigFile(ARGS[1], parse(Int64, ARGS[2])*1024*1024)

dfMovies = DataFrame(CSV.File("../data/movies.csv"))
select!(dfMovies, [:movieId, :genres])

stats = Dict{String, Stats}()

files = filter!(x->contains(x, r"^tmp.*\.csv$"), readdir(ARGS[3]))
Threads.@threads for file in files
    processFile!(file, dfMovies, stats)
end

for (key, value) in stats
    @printf("%s,%f\n", key, value.rating/value.observations)
end
