package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func limpiarTexto(texto string) string {
	// Reemplazar acentos por letras sin acento
	texto = strings.ToLower(texto)
	reemplazos := map[string]string{
		"á": "a", "é": "e", "í": "i", "ó": "o", "ú": "u",
		"Á": "A", "É": "E", "Í": "I", "Ó": "O", "Ú": "U",
		"ñ": "n", "Ñ": "N",
	}
	for acento, sinAcento := range reemplazos {
		texto = strings.ReplaceAll(texto, acento, sinAcento)
	}

	// Eliminar comas, comillas y puntos
	reg := regexp.MustCompile(`[[:punct:]\d¿¡]`)

	texto = reg.ReplaceAllString(texto, "")

	// Eliminar espacios extra
	texto = strings.TrimSpace(texto)

	return texto
}

func main() {

	//start := time.Now()

	name := "LM"
	start := time.Now()
	nameFile := fmt.Sprintf("%s.csv", name)
	file, err := os.Open(nameFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	reader.FieldsPerRecord = -1
	//reader.LazyQuotes = true

	var datos [][]string
	paltot := 0
	conteo := make(map[string]int)

	for {
		// Leer una fila
		record, err := reader.Read()
		if err != nil {
			// Verificar si se alcanzó el final del archivo
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}
		// Limpiar y filtrar celdas vacías
		var filaLimpia []string
		for _, celda := range record {
			textoLimpio := limpiarTexto(celda)
			if textoLimpio != "" { // Solo agregar si no está vacío
				conteo[textoLimpio]++
				paltot++
				filaLimpia = append(filaLimpia, textoLimpio)
			}
		}

		// Agregar la fila limpia solo si no está vacía
		if len(filaLimpia) > 0 {
			datos = append(datos, filaLimpia)
		}
	}
	guarda(&datos)

	// Convertir el mapa en un slice de pares clave-valor
	var pares []struct {
		clave string
		valor int
	}

	for k, v := range conteo {
		pares = append(pares, struct {
			clave string
			valor int
		}{k, v})
	}

	for k, v := range conteo {
		if v == 1 {
			fmt.Println(k)
		}
	}
	// Ordenar de mayor a menor según el valor
	sort.Slice(pares, func(i, j int) bool {
		return pares[i].valor > pares[j].valor
	})

	fmt.Println("Top 100 valores más grandes:")
	for i := 0; i < len(pares) && i < 100; i++ {
		fmt.Printf("%d) %s: %d\n", i+1, pares[i].clave, pares[i].valor)
	}

	fmt.Println("Ultimos 100 valores menos grandes:")
	for i := 0; i < len(pares) && i < 100; i++ {
		fmt.Printf("%d) %s: %d\n", i+1, pares[len(pares)-i-1].clave, pares[len(pares)-i-1].valor)
	}

	fmt.Println("Numero de palabras y numeros distintas en el texto:", len(pares))
	fmt.Println("Numero TOTAL de palabras en el texto:", paltot)
	duration := time.Since(start)
	fmt.Println("Duration = ", duration)

}

func guarda(datos *[][]string) {

	nuevoArchivo := "LM_limpio.csv"
	outputFile, err := os.Create(nuevoArchivo)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	err = writer.WriteAll(*datos)
	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}
