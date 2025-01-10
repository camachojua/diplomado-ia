using Graphs
using GraphPlot
using Random
using Colors
using Plots


function random_geometric_graph(n, radius)
    pos = [(rand(), rand()) for _ in 1:n]
    g = SimpleGraph(n)
    for i in 1:n, j in 1:n
        if i < j
            xi, yi = pos[i]
            xj, yj = pos[j]
            if sqrt((xi - xj)^2 + (yi - yj)^2) <= radius
                add_edge!(g, i, j)
            end
        end
    end
    return g, pos
end

function get_ncenter(pos)
    dmin = 1
    ncenter = 0
    for (i, (x, y)) in enumerate(pos)
        d = (x - 0.5)^2 + (y - 0.5)^2
        if d < dmin
            ncenter = i
            dmin = d
        end
    end
    return ncenter
end

function main()
    
    # Grafo geométrico aleatorio con 200 nodos y radio 0.125
    two_random_values = rand(), rand()
    n = 200
    G, pos = random_geometric_graph(n, 0.125)

    # Colores para los nodos según la distancia
    distances = dijkstra_shortest_paths(G, get_ncenter(pos)).dists
    colors = [RGB(1.0 - d / maximum(distances), 0.3, 0.3) for d in distances]

    # Nodos el grafo
    s = scatter([p[1] for p in pos],
                [p[2] for p in pos],
                color=colors,
                legend=false,
                ms=5,
                title="Random Geometric Graph")
    
    # Aristas del grafo
    for i in 1:n
        for j in neighbors(G, i)
            plot!([pos[i][1], pos[j][1]], [pos[i][2], pos[j][2]], lw=0.2, color=:black)
        end
    end

    println("Guardando imagen...")
    img_name = "random_geometric_graph.png"
    savefig(img_name)
    println("Imagen $(img_name) guardada.")
end

main()
