#Autor Martha Rico Diener
#Equipo Los cantantes
#Este programa procesa archivos CSV que contienen calificaciones de películas y 
#estadísticas sobre géneros. Utiliza concurrencia para manejar múltiples archivos 
#de calificaciones y calcula el promedio de calificaciones por género.

using Pkg
Pkg.add("CSV")
Pkg.add("DataFrames")
using CSV
using DataFrames
using Base.Threads
using Printf

# Tiempo de inicio
t1 = time()

# Cargar el archivo CSV
df = CSV.File("ratings.csv") |> DataFrame

# Definir el número de partes en las que quieres dividir el archivo
N = 10

# Calcular el tamaño de cada parte
num_rows = nrow(df)
rows_per_part = div(num_rows, N)

# Crear un loop para dividir el DataFrame y guardar cada parte
for i in 0:(N-1)
    start_row = i * rows_per_part + 1
    end_row = (i == N - 1) ? num_rows : (start_row + rows_per_part - 1)

    # Crear un DataFrame para la parte actual
    part_df = df[start_row:end_row, :]

    # Guardar la parte como un nuevo archivo CSV
    CSV.write("ratings_part_$(i+1).csv", part_df)
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

# Función para contar las calificaciones por género
function count_ratings_by_genre(ratings, movies)
    # Unir las tablas por 'movieId'
    data = innerjoin(ratings, movies, on=:movieId)
   
    # Inicializar un diccionario para contar las calificaciones por género
    genre_count = Dict{String, Int}()
   
    for row in eachrow(data)
        genres = split(row.genres, "|")  # Dividir los géneros por '|'
        for genre in genres
            genre_count[genre] = get(genre_count, genre, 0) + 1
        end
    end
    return genre_count
end

# Cargar la lista de películas
movies = CSV.read("movies.csv", DataFrame)

# Arreglo para almacenar los resultados de cada hilo
results = Vector{Dict{String, Int}}(undef, N)

# Usar Threads.@threads para ejecutar en paralelo
Threads.@threads for i in 1:N
    results[i] = cargar_y_contar(i, movies)
end

# Combinar los resultados de todos los hilos
final_result = Dict{String, Int}()
for result in results
    for (genre, count) in result
        final_result[genre] = get(final_result, genre, 0) + count
    end
end

function format_with_commas(n::Int)
    s = string(n)
    len = length(s)

    # Si la longitud es menor o igual a 3, solo retorna el número
    if len <= 3
        return s
    end

    # Crear una cadena con comas
    result = ""
    for (i, digit) in enumerate(reverse(s))
        if i > 1 && (i - 1) % 3 == 0
            result *= ","
        end
        result *= digit
    end

    return reverse(result)
end

# Iterar sobre el diccionario e imprimir el resultado formateado
for (genre, count) in final_result
    formated_count = format_with_commas(count)
    println("$genre: \t $formated_count calificaciones")
end
