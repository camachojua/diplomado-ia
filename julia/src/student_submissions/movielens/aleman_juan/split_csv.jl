using CSV
using DataFrames
using Parquet

function split_csv_to_parquet(input_csv::String, output_dir::String, num_splits::Int)
    df = CSV.read(input_csv, DataFrame)
    n = nrow(df)

    # Calcular el tamaño de cada fragmento
    chunk_size = ceil(Int, n / num_splits)

    if !isdir(output_dir)
        mkpath(output_dir)
    end

    # Dividir y guardar cada fragmento como archivo Parquet
    for i in 1:num_splits
        start_row = (i - 1) * chunk_size + 1
        end_row = min(i * chunk_size, n)

        df_chunk = df[start_row:end_row, :]

        # Guardar como archivo Parquet
        output_file = joinpath(output_dir, "ratings_part_$i.parquet")
        Parquet.write_parquet(output_file, df_chunk)
        println("Fragmento $i guardado en $output_file")
    end
end
