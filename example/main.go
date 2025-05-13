package main

import (
	"fmt"

	"github.com/droqsic/glint"
)

func main() {

	// Check if terminal supports colors
	if glint.IsColorSupported() {
		fmt.Println("Terminal supports colors")
	} else {
		fmt.Println("Terminal does not support colors")
	}

	// Get a human-readable description of the color support level
	colorLevel := glint.IsColorSupportedLevel()
	fmt.Printf("Color support level: %s\n", colorLevel)

	// Force color support
	fmt.Println("Before forcing:", glint.IsColorSupported())
	glint.ForceColorSupport()
	fmt.Println("After forcing:", glint.IsColorSupported())

}
