using CSV, DataFrames

# Función para cargar el archivo de películas
function load_movies(filename::String)
    df = CSV.read(filename, DataFrame)
    movies = Dict(row.movieId => row.genres for row in eachrow(df))
    return movies
end
