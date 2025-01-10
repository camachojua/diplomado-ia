using CSV
using DataFrames
using Plots
using GLM
using Statistics
using Distributions
using Random
using MultivariateStats
using MLBase
using Makie
using GLMakie
using Plots

function SplittrainTestSets(df,split)
    nrows = size(df)[1]
    nrowsTrain = round(Int, (nrows*split))
    nrowsTest = round(Int, (nrows - nrowsTrain))
    return df[1:nrowsTrain, :], df[nrowsTrain+1 : nrows,:]
end

     
function main()

    df = CSV.read("Churn_Modelling.csv", DataFrame)
    @show(describe(df))
    @show(GLM.countmap(df.Surname))
    
    # One Hot Encoding
    replace!(df.Gender,"Female" => "1")
    replace!(df.Gender,"Male" => "0")
    replace!(df.Geography, "France"  => "2")
    replace!(df.Geography, "Germany" => "1")
    replace!(df.Geography, "Spain"   => "0")

    #Prediccion
    dftrn , dftst = SplittrainTestSets(df, 0.75)
    fm = @formula(
        Exited ~ CreditScore
        + Age
        + Tenure
        + Balance
        + NumOfProducts
        + HasCrCard
        + IsActiveMember
        + EstimatedSalary
        + Gender
        + Geography
    )
    prediction = predict(glm(fm, dftrn, Binomial(), ProbitLink()), dftst)
    pred_df = DataFrame(
        y_actual=dftst.Exited,
        y_predicted=[if x < 0.5 0 else 1 end for x in prediction],
        prob_predicted=prediction
    )
    pred_df.correctly_classified = pred_df.y_actual .== pred_df.y_predicted
    accuracy = mean(pred_df.correctly_classified)
    println("acc: ", accuracy)
    
    # Matriz de Confusion
    cm = MLBase.roc(pred_df.y_actual, pred_df.y_predicted)
    acc = (cm.tp + cm.tn) / (cm.p + cm.n)
    println("acc: ", acc)

end

main()
