using CSV
using Printf

function read_ratings_csv_file(filename::String)
    data = CSV.File(filename) |> collect
    return [(row.movieId, row.rating) for row in data]
end

function consolidate_results(genres, count_array, sum_array, start_time)
    n_genres = length(genres)
    global_counts = zeros(Int, n_genres)
    global_sums = zeros(Float64, n_genres)
    global_averages = zeros(Float64, n_genres)

    for i in 1:n_genres
        global_counts[i] = sum(count_array[i, :])
        global_sums[i] = sum(sum_array[i, :])

        # Calcular el promedio
        global_averages[i] = global_counts[i] > 0 ? global_sums[i] / global_counts[i] : 0.0
    end

    println("\nResultados finales:")
    for i in 1:n_genres
        @printf("%2d  %-20s  %10d  %12.2f\n", i-1, genres[i], global_counts[i], global_averages[i])
    end

    # Mostrar el tiempo total de ejecución
    duration = time() - start_time
    @printf("Duration = %.9fs\n", duration)
end
