using Images, LinearAlgebra

function load_grayscale_image(filepath) # Upload grayscale images
    img = Images.load(filepath)
    gray_img = Gray.(img)
    return Float64.(gray_img) # Convert to array of float values
end

function perform_svd(img_matrix)
    U, S, V = svd(img_matrix)
    return U, S, V
end

function reconstruct_image(U, S, V, num_components)
    U_reduced = U[:, 1:num_components]
    S_reduced = Diagonal(S[1:num_components])
    V_reduced = V[:, 1:num_components]
    return U_reduced * S_reduced * V_reduced'
end

function normalize_image(img_matrix) # Function to normalize the image to values ​​within the range [0, 1]
    min_val = minimum(img_matrix)
    max_val = maximum(img_matrix)
    normalized_img = (img_matrix .- min_val) ./ (max_val - min_val)
    return normalized_img
end

input_dir = "/Users/michelletorres/Desktop/Homeworks AI/"
output_dir = "/Users/michelletorres/Desktop/Homeworks AI/Reconstructed/"
isdir(output_dir) || mkdir(output_dir) # Create the output directory if it does not exist

image_names = ["Imagen 1.jpeg", "Imagen 2.jpeg", "Imagen 3.jpeg", 
               "Imagen 4.jpeg", "Imagen 5.jpeg", "Imagen 6.jpeg", 
               "Imagen 7.jpeg"]

num_components = 50  # Num for reconstruction

for image_name in image_names # Process each image
    try
        image_path = joinpath(input_dir, image_name)
        println("Procesando: $image_path")

        img_matrix = load_grayscale_image(image_path) # Load the original image
        println("Dimensiones de la imagen: ", size(img_matrix)) 

        U, S, V = perform_svd(img_matrix) # Perform SVD
        println("SVD completado con éxito.")

        reconstructed_img = reconstruct_image(U, S, V, num_components)
        println("Dimensiones de la imagen reconstruida: ", size(reconstructed_img))

        normalized_img = normalize_image(reconstructed_img) # Normalize
        println("Imagen normalizada.")

        output_path = joinpath(output_dir, "reconstructed_$image_name") # Save
        save(output_path, Gray.(normalized_img))
        println("Imagen guardada con éxito en '$output_path'.")
    catch e
        println("Error al procesar $image_name: ", e)
    end
end
