package main

import (
	"fmt"

	"github.com/droqsic/glint"
)

func main() {
	// Check if terminal supports colors
	fmt.Println("Terminal supports colors:", glint.ColorSupport())

	// Get color support level
	fmt.Println("Color support level:", glint.ColorLevel())

	// Force color support
	glint.ForceColor(true)

	// Reset color support
	glint.ResetColor()
}
