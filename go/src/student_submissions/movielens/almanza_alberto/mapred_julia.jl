using CSV
using DataFrames


function split_csv()
    n = Threads.nthreads()
    df = CSV.read("./test/ratings.csv", DataFrame)    
    num_filas = size(df, 1)
    parte_tamaño = ceil(Int, num_filas/n)
    Threads.@threads for i in 1:n
        df_i = df[(i-1)*parte_tamaño + 1 : min(i*parte_tamaño, num_filas), :]
        CSV.write(string("./test/ratings_part_", i, ".csv"), df_i)
    end
end

function count_ratings()
    movies_dict = Dict{String, String}()
    open("./test/movies.csv") do mvs_file
        for mv in eachline(mvs_file)
            mv_spl = split(mv, ',')
            movies_dict[mv_spl[1]] = mv_spl[lastindex(mv_spl)]
        end
    end

    results = []

    Threads.@threads for i in 1:(Threads.nthreads())
        res_dict = Dict{String, Array{Float64}}()
        open(string("./test/ratings_part_", i, ".csv")) do rtg_file
            for rtg in eachline(rtg_file)
                gens = movies_dict[string(split(rtg, ",")[2])]
                for gen in split(gens, '|')
                    if gen ∉ keys(res_dict)
                        res_dict[gen] = [0, 0]
                    end
                    res_dict[gen][1] += 1
                    try
                        res_dict[gen][2] += parse(Float64, split(rtg, ",")[3])
                    catch
                        res_dict[gen][2] += 0
                    end
                end
            end
        end
        push!(results, res_dict)
    end
    results_dict = Dict{String, Array{Float64}}()
    for res in results      
        for (gen, count) in res
            if gen ∉ keys(results_dict)
                results_dict[gen] = [0, 0]
            end
            results_dict[gen][1] += count[1]
            results_dict[gen][2] += count[2]
        end
    end
    for (gen, count) in results_dict
        println(gen, ": ", Int64(count[1]), " - ", round(count[2]/count[1], digits=2))

    end
end

println("\n Partiendo archivo...\n")
split_csv()

println(" Contando ratings...\n")
@time count_ratings()

