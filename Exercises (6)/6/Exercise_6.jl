using LightGraphs, Plots, Random, LinearAlgebra 

n = 200  # Number of nodes
radius = 0.125  # Connection radius
positions = [(rand(), rand()) for _ in 1:n] # Create random positions for the nodes

g = SimpleGraph(n) # Create graph
for i in 1:n
    for j in (i+1):n
        if norm(positions[i] .- positions[j]) ≤ radius
            add_edge!(g, i, j)
        end
    end
end

center = (0.5, 0.5) # Find the node closest to the center (0.5, 0.5)
distances = [norm(p .- center) for p in positions]
center_node = argmin(distances)
path_lengths = dijkstra_shortest_paths(g, center_node).dists # Calculate the shortest path lengths from the central node

max_length = maximum(path_lengths) # Normalize the path lengths for color assignment
min_length = minimum(path_lengths)
scaled_lengths = (path_lengths .- min_length) * 10  # Avoid normalization over a too-small range -> scale distances to increase variability; multiplying by 10 to amplify the difference

max_scaled_length = maximum(scaled_lengths) # Recalculate max and min after scaling
min_scaled_length = minimum(scaled_lengths)

if max_scaled_length == min_scaled_length # Ensure the range is not zero
    normalized_lengths = zeros(n)
else
    normalized_lengths = (scaled_lengths .- min_scaled_length) / (max_scaled_length - min_scaled_length)
end

colormap = cgrad([:red, :white], 256)  # 256 shades between red and white
normalized_indices = round.(Int, normalized_lengths * 255) .+ 1  # Convert normalized values to color indices

p = plot()
scatter!(
    [p[1] for p in positions], [p[2] for p in positions],
    c=[colormap[i] for i in normalized_indices],  # Using the normalized colors
    legend=false, markersize=8
)

for edge in edges(g) # Draw edges
    x1, y1 = positions[src(edge)]
    x2, y2 = positions[dst(edge)]
    plot!(p, [x1, x2], [y1, y2], color=:gray, alpha=0.5, linewidth=0.5)
end

scatter!(p, [positions[center_node][1]], [positions[center_node][2]],
         color=:red, label="Nodo central", markersize=10)

xlims!(p, -0.05, 1.05)
ylims!(p, -0.05, 1.05)
title!(p, "Visualization")

display(p)
savefig(p, "random_geometric_graph.png")