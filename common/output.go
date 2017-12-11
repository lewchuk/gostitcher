package common

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
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

func WriteConfig(root string, config ConfigFile) error {
	configJson, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot serialize config json %s: %s", config, err)
	}

	configPath := fmt.Sprintf("%s/config.json", root)
	err = ioutil.WriteFile(configPath, configJson, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write config json to %s: %s", configPath, err)
	}

	return nil
}
