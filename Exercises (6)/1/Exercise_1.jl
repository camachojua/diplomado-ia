using CSV, DataFrames, Statistics, StatsPlots

function save_report(file_path::String, content::String)
    open(file_path, "w") do f
        write(f, content)
    end
end

# Read
function read_csv(file_path::String)
    return CSV.read(file_path, DataFrame)
end

# Ensure data is numeric where applicable, handling missing values properly
function ensure_numeric(data::DataFrame)
    println("Tipos iniciales de las columnas:")
    for col in names(data)
        println("$col: $(eltype(data[!, col]))")
    end

    for col in names(data)
        # Si la columna ya es numérica, pero tiene valores missing, no hacemos conversión
        if eltype(data[!, col]) <: Union{Missing, Number}
            continue
        elseif eltype(data[!, col]) <: AbstractString
            try
                data[!, col] = parse.(Float64, replace.(data[!, col], missing => "NaN"))
                println("Columna $col convertida a Float64.")
            catch e
                println("No se pudo convertir la columna $col a Float64: $e")
            end
        else
            println("Columna $col no es de tipo compatible para conversión.")
        end
    end

    println("Tipos finales de las columnas:")
    for col in names(data)
        println("$col: $(eltype(data[!, col]))")
    end

    return data
end


# Calculate missing ones
function missing_percentage(data::DataFrame)
    total_rows = nrow(data)
    return Dict(col => count(ismissing, data[!, col]) / total_rows * 100 for col in names(data))
end

# Eliminate columns with missing data above threshold
function deleteColumns(data::DataFrame, threshold::Float64)
    missing_perc = missing_percentage(data)
    keep_cols = [col for col in names(data) if missing_perc[col] <= threshold]
    deleted_cols = setdiff(names(data), keep_cols)  # Capture eliminated columns 
    return data[:, keep_cols], deleted_cols
end

# Correlation matrix
function cal_correlation(data::DataFrame)
    numeric_cols = names(data)[map(c -> eltype(data[!, c]) <: Number, names(data))]
    numeric_data = data[:, numeric_cols]
    return cor(Matrix(numeric_data))
end

# Show correlation matrix using heatmap and save image ...
function display_correlation(data::DataFrame, img_path::String)
    corr_matrix = cal_correlation(data)
    heatmap(corr_matrix, title="Matriz de correlación", xlabel="Columnas", ylabel="Columnas")
    savefig(img_path)  # Save image in file
    display(heatmap(corr_matrix))  # Show heatmap
end

# Eliminate outliers with interquartile range
function remove_outliers_IQR(data::DataFrame)
    numeric_cols = names(data)[map(c -> eltype(data[!, c]) <: Number, names(data))]
    original_rows = nrow(data)  # Num of rows before removing outliers
    for col in numeric_cols
        q1, q3 = quantile(data[!, col], [0.25, 0.75])
        iqr = q3 - q1
        lower_bound, upper_bound = q1 - 1.5 * iqr, q3 + 1.5 * iqr
        data = filter(row -> (row[col] ≥ lower_bound) && (row[col] ≤ upper_bound), data)
    end
    deleted_rows = original_rows - nrow(data)  # Num of rows deleted
    println("Se eliminaraon $deleted_rows filas por outliers con el IQR.")
    return data, deleted_rows
end

# Describe data set
function describe_data(data::DataFrame)
    return DataFrames.describe(data)
end

# Principal processing function
function process_csv(file_path::String, missing_threshold::Float64)
    data = read_csv(file_path)
    println("Archivo leído con ", nrow(data), " filas y ", ncol(data), " columnas.")
    
    data = ensure_numeric(data)  # Asegurar datos numéricos
    println("Tipos de datos después de asegurar numéricos: ", eltype.(eachcol(data)))
    
    println("Porcentaje de datos faltantes por columna:")
    missing_percentages = missing_percentage(data)
    println(missing_percentages)
    
    data, deleted_cols = deleteColumns(data, missing_threshold)  # Eliminate missing data
    println("Columnas eliminadas: ", deleted_cols)
    
    return data, missing_percentages, deleted_cols
end

file_path = "/Users/michelletorres/Desktop/Homeworks AI/archive/bottle.csv"
missing_threshold = 10.0  # Max percentage

# Process data
processed_data, missing_percentages, deleted_cols = process_csv(file_path, missing_threshold)

# Data set description
description = describe_data(processed_data)

# Eliminate outliers
processed_data, deleted_rows = remove_outliers_IQR(processed_data)

# Save img heatmap
heatmap_img_path = "/Users/michelletorres/Desktop/heatmap.png"
display_correlation(processed_data, heatmap_img_path)

# Report
report_content = """
Informe de Análisis Exploratorio de Datos (EDA)

1. Información general:
El archivo contiene $(nrow(processed_data)) filas y $(ncol(processed_data)) columnas.
   
2. Porcentaje de datos faltantes por columna:
   $(missing_percentages)

3. Columnas eliminadas debido a datos faltantes:
   $(deleted_cols)

4. Descripción de los datos:
   $(description)

5. Outliers eliminados con el IQR:
$deleted_rows

6. Matriz de correlación:
$heatmap_img_path
   
7. Conclusión:
Se completó el análisis de manera satisfactoria, limpiando los datos faltantes, eliminando outliers y mostrando la matriz de correlación entre las variables.

"""

# Save it 
report_path = "/Users/michelletorres/Desktop/Homeworks AI/EDA_report.txt"
save_report(report_path, report_content)



