package main

import (
	"fmt"
	"strconv"
	"time"
)

// Poner en True si se quiere hacer la fragmentación. En caso de ya tenerla no es necesario rehacerla y se puede omitir con False
const shred bool = false
const count_rating = true

func main() {

	start := time.Now()

	if shred {
		//Ruta del archivo a fragmentar
		var rating_pth string = "C:/Users/felip/Desktop/GoLang/MovieLense/ratings.csv"

		iTotalRatings, _ := getFileLength(rating_pth)

		iNumCPU, sliceSubFileSize := getSubFilesInfo(iTotalRatings)

		arrIndex := getFirstnLastIndex(sliceSubFileSize)

		for i := 0; i < iNumCPU; i++ {

			arrFileContent, _ := getFileContenToObj(rating_pth, arrIndex[i][0], arrIndex[i][1])
			getParquetFromObj("Ratings_"+strconv.Itoa(i+1)+".parquet", arrFileContent)
		}
	}

	if count_rating {
		// Ruta de prueba para archivos de películas y calificaciones
		var sMovieParquetPath = "C:/Users/felip/Desktop/GoLang/MovieLense/movie_go.parquet"
		sRatingParquetPaths := []string{"C:/Users/felip/Desktop/GoLang/MovieLense/Ratings_1.parquet",
			"C:/Users/felip/Desktop/GoLang/MovieLense/Ratings_2.parquet",
			"C:/Users/felip/Desktop/GoLang/MovieLense/Ratings_3.parquet",
			"C:/Users/felip/Desktop/GoLang/MovieLense/Ratings_4.parquet",
			"C:/Users/felip/Desktop/GoLang/MovieLense/Ratings_5.parquet",
			"C:/Users/felip/Desktop/GoLang/MovieLense/Ratings_6.parquet",
			"C:/Users/felip/Desktop/GoLang/MovieLense/Ratings_7.parquet",
			"C:/Users/felip/Desktop/GoLang/MovieLense/Ratings_8.parquet",
		}

		// Llamamos a la función para contar y sumar los ratings por género
		//genreRatingData := countAndSumMovieRatings(sMovieParquetPath, sRatingParquetPaths)
		genreRatingData := countAndSumMovieRatings_limited(sMovieParquetPath, sRatingParquetPaths)

		// Imprimimos los resultados de conteo, suma y promedio por género
		fmt.Println("Datos de Rating por Género:")
		for genre, data := range genreRatingData {
			fmt.Printf("Género: %s, Conteo: %v, Suma: %v, Promedio: %.2f\n", genre, data["count"], data["sum"], data["average"])
		}
	}

	dura := time.Since(start)
	fmt.Println("Tiempo de ejecución: %v\n", dura)

}
