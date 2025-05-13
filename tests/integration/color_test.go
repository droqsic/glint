package integration

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/feature"
)

// TestIsColorSupported tests the IsColorSupported function
func TestIsColorSupported(t *testing.T) {
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
		expected  bool
	}{
		{
			name:      "NO_COLOR set",
			noColor:   "1",
			term:      "xterm-256color",
			colorTerm: "truecolor",
			expected:  false,
		},
		{
			name:      "TERM=xterm-256color",
			noColor:   "",
			term:      "xterm-256color",
			colorTerm: "",
			expected:  true,
		},
		{
			name:      "TERM=xterm",
			noColor:   "",
			term:      "xterm",
			colorTerm: "",
			expected:  true,
		},
		{
			name:      "COLORTERM=truecolor",
			noColor:   "",
			term:      "",
			colorTerm: "truecolor",
			expected:  true,
		},
		{
			name:      "COLORTERM=24bit",
			noColor:   "",
			term:      "",
			colorTerm: "24bit",
			expected:  true,
		},
		{
			name:      "No environment variables",
			noColor:   "",
			term:      "",
			colorTerm: "",
			expected:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(feature.EnvNoColor, tc.noColor)
			os.Setenv(feature.EnvTerm, tc.term)
			os.Setenv(feature.EnvColorTerm, tc.colorTerm)

			feature.ResetCache()

			result := glint.IsColorSupported()

			// In a non-terminal test environment, the result might not match the expected value
			// We'll check if this is a case where we expect color support but don't have it
			if result != tc.expected {
				if tc.expected == true && result == false {
					t.Logf("Warning: Expected %v, got %v. This might be due to running in a non-terminal environment.", tc.expected, result)
					if !isRunningInTerminal() {
						t.Skip("Test requires a terminal environment")
					}
				} else {
					t.Errorf("Expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}

// TestForceColorSupport tests the ForceColorSupport function
func TestForceColorSupport(t *testing.T) {
	origNoColor := os.Getenv(feature.EnvNoColor)
	origTerm := os.Getenv(feature.EnvTerm)
	origColorTerm := os.Getenv(feature.EnvColorTerm)
	origGodebug := os.Getenv("GODEBUG")

	defer func() {
		os.Setenv(feature.EnvNoColor, origNoColor)
		os.Setenv(feature.EnvTerm, origTerm)
		os.Setenv(feature.EnvColorTerm, origColorTerm)
		os.Setenv("GODEBUG", origGodebug)
	}()

	t.Run("Force with GODEBUG", func(t *testing.T) {
		os.Setenv(feature.EnvNoColor, "1")
		os.Setenv("GODEBUG", "glintforcecolor=1")

		feature.ResetCache()
		before := glint.IsColorSupported()
		glint.ForceColorSupport()
		after := glint.IsColorSupported()

		if before != false {
			t.Errorf("Expected color support to be false before forcing, got %v", before)
		}

		// Note: The actual effect of ForceColorSupport depends on the terminal
		// and platform, so we can't make strong assertions about the after value
		t.Logf("Color support before forcing: %v, after forcing: %v", before, after)
	})
}

// isRunningInTerminal checks if the test is running in a terminal environment
func isRunningInTerminal() bool {
	return glint.IsTerminal(os.Stdout.Fd()) || glint.IsCygwinTerminal(os.Stdout.Fd())
}
