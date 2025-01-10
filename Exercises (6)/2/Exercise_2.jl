using CSV, DataFrames, Statistics, GLM, StatsPlots, StatsModels, StatsBase

file_path = "/Users/michelletorres/Desktop/Homeworks AI/archive/bottle.csv"
data = CSV.read(file_path, DataFrame)

names(data) .= strip.(names(data)) # Remove extra spaces from column names
println("Columnas disponibles:") 
for col in names(data)
    println("`$col`")
end

columns_of_interest = [:T_degC, :Salnty, :Depthm, :O2ml_L] # Necessary columns as symbols
missing_columns = setdiff(columns_of_interest, Symbol.(names(data))) # Check if required columns are present (convert DataFrame columns to symbols)
if !isempty(missing_columns)
    println("Faltan las siguientes columnas en DataFrame: $missing_columns")
    error("Faltan columnas necesarias")
end

filtered_data = data[:, columns_of_interest] # Filter columns
filtered_data = dropmissing(filtered_data) # Ensure there are no missing values in the selected columns
println("Datos después de filtrado:") # Verify that the columns have been loaded correctly

data_model = @formula(T_degC ~ Salnty + Depthm + O2ml_L) # Linear regression with GLM
lm_model = lm(data_model, filtered_data)
println("Resumen del modelo")
println(coef(lm_model))
println(summary(lm_model))

function calculate_rmse(model, data)
    predictions = StatsBase.predict(model, data)  # Utiliza StatsBase.predict
    residuals = data[:, :T_degC] .- predictions
    return sqrt(mean(residuals .^ 2))
end

rmse = calculate_rmse(lm_model, filtered_data)
println("RMSE del modelo: $rmse")

histogram(filtered_data[!, :T_degC], title="Distribución T_degC", xlabel="T_degC", ylabel="Frecuencia") # For each variable
histogram(filtered_data[!, :Salnty], title="Distribución Salnty", xlabel="Salnty", ylabel="Frecuencia")
histogram(filtered_data[!, :Depthm], title=" Distribución Depthm", xlabel="Depthm", ylabel="Frecuencia")
histogram(filtered_data[!, :O2ml_L], title="Distribución O2ml_L", xlabel="O2ml_L", ylabel="Frecuencia")

#############################
combinations = [ # List of independent variable combinations
    [:Salnty, :Depthm, :O2ml_L],
    [:Salnty, :Depthm],
    [:Salnty, :O2ml_L],
    [:Depthm, :O2ml_L],
    [:Salnty],
    [:Depthm],
    [:O2ml_L]
]

best_rmse = Inf
best_model = nothing
best_combination = nothing
##################################
names(filtered_data)
println(names(filtered_data))  # Verify column names
#################################
for combination in combinations
    formula = @eval @formula(T_degC ~ $(Expr(:call, :+, combination...)))     # @eval to correctly construct the formula
    lm_model = lm(formula, filtered_data)
    rmse = calculate_rmse(lm_model, filtered_data)    # Calculate RMSE for the current combination
    println("RMSE para combinación $combination: $rmse")
    if rmse < best_rmse     # If RMSE is better, save the model
        best_rmse = rmse
        best_model = lm_model
        best_combination = combination
    end
end

println("Mejor variable para combinación: $best_combination con RMSE: $best_rmse")
correlation_matrix = cor(Matrix(filtered_data[:, [:T_degC, :Salnty, :Depthm, :O2ml_L]])) # Correlation between variables
heatmap(correlation_matrix, xlabel="Variables", ylabel="Variables", title="Matriz de correlación")
