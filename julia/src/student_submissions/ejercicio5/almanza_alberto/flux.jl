using CSV, Flux, Statistics, MLDatasets, DataFrames, OneHotArrays

const classes = ["Iris-setosa", "Iris-versicolor", "Iris-virginica"];

function flux_loss(flux_model, features, labels_onehot)
           ŷ = flux_model(features)
           Flux.logitcrossentropy(ŷ, labels_onehot)
end;

function custom_onecold(labels_onehot)
        max_idx = [x[1] for x in argmax(labels_onehot; dims=1)]
        return vec(classes[max_idx])
end;

function train_custom_model!(f_loss, weights, biases, features, labels_onehot)
    dLdW, dLdb, _, _ = gradient(f_loss, weights, biases, features, labels_onehot)
    weights .= weights .- 0.1 .* dLdW
    biases .= biases .- 0.1 .* dLdb
end;

function main()

    df = CSV.read("Churn_Modelling.csv", DataFrame)
    @show(describe(df))

    # One Hot Encoding
    replace!(df.Gender,"Female" => "1")
    replace!(df.Gender,"Male" => "0")
    replace!(df.Geography, "France"  => "2")
    replace!(df.Geography, "Germany" => "1")
    replace!(df.Geography, "Spain"   => "0")
    
    y = select(df, [:Exited])    
    x = select(df, [:CreditScore,
                    :Age,
                    :Tenure,
                    :Balance,
                    :NumOfProducts,
                    :HasCrCard,
                    :IsActiveMember,
                    :EstimatedSalary])

    x, y = Iris(as_df=false)[:];
    @show typeof(x)
    @show typeof(y)
    
    x = Float32.(x);
    @show typeof(x)
    y = vec(y);
    @show typeof(y)
    
    custom_y_onehot = unique(y) .== permutedims(y)
    println("custom_y_onehot: ", custom_y_onehot)
    

    flux_y_onehot = onehotbatch(y, classes)
    println("flux_y_onehot: ", flux_y_onehot)
    
    # Building a model
    m(W, b, x) = W*x .+ b
    W = rand(Float32, 3, 4);
    b = [0.0f0, 0.0f0, 0.0f0];
    
    custom_softmax(x) = exp.(x) ./ sum(exp.(x), dims=1)
    custom_model(W, b, x) = m(W, b, x)
    @show custom_model(W, b, x)
    @show all(0 .<= custom_model(W, b, x) .<= 1)
    @show sum(custom_model(W, b, x), dims=1)
    @show flux_model = Chain(Dense(4 => 3), softmax)
    @show flux_model[1].weight, flux_model[1].bias
    

    # Loss and accuracy
    custom_logitcrossentropy(ŷ, y) = mean(.-sum(y .* logsoftmax(ŷ; dims = 1); dims = 1));
    function custom_loss(weights, biases, features, labels_onehot)
           ŷ = custom_model(weights, biases, features)
           custom_logitcrossentropy(ŷ, labels_onehot)
    end;

    @show custom_loss(W, b, x, custom_y_onehot)
    @show flux_loss(flux_model, x, flux_y_onehot)
    @show argmax(custom_y_onehot, dims=1) 
    max_idx = [x[1] for x in argmax(custom_y_onehot; dims=1)]
    println(max_idx)    
    @show(custom_onecold(custom_y_onehot))
    istrue = Flux.onecold(flux_y_onehot, classes) .== custom_onecold(custom_y_onehot);
    println(istrue)
    @show(all(istrue))
    custom_accuracy(W, b, x, y) = mean(custom_onecold(custom_model(W, b, x)) .== y);
    @show(custom_accuracy(W, b, x, y))
    flux_accuracy(x, y) = mean(Flux.onecold(flux_model(x), classes) .== y);
    @show(flux_accuracy(x, y))
    
    
    # Training the model
    dLdW, dLdb, _, _ = gradient(custom_loss, W, b, x, custom_y_onehot);
    W .= W .- 0.1 .* dLdW;
    b .= b .- 0.1 .* dLdb;
    @show(custom_loss(W, b, x, custom_y_onehot))
    for i = 1:500
        train_custom_model!(custom_loss, W, b, x, custom_y_onehot);
        custom_accuracy(W, b, x, y) >= 0.98 && break
    end

    @show custom_accuracy(W, b, x, y);
    @show custom_loss(W, b, x, custom_y_onehot)
    @show flux_loss(flux_model, x, flux_y_onehot)
    @show flux_accuracy(x, y);
    @show flux_loss(flux_model, x, flux_y_onehot)
end
main()



