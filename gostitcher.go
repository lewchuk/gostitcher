package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// The type for marshalling a config.json file from a folder with images.
type ImageConfig struct {
	Filename string `json:"filename"`
	Filter   string `json:"filter"`
}

type ConfigFile struct {
	Files []ImageConfig `json:"files"`
}

func main() {
	inputPath := os.Args[1]
	fmt.Printf("Processing image at: %s\n", inputPath)

	configS, err := ioutil.ReadFile(path.Join(inputPath, "config.json"))

	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(-1)
	}

	config := ConfigFile{}
	err2 := json.Unmarshal(configS, &config)

	if err2 != nil {
		fmt.Println("Error parsing config:", err)
		os.Exit(-1)
	}

	fmt.Println(config.Files[0].Filename)
	fmt.Println(config.Files[0].Filter)
}
