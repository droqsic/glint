package feature

import (
	"runtime"
	"strings"
	"sync"

	"github.com/droqsic/glint/platform"
)

// Level represents the color support level of the terminal.
type Level int

const (
	LevelNone Level = iota // Indicates no color support (0 colors)
	Level16                // Indicates basic ANSI color support (16 colors)
	Level256               // Indicates extended color support (256 colors)
	LevelTrue              // Indicates 24-bit RGB color support (16,777,216 colors)
)

// Pre-allocated color support level descriptions to avoid string allocations
var (
	LevelNoneDesc    = "No color support"
	Level16Desc      = "Basic ANSI color support (16 colors)"
	Level256Desc     = "Extended color support (256 colors)"
	LevelTrueDesc    = "24-bit RGB color support (16,777,216 colors)"
	LevelUnknownDesc = "Unknown color support level"
	LevelErrorDesc   = "Error: Output is not a terminal."
)

var (
	detectedLevel Level     // Caches the detected terminal color support level
	detectOnce    sync.Once // Ensures that the color detection is performed only once
)

// DetectColorSupport detects the color support level of the terminal.
// It returns a Level representing the level of color support.
// This function is optimized for zero allocations and maximum performance.
// The result is cached after the first call for subsequent calls to be extremely fast.
func DetectColorSupport() Level {
	detectOnce.Do(detectColorLevel)
	return detectedLevel
}

// GetColorDescription returns a human-readable description of the color support level.
// It takes a Level and returns a string describing that level.
// This function is optimized for zero allocations by using pre-allocated strings.
func GetColorDescription(level Level) string {
	switch detectedLevel {
	case LevelNone:
		return LevelNoneDesc
	case Level16:
		return Level16Desc
	case Level256:
		return Level256Desc
	case LevelTrue:
		return LevelTrueDesc
	default:
		return LevelUnknownDesc
	}
}

// detectColorLevel determines the color support level of the terminal.
// This function is called automatically by DetectColorSupport and is not meant to be called directly.
// It sets the detectedLevel variable based on environment variables and terminal capabilities.
func detectColorLevel() {
	envInitOnce.Do(initEnvCache)

	if GetEnvCache(EnvNoColor) != "" {
		detectedLevel = LevelNone
		return
	}

	term := GetEnvCache(EnvTerm)
	colorTerm := GetEnvCache(EnvColorTerm)

	if level, ok := terminalCheck(term, colorTerm); ok {
		detectedLevel = level
		return
	}

	if containsIgnoreCase(colorTerm, "truecolor") || containsIgnoreCase(colorTerm, "24bit") {
		detectedLevel = LevelTrue
		return
	}

	if containsIgnoreCase(term, "256color") {
		detectedLevel = Level256
		return
	}

	if term != "" {
		detectedLevel = Level16
		return
	}

	if runtime.GOOS == "windows" {
		platform.EnableVirtualTerminal()
		detectedLevel = Level16
	} else {
		detectedLevel = LevelNone
	}
}

// terminalCheck provides quick checks for common terminal types
// to avoid expensive string operations.
// It returns a Level and a boolean indicating if a match was found.
func terminalCheck(term, colorTerm string) (Level, bool) {
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return LevelTrue, true
	}

	if term == "xterm-256color" || term == "screen-256color" {
		return Level256, true
	}

	if term == "xterm" || term == "screen" || term == "vt100" {
		return Level16, true
	}

	return LevelNone, false
}

// containsIgnoreCase checks if a string contains a substring, ignoring case.
// This is optimized for performance with minimal comparisons.
// It returns true if the substring is found, false otherwise.
func containsIgnoreCase(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}

	if len(substr) > len(s) {
		return false
	}

	strLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(strLower, substrLower)
}
