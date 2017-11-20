package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lewchuk/gostitcher/algv1masking"
	"github.com/lewchuk/gostitcher/algv2blending"
	"github.com/lewchuk/gostitcher/algv3aligning"
	"github.com/lewchuk/gostitcher/common"
	"github.com/lewchuk/gostitcher/opus"
	"io/ioutil"
	"os"
	"path"
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

	configPath := path.Join(inputPath, "config.json")
	configS, err := ioutil.ReadFile(configPath)

	if err != nil {
		return err
	}

	config := common.ConfigFile{}

	if err := json.Unmarshal(configS, &config); err != nil {
		return fmt.Errorf("parsing config %s: %s", configPath, err)
	}

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
	pathPtr := flag.String("path", "", "path to a local folder with images and config.json")
	apiPtr := flag.Bool("api", false, "run using OPUS API")

	flag.Parse()

	var err error
	if (*apiPtr) {
		err = opus.CombineImages()
	} else {
		err = processImages(*pathPtr)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
