package main

import (
    "os"
    "fmt"
    "log"
    "sync"
    "time"
    "strings"
    "strconv"
    "encoding/csv"
    
    "github.com/kfultz07/go-dataframe"
)

type RatingObj struct {
    UserId    int64
    movieId   int64
    rating    float64
    Timestamp int64
}

func main() {
    var dir, file string
    var chunkNumber int
    var chunkData [][]string
    
    var startTotalTime = time.Now()
    dir = "movielens"
    directory := MakeDir(dir)
    fmt.Println("Los archivos de salida se guardarán en el directorio:", dir)

    var startTime = time.Now()
    file = "ratings"
    data := readCSV(file)
    fmt.Printf("El archivo a particionar es '%s', contiene '%d' lineas en total", file, len(data))
    fmt.Println("\ny tomó", time.Since(startTime).Seconds(), "segundos en leer el numero de lineas.")
    
    chunkNumber = 10
    numLineas := len(data)/chunkNumber
    fmt.Printf("\nEl archivo '%s' se dividirá en '%d' subarchivos",file,chunkNumber)
    fmt.Printf("\ny cada subarchivo tendrá '%d' lineas.\n",numLineas)
    
    var wg sync.WaitGroup    // Numero de gorutinas trabajando
    startTime = time.Now()
    chunkName := make(chan string)
    for i:=0; i<chunkNumber; i++ {
        wg.Add(2)
        go createCSV(i, file, directory, chunkName, &wg)
        if i == chunkNumber-1 {
            chunkData = data[(numLineas*i):]
        } else {
            chunkData = data[(numLineas*i):(numLineas*(i+1))]
        }
        go writeCSV(chunkData, chunkName, &wg)
        wg.Wait()
    }
    
    fmt.Println("\n¡Listo! El particionado de archivos tardó", time.Since(startTime).Seconds(), "segundos.")
    
    Mt_FindRatingsMaster(file, directory, chunkNumber)

    fmt.Println("\nTiempo total de análisis", time.Since(startTotalTime), "minutos.")
}

// Crear directorio donde almacenar todos los archivos particionados
func MakeDir(dir string) string {
    directory := "./" + dir + "/"
    err := os.Mkdir(directory, 0777)
    if err != nil {
        log.Fatalf("Error al crear el directorio: %s", err)
    }
    return directory
}

// Lectura y almacenaje de los datos contenidos en el archivo principal
func readCSV(file string) (data [][]string) {
    filename, err := os.Open(file + ".csv") 
    if err != nil {
        log.Fatalf("Error al acceder al archivo: %s", err)
    }
    defer filename.Close()
    
    csvReader := csv.NewReader(filename)

    csvReader.Comment = 'u'  // Removiendo el encabezado, si existe.
    data, err = csvReader.ReadAll()
    if err != nil {
        log.Fatalf("Error al leer los datos del archivo: %s", err)
    }
    return data
}

// Creación archivos particionados
func createCSV(i int, file string, directory string, chunkName chan string, wg *sync.WaitGroup) {
    defer wg.Done()
    Name := directory + file + "_" + strconv.Itoa(i) + ".csv"
    chunkName <- Name
}

// Escritura de datos en cada archivo particionado
func writeCSV(chunkData [][]string, chunkName chan string, wg *sync.WaitGroup) {
    defer wg.Done()
    Name := <- chunkName
    csvFile, err := os.Create(Name)
    if err != nil {
        log.Fatalf("Error al crear los archivos particionados: %s", err)
    }
    defer csvFile.Close()
    
    dataWriter := csv.NewWriter(csvFile)
    defer dataWriter.Flush()
    
    // Escribir encabezados
    dataWriter.Write([]string{"UserId", "movieId", "rating", "Timestamp"})
    
    // Procesar y escribir cada registro en el archivo
    for _, record := range chunkData {
        userId, _ := strconv.ParseInt(record[0], 10, 64)
        movieId, _ := strconv.ParseInt(record[1], 10, 64)
        rating, _ := strconv.ParseFloat(record[2], 64)
        timestamp, _ := strconv.ParseInt(record[3], 10, 64)
        
        data := RatingObj{
            UserId:    userId,
            movieId:   movieId,
            rating:    rating,
            Timestamp: timestamp,
        }
        
        // Convertir el objeto Rating en un registro CSV
        recordToWrite := []string {
            strconv.FormatInt(data.UserId, 10),
            strconv.FormatInt(data.movieId, 10),
            strconv.FormatFloat(data.rating, 'f', 1, 64),
            strconv.FormatInt(data.Timestamp, 10),
        }
        dataWriter.Write(recordToWrite)
    }
    // dataWriter.WriteAll(chunkData)
    
    if err != nil {
        log.Fatalf("Error al escribir los datos en los archivos particionados: %s", err)
    }
}

func Mt_FindRatingsMaster(file string, directory string, chunkNumber int) {
    fmt.Println("\nAhora inicia 'MtFindRatingsMaster'\n")
    start := time.Now()
    nf := chunkNumber
    
    kg := []string{"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary",
                   "Drama", "Fantasy", "Film-Noir", "Horror", "IMAX", "Musical", "Mystery", "Romance",
                   "Sci-Fi", "Thriller", "War", "Western", "(no genres listed)"}
    
    ng := len(kg)
        
    ra := make([][]float64, ng)
    ca := make([][]int, ng)
    
    for i := 0; i < ng; i++ {
        ra[i] = make([]float64, nf)
        ca[i] = make([]int, nf)
    }
    var ci = make(chan int)
    movies := ReadMoviesCsvFile("./","movies.csv")
    
    for i := 0; i < nf; i++ {
        go Mt_FindRatingsWorker(i, file, directory, ci, kg, &ca, &ra, movies)
    }
    
    iMsg := 0
    go func() {
        for {
            i := <-ci
            iMsg += i
        }
    }()
    for {
        if iMsg == 10 {
            break
        }
    }
    
    locCount := make([]int, ng)
    locVals := make([]float64, ng)
    for i := 0; i < ng; i++ {
        for j := 0; j < nf; j++ {
            locCount[i] += ca[i][j]
            locVals[i] += ra[i][j]
        }
    }
    
    fmt.Println("\n======================> RESULTADOS <======================")
    // Calcular y mostrar el promedio de ratings por género
    for i := 0; i < ng; i++ {
        avgRating := 0.0
        if locCount[i] > 0 {
            avgRating = locVals[i] / float64(locCount[i])
        }
        fmt.Println(fmt.Sprintf("%2d", i), "  ", fmt.Sprintf("%20s", kg[i]), "  ", fmt.Sprintf("%8d", locCount[i]), "  ", fmt.Sprintf("Avg Rating: %.2f", avgRating))
    }
    
    duration := time.Since(start)
    fmt.Println("\n\nDuración de ejecución de 'Mt_FindRatingsMaster'", duration)
    println("'Mt_FindRatingsMaster' terminó")
}

func Mt_FindRatingsWorker(w int, file string, directory string, ci chan int, kg []string, ca *[][]int, va *[][]float64, movies dataframe.DataFrame) {
    aFileName := file + "_" + fmt.Sprintf("%01d", w) + ".csv"
    println("El worker " , fmt.Sprintf("%01d", w), " está procesando el archivo ", aFileName)
    
    ratings := ReadRatingsCsvFile(directory, aFileName)
    ng := len(kg)
    start := time.Now()
        
    ratings.Merge(&movies, "movieId", "genres")

    grcs := [2]string{"genres", "rating"}
    grDF := ratings.KeepColumns(grcs[:])
    for ig := 0; ig < ng; ig++ {
        for _, row := range grDF.FrameRecords {
            if strings.Contains(row.Data[0], kg[ig]) {
                (*ca)[ig][w] += 1
                v, _ := strconv.ParseFloat((row.Data[1]), 32)
                (*va)[ig][w] += v
            }
        }
    }
    duration := time.Since(start)
    fmt.Println("El worker ", w," terminó con una duración de ", duration)
    
    ci <- 1
}

// Función para leer el archivo de películas
func ReadMoviesCsvFile(filePath string, fileOb string) dataframe.DataFrame {
    // Abrir el archivo
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Error al abrir el archivo %s: %v", filePath, err)
    }
    defer file.Close()
    
    // Leer el CSV y devolver el DataFrame
    return dataframe.CreateDataFrame(filePath, fileOb)
}

// Función para leer el archivo de ratings
func ReadRatingsCsvFile(directory string, fileOb string) dataframe.DataFrame {
    // Abrir el archivo
    file, err := os.Open(directory)
    if err != nil {
        log.Fatalf("Error al abrir el archivo %s: %v", directory, err)
    }
    defer file.Close()
    
    // Leer el CSV y devolver el DataFrame
    return dataframe.CreateDataFrame(directory, fileOb)
}

func spinner(delay time.Duration) {
    for {
        for _, r := range `-\|/` {
            fmt.Printf("\r %c", r)
            time.Sleep(delay)
        }
    }
}