using Random
using Statistics
using StatsBase
using Plots
using StatsPlots
using CSV
using DataFrames


function deleteColumns(df, df_res, threshold)
    println("Borrando columnas con más de ", threshold, "% de datos faltantes...")
    
    cols = df_res[df_res.prcnt .> threshold, :]
    select!(df, Not(cols.variable))
    
    select!(df, Not(:Depth_ID))
    select!(df, Not(:Sta_ID))
    return df
end


function filterColumnsByCorrelation(df, target, threshold, relation)
    println("Borrando columnas correlacionadas ", threshold, " ", relation, " con: ", target)

    for col_name in names(df)
        
        if col_name == target
            continue
        end
        
        if col_name ∉ names(df)
            continue
        end
        
        cor_ = cor(df[:, col_name], df[:, target])
        if ((relation === false) & (abs(cor_) <= abs(threshold))) |
            ((relation === true) & (abs(cor_) > abs(threshold)))
            df = select!(df, Not(col_name))
        end
    end

    return df
end

function removeOutliersIQR(df)
    println("Borrando outliers...")
    
    for col in eachcol(df)
        q_25 = quantile(col, 0.25)
        q_75 = quantile(col, 0.75)
        iqr = q_75 - q_25
        iqr_1_5 = iqr * 1.5
        out_1 = q_25 - iqr_1_5
        out_3 = q_75 + iqr_1_5

        df = deleteat!(df, findall(<(out_1), col))
        df = deleteat!(df, findall(>(out_3), col))
        
        box_c1 = boxplot(col)
        display(box_c1)
        readline()
    end
    
    return df
end    

function deleteRow(df, col)
    println("Borrando registros con campos faltantes en ", col, "...")
    
    return dropmissing(df, col)
end

function linear_regression(x, y)
    n = length(x)
    mean_x = mean(x)
    mean_y = mean(y)
    
    slope = sum((x .- mean_x) .* (y .- mean_y)) / sum((x .- mean_x).^2)
    println("Pendiente (slope): ", slope)
    
    intercept = mean_y - slope * mean_x
    println("Intercepto (intercept): ", intercept)

    println("Mostrando regresión...")
    s = scatter(x, y, label="Datos", xlabel="x", ylabel="y", title="Regresión Lineal")
    y_pred = slope .*x .+ intercept
    p = plot!(x, y_pred, label="Modelo ajustado", lw=1)
    display(s)
    readline()
    
    return slope, intercept
end

function dataMissingPercentage(df)
    df_resumen = describe(df, :nmissing, :eltype)
    df_resumen.prcnt = (df_resumen.nmissing ./ size(df)[1]) * 100
    return df_resumen
end

function displayCorrelation(df)
    m_corr = cor(Matrix(df))
    println("Mostrando matriz de correlación...")
    hm = Plots.heatmap(m_corr, x=names(df), y=names(df))
    display(hm)
    readline()
end

function muestra_df(df)
    println("DF: ", first(df, 10))
    println("shape: " , size(df))
    readline()
end

function read_csv()
    df = CSV.read("./bottle.csv", DataFrame)    
    println(size(df))  # Data shape
    df_resumen = dataMissingPercentage(df)
    println(df_resumen)
    readline()

    df = deleteColumns(df, df_resumen, 7)
    muestra_df(df)

    df = deleteRow(df, "T_degC")
    for c_name in names(df)
        df = deleteRow(df, c_name)
    end
    muestra_df(df)

    displayCorrelation(df)

    df = filterColumnsByCorrelation(df, "T_degC", 0.1, false)
    muestra_df(df)
    
    df = filterColumnsByCorrelation(df, "T_degC", 0.8, true)
    muestra_df(df)

    displayCorrelation(df)

    df = removeOutliersIQR(df)
    muestra_df(df)
end

read_csv()
