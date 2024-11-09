using CSV
using DataFrames
using Base.Threads: @threads
using Printf

function dividir_csv(file_path::String, n::Int)
    # Leer el archivo CSV completoç
    data = CSV.File(file_path) |> DataFrame
    
    # Calcular el tamaño de cada parte
    total_rows = size(data, 1)
    rows_per_part = ceil(Int, total_rows / n)
    
    # Función para guardar una parte del DataFrame en un archivo CSV
    function guardar_parte(data_part::DataFrame, index::Int)
        part_file = "Julia/Output/ratings_$index.csv"
        CSV.write(part_file, data_part)
        println("Parte $index guardada como $part_file")
    end

    # Crear tareas concurrentes para dividir y guardar cada parte
    @threads for i in 1:n
        start_idx = (i - 1) * rows_per_part + 1
        end_idx = min(i * rows_per_part, total_rows)
        data_part = data[start_idx:end_idx, :]
        guardar_parte(data_part, i)
    end
end

dividir_csv("Julia/ratings.csv", 10)

# Worker to proscces the data for the inner join between movies and ratings
function FindRatingsWorker( wId::Integer, numGen::Integer, gArr::Array, dfmv::DataFrame, dfrt::DataFrame )
    println("In Worker  ", wId, "\n")
    raVal= zeros(numGen)    # is an 1D array for keeping the values of the Ratings for each genre
    caCnt= zeros(numGen)    # is an 1D array to keep the number of Ratings for each genre
  
    # the innerjoin will have the following columns: {movieId, genre, rating}
    ij = innerjoin(dfmv, dfrt, on = :movieId)
    nij = size(ij,1)
    println("Size of inner-join ij = ", nij)
  
    # println("nij = ", nij)
    # ng = size(kg,1)
    for i = 1:numGen
      for j = 1:nij
        r = ij[j,:]       # get all columns for row j. gender is col=2 of the row
        g = r[2] 
        if ( contains( g , gArr[i]) == true)
            caCnt[i] += 1      # keep the count of ratings for this genre
            raVal[i] += r[4]   #add the value for this genre
        end
      end
    end
   return raVal, caCnt
  end


  function FindRatingsMaster()

    nF = 10 # number of files with ratings
    gArr = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
          "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
          "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)" ]
  
  
    numGen = size(gArr,1)       # ng is just the number of rows in gArr
    raVal = zeros(numGen,nF)     # ra is  2D arrayof
    caCnt = zeros(numGen,nF)     # ra is  2D arrayof
  
    dfmv = CSV.File("Julia/movies_large.csv") |> DataFrame
    dfmv = dfmv[: , [:movieId, :genres] ] |> DataFrame
    
    ##slide 29
    dfr_v = [DataFrame() for _ in 1:nF]
    @threads  for i=1:nF
      # for i= 1:10
      #a for cycle for this
      rfn = "Julia/Output/ratings_" * string(i) * ".csv"
      println(rfn) 
      dfr_v[i] = CSV.File(rfn) |> DataFrame
      raVal[:,i] , caCnt[:,i] = FindRatingsWorker( i, numGen, gArr, dfmv, dfr_v[i])
    end # @threads for 
  
      # end # @everywhere  
      # sra is an 1D array for summing the values of the Ratings for each genre
      sra = zeros(numGen)     
      # sca is an 1D array for summing the counts of the Ratings for each genre
      sca = zeros(numGen)     
      
      @sync for i =1:numGen
              for j = 1:nF
                sra[i] += raVal[i,j]
                sca[i] += caCnt[i,j]
                
              end
            end
  
      @sync for i =1:numGen
          @printf("sca = %14.2f   sra = %14.2f  avgra = %14.2f  genre = %s  \n", sca[i], sra[i], (1.0*sra[i])/(sca[i]), gArr[i])
      end
  
  end 

  FindRatingsMaster()


