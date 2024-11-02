import Pkg;
Pkg.offline(true)
using Base.Threads
using BenchmarkTools
using Parquet
using CSV
using DataFrames
#Pkg.add("Query")
using Query
using Base.Threads
using Printf: @printf

#include("/home/luna/Schreibtisch/IADC/Julia/FileProcess.jl")
# Directories
mlDir = "/home/luna/Downloads/ml-25m/"
prqDir = "/home/luna/Schreibtisch/IADC/Julia/prq/"

# 1
# Read CSV files
fileMovies = "movies"
fileRatings = "ratings"

function filesToParquet(filename::String, fileDestiny::String)
#  Create Parquet files for movies.csv and ratings.csv
    println("Creating parquet file for: ", filename)
    df =  CSV.read( filename, DataFrame)
    Parquet.write_parquet(fileDestiny, compression_codec = "ZSTD", df)
    println("Save file as:", fileDestiny )
end

function ReadParquetAsDf(filename::String)
    df = DataFrame(read_parquet( filename)) #, rows=nRows))
   rows = nrow(df)
 #println("Numero de filas: ", rows)
    return df, rows

end

# 2 
# Read and divide into N pieces.
function readDivideRatings(filename ::String, chunkSize, fileResult ::String)
    # Open file for lecture
        dfRatings = DataFrame(read_parquet(filename))
       # Read all file lines
        lines = eachrow(dfRatings)
        totalLines = nrow(dfRatings)
        
        # Create divisions for every file.
        linesPPart = div(totalLines, chunkSize)
        
       @threads for i in 1:chunkSize
            # Calcular el rango de líneas para la parte actual
            first   = (i - 1) * linesPPart + 1
            last = (i == 10) ? totalLines : i * linesPPart # min(i * linesPPart, totalLines)
            part = dfRatings[first:last, :]

            # Crear un nuevo DataFrame con la cabecera como la primera fila
           # header = dfRatings[1,:] #DataFrame(collect(keys(part[1, :])))
          #println("Creating parquet file for: ", header)
            #partWHeader = vcat(header, part)    
            # Crear el nombre del nuevo archivo
            fileRes = fileResult * "_$(i).parquet"
            Parquet.write_parquet(fileRes, compression_codec = "ZSTD", part)
        println("Creating parquet file for: ", fileRes)
        end # @threads for
               
end

function findRatingsWorker( w::Integer, ng::Integer, kg::Array, dfm::DataFrame, dfr::DataFrame )
    println("In Worker ", w, "\n")
    ra = zeros(ng) # ra is an 1D array for keeping the values of the Ratings for each genre
    ca = zeros(ng) # ca is an 1D array to keep the number of Ratings for each genre
    # the innerjoin will have the following columns: {movieId, genre, rating}
    ij = innerjoin(dfm, dfr, on = :movieId)
    nij = size(ij,1)
    println("Size of inner-join ij = ", nij)
    # println("nij = ", nij)
    # ng = size(kg,1)
        for i = 1:ng
            for j = 1:nij
            r = ij[j,:] # get all columns for row j. gender is col=2 of the row
            g = r[2]
                if ( contains( g , kg[i]) == true)
                ca[i] += 1 # keep the count of ratings for this genre
                ra[i] += r[3] #add the value for this genre
                end
            end
        end
    return ra, ca
end

function findRatingsMaster(filename)
    println("In master")
    nF = 10 # number of files with ratings
    # kg is a 1D array that contains the Known Genders
    kg = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
    "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
    "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)" ]
    ng = size(kg,1) # ng is just the number of rows in kg
    ra = zeros(ng,nF) # ra is 2D arrayof
    ca = zeros(ng,nF) # ra is 2D arrayof
    # dfm has all rows from Movies with cols :movieId, :genres
    dfm = DataFrame(read_parquet( movies))
    dfm = dfm[: , [:movieId, :genres] ]
    dfr_v = [DataFrame() for _ in 1:nF]
    @threads for i=1:nF
    #    for i=1:nF
            #sleep(1)
        rfn = filename * string(i, pad = 1) * ".parquet"
        println( rfn )
        println(
            "word",
            "\tthread: $(Threads.threadid())"
        )
        dfr_v[i] = DataFrame(read_parquet( rfn ))
        ra[:,i] , ca[:,i] = findRatingsWorker( i, ng, kg, dfm, dfr_v[i])
    end # @threads for
    # end # @everywhere
    # sra is an 1D array for summing the values of the Ratings for each genre
    sra = zeros(Float64, ng)
    # sca is an 1D array for summing the counts of the Ratings for each genre
    sca = zeros(ng)
    mean  = zeros(ng)
    @sync for i =1:ng
        for j = 1:nF
            sra[i] += ra[i,j]
            sca[i] += ca[i,j]
            mean[i] = div(sra[i], Float64(sca[i]))
        end
    end
    @sync for i =1:ng
    @printf("ca = %14.2f ra = %14.2f genre = %s mean = %14.2f \n", sca[i], sra[i], kg[i], mean[i])
    end
end #FindRatingsMaster()
JULIA_NUM_THREADS=6
println("Threads disponibles en Julia: ", Threads.nthreads())

function createParquet()
#Create Parquet file for movies.csv
filename = mlDir * fileMovies * ".csv"
fileDestiny = prqDir * fileMovies * ".parquet"
#c_fastmath_bench = @benchmark $filesToParquet(filename, fileDestiny);
tiempo = @elapsed filesToParquet(filename, fileDestiny);
println("Tiempo de ejecución: $tiempo segundos")

#Create Parquet files for movies.csv and ratings.csv
filename = mlDir * fileRatings * ".csv"
fileDestiny = prqDir * fileRatings * ".parquet"
#c_fastmath_bench = @benchmark $filesToParquet(filename, fileDestiny);
tiempo = @elapsed filesToParquet(filename, fileDestiny);
println("Tiempo de ejecución: $tiempo segundos")
end 

function splitRatings()

    filename = prqDir * fileRatings * ".parquet"
    fileResult = prqDir * fileRatings
    # c_fastmath_bench = @benchmark $readDivideRatings(filename, 10);

    tiempo = @elapsed readDivideRatings(filename, 10, fileResult);
    println("Tiempo de ejecución: $tiempo segundos")
    #println("C: Fastest time was $(minimum(c_fastmath_bench.times) / 1e9) sec") # in mi
end

# Final paths
    movies = prqDir * "movies.parquet"
    ratings = prqDir * "ratings_"

	println("Comenzaremos a leer el archivo 'movies.parquet'")
	#= Read Movies.parquet =#
	dfMovies, rows = ReadParquetAsDf(movies)
	println("Se termino de leer el archivo 'movies.csv'")
	println("Numeros de registros segun DataFrame: ", rows)

 #DONE 	
     tiempo = @elapsed findRatingsMaster(ratings) #// <= Este es el orquestador
    println("Tiempo de ejecución: $tiempo segundos")