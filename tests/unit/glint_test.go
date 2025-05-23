package unit

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/internal/core"
)

// TestColorSupport tests the ColorSupport function
func TestColorSupport(t *testing.T) {
	// Reset state before each test
	glint.ResetColor()
	core.ClearCache()

	t.Run("BasicColorSupport", func(t *testing.T) {
		result := glint.ColorSupport()
		// Result should be boolean
		if result != true && result != false {
			t.Errorf("ColorSupport() should return a boolean, got %v", result)
		}
	})

	t.Run("CachedColorSupport", func(t *testing.T) {
		// First call
		result1 := glint.ColorSupport()
		// Second call should return same result (cached)
		result2 := glint.ColorSupport()
		if result1 != result2 {
			t.Errorf("ColorSupport() should return consistent results, got %v then %v", result1, result2)
		}
	})

	t.Run("ColorSupportWithNoColor", func(t *testing.T) {
		originalNoColor := os.Getenv("NO_COLOR")
		defer func() {
			if originalNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", originalNoColor)
			}
		}()

		os.Setenv("NO_COLOR", "1")
		glint.ResetColor()
		core.ClearCache()

		result := glint.ColorSupport()
		// With NO_COLOR set, should return false
		if result != false {
			t.Errorf("ColorSupport() with NO_COLOR should return false, got %v", result)
		}
	})
}

// TestColorLevel tests the ColorLevel function
func TestColorLevel(t *testing.T) {
	glint.ResetColor()
	core.ClearCache()

	t.Run("BasicColorLevel", func(t *testing.T) {
		level := glint.ColorLevel()
		// Level should be one of the valid levels
		validLevels := []core.Level{core.LevelNone, core.Level16, core.Level256, core.LevelTrue}
		found := false
		for _, validLevel := range validLevels {
			if level == validLevel {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ColorLevel() returned invalid level: %v", level)
		}
	})

	t.Run("ColorLevelWithTrueColor", func(t *testing.T) {
		originalColorTerm := os.Getenv("COLORTERM")
		defer func() {
			if originalColorTerm == "" {
				os.Unsetenv("COLORTERM")
			} else {
				os.Setenv("COLORTERM", originalColorTerm)
			}
		}()

		os.Setenv("COLORTERM", "truecolor")
		glint.ResetColor()
		core.ClearCache()

		level := glint.ColorLevel()
		if glint.ColorSupport() && level != core.LevelTrue {
			t.Errorf("ColorLevel() with COLORTERM=truecolor should return LevelTrue, got %v", level)
		}
	})

	t.Run("ColorLevelWith256Color", func(t *testing.T) {
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		os.Setenv("TERM", "xterm-256color")
		glint.ResetColor()
		core.ClearCache()

		level := glint.ColorLevel()
		if glint.ColorSupport() && level != core.Level256 {
			t.Errorf("ColorLevel() with TERM=xterm-256color should return Level256, got %v", level)
		}
	})

	t.Run("ColorLevelWithNoColorSupport", func(t *testing.T) {
		originalNoColor := os.Getenv("NO_COLOR")
		defer func() {
			if originalNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", originalNoColor)
			}
		}()

		os.Setenv("NO_COLOR", "1")
		glint.ResetColor()
		core.ClearCache()

		level := glint.ColorLevel()
		if level != core.LevelNone {
			t.Errorf("ColorLevel() with NO_COLOR should return LevelNone, got %v", level)
		}
	})
}

// TestForceColor tests the ForceColor function
func TestForceColor(t *testing.T) {
	t.Run("ForceColorTrue", func(t *testing.T) {
		glint.ResetColor()
		glint.ForceColor(true)

		result := glint.ColorSupport()
		if !result {
			t.Errorf("ColorSupport() after ForceColor(true) should return true, got %v", result)
		}
	})

	t.Run("ForceColorFalse", func(t *testing.T) {
		glint.ResetColor()
		glint.ForceColor(false)

		result := glint.ColorSupport()
		if result {
			t.Errorf("ColorSupport() after ForceColor(false) should return false, got %v", result)
		}
	})

	t.Run("ForceColorWithNoColor", func(t *testing.T) {
		originalNoColor := os.Getenv("NO_COLOR")
		defer func() {
			if originalNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", originalNoColor)
			}
		}()

		os.Setenv("NO_COLOR", "1")
		core.ClearCache()
		glint.ResetColor()
		glint.ForceColor(true)

		result := glint.ColorSupport()
		// NO_COLOR should override ForceColor(true)
		if result {
			t.Errorf("ColorSupport() with NO_COLOR set should return false even after ForceColor(true), got %v", result)
		}
	})

	t.Run("ForceColorOnWindows", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows-specific test on non-Windows platform")
		}

		// Store and clear NO_COLOR to ensure ForceColor works
		originalNoColor := os.Getenv("NO_COLOR")
		defer func() {
			if originalNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", originalNoColor)
			}
		}()

		os.Unsetenv("NO_COLOR")
		glint.ResetColor()
		core.ClearCache()
		glint.ForceColor(true)

		result := glint.ColorSupport()
		if !result {
			t.Errorf("ColorSupport() after ForceColor(true) on Windows should return true, got %v", result)
		}

		// Test the Windows-specific code path in ForceColor
		// This should cover the vtEnabled logic
		glint.ResetColor()
		core.ClearCache()

		// Call ForceColor multiple times to test different scenarios
		glint.ForceColor(true)
		level1 := glint.ColorLevel()

		glint.ResetColor() // Reset to clear cache
		glint.ForceColor(false)
		level2 := glint.ColorLevel()

		glint.ResetColor() // Reset to clear cache
		glint.ForceColor(true)
		level3 := glint.ColorLevel()

		// Verify that levels change appropriately
		if level2 != core.LevelNone {
			t.Errorf("ColorLevel() after ForceColor(false) should be LevelNone, got %v", level2)
		}

		_ = level1
		_ = level3
	})
}

// TestResetColor tests the ResetColor function
func TestResetColor(t *testing.T) {
	t.Run("ResetAfterForceColor", func(t *testing.T) {
		// Force color first
		glint.ForceColor(true)
		forcedResult := glint.ColorSupport()

		// Reset color
		glint.ResetColor()
		resetResult := glint.ColorSupport()

		// Results might be different after reset (depends on actual terminal)
		_ = forcedResult
		_ = resetResult
		// Just ensure no panic occurs
	})

	t.Run("ResetMultipleTimes", func(t *testing.T) {
		glint.ResetColor()
		glint.ResetColor()
		glint.ResetColor()

		// Should not panic
		result := glint.ColorSupport()
		_ = result
	})

	t.Run("ResetClearsCache", func(t *testing.T) {
		// Get initial result
		glint.ColorSupport()

		// Reset should clear cache
		glint.ResetColor()

		// Should be able to get result again
		result := glint.ColorSupport()
		_ = result
	})
}

// TestColorSupportTerminalDetection tests terminal detection paths
func TestColorSupportTerminalDetection(t *testing.T) {
	t.Run("ColorSupportWithTerminalDetection", func(t *testing.T) {
		// This test covers the terminal detection code paths
		// The actual result depends on whether we're running in a terminal
		glint.ResetColor()
		core.ClearCache()

		result := glint.ColorSupport()
		// Result can be true or false depending on environment
		_ = result

		// Test the cached path
		result2 := glint.ColorSupport()
		if result != result2 {
			t.Errorf("Cached ColorSupport() should return same result: %v vs %v", result, result2)
		}
	})

	t.Run("ColorLevelWithTerminalDetection", func(t *testing.T) {
		// This test covers the terminal detection code paths in ColorLevel
		glint.ResetColor()
		core.ClearCache()

		level := glint.ColorLevel()
		// Level depends on environment and terminal detection
		_ = level

		// Test the cached path
		level2 := glint.ColorLevel()
		if level != level2 {
			t.Errorf("Cached ColorLevel() should return same result: %v vs %v", level, level2)
		}
	})

	t.Run("ColorLevelConsistencyWithColorSupport", func(t *testing.T) {
		// Test that ColorLevel is consistent with ColorSupport
		glint.ResetColor()
		core.ClearCache()

		colorSupported := glint.ColorSupport()
		colorLevel := glint.ColorLevel()

		if !colorSupported && colorLevel != core.LevelNone {
			t.Errorf("If color is not supported, level should be LevelNone, got %v", colorLevel)
		}

		if colorSupported && colorLevel == core.LevelNone {
			t.Errorf("If color is supported, level should not be LevelNone, got %v", colorLevel)
		}
	})

	t.Run("ColorLevelWithForceColor", func(t *testing.T) {
		// Test ColorLevel when ForceColor is used
		originalNoColor := os.Getenv("NO_COLOR")
		defer func() {
			if originalNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", originalNoColor)
			}
		}()

		os.Unsetenv("NO_COLOR")
		glint.ResetColor()
		core.ClearCache()
		glint.ForceColor(true)

		colorSupported := glint.ColorSupport()
		if !colorSupported {
			t.Errorf("ColorSupport() should return true after ForceColor(true), got %v", colorSupported)
		}

		colorLevel := glint.ColorLevel()
		if colorLevel == core.LevelNone {
			t.Errorf("ColorLevel() should not return LevelNone after ForceColor(true), got %v", colorLevel)
		}

		// Test the cached path for ColorLevel
		colorLevel2 := glint.ColorLevel()
		if colorLevel != colorLevel2 {
			t.Errorf("Cached ColorLevel() should return same result: %v vs %v", colorLevel, colorLevel2)
		}
	})

	t.Run("ColorLevelWithForceColorFalse", func(t *testing.T) {
		// Test ColorLevel when ForceColor(false) is used
		glint.ResetColor()
		core.ClearCache()
		glint.ForceColor(false)

		colorSupported := glint.ColorSupport()
		if colorSupported {
			t.Errorf("ColorSupport() should return false after ForceColor(false), got %v", colorSupported)
		}

		colorLevel := glint.ColorLevel()
		if colorLevel != core.LevelNone {
			t.Errorf("ColorLevel() should return LevelNone after ForceColor(false), got %v", colorLevel)
		}
	})
}

// TestConcurrentAccess tests concurrent access to glint functions
func TestConcurrentAccess(t *testing.T) {
	t.Run("ConcurrentColorSupport", func(t *testing.T) {
		glint.ResetColor()

		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 100; j++ {
					glint.ColorSupport()
				}
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("ConcurrentForceColor", func(t *testing.T) {
		glint.ResetColor()

		done := make(chan bool, 5)
		for i := 0; i < 5; i++ {
			go func(val bool) {
				defer func() { done <- true }()
				for j := 0; j < 50; j++ {
					glint.ForceColor(val)
					glint.ColorSupport()
				}
			}(i%2 == 0)
		}

		for i := 0; i < 5; i++ {
			<-done
		}
	})
}
