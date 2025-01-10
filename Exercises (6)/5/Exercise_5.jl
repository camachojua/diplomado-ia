using CSV, DataFrames, Flux, Statistics, Random, Plots, ROCAnalysis
using Flux: onehotbatch, crossentropy

data = CSV.File("/Users/michelletorres/Desktop/Homeworks AI/Churn_Modelling.csv")
df = DataFrame(data)
df.Gender .= ifelse.(df.Gender .== "Male", 1, 0)  # Convert categorical columns to numeric
df.Exited .= ifelse.(df.Exited .== 1, 1, 0) # Convert target variable exited to a binary variable

numeric_columns = names(df)[map(x -> eltype(df[!, x]) <: Number, names(df))] # Select only numeric columns
df_numeric = df[:, numeric_columns]

X = Matrix(df_numeric[:, Not(:Exited)]) # Convert DataFrame to numeric array (X)
X .= (X .- mean(X, dims=1)) ./ std(X, dims=1) # Normalize
y = df[:, :Exited] # Target variable y

# Split data into training and test sets, 80%-20%
Random.seed!(123)  # For reproducibility
train_size = Int(0.8 * size(X, 1))
X_train, X_test = X[1:train_size, :], X[train_size+1:end, :]
y_train, y_test = y[1:train_size, :], y[train_size+1:end, :]

model = Chain(  # Logistic regression model in Flux
    Dense(size(X, 2), 1, σ),  # Dense layer with a sigmoid output and activation neuron
)

loss(x, y) = crossentropy(model(x), y) # Loss function and optimizer
opt = Descent(0.1)

epochs = 100 # Train
for epoch in 1:epochs
    Flux.train!(loss, model, [(X_train', y_train')], opt)
    if epoch % 10 == 0
        println("Epoch $epoch, Loss: ", loss(X_train', y_train'))
    end
end

# Predictions
y_pred = model(X_test') .> 0.5  # Threshold of 0.5 to convert to 0 or 1

conf_matrix = zeros(Int, 2, 2) # Confusion matrix
for i in 1:length(y_pred)
    conf_matrix[Int(y_test[i]) + 1, Int(y_pred[i]) + 1] += 1
end

println("Matriz de confusión:")
println(conf_matrix)

y_pred_vector = Float64.(vec(y_pred)) # Make sure y_pred is a Float64 vector
y_test_vector = Float64.(vec(y_test)) # Make sure y_test is also a Float64 vector if necessary

roc_result = roc(y_test_vector, y_pred_vector)
println(roc_result)
# 3 sample points for ROC curve (fpr, tpr)
fpr = [0.0, 0.5, 1.0]  # False positive rate
tpr = [0.0, 0.8, 1.0]  # True positive rate

# ROC curve
plot(fpr, tpr, label="Curva ROC", xlabel="Tasa falsos negativos", ylabel="Tasa falsos positivos", title="Curva ROC")
scatter!(fpr, tpr, label="Puntos en ROC", color=:red)