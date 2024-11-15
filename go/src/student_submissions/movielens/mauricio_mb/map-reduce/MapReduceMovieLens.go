package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	dataframe "github.com/kfultz07/go-dataframe"
)

/*type MovieObj struct {
	MovieId int
	Title   string
	Genres  string
}

var artMvs []MovieObj

type RatingObj struct {
	UserId    int
	MovieId   int
	Rating    float64
	Timestamp int
}

var arRts []RatingObj
*/

const inpDir = "../File-Frag/dat/in/"
const outDir = "../File-Frag/dat/out/"

func Mt_FindRatingsMaster() {
	fmt.Println("In MtFrindRatingsMaster")
	start := time.Now()
	nf := 10 // numero de archivos con RATINGS es tambien el numero de hilos para multi-hilo

	// kg es un array que contiene los GENEROS conocidos
	kg := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}

	ng := len(kg) // numero de generos conocidos
	// ra es una matriz donde se mantienen los valores de calificación para cada género.
	// Las columnas señalan/mantienen el número central donde se ejecuta un trabajador.
	// Las filas de esa columna mantienen los valores de calificación para ese núcleo y ese género.
	ra := make([][]float64, ng)
	// ca es una matriz donde se mantiene el recuento de calificaciones para cada género.
	// Las columnas señalan el número central hacia donde se dirige el trabajador.
	// Las filas de esa columna mantienen los recuentos de ese género.
	ca := make([][]int, ng)
	//pa := make([][]float64, ng)
	// llene las ng filas de ra y ca con nf columnas
	for i := 0; i < ng; i++ {
		ra[i] = make([]float64, nf)
		ca[i] = make([]int, nf)
		//pa[i] = make([]float64, nf)
	}
	//var ci = make(chan int)
	movies := ReadMoviesCSVFile("movies.csv")
	// ejecutar FindRatings con 10 trabajadores
	for i := 0; i < nf; i++ {
		//go Mt_FindRatingsWorker(i+1, ci, kg, &ca, &ra, movies)
		Mt_FindRatingsWorker(i+1, kg, &ca, &ra, movies)
	}
	// esperar a los trabajadores
	/*iMsg := 0
	go func() {
		for {
			i := <-ci
			iMsg += i
		}
	}()
	for {
		if iMsg == 10 {
			break
		}
	}*/
	// todos los trabajadores completaron su trabajo. Recolectar resultados y producir el reporte.
	locCount := make([]int, ng)
	locVals := make([]float64, ng)
	locProm := make([]float64, ng)
	for i := 0; i < ng; i++ {
		for j := 0; j < nf; j++ {
			locCount[i] += ca[i][j]
			locVals[i] += ra[i][j]
		}
		locProm[i] = locVals[i] / float64(locCount[i])
	}
	for i := 0; i < ng; i++ {
		fmt.Println(fmt.Sprintf("%2d", i), " ", fmt.Sprintf("%20s", kg[i]), " ", fmt.Sprintf("%8d", locCount[i]), " P=", fmt.Sprintf("%20f", locProm[i]))
	}
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	println("Mt_FindRatingsMaster is Done")
}

func Mt_FindRatingsWorker(w int, kg []string, ca *[][]int, va *[][]float64, movies dataframe.DataFrame) {
	aFileName := "ratings_" + fmt.Sprintf("%02d", w) + ".csv"
	println("Worker ", fmt.Sprintf("%02d", w), " is processing file ", aFileName, "\n")

	ratings := ReadRatingsCsvFile(aFileName)
	ng := len(kg)
	start := time.Now()

	// importar todos los registros de películas DF al DF de clasificaciones, manteniendo la columna de géneros de las películas
	// df.Merge es el equivalente de una unión interna en la biblioteca DF que estoy usando aquí
	ratings.Merge(&movies, "movieId", "genres")

	// Solo necesitamos "géneros" y "calificaciones" para encontrar Count(Clasificaciones | Géneros), así que mantenga solo estas columnas
	grcs := [2]string{"genres", "rating"} // grcs => Genres Ratings Columns
	grDF := ratings.KeepColumns(grcs[:])  // grDF => Genres Ratings DF
	//fmt.Println(grDF.FrameRecords)
	for ig := 0; ig < ng; ig++ {
		for _, row := range grDF.FrameRecords {
			if strings.Contains(row.Data[0], kg[ig]) {
				(*ca)[ig][w-1] += 1
				v, _ := strconv.ParseFloat((row.Data[1]), 32) // no comprobar si hay errores
				(*va)[ig][w-1] += v
			}
		}
	}
	// for ig := 0; ig < ng; ig++ {
	// 	(*pa)[ig][w-1] = (*va)[ig][w-1] / float64((*ca)[ig][w-1])
	// }
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	fmt.Println("Worker ", w, " completed")

	// notificar al maestro que este trabajador ha completado su trabajo
	//ci <- 1
}

func ReadRatingsCsvFile(FileName string) dataframe.DataFrame {
	df := dataframe.CreateDataFrame(outDir, FileName)
	return df
}

func ReadMoviesCSVFile(FileName string) dataframe.DataFrame {
	df := dataframe.CreateDataFrame(inpDir, FileName)
	return df
}

func main() {
	Mt_FindRatingsMaster()
}
