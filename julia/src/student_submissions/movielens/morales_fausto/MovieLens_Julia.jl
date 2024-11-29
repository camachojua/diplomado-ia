Pkg.add("CSV")
Pkg.add("DataFrames")
Pkg.add("Dates")
Pkg.add("BenchmarkTools") 
Pkg.add("Printf")

using Pkg
using CSV
using DataFrames
using Dates
using BenchmarkTools
using Printf
using Base.Threads
using Dates

function crear10CSV()
  df = CSV.read("./ratings.csv", DataFrame)
  #Validar estructura del df
  #first(df, 5)

  # Contando el numero de registros del archivo
  num_rows = nrow(df)

  #Divisiones al archivo
  n = 10
  filas_archivo = num_rows/n
  Num_archivos = num_rows/filas_archivo


  #Generar nuevos archivos csv
  for i in 1:10
          inicio = (i-1)*filas_archivo+1
          if i < 10
              final = i*filas_archivo
          else i=10
              final = num_rows
          end
          inicio = ceil(Int, inicio)
          final = ceil(Int, final)
          DF_aux = df[Int(inicio):Int(final), :]
          CSV.write("Archivo_$(i).csv",DF_aux)
          @printf("Se creo: Archivo_%2d.csv \n",i)
  end

end

function FindRatingsWorker(w::Integer, ng::Integer, kg::Array, dfm::DataFrame, dfr::DataFrame)
    println("In Worker ", w, "\n")

    ra = zeros(ng) # ra is an 1D array for keeping the values of the Ratings for each genre
    ca = zeros(ng) # ca is an 1D array to keep the number of Ratings for each genre

    #println("local ndfr after resize =", size(dfr, 1))

    # The inner join will have the following columns: {movieId, genre, rating}
    ij = innerjoin(dfm, dfr, on = :movieId)
    nij = size(ij, 1)
    #println(string(nij))
    # ng = 20
    #println("nij = ", nij)
    # ng = size(kg, 1)
    for i = 1:ng
        #@printf("%d\n",i)
        for j = 1:nij
            r = ij[j,:] # get all columns for row j, gender is col=2 of the row
            g = r[2]
            #@printf("%d\n",j)
            #println(kg[j])
            if ( contains(g, kg[i]) == true)
                #println(string(j), string(kg(j)))
                ca[i] += 1    # keep the count of ratings for thin genre
                ra[i] += r[4] # add the value for this genre
            end
        end
    end
    println("Done Worker ", w, "\n")
 
    return ra, ca
end

function FindRatingsMaster()

    nF = 10 # number of files with ratings
    # kg is a 1D array that contains the Known Genders
    kg = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
      "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
       "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)" ]
  
    ng = size(kg,1)       # ng is just the number of rows in kg
    ra = zeros(ng,nF)     # ra is  2D arrayof
    ca = zeros(ng,nF)     # ra is  2D arrayof

    #####
    #ra = zeros(ng) # ra is an 1D array for keeping the values of the Ratings for each genre
    #ca = zeros(ng) # ca is an 1D array to keep the number of Ratings for each genre

    ####
  
    # dfm has all rows from Movies with cols :movieId, :genres 
    dfm = CSV.read("./movies.csv", DataFrame)
    dfm = dfm[: , [:movieId, :genres] ]

    dfr_v = [DataFrame() for _ in 1:nF]
    first(dfr_v, 5)
    @threads  for i=1:nF
    #for i=1:nF
      #rfn = prqDir * "ratings_" * string(i, pad = 2) * ".parquet"
      rfn = CSV.read("./"*"Archivo_"*string(i)*".csv", DataFrame)
      dfr_v[i] = rfn
#      first(dfr_v[i], 5)
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

  @printf("count ,   rating ,   genre , prom \n")
  @sync for i =1:ng
     @printf("%14.2f ,   %14.2f,   %s , %14.2f \n", sca[i], sra[i], kg[i], sra[i]/sca[i])
  end
end #FindRatingsMaster()

##Codigo

#dfm = CSV.read("./movies.csv", DataFrame)
#dfm = dfm[: , [:movieId, :genres] ]
#rfn = CSV.read("./"*"Archivo_"*string(1)*".csv", DataFrame)
#ij = innerjoin(dfm, rfn, on = :movieId)
#r=ij[1,:]
#println(string(r[4]))


inicio = now()
println("Inicio del proceso: $inicio")
crear10CSV()
FindRatingsMaster()
# Obtener el tiempo de fin
fin = now()
println("Fin del proceso: $fin")

# Calcular la duración del proceso
duracion = fin - inicio
println("Duración del proceso: $duracion")

  
  