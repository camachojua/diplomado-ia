
using CSV
using DataFrames
using Base.Threads
using Printf

# Cargar el archivo `movies.csv` y mapear `movieId` a `genres`
movies = CSV.File("movies.csv") |> DataFrame
movies_dict = Dict(movies.movieId .=> movies.genres)

# Función para contar y sumar calificaciones por género
function count_and_sum_ratings_by_genre(ratings, movies_dict)
    genre_count = Dict{String, Int}()
    genre_sum = Dict{String, Float64}()
    
    for row in eachrow(ratings)
        genres = get(movies_dict, row.movieId, "")
        if !isempty(genres)
            for genre in split(genres, "|")
                genre_count[genre] = get(genre_count, genre, 0) + 1
                genre_sum[genre] = get(genre_sum, genre, 0.0) + row.rating
            end
        end
    end
    return genre_count, genre_sum
end

# Función para cargar y procesar cada archivo CSV en paralelo
function load_count_and_sum(i, movies_dict)
    ratings = CSV.read("ratings_part_$i.csv", DataFrame; select = [:movieId, :rating])
    count_and_sum_ratings_by_genre(ratings, movies_dict)
end

# Arreglos para almacenar resultados de cada hilo
results_count = Vector{Dict{String, Int}}(undef, 10)
results_sum = Vector{Dict{String, Float64}}(undef, 10)

# Procesamiento en paralelo con Threads
Threads.@threads for i in 1:10
    results_count[i], results_sum[i] = load_count_and_sum(i, movies_dict)
end

# Combinar los resultados de todos los hilos
final_count = Dict{String, Int}()
final_sum = Dict{String, Float64}()

for i in 1:10
    for (genre, count) in results_count[i]
        final_count[genre] = get(final_count, genre, 0) + count
    end
    for (genre, sum_rating) in results_sum[i]
        final_sum[genre] = get(final_sum, genre, 0.0) + sum_rating
    end
end

# Calcular el promedio y almacenar los resultados en un DataFrame
results_df = DataFrame(Genre=String[], Count=Int[], AverageRating=Float64[])

for (genre, count) in final_count
    avg_rating = final_sum[genre] / count
    push!(results_df, (Genre=genre, Count=count, AverageRating=avg_rating))
end

# Guardar los resultados en un archivo CSV
CSV.write("genre_ratings_summary.csv", results_df)

println("Archivo 'genre_ratings_summary.csv' creado con éxito.")

