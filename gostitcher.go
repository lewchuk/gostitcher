package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
)

// https://space.stackexchange.com/questions/12510/cassinis-camera-continuum-band-filters
// A map of filter names to effective wavelengths.
var FilterMap = map[string]int{
	"BL1": 463,
	"GRN": 568,
	"RED": 647,
}

// The type for marshalling a config.json file from a folder with images.
type ImageConfig struct {
	Filename string `json:"filename"`
	Filter   string `json:"filter"`
}

type ConfigFile struct {
	Files []ImageConfig `json:"files"`
}

func main() {
	inputPath := os.Args[1]
	fmt.Printf("Processing image at: %s\n", inputPath)

	configS, err := ioutil.ReadFile(path.Join(inputPath, "config.json"))

	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(-1)
	}

	config := ConfigFile{}
	err2 := json.Unmarshal(configS, &config)

	if err2 != nil {
		fmt.Println("Error parsing config:", err)
		os.Exit(-1)
	}

	for _, imageConfig := range config.Files {
		fmt.Println("Processing: ", imageConfig.Filename)
		fullPath := path.Join(inputPath, imageConfig.Filename)

		f, err3 := os.Open(fullPath)
		if err3 != nil {
			fmt.Println("Error reading iamge:", err3)
			os.Exit(1)
		}

		image, err4 := jpeg.Decode(f)

		if err4 != nil {
			fmt.Println("Error decoding image:", err4)
			os.Exit(1)
		}

		if image.ColorModel() != color.GrayModel {
			fmt.Println("Image not grayscale, can't convert")
			os.Exit(1)
		}

		fmt.Println(image.Bounds())
	}
}
