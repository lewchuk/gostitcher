package common

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
)

// loadImage, loads an image from a path and validates that it is a grayscale image.
// It returns the image as a properly typed Gray image and any error encountered.
func loadImage(imagePath string) (*image.Gray, error) {
	f, err := os.Open(imagePath)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(f)

	if err != nil {
		return nil, fmt.Errorf("decoding image %s: %s", imagePath, err)
	}

	grayImg, ok := img.(*image.Gray)
	if !ok {
		return nil, fmt.Errorf("image %s: not grayscale (%s), can't process", imagePath, img.ColorModel())
	}

	return grayImg, nil
}

// LoadImages loads all images based on a config file and a filesystem root.
// The set of images will be validated as grayscale images and to ensure they
// represent a complete RGB set of images.
// It returns a map of filters to images and any errors encountered.
func LoadImages(config ConfigFile, root string) (map[string]image.Gray, error) {
	var imageBounds image.Rectangle
	imageMap := make(map[string]image.Gray)

	for _, imageConfig := range config.Files {
		fmt.Println("Reading: ", imageConfig.Filename)
		fullPath := path.Join(root, imageConfig.Filename)
		img, err := loadImage(fullPath)
		if err != nil {
			return nil, err
		}

		newBounds := img.Bounds()
		if imageBounds != image.ZR && imageBounds != newBounds {
			return nil, fmt.Errorf("image %s: has different bounds (%s) than other images (%s)",
				fullPath, newBounds, imageBounds)
		}
		imageBounds = newBounds
		imageMap[imageConfig.Filter] = *img
		fmt.Println("Filter:", imageConfig.Filter)
	}

	for _, filter := range Filters {
		if _, ok := imageMap[filter]; !ok {
			var filters []string
			for k := range imageMap {
				filters = append(filters, k)
			}
			return nil, fmt.Errorf("images in %s: missing one or more RGB filters: %s in %s", root, filter, filters)
		}
	}
	return imageMap, nil
}
