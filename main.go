package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type config struct {
	ImagesInput  string `json:"imagesInput"`
	ImagesOutput string `json:"imagesOutput"`
	WidthImage   int    `json:"widthImage"`
}

func main() {
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("- Iniciando algoritmo E.A.M.B (Escalo Automatico de Mapa de Bits ")
	fmt.Println("-----------------------------------------------------------------\n")
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
		var width, height, newWidth, newHeight, difWidth int
		var porDifWidth float32

		img, err := imgio.Open(config.ImagesInput + files[i].Name())

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
		for x, v := range files[i].Name() {
			if v == '.' {
				nameFile = files[i].Name()[:x]
			}
		}

		if err := imgio.Save(config.ImagesOutput+nameFile+".png", resized, imgio.PNG); err != nil {
			fmt.Println("\nA ocurrido un error con la Imagen #", i+1, " Nombre:", files[i].Name())
			panic(err)
		}

		fmt.Println("Imagen #", i+1, " Lista! +", "Nombre :", files[i].Name(), ";")
	}

	fmt.Println("\n-----------------------")
	fmt.Println("- Ejecusion terminada -")
	fmt.Println("-----------------------\n")
}
