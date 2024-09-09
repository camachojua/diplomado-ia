package main

import "fmt"

// creo una función que me regresa si la nueva propuesta es mejor que la de ahora
func matcher(sNewproposal, sCurrentProposal string, receiverPreference []string) bool {

	for _, preference := range receiverPreference {

		//se revisa las preferencias 1:1, si hay una mejor entonces aviso con un true

		if preference == sNewproposal {
			return true
		}
		if preference == sCurrentProposal {
			return false
		}
	}

	return false

}

// creo una función para crear slices de keys para maps
func keyStagger(mapa map[string][]string) []string {

	kSlice := []string{}
	for k := range mapa {
		kSlice = append(kSlice, k)
	}

	return kSlice
}

// keys to values: hago una función para cambiar el orden de la solución
func transposeSol(sol map[string]string) map[string]string {
	tSol := make(map[string]string)
	//para cada k,v devuelvo un mapeo invertido
	for k, v := range sol {
		tSol[v] = k
	}
	return tSol
}

func stabilizeMarriages(receivers, senders map[string][]string) map[string]string {

	//mapa para almacenar a los solteros
	singleSend := keyStagger(senders)

	//mapa para almacenar los keys de los receivers
	receiversKeys := keyStagger(receivers)

	//es un mapa para almacenar el último commit de un sender
	receiverComm := make(map[string]string)

	for _, offer := range receiversKeys {
		receiverComm[offer] = ""
	}

	for len(singleSend) > 0 {

		//tomo el primer soltero de los senders
		sender := singleSend[0]

		//actualizo el valor de la lista para darle chance de intentar a otro soltero
		if len(singleSend) != 0 {
			singleSend = singleSend[1:]
		} else {
			break
		}

		//en un loop reviso realizo el emparejamiento
		for _, proposal := range senders[sender] {

			//reviso si el sender tiene un buen match
			if matcher(sender, receiverComm[proposal], receivers[proposal]) {

				//si hace match y el valor no es el inicial, entonces el sender vuelve a la cola de singles
				if len(receiverComm[proposal]) > 0 {
					singleSend = append(singleSend, receiverComm[proposal])
				}

				//actualizo la lista de matches
				receiverComm[proposal] = sender

				//saliendo del bucle para ir con otro sender
				break

			} else {

				//si el match devuelve false entonces añado el sender a la cola de singles
				singleSend = append(singleSend, sender)
				//actualizo el registro del rank de receivers para el sender (esto protegerá su corazón)
				senders[sender] = senders[sender][1:]
				break
			}

		}

	}

	return receiverComm
}

func main() {

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

	res := transposeSol(stabilizeMarriages(receivers, senders))

	fmt.Println(res)

}
