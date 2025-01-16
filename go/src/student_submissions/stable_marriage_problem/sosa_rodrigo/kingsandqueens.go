package main

import "fmt"

// Función para  encontrar un emparejamiento estable usando el algorito Gale-Shapley

func stableMarriage(senders, receivers map[string][]string) map[string]string {
	// Almacena los emparejamientos actuales
	matches := make(map[string]string)
	// Rastreo de quién esta libre
	freeSenders := make([]string, 0)

	// Rastreo del próximo receiver al que cada sender se propondrá
	nextProposal := make(map[string]int)
	for sender := range senders {
		freeSenders = append(freeSenders, sender)
		nextProposal[sender] = 0
	}

	// Invierte la preferencia de los recerivers para facilitar la comparación
	rank := make(map[string]map[string]int)
	for receiver, prefs := range receivers {
		rank[receiver] = make(map[string]int)
		for i, sender := range prefs {
			rank[receiver][sender] = i
		}
	}

	for len(freeSenders) > 0 {
		// El sender escoge el siguiente receiver en su lista de preferncia
		sender := freeSenders[0]
		receiver := senders[sender][nextProposal[sender]]
		nextProposal[sender]++

		// Verifica si el recibidor ya está emparejado
		currentMatch, matched := matches[receiver]

		if !matched {
			// Si el receiver no está emparejado, lo empareja con el sender
			matches[receiver] = sender
			// Elimina al sender de la lista de senders libres
			freeSenders = freeSenders[1:]
		} else {
			// Si el receiver está emparejado, compara el sender con el emparejamiento actual
			if rank[receiver][sender] < rank[receiver][currentMatch] {
				// Si el nuevo sender es preferido, se emparejan y libera al sender previo
				matches[receiver] = sender
				freeSenders = freeSenders[1:]
				freeSenders = append(freeSenders, currentMatch)
			}
		}
	}

	//Regresa los emparejamientos
	return matches
}

func main() {

	// Preferencias de senders

	senders := map[string][]string{
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

	// Preferencias para receivers

	receivers := map[string][]string{
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

	//Encontrar los e mparejamientos estables

	matches := stableMarriage(senders, receivers)

	// Imprime las parejas
	fmt.Println("Stable Matches:")
	for receiver, sender := range matches {
		fmt.Printf("%s is matched with %s\n", sender, receiver)
	}

}
