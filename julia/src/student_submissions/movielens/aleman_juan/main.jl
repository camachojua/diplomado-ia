using Printf

# Cargar otros archivos
include("split_csv.jl")
include("movies.jl")
include("workers.jl")
include("utils.jl")

function main()
    println("Iniciando el proceso de división del CSV en archivos Parquet...")
    start_time = time()

    # Dividir el CSV grande en archivos Parquet
    input_csv = "ratings.csv"
    output_dir = "output_parquet"
    num_splits = 10

    split_csv_to_parquet(input_csv, output_dir, num_splits)

    println("División completada. Iniciando el procesamiento en paralelo...")

    genres = ["Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"]
    movies = load_movies("movies.csv")

    n_workers = 10
    n_genres = length(genres)
    count_array = zeros(Int, n_genres, n_workers)
    sum_array = zeros(Float64, n_genres, n_workers)

    @sync begin
        for worker_id in 1:n_workers
            @async find_ratings_worker(worker_id, genres, movies, count_array, sum_array)
        end
    end

    # Consolidar resultados y mostrar el reporte
    consolidate_results(genres, count_array, sum_array, start_time)
end

main()
