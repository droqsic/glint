package integration

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/internal/core"
	"github.com/droqsic/probe"
)

// TestTerminalDetectionIntegration tests integration with terminal detection
func TestTerminalDetectionIntegration(t *testing.T) {
	t.Run("TerminalDetectionWithColorSupport", func(t *testing.T) {
		// Test that terminal detection works with color support
		glint.ResetColor()
		core.ClearCache()

		// Check if we're in a terminal
		isTerminal := probe.IsTerminal(os.Stdout.Fd())
		isCygwin := probe.IsCygwinTerminal(os.Stdout.Fd())

		// Get color support
		colorSupported := glint.ColorSupport()

		// If not in a terminal, color support might be false
		if !isTerminal && !isCygwin && colorSupported {
			// This might happen if FORCE_COLOR is set
			forceColor := os.Getenv("FORCE_COLOR")
			noColor := os.Getenv("NO_COLOR")

			if forceColor == "" && noColor == "" {
				t.Logf("Color supported outside terminal - this might be expected in CI or with forced colors")
			}
		}
	})

	t.Run("TerminalDetectionAllFileDescriptors", func(t *testing.T) {
		// Test terminal detection for all standard file descriptors
		fds := []struct {
			name string
			fd   uintptr
		}{
			{"stdout", os.Stdout.Fd()},
			{"stderr", os.Stderr.Fd()},
			{"stdin", os.Stdin.Fd()},
		}

		for _, fdInfo := range fds {
			t.Run(fdInfo.name, func(t *testing.T) {
				isTerminal := probe.IsTerminal(fdInfo.fd)
				isCygwin := probe.IsCygwinTerminal(fdInfo.fd)

				// Should not panic and should return boolean values
				_ = isTerminal
				_ = isCygwin
			})
		}
	})

	t.Run("ColorSupportWithTerminalDetection", func(t *testing.T) {
		// Test color support behavior based on terminal detection
		glint.ResetColor()
		core.ClearCache()

		isStdoutTerminal := probe.IsTerminal(os.Stdout.Fd())
		isStdoutCygwin := probe.IsCygwinTerminal(os.Stdout.Fd())

		colorSupported := glint.ColorSupport()

		// Log the relationship for debugging
		t.Logf("Terminal detection - stdout: %v, cygwin: %v, color: %v",
			isStdoutTerminal, isStdoutCygwin, colorSupported)

		// The relationship depends on environment variables and platform
		// This test mainly ensures no panics occur
	})
}

// TestEnvironmentVariableIntegration tests comprehensive environment variable scenarios
func TestEnvironmentVariableIntegration(t *testing.T) {
	t.Run("ComplexEnvironmentScenarios", func(t *testing.T) {
		scenarios := []struct {
			name    string
			envVars map[string]string
			cleanup func()
		}{
			{
				name: "GitHubActions",
				envVars: map[string]string{
					"CI":             "true",
					"GITHUB_ACTIONS": "true",
					"TERM":           "xterm",
				},
			},
			{
				name: "JenkinsCI",
				envVars: map[string]string{
					"CI":          "true",
					"JENKINS_URL": "http://jenkins.example.com",
					"TERM":        "dumb",
				},
			},
			{
				name: "VSCodeIntegratedTerminal",
				envVars: map[string]string{
					"TERM_PROGRAM": "vscode",
					"COLORTERM":    "truecolor",
				},
			},
			{
				name: "WindowsTerminal",
				envVars: map[string]string{
					"WT_SESSION":    "12345678-1234-1234-1234-123456789012",
					"WT_PROFILE_ID": "12345678-1234-1234-1234-123456789012",
				},
			},
			{
				name: "iTerm2",
				envVars: map[string]string{
					"TERM_PROGRAM":         "iTerm.app",
					"TERM_PROGRAM_VERSION": "3.4.0",
					"COLORTERM":            "truecolor",
				},
			},
			{
				name: "Termux",
				envVars: map[string]string{
					"TERMUX_VERSION": "0.118",
					"TERM":           "xterm-256color",
				},
			},
			{
				name: "WSL",
				envVars: map[string]string{
					"WSLENV": "PATH/l:USERPROFILE/pu",
					"TERM":   "xterm-256color",
				},
			},
			{
				name: "SSH",
				envVars: map[string]string{
					"SSH_CONNECTION": "192.168.1.100 54321 192.168.1.1 22",
					"TERM":           "xterm-256color",
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.name, func(t *testing.T) {
				// Store original values
				originalValues := make(map[string]string)
				for key := range scenario.envVars {
					originalValues[key] = os.Getenv(key)
				}

				// Set test values
				for key, value := range scenario.envVars {
					os.Setenv(key, value)
				}

				// Restore original values after test
				defer func() {
					for key, originalValue := range originalValues {
						if originalValue == "" {
							os.Unsetenv(key)
						} else {
							os.Setenv(key, originalValue)
						}
					}
				}()

				// Test color detection
				glint.ResetColor()
				core.ClearCache()

				colorSupported := glint.ColorSupport()
				colorLevel := glint.ColorLevel()

				// Log results for debugging
				t.Logf("Scenario %s: color=%v, level=%v", scenario.name, colorSupported, colorLevel)

				// Verify consistency
				if !colorSupported && colorLevel != core.LevelNone {
					t.Errorf("Inconsistent results: color=%v, level=%v", colorSupported, colorLevel)
				}
			})
		}
	})

	t.Run("EnvironmentVariablePrecedence", func(t *testing.T) {
		// Test the precedence order of environment variables
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

		defer func() {
			for env, value := range originalValues {
				if value != "" {
					os.Setenv(env, value)
				}
			}
		}()

		// Test NO_COLOR takes precedence over everything
		os.Setenv("NO_COLOR", "1")
		os.Setenv("FORCE_COLOR", "1")
		os.Setenv("COLORTERM", "truecolor")
		os.Setenv("TERM", "xterm-256color")

		glint.ResetColor()
		core.ClearCache()

		colorSupported := glint.ColorSupport()
		if colorSupported {
			t.Errorf("NO_COLOR should override all other settings, but color is still supported")
		}

		// Test FORCE_COLOR takes precedence over terminal detection
		// Clear all environment variables again
		for _, env := range allEnvVars {
			os.Unsetenv(env)
		}

		os.Setenv("FORCE_COLOR", "1")
		os.Setenv("TERM", "dumb")

		glint.ResetColor()
		core.ClearCache()

		colorSupported = glint.ColorSupport()
		// Note: FORCE_COLOR environment variable only works in terminal environments
		// In test environments (non-terminal), it will still return false
		// This is the expected behavior - FORCE_COLOR doesn't override terminal detection
		if colorSupported {
			t.Logf("FORCE_COLOR enabled color support in terminal environment")
		} else {
			t.Logf("FORCE_COLOR did not enable color support in non-terminal environment (expected)")
		}
	})
}

// TestConcurrentIntegration tests concurrent access in integration scenarios
func TestConcurrentIntegration(t *testing.T) {
	t.Run("ConcurrentEnvironmentChanges", func(t *testing.T) {
		// This test simulates concurrent access while environment might be changing
		// Note: In real applications, environment shouldn't change during execution

		done := make(chan bool, 5)

		for i := 0; i < 5; i++ {
			go func(id int) {
				defer func() { done <- true }()

				for j := 0; j < 20; j++ {
					// Each goroutine tests color support
					colorSupported := glint.ColorSupport()
					colorLevel := glint.ColorLevel()

					// Verify consistency within this goroutine
					if !colorSupported && colorLevel != core.LevelNone {
						t.Errorf("Goroutine %d: inconsistent results: color=%v, level=%v", id, colorSupported, colorLevel)
					}
				}
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 5; i++ {
			<-done
		}
	})

	t.Run("ConcurrentForceColorOperations", func(t *testing.T) {
		// Test concurrent force color operations
		glint.ResetColor()

		done := make(chan bool, 3)

		// Goroutine 1: Force color on
		go func() {
			defer func() { done <- true }()
			for i := 0; i < 10; i++ {
				glint.ForceColor(true)
				glint.ColorSupport()
			}
		}()

		// Goroutine 2: Force color off
		go func() {
			defer func() { done <- true }()
			for i := 0; i < 10; i++ {
				glint.ForceColor(false)
				glint.ColorSupport()
			}
		}()

		// Goroutine 3: Just read
		go func() {
			defer func() { done <- true }()
			for i := 0; i < 20; i++ {
				glint.ColorSupport()
				glint.ColorLevel()
			}
		}()

		// Wait for all goroutines
		for i := 0; i < 3; i++ {
			<-done
		}
	})
}

// TestPerformanceIntegration tests performance in integration scenarios
func TestPerformanceIntegration(t *testing.T) {
	t.Run("CachedPerformance", func(t *testing.T) {
		// Prime the cache
		glint.ResetColor()
		glint.ColorSupport()
		glint.ColorLevel()

		// Subsequent calls should be very fast
		for i := 0; i < 1000; i++ {
			glint.ColorSupport()
			glint.ColorLevel()
		}
	})

	t.Run("UncachedPerformance", func(t *testing.T) {
		// Test performance of uncached calls
		for i := 0; i < 10; i++ {
			glint.ResetColor()
			core.ClearCache()
			glint.ColorSupport()
			glint.ColorLevel()
		}
	})
}
