using CSV, DataFrames, Statistics, StatsBase, CairoMakie

#Esta función devuelve el número de filas y columnas de un DataFrame
function dataShape(data::DataFrame)
    rows = nrow(data)
    cols = ncol(data)
    return rows, cols
end

#Esta función devuelve un DataFrame que contiene el nombre de las variables de un determinado DataFrame y el tipo de datos que contiene. 
function dataType(data::DataFrame)
    types = eltype.(eachcol(data))
    cols = names(data)
    return DataFrame(Column = cols, Data_type = types)
end

#Esta función cuenta el número de registros faltantes (missing) para una determinada columna de datos
function count_missing(col_data)
    return count(ismissing, col_data)
end

#Esta función determina el porcentaje de datos faltantes dado un conjunto de datos
function dataMissingPercentage(col_data)
    return count_missing(col_data) / length(col_data) * 100
end

#Esta función elimina las columnas de un DataFrame de acuerdo a un valor umbral del porcentaje de datos faltantes
function deleteColumns(data::DataFrame, threshold)
    cols_deleted = filter(col_name -> dataMissingPercentage(data[!, col_name]) >= threshold, names(data))
    filter_data = select(data, Not(cols_deleted))
    println("Se eliminaron $(length(cols_deleted)) columnas: $(cols_deleted)")
    return filter_data
end

#Esta función elimina las filas de un DataFrame en las cuales existe un valor ausente
function deleteRows(data::DataFrame)
    return dropmissing(data)
end

#Esta función elimina los valores atípicos de un conjunto de datos usando el rango intercuartílico
function removeOutliersIQR(col_name, data::DataFrame)
    len1, = dataShape(data)
    iq = iqr(data[!,col_name])
    Q1 = quantile(data[!,col_name], 0.25)
    Q3 = quantile(data[!,col_name], 0.75)
    upper_lim = Q3 + 1.5*iq 
    lower_lim = Q1 - 1.5*iq
    filter!(col_name => x -> lower_lim < x < upper_lim, data)
    len2, = dataShape(data)
    println("Se removieron $(len1-len2) outliers de la columna $(col_name)") 
end

#Esta función crea la matriz de correlación de los datos de un DataFrame
function calculateCorrelation(data::DataFrame)
    return cor(Matrix(data))
end

#Esta función muestra gráficamente (heatmap) una matriz de correlación
function displayCorrelation(data::DataFrame)
    cor_matrix = calculateCorrelation(data)
    var_names = names(data)
    fig = Figure(size=(1000, 800))
    ax = Axis(
        fig[1,1], 
        title = "Matríz de correlación", 
        xticks = (1:length(var_names), var_names), 
        yticks = (1:length(var_names), var_names),
        xticklabelrotation = 45,
        aspect = DataAspect()
    )
    pltobj = heatmap!(cor_matrix; colorrange=(-1, 1), colormap=Reverse(:viridis))
    Colorbar(fig[1, 2], pltobj)
    colsize!(fig.layout, 1, Aspect(1, 1.0))

    # Agregar valores de la correlación en cada celda
    for i in 1:size(cor_matrix, 1)
        for j in 1:size(cor_matrix, 2)
            text!(
                ax,
                j, i,  # Coordenadas (columna, fila)
                text = string(round(cor_matrix[i, j], digits=2)),  
                align = (:center, :center), 
                fontsize = 12,               
                color = :black               
            )
        end
    end
    
    return fig
end