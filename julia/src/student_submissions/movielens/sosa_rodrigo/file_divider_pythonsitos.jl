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

function dividir_archivo_csv_linea_por_linea(archivo, filas_por_fragmento)
    start_time = now()  # Tiempo de inicio

    # Inicializar variables para el procesamiento
    i = 1
    contador_filas = 0
    df_fragmento = DataFrame()

    # Abrir el archivo CSV y leer línea por línea
    for row in CSV.File(archivo)
        push!(df_fragmento, row)
        contador_filas += 1

        # Guardar el fragmento cuando se alcanza el número de filas por fragmento
        if contador_filas >= filas_por_fragmento
            nombre_parquet = "ratings_" * lpad(i, 2, '0') * ".parquet"
            Parquet.write_parquet(nombre_parquet, df_fragmento)
            println("Fragmento guardado: $nombre_parquet")

            # Reiniciar el DataFrame y el contador de filas
            df_fragmento = DataFrame()
            contador_filas = 0
            i += 1
        end
    end

    # Guardar cualquier fila restante en un último fragmento
    if nrow(df_fragmento) > 0
        nombre_parquet = "ratings_" * lpad(i, 2, '0') * ".parquet"
        Parquet.write_parquet(nombre_parquet, df_fragmento)
        println("Último fragmento guardado: $nombre_parquet")
    end

    end_time = now()  # Tiempo de finalización
    println("Tiempo total para dividir el archivo: ", end_time - start_time)
end

# Llamar a la función con el nombre del archivo y el tamaño del fragmento
dividir_archivo_csv_linea_por_linea("ratings.csv", 2500010)