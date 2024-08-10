# Instalación

La instalación de Julia (en linux) es relativamente sencilla, [siguiendo la](https://julialang.org/downloads/) sólo debemos ejecutar el siguiente comando:

```bash
curl -fsSL https://install.julialang.org | sh
```

Una vez que el comando termine de ejecutarse podemos probar que la instalación fue exitosa al ejecutar:

```bash
julia -v
```
Lo cual debe regresar algo como `julia Version 1.10.4`

# Mi primer lenguaje en julia

Para nuestro primer programa en julia debemos crear una carpeta para contener
nuestro código:

```
mkdir -p ~/hello_julia
cd ~/hello_julia
```

Ahora abriremos nuestro editor de código para crear un archivo llamado `hello_julia.jl` con el siguiente contenido:

```julia
println("¡Hola Mundo!")
```

Al guardar el archivo podemos ejecutarlo en la terminal con el comando:

```bash
julia hello_julia.jl
```

Esto imprimirá la cadena `¡Hola Mundo!` en la terminal.
