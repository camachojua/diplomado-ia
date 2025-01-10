using Images, FileIO, LinearAlgebra

function reducir_svd(img_matrix, k)    
    U, S, V = svd(img_matrix)
    return U[:, 1:k] * Diagonal(S[1:k]) * V[:, 1:k]'
end


function procesa_imagen(name, k)
    file_name = "./imagen/$(name).jpeg"
    println("\nProcesando ($(file_name))...")
    
    img = channelview(load(file_name))
    println("Imagen cargada. size(img) = $(size(img))")

    println("Reduciendo imagen...")
    img_reducida = zeros(size(img))
    for c in 1:size(img, 1)
        img_reducida[c, :, :] = reducir_svd(img[c, :, :], k)
    end
    println("Imagen reducida con $k valores singulares. size(img): $(size(img_reducida))")

    println("Guardando imagen...")
    ruta_salida = "./imagen/$(name)_red_k$(k).jpg"
    save(ruta_salida, colorview(RGB, clamp.(img_reducida, 0, 1)))
    println("Imagen reducida guardada en $ruta_salida\n")
end

function main()
    for i in 1:10
        for k in 5:5:50
            ruta_entrada = "img$(i)"
            procesa_imagen(ruta_entrada, k)
        end
    end
end

main()
