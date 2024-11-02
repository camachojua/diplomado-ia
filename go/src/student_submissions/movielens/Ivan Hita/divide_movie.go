package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

func Dividir(nombreArchivo string, numeroParticiones int, directorio string) ([]string, error) {
	rutaArchivo := fmt.Sprintf("%s/%s", directorio, nombreArchivo)
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		return nil, fmt.Errorf("el directorio no existe: %s", directorio)
	}

	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer archivo.Close()

	lectorCSV := csv.NewReader(archivo)
	particiones := make(map[string][][]string)
	encabezados, err := lectorCSV.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("archivo vacío: %s", nombreArchivo)
		}
		return nil, fmt.Errorf("no se pudieron leer los encabezados: %w", err)
	}

	numRegistro := 0
	for {
		registro, err := lectorCSV.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("no se pudo leer el registro: %w", err)
		}
		procesar(numeroParticiones, numRegistro, registro, particiones)
		numRegistro++
	}

	var wg sync.WaitGroup
	for i := 0; i < len(particiones); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := guardarParticion(rutaArchivo, i, particiones, encabezados); err != nil {
				log.Printf("no se pudo guardar la partición %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	var archivosParticiones []string
	for clave := range particiones {
		nombreArchivoParticion := fmt.Sprintf("%s_particion_%s.csv", nombreArchivo, clave)
		archivosParticiones = append(archivosParticiones, nombreArchivoParticion)
	}
	return archivosParticiones, nil
}

func procesar(numeroParticiones, numRegistro int, registro []string, particiones map[string][][]string) {
	particion := strconv.Itoa(numRegistro % numeroParticiones + 1)
	particiones[particion] = append(particiones[particion], registro)
}

func guardarParticion(rutaArchivo string, numTrabajador int, particiones map[string][][]string, encabezados []string) error {
	nombreArchivoParticion := fmt.Sprintf("%s_particion_%d.csv", rutaArchivo, numTrabajador+1)
	particion := strconv.Itoa(numTrabajador + 1)
	filasParticion := particiones[particion]

	filasParticionConEncabezados := append([][]string{encabezados}, filasParticion...)
	archivoCSV, err := os.Create(nombreArchivoParticion)
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo %s: %w", nombreArchivoParticion, err)
	}
	defer archivoCSV.Close()

	escritorCSV := csv.NewWriter(archivoCSV)
	if err := escritorCSV.WriteAll(filasParticionConEncabezados); err != nil {
		return fmt.Errorf("no se pudo escribir en el archivo %s: %w", nombreArchivoParticion, err)
	}
	escritorCSV.Flush()
	return nil
}
