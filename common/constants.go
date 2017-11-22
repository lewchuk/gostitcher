package common

import (
	"image"
)

const (
	BLUE  = "BL1"
	GREEN = "GRN"
	RED   = "RED"
)

var Filters = [3]string{BLUE, GREEN, RED}

// The type for marshalling a config.json file from a folder with images.
type ImageConfig struct {
	Filename string `json:"filename"`
	Filter   string `json:"filter"`
}

type ConfigFile struct {
	Files []ImageConfig `json:"files"`
}

// ImageMap is a type alias for a map of filter strings to gray images.
type ImageMap = map[string]image.Gray
