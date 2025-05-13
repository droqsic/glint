package glint

import (
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/droqsic/glint/feature"
	"github.com/droqsic/glint/platform"
	"github.com/droqsic/probe"
)

var (
	debugForceColor string    // GODEBUG setting for forcing color support
	debugInitOnce   sync.Once // Ensures debug settings are initialized only once
)

var (
	colorLevelResult           string    // Caches the result of IsColorSupportedLevel
	colorLevelResultOnce       sync.Once // Ensures color level detection is performed only once
	isColorSupportedResult     bool      // Caches the result of IsColorSupported
	isColorSupportedResultOnce sync.Once // Ensures color support detection is performed only once
)

// IsTerminal checks if the file descriptor is a terminal.
// It returns true if the file descriptor is a terminal, or false otherwise.
// The parameter should be a file descriptor (e.g., os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()).
func IsTerminal(fd uintptr) bool {
	return probe.IsTerminal(fd)
}

// IsCygwinTerminal checks if the file descriptor is a Cygwin/MSYS2 terminal.
// It returns true if the file descriptor is a Cygwin/MSYS2 terminal, or false otherwise.
// The parameter should be a file descriptor (e.g., os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()).
func IsCygwinTerminal(fd uintptr) bool {
	return probe.IsCygwinTerminal(fd)
}

// IsColorSupported checks if the terminal supports color.
// It returns true if the terminal supports at least 16 colors.
// It returns false if the terminal does not support color or if the output is not a terminal.
// It is optimized for zero allocations and maximum performance, using caching and sync.Once after the first call.
func IsColorSupported() bool {
	isColorSupportedResultOnce.Do(isColorSupported)
	return isColorSupportedResult
}

// IsColorSupportedLevel checks the color support level of the terminal.
// It returns a string describing the color support level.
// It returns an error message if the output is not a terminal.
// It is optimized for zero allocations and maximum performance, using caching and sync.Once after the first call.
func IsColorSupportedLevel() string {
	colorLevelResultOnce.Do(func() {
		if !IsTerminal(os.Stdout.Fd()) && !IsCygwinTerminal(os.Stdout.Fd()) {
			colorLevelResult = feature.LevelErrorDesc
			return
		}

		level := feature.DetectColorSupport()
		colorLevelResult = feature.GetColorDescription(level)
	})

	return colorLevelResult
}

// ForceColorSupport forces the terminal to support color.
// It should be called before any color output.
// It is optimized for zero allocations and maximum performance, using caching and sync.Once after the first call.
// On Windows, it enables virtual terminal processing to support ANSI escape sequences.
// On other platforms, it's a no-op as ANSI escape sequences are supported by default.
func ForceColorSupport() {
	debugInitOnce.Do(initDebugSettings)
	if debugForceColor == "1" {
		feature.DetectColorSupport()
		return
	}

	if !IsTerminal(os.Stdout.Fd()) && !IsCygwinTerminal(os.Stdout.Fd()) {
		return
	}

	feature.DetectColorSupport()

	if runtime.GOOS == "windows" {
		platform.EnableVirtualTerminal()
	}
}

// isColorSupported detects if the terminal supports color.
// It is called once and caches the result for future use.
// It checks the NO_COLOR environment variable, the terminal type, and the color support level.
func isColorSupported() {
	if feature.GetEnvCache(feature.EnvNoColor) != "" {
		isColorSupportedResult = false
		return
	}

	fd := os.Stdout.Fd()
	if !IsTerminal(fd) && !IsCygwinTerminal(fd) {
		isColorSupportedResult = false
		return
	}

	level := feature.DetectColorSupport()
	isColorSupportedResult = level >= feature.Level16

}

// initDebugSettings initializes debug settings from GODEBUG environment variable.
// It parses the GODEBUG string to extract glint-specific debug settings.
// This function is called automatically when needed and is thread-safe.
func initDebugSettings() {
	godebug := os.Getenv("GODEBUG")
	if godebug == "" {
		return
	}

	parts := strings.Split(godebug, ",")
	for _, part := range parts {
		if strings.HasPrefix(part, "glintforcecolor=") {
			debugForceColor = strings.TrimPrefix(part, "glintforcecolor=")
			break
		}
	}
}
