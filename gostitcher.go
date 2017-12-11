package main

import (
	"flag"
	"fmt"
	"github.com/lewchuk/gostitcher/algv1masking"
	"github.com/lewchuk/gostitcher/algv2blending"
	"github.com/lewchuk/gostitcher/algv3aligning"
	"github.com/lewchuk/gostitcher/common"
	"github.com/lewchuk/gostitcher/opus"
	"os"
)

// https://space.stackexchange.com/questions/12510/cassinis-camera-continuum-band-filters
// A map of filter names to effective wavelengths.
var FilterMap = map[string]int{
	common.BLUE:  463,
	common.GREEN: 568,
	common.RED:   647,
}

func processImages(inputPath string) error {
	fmt.Printf("Processing: %s\n", inputPath)

	config, err := common.LoadConfig(inputPath)

	imageMap, err := common.LoadImages(config, inputPath)
	if err != nil {
		return err
	}

	if err = algv1masking.CombineImages(imageMap, inputPath); err != nil {
		return err
	}

	if err = algv2blending.CombineImages(imageMap, inputPath); err != nil {
		return err
	}

	if err = algv3aligning.AlignImages(imageMap, inputPath); err != nil {
		return err
	}

	return nil
}

func main() {
	pathPtr := flag.String("path", "", "path to a local folder with images and config.json. " +
		"Not compatible with the --api flag and will override any other flags if present.")
	apiPtr := flag.String("api", "", "use the OPUS API to pull down Cassini images to combine. Provide the output folder to place the images in.")
	cameraPtr := flag.String("camera", "narrow", "either 'narrow' (default) or 'wide' to select which Cassini camera. The same observation often includes images from both cameras so they cannot be fetched at once.")
	targetPtr := flag.String("target", "", "the target filter for the OPUS API (optional).")
	observationPtr := flag.String("observation", "", "the observation name for the OPUS API (optional).")
	extraPtr := flag.String("extra", "", "extra filters to add to the search URL, e.g. planet=Jupiter.")

	flag.Parse()

	var err error
	if *pathPtr != "" {
		err = processImages(*pathPtr)
	} else if *apiPtr != "" {
		if *cameraPtr != "narrow" && *cameraPtr != "wide" {
			err = fmt.Errorf("--camera must be either 'narrow' or 'wide': %s", *cameraPtr)
		} else {
			err = opus.ProcessImages(*apiPtr, *cameraPtr, *targetPtr, *observationPtr, *extraPtr)
		}
	} else {
		err = fmt.Errorf("Either --path parameter or --api flag must be provided.")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
