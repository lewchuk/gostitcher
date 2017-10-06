package common

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

// WriteImages writes out an image as a jpeg to a given path.
// Returns any errors from the write.
func WriteImage(path string, img image.Image) error {
	fmt.Println("Writing image to:", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return jpeg.Encode(f, img, nil)
}
