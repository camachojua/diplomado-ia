using CSV
using DataFrames
using Base.Threads

function dividir_csv(file_path::String, n::Int)
    # Leer el archivo CSV completo
    data = CSV.File(file_path) |> DataFrame
    
    # Calcular el tamaño de cada parte
    total_rows = size(data, 1)
    rows_per_part = ceil(Int, total_rows / n)
    
    # Función para guardar una parte del DataFrame en un archivo CSV
    function guardar_parte(data_part::DataFrame, index::Int)
        part_file = "ratings_parte_$index.csv"
        CSV.write(part_file, data_part)
        println("Parte $index guardada como $part_file")
    end

    # Crear tareas concurrentes para dividir y guardar cada parte
    @threads for i in 1:n
        start_idx = (i - 1) * rows_per_part + 1
        end_idx = min(i * rows_per_part, total_rows)
        data_part = data[start_idx:end_idx, :]
        guardar_parte(data_part, i)
    end
end

dividir_csv("ratings_large.csv", 10)