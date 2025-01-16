package main

import (
	"encoding/csv"

	"fmt"

	"log"

	"os"

	"strconv"

	"gonum.org/v1/plot"

	"gonum.org/v1/plot/plotter"

	"gonum.org/v1/plot/vg"

	"image/color"
)

func main() {

	// Abrir el archivo CSV

	file, err := os.Open("Advertising.csv")

	if err != nil {

		log.Fatalf("Error abriendo el archivo: %v", err)

	}

	defer file.Close()

	// Leer los datos del archivo

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {

		log.Fatalf("Error leyendo el archivo: %v", err)

	}

	// Asegurarse de que hay datos

	if len(records) < 2 {

		log.Fatal("El archivo no contiene suficientes datos.")

	}

	// Mapear nombres de columnas con índices

	columns := map[string]int{

		"TV": 1,

		"Radio": 2,

		"Newspaper": 3,

		"Sales": 4,
	}

	// Generar las gráficas para cada variable vs Sales

	for variable, colIndex := range columns {

		if variable == "Sales" {

			continue // Saltar la variable dependiente

		}

		var xValues, sales []float64

		// Leer las columnas de datos (ignorando la cabecera)

		for i, record := range records {

			if i == 0 {

				// Cabecera, ignorar

				continue

			}

			xValue, err := strconv.ParseFloat(record[colIndex], 64)

			if err != nil {

				log.Printf("Error convirtiendo el valor de %s: %v", variable, err)

				continue

			}

			salesValue, err := strconv.ParseFloat(record[columns["Sales"]], 64)

			if err != nil {

				log.Printf("Error convirtiendo el valor de Sales: %v", err)

				continue

			}

			xValues = append(xValues, xValue)

			sales = append(sales, salesValue)

		}

		// Crear puntos para graficar

		points := make(plotter.XYs, len(xValues))

		for i := range xValues {

			points[i].X = xValues[i]

			points[i].Y = sales[i]

		}

		// Calcular la regresión lineal

		slope, intercept := linearRegression(xValues, sales)

		// Crear un gráfico

		p := plot.New()

		//if err != nil {

		//	log.Fatalf("Error creando el gráfico: %v", err)

		//}

		p.Title.Text = fmt.Sprintf("%s vs Sales", variable)

		p.X.Label.Text = variable

		p.Y.Label.Text = "Sales"

		// Crear los puntos y configurar el color a rojo

		scatter, err := plotter.NewScatter(points)

		if err != nil {

			log.Fatalf("Error creando puntos para el gráfico: %v", err)

		}

		scatter.GlyphStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Rojo

		p.Add(scatter)

		// Dibujar la línea de regresión

		line := plotter.NewFunction(func(x float64) float64 {

			return slope*x + intercept

		})

		line.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Azul

		line.Width = vg.Points(2)

		//line.Dashes = draw.Dashes{1, 1} // Línea sólida

		p.Add(line)

		// Guardar el gráfico como imagen

		outputFile := fmt.Sprintf("%s_vs_Sales.png", variable)

		if err := p.Save(8*vg.Inch, 6*vg.Inch, outputFile); err != nil {

			log.Fatalf("Error guardando el gráfico: %v", err)

		}

		fmt.Printf("El gráfico %s se ha generado y guardado como '%s'\n", p.Title.Text, outputFile)

	}

	fmt.Println("Todas las gráficas se han generado con éxito.")

}

// linearRegression calcula la pendiente (slope) y el intercepto para una regresión lineal simple

func linearRegression(x, y []float64) (float64, float64) {

	n := len(x)

	if n != len(y) || n == 0 {

		log.Fatal("Los datos para regresión lineal no son válidos.")

	}

	var sumX, sumY, sumXY, sumX2 float64

	for i := 0; i < n; i++ {

		sumX += x[i]

		sumY += y[i]

		sumXY += x[i] * y[i]

		sumX2 += x[i] * x[i]

	}

	// Calcular la pendiente y el intercepto

	slope := (float64(n)*sumXY - sumX*sumY) / (float64(n)*sumX2 - sumX*sumX)

	intercept := (sumY - slope*sumX) / float64(n)

	return slope, intercept

}
