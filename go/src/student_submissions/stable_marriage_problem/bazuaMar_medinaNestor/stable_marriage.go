package main

import (
	"fmt"
)

type stable_marriage struct {
	Receiver string
	Sender   string
}

// Las preferencias de los receivers y senders
func main() {
	receivers := map[string][10]string{
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

	senders := map[string][10]string{
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

	// Lista de estado de los receivers y senders
	receivers_status := make(map[string]string)
	senders_status := make(map[string]string)
	for r := range receivers {
		receivers_status[r] = "single"
	}
	for s := range senders {
		senders_status[s] = "single"
	}

	// Función auxiliar para encontrar el índice de un elemento en un array
	findIndex := func(arr [10]string, value string) int {
		for i, v := range arr {
			if v == value {
				return i
			}
		}
		return -1
	}

	// Lista de emparejamientos
	marriageList := make(map[string]string)

	// Mientras haya senders solteros
	for {
		allPaired := true
		for s, status := range senders_status {
			if status == "single" {
				allPaired = false
				// El sender hace una propuesta a su preferencia más alta disponible
				for _, r := range senders[s] {
					// Si el receiver está soltero, emparejar
					if receivers_status[r] == "single" {
						marriageList[s] = r
						receivers_status[r] = "engaged"
						senders_status[s] = "engaged"
						break
					} else {
						// Si el receiver está comprometido, revisar si prefiere al nuevo sender
						current_sender := ""
						for key, val := range marriageList {
							if val == r {
								current_sender = key
								break
							}
						}

						// Si el receiver prefiere al nuevo sender
						if findIndex(receivers[r], s) < findIndex(receivers[r], current_sender) {
							// Romper el emparejamiento anterior y emparejar con el nuevo sender
							senders_status[current_sender] = "single"
							marriageList[s] = r
							senders_status[s] = "engaged"
							break
						}
					}
				}
			}
		}

		if allPaired {
			break
		}
	}

	// Imprimir los emparejamientos finales
	fmt.Println("Emparejamientos finales:")
	for sender, receiver := range marriageList {
		fmt.Printf("%s está emparejado con %s\n", sender, receiver)
		fmt.Printf("Las preferencias de sender %s son: %s \n", sender, senders[sender])
		fmt.Printf("Las preferencias de receiver %s son: %s \n", receiver, receivers[receiver])
		fmt.Printf("\n")
	}
}
