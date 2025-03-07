package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"log"
	"os"
	"strings"
	"time"
	"io/ioutil"
)

type Pelicula struct {
	MovieId string
	Genero  string
}

type Peliculas []Pelicula

func carga_peliculas() Peliculas {
	file, err := os.Open("movies.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	csv_file := csv.NewReader(file)
	csv_file.FieldsPerRecord = -1
	data, err := csv_file.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var peliculas = Peliculas{}
	for _, record := range data {
		for _, gen := range strings.Split(record[2], "|") {
			peliculas = append(peliculas,
				Pelicula{MovieId: strings.TrimSpace(record[0]),
					Genero: strings.TrimSpace(gen)})
		}
	}	
	return peliculas
}

func main() {
	tiempo_inicial := time.Now()
	
	archivo_ratings, _ := ioutil.ReadFile("ratings.csv")
	ratings := strings.Split(string(archivo_ratings), "\n")

	f1, err := os.Create("ratings1.csv")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	ratings_1 := ratings[:500]
	f1.WriteString(strings.Join(ratings_1, "\n"))

	file1, err := os.Open("ratings1.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()
	fmt.Println("lenf1:", file1)

	
	f2, err := os.Create("ratings2.csv")
	if err != nil {
		panic(err)
	}
	defer f2.Close()
	ratings_2 := ratings[500:]
	ratings_2h := "userId,MovieId,rating,timestamp\n" + strings.Join(ratings_2, "\n")
	f2.WriteString(ratings_2h)
	file2, err := os.Open("ratings2.csv")	
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()
	fmt.Println("lenf2:", file2)
	
	df_peliculas := dataframe.LoadStructs(carga_peliculas())
	fmt.Println("\nDataframe de movies.csv: ", df_peliculas)	

	c := make(chan int)
	go func() {

		df_ratings1 := dataframe.ReadCSV(file1)
		fmt.Println("DataFrame de raitings1.csv:", df_ratings1)
	
		df_join := df_ratings1.InnerJoin(df_peliculas, "MovieId")
		fmt.Println("DataFrame InnerJoin [movies-ratings]:", df_join)
		
		df_groupby := df_join.GroupBy("Genero")
		
		df_count := df_groupby.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_COUNT},
			[]string{"Genero"})
		fmt.Println("Dataframe final: ", df_count)
		c <- 1
	} ()
	go func() {

		df_ratings2 := dataframe.ReadCSV(file2)
		fmt.Println("DataFrame de raitings2.csv:", df_ratings2)
	
		df_join := df_ratings2.InnerJoin(df_peliculas, "MovieId")
		fmt.Println("DataFrame InnerJoin [movies-ratings]:", df_join)
		
		df_groupby := df_join.GroupBy("Genero")
		
		df_count := df_groupby.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_COUNT},
			[]string{"Genero"})
		fmt.Println("Dataframe final: ", df_count)
		c <- 1
	} ()
	
	for i := 0; i<2; i++ { <-c }
	
	fmt.Printf("Tiempo transcurrido: %f\n\n", time.Now().Sub(tiempo_inicial).Seconds())
}