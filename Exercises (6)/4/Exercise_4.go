package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"os"
	"path/filepath"

	"gonum.org/v1/gonum/mat"
)

func load_grayscale_image(filepath string) (*mat.Dense, error) { // Loads a grayscale image and converts it to a gonum array
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file) // Decode the image into JPEG format
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds() // Image to grayscale
	width, height := bounds.Max.X, bounds.Max.Y
	grayscaleImg := make([]float64, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the pixel value at position (x, y)
			r, g, b, _ := img.At(x, y).RGBA()
			// Calculate the gray value using formula
			grayscaleImg[y*width+x] = float64((r+g+b)/3) / 65535.0
		}
	}
	grayscale_matrix := mat.NewDense(height, width, grayscaleImg)
	return grayscale_matrix, nil
}

func save_grayscale_image(filepath string, imgMatrix *mat.Dense) error {
	height, width := imgMatrix.Dims()
	img := image.NewGray(image.Rect(0, 0, width, height)) // Creates a new blank image
	for y := 0; y < height; y++ {                         // Converts grayscale matrix to pixels
		for x := 0; x < width; x++ {
			grayValue := uint8(imgMatrix.At(y, x) * 255) // Convert the float value of the array to uint8 for the pixels
			img.Set(x, y, color.Gray{Y: grayValue})
		}
	}

	file, err := os.Create(filepath) // Creats output file
	if err != nil {
		return err
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil) // Save the image like JPEG
	return err
}

func main() {
	inputDir := "/Users/michelletorres/Desktop/Homeworks AI/"
	outputDir := "/Users/michelletorres/Desktop/Homeworks AI/Reconstructed/"

	imageNames := []string{"Imagen 1.jpeg", "Imagen 2.jpeg", "Imagen 3.jpeg, Imagen 4.jpeg,Imagen 5.jpeg, Imagen 6.jpeg, Imagen 7.jpeg"}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) { // Create the output directory if it does not exist
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			fmt.Println("Error al crear el directorio:", err)
			return
		}
		fmt.Println("Directorio de salida creado:", outputDir)
	} else {
		fmt.Println("El directorio ya existe:", outputDir)
	}

	for _, imageName := range imageNames { // Process each image
		imagePath := filepath.Join(inputDir, imageName)
		fmt.Printf("Procesando: %s\n", imagePath)

		imgMatrix, err := load_grayscale_image(imagePath) // Load original image
		if err != nil {
			fmt.Printf("Error al cargar la imagen %s: %v\n", imageName, err)
			continue
		}

		rows, cols := imgMatrix.Dims() //Dimensions
		fmt.Printf("Dimensiones de la imagen: %dx%d\n", rows, cols)
		outputImagePath := filepath.Join(outputDir, "reconstructed_"+imageName)
		err = save_grayscale_image(outputImagePath, imgMatrix)
		if err != nil {
			fmt.Printf("Error al guardar imagen %s: %v\n", imageName, err)
		} else {
			fmt.Printf("Imagen reconstruida guardada en: %s\n", outputImagePath)
		}
	}
}
