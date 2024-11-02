using Parquet
using DataFrames

# Función para procesar cada archivo Parquet en paralelo
function find_ratings_worker(worker_id::Int, genres::Vector{String}, movies::Dict, count_array::Array{Int, 2}, sum_array::Array{Float64, 2})
    filename = "output_parquet/ratings_part_$worker_id.parquet"
    println("$worker_id está procesando el archivo $filename")

    # Leer el archivo Parquet y convertirlo a DataFrame
    df_chunk = DataFrame(Parquet.read_parquet(filename))

    for row in eachrow(df_chunk)
        movie_id = row[:movieId]
        rating_value = row[:rating]
        
        if haskey(movies, movie_id)
            movie_genres = movies[movie_id]

            # Contar calificaciones por género
            for (i, genre) in enumerate(genres)
                if contains(movie_genres, genre)
                    count_array[i, worker_id] += 1
                    sum_array[i, worker_id] += rating_value
                end
            end
        end
    end

    println("$worker_id completó el procesamiento.")
end
