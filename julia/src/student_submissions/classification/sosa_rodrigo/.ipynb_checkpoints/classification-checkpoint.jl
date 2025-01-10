### A Pluto.jl notebook ###
# v0.19.47

using Markdown
using InteractiveUtils

# ╔═╡ fa941e30-d8e0-4615-851f-d4560686ab55
begin
	using Pkg
	Pkg.add("CategoricalArrays")
	Pkg.add("MLJ")
	Pkg.add("MLJLinearModels")
	Pkg.add("ROCAnalysis")
end

# ╔═╡ ab3acd8c-c89c-11ef-091c-a52987655355
begin

using CSV, DataFrames, Plots, CategoricalArrays, MLJ, MLJLinearModels, ROCAnalysis

# Load dataset
data = CSV.read("Smarket.csv", DataFrame)

end

# ╔═╡ beac84b1-4382-43e4-8430-4d75ebc96ead
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

# ╔═╡ 42e542ed-67fe-4e87-aef7-7c85d978d487
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

# ╔═╡ 9cb9173d-4f5c-4db5-8321-0a8c8d1b9f55
# Conversión de la columna 'Direction' a categórica
data.Direction = categorical(data.Direction)

# ╔═╡ 6477f6be-2f5b-4aef-a867-c4c29e131c51
# Convertir la columna 'Direction' a valores binarios
data.BinaryDirection = ifelse.(data.Direction .== "Up", 1.0, 0.0)

# ╔═╡ 3ec855c1-d4b2-45fb-b0fe-e8cf8c5af5dd
describe(data)

# ╔═╡ 44062985-436e-44bf-ba7e-60dc308937b1
begin
	# Dividir los datos en entrenamiento y prueba
	train_indices = data.Year .< 2005
	test_indices = data.Year .== 2005
end

# ╔═╡ 8936c5d4-c8b1-4b08-b758-ab1334de3fd7
begin
	train_data = data[train_indices, :]
	test_data = data[test_indices, :]
end

# ╔═╡ f3d29670-3b1d-4cd1-a47d-4b48019911b8
begin
	# Definir las variables predictoras y objetivo
	X_train = select(train_data, Not([:Direction, :BinaryDirection]))
	y_train = categorical(train_data.BinaryDirection)
end

# ╔═╡ 248c5997-440a-499e-8b4d-2913ed04eae9
begin
	X_test = select(test_data, Not([:Direction, :BinaryDirection]))
	y_test = categorical(test_data.BinaryDirection)
end

# ╔═╡ f5b1b775-2f39-4fc7-b00a-d8f1339009ad
md"""
Datos preparados para entrenamiento y prueba
"""

# ╔═╡ 6c19a8f0-a7e5-4dd7-a8ec-36381f2c4dfd
md" # Implementación de LASSO"

# ╔═╡ 7cb936ab-e628-46ed-a1dd-359fdef977a6
begin
	lasso_model = @load LogisticClassifier pkg=MLJLinearModels
	model = lasso_model(penalty=:l1, lambda=0.1) # LASSO usa penalización L1 
end

# ╔═╡ a93a79d8-c5bf-4c2a-91b1-8c4a36006c54
begin
	mach = machine(model, X_train, y_train)
	fit!(mach)
end

# ╔═╡ 2f2046f4-592b-4ba7-a2dd-56bc1f028f23
begin
	# Predicciones y evaluación
	probabilities = predict(mach, X_test)
	probabilities = pdf.(probabilities, 1.0)  # Convertir a escala numérica
end

# ╔═╡ a3bf0e6d-85c6-4dc4-9d42-264ac8fff09b
# Convertir y_test a categórico explícito con niveles definidos
y_test_categorical = categorical(y_test, levels=[0.0, 1.0])

# ╔═╡ 9b046544-3940-456f-9bf5-9963e8aea83f
# Generar etiquetas predichas como categóricas binarias
predicted_labels = categorical(probabilities .> 0.5, levels=[0.0, 1.0])

# ╔═╡ 91435e78-24c5-4c56-aa9c-0d3a8bb5ffa7
begin
	#predictions = probabilities .> 0.5  # Convertir probabilidades en clases binarias
	conf_mat = confusion_matrix(y_test, predicted_labels)
	println("Matriz de confusión: ", conf_mat)
end

# ╔═╡ fcb93b3f-1911-488f-a94c-0b9445dc55fe
begin
	#scores = predict(mach, X_test)
	#y_test_categorical = categorical(y_test, levels=[0.0, 1.0]) # Convertir y_test a categórico
	roc = roc_curve(y_test_categorical, predicted_labels)
	plot(roc, xlabel="False Positive Rate", ylabel="True Positive Rate", title="ROC Curve - LASSO")
end

# ╔═╡ Cell order:
# ╠═fa941e30-d8e0-4615-851f-d4560686ab55
# ╠═ab3acd8c-c89c-11ef-091c-a52987655355
# ╠═beac84b1-4382-43e4-8430-4d75ebc96ead
# ╠═42e542ed-67fe-4e87-aef7-7c85d978d487
# ╠═9cb9173d-4f5c-4db5-8321-0a8c8d1b9f55
# ╠═6477f6be-2f5b-4aef-a867-c4c29e131c51
# ╠═3ec855c1-d4b2-45fb-b0fe-e8cf8c5af5dd
# ╠═44062985-436e-44bf-ba7e-60dc308937b1
# ╠═8936c5d4-c8b1-4b08-b758-ab1334de3fd7
# ╠═f3d29670-3b1d-4cd1-a47d-4b48019911b8
# ╠═248c5997-440a-499e-8b4d-2913ed04eae9
# ╟─f5b1b775-2f39-4fc7-b00a-d8f1339009ad
# ╟─6c19a8f0-a7e5-4dd7-a8ec-36381f2c4dfd
# ╠═7cb936ab-e628-46ed-a1dd-359fdef977a6
# ╠═a93a79d8-c5bf-4c2a-91b1-8c4a36006c54
# ╠═2f2046f4-592b-4ba7-a2dd-56bc1f028f23
# ╠═a3bf0e6d-85c6-4dc4-9d42-264ac8fff09b
# ╠═9b046544-3940-456f-9bf5-9963e8aea83f
# ╠═91435e78-24c5-4c56-aa9c-0d3a8bb5ffa7
# ╠═fcb93b3f-1911-488f-a94c-0b9445dc55fe
