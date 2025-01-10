package main

import (
	"fmt"
)

//Determina el indice de un elemento en partícular dentro de un slice
func indexOf(slice []string, element string) int {
	for i, v := range slice {
		if v == element {
			return i // Retorna el índice si encuentra el elemento
		}
	}
	return -1 // Retorna -1 si no encuentra el elemento
}

//Función auxiliar para encontrar el primer elemento en un slice 1, de un conjunto de elementos en un slice 2
func findFirstKing(proposals, preferences []string) string {
    for _, elem := range preferences {
        // Si el elemento está en proposals, devolverlo
        if indexOf(proposals, elem) != -1 {
            return elem
        }
    }
    return "none" // Si no se encuentra ningún elemento
}

func main() {
	//Senders (aka kings)
	var senders = map[string][]string{
		"s01": {"r01", "r05", "r03", "r09", "r10", "r04", "r06", "r02", "r08", "r07"},
		"s02": {"r03", "r08", "r01", "r04", "r05", "r06", "r02", "r10", "r09", "r07"},
		"s03": {"r08", "r05", "r01", "r04", "r02", "r06", "r09", "r07", "r03", "r10"},
		"s04": {"r09", "r06", "r04", "r07", "r08", "r05", "r10", "r02", "r03", "r01"},
		"s05": {"r10", "r04", "r02", "r03", "r06", "r05", "r01", "r09", "r08", "r07"},
		"s06": {"r02", "r01", "r04", "r07", "r05", "r09", "r03", "r10", "r08", "r06"},
		"s07": {"r07", "r05", "r09", "r02", "r03", "r01", "r04", "r08", "r10", "r06"},
		"s08": {"r01", "r05", "r08", "r06", "r09", "r03", "r10", "r02", "r07", "r04"},
		"s09": {"r08", "r03", "r04", "r07", "r02", "r01", "r06", "r09", "r10", "r05"},
		"s10": {"r01", "r06", "r10", "r07", "r05", "r02", "r04", "r03", "r09", "r08"},
	}

	//Receivers (aka queens)
	var receivers = map[string][]string{
		"r01": {"s02", "s06", "s10", "s07", "s09", "s01", "s04", "s05", "s03", "s08"},
		"r02": {"s02", "s01", "s03", "s06", "s07", "s04", "s09", "s05", "s10", "s08"},
		"r03": {"s06", "s02", "s05", "s07", "s08", "s03", "s09", "s01", "s04", "s10"},
		"r04": {"s06", "s10", "s03", "s01", "s09", "s08", "s07", "s04", "s02", "s05"},
		"r05": {"s10", "s08", "s06", "s04", "s01", "s07", "s03", "s05", "s09", "s02"},
		"r06": {"s02", "s01", "s05", "s09", "s10", "s04", "s06", "s07", "s03", "s08"},
		"r07": {"s10", "s07", "s08", "s06", "s02", "s01", "s03", "s05", "s04", "s09"},
		"r08": {"s07", "s10", "s02", "s01", "s09", "s04", "s08", "s05", "s03", "s06"},
		"r09": {"s09", "s03", "s08", "s07", "s06", "s02", "s01", "s05", "s10", "s04"},
		"r10": {"s05", "s08", "s07", "s01", "s02", "s10", "s03", "s09", "s06", "s04"},
	}

	//Este map almacena los posibles matches de cada Queen.
	//En un principio cada Queen no tiene match, lo que se representa con un slice vacío.
	var matches_Q = map[string][]string{
		"r01": {},
		"r02": {},
		"r03": {},
		"r04": {},
		"r05": {},
		"r06": {},
		"r07": {},
		"r08": {},
		"r09": {},
		"r10": {},
	}

	//Este map indica si cada sender está actualmente emparejado con alguién
	var matches_K = map[string]string{
		"s01": "unmatch",
		"s02": "unmatch",
		"s03": "unmatch",
		"s04": "unmatch",
		"s05": "unmatch",
		"s06": "unmatch",
		"s07": "unmatch",
		"s08": "unmatch",
		"s09": "unmatch",
		"s10": "unmatch",
	}

	for {
		for i := 0; i < 10; i++ {
			//Etapa de propuestas
			for king_i, s0i := range senders {
				//si king_i no está emparejado con alguién, entonces lanza propuesta
				if matches_K[king_i] == "unmatch" {
					queen_j := s0i[i] //es a quien le propone match
					matches_Q[queen_j] = append(matches_Q[queen_j], king_i) //se agrega a king_i en la lista de pretendientes de queen_j
				}
			}
			//Etapa de evaluación
			for queen_j, kings := range matches_Q {
				//Si queen_j solo recibió una propuesta, se empareja y el status del king cambia a match
				if len(kings) == 1{
					matches_K[kings[0]] = "match"
				} else if len(kings) > 1 {
					//Se compara a los pretendientes y hace match con el mejor posicionado segun su slice de preferencias
					king_j := findFirstKing(kings, receivers[queen_j])
					matches_K[kings[0]] = "unmatch" //El actual match de queen se desmarca
					matches_Q[queen_j] = []string{king_j} //Se asigana el nuevo match
					matches_K[king_j] = "match" //El status del king se cambia
				}
			}
		}

		//¿Todos los senders están ya emparejados?
		n_pairs := 0
		for _, match := range matches_K {
			if match == "match" {
				n_pairs ++
			}
		}
		if n_pairs == 10 {
			break
		} 
	}
	fmt.Println("Sender / Receiver")
	for queen, king := range matches_Q {
		fmt.Println(king[0], queen)
	}
}

/* 
Sender / Receiver
s03 r04
s08 r05
s09 r08
s05 r10
s01 r09
s10 r01
s06 r02
s02 r03
s04 r06
s07 r07 
*/