package main

import (
	"encoding/json"
	"fmt"
	"image"
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

func loadImage(imagePath string) (image.Image, error) {
	f, err := os.Open(imagePath)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(f)

	if err != nil {
		return nil, fmt.Errorf("decoding image %s: %s", imagePath, err)
	}

	if img.ColorModel() != color.GrayModel {
		return nil, fmt.Errorf("image %s: not grayscale, can't process", imagePath)
	}

	return img, nil
}

func combineImages1(config ConfigFile, root string) error {
	// var imageBounds image.Rectangle
	imageMap := make(map[string]image.Image)

	for _, imageConfig := range config.Files {
		fmt.Println("Reading: ", imageConfig.Filename)
		fullPath := path.Join(root, imageConfig.Filename)
		img, err := loadImage(fullPath)
		if err != nil {
			return err
		}

		// imageBounds = img.Bounds()
		imageMap[imageConfig.Filter] = img
	}

	if imageMap["BL1"] == nil || imageMap["GRN"] == nil || imageMap["RED"] == nil {
		var filters []string
		for k := range imageMap {
		    filters = append(filters, k)
		}
		return fmt.Errorf("images in %s: missing one or more RGB filters: %s", root, filters)
	}

	return nil
}

func processImages(inputPath string) error {
	fmt.Printf("Processing: %s\n", inputPath)

	configPath := path.Join(inputPath, "config.json")
	configS, err := ioutil.ReadFile(configPath)

	if err != nil {
		return err
	}

	config := ConfigFile{}

	if err := json.Unmarshal(configS, &config); err != nil {
		return fmt.Errorf("parsing config %s: %s", configPath, err)
	}

	return combineImages1(config, inputPath)
}

func main() {
	err := processImages(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
