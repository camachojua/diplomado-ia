using CSV
using DataFrames
using Base.Threads
using Printf

using CSV
using DataFrames
using Base.Threads: @threads

function dividir_archivo_en_partes_threaded(input_file::String, num_partes::Int)
    # Primero, contamos el total de filas
    total_filas = 0
    for _ in CSV.File(input_file)
        total_filas += 1
    end
    
    # Calculamos el tamaño de cada parte
    filas_por_parte = div(total_filas, num_partes)
    filas_restantes = total_filas % num_partes  # Filas adicionales para distribución

    # Creamos una tarea para cada parte usando @threads
    @threads for i in 1:num_partes
        # Calcular el inicio de cada partef
        skipto = (i - 1) * filas_por_parte + 1
        # Definir el límite de filas para cada parte
        num_filas = filas_por_parte + (i <= filas_restantes ? 1 : 0)

        # Leer y escribir la parte actual
        chunk = CSV.File(input_file; skipto=skipto, limit=num_filas) |> DataFrame
        output_file = "ratings_parte_$i.csv"
        CSV.write(output_file, chunk)
        
        println("Parte $i guardada en $output_file con $num_filas filas")
    end
end




##slide 31
function FindRatingsWorker( w::Integer, ng::Integer, kg::Array, dfm::DataFrame, dfr::DataFrame )
  println("In Worker  ", w, "\n")
  ra= zeros(ng)    # ra is an 1D array for keeping the values of the Ratings for each genre
  ca= zeros(ng)    # ca is an 1D array to keep the number of Ratings for each genre

  # the innerjoin will have the following columns: {movieId, genre, rating}
  ij = innerjoin(dfm, dfr, on = :movieId)
  nij = size(ij,1)
  println("Size of inner-join ij = ", nij)

  # println("nij = ", nij)
  # ng = size(kg,1)
  for i = 1:ng
    for j = 1:nij
      r = ij[j,:]       # get all columns for row j. gender is col=2 of the row
      g = r[2] 
      if ( contains( g , kg[i]) == true)
          ca[i] += 1      # keep the count of ratings for this genre
          ra[i] += r[3]   #add the value for this genre
      end
    end
  end
 return ra, ca
end

function FindRatingsMaster()

  nF = 10 # number of files with ratings
  kg = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
        "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
        "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)" ]


  ng = size(kg,1)       # ng is just the number of rows in kg
  ra = zeros(ng,nF)     # ra is  2D arrayof
  ca = zeros(ng,nF)     # ra is  2D arrayof

  dfm = CSV.File("movies_large.csv") |> DataFrame
  dfm = dfm[: , [:movieId, :genres] ] |> DataFrame
  
  ##slide 29
  dfr_v = [DataFrame() for _ in 1:nF]
  @threads  for i=1:nF
    # for i= 1:10
    #a for cycle for this
    rfn = "ratings_parte_" * string(i) * ".csv"
    println(rfn) 
    dfr_v[i] = CSV.File(rfn) |> DataFrame
    ra[:,i] , ca[:,i] = FindRatingsWorker( i, ng, kg, dfm, dfr_v[i])
  end # @threads for 

    # end # @everywhere  
    # sra is an 1D array for summing the values of the Ratings for each genre
    sra = zeros(ng)     
    # sca is an 1D array for summing the counts of the Ratings for each genre
    sca = zeros(ng)     
    @sync for i =1:ng
            for j = 1:nF
              sra[i] += ra[i,j]
              sca[i] += ca[i,j]
            end
          end

    @sync for i =1:ng
        @printf("sca = %14.2f   sra = %14.2f   genre = %s  \n", sca[i], sra[i], kg[i])
    end

end #FindRatingsMaster()

# Uso: dividir el archivo en 10 partes
dividir_archivo_en_partes_threaded("ratings_large.csv", 10)
# Encuentro los ratings
FindRatingsMaster()