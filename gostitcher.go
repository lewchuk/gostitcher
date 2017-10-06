package main

import (
	"encoding/json"
	"fmt"
	"gostitcher/algv1masking"
	"gostitcher/common"
	"io/ioutil"
	"os"
	"path"
)

// https://space.stackexchange.com/questions/12510/cassinis-camera-continuum-band-filters
// A map of filter names to effective wavelengths.
var FilterMap = map[string]int{
	common.BLUE: 463,
	common.GREEN: 568,
	common.RED: 647,
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

	return algv1masking.CombineImages(config, inputPath)
}

func main() {
	err := processImages(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
