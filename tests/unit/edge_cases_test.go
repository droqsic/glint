package unit

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/internal/core"
)

// TestEdgeCases tests edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	t.Run("EmptyStringEnvironmentValues", func(t *testing.T) {
		envVars := []string{
			"TERM", "COLORTERM", "NO_COLOR", "FORCE_COLOR",
			"TERM_PROGRAM", "WT_SESSION", "ANSICON", "ConEmuANSI",
		}

		// Store original values
		originalValues := make(map[string]string)
		for _, env := range envVars {
			originalValues[env] = os.Getenv(env)
			os.Setenv(env, "") // Set to empty string
		}

		defer func() {
			for env, value := range originalValues {
				if value == "" {
					os.Unsetenv(env)
				} else {
					os.Setenv(env, value)
				}
			}
		}()

		glint.ResetColor()
		core.ClearCache()

		// Should handle empty strings gracefully
		colorSupported := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		_ = colorSupported
		_ = colorLevel
	})

	t.Run("WhitespaceEnvironmentValues", func(t *testing.T) {
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		// Test with whitespace values
		testValues := []string{" ", "\t", "\n", "  \t\n  "}
		for _, value := range testValues {
			os.Setenv("TERM", value)
			glint.ResetColor()
			core.ClearCache()

			colorSupported := glint.ColorSupport()
			colorLevel := glint.ColorLevel()

			_ = colorSupported
			_ = colorLevel
		}
	})

	t.Run("VeryLongEnvironmentValues", func(t *testing.T) {
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		// Test with very long environment value
		longValue := ""
		for i := 0; i < 1000; i++ {
			longValue += "x"
		}

		os.Setenv("TERM", longValue)
		glint.ResetColor()
		core.ClearCache()

		// Should handle long values gracefully
		colorSupported := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		_ = colorSupported
		_ = colorLevel
	})

	t.Run("SpecialCharactersInEnvironment", func(t *testing.T) {
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		// Test with special characters
		specialValues := []string{
			"term-with-unicode-🌈",
			"term\x00with\x00nulls",
			"term\twith\ttabs",
			"term\nwith\nnewlines",
			"term with spaces",
			"TERM_WITH_CAPS",
			"term-with-dashes",
			"term_with_underscores",
			"term.with.dots",
		}

		for _, value := range specialValues {
			os.Setenv("TERM", value)
			glint.ResetColor()
			core.ClearCache()

			colorSupported := glint.ColorSupport()
			colorLevel := glint.ColorLevel()

			_ = colorSupported
			_ = colorLevel
		}
	})
}

// TestBoundaryConditions tests boundary conditions
func TestBoundaryConditions(t *testing.T) {
	t.Run("ForceColorBoundaryValues", func(t *testing.T) {
		originalForceColor := os.Getenv("FORCE_COLOR")
		defer func() {
			if originalForceColor == "" {
				os.Unsetenv("FORCE_COLOR")
			} else {
				os.Setenv("FORCE_COLOR", originalForceColor)
			}
		}()

		// Test various FORCE_COLOR values
		testValues := []string{
			"0", "1", "2", "3", "4", "true", "false", "yes", "no",
			"TRUE", "FALSE", "YES", "NO", "on", "off", "ON", "OFF",
		}

		for _, value := range testValues {
			os.Setenv("FORCE_COLOR", value)
			glint.ResetColor()
			core.ClearCache()

			colorSupported := glint.ColorSupport()
			colorLevel := glint.ColorLevel()

			_ = colorSupported
			_ = colorLevel
		}
	})

	t.Run("NoColorBoundaryValues", func(t *testing.T) {
		originalNoColor := os.Getenv("NO_COLOR")
		defer func() {
			if originalNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", originalNoColor)
			}
		}()

		// Test various NO_COLOR values
		testValues := []string{
			"0", "1", "true", "false", "yes", "no", "anything", "",
		}

		for _, value := range testValues {
			if value == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", value)
			}
			glint.ResetColor()
			core.ClearCache()

			colorSupported := glint.ColorSupport()
			if value != "" && colorSupported {
				t.Errorf("With NO_COLOR=%q, color support should be false, got %v", value, colorSupported)
			}
		}
	})

	t.Run("TerminalTypeEdgeCases", func(t *testing.T) {
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		// Test edge case terminal types
		edgeCases := []string{
			"",                     // Empty
			"dumb",                 // Explicitly dumb
			"unknown",              // Unknown terminal
			"xterm-",               // Incomplete
			"256color",             // Missing prefix
			"xterm-256",            // Missing suffix
			"XTERM-256COLOR",       // Wrong case
			"xterm-256color-extra", // Extra suffix
		}

		for _, termType := range edgeCases {
			if termType == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", termType)
			}
			glint.ResetColor()
			core.ClearCache()

			colorSupported := glint.ColorSupport()
			colorLevel := glint.ColorLevel()

			_ = colorSupported
			_ = colorLevel
		}
	})
}

// TestRaceConditions tests potential race conditions
func TestRaceConditions(t *testing.T) {
	t.Run("ConcurrentResetAndAccess", func(t *testing.T) {
		done := make(chan bool, 10)

		// Start multiple goroutines that reset and access
		for i := 0; i < 5; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 20; j++ {
					glint.ResetColor()
					glint.ColorSupport()
					glint.ColorLevel()
				}
			}()
		}

		// Start multiple goroutines that only access
		for i := 0; i < 5; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 50; j++ {
					glint.ColorSupport()
					glint.ColorLevel()
				}
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("ConcurrentForceColorAndAccess", func(t *testing.T) {
		done := make(chan bool, 10)

		// Start multiple goroutines that force color
		for i := 0; i < 5; i++ {
			go func(val bool) {
				defer func() { done <- true }()
				for j := 0; j < 20; j++ {
					glint.ForceColor(val)
					glint.ColorSupport()
				}
			}(i%2 == 0)
		}

		// Start multiple goroutines that only access
		for i := 0; i < 5; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 50; j++ {
					glint.ColorSupport()
					glint.ColorLevel()
				}
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("ConcurrentCacheOperations", func(t *testing.T) {
		done := make(chan bool, 10)

		// Start multiple goroutines that clear cache
		for i := 0; i < 3; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 10; j++ {
					core.ClearCache()
					core.SetEnvCache()
				}
			}()
		}

		// Start multiple goroutines that access cache
		for i := 0; i < 7; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 30; j++ {
					core.GetEnvCache(core.EnvTerm)
					core.GetEnvCache(core.EnvColorTerm)
				}
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// TestMemoryLeaks tests for potential memory leaks
func TestMemoryLeaks(t *testing.T) {
	t.Run("RepeatedOperations", func(t *testing.T) {
		// Perform many operations to check for memory leaks
		for i := 0; i < 1000; i++ {
			glint.ResetColor()
			glint.ColorSupport()
			glint.ColorLevel()
			glint.ForceColor(true)
			glint.ForceColor(false)
			core.ClearCache()
			core.SetEnvCache()
			core.GetEnvCache(core.EnvTerm)
		}
	})

	t.Run("RepeatedEnvironmentChanges", func(t *testing.T) {
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		// Repeatedly change environment and test
		termValues := []string{"xterm", "xterm-256color", "screen", "tmux", "dumb"}
		for i := 0; i < 200; i++ {
			termValue := termValues[i%len(termValues)]
			os.Setenv("TERM", termValue)
			glint.ResetColor()
			core.ClearCache()
			glint.ColorSupport()
			glint.ColorLevel()
		}
	})
}

// TestErrorRecovery tests error recovery scenarios
func TestErrorRecovery(t *testing.T) {
	t.Run("RecoveryAfterPanic", func(t *testing.T) {
		// This test ensures the system can recover after any potential panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()

		// Try various operations that might cause issues
		glint.ResetColor()
		glint.ForceColor(true)
		glint.ForceColor(false)
		glint.ResetColor()

		core.ClearCache()
		core.SetEnvCache()
		core.GetEnvCache("INVALID_KEY")
		core.ClearCache()

		glint.ColorSupport()
		glint.ColorLevel()
	})

	t.Run("StateConsistencyAfterErrors", func(t *testing.T) {
		// Ensure state remains consistent even after potential errors
		glint.ResetColor()

		// Get baseline
		baseline1 := glint.ColorSupport()
		baseline2 := glint.ColorLevel()

		// Perform various operations
		glint.ForceColor(true)
		glint.ForceColor(false)
		glint.ResetColor()

		// Check consistency
		result1 := glint.ColorSupport()
		result2 := glint.ColorLevel()

		if baseline1 != result1 {
			t.Errorf("ColorSupport consistency check failed: %v vs %v", baseline1, result1)
		}
		if baseline2 != result2 {
			t.Errorf("ColorLevel consistency check failed: %v vs %v", baseline2, result2)
		}
	})
}
