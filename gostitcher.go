package main

import (
	"fmt"
	"os"
)

func main() {
	inputPath := os.Args[1]
	fmt.Printf("Processing image at: %s\n", inputPath)
}
