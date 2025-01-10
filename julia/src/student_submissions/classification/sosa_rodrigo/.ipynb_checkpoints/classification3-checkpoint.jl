#using Pkg

#Pkg.add("StatsModels")
#Pkg.add("GLM")

using CSV, DataFrames, Plots, MultivariateStats, MLJ, ROCAnalysis, GLM
using GLMNet, LinearAlgebra, StatsBase, CategoricalArrays, StatsModels

# Load dataset
df = CSV.read("Smarket.csv", DataFrame)

# Convertir la columna Direction a valores numéricos
df.Direction = ifelse.(df.Direction .== "Up", 1, 0)

# Normalizar las variables predictoras
for col in [:Lag1, :Lag2, :Lag3, :Lag4, :Lag5, :Volume]
    df[!, col] = (df[!, col] .- mean(df[!, col])) ./ std(df[!, col])
end

# Dividir en conjuntos de entrenamiento y prueba
train_ratio = 0.8
n_train = Int(train_ratio * nrow(df))
train_data = df[1:n_train, :]
test_data = df[n_train+1:end, :]

X_train = select(train_data, Not(:Direction)) |> Matrix
y_train = train_data.Direction

X_test = select(test_data, Not(:Direction)) |> Matrix
y_test = test_data.Direction

# Función para evaluar el modelo
function evaluate_model(predictions, probabilities, y_test)
    cm = confusion_matrix(y_test, predictions)
    println("Confusion Matrix:\n", cm)
    
    # Calcular precisión
    tp = cm[1, 1]  # Verdaderos positivos
    fp = cm[1, 2]  # Falsos positivos
    fn = cm[2, 1]  # Falsos negativos
    tn = cm[2, 2]  # Verdaderos negativos
    acc = (tp + tn) / (tp + tn + fp + fn)
    println("Accuracy: ", acc)

    roc = roc_curve(y_test, probabilities)
    auc = auc(roc)
    println("AUC: ", auc)
    plot(roc, title="ROC Curve", xlabel="False Positive Rate", ylabel="True Positive Rate")
end

# LASSO
lasso_model = glmnet(X_train, y_train, alpha=1.0, lambda=[0.1]) # α=1.0 para LASSO
lasso_probs = GLMNet.predict(lasso_model, X_test)[:, 1]
lasso_preds = round.(lasso_probs)
println("\nLASSO")
evaluate_model(lasso_preds, lasso_probs, y_test)
