/* *****************************************************
 *
 * Diplomado Inteligencia Artificial y Ciencia de Datos
 * 		07 de septiembre de 2024
 *
 *  Code Challenge:
 *		Implementar el algoritmo de Gale y Shapley para
 * 		el Stable Matching Problem.
 *
 *	Nuestro programa asume que existen los archivos
 *  data/kings.csv y data/queens.csv relativos al
 *	ejecutable.
 *
 * 	Compilado:
 *			go build
 *  Ejecución:
 *			./smp
 *  Ejecución directa:
 *			go run smp.go
 *
 *	Equipo:
 * 		Guadalupe Villeda Gómez
 * 			<lupis_act@ciencias.unam.mx>
 * 		Ivan Pompa-García
 *			<ivanjpg@ekbt.nl>
 *
 * *****************************************************
 */

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
)

// Definimos la estructura Rating que sirve de base para
// la escritura hacia parquet.
type Rating struct {
	UserID    int32   `parquet:"name=user_id, type=INT32"`
	MovieID   int32   `parquet:"name=movie_id, type=INT32"`
	Rating    float32 `parquet:"name=rating, type=FLOAT"`
	Timestamp int64   `parquet:"name=timestamp, type=INT64"`
}

// La función perror revisa si el
// tipo de dato del argumento que recibe.
// Si es de tipo «error» se termina la ejecución
// de todo el programa.
// Si es de tipo string, se genera un tipo error a
// partir de él y se termina el programa.
// En otro caso se asume que no hay error y la función
// simplemente retorna.
// El manejo de errores puede ser más fino, pero se
// deja así por simplicidad.

func perror(err interface{}) {
	switch v := err.(type) {
	case error:
		panic(err)
	case string:
		panic(errors.New(v))
	default:
		return
	}
}

func readCSV(path string) [][]string {
	fmt.Printf("Leyendo el archivo *%s*...\n", path)

	// Abrimos el CSV.
	csvFile, err := os.Open(path)
	perror(err)

	// Generamos un nuevo lector de CSV a partir del
	// archivo abierto.
	reader := csv.NewReader(csvFile)

	// Evitamos la verificación del número de columnas
	// de cada registro del CSV. Cada columna puede tener
	// el número de registros que sea.
	reader.FieldsPerRecord = -1

	// Leemos todo el contenido del CSV.
	rows, err := reader.ReadAll()
	perror(err)

	// Ya hemos leído los datos, no es necesario dejar
	// abierto el archivo.
	csvFile.Close()

	return rows
}

// Imprime los mapas de preferencias.
func printPrefs(prefs map[string][]string) {
	for i, plist := range prefs {
		fmt.Printf("[%s] -> [", i)

		for _, v := range plist {
			fmt.Printf(" %s ", v)
		}

		fmt.Println("]")
	}
}

// Toma las filas de preferencias y las asigna al mapa
// de strings de Go.
func prefsFromCSV(rows [][]string, prefs map[string][]string) {
	for _, row := range rows {
		var rowKey string

		for j, col := range row {
			if j == 0 {
				rowKey = col
			} else {
				prefs[rowKey] = append(prefs[rowKey], col)
			}
		}
	}
}

// Imprime el mapa de matches.
func printMatches(matches map[string]string) {
	fmt.Println("\n==== MATCHES ====")

	for king, queen := range matches {
		fmt.Printf("[%s,%s]\n", king, queen)
	}

	fmt.Println()
}

// Verifica si dentro de un map[string]string existe
// un valor «val» dado.
func isValInMatches(matches map[string]string, val string) bool {
	for _, v := range matches {
		if v == val {
			return true
		}
	}

	return false
}

// Verifica si hay un rey soltero, es decir, si
// existe alguna cadena vacía ("") en los valores
// del map «matches».
func areStillSingles(matches map[string]string) bool {
	return isValInMatches(matches, "")
}

// Verifica si el rey «king» está en los valores del
// map «matches», es decir, si tiene pareja.
func isKingMatched(matches map[string]string, king string) bool {
	return isValInMatches(matches, king)
}

// Verificamos si un match es inestable.
func isMatchUnstable(prefsQueens map[string][]string, prefsKings map[string][]string, matches map[string]string, matchedQueen string, matchedKing string) bool {
	// Obtenemos la lista de preferencias del rey
	// de este emparejamiento.
	kingPrefs := prefsKings[matchedKing]

	// Obtenemos el "ranking" de la reina emparejada en
	// las preferencias del rey.
	idxQueenMatched := slices.Index(kingPrefs, matchedQueen)

	// Recorremos la lista de preferencias del rey
	// emparejado.
	for idxQueenCandidate, queenCandidate := range kingPrefs {
		// Si la reina candidata es la misma reina con la que
		// está emparejado, continuamos con la siguiente reina.
		if idxQueenCandidate == idxQueenMatched {
			continue
		}

		// Obtenemos el ranking del rey emparejado de acuerdo
		// a la lista de preferencias de la reina candidata.
		idxKingInQueenCandidate := slices.Index(prefsQueens[queenCandidate], matchedKing)

		// Obtenemos el rey con el que está emparejada la
		// reina candidata.
		matchedKingOfQueenCandidate := matches[queenCandidate]

		// Obtenemos el raking del rey con el que está
		// emparejado la reina candidata.
		idxKingMatchedOfQueenCandidate := slices.Index(prefsQueens[queenCandidate], matchedKingOfQueenCandidate)

		candidateQueenPrefersAnotherMatch := idxKingInQueenCandidate < idxKingMatchedOfQueenCandidate
		currentKingPrefersAnotherMatch := idxQueenCandidate < idxQueenMatched

		return candidateQueenPrefersAnotherMatch && currentKingPrefersAnotherMatch
	}

	return false
}

func main() {
	// Definimos los maps para guardar las listas
	// de preferencias de reinas y reyes.
	prefsQueens := make(map[string][]string)
	prefsKings := make(map[string][]string)

	// Definimos un map para guardar los matches.
	// La clave es la reina y el valor es el rey
	// con el que está emparejada.
	matches := make(map[string]string)

	// Leemos el archivo queens
	rowsQueens := readCSV("data/queens.csv")
	// Informamos cuántos registros leímos.
	fmt.Println("El archivo queens tiene", len(rowsQueens), "registros.")

	// Leemos el archivo kings
	rowsKings := readCSV("data/kings.csv")
	// Informamos cuántos registros leímos.
	fmt.Println("El archivo kings tiene", len(rowsKings), "registros.")

	// Verificamos que exista el mismo número de
	// reinas y reyes.
	if len(rowsQueens) != len(rowsKings) {
		perror("¡El número de reinas y reyes es distinto!")
	}

	// Llenamos los mapas de preferencias con la información
	// proveniente de los CSVs.
	prefsFromCSV(rowsQueens, prefsQueens)
	prefsFromCSV(rowsKings, prefsKings)

	// Imprimimos las preferencias.
	fmt.Println("==== prefsQueens ====")
	printPrefs(prefsQueens)

	fmt.Println("\n==== prefsKings ====")
	printPrefs(prefsKings)

	// Inicializamos la estructura «matches» con las
	// claves de las reinas.
	for queen := range prefsQueens {
		matches[queen] = ""
	}

	// Imprimimos la estructura de matches inicial.
	printMatches(matches)

	// Iniciamos el algoritmo de SMP

	for areStillSingles(matches) {
		// De cada rey, obtenemos su «nombre» (king) y su
		// lista de reinas en orden de preferencia (queens).
		for king, queens := range prefsKings {
			if isKingMatched(matches, king) {
				continue
			}

			// Recorremos cada reina (candidateQueen) en la
			// lista de preferencias del rey (king).
			for _, candidateQueen := range queens {
				// Buscamos al rey (kingMatch) que está emparejado
				// con la reina actual (candidateQueen)
				kingMatch := matches[candidateQueen]

				// Si NO hay rey emparejado (kingMatch) con la
				// reina actual (candidateQueen), entonces los
				// emparejamos.
				if kingMatch == "" {
					matches[candidateQueen] = king

					// Como logramos el emparejamiento, dejamos
					// de iterar sobre las reinas candidatas para
					// pasar al siguiente rey.
					break
				}

				// Si la reina candidata ya tenía una
				// pareja (kingMatch)...

				// Obtenemos el "ranking" del rey
				// pretendiente (king) y del rey pareja (kingMatch).
				idxMatch := slices.Index(prefsQueens[candidateQueen], kingMatch)
				idxKing := slices.Index(prefsQueens[candidateQueen], king)

				// Si el rey pretendiente tiene un mejor ranking
				// que el rey pareja, entonces emparejamos a la
				// reina candidata con el pretendiente.
				if idxKing < idxMatch {
					matches[candidateQueen] = king

					// Como logramos el emparejamiento, dejamos
					// de iterar sobre las reinas candidatas para
					// pasar al siguiente rey.
					break
				}

				// Si llegamos a este punto, quiere decir que el
				// rey pareja es el mejor para la reina candidata.
				// Entonces, se procede a verificar la siguiente
				// reina.
			}
		}
	}

	printMatches(matches)

	// Verificamos la estabilidad de los matches.
	for queen, king := range matches {
		isUnstable := isMatchUnstable(prefsKings, prefsQueens, matches, queen, king)

		if isUnstable {
			fmt.Printf("[%s,%s] es inestable.\n", queen, king)
		}
	}
}
