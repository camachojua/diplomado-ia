# Autor: Martha Rico Diener
# Equipo: Los cantantes
# Este programa procesa archivos CSV que contienen calificaciones de películas y 
# estadísticas sobre géneros. Utiliza concurrencia para manejar múltiples archivos 
# de calificaciones y calcula el promedio de calificaciones por género.

# Instalar los paquetes necesarios si no están instalados
using Pkg
# Verifica si los paquetes "CSV" y "DataFrames" están instalados
if !haskey(Pkg.installed(), "CSV")
    Pkg.add("CSV")  # Agrega el paquete si no está instalado
else
    println("El paquete 'CSV' ya está instalado.")
end
if !haskey(Pkg.installed(), "DataFrames")
    Pkg.add("DataFrames") # Agrega el paquete si no está instalado
else
    println("El paquete 'DataFrames' ya está instalado.")
end
using Pkg
Pkg.add("CSV")
Pkg.add("DataFrames")
using CSV
using DataFrames
using Base.Threads
using Printf

# Tiempo de inicio
t1 = time() #Lo utilizo para calcular el tiempo total de ejecución.

# Cargar el archivo CSV en un DataFrame
df = CSV.File("ratings.csv") |> DataFrame

# Definir el número de partes en las que quieres dividir el archivo. 
# Hacemos esto para poder trabajar ls archivos resultantes en paralelo. 
N = 10

# Calcular el tamaño de cada parte
num_rows = nrow(df)
rows_per_part = div(num_rows, N) + (num_rows % N != 0 ? 1 : 0)  # Asegura que todas las filas se procesen

# Crear un loop para dividir el DataFrame y guardar cada parte
for i in 0:(N-1)
    start_row = i * rows_per_part + 1
    end_row = min((i == N - 1) ? num_rows : (start_row + rows_per_part - 1), num_rows)

    if start_row > num_rows
        break  # Salir si el índice inicial es mayor que el número de filas
    end

    # Crear un DataFrame para la parte actual
    part_df = df[start_row:end_row, :]

    # Guardar la parte como un nuevo archivo CSV
    CSV.write("ratings_part_$(i + 1).csv", part_df)
end

println("El archivo ha sido dividido en $N partes.")
elapsed_time = time() - t1
println("Tiempo transcurrido en dividir: ", elapsed_time, " seconds")

# Asegúrate de que el número de hilos esté configurado
println("Número de hilos disponibles: ", nthreads())

# Función para cargar los archivos CSV y contar calificaciones por género
function cargar_y_contar(i, movies)
    ratings = CSV.read("ratings_part_" * string(i) * ".csv", DataFrame)
    return count_ratings_by_genre(ratings, movies)
end

# Función para contar las calificaciones por género y calcular el promedio
function count_ratings_by_genre(ratings, movies)
    # Unir las tablas por 'movieId'
    data = innerjoin(ratings, movies, on=:movieId)
   
    # Inicializar un diccionario para acumular calificaciones y contar por género
    genre_data = Dict{String, Tuple{Float64, Int}}()  # (suma de calificaciones, conteo)
   
    for row in eachrow(data)
        genres = split(row.genres, "|")  # Dividir los géneros por '|'
        rating = row.rating  # Suponiendo que la calificación está en la columna 'rating'
        
        for genre in genres
            if haskey(genre_data, genre)
                genre_data[genre] = (genre_data[genre][1] + rating, genre_data[genre][2] + 1)
            else
                genre_data[genre] = (rating, 1)
            end
        end
    end

    return genre_data  # Cambiado para devolver el diccionario completo
end

# Cargar la lista de películas
movies = CSV.read("movies.csv", DataFrame)

# Arreglo para almacenar los resultados de cada hilo
results = Vector{Dict{String, Tuple{Float64, Int}}}(undef, N)

# Usar Threads.@threads para ejecutar en paralelo
Threads.@threads for i in 1:N
    results[i] = cargar_y_contar(i, movies)
end

# Combinar los resultados de todos los hilos
final_counts = Dict{String, Tuple{Float64, Int}}()  # (suma de calificaciones, conteo total)

for result in results
    for (genre, (sum_ratings, count)) in result
        if haskey(final_counts, genre)
            final_counts[genre] = (final_counts[genre][1] + sum_ratings, final_counts[genre][2] + count)
        else
            final_counts[genre] = (sum_ratings, count)
        end
    end
end
using Printf

println("Género\t\t\tConteo\tPromedio de calificaciones")
println("----------------------------------------------------------")

# Función para dividir los núemros grandes con comas
function format_with_commas(n::Int)
    s = string(n)
    len = length(s)

    if len <= 3
        return s
    end

    result = ""
    for (i, digit) in enumerate(reverse(s))
        if i > 1 && (i - 1) % 3 == 0
            result *= ","
        end
        result *= digit
    end

    return reverse(result)
end

# Asegurarse de que final_counts esté definido y no esté vacío
if isempty(final_counts)
    println("No se encontraron géneros.")
else
    sorted_genres = sort(collect(keys(final_counts)))

    for genre in sorted_genres
        (sum_ratings, count) = final_counts[genre]
        average = count > 0 ? sum_ratings / count : 0.0  # Evitar división por cero
        formated_count = format_with_commas(count)
       # println(@sprintf("%-20s\t%s\t\t%.2f", genre, formated_count, average))
       println(@sprintf("%-20s\t%20s\t%10.2f", genre, formated_count, average))
    end
end

final_time = time() - t1
println("Tiempo transcurrido total: ", final_time, " seconds")
count_time = final_time - elapsed_time
println("Tiempo transcurrido en contar y promediar: ", count_time, " seconds")

