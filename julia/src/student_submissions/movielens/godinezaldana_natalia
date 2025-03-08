using CSV
using DataFrames
using Printf
using Base.Threads
using Dates  # medición del tiempo

# Particionar un archivo CSV en varios archivos
function particion(archivo::String, num_archivos::Int)
    df = CSV.read(archivo, DataFrame)  # Leer el archivo completo
    n = nrow(df)  # Obtener el número de filas
    filas_por_archivo = div(n, num_archivos)  # Dividir el número total de filas

    for i in 0:num_archivos-1
        inicio = i * filas_por_archivo + 1  
        fin = (i == num_archivos - 1) ? n : (i + 1) * filas_por_archivo  # Último archivo obtiene el resto
        lineas = df[inicio:fin, :]  # Seleccionar las filas
        nombre = "ratings_" * string(i + 1) * ".csv"  # Crear el nombre del archivo
        CSV.write(nombre, lineas)  # Guardar el nuevo archivo CSV
    end

end

# Función worker
function find_ratings_worker(w::Int, kg::Vector{String}, dfm::DataFrame)
    start_time = time()  # Inicio del tiempo para el worker
    println("\nEl Worker ", lpad(w, 2), " está procesando el archivo ratings_", w, ".csv")
    rfn = "ratings_" * string(w) * ".csv" 

    # Leer el archivo de ratings CSV directamente
    ratings = CSV.read(rfn, DataFrame)

    # Crear un mapa para buscar los géneros por movieId rápidamente
    movie_genres = Dict(row.movieId => row.genres for row in eachrow(dfm))

    ng = length(kg)
    ra = zeros(ng)
    ca = zeros(ng)

    # Procesar los ratings
    for rating in eachrow(ratings)
        genres = get(movie_genres, rating.movieId, "")
        for (i, genre) in enumerate(kg)
            if occursin(genre, genres)
                ca[i] += 1
                ra[i] += rating.rating
            end
        end
    end

    elapsed_time = time() - start_time  # Tiempo transcurrido
    println("El Worker ", lpad(w, 2), " terminó en ", round(elapsed_time, digits=2), " segundos.")
    return ra, ca
end

# Función principal
function find_ratings_master()
    println("En el Master")
    start_time = time()  # Inicio del tiempo para el master
    nF = 10  # número de archivos de ratings
    kg = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
          "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
          "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"]

    ng = length(kg)
    ra = zeros(ng, nF)
    ca = zeros(ng, nF)

    # Leer el archivo de películas
    dfm = CSV.read("movies.csv", DataFrame)[:, [:movieId, :genres]]

    # Procesar cada archivo usando threads
    @threads for i in 1:nF
        ra[:, i], ca[:, i] = find_ratings_worker(i, kg, dfm)
    end

    # Sumar los resultados
    sra = zeros(ng)
    sca = zeros(ng)

    for i in 1:ng
        for j in 1:nF
            sra[i] += ra[i, j]
            sca[i] += ca[i, j]
        end
    end

    # Crear un DataFrame para los resultados
    results_df = DataFrame(Género = kg, Total_Ratings = sca, Rating_Promedio = zeros(ng))

    # Calcular y almacenar promedios en el DataFrame
    for i in 1:ng
        results_df.Rating_Promedio[i] = sca[i] > 0 ? sra[i] / sca[i] : 0.0  # Aquí se evita la división por cero
    end

    # Mostrar el DataFrame
    println("\nResultados por Género:")
    display(results_df)

    # Calcular el tiempo total del master
    elapsed_time_master = time() - start_time
    println("\nEl Master terminó en ", round(elapsed_time_master, digits=2), " segundos.")
end

# Ejecución de las funciones
particion("ratings.csv", 10)
find_ratings_master()
