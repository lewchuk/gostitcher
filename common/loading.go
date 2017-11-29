package common

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path"
)

func LoadImage(data io.Reader) (*image.Gray, error) {
	img, err := jpeg.Decode(data)

	if err != nil {
		return nil, fmt.Errorf("error decoding image: %s", err)
	}

	grayImg, ok := img.(*image.Gray)
	if !ok {
		return nil, fmt.Errorf("image not grayscale (%s), can't process", img.ColorModel())
	}

	return grayImg, nil
}

// LoadImageFromPath, loads an image from a path and validates that it is a grayscale image.
// It returns the image as a properly typed Gray image and any error encountered.
func LoadImageFromPath(imagePath string) (*image.Gray, error) {
	f, err := os.Open(imagePath)
	defer f.Close()

	if err != nil {
		return nil, fmt.Errorf("error opening image %s: %s", imagePath, err)
	}

	image, err := LoadImage(f)
	if err != nil {
		return nil, fmt.Errorf("error loading image %s: %s", imagePath, err)
	}

	return image, nil
}

func ValidateImageMap(imageMap map[string]string) error {
	for _, filter := range Filters {
		if _, ok := imageMap[filter]; !ok {
			var filters []string
			for k := range imageMap {
				filters = append(filters, k)
			}
			return fmt.Errorf("images missing one or more RGB filters: %s in %s", filter, filters)
		}
	}
	return nil
}

// LoadImages loads all images based on a config file and a filesystem root.
// The set of images will be validated as grayscale images and to ensure they
// represent a complete RGB set of images.
// It returns a map of filters to images and any errors encountered.
func LoadImages(config ConfigFile, root string) (ImageMap, error) {
	var imageBounds image.Rectangle
	imageMap := make(ImageMap)
	filenameMap := make(ImageFilenameMap)

	for _, imageConfig := range config.Files {
		fmt.Println("Reading: ", imageConfig.Filename)
		fullPath := path.Join(root, imageConfig.Filename)
		img, err := LoadImageFromPath(fullPath)
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
		filenameMap[imageConfig.Filter] = imageConfig.Filename
		fmt.Println("Filter:", imageConfig.Filter)
	}

	if err := ValidateImageMap(filenameMap); err != nil {
		return nil, fmt.Errorf("%s: %s", root, err)
	}

	return imageMap, nil
}
