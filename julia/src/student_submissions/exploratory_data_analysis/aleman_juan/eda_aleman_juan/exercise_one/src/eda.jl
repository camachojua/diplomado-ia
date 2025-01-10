using CSV
using DataFrames
using Statistics
using StatsBase
using Plots

# Función para calcular el porcentaje de valores faltantes por columna
# @param df: DataFrame
# @return: Diccionario con el nombre de la columna como clave y el porcentaje de valores faltantes como valor
function porcentaje_faltantes(df::DataFrame)
    porcentajes_faltantes = Dict()
    for col in names(df)
        conteo_faltantes = count(ismissing, df[:, col])
        porcentajes_faltantes[col] = (conteo_faltantes / nrow(df)) * 100
    end
    return porcentajes_faltantes
end

# Función para eliminar columnas con un porcentaje de valores faltantes mayor al umbral
# @param df: DataFrame
# @param umbral: Umbral de porcentaje de valores faltantes para eliminar columnas
# @return: DataFrame con las columnas que tienen un porcentaje de valores faltantes menor o igual al umbral
function eliminar_columnas(df::DataFrame, umbral::Float64)
    println("Eliminando columnas con porcentaje de valores faltantes mayor al umbral $umbral...")
    porcentajes_faltantes = porcentaje_faltantes(df)
    columnas_a_conservar = filter(col -> porcentajes_faltantes[col] <= umbral, names(df))
    println("Columnas conservadas: ", columnas_a_conservar)
    return df[:, columnas_a_conservar]
end

# Función para calcular la matriz de correlación
# @param df: DataFrame
# @return: Matriz de correlación y nombres de las columnas numéricas
function calcular_correlacion(df::DataFrame)
    println("Calculando matriz de correlación...")
    columnas_numericas = filter(c -> eltype(df[!, c]) <: Number, names(df))
    df_numerico = df[:, columnas_numericas]
    matriz_correlacion = pairwise(cor, eachcol(df_numerico))
    println("Matriz de correlación:")
    println(matriz_correlacion)
    return matriz_correlacion, columnas_numericas
end

# Función para mostrar la matriz de correlación como un mapa de calor
# @param matriz_correlacion: Matriz de correlación
# @param nombres_columnas: Nombres de las columnas
function mostrar_correlacion(matriz_correlacion, nombres_columnas)
    println("Mostrando mapa de calor de la matriz de correlación...")
    heatmap(
        matriz_correlacion,
        xticks=(1:length(nombres_columnas), nombres_columnas),
        yticks=(1:length(nombres_columnas), nombres_columnas),
        color=:coolwarm,
        title="Mapa de Calor de Correlación",
        xlabel="Variables",
        ylabel="Variables"
    )
end

# Función para eliminar outliers usando el IQR
# @param df: DataFrame
function eliminar_outliers_IQR!(df::DataFrame)
    println("Eliminando outliers usando IQR...")
    columnas_numericas = filter(c -> eltype(df[!, c]) <: Number, names(df))
    for col in columnas_numericas
        Q1 = quantile(skipmissing(df[!, col]), 0.25)
        Q3 = quantile(skipmissing(df[!, col]), 0.75)
        IQR = Q3 - Q1
        limite_inferior = Q1 - 1.5 * IQR
        limite_superior = Q3 + 1.5 * IQR
        df[!, col] = [x < limite_inferior || x > limite_superior ? missing : x for x in df[!, col]]
    end
    println("Outliers reemplazados con valores faltantes.")
end

# Función para eliminar filas con valores nulos en una columna específica
# @param df: DataFrame
# @param columna: Nombre de la columna
# @return: DataFrame sin filas que tienen valores nulos en la columna especificada
function eliminar_filas_con_nulos(df::DataFrame, columna::Symbol)
    println("Eliminando filas con valores nulos en la columna $columna...")
    nuevo_df = dropmissing(df, [columna])
    println("Filas restantes: ", nrow(nuevo_df))
    return nuevo_df
end

# Función para filtrar columnas según la correlación con una columna objetivo
# @param df: DataFrame
# @param objetivo: Nombre de la columna objetivo
# @param umbral: Umbral de correlación
# @param relacion: Relación de correlación ("mayor" o "menor")
# @return: DataFrame con las columnas filtradas
function filtrar_columnas_por_correlacion(df::DataFrame, objetivo::Symbol, umbral::Float64, relacion::String)
    println("Filtrando columnas según la correlación con $objetivo (umbral: $umbral, relación: $relacion)...")
    columnas_numericas = filter(c -> eltype(df[!, c]) <: Number, names(df))
    correlaciones = Dict(col => cor(skipmissing(df[!, objetivo]), skipmissing(df[!, col])) for col in columnas_numericas if col != objetivo)
    if relacion == "mayor"
        columnas_filtradas = filter(col -> correlaciones[col] > umbral, keys(correlaciones))
    elseif relacion == "menor"
        columnas_filtradas = filter(col -> correlaciones[col] < umbral, keys(correlaciones))
    else
        error("La relación debe ser 'mayor' o 'menor'.")
    end
    println("Columnas filtradas: ", columnas_filtradas)
    return df[:, vcat(objetivo, collect(columnas_filtradas))]
end

# Función para mostrar estadísticas descriptivas de cada columna
# @param df: DataFrame
function describir_datos(df::DataFrame)
    println("Estadísticas descriptivas de cada columna:")
    println(describe(df))
end

# Ejecución del análisis
println("Cargando dataset...")
df = CSV.read("julia/src/student_submissions/exploratory_data_analysis/aleman_juan/bottle_cleaned.csv", DataFrame)

# 1. Mostrar dimensiones del dataset
println("Dimensiones del dataset: ", size(df))

# 2. Mostrar tipos de datos
println("Tipos de datos de cada columna:")
for col in names(df)
    println("$col: $(eltype(df[!, col]))")
end

# 3. Contar y mostrar valores faltantes
println("Porcentaje de valores faltantes por columna:")
porcentajes_faltantes = porcentaje_faltantes(df)
for (col, perc) in porcentajes_faltantes
    println("$col: $perc%")
end

# 4. Eliminar columnas con muchos valores faltantes (ejemplo: umbral 50%)
df_limpio = eliminar_columnas(df, 50.0)

# 5. Calcular y mostrar la matriz de correlación
matriz_correlacion, nombres_columnas = calcular_correlacion(df_limpio)

# 6. Visualizar la matriz de correlación
mostrar_correlacion(matriz_correlacion, nombres_columnas)

# 7. Eliminar outliers
eliminar_outliers_IQR!(df_limpio)

# 8. Eliminar filas con valores nulos en una columna específica (ejemplo: "Temperature")
if :Temperature in names(df_limpio)
    df_limpio = eliminar_filas_con_nulos(df_limpio, :Temperature)
end

# 9. Filtrar columnas según correlación (ejemplo: correlación > 0.7 con "Salinity")
if :Salinity in names(df_limpio)
    df_limpio = filtrar_columnas_por_correlacion(df_limpio, :Salinity, 0.7, "mayor")
end

# 10. Mostrar estadísticas descriptivas
describir_datos(df_limpio)