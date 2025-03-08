# Importar los paquetes necesarios
using CSV, DataFrames, Lasso, GLMNet, DecisionTree, NearestNeighbors, LIBSVM, MLBase, Plots, Random, Statistics

# Paso 1: Cargar y preprocesar los datos
println("Cargando y limpiando los datos...")
data = CSV.read("/home/nat/Ejercicios diplomado/exercise_3/dat/Smarket.csv", DataFrame)
data = dropmissing(data)  # Eliminar filas con valores faltantes
data.Direction = ifelse.(data.Direction .== "Up", 1.0, 0.0)  # Codificar la variable objetivo como binaria (1 para "Up", 0 para "Down")

# Normalizar las columnas numéricas (excluir la columna 'Direction')
for col in names(data)
    if eltype(data[!, col]) <: Number && col != "Direction"
        data[!, col] .= (data[!, col] .- mean(data[!, col])) ./ std(data[!, col])
    end
end

# Dividir los datos en conjuntos de entrenamiento y prueba
function train_test_split(data, ratio=0.8)
    n = size(data, 1)
    indices = shuffle(1:n)
    train_idx = indices[1:floor(Int, n * ratio)]
    test_idx = indices[floor(Int, n * ratio) + 1:end]
    return data[train_idx, :], data[test_idx, :]
end

train_data, test_data = train_test_split(data)
X_train, y_train = select(train_data, Not(:Direction)) |> Matrix, train_data.Direction
X_test, y_test = select(test_data, Not(:Direction)) |> Matrix, test_data.Direction

# Paso 2: Entrenar modelos
# LASSO
println("Entrenando modelo LASSO...")
lasso_model = fit(LassoPath, X_train, y_train; standardize=true)
lasso_predictions = Lasso.predict(lasso_model, X_test)[:, end]  # Modelo final con lambda

# Ridge Regression
println("Entrenando modelo Ridge...")
ridge_model = glmnet(X_train, y_train; alpha=0.0)
ridge_predictions = GLMNet.predict(ridge_model, X_test)

# Elastic Net
println("Entrenando modelo Elastic Net...")
elastic_net_model = fit(LassoPath, X_train, y_train; standardize=true)
elastic_net_predictions = Lasso.predict(elastic_net_model, X_test)[:, end]

# Decision Tree
println("Entrenando modelo Decision Tree...")
tree_model = DecisionTreeClassifier(max_depth=4)
DecisionTree.fit!(tree_model, X_train, y_train)
tree_predictions = DecisionTree.predict(tree_model, X_test)

# Random Forest
println("Entrenando modelo Random Forest...")
forest_model = RandomForestClassifier(n_trees=100, max_depth=4)
DecisionTree.fit!(forest_model, X_train, y_train)
forest_predictions = DecisionTree.predict(forest_model, X_test)

# k-Nearest Neighbors
println("Entrenando modelo k-Nearest Neighbors...")
function knn_classify(X_train, y_train, X_test; k=5)
    kd_tree = KDTree(X_train)
    predictions = []
    for i in 1:size(X_test, 1)
        idxs, _ = knn(kd_tree, X_test[i, :], k)
        neighbor_labels = y_train[idxs]
        push!(predictions, mode(neighbor_labels))
    end
    return predictions
end
knn_predictions = knn_classify(X_train, y_train, X_test; k=5)

# SVM
println("Entrenando modelo SVM...")

# Asegurarse de que y_train esté en un formato adecuado para LIBSVM
y_train_svm = Int.(2 .* y_train .- 1)  # Convertir 1.0 -> 1, 0.0 -> -1
println("Valores en y_train_svm: ", unique(y_train_svm))  # Depuración

# Entrenar el modelo SVM
svm_model = svmtrain(X_train, y_train_svm; kernel=LIBSVM.Kernel.RadialBasis, cost=1.0)

# Predecir con el modelo SVM
svm_predictions = LIBSVM.predict(svm_model, X_test)

# Paso 3: Evaluar los modelos
function evaluate_model(y_true, y_pred, model_name)
    cm = confusion_matrix(y_true, round.(y_pred))
    println("Matriz de confusión para $model_name:\n$cm")

    # Calcular la curva ROC
    fpr, tpr, _ = roc(y_true, y_pred)
    plot(fpr, tpr, label=model_name, xlabel="Tasa de Falsos Positivos", ylabel="Tasa de Verdaderos Positivos", title="Curva ROC")
end

evaluate_model(y_test, lasso_predictions, "LASSO")
evaluate_model(y_test, ridge_predictions, "Ridge")
evaluate_model(y_test, elastic_net_predictions, "Elastic Net")
evaluate_model(y_test, tree_predictions, "Decision Tree")
evaluate_model(y_test, forest_predictions, "Random Forest")
evaluate_model(y_test, knn_predictions, "k-NN")
evaluate_model(y_test, svm_predictions, "SVM")

# Guardar curvas ROC
savefig("/home/nat/Ejercicios diplomado/exercise_3/fig/roc_curves.png")

# Paso 4: Generar el informe en Markdown
println("Generando el informe...")
report_text = """
# Informe del Ejercicio 3: Clasificación

## Resumen
Este informe evalúa el desempeño de siete algoritmos de clasificación en el conjunto de datos del mercado de valores.

## Algoritmos
1. LASSO
2. Ridge
3. Elastic Net
4. Decision Tree
5. Random Forest
6. k-Nearest Neighbors
7. Support Vector Machines

## Resultados
Se generaron las matrices de confusión y las curvas ROC para cada modelo.

## Conclusión
El desempeño de cada modelo varía según su configuración y las características del conjunto de datos. El análisis completo y las visualizaciones se encuentran en los archivos asociados.
"""
write("/home/nat/Ejercicios diplomado/exercise_3/report/aReport.md", report_text)

println("Ejercicio completado. Resultados e informe guardados.")
