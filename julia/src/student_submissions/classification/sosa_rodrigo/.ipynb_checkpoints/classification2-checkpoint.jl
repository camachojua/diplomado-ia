### A Pluto.jl notebook ###
# v0.19.47

using Markdown
using InteractiveUtils

# ╔═╡ 966ca760-c935-11ef-3581-97f58714ff68
begin
	using Pkg
	Pkg.add("CategoricalArrays")
	Pkg.add("MLJ")
	Pkg.add("MLJBase")
	Pkg.add("MLJLinearModels")
	Pkg.add("ROCAnalysis")
	Pkg.add("GLMNet")
	Pkg.add("MultivariateStats")
	Pkg.add("StatisticalMeasures")
end

# ╔═╡ f4c05b75-8ab2-4dd0-a50c-f509da611aff
begin

using CSV, DataFrames, Plots, MultivariateStats, MLJBase, ROCAnalysis, MLJ
using GLMNet, LinearAlgebra, StatsBase, CategoricalArrays, StatisticalMeasures

# Load dataset
data = CSV.read("Smarket.csv", DataFrame)

end


# ╔═╡ e47511d6-ea91-4d24-9543-d324e792b22d
begin
	# Function to get the shape of the dataset
	function dataShape(data)
	    rows, cols = size(data)
	    println("No. of rows: $rows \nNo. of Columns: $cols \n")
	end
	
	# Mostrar la forma del dataset (número de filas y columnas)
	dataShape(data)

	# Function to get data types of each column in a DataFrame
	function dataType(df::DataFrame)
	    results = DataFrame(column_name = names(df), data_type = eltype.(eachcol(df)))
	    return results    
	end
	
	# Mostrar los tipos de datos de cada columna
	result = dataType(data)
	println(result)
	
	# Resumen estadístico de las columnas numéricas
	#println(describe(data))
	describe(data)
end

# ╔═╡ add989fb-1e5b-4183-b343-667fbdcc5813
begin
	# Function to count the number of missing values in a given column
	function count_missing(df::DataFrame, col::String)
	    if col in names(df)
	        count = 0
	        for value in df[!, col]
	            if ismissing(value)
	                count += 1
	            end
	        end
	        #println("Column: $col, Missing Values: $count")
	        return (col, count)
	    else
	        println("Column: '$col' does not exist in the DataFrame.")
	        return nothing
	    end
	end
	
	# Function to calculate missing percentage for each column
	function dataMissingPercentage(df::DataFrame)
	    # Retrieve the data types
	    data_types = dataType(df)
	
	    # Add a column for missing percentages
	    data_types[!, :missing_percent] .= 0.0
	
	    #Calculate missing percentages for each column
	    for col in names(df)
	        _, missing_count = count_missing(df, col)
	        missing_percentage = (missing_count * 100) / nrow(df) #Calcula el porcentaje
	        row_idx = findfirst(==(col), data_types[!, :column_name]) # Encuentra el índice
	        data_types[row_idx, :missing_percent] = missing_percentage
	    end
	
	    return data_types
	end
	
	missing_percentages_df = dataMissingPercentage(data)
	println(missing_percentages_df)
end

# ╔═╡ e6a5c95b-4097-4517-8e34-79ea8d295e68
# Convertir la columna Direction a valores numéricos
data.Direction = ifelse.(data.Direction .== "Up", 1, 0)

# ╔═╡ 6b31722b-9e81-4e3b-babd-38f970c85dba
# Normalizar las variables predictoras
for col in [:Lag1, :Lag2, :Lag3, :Lag4, :Lag5, :Volume]
    data[!, col] = (data[!, col] .- mean(data[!, col])) ./ std(data[!, col])
end

# ╔═╡ 02fce5d4-dec7-4c11-9198-d73b6a8587b0
describe(data)

# ╔═╡ 6cd46408-4661-4aef-9557-cc155e956027
begin
	# Dividir en conjuntos de entrenamiento y prueba
	train_ratio = 0.8
	n_train = Int(train_ratio * nrow(data))
	train_data = data[1:n_train, :]
	test_data = data[n_train+1:end, :]
end

# ╔═╡ affc0417-db08-4a75-901f-ba0a8fbf516b
begin
	X_train = DataFrames.select(train_data, Not(:Direction)) |> Matrix
	y_train = train_data.Direction
end

# ╔═╡ 0240dcc1-0fec-470f-ad40-a4534fdf3009
begin
	X_test = DataFrames.select(test_data, Not(:Direction)) |> Matrix
	y_test = test_data.Direction
end

# ╔═╡ cbdaa39c-2cdf-4fba-a244-11de3600c945
# Función para evaluar el modelo
function evaluate_model(predictions, probabilities, y_test)

	#describe(y_test)
	#describe(probabilities)
	
    cm = StatisticalMeasures.confusion_matrix(y_test, predictions)
    println("Confusion Matrix:\n", cm)
    #acc = sum(diagm(cm)) / sum(cm)
    #println("Accuracy: ", acc)

	# Calcular precisión
    #tp = cm[1, 1]  # Verdaderos positivos
    #fp = cm[1, 2]  # Falsos positivos
    #fn = cm[2, 1]  # Falsos negativos
    #tn = cm[2, 2]  # Verdaderos negativos
    #acc = (tp + tn) / (tp + tn + fp + fn)
    #println("Accuracy: ", acc)

	# Eliminar los valores missing de y_test y probabilities
    #y_test_clean = DataFrames.remove_missing(y_test)
    #probabilities_clean = DataFrames.remove_missing(probabilities)
	
    roc_result = StatisticalMeasures.roc_curve(y_test, probabilities)
    auc_value = ROCAnalysis.auc(roc_result)
    println("AUC: ", auc_value)
    plot(roc_result, title="ROC Curve", xlabel="False Positive Rate", ylabel="True Positive Rate")
end

# ╔═╡ bad40b4a-e939-4d9a-8e13-3f559548d018
# ╠═╡ disabled = true
#=╠═╡
begin
	using Printf  # Para formatear números
	
	function evaluate_model(predictions, probabilities, y_test)
    # Convertir y_test y predictions a cadenas con formato decimal
    y_test_str = map(x -> @sprintf("%.1f", x), y_test)
    predictions_str = map(x -> @sprintf("%.1f", x), predictions)

    # Convertir a categóricos
    y_test_categorical = categorical(y_test_str, levels=["0.0", "1.0"])
    predictions_categorical = categorical(predictions_str, levels=["0.0", "1.0"])

    # Obtener matriz de confusión
    cm = confusion_matrix(y_test_categorical, predictions_categorical)
    println("Confusion Matrix:\n", cm)

    # Calcular precisión
    tp = cm[1, 1]  # Verdaderos positivos
    fp = cm[1, 2]  # Falsos positivos
    fn = cm[2, 1]  # Falsos negativos
    tn = cm[2, 2]  # Verdaderos negativos
    acc = (tp + tn) / (tp + tn + fp + fn)
    println("Accuracy: ", acc)

    # Calcular curva ROC y AUC
    y_test_binary = convert(Vector{Int}, y_test .== 1)  # Convertir a 0/1
    roc = roc_curve(y_test, probabilities)
    auc_score = auc(roc)
    println("AUC: ", auc_score)

    # Graficar curva ROC
    plot(roc, title="ROC Curve", xlabel="False Positive Rate", ylabel="True Positive Rate")
	end
end
  ╠═╡ =#

# ╔═╡ 071d7134-3e4f-4bba-b068-1beead08ec2c
begin
	# LASSO
	lasso_model = glmnet(X_train, y_train, alpha=1.0, lambda=[0.1]) # α=1.0 para LASSO
	lasso_probs = GLMNet.predict(lasso_model, X_test)[:, 1]
	lasso_preds = round.(lasso_probs)
	println("\nLASSO")
	evaluate_model(lasso_preds, lasso_probs, y_test)
end

# ╔═╡ 95e2c8d7-803e-4320-b4cf-95581f204e5f
typeof(y_test)

# ╔═╡ Cell order:
# ╠═966ca760-c935-11ef-3581-97f58714ff68
# ╠═f4c05b75-8ab2-4dd0-a50c-f509da611aff
# ╠═e47511d6-ea91-4d24-9543-d324e792b22d
# ╠═add989fb-1e5b-4183-b343-667fbdcc5813
# ╠═e6a5c95b-4097-4517-8e34-79ea8d295e68
# ╠═6b31722b-9e81-4e3b-babd-38f970c85dba
# ╠═02fce5d4-dec7-4c11-9198-d73b6a8587b0
# ╠═6cd46408-4661-4aef-9557-cc155e956027
# ╠═affc0417-db08-4a75-901f-ba0a8fbf516b
# ╠═0240dcc1-0fec-470f-ad40-a4534fdf3009
# ╠═cbdaa39c-2cdf-4fba-a244-11de3600c945
# ╠═bad40b4a-e939-4d9a-8e13-3f559548d018
# ╠═071d7134-3e4f-4bba-b068-1beead08ec2c
# ╠═95e2c8d7-803e-4320-b4cf-95581f204e5f
