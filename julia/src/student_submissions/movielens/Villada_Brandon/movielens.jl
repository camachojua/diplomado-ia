using DataFrames, CSV
using Base.Threads: @threads

println("Numero de hilos:",Threads.nthreads())

#Funcion separar archivos
function filesharding()
    archivo = "ratings.csv"
num_procesos = 10

# Función para contar las líneas del archivo
function contar_lineas(archivo)
    count = 0
    open(archivo, "r") do f
        for _ in eachline(f)
            count += 1
        end
    end
    print("Numero de lineas en el archivo: ", count)
    return count
end

# Leer el encabezado del archivo
function leer_header(archivo)
    open(archivo, "r") do f
        return readline(f)  # Leer solo la primera línea como encabezado
    end
end

# Parámetros de partición
size_of_file = contar_lineas(archivo) - 1  # Excluyendo el encabezado
number_of_chunks = ceil(Int, size_of_file / num_procesos)
header = leer_header(archivo)

# Función para generar archivos pequeños con encabezado
function generate_small_file(archivo, i, header, num_procesos, number_of_chunks)
    start_line = (i - 1) * number_of_chunks + 2  # Saltamos el encabezado
    end_line = min(i * number_of_chunks + 1, size_of_file + 1)  # Incluir límites
    output_file = "ratings_$i.csv"
    
    open(archivo, "r") do f
        open(output_file, "w") do out
            println(out, header)  # Escribir el encabezado
            for (line_num, line) in enumerate(eachline(f))
                if line_num >= start_line && line_num <= end_line
                    println(out, line)
                end
            end
        end
    end
end

# Medir el tiempo de ejecución
start_time = time()

@threads for i in 1:num_procesos
    generate_small_file(archivo, i, header, num_procesos, number_of_chunks)
end

end_time = time()

println("Tiempo transcurrido separando archivos: $(end_time - start_time) segundos")
    
end

filesharding()

# Leer los datos de películas
movies = DataFrame(CSV.File("movies.csv"))

# Lista de géneros conocidos
kg = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
      "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
      "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"]

ng = length(kg) # Número de géneros conocidos

# Arrays para almacenar calificaciones y conteos para cada género
ra = zeros(Float64, ng)  # Sumas de calificaciones
ca = zeros(Int, ng)       # Contadores

# Función para actualizar los contadores y sumas
function update_genre_stats(ratings::DataFrame, ca, ra, i)
    # Medir el tiempo de inicio
    start_time = time()
    println("Worker $i started")
    for row in eachrow(ratings)
        genres = split(row.genres, "|")
        rating = row.rating

        for (ig, genre) in enumerate(kg)
            if genre in genres
                @lock lock ca[ig] += 1              # Incrementar contador para el género
                @lock lock ra[ig] += rating         # Sumar la calificación
            end
        end
    end
    # Medir el tiempo de finalización
    end_time = time()
    println("Duration = $(end_time - start_time)s")
    println("Worker $i completed")
end


# Medir el tiempo de inicio
start_time = time()
lock = ReentrantLock()
# Procesar archivos de ratings del 1 al 10
@threads for i in 1:10
    filename = "ratings_$i.csv"
    ratings = DataFrame(CSV.File(filename))

    # Hacer un inner join entre películas y calificaciones
    ratings = innerjoin(movies, ratings, on = :movieId)

    # Filtrar solo las columnas necesarias
    select!(ratings, [:genres, :rating])

    # Procesar las calificaciones

    update_genre_stats(ratings, ca, ra, i)  # Usar solo un hilo para simplificar
end



# Imprimir resultados
for ig in 1:ng
    println("$ig    $(kg[ig])   $(ca[ig])")
end

# Calcular promedios
println()
println("ID    Genre    Avg Rating")

for ig in 1:ng
    if ca[ig] > 0
        avg_rating = ra[ig] / ca[ig]   # Calcular promedio
        println("$ig    $(kg[ig])   $(avg_rating)")
    end
end

# Medir el tiempo de finalización
end_time = time()
println("Duration = $(end_time - start_time)s")
