// A package containing the functions for v3 algorithm using blending and an attempt to align the images.
package algv3aligning

import (
	"fmt"
	"image"
	"image/color"
	"github.com/lewchuk/gostitcher/common"
	"math"
	"path"
)

func subtractImages(imageMap common.ImageMap, base string, layer string, xOffset int, yOffset int) (image.Image, int) {
	baseImage := imageMap[base]
	layerImage := imageMap[layer]

	offsetPoint := image.Pt(xOffset, yOffset)
	bounds := baseImage.Bounds()
	offsetBounds := bounds.Add(offsetPoint)
	overlapBounds := bounds.Intersect(offsetBounds)
	composedImage := image.NewGray(overlapBounds)
	totalDelta := 0
	maxD := 0
	for x := overlapBounds.Min.X; x < overlapBounds.Max.X; x++ {
		for y := overlapBounds.Min.Y; y < overlapBounds.Max.Y; y++ {
			delta := int(baseImage.GrayAt(x, y).Y) - int(layerImage.GrayAt(x - xOffset, y - yOffset).Y)
			if delta < 0 {
				delta *= -1
			}
			deltaPixel := color.Gray{uint8(delta)}
			composedImage.SetGray(x, y, deltaPixel)
			if (delta > maxD) {
				maxD = delta
			}
			// Ignore middle deltas which probably represent the background and overshaddow the delta of the image.
			// Want to minimize the extreme differences of image features.
			// Values chosen from a minimum max delta of 134 and then split into quarters preserving values in top and bottom 25%.
			if (delta > 32 && delta < 96) {
				continue
			}
			totalDelta += delta
		}
	}
	// fmt.Println(base, layer, xOffset, yOffset, totalDelta, maxD)

	return composedImage, totalDelta
}

func AlignImages(imageMap common.ImageMap, root string) error {
	bgX, bgY, brX, brY := 0, 0, 0, 0
	minBg, minBr := math.MaxInt32, math.MaxInt32
	for x := -5; x < 5; x++ {
		for y := -5; y < 5; y++ {
			_, bgVal := subtractImages(imageMap, common.BLUE, common.GREEN, x, y)
			_, brVal := subtractImages(imageMap, common.BLUE, common.RED, x, y)
			if (bgVal < minBg) {
				minBg = bgVal
				bgX = x
				bgY = y
			}

			if (brVal < minBr) {
				minBr = brVal
				brX = x
				brY = y
			}
 		}
 	}

 	bgzImg, bgzVal := subtractImages(imageMap, common.BLUE, common.GREEN, 0, 0)
	brzImg, brzVal := subtractImages(imageMap, common.BLUE, common.RED, 0, 0)
	bgoImg, bgoVal := subtractImages(imageMap, common.BLUE, common.GREEN, bgX, bgY)
	broImg, broVal := subtractImages(imageMap, common.BLUE, common.RED, brX, brY)

	err := common.WriteImage(path.Join(root, fmt.Sprintf("output_v3_bg_align_%d%d.jpg", 0, 0)), bgzImg)
 	if err != nil { return err }

 	err = common.WriteImage(path.Join(root, fmt.Sprintf("output_v3_br_align_%d%d.jpg", 0, 0)), brzImg)
	if err != nil { return err }

	err = common.WriteImage(path.Join(root, fmt.Sprintf("output_v3_bg_align_%d%d.jpg", bgX, bgY)), bgoImg)
 	if err != nil { return err }

 	err = common.WriteImage(path.Join(root, fmt.Sprintf("output_v3_br_align_%d%d.jpg", brX, brY)), broImg)
	if err != nil { return err }

 	fmt.Printf("0-0 BG value: %d, Min %d-%d BG value: %d\n", bgzVal, bgX, bgY, bgoVal)
 	fmt.Printf("0-0 BR value: %d, Min %d-%d BR value: %d\n", brzVal, brX, brY, broVal)

 	return nil
}
