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
	OffsetX  int    `json:"offsetX"`
	OffsetY  int    `json:"offsetY"`
}

type ConfigFile struct {
	Files     []ImageConfig `json:"files"`
	MaxOffset int           `json:"maxOffset"`
}

type LoadedConfig struct {
	Config ImageConfig
	Image  image.Gray
}

// ImageMap is a type alias for a map of filter strings to gray images.
type ImageMap = map[string]LoadedConfig

type ImageFilenameMap = map[string]string
