using LinearAlgebra, MLBase, GLMNet, DataFrames, CSV, Random, Plots, DecisionTree, RDatasets, PrettyTables, DataStructures, NearestNeighbors, Distances, LIBSVM, Distances

categorize_value(pred_val) = argmin(abs.(pred_val .- 1))
calculate_accuracy(predicted_vals, actual_vals) = sum(predicted_vals .== actual_vals) / length(actual_vals)

function split_by_year(data_years, proportion)
    unique_years = unique(data_years)
    selected_ids = []
    for year in unique_years
        indices = findall(data_years .== year)
        sampled_ids = randsubseq(indices, proportion)
        append!(selected_ids, sampled_ids...)
    end
    return selected_ids
end

stock_data = CSV.read("/Users/michelletorres/Desktop/Homeworks AI/Smarket.csv", DataFrame)
@show(stock_data)
@show(size(stock_data))

feature_matrix = Matrix(stock_data[:, 3:8])
println(names(stock_data))
labels = stock_data[:, 9]

label_mapping = MLBase.labelmap(labels)
encoded_labels = labelencode(label_mapping, labels)

training_indices = split_by_year(feature_matrix[:, 2], 0.7)
@show size(training_indices)
testing_indices = setdiff(1:length(feature_matrix[:, 2]), training_indices)
@show size(testing_indices)

# Lasso Regression
lasso_model = glmnet(feature_matrix[training_indices, :], encoded_labels[training_indices])
lasso_cv = glmnetcv(feature_matrix[training_indices, :], encoded_labels[training_indices])
optimal_lambda_lasso = lasso_model.lambda[argmin(lasso_cv.meanloss)]
lasso_model = glmnet(feature_matrix[training_indices, :], encoded_labels[training_indices], lambda=[optimal_lambda_lasso])
lasso_predictions = GLMNet.predict(lasso_model, feature_matrix[testing_indices, :])
lasso_predictions = categorize_value.(lasso_predictions)
calculate_accuracy(lasso_predictions, encoded_labels[testing_indices])

# Ridge Regression
ridge_model = glmnet(feature_matrix[training_indices, :], encoded_labels[training_indices], alpha=0)
ridge_cv = glmnetcv(feature_matrix[training_indices, :], encoded_labels[training_indices], alpha=0)
optimal_lambda_ridge = ridge_model.lambda[argmin(ridge_cv.meanloss)]
ridge_model = glmnet(feature_matrix[training_indices, :], encoded_labels[training_indices], alpha=0, lambda=[optimal_lambda_ridge])
ridge_predictions = GLMNet.predict(ridge_model, feature_matrix[testing_indices, :])
ridge_predictions = categorize_value.(ridge_predictions)
calculate_accuracy(ridge_predictions, encoded_labels[testing_indices])

# Elastic Net
elastic_net_model = glmnet(feature_matrix[training_indices, :], encoded_labels[training_indices], alpha=0.5)
elastic_net_cv = glmnetcv(feature_matrix[training_indices, :], encoded_labels[training_indices], alpha=0.5)
optimal_lambda_en = elastic_net_model.lambda[argmin(elastic_net_cv.meanloss)]
elastic_net_model = glmnet(feature_matrix[training_indices, :], encoded_labels[training_indices], alpha=0.5, lambda=[optimal_lambda_en])
en_predictions = GLMNet.predict(elastic_net_model, feature_matrix[testing_indices, :])
en_predictions = categorize_value.(en_predictions)
calculate_accuracy(en_predictions, encoded_labels[testing_indices])

# Decision Tree
tree_model = DecisionTreeClassifier(max_depth=2)
DecisionTree.fit!(tree_model, feature_matrix[training_indices, :], encoded_labels[training_indices])
tree_predictions = DecisionTree.predict(tree_model, feature_matrix[testing_indices, :])
calculate_accuracy(tree_predictions, encoded_labels[testing_indices])

# Random Forest
rf_model = DecisionTree.RandomForestClassifier(n_trees=20)
DecisionTree.fit!(rf_model, feature_matrix[training_indices, :], encoded_labels[training_indices])
rf_predictions = DecisionTree.predict(rf_model, feature_matrix[testing_indices, :])
calculate_accuracy(rf_predictions, encoded_labels[testing_indices])

# k-Nearest Neighbors
train_features = feature_matrix[training_indices, :]
train_labels = encoded_labels[training_indices]
kdtree_model = KDTree(train_features')
query_points = feature_matrix[testing_indices, :]
nearest_indices, distances = knn(kdtree_model, query_points', 5, true)
neighbors_labels = train_labels[hcat(nearest_indices...)]
label_counts = map(i -> counter(neighbors_labels[:, i]), 1:size(neighbors_labels, 2))
knn_predictions = map(i -> parse(Int, string(argmax(label_counts[i]))), 1:size(neighbors_labels, 2))
calculate_accuracy(knn_predictions, encoded_labels[testing_indices])

svm_model = svmtrain(train_features', train_labels) # Support vector machine
svm_predictions, decision_values = svmpredict(svm_model, feature_matrix[testing_indices, :]')
calculate_accuracy(svm_predictions, encoded_labels[testing_indices])

accuracy_scores = zeros(7) # Compare overall accuracies
methods_list = ["Lasso", "Ridge", "ElasticNet", "DecisionTree", "RandomForest", "kNN", "SVM"]
test_labels = encoded_labels[testing_indices]
accuracy_scores[1] = calculate_accuracy(lasso_predictions, test_labels)
accuracy_scores[2] = calculate_accuracy(ridge_predictions, test_labels)
accuracy_scores[3] = calculate_accuracy(en_predictions, test_labels)
accuracy_scores[4] = calculate_accuracy(tree_predictions, test_labels)
accuracy_scores[5] = calculate_accuracy(rf_predictions, test_labels)
accuracy_scores[6] = calculate_accuracy(knn_predictions, test_labels)
accuracy_scores[7] = calculate_accuracy(svm_predictions, test_labels)
hcat(methods_list, accuracy_scores)

# Confusion matrices
println("Confusion matrix Lasso")
pretty_table(confusmat(2, test_labels, lasso_predictions[:]))
println("Confusion matrix Ridge")
pretty_table(confusmat(2, test_labels, ridge_predictions[:]))
println("Confusion matrix Elastic Net")
pretty_table(confusmat(2, test_labels, en_predictions[:]))
println("Confusion matrix Decision Tree")
pretty_table(confusmat(2, test_labels, tree_predictions[:]))
println("Confusion matrix Random Forest")
pretty_table(confusmat(2, test_labels, rf_predictions[:]))
println("Confusion matrix kNN")
pretty_table(confusmat(2, test_labels, knn_predictions[:]))
println("Confusion matrix SVM")
pretty_table(confusmat(2, test_labels, svm_predictions[:]))
