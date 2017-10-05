package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
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

// A map from filter colors to a naive RGB color scheme for the V1 algorithm.
var FilterMap_V1 = map[string]color.Color{
	"BL1": color.RGBA{0, 0, 255, 255},
	"GRN": color.RGBA{0, 255, 0, 255},
	"RED": color.RGBA{255, 0, 0, 255},
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

func convertToAlpha(grayImage image.Image) image.Image {
	bounds := grayImage.Bounds()
	mask := image.NewRGBA(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			grayPixel := grayImage.At(x, y).(color.Gray)
			rgbaPixel := color.RGBA{0, 0, 0, grayPixel.Y}
			mask.Set(x, y, rgbaPixel)
		}
	}

	return mask
}

func layerColor(dst draw.Image, grayImage image.Image, layerColor color.Color) {
	src := &image.Uniform{layerColor}
	mask := convertToAlpha(grayImage)
	draw.DrawMask(dst, grayImage.Bounds(), src, image.ZP, mask, image.ZP, draw.Over)
}

func combineImages1(config ConfigFile, root string) error {
	var imageBounds image.Rectangle
	imageMap := make(map[string]image.Image)

	for _, imageConfig := range config.Files {
		fmt.Println("Reading: ", imageConfig.Filename)
		fullPath := path.Join(root, imageConfig.Filename)
		img, err := loadImage(fullPath)
		if err != nil {
			return err
		}

		newBounds := img.Bounds()
		if imageBounds != image.ZR && imageBounds != newBounds {
			return fmt.Errorf("image %s: has different bounds (%s) than other images (%s)",
				fullPath, newBounds, imageBounds)
		}
		imageBounds = newBounds
		imageMap[imageConfig.Filter] = img
	}

	if imageMap["BL1"] == nil || imageMap["GRN"] == nil || imageMap["RED"] == nil {
		var filters []string
		for k := range imageMap {
			filters = append(filters, k)
		}
		return fmt.Errorf("images in %s: missing one or more RGB filters: %s", root, filters)
	}

	fmt.Println("Loaded images of size:", imageBounds)

	composedImage := image.NewRGBA(imageBounds)

	layerColor(composedImage, imageMap["BL1"], FilterMap_V1["BL1"])
	layerColor(composedImage, imageMap["GRN"], FilterMap_V1["GRN"])
	layerColor(composedImage, imageMap["RED"], FilterMap_V1["RED"])

	outPath := path.Join(root, "output_v1.jpg")
	fmt.Println("Writing image to:", outPath)
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}

	return jpeg.Encode(f, composedImage, nil)
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
