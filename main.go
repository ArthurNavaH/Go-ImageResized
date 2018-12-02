package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type config struct {
	ImagesInput  string `json:"imagesInput"`
	ImagesOutput string `json:"imagesOutput"`
	WidthImage   int    `json:"widthImage"`
}

var wg sync.WaitGroup

func main() {
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	PrintStart()
	defer PrintEnd()

	var config config

	cj, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(cj, &config)
	if err != nil {
		fmt.Println(err)
	}

	files, _ := ioutil.ReadDir(config.ImagesInput)

	for i := 0; i < len(files); i++ {
		wg.Add(1)
		go SaveImage(config, files[i], i)
	}

	wg.Wait()
}

func PrintStart() {
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("- Iniciando algoritmo E.A.M.B (Escalo Automatico de Mapa de Bits ")
	fmt.Println("-----------------------------------------------------------------\n")
}

func PrintEnd() {
	fmt.Println("\n-----------------------")
	fmt.Println("- Ejecusion terminada -")
	fmt.Println("-----------------------\n")
}

func SaveImage(config config, file os.FileInfo, index int) {
	defer wg.Done()

	var width, height, newWidth, newHeight, difWidth int
	var porDifWidth float32

	fmt.Println(file.Name())

	img, err := imgio.Open(config.ImagesInput + file.Name())

	if err != nil {
		panic(err)
	}
	width = img.Bounds().Max.X
	height = img.Bounds().Max.Y

	newWidth = config.WidthImage

	difWidth = width - newWidth
	porDifWidth = float32(difWidth) / float32(width)

	newHeight = height - int(float32(height)*porDifWidth)

	resized := transform.Resize(img, newWidth, newHeight, transform.Linear)

	var nameFile string
	for x, v := range file.Name() {
		if v == '.' {
			nameFile = file.Name()[:x]
		}
	}

	if err := imgio.Save(config.ImagesOutput+nameFile+".png", resized, imgio.PNG); err != nil {
		fmt.Println("\nA ocurrido un error con la Imagen #", index+1, " Nombre:", file.Name())
		panic(err)
	}

	fmt.Println("Imagen #", index+1, " Lista! +", "Nombre :", file.Name(), ";")

}
