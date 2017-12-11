// A package containing the functions for v3 algorithm using blending and an attempt to align the images.
package algv3aligning

import (
	"fmt"
	"github.com/lewchuk/gostitcher/common"
	"image"
	"image/color"
	"math"
	"path"
)

type ImageOffsets struct {
	bgX int
	bgY int
	brX int
	brY int
}

func getPixel(image common.LoadedConfig, x, y int) uint8 {
	return image.Image.GrayAt(x-image.Config.OffsetX, y-image.Config.OffsetY).Y
}

func subtractImages(baseImage, layerImage common.LoadedConfig) (image.Image, int) {
	// TODO: Figure out if this actually works for the case where baseImages offsets are non 0.
	xOffset := layerImage.Config.OffsetX - baseImage.Config.OffsetX
	yOffset := layerImage.Config.OffsetY - baseImage.Config.OffsetY

	offsetPoint := image.Pt(xOffset, yOffset)
	bounds := baseImage.Image.Bounds()
	offsetBounds := bounds.Add(offsetPoint)
	overlapBounds := bounds.Intersect(offsetBounds)
	composedImage := image.NewGray(overlapBounds)
	totalDelta := 0
	maxD := 0
	for x := overlapBounds.Min.X; x < overlapBounds.Max.X; x++ {
		for y := overlapBounds.Min.Y; y < overlapBounds.Max.Y; y++ {
			delta := int(getPixel(baseImage, x, y)) - int(getPixel(layerImage, x, y))
			if delta < 0 {
				delta *= -1
			}
			deltaPixel := color.Gray{uint8(delta)}
			composedImage.SetGray(x, y, deltaPixel)
			if delta > maxD {
				maxD = delta
			}
			// Ignore middle deltas which probably represent the background and overshaddow the delta of the image.
			// Want to minimize the extreme differences of image features.
			// Values chosen from a minimum max delta of 134 and then split into quarters preserving values in top and bottom 25%.
			if delta > 32 && delta < 96 {
				continue
			}
			totalDelta += delta
		}
	}

	return composedImage, totalDelta
}

func AlignImages(imageMap *common.ImageMap, maxOffset int) {
	bgX, bgY, brX, brY := 0, 0, 0, 0
	minBg, minBr := math.MaxInt32, math.MaxInt32
	blueImage := (*imageMap)[common.BLUE]
	for x := -1 * maxOffset; x < maxOffset; x++ {
		for y := -1 * maxOffset; y < maxOffset; y++ {
			greenImage := (*imageMap)[common.GREEN]
			redImage := (*imageMap)[common.RED]
			greenImage.Config.OffsetX = x
			greenImage.Config.OffsetY = y
			redImage.Config.OffsetX = x
			redImage.Config.OffsetY = y
			_, bgVal := subtractImages(blueImage, greenImage)
			_, brVal := subtractImages(blueImage, redImage)
			if bgVal < minBg {
				minBg = bgVal
				bgX = x
				bgY = y
			}

			if brVal < minBr {
				minBr = brVal
				brX = x
				brY = y
			}
		}
	}

	// Update the green and red configs to have proper offsets.
	greenImage := (*imageMap)[common.GREEN]
	greenImage.Config.OffsetX = bgX
	greenImage.Config.OffsetY = bgY
	(*imageMap)[common.GREEN] = greenImage
	redImage := (*imageMap)[common.RED]
	redImage.Config.OffsetX = brX
	redImage.Config.OffsetY = brY
	(*imageMap)[common.RED] = redImage
}

func CombineImages(imageMap common.ImageMap) image.Image {
	blueImage := imageMap[common.BLUE]
	greenImage := imageMap[common.GREEN]
	redImage := imageMap[common.RED]

	bounds := blueImage.Image.Bounds()
	composedImage := image.NewRGBA(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			rgbaPixel := color.RGBA{
				getPixel(redImage, x, y),
				getPixel(greenImage, x, y),
				getPixel(blueImage, x, y),
				255}
			composedImage.Set(x, y, rgbaPixel)
		}
	}

	return composedImage
}

func OutputImageDiffs(imageMap common.ImageMap, root string) error {
	bgImg, _ := subtractImages(imageMap[common.BLUE], imageMap[common.GREEN])
	brImg, _ := subtractImages(imageMap[common.BLUE], imageMap[common.RED])

	err := common.WriteImage(path.Join(root, fmt.Sprintf("output_v3_bg_align_%d%d.jpg", imageMap[common.GREEN].Config.OffsetX, imageMap[common.GREEN].Config.OffsetY)), bgImg)
	if err != nil {
		return err
	}

	err = common.WriteImage(path.Join(root, fmt.Sprintf("output_v3_br_align_%d%d.jpg", imageMap[common.RED].Config.OffsetX, imageMap[common.RED].Config.OffsetY)), brImg)
	if err != nil {
		return err
	}

	return nil
}

func CombineAndAlignImages(config common.ConfigFile, imageMap common.ImageMap, maxOffset int, root string) error {
	if maxOffset > config.MaxOffset {
		// Output unalingned/last best aligned diffs.
		OutputImageDiffs(imageMap, root)

		// Run the alignment algorithm to update the imageMap.
		AlignImages(&imageMap, maxOffset)

		// Update the config file with new data
		config.MaxOffset = maxOffset
		for i, sourceConfig := range config.Files {
			config.Files[i] = imageMap[sourceConfig.Filter].Config
		}

		if err := common.WriteConfig(root, config); err != nil {
			return err
		}
	}

	// Output new/current best diffs.
	OutputImageDiffs(imageMap, root)

	// Create combined colour image.
	composedImage := CombineImages(imageMap)

	err := common.WriteImage(path.Join(root, "output_v3.jpg"), composedImage)
	if err != nil {
		return err
	}

	return nil
}
