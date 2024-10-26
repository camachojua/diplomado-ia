//Authors
//Mar Bazúa & Néstor Medina

package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"golang.org/x/image/colornames"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Lee el CSV y lo convierte en un slice de slices de strings para que sea parecido a un dataframe
func readCSV(inputDirectory string, fileName string) [][]string {
	file, err := os.Open(inputDirectory + fileName)
	if err != nil {
		log.Fatalf("Error al abrir el archivo CSV: %v", err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatalf("Error al leer el archivo CSV: %v", err)
	}

	return records
}

// Convierte las columnas del CSV a valores flotantes
func getColumns(data [][]string, xCol, yCol int) plotter.XYs {
	pts := make(plotter.XYs, len(data)-1) // Ignorar la primera fila (headers)
	for i := 1; i < len(data); i++ {
		x, _ := strconv.ParseFloat(data[i][xCol], 64)
		y, _ := strconv.ParseFloat(data[i][yCol], 64)
		pts[i-1].X = x
		pts[i-1].Y = y
	}
	return pts
}

// Calcular la regresión lineal a manita
func linearRegression(pts plotter.XYs) (slope, intercept float64) {
	var sumX, sumY, sumXY, sumX2 float64
	n := float64(len(pts))
	for _, pt := range pts {
		sumX += pt.X
		sumY += pt.Y
		sumXY += pt.X * pt.Y
		sumX2 += pt.X * pt.X
	}
	slope = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept = (sumY - slope*sumX) / n
	return
}

// Genera la línea de regresión aplicando pendiente y ordenada al origen
func regressionLine(pts plotter.XYs, slope, intercept float64) plotter.XYs {
	line := make(plotter.XYs, 2)
	minX, maxX := pts[0].X, pts[0].X
	for _, pt := range pts {
		if pt.X < minX {
			minX = pt.X
		}
		if pt.X > maxX {
			maxX = pt.X
		}
	}
	line[0].X = minX
	line[0].Y = slope*minX + intercept
	line[1].X = maxX
	line[1].Y = slope*maxX + intercept
	return line
}

// Funcion para crear la gráfica una por una en el formato deseado
func createPlot(title, xLabel, yLabel string, pts plotter.XYs, slope, intercept float64) *plot.Plot {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	// Gráfico de dispersión
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatalf("Error al crear scatter plot: %v", err)
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Rojo
	p.Add(scatter)

	// Línea de regresión
	line := regressionLine(pts, slope, intercept)
	regression, err := plotter.NewLine(line)
	if err != nil {
		log.Fatalf("Error al crear la línea de regresión: %v", err)
	}
	regression.LineStyle.Color = colornames.Blue // Azul
	regression.LineStyle.Width = vg.Points(2)    // Grosor
	p.Add(regression)

	return p
}

func main() {
	// Leer los datos del CSV
	data := readCSV("./", "Advertising.csv")

	// Variables (TV, Radio, Newspaper vs Sales)
	variables := []struct {
		xCol   int
		yCol   int
		xLabel string
		yLabel string
	}{
		{1, 4, "TV", "Sales"},
		{2, 4, "Radio", "Sales"},
		{3, 4, "Newspaper", "Sales"},
	}

	// Guardar cada gráfico
	for i, v := range variables {
		pts := getColumns(data, v.xCol, v.yCol)
		slope, intercept := linearRegression(pts)
		p := createPlot(fmt.Sprintf("%s vs %s", v.xLabel, v.yLabel), v.xLabel, v.yLabel, pts, slope, intercept)

		imagePath := fmt.Sprintf("./plots/plot_%d.png", i+1)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, imagePath); err != nil {
			log.Fatalf("Error al guardar la gráfica: %v", err)
		}
	}

	fmt.Println("Gráficas creadas, combinadas y guardadas exitosamente.")
}
