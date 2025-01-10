using CSV, DataFrames, Statistics, Random, Plots, StatsPlots

dataShape=size

function dataType(df::DataFrame)
    return Dict(column => eltype(dropmissing(df)[:,column]) for column in names(df))
end

function count_missing(df::DataFrame,col)
    return count(x -> ismissing(x), df[!, col])
end

function dataMissingPercentage(df::DataFrame)
    # Create a DataFrame with the column names and their missing value counts
    nrows=dataShape(df)[1]
    df=DataFrame(Column = names(df), MissingCount = map(column -> count(ismissing, df[!, column]), names(df)))
    df[!,"MissingPercentage"]=df[:,"MissingCount"]*100/nrows
    return df
end

function deleteColumns(df::DataFrame,threshold)
    return select(df,Not(filter(:MissingPercentage => x->x>threshold, dataMissingPercentage(df))[:,:Column]))
end

function deleteColumns!(df::DataFrame,threshold)
    return select!(df,Not(filter(:MissingPercentage => x->x>threshold, dataMissingPercentage(df))[:,:Column]))
end

function calculateCorrelation(df::DataFrame)
    num=filter(((k,v),) -> v <:Number, dataType(df))
    c=length(num)
    corr=ones(c,c)
    df=select(df,collect(keys(num)))
    for i in 1:c
        for j in 1:i
            if i!=j
                cols=select(df,i,j)
                cols=dropmissing(cols)
                try
                    corr[i,j]=cor(cols[:,1],cols[:,2])
                catch 
                    corr[i,j]=NaN
                end
            end
            corr[j,i]=corr[i,j]
            
        end
    end
    return corr, df
end

function displayCorrelation(df::DataFrame)
    matrix,data=calculateCorrelation(df)
    len=size(matrix)[1]
    heatmap(1:len,
    1:len,
    (x,y)->matrix[len+1-y,x],
    xticks=(1:len,names(data)),
    yticks=(1:len,reverse(names(data))),
    xrotation=90,
    size=(1000,1000))
end

function removeOutliersIQR(df::DataFrame)
    limits(q)=[q[1]-(1.5*(q[2]-q[1])),q[2]+(1.5*(q[2]-q[1]))]
    num=filter(((k,v),) -> v <:Number, dataType(df))
    quantiles=Dict(column=>limits(quantile(collect(skipmissing(df[:,column])),[0.25,0.75])) for column in collect(keys(num)))
    for column in collect(keys(num))
        df=filter(Symbol(column)=>x->(ismissing(x))||(quantiles[column][1]<=x<=quantiles[column][2]),df)
        println(size(df))
    end
    return df
end

function removeOutliersIQR!(df::DataFrame)
    limits(q)=[q[1]-(1.5*(q[2]-q[1])),q[2]+(1.5*(q[2]-q[1]))]
    num=filter(((k,v),) -> v <:Number, dataType(df))
    quantiles=Dict(column=>limits(quantile(collect(skipmissing(df[:,column])),[0.25,0.75])) for column in collect(keys(num)))
    for column in collect(keys(num))
        filter!(Symbol(column)=>x->(ismissing(x))||(quantiles[column][1]<=x<=quantiles[column][2]),df)
        println(size(df))
    end
    return df
end

function deleteRow(df::DataFrame,column)
    return filter(Symbol(column)=>!ismissing,df)
end

function deleteRow!(df::DataFrame,column)
    return filter!(Symbol(column)=>!ismissing,df)
end

function calculateCorrelation(df::DataFrame)
    num=filter(((k,v),) -> v <:Number, dataType(df))
    c=length(num)
    corr=ones(c,c)
    df=select(df,collect(keys(num)))
    for i in 1:c
        for j in 1:i
            if i!=j
                cols=select(df,i,j)
                cols=dropmissing(cols)
                try
                    corr[i,j]=cor(cols[:,1],cols[:,2])
                catch 
                    corr[i,j]=NaN
                end
            end
            corr[j,i]=corr[i,j]
            
        end
    end
    return corr, df
end

function filterColumnsByCorrelation(df::DataFrame,target,threshold,relation)
    filteredColumns=[]
    num=filter(((k,v),) -> v <:Number, dataType(df))
    cols=collect(keys(num))
    filter!(x->x!=target,cols)
    for c in cols
        data=dropmissing(select(df,target,c))
        corr=cor(data[:,1],data[:,2])
        if (abs(corr)>threshold) && (relation=="greater")
            push!(filteredColumns,c)
        elseif (abs(corr)<threshold) && (relation=="lesser")
            push!(filteredColumns,c)
        end
    end
    return select(df,Not(filteredColumns))
end

function filterColumnsByCorrelation!(df::DataFrame,target,threshold,relation)
    filteredColumns=[]
    num=filter(((k,v),) -> v <:Number, dataType(df))
    cols=collect(keys(num))
    filter!(x->x!=target,cols)
    for col in cols
        data=dropmissing(select(df,target,col))
        corr=cor(data[:,1],data[:,2])
        if (abs(corr)>threshold) && (relation=="greater")
            push!(filteredColumns,col)
        elseif (abs(corr)<threshold) && (relation=="lesser")
            push!(filteredColumns,col)
        end
    end
    return select!(df,Not(filteredColumns))
end