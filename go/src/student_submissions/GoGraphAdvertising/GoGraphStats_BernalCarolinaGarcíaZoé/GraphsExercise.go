// Using advertising.csv graph
// Duo Dinamita :) Bernal Rodriguez Carolina ; García Martínez Zoé Ariel

package main

import (
	"image/color"
	"strconv"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"

	//"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	//"gonum.org/v1/plot/vg/draw"
	//"github.com/gonum/plot"
	//"github.com/gonum/plot/plotter"
	//"github.com/gonum/plot/vg"
	//"github.com/gonum/plot/vg/draw"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("\nReady? ,Here we go! ")

	fileName := "Advertising.csv"
	directory := "../GoGraphStats"

	headRows := 5

	records := readCSV(fileName, directory)
	fmt.Println("\nHere is a look of your data: ")

	for i, record := range records {
		if i >= headRows {
			break
		}
		fmt.Println(record)
	}

	tvIndex := 1
	radioIndex := 2
	newspaperIndex := 3
	salesIndex := 4

	//fmt.Printf("\nCol %d data:", tvIndex)
	//fmt.Printf("\nCol %d data:", salesIndex)

	tv := getMediaData(records, tvIndex)
	radio := getMediaData(records, radioIndex)
	newspaper := getMediaData(records, newspaperIndex)
	sales := getMediaData(records, salesIndex)
	/**
		for i, record := range tv {
			if i >= headRows {
				break
			}
			fmt.Println(record)
		}
	**/
	//fmt.Println(tv[1])
	//fmt.Println(radio[1])
	//fmt.Println(newspaper[1])

	//--- TV Sales Graph ----

	p := plot.New()

	dataTV := scatterPoints(tv, sales)

	// Add a scatter plot to the plot
	scatter, err := plotter.NewScatter(dataTV)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(2)

	// Add the scatter plot to the plot and set the axes labels
	p.Add(scatter)
	p.Title.Text = "TV vs Sales"
	p.X.Label.Text = "TV"
	p.Y.Label.Text = "Sales"

	slope, intercept := linearGraph(tv, sales)

	// Linear adjusment to the data

	linearAdjust := plotter.NewFunction(func(x float64) float64 {
		return slope*x + intercept
	})

	fmt.Printf("Ajuste lineal TV: y = %.2fx + %.2f\n", slope, intercept)

	linearAdjust.LineStyle.Color = color.RGBA{B: 255, A: 255}
	p.Add(linearAdjust)

	// Save the plot to a PNG file
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "./plots/TV_Sales.png"); err != nil {
		panic(err)
	}

	// --- Newspaper Sales Graph -----

	g := plot.New()

	dataNews := scatterPoints(newspaper, sales)

	scatterNews, err := plotter.NewScatter(dataNews)
	if err != nil {
		panic(err)
	}
	scatterNews.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	scatterNews.GlyphStyle.Radius = vg.Points(2)

	// Add the scatter plot to the plot and set the axes labels
	g.Add(scatterNews)
	g.Title.Text = "Newspaper vs Sales"
	g.X.Label.Text = "Newspaper"
	g.Y.Label.Text = "Sales"

	slopeNews, interceptNews := linearGraph(newspaper, sales)

	linearAdjustNews := plotter.NewFunction(func(x float64) float64 {
		return slopeNews*x + interceptNews
	})

	fmt.Printf("Ajuste lineal Newspaper: y = %.2fx + %.2f\n", slopeNews, interceptNews)

	linearAdjustNews.LineStyle.Color = color.RGBA{B: 255, A: 255}
	g.Add(linearAdjustNews)

	// Save the plot to a PNG file
	if err := g.Save(4*vg.Inch, 4*vg.Inch, "./plots/Newspaper_Sales.png"); err != nil {
		panic(err)
	}

	//--- Radio Sales Graph ---

	q := plot.New()

	dataRadio := scatterPoints(radio, sales)

	scatterRadio, err := plotter.NewScatter(dataRadio)
	if err != nil {
		panic(err)
	}
	scatterRadio.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	scatterRadio.GlyphStyle.Radius = vg.Points(2)

	// Add the scatter plot to the plot and set the axes labels
	q.Add(scatterRadio)
	q.Title.Text = "Radio vs Sales"
	q.X.Label.Text = "Radio"
	q.Y.Label.Text = "Sales"

	slopeRadio, interceptRadio := linearGraph(radio, sales)

	linearAdjustRadio := plotter.NewFunction(func(x float64) float64 {
		return slopeRadio*x + interceptRadio
	})

	fmt.Printf("Ajuste lineal Radio: y = %.2fx + %.2f\n", slopeRadio, interceptRadio)

	linearAdjustRadio.LineStyle.Color = color.RGBA{B: 255, A: 255}
	q.Add(linearAdjustRadio)

	// Save the plot to a PNG file
	if err := q.Save(4*vg.Inch, 4*vg.Inch, "./plots/Radio_Sales.png"); err != nil {
		panic(err)
	}

	fmt.Println("\nDone :)")
}

// fileName is the name of the file; directory is the directory where the file is stored
func readCSV(fileName string, directory string) [][]string {

	// Join to get complete path
	fileName = filepath.Join(directory, fileName)

	file, err := os.Open(fileName)
	if err != nil {

		fmt.Println("cannot find the path to read the files, sorry :(")
		fmt.Println(fileName)
		return nil
	}
	defer file.Close()

	//Reading CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil
	}

	return records
}

func getMediaData(records [][]string, indexCol int) []string {
	var mediaTV []string

	for _, record := range records {
		if len(record) > indexCol {
			mediaTV = append(mediaTV, record[indexCol])
		} else {
			mediaTV = append(mediaTV, "")
		}
	}

	return mediaTV

}

func linearGraph(ColX []string, ColY []string) (float64, float64) {

	fileName := "Advertising.csv"
	directory := "../GoGraphStats"
	records := readCSV(fileName, directory)

	data := make(plotter.XYs, len(records))

	var xPoints []float64
	var yPoints []float64

	for i := range data {
		if i > 0 {
			x, err := strconv.ParseFloat(ColX[i], 64)
			if err != nil {
				panic(err)
			}
			xPoints = append(xPoints, x)

			y, err := strconv.ParseFloat(ColY[i], 64)
			if err != nil {
				panic(err)
			}
			yPoints = append(yPoints, y)
		}
	}

	//Linear Regresion
	intercept, slope := stat.LinearRegression(xPoints, yPoints, nil, false)
	//fmt.Printf("Ajuste lineal: y = %.2fx + %.2f\n", slope, intercept)

	return slope, intercept

}

func scatterPoints(ColX []string, ColY []string) plotter.XYs {

	fileName := "Advertising.csv"
	directory := "../GoGraphStats"
	records := readCSV(fileName, directory)

	data := make(plotter.XYs, len(records))

	for i := range data {
		if i > 0 {
			x, err := strconv.ParseFloat(ColX[i], 64)
			if err != nil {
				panic(err)
			}
			data[i].X = x

			y, err := strconv.ParseFloat(ColY[i], 64)
			if err != nil {
				panic(err)
			}
			data[i].Y = y
		}
	}

	return data
}
