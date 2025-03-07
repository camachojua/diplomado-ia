package main

import (
	"fmt"
	"time"
	"sync"
	"strings"
	"io/ioutil"
	"strconv"
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

func reduce_ratings(resultados []map[string]*[2]float64) {
	generos_resultado := make(map[string]*[2]float64)
	for _, generos_peli := range resultados {
		for gen, cont := range generos_peli {
			_, exists := generos_resultado[gen]
			if !exists {generos_resultado[gen] = &[2]float64{0, 0}}
			
			generos_resultado[gen][0] = generos_resultado[gen][0] + cont[0]
			
			generos_resultado[gen][1] = generos_resultado[gen][1] + cont[1]
		}
	}
	for gen, cont := range generos_resultado {
		fmt.Printf("%s: %d - %.2f\n", gen, int(cont[0]), cont[1]/cont[0])
	}
}

func map_ratings(num_partes int) {

	var resultados []map[string]*[2]float64
	
	tiempo_inicial := time.Now()
	archivo_ratings, _ := ioutil.ReadFile("./test/ratings.csv")
	ratings := strings.Split(string(archivo_ratings), "\n")
	num_ratings := len(ratings)
	peliculas := csv2map_pelis()
	
	c := make(chan int)
	mutex := sync.RWMutex{}
	for i := 1; i <= num_partes; i++ {
		parte_medida := num_ratings / num_partes 
		fin_parte := parte_medida * i
		if fin_parte > num_ratings { fin_parte = num_ratings }
		parte_data := ratings[(i-1)*parte_medida + 1:fin_parte]
		if len(parte_data)== 0 { break }
		go func() {
			generos := make(map[string]*[2]float64)
			for _, peli := range parte_data {
				if peli == "" {break}
				movie_id := strings.Split(strings.TrimSpace(peli), ",")[1]
				rating := strings.Split(strings.TrimSpace(peli), ",")[2]
				generos_pelicula := strings.TrimSpace(peliculas[movie_id])
				for _, genero := range strings.Split(generos_pelicula, "|") {
					_, exists := generos[genero]
					if !exists {generos[genero] = &[2]float64{0, 0}}
					generos[genero][0] = generos[genero][0] + 1
					prm, _ := strconv.ParseFloat(rating, 64)
					generos[genero][1] = generos[genero][1] + float64(prm)
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
