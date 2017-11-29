// A package containing the functions for v2 algorithm using blending.
package algv2blending

import (
	"github.com/lewchuk/gostitcher/common"
	"image"
	"image/color"
	"path"
)

// blendImage combines separte RGB grayscale images into a single RGB image.
// It returns the generated image.
func BlendImage(imageMap common.ImageMap) image.Image {
	blueImage := imageMap[common.BLUE]
	greenImage := imageMap[common.GREEN]
	redImage := imageMap[common.RED]

	bounds := blueImage.Bounds()
	composedImage := image.NewRGBA(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			rgbaPixel := color.RGBA{
				redImage.GrayAt(x, y).Y,
				greenImage.GrayAt(x, y).Y,
				blueImage.GrayAt(x, y).Y,
				255}
			composedImage.Set(x, y, rgbaPixel)
		}
	}

	return composedImage
}

// CombineImages runs the v2 blending algorithm to combine a set of grayscale images into
// a "true" color image.
func CombineImages(imageMap common.ImageMap, root string) error {
	composedImage := BlendImage(imageMap)

	err := common.WriteImage(path.Join(root, "output_v2_alpha.jpg"), composedImage)
	if err != nil {
		return err
	}

	return nil
}
