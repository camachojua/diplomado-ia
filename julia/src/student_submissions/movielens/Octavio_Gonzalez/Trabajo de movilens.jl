#import Pkg; Pkg.add("CSV")
#import Pkg; Pkg.add("DataFrames")

using DataFrames
using CSV
using Parquet
using Printf

# Función para dividir el archivo ratings en 10 partes y guardarlas en formato Parquet
function generate_small_files(ratings_file::String, output_prefix::String, output_dir::String)
    println("Dividiendo archivo de ratings en partes...")
    
    # Leer el archivo completo de ratings
    data = CSV.read(ratings_file, DataFrame)
    total_rows = nrow(data)
    num_chunks = 10  # Dividir en 10 partes
    chunk_size = ceil(Int, total_rows / num_chunks)

    # Crear y guardar cada chunk en archivos separados
    for i in 1:num_chunks
        start_row = (i - 1) * chunk_size + 1
        end_row = min(i * chunk_size, total_rows)
        chunk = data[start_row:end_row, :]

        output_path = joinpath(output_dir, "$(output_prefix)_ratings$(lpad(i, 2, '0')).parquet")
        Parquet.write_parquet(output_path, chunk)
        println("Archivo guardado: $output_path con $(nrow(chunk)) filas")
    end
end

# Función principal para procesar archivos de ratings y cruzarlos con movies
function find_ratings_master(input_dir::String, output_dir::String)
    nF = 10  # Número de archivos ratings
    prqDir = output_dir

    # Lista de géneros de películas
    genres = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
              "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
              "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"]
    ng = length(genres)

    # Arrays para acumular los resultados de calificaciones por género
    rating_sum = zeros(ng, nF)
    count_sum = zeros(Int, ng, nF)

    # Leer el archivo movies.csv y mantener solo las columnas necesarias
    movies_path = joinpath(input_dir, "movies.csv")
    df_movies = CSV.read(movies_path, DataFrame)
    df_movies = df_movies[:, [:movieId, :genres]]

    # Procesar cada archivo de ratings
    for i in 1:nF
        rating_file = joinpath(prqDir, "ratings_ratings$(lpad(i, 2, '0')).parquet")
        println("Procesando archivo: $rating_file")

        if isfile(rating_file)
            df_ratings = DataFrame(read_parquet(rating_file))
            rating_sum[:, i], count_sum[:, i] = process_ratings(ng, genres, df_movies, df_ratings)
        else
            println("Archivo no encontrado: $rating_file")
        end
    end

    # Sumar resultados finales por género y mostrarlos
    for i in 1:ng
        total_rating = sum(rating_sum[i, :])
        total_count = sum(count_sum[i, :])
        promedio = total_rating/total_count 
        @printf("Género: %s   Total calificaciones: %.2f   Total conteo: %d   Promedio: %.2f\n", genres[i], total_rating, total_count,promedio) 
    end
end

# Función para procesar cada archivo ratings y acumular resultados por género
function process_ratings(ng::Int, genres::Vector{String}, df_movies::DataFrame, df_ratings::DataFrame)
    rating_accum = zeros(ng)
    count_accum = zeros(Int, ng)

    # Hacer un inner join entre movies y ratings
    joined_df = innerjoin(df_movies, df_ratings, on=:movieId)

    # Calcular sumas y conteos por cada género
    for i in 1:ng
        genre_rows = joined_df[occursin.(genres[i], joined_df.genres), :]
        count_accum[i] = nrow(genre_rows)
        rating_accum[i] = sum(genre_rows.rating)
    end

    return rating_accum, count_accum
end

# Rutas de entrada y salida
input_dir = "C:\\Users\\Alexis\\Documents\\Practica_Julia"
output_dir = "C:\\Users\\Alexis\\Documents\\Practica_Julia"

# Dividir el archivo ratings y guardar en formato Parquet
generate_small_files(joinpath(input_dir, "ratings.csv"), "ratings", output_dir)

# Procesar los archivos de ratings y cruzarlos con movies
find_ratings_master(input_dir, output_dir)
