// A package containing the functions for v1 algorithm using masking.
package algv1masking

import (
	"github.com/lewchuk/gostitcher/common"
	"image"
	"image/color"
	"image/draw"
	"path"
)

// A map from filter colors to a naive RGB color scheme.
var filterMap = map[string]color.Color{
	common.BLUE:  color.RGBA{0, 0, 255, 255},
	common.GREEN: color.RGBA{0, 255, 0, 255},
	common.RED:   color.RGBA{255, 0, 0, 255},
}

// convertToAlpha takes a Gray image and converts it to an RGBA image.
// The value of the gray pixel is used to create a black pixel with an alpha
// value equal to the gray pixel.
// It returns an RGBA image suitable for use as a mask.
func convertToAlpha(grayImage image.Gray) image.Image {
	bounds := grayImage.Bounds()
	mask := image.NewRGBA(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			rgbaPixel := color.RGBA{0, 0, 0, grayImage.GrayAt(x, y).Y}
			mask.Set(x, y, rgbaPixel)
		}
	}

	return mask
}

// layerColor uses a grayscale image as a mask to draw a layer of color onto another image.
func layerColor(dst draw.Image, grayImage image.Gray, layerColor color.Color) {
	src := &image.Uniform{layerColor}
	mask := convertToAlpha(grayImage)
	draw.DrawMask(dst, grayImage.Bounds(), src, image.ZP, mask, image.ZP, draw.Over)
}

// CombineImages runs the v1 masking algorithm to combine a set of grayscale images into
// a "true" color image.
func CombineImages(imageMap common.ImageMap, root string) error {
	// LoadImages validates the presence of RGB images and that they all share the same bounds.
	blueImage := imageMap[common.BLUE].Image
	imageBounds := blueImage.Bounds()

	composedImage := image.NewRGBA(imageBounds)

	layerColor(composedImage, imageMap[common.BLUE].Image, filterMap[common.BLUE])
	layerColor(composedImage, imageMap[common.GREEN].Image, filterMap[common.GREEN])
	layerColor(composedImage, imageMap[common.RED].Image, filterMap[common.RED])

	err := common.WriteImage(path.Join(root, "output_v1_alpha.jpg"), composedImage)
	if err != nil {
		return err
	}

	composedImage2 := image.NewRGBA(imageBounds)

	layerColor(composedImage2, imageMap[common.RED].Image, filterMap[common.RED])
	layerColor(composedImage2, imageMap[common.GREEN].Image, filterMap[common.GREEN])
	layerColor(composedImage2, imageMap[common.BLUE].Image, filterMap[common.BLUE])

	return common.WriteImage(path.Join(root, "output_v1_beta.jpg"), composedImage2)
}
