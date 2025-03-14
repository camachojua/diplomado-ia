using CSV
using DataFrames
using StatsBase

tiempo=time()
archivo="ratings.csv"
N=10
v_datos=Vector{Vector{String}}(undef, N)
movies=CSV.read("movies.csv", DataFrame)
DF=DataFrame(Genero=String[], Numero_de_rankings=Int[])
function dividir(archivo,N)
	df=CSV.read(archivo, DataFrame)
	x=nrow(df)
	partes=div(x,N)
	for i in 0:N-1
		inicio= (i* partes)+1
		fin=(i+1)*partes
		if i==N-1
			resto=x-fin
			fin+=resto
		end
		df_div=df[inicio:fin, :]
		CSV.write("rating$(i+1).csv", df_div)
	end
	t_d=time() - tiempo
	println("Tiempo de division de archivos= $t_d s")
end

function resultados(n, movies)
	rating_parts=CSV.read("rating$n.csv", DataFrame)
	return contar(rating_parts,movies)
end

function contar(rating_parts,movies)
	datos=innerjoin(rating_parts,movies, on=:movieId)
	generos=[split(row,"|") for row in datos.genres]
	v_generos=reduce(vcat, generos)
	return v_generos
end
dividir(archivo, N)
Threads.@threads for n in 1:N
		v_datos[n]=resultados(n,movies)
			
end
for i in v_datos
	conteo_g=countmap(i)
	for (genero, conteo) in conteo_g
		push!(DF,(genero, conteo))
	end
end

Datos=combine(groupby(DF, :Genero), :Numero_de_rankings => sum => :Numero_de_rankings)
sort!(Datos, :Numero_de_rankings, rev=true)
println(Datos)
duracion=time()-tiempo
println("Duracion total= $duracion s")