using Pkg
Pkg.add(["CSV", "DataFrames"])
using CSV
using DataFrames
using Distributed 
using Base.Threads

struct RecordWithIndex
	record::Vector{String}
	index::Int
end

function main()
	@time split_csv_file("ratings.csv", "movies.csv", 10)
end

function split_csv_file(ratings_file::String, movies_file::String, num_files::Int)
	record_channel = Channel{RecordWithIndex}(100)
	@sync for i in 1:num_files
		Threads.@spawn write_csv_file("ratings_$i.csv", record_channel, num_files, i - 1)
	end
	
	read_ratings_csv_file(ratings_file, record_channel)
	close(record_channel)
	
	genres_count = count_ratings_by_genre(ratings_file, movies_file)
	
	for (genre, count) in genres_count
		println("Genre: $genre, Ratings: $count")
	end
end

function read_ratings_csv_file(filename::String, record_channel::Channel{RecordWithIndex})
	df = CSV.File(filename, header = true) |> DataFrame
	
	#Send thhe header
	put!(record_channel, RecordWithIndex(names(df), -1))
	
	#Send 
	for (i, row) in enumerate(eachrow(df))
		record = [string(row[col]) for col in names(df)]
		put!(record_channel, RecordWithIndex(record, i))
	end
end


function write_csv_file(filename::String, record_channel::Channel{RecordWithIndex}, num_files::Int, index::Int)

	rowns = DataFrame()
	for record_info in record_channel
			if record_info.index == -1
				rows = DataFrame([record_info.record], Symbol.(record_info.record))

			elseif record_info.index % num_files == index
				push!(rows, rec.record)
			end
		end
		CSV.write(filename, rows)
	end

function count_ratings_by_genre(ratings_file::String, movies_file::String)
	genres_count = Dict{String, Int}()
	movies_df = CSV.File(movies_file, header = true) |> DataFrame
	movie_genres = Dict{String, Vector{String}}()

	for row in eachrow(movies_df)
		movie_id = string(row["movieId"])
		genres = split(string(row["genres"]), "|")
		movie_genres[movie_id] = genres
	end

	ratings_df = CSV.File(ratings_file, header = true) |> DataFrame

	for row in eachrow(ratings_df)
		movie_id = string(row["movieId"])
		if haskey(movie_genres, movie_id)
			for genre in movie_genres[movie_id]
				genres_count[genre] = get(genres_count, genre, 0) + 1
			end
		end
	end
	return genres_count
end

main()

