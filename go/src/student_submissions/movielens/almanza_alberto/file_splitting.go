package main

import (
	"fmt"
	"time"
	"sync"
	"strings"
	"io/ioutil"
)

type Pelicula struct {
	MovieId string
	Genero  string
}
type Peliculas []Pelicula

func csv2map_pelis() map[string]string {
	dta, _ := ioutil.ReadFile("./test/movies.csv")	
	data := strings.Split(string(dta), "\n")
	pelis := make(map[string]string)
	for _, record := range data {
		if record == "" { break }
		d_split := strings.Split(record, ",")
		pelis[strings.TrimSpace(d_split[0])] = strings.TrimSpace(d_split[len(d_split)-1])
	}
	return pelis
}

func reduce_ratings(resultados []map[string]int) {
	generos_resultado := make(map[string]int)
	for _, generos_peli := range resultados {
		for gen, cont := range generos_peli {
			v, exists := generos_resultado[gen]
			if !exists { generos_resultado[gen] = 0 }
			generos_resultado[gen] = v + cont
		}
	}
	for gen, cont := range generos_resultado {
		fmt.Printf("%s: %d\n", gen, cont)
	}
}

func map_ratings(num_partes int) {

	var resultados []map[string]int
	
	tiempo_inicial := time.Now()
	archivo_ratings, _ := ioutil.ReadFile("./test/ratings.csv")
	ratings := strings.Split(string(archivo_ratings), "\n")
	num_ratings := len(ratings)
	peliculas := csv2map_pelis()
	
	c := make(chan int)
	mutex := sync.RWMutex{}
	for i := 0; i <= num_partes; i++ {
		parte_medida := num_ratings / num_partes 
		fin_parte := parte_medida * (i + 1)
		if fin_parte > num_ratings { fin_parte = num_ratings }
		parte_data := ratings[i*parte_medida:fin_parte]
		if len(parte_data)== 0 { break }
		go func() {
			generos := make(map[string]int)
			for _, peli := range parte_data {
				if peli == "" {break}
				movie_id := strings.Split(strings.TrimSpace(peli), ",")[1]
				generos_pelicula := strings.TrimSpace(peliculas[movie_id])
				for _, genero := range strings.Split(generos_pelicula, "|") {
					v, exists := generos[genero]
					if !exists {generos[genero] = 0}
					generos[genero] = v + 1
				}
			}
			mutex.Lock()
			resultados = append(resultados, generos)
			mutex.Unlock()
			c <-1
		}()
	}
	for i := 0; i<num_partes; i++ { <-c }

	fmt.Printf("\nArchivo ratings.csv con %d registros\n", num_ratings)
	fmt.Printf("Número de partes: %d\n", num_partes)
	fmt.Printf("Archivo movies.csv con %d registros\n", len(peliculas))
	fmt.Printf("Tiempo transcurrido: %f\n\n", time.Now().Sub(tiempo_inicial).Seconds())
	reduce_ratings(resultados)
}

func main() {
	map_ratings(10)
}
