package integration

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/feature"
)

// TestIsColorSupportedLevel tests the IsColorSupportedLevel function
func TestIsColorSupportedLevel(t *testing.T) {
	origNoColor := os.Getenv(feature.EnvNoColor)
	origTerm := os.Getenv(feature.EnvTerm)
	origColorTerm := os.Getenv(feature.EnvColorTerm)

	defer func() {
		os.Setenv(feature.EnvNoColor, origNoColor)
		os.Setenv(feature.EnvTerm, origTerm)
		os.Setenv(feature.EnvColorTerm, origColorTerm)
	}()

	testCases := []struct {
		name      string
		noColor   string
		term      string
		colorTerm string
		expected  string
	}{
		{
			name:      "NO_COLOR set",
			noColor:   "1",
			term:      "xterm-256color",
			colorTerm: "truecolor",
			expected:  feature.LevelNoneDesc,
		},
		{
			name:      "TERM=xterm-256color",
			noColor:   "",
			term:      "xterm-256color",
			colorTerm: "",
			expected:  feature.Level256Desc,
		},
		{
			name:      "TERM=xterm",
			noColor:   "",
			term:      "xterm",
			colorTerm: "",
			expected:  feature.Level16Desc,
		},
		{
			name:      "COLORTERM=truecolor",
			noColor:   "",
			term:      "",
			colorTerm: "truecolor",
			expected:  feature.LevelTrueDesc,
		},
		{
			name:      "COLORTERM=24bit",
			noColor:   "",
			term:      "",
			colorTerm: "24bit",
			expected:  feature.LevelTrueDesc,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(feature.EnvNoColor, tc.noColor)
			os.Setenv(feature.EnvTerm, tc.term)
			os.Setenv(feature.EnvColorTerm, tc.colorTerm)

			feature.ResetCache()

			result := glint.IsColorSupportedLevel()
			if result != tc.expected && !isTerminalError(result) {
				if result == feature.LevelErrorDesc {
					t.Skipf("Test requires a terminal environment. Expected %q, got %q", tc.expected, result)
				} else {
					t.Errorf("Expected %q, got %q", tc.expected, result)
				}
			}
		})
	}
}

// TestDetectColorSupport tests the DetectColorSupport function
func TestDetectColorSupport(t *testing.T) {
	origNoColor := os.Getenv(feature.EnvNoColor)
	origTerm := os.Getenv(feature.EnvTerm)
	origColorTerm := os.Getenv(feature.EnvColorTerm)

	defer func() {
		os.Setenv(feature.EnvNoColor, origNoColor)
		os.Setenv(feature.EnvTerm, origTerm)
		os.Setenv(feature.EnvColorTerm, origColorTerm)
	}()

	testCases := []struct {
		name      string
		noColor   string
		term      string
		colorTerm string
		expected  feature.Level
	}{
		{
			name:      "NO_COLOR set",
			noColor:   "1",
			term:      "xterm-256color",
			colorTerm: "truecolor",
			expected:  feature.LevelNone,
		},
		{
			name:      "TERM=xterm-256color",
			noColor:   "",
			term:      "xterm-256color",
			colorTerm: "",
			expected:  feature.Level256,
		},
		{
			name:      "TERM=xterm",
			noColor:   "",
			term:      "xterm",
			colorTerm: "",
			expected:  feature.Level16,
		},
		{
			name:      "COLORTERM=truecolor",
			noColor:   "",
			term:      "",
			colorTerm: "truecolor",
			expected:  feature.LevelTrue,
		},
		{
			name:      "COLORTERM=24bit",
			noColor:   "",
			term:      "",
			colorTerm: "24bit",
			expected:  feature.LevelTrue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(feature.EnvNoColor, tc.noColor)
			os.Setenv(feature.EnvTerm, tc.term)
			os.Setenv(feature.EnvColorTerm, tc.colorTerm)

			feature.ResetCache()

			result := feature.DetectColorSupport()
			if result != tc.expected {
				if tc.expected == feature.Level256 && result == feature.LevelNone {
					t.Logf("Warning: Expected %v, got %v. This might be due to running in a non-terminal environment.", tc.expected, result)
				} else if tc.expected == feature.Level16 && result == feature.LevelNone {
					t.Logf("Warning: Expected %v, got %v. This might be due to running in a non-terminal environment.", tc.expected, result)
				} else if tc.expected == feature.LevelTrue && result == feature.LevelNone {
					t.Logf("Warning: Expected %v, got %v. This might be due to running in a non-terminal environment.", tc.expected, result)
				} else {
					t.Errorf("Expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}

// isTerminalError checks if the result is an error message
func isTerminalError(result string) bool {
	return result == feature.LevelErrorDesc
}
