### Movie Lens en julia
### Mar Bazúa 
using CSV
using DataFrames
using Base.Threads
using BenchmarkTools
using Printf

# Procesa registros en chunks, combina géneros y guarda en archivos de salida
function process_records(records, output_index)
    # Guardar el chunk en un archivo CSV
    output_filename = "./output/ratings_$(lpad(output_index, 2, '0')).csv"
    CSV.write(output_filename, records)
end

# Función principal que divide el archivo de ratings en chunks y los procesa en paralelo
function Split_Ratings(total_jobs = 10)
    # Leer archivos CSV
    ratings_df = CSV.read("./ml-25m/ratings.csv", DataFrame)

    # Determinar tamaño de cada chunk
    size_range = div(nrow(ratings_df), total_jobs)

    # Procesar en paralelo utilizando múltiples hilos
    @threads for i in 1:total_jobs
        println("In Worker ", i, " to split Ratings \n")
        start_idx = (i - 1) * size_range + 1
        end_idx = min(i * size_range, nrow(ratings_df))
        records_chunk = ratings_df[start_idx:end_idx, :]
        process_records(records_chunk, i)
    end
end

function FindRatingsMaster(nF = 10)
    #nF number of files with ratings
    # kg is a 1D array that contains the Known Genders
    kg = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
      "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
       "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"]
  
    ng = size(kg,1)       # ng is just the number of rows in kg
    ra = zeros(ng,nF)     # ra is  2D arrayof
    ca = zeros(ng,nF)
  
    # dfm has all rows from Movies with cols :movieId, :genres 
    dfm = CSV.read("./ml-25m/movies.csv", DataFrame)
    dfm = dfm[: , [:movieId, :genres] ]
  
    dfr_v = [DataFrame() for _ in 1:nF]
    @threads  for i=1:nF
    #for i=1:nF
        #rfn = CSV.read("./output/ratings_$(lpad(i, 2, '0')).csv", DataFrame)
        #println( rfn ) 
        dfr_v[i] = CSV.read("./output/ratings_$(lpad(i, 2, '0')).csv", DataFrame)
        ra[:,i] , ca[:,i] = FindRatingsWorker( i, ng, kg, dfm, dfr_v[i])
    end # @threads for 

    # end # @everywhere  
    # sra is an 1D array for summing the values of the Ratings for each genre
    sra = zeros(ng)     
    sca = zeros(ng)     
    @sync for i =1:ng
            for j = 1:nF
            sra[i] += ra[i,j]
            sca[i] += ca[i,j]
            end
        end

    @sync for i in 1:ng
    @printf("count = %14.2f   average = %14.2f   genre = %s\n", sca[i], sra[i]/sca[i], kg[i])
  end

end #FindRatingsMaster()

function FindRatingsWorker(w::Integer, ng::Integer, kg::Array, dfm::DataFrame, dfr::DataFrame)
    println("In Worker ", w, " to process Ratings with Movielens \n")

    ra = zeros(ng) # ra is an 1D array for keeping the values of the Ratings for each genre
    ca = zeros(ng) # ca is an 1D array to keep the number of Ratings for each genre

    # The inner join will have the following columns: {movieId, genre, rating}
    ij = innerjoin(dfm, dfr, on = :movieId)
    #println("El encabezado es: ", first(ij))
    nij = size(ij, 1)

    for i = 1:ng
        for j = 1:nij
            r = ij[j,:] # get all columns for row j, gender is col=2 of the row
            g = r[2]
            if (contains(g, kg[i]) == true)
                ca[i] += 1    # keep the count of ratings for thin genre
                ra[i] += r[4] # add the value for this genre
            end
        end
    end
    
    return ra, ca
end

function main(execute_split::Bool)
    if execute_split == true
        println("Split Ratings will be executed: \n")
        @time Split_Ratings()
    end
    
    println("Join Ratings with Movielens using threads will be executed: \n")
    @time FindRatingsMaster()
end

# Main Function 
@time main(true)