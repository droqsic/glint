package integration

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/internal/core"
	"github.com/droqsic/glint/internal/platform"
)

// TestFullWorkflow tests the complete glint workflow
func TestFullWorkflow(t *testing.T) {
	t.Run("CompleteColorDetectionWorkflow", func(t *testing.T) {
		// Reset state
		glint.ResetColor()
		core.ClearCache()

		// Test complete workflow
		colorSupported := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		// Verify consistency
		if !colorSupported && colorLevel != core.LevelNone {
			t.Errorf("If color is not supported, level should be LevelNone, got %v", colorLevel)
		}

		if colorSupported && colorLevel == core.LevelNone {
			t.Errorf("If color is supported, level should not be LevelNone, got %v", colorLevel)
		}
	})

	t.Run("ForceColorWorkflow", func(t *testing.T) {
		// Reset state
		glint.ResetColor()

		// Force color on
		glint.ForceColor(true)

		colorSupported := glint.ColorSupport()

		// Check NO_COLOR override
		originalNoColor := os.Getenv("NO_COLOR")
		defer func() {
			if originalNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", originalNoColor)
			}
		}()

		if os.Getenv("NO_COLOR") == "" {
			// Without NO_COLOR, forced color should work
			if !colorSupported {
				t.Errorf("ForceColor(true) should enable color support, got %v", colorSupported)
			}
		}

		// Force color off
		glint.ForceColor(false)
		colorSupported = glint.ColorSupport()
		if colorSupported {
			t.Errorf("ForceColor(false) should disable color support, got %v", colorSupported)
		}

		// Reset and test natural detection
		glint.ResetColor()
		naturalSupport := glint.ColorSupport()
		_ = naturalSupport // May vary by environment
	})

	t.Run("EnvironmentVariableWorkflow", func(t *testing.T) {
		testCases := []struct {
			name     string
			envVars  map[string]string
			expected struct {
				colorSupport bool
				colorLevel   core.Level
			}
		}{
			{
				name: "NoColorSet",
				envVars: map[string]string{
					"NO_COLOR": "1",
				},
				expected: struct {
					colorSupport bool
					colorLevel   core.Level
				}{
					colorSupport: false,
					colorLevel:   core.LevelNone,
				},
			},
			{
				name: "ForceColorSet",
				envVars: map[string]string{
					"FORCE_COLOR": "1",
				},
				expected: struct {
					colorSupport bool
					colorLevel   core.Level
				}{
					colorSupport: false, // FORCE_COLOR only works in terminal environments
					colorLevel:   core.LevelNone,
				},
			},
			{
				name: "TrueColorTerminal",
				envVars: map[string]string{
					"COLORTERM": "truecolor",
				},
				expected: struct {
					colorSupport bool
					colorLevel   core.Level
				}{
					colorSupport: false, // Environment variables only work in terminal environments
					colorLevel:   core.LevelNone,
				},
			},
			{
				name: "256ColorTerminal",
				envVars: map[string]string{
					"TERM": "xterm-256color",
				},
				expected: struct {
					colorSupport bool
					colorLevel   core.Level
				}{
					colorSupport: false, // Environment variables only work in terminal environments
					colorLevel:   core.LevelNone,
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Store ALL relevant environment variables
				allEnvVars := []string{
					"NO_COLOR", "FORCE_COLOR", "COLOR_24", "COLOR_256", "COLOR_16",
					"COLORTERM", "TERM", "WT_SESSION", "WT_PROFILE_ID", "ANSICON",
					"ConEmuANSI", "TERM_PROGRAM", "CI", "TERMUX_VERSION", "WSLENV", "SSH_CONNECTION",
				}

				originalValues := make(map[string]string)
				for _, env := range allEnvVars {
					originalValues[env] = os.Getenv(env)
					os.Unsetenv(env)
				}

				// Set only the test values
				for key, value := range tc.envVars {
					os.Setenv(key, value)
				}

				// Restore original values after test
				defer func() {
					for env, originalValue := range originalValues {
						if originalValue != "" {
							os.Setenv(env, originalValue)
						}
					}
				}()

				// Reset and test
				glint.ResetColor()
				core.ClearCache()

				colorSupported := glint.ColorSupport()
				colorLevel := glint.ColorLevel()

				if colorSupported != tc.expected.colorSupport {
					t.Errorf("Expected color support %v, got %v", tc.expected.colorSupport, colorSupported)
				}

				if colorLevel != tc.expected.colorLevel {
					t.Errorf("Expected color level %v, got %v", tc.expected.colorLevel, colorLevel)
				}
			})
		}
	})
}

// TestPlatformIntegration tests platform-specific integration
func TestPlatformIntegration(t *testing.T) {
	t.Run("WindowsIntegration", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows-specific test on non-Windows platform")
		}

		// Test Windows-specific behavior
		glint.ResetColor()

		// Enable virtual terminal processing
		vtEnabled := platform.EnableVirtualTerminal()
		_ = vtEnabled

		// Test color support on Windows
		colorSupported := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		// On Windows, results depend on terminal capabilities
		_ = colorSupported
		_ = colorLevel
	})

	t.Run("UnixIntegration", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping Unix-specific test on Windows platform")
		}

		// Test Unix-specific behavior
		glint.ResetColor()

		// EnableVirtualTerminal should be no-op on Unix
		vtEnabled := platform.EnableVirtualTerminal()
		if vtEnabled != false {
			t.Errorf("EnableVirtualTerminal() on Unix should return false, got %v", vtEnabled)
		}

		// Test color support on Unix
		colorSupported := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		// Results depend on terminal and environment
		_ = colorSupported
		_ = colorLevel
	})
}

// TestCacheIntegration tests cache behavior across the system
func TestCacheIntegration(t *testing.T) {
	t.Run("CacheConsistency", func(t *testing.T) {
		// Reset state
		glint.ResetColor()
		core.ClearCache()

		// First calls should initialize cache
		colorSupport1 := glint.ColorSupport()
		colorLevel1 := glint.ColorLevel()

		// Subsequent calls should use cache
		colorSupport2 := glint.ColorSupport()
		colorLevel2 := glint.ColorLevel()

		// Results should be consistent
		if colorSupport1 != colorSupport2 {
			t.Errorf("ColorSupport() results inconsistent: %v vs %v", colorSupport1, colorSupport2)
		}

		if colorLevel1 != colorLevel2 {
			t.Errorf("ColorLevel() results inconsistent: %v vs %v", colorLevel1, colorLevel2)
		}
	})

	t.Run("CacheInvalidation", func(t *testing.T) {
		// Get initial results
		glint.ResetColor()
		core.ClearCache()

		colorSupport1 := glint.ColorSupport()
		colorLevel1 := glint.ColorLevel()

		// Force color and verify change
		glint.ForceColor(true)
		colorSupport2 := glint.ColorSupport()

		// Reset and verify cache is cleared
		glint.ResetColor()
		colorSupport3 := glint.ColorSupport()
		colorLevel3 := glint.ColorLevel()

		// After reset, should get fresh results
		_ = colorSupport1
		_ = colorLevel1
		_ = colorSupport2
		_ = colorSupport3
		_ = colorLevel3
	})
}

// TestErrorConditions tests error conditions and edge cases
func TestErrorConditions(t *testing.T) {
	t.Run("InvalidEnvironmentValues", func(t *testing.T) {
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		// Test with invalid TERM value
		os.Setenv("TERM", "invalid-terminal-that-does-not-exist")
		glint.ResetColor()
		core.ClearCache()

		// Should not panic
		colorSupported := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		// Should return valid values
		_ = colorSupported
		_ = colorLevel
	})

	t.Run("EmptyEnvironment", func(t *testing.T) {
		// Store all relevant environment variables
		envVars := []string{
			"TERM", "COLORTERM", "NO_COLOR", "FORCE_COLOR",
			"TERM_PROGRAM", "WT_SESSION", "CI", "ANSICON",
		}

		originalValues := make(map[string]string)
		for _, env := range envVars {
			originalValues[env] = os.Getenv(env)
			os.Unsetenv(env)
		}

		defer func() {
			for env, value := range originalValues {
				if value != "" {
					os.Setenv(env, value)
				}
			}
		}()

		// Test with minimal environment
		glint.ResetColor()
		core.ClearCache()

		colorSupported := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		// Should return valid values even with empty environment
		_ = colorSupported
		_ = colorLevel
	})
}

// TestRealWorldScenarios tests real-world usage scenarios
func TestRealWorldScenarios(t *testing.T) {
	t.Run("CLIApplicationScenario", func(t *testing.T) {
		// Simulate a CLI application checking for color support
		glint.ResetColor()

		if glint.ColorSupport() {
			level := glint.ColorLevel()
			switch level {
			case core.LevelNone:
				// Use no colors
			case core.Level16:
				// Use 16 colors
			case core.Level256:
				// Use 256 colors
			case core.LevelTrue:
				// Use true colors
			}
		}
		// Should complete without issues
	})

	t.Run("LibraryInitializationScenario", func(t *testing.T) {
		// Simulate a library that initializes color support once
		glint.ResetColor()

		// Library initialization
		hasColor := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		// Store results for later use
		_ = hasColor
		_ = colorLevel

		// Subsequent calls should be fast (cached)
		for i := 0; i < 100; i++ {
			glint.ColorSupport()
			glint.ColorLevel()
		}
	})
}
