// A package containing the functions for v3 algorithm using blending and an attempt to align the images.
package algv3aligning

import (
	"fmt"
	"image"
	"image/color"
	"github.com/lewchuk/gostitcher/common"
	"path"
)

func subtractImages(imageMap common.ImageMap, base string, layer string, xOffset int, yOffset int) image.Image {
	baseImage := imageMap[base]
	layerImage := imageMap[layer]

	offsetPoint := image.Pt(xOffset, yOffset)
	bounds := baseImage.Bounds()
	offsetBounds := bounds.Add(offsetPoint)
	overlapBounds := bounds.Intersect(offsetBounds)
	composedImage := image.NewGray(overlapBounds)
	totalDelta := 0
	for x := overlapBounds.Min.X; x < overlapBounds.Max.X; x++ {
		for y := overlapBounds.Min.Y; y < overlapBounds.Max.Y; y++ {
			delta := int(baseImage.GrayAt(x, y).Y) - int(layerImage.GrayAt(x - xOffset, y - yOffset).Y)
			if delta < 0 {
				delta *= -1
			}
			totalDelta += delta
			deltaPixel := color.Gray{uint8(delta)}
			composedImage.SetGray(x, y, deltaPixel)
		}
	}
	fmt.Println(base, layer, xOffset, yOffset, totalDelta)

	return composedImage
}

func AlignImages(imageMap common.ImageMap, root string) error {
	err := common.WriteImage(path.Join(root, "output_v3_bg_align_00.jpg"), subtractImages(imageMap, common.BLUE, common.GREEN, 0, 0))
 	if err != nil { return err }
 	err = common.WriteImage(path.Join(root, "output_v3_bg_align_11.jpg"), subtractImages(imageMap, common.BLUE, common.GREEN, -1, 1))
 	if err != nil { return err }
 	err = common.WriteImage(path.Join(root, "output_v3_bg_align_22.jpg"), subtractImages(imageMap, common.BLUE, common.GREEN, -2, 2))
 	if err != nil { return err }

 	err = common.WriteImage(path.Join(root, "output_v3_br_align_00.jpg"), subtractImages(imageMap, common.BLUE, common.RED, 0, 0))
 	if err != nil { return err }
 	err = common.WriteImage(path.Join(root, "output_v3_bg_align_11.jpg"), subtractImages(imageMap, common.BLUE, common.RED, -1, 1))
 	if err != nil { return err }
 	err = common.WriteImage(path.Join(root, "output_v3_bg_align_22.jpg"), subtractImages(imageMap, common.BLUE, common.RED, -2, 2))
 	if err != nil { return err }

 	return nil
}

// CombineImages runs the v2 blending algorithm to combine a set of grayscale images into
// a "true" color image.
// func CombineImages(imageMap common.ImageMap, root string) error {
// 	composedImage := blendImage(imageMap)

// 	err := common.WriteImage(path.Join(root, "output_v2_alpha.jpg"), composedImage)
// 	if err != nil { return err }

// 	return nil
// }
