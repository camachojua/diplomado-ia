package main

import (
	"fmt"
	"encoding/csv"
	"image/color"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

//Esta función devuelve los plotter asociados a los datos (puntos) y a su regresión lineal
//Toma por argumentos el nombre del archivo CSV, y el número de las dos columnas/variables  
func scatter_lr_plotter(csv_file string, x_data int, y_data int) (*plotter.Scatter, *plotter.Function) {
	//Se abre el archivo
	file, err := os.Open(csv_file)
	if err != nil {
		//Muestra el error y termina el programa si no se puede leer el archivo CSV
		log.Fatal("No se pudo abrir el archivo CSV: ", err)
	}
	defer file.Close()

	//Se lee el archivo
	reader := csv.NewReader(file)
	//se omite la primera linea (encabezados)
	if _, err := reader.Read(); err != nil {
		panic(err)
	}
	CSV_data, err := reader.ReadAll()
	if err != nil {
		log.Fatal("No se pudo leer el archivo CSV: ", err)
	}

	//Se convierten los datos del CSV a plotter.values
	data := make(plotter.XYs, len(CSV_data))
	var data_x, data_y []float64
	for i, record := range CSV_data {
		x, err := strconv.ParseFloat(record[x_data], 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseFloat(record[y_data], 64)
		if err != nil {
			panic(err)
		}
		data[i].X = x
		data[i].Y = y
		data_x = append(data_x, x)
		data_y = append(data_y, y)
	}

	//Scatter plot
	scatter, err := plotter.NewScatter(data)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(2)

	//Regresión lineal
	b, a := stat.LinearRegression(data_x, data_y, nil, false)
	//fmt.Printf("El ajuste lineal es %.4v x + %.4v", a, b)
	line := plotter.NewFunction(func(x float64) float64 { return a*x + b })
	line.Color = color.RGBA{B: 255, A: 255}

	return scatter, line
}

//Esta función realiza la gráfica y la guarda
func scatter_lr_plot(scat_plotter *plotter.Scatter, line_plotter *plotter.Function, title, x_label, y_label string) {
	//Se crea el plot y se definen los labels
	p := plot.New()
	p.Add(scat_plotter, line_plotter)
	p.Title.Text = title
	p.X.Label.Text = x_label
	p.Y.Label.Text = y_label

	//Se guarda la imagen
	image_title := fmt.Sprintf("%s_%s.png", x_label, y_label)
	if err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, image_title); err != nil {
		panic(err)
	}
}

func main() {

	scatter_tv, line_tv := scatter_lr_plotter("Advertising.csv", 1, 4)
	scatter_lr_plot(scatter_tv, line_tv, "Advertising data set (1)", "TV", "sales")

	scatter_radio, line_radio := scatter_lr_plotter("Advertising.csv", 2, 4)
	scatter_lr_plot(scatter_radio, line_radio, "Advertising data set (2)", "radio", "sales")

	scatter_news, line_news := scatter_lr_plotter("Advertising.csv", 3, 4)
	scatter_lr_plot(scatter_news, line_news, "Advertising data set (3)", "news", "sales")
}
