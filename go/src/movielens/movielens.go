package main
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func obtenerPrimeraColumna(matriz [][]string, i int) []string {
	var primeraColumna []string

	for _, fila := range matriz {
		if len(fila) > 0 { // Verifica que la fila no esté vacía
			primeraColumna = append(primeraColumna, fila[i]) // Agrega el primer elemento de cada fila
		}
	}

	return primeraColumna
}

func elementosUnicosConConteo(lista []string) map[string]int {
	conteo := make(map[string]int)

	for _, elemento := range lista {
		conteo[elemento]++ // Incrementa el conteo de cada elemento
	}
	return conteo
}

func find_ratings(ch chan int, conteo *[][]int, reader *[][]string, names []string, m *map[string]int, s int, it int) {

	for p := 1; p < s; p++ {
		rec := (*reader)[p]
		if len(rec) == 0 {
			continue
		}
		an := strings.Split(rec[2], "|")
		for i := 0; i < len(an); i++ {
			for j, v := range names {
				if v == an[i] {
					(*conteo)[j][it] += 1 * (*m)[rec[0]]
				}
			}
		}
	}
	fmt.Println("termine")
	ch <- 1
}

func find_ratings_main(name string, n int) {
	start := time.Now()
	kg := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
		"Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
		"Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}
	nk := len(kg)

	var m = make([]map[string]int, n)

	for i := 0; i < n; i++ {
		fmt.Println(i)
		name_file := fmt.Sprintf("%s_%d.csv", "ratings", i)
		file, _ := os.Open(name_file)
		read := csv.NewReader(file)
		r, _ := read.ReadAll()
		m[i] = elementosUnicosConConteo(obtenerPrimeraColumna(r, 1))
		defer file.Close()
		r = nil
	}

	name_fil := fmt.Sprintf("%s.csv", name)
	fil, _ := os.Open(name_fil)
	reader := csv.NewReader(fil)
	data, _ := reader.ReadAll()

	s := len(obtenerPrimeraColumna(data, 0))

	var ch = make(chan int)
	var conteo = make([][]int, nk)
	for i := 0; i < nk; i++ {
		conteo[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		go find_ratings(ch, &conteo, &data, kg, &m[i], s, i)
	}

	iMsg := 0
	go func() {
		for {
			i := <-ch
			iMsg += i
		}
	}()
	for {
		if iMsg == n {
			break
		}
	}
	locVals := make([]int, nk)
	for i := 0; i < nk; i++ {
		for j := 0; j < n; j++ {
			locVals[i] += conteo[i][j]
		}
	}

	for i := 0; i < nk; i++ {
		fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", kg[i]), "  ", fmt.Sprintf("%8d", locVals[i]))
	}
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)
	println("Mt_FindRatingsMaster is Done")
}

func guarda(c chan int, data *[][]string, w *csv.Writer, it int, filas int, id int) {
	if id == 0 {
		err := w.WriteAll((*data)[(it * filas):((it + 1) * filas)])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := w.WriteAll((*data)[id:filas])
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
	c <- 1
}

func filepartition(n int, name string) {
	start := time.Now()
	nameFile := fmt.Sprintf("%s.csv", name)
	file, err := os.Open(nameFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, _ := reader.Read()

	outputFiles := make([]*os.File, n)
	writers := make([]*csv.Writer, n)
	// Crear archivos de salida y sus escritores CSV
	for i := 0; i < n; i++ {
		outputFileName := fmt.Sprintf("%s_%d.csv", name, i)
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		// Crear un escritor CSV para cada archivo de salida
		writer := csv.NewWriter(outputFile)
		writers[i] = writer

		// Escribir la cabecera en cada archivo de salida (si existe)
		if len(header) > 0 && i == 0 {
			writer.Write(header)
			writer.Flush()
		}

		outputFiles[i] = outputFile
	}

	var datos [][]string
	//headers, err := reader.Read()
	for {
		// Leer una fila (slice de strings)
		record, err := reader.Read()
		if err != nil {
			// Verificar si se alcanzó el final del archivo
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}
		//a := record[0]
		//fmt.Println(a)
		datos = append(datos, record)
	}

	tamano := len(datos)

	var chans = make(chan int)

	filas := int(tamano / n)

	for i := 0; i < n-1; i++ {
		go guarda(chans, &datos, writers[i], i, filas, 0)
	}

	go guarda(chans, &datos, writers[n-1], n-1, tamano, filas*(n-1))

	iMsg := 0
	go func() {
		for {
			i := <-chans
			iMsg += i
		}
	}()
	for {
		if iMsg == n {
			break
		}
	}

	duration := time.Since(start)
	fmt.Println("Duration = ", duration)

}
func main() {
	//filepartition(10, "ratings")
	find_ratings_main("movies", 1)
}
