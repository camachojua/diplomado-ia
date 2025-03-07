using CSV
using DataFrames
using Base.Threads
using Printf

# Guarda los registros procesados en un archivo CSV
function save_chunk(records, index)
    filename = "./output/ratings_$(lpad(index, 2, '0')).csv"
    CSV.write(filename, records)
end

# Divide el archivo de calificaciones en partes y procesa en paralelo
function split_ratings(total_chunks = 10)
    df = CSV.read("./ml-25m/ratings.csv", DataFrame)
    chunk_size = div(nrow(df), total_chunks)

    @threads for i in 1:total_chunks
        start_idx = (i - 1) * chunk_size + 1
        end_idx = min(i * chunk_size, nrow(df))
        chunk = df[start_idx:end_idx, :]
        save_chunk(chunk, i)
    end
end

# Procesa calificaciones y géneros para todos los archivos de calificaciones
function process_all_ratings(num_files = 10)
    genres = [
        "Action", "Adventure", "Animation", "Children", "Comedy", "Crime", 
        "Documentary", "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", 
        "Musical", "Mystery", "Romance", "Sci-Fi", "Thriller", "War", "Western", 
        "(no genres listed)"
    ]

    num_genres = length(genres)
    ratings_sum = zeros(num_genres, num_files)
    ratings_count = zeros(num_genres, num_files)

    movies_df = CSV.read("./ml-25m/movies.csv", DataFrame)
    movies_df = movies_df[:, [:movieId, :genres]]

    data_chunks = [DataFrame() for _ in 1:num_files]

    @threads for i in 1:num_files
        data_chunks[i] = CSV.read("./output/ratings_$(lpad(i, 2, '0')).csv", DataFrame)
        ratings_sum[:, i], ratings_count[:, i] = process_chunk(i, num_genres, genres, movies_df, data_chunks[i])
    end

    summarize_ratings(ratings_sum, ratings_count, genres)
end

# Procesa las calificaciones de un archivo y las agrupa por género
function process_chunk(index, num_genres, genres, movies, ratings)
    println("Processing chunk $index")
    ratings_sum = zeros(num_genres)
    ratings_count = zeros(num_genres)

    joined_data = innerjoin(movies, ratings, on=:movieId)
    for row in eachrow(joined_data)
        movie_genres = split(row.genres, "|")
        for genre in movie_genres
            if genre in genres
                idx = findfirst(==(genre), genres)
                ratings_count[idx] += 1
                ratings_sum[idx] += row.rating
            end
        end
    end

    return ratings_sum, ratings_count
end

# Calcula y muestra el promedio de calificaciones por género
function summarize_ratings(ratings_sum, ratings_count, genres)
    total_sum = sum(ratings_sum, dims=2)
    total_count = sum(ratings_count, dims=2)

    for i in 1:length(genres)
        average_rating = total_count[i] > 0 ? total_sum[i] / total_count[i] : 0.0
        @printf("Genre: %-20s Count: %10d Average: %6.2f\n", genres[i], total_count[i], average_rating)
    end
end

# Función principal para ejecutar el programa
function main(split::Bool = true)
    if split
        println("Splitting ratings into chunks...")
        @time split_ratings()
    end

    println("Processing ratings and genres...")
    @time process_all_ratings()
end

# Ejecutar el programa
@time main(true)
