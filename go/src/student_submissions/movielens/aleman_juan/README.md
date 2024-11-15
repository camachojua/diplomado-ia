# MovieLens Data Processing Project in Go

Este proyecto procesa datos de MovieLens para calcular el número de calificaciones y el promedio de calificación por género. Está implementado en Go y utiliza goroutines para paralelizar el procesamiento de archivos grandes.

## Estructura del Proyecto

- `movielens.go`: Punto de entrada principal. Realiza la división de datos, lanza los trabajadores en paralelo y consolida los resultados.
- `split_csv.go`: Contiene la función para dividir un archivo CSV grande en varios archivos CSV más pequeños.

## Requisitos

Asegúrate de tener Go instalado. Puedes descargarlo de [Go.dev](https://go.dev/dl/).

## Configuración y Ejecución

1. **Dividir el archivo CSV grande en archivos más pequeños**:
   - `split_csv.go` contiene la función que divide `ratings.csv` en archivos más pequeños.
   - Coloca el archivo `ratings.csv` en el mismo directorio que `movielens.go`.

2. **Ejecución del Proyecto**:
   - Ejecuta el archivo `movielens.go`:

     ```bash
     go run movielens.go
     ```

   - `movielens.go` realizará los siguientes pasos:
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
