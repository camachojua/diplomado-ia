#este código propone preprocesar los datos y aprovechar las fortalezas de la librería DataFrames de Julia
# en la función "countGenres" se comentan los pasos de esta diferente metodología

using Pkg
using BenchmarkTools
using DataFrames
using CSV
using Base.Threads
using Parquet

function readStuff(filename,has_header)
    if lowercase(reverse(reverse(filename)[1:4])) == ".csv"
        if has_header == true
            ratings = DataFrame(CSV.File(filename))
        else
            ratings = DataFrame(CSV.File(filename,header=false))
            ratings = rename(ratings, :Column1 => :userId, :Column2 => :movieId,:Column3 => :rating)
        end
    elseif lowercase(reverse(reverse(filename)[1:8])) == ".parquet"
        ratings = DataFrame(Parquet.read_parquet(filename))
    else
        println("pasame un csv plis")
        return 
    end
    return ratings
end

function countGenres(ratingsIn,moviesIn)
    #keep only relevant columns
    ratings = select(ratingsIn,:movieId,:rating)
    #separate string containing n genres into an array of length n with each entry having one string corresponding to a genre
    movies = transform!(moviesIn, :genres => ByRow(x -> ismissing(x) ? [missing] : string.(split(x, "|"))) => :flattened_genres)
    #expand movies dataframe genres column i.e. instead of 1 row of a movie A with n genres, we get n rows of movie A with each of its genres
    movies = flatten(movies,:flattened_genres)
    #we only keep movie id, and genre columns 
    movies = select(movies,:movieId,:flattened_genres =>:genres)
    #we do an inner join on movie id, this will result in each rating of movie A, matching to its n rows
    #the end result dataframe will contain all the reviews per genre
    movies = innerjoin(movies,ratings, on = :movieId)
    #we group by genre and count each ocurrence, as well as obtain its average rating
    movies = combine(groupby(movies,:genres), nrow => :count,:rating => mean => :rating)
    #finally we sort by genre
    movies = sort(movies,:genres)
    #we return the dataframe
    return movies
end

function  chunkPostProcessing(finres)
    finres = select(finres,:genres,:count,[:count,:rating] => ((cnt,rtng) -> cnt.*rtng) => :unweighted_mean)
    finres = sort(combine(groupby(finres,[:genres]), :count => sum => :count, :unweighted_mean => sum => :unweighted_mean),:genres)
    finres = select(finres,:genres,:count,[:count,:unweighted_mean]=>((cnt,umn) -> umn./cnt)=> :rating)
    return finres
end

function bufferRatings(numberOfChunks,start_w_zero,format)
    res = [DataFrame() for _ in 1:10]
    movieso = DataFrame(CSV.File("movies.csv"))
    if start_w_zero == true
        eff_range = range(0, step=1, length=numberOfChunks)
        offset = 1
    else
        eff_range = range(1, step=1, length=numberOfChunks)
        offset = 0
    end
    
    @threads for i in eff_range
        if format == ".parquet"
            filenameCounter = "_"*lpad(i, 2, '0')
        else
            filenameCounter = i
        end
        chunkFilename = string("ratings",filenameCounter,format)
        res[i+offset] = countGenres(readStuff(chunkFilename),movieso)
    end
    
    finres = DataFrame()
    for i in eff_range
        finres = [finres;res[i+offset]]
    end

    return chunkPostProcessing(finres)
end
#the function bufferRatings will analyze the data and tally the results if it is split apart, but the function countGenres can do the whole file in one group
#10 csv files
@btime processInCSVChunks = bufferRatings(10,true,".csv")
#10 parquet files
@btime processInPARQUETChunks = bufferRatings(10,false,".parquet")
# 1 csv file
@btime processOneCSVChunk = countGenres(readStuff("ratings.csv",true),readStuff("movies.csv",true))