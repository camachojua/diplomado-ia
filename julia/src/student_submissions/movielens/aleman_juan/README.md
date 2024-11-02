# MovieLens Data Processing Project in Julia

Este proyecto procesa datos de MovieLens, permitiendo calcular el número de calificaciones y el promedio de calificación por género. El proyecto está implementado en Julia y utiliza procesamiento en paralelo para mejorar la eficiencia al dividir y analizar los datos.

## Estructura del Proyecto

- `main.jl`: Punto de entrada principal del proyecto. Realiza la división de datos, lanza los trabajadores en paralelo y consolida los resultados.
- `split_csv.jl`: Lee un archivo CSV grande y lo divide en 10 archivos Parquet más pequeños.
- `movies.jl`: Contiene la función para cargar los datos de películas y asociarlos con los géneros.
- `workers.jl`: Define la función `find_ratings_worker`, que procesa cada archivo Parquet en paralelo y calcula las calificaciones y el conteo por género.
- `utils.jl`: Funciones auxiliares, incluida la consolidación de resultados para mostrar el conteo y promedio de calificaciones por género.

## Requisitos

Asegúrate de tener Julia instalado. Puedes descargarlo de [JuliaLang.org](https://julialang.org/downloads/).

### Librerías necesarias

Este proyecto requiere algunas librerías adicionales de Julia. Instálalas desde el modo de paquetes de Julia (`]`):

```julia
] add CSV
] add DataFrames
] add Parquet
] add Printf
```

2. **Ejecución del Proyecto**:
   - Ejecuta el archivo `main.jl`:

     ```bash
     julia main.jl
     ```

   - `main.jl` realizará los siguientes pasos:
     - Dividir el archivo `ratings.csv` en múltiples archivos más pequeños.
     - Cargar los datos de las películas desde `movies.csv`.
     - Procesar cada archivo en paralelo usando goroutines para calcular el número de calificaciones y el promedio de calificación por género.
     - Consolidar los resultados y mostrar el informe final.

## Ejemplo de Salida

La salida mostrará el índice del género, el nombre del género, el conteo total de calificaciones y el promedio de calificación por género, junto con el tiempo total de procesamiento:

```plaintext
 0  Action               7446893      3.56
 1  Adventure            5832400      4.12
 2  Animation            1630979      4.01
 3  Children             2124250      3.89
 ...
Duración total = 13.747019914s