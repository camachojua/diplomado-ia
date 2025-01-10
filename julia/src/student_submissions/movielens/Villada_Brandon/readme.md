# MovieLens Project

Este proyecto utiliza el conjunto de datos MovieLens para realizar análisis y recomendaciones de películas.

## Requisitos

- Julia 1.6 o superior
- DataFrames.jl
- CSV.jl
- Base.Threads

## Uso

Para ejecutar el script con 10 threads, utiliza el siguiente comando:

```sh
julia -t 10 movielens.jl
```

## Descripción

El script `movielens.jl` realiza las siguientes tareas:
- Carga y preprocesa los datos de MovieLens en archivos mas pequeños.
- Realiza análisis exploratorio de los datos.
- Implementa un sistema de recomendación básico.

## Estructura del Proyecto

- `movielens.jl`: Script principal para ejecutar el análisis y las recomendaciones.

## Créditos
Este proyecto fue desarrollado como parte del Diplomado en Inteligencia Artificial.
