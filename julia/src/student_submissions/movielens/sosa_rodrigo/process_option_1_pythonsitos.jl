using CSV
using DataFrames
using BenchmarkTools
using Tables
using Statistics
using Query
using CUDA
using Dates
using Printf
using Parquet

# Cargar el archivo CSV de películas (movies.csv)
movies_df = CSV.read("movies.csv", DataFrame)
println("Archivo de películas cargado con éxito.")

function procesar_fragmento(worker_id, generos_conocidos, arreglo_conteo, arreglo_valor, movies_df)
    nombre_fragmento = "ratings_" * lpad(worker_id, 2, '0') * ".parquet"
    println("Intentando procesar fragmento: $nombre_fragmento")

    try
        # Cargar el archivo Parquet directamente como un DataFrame
        ratings_df = Parquet.read_parquet(nombre_fragmento) |> DataFrame
        println("Fragmento cargado exitosamente: $nombre_fragmento")
    catch e
        println("Error al cargar $nombre_fragmento: $e")
        return
    end

    try
        # Usar innerjoin en lugar de join!
        ratings_df = innerjoin(ratings_df, movies_df, on=:movieId)
        println("Join completado para fragmento: $nombre_fragmento")
    catch e
        println("Error en join de $nombre_fragmento: $e")
        return
    end

    ratings_df = select(ratings_df, [:genres, :rating])
    println("Columnas filtradas para fragmento: $nombre_fragmento")

    for row in eachrow(ratings_df)
        for i in 1:length(generos_conocidos)
            if occursin(generos_conocidos[i], row.genres)
                arreglo_conteo[i][worker_id] += 1
                arreglo_valor[i][worker_id] += row.rating
            end
        end
    end

    println("Worker $worker_id ha terminado de procesar $nombre_fragmento")
end

# Parámetros de prueba
generos_conocidos = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
                     "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
                     "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"]
numero_generos = length(generos_conocidos)

arreglo_conteo = [zeros(Int, 1) for _ in 1:numero_generos]
arreglo_valor = [zeros(Float64, 1) for _ in 1:numero_generos]

# Llamar a procesar_fragmento solo para el fragmento 1
procesar_fragmento(1, generos_conocidos, arreglo_conteo, arreglo_valor, movies_df)



# Función para consolidar resultados y procesar los archivos en paralelo
function procesar_archivo_multihilo(num_procesos)
    println("El orquestador del proceso ha iniciado su ejecución.")
    
    generos_conocidos = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
                         "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
                         "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"]
    numero_generos = length(generos_conocidos)

    arreglo_conteo = [zeros(Int, num_procesos) for _ in 1:numero_generos]
    arreglo_valor = [zeros(Float64, num_procesos) for _ in 1:numero_generos]

    Threads.@threads for i in 1:num_procesos
        procesar_fragmento(i, generos_conocidos, arreglo_conteo, arreglo_valor, movies_df)
    end

    locCount = zeros(Int, numero_generos)
    locVals = zeros(Float64, numero_generos)
    locAvg = zeros(Float64, numero_generos)  # Nuevo arreglo para el promedio

    for i in 1:numero_generos
        for j in 1:num_procesos
            locCount[i] += arreglo_conteo[i][j]
            locVals[i] += arreglo_valor[i][j]
        end
        # Calcular el promedio solo si hay calificaciones para el género
        locAvg[i] = locCount[i] > 0 ? locVals[i] / locCount[i] : 0
    end

    # Imprimir resultados finales
    for i in 1:numero_generos
        println("Género: ", generos_conocidos[i], " | Calificaciones: ", locCount[i], 
                " | Suma de Ratings: ", locVals[i], " | Promedio de Rating: ", locAvg[i])
    end

    println("Fin del orquestador.")
end

# Llamar a la función de procesamiento en paralelo
@btime procesar_archivo_multihilo(10)  # Cambia el número según el número de fragmentos