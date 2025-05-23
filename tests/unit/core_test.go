package unit

import (
	"os"
	"testing"

	"github.com/droqsic/glint/internal/core"
)

// TestLevelConstants tests the Level constants
func TestLevelConstants(t *testing.T) {
	tests := []struct {
		level    core.Level
		expected int8
		name     string
	}{
		{core.LevelNone, 0, "LevelNone"},
		{core.Level16, 1, "Level16"},
		{core.Level256, 2, "Level256"},
		{core.LevelTrue, 3, "LevelTrue"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if int8(test.level) != test.expected {
				t.Errorf("%s should be %d, got %d", test.name, test.expected, int8(test.level))
			}
		})
	}
}

// TestLevelString tests the String method for Level constants
func TestLevelString(t *testing.T) {
	tests := []struct {
		level    core.Level
		expected string
		name     string
	}{
		{core.LevelNone, "No color support detected", "LevelNone"},
		{core.Level16, "ANSI color support (16 colors)", "Level16"},
		{core.Level256, "ANSI extended color support (256 colors)", "Level256"},
		{core.LevelTrue, "TrueColor support (24-bit RGB)", "LevelTrue"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.level.String()
			if result != test.expected {
				t.Errorf("%s.String() should return %q, got %q", test.name, test.expected, result)
			}
		})
	}

	// Test invalid level
	t.Run("InvalidLevel", func(t *testing.T) {
		invalidLevel := core.Level(99)
		result := invalidLevel.String()
		expected := "Unknown terminal color capability"
		if result != expected {
			t.Errorf("Invalid level should return %q, got %q", expected, result)
		}
	})
}

// TestSetEnvCache tests the SetEnvCache function
func TestSetEnvCache(t *testing.T) {
	t.Run("SetEnvCacheBasic", func(t *testing.T) {
		core.ClearCache()
		core.SetEnvCache()
		// Should not panic
	})

	t.Run("SetEnvCacheMultipleTimes", func(t *testing.T) {
		core.ClearCache()
		core.SetEnvCache()
		core.SetEnvCache()
		core.SetEnvCache()
		// Should not panic and should be idempotent
	})
}

// TestGetEnvCache tests the GetEnvCache function
func TestGetEnvCache(t *testing.T) {
	t.Run("GetEnvCacheKnownKey", func(t *testing.T) {
		core.ClearCache()
		result := core.GetEnvCache(core.EnvTerm)
		// Should return a string (might be empty)
		_ = result
	})

	t.Run("GetEnvCacheUnknownKey", func(t *testing.T) {
		core.ClearCache()
		result := core.GetEnvCache("UNKNOWN_ENV_VAR")
		// Should return empty string for unknown keys
		if result != "" {
			t.Errorf("GetEnvCache for unknown key should return empty string, got %q", result)
		}
	})

	t.Run("GetEnvCacheWithValue", func(t *testing.T) {
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		testValue := "test-terminal"
		os.Setenv("TERM", testValue)
		core.ClearCache()

		result := core.GetEnvCache(core.EnvTerm)
		if result != testValue {
			t.Errorf("GetEnvCache(EnvTerm) should return %q, got %q", testValue, result)
		}
	})

	t.Run("GetEnvCacheAllKnownKeys", func(t *testing.T) {
		core.ClearCache()
		knownKeys := []string{
			core.EnvTerm,
			core.EnvColorTerm,
			core.EnvNoColor,
			core.EnvForceColor,
			core.EnvTermProgram,
			core.EnvTermProgramVer,
			core.EnvWTSession,
			core.EnvWTProfileID,
			core.EnvANSICON,
			core.EnvConEmuANSI,
			core.EnvCI,
			core.EnvSSHConnection,
			core.EnvWSLEnv,
			core.EnvTermuxVersion,
			core.EnvCustomColor16,
			core.EnvCustomColor256,
			core.EnvCustomColor24,
		}

		for _, key := range knownKeys {
			result := core.GetEnvCache(key)
			// Should not panic and should return a string
			_ = result
		}
	})
}

// TestClearCache tests the ClearCache function
func TestClearCache(t *testing.T) {
	t.Run("ClearCacheBasic", func(t *testing.T) {
		core.SetEnvCache()
		core.ClearCache()
		// Should not panic
	})

	t.Run("ClearCacheMultipleTimes", func(t *testing.T) {
		core.ClearCache()
		core.ClearCache()
		core.ClearCache()
		// Should not panic
	})

	t.Run("ClearCacheResetsState", func(t *testing.T) {
		// Set up environment
		originalTerm := os.Getenv("TERM")
		defer func() {
			if originalTerm == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", originalTerm)
			}
		}()

		testValue1 := "test-terminal-1"
		os.Setenv("TERM", testValue1)
		core.ClearCache()

		result1 := core.GetEnvCache(core.EnvTerm)
		if result1 != testValue1 {
			t.Errorf("First GetEnvCache should return %q, got %q", testValue1, result1)
		}

		// Change environment and clear cache
		testValue2 := "test-terminal-2"
		os.Setenv("TERM", testValue2)
		core.ClearCache()

		result2 := core.GetEnvCache(core.EnvTerm)
		if result2 != testValue2 {
			t.Errorf("After cache clear, GetEnvCache should return %q, got %q", testValue2, result2)
		}
	})
}

// TestTerminalColorLevel tests the TerminalColorLevel function
func TestTerminalColorLevel(t *testing.T) {
	t.Run("TerminalColorLevelBasic", func(t *testing.T) {
		core.ClearCache()
		level := core.TerminalColorLevel()

		validLevels := []core.Level{core.LevelNone, core.Level16, core.Level256, core.LevelTrue}
		found := false
		for _, validLevel := range validLevels {
			if level == validLevel {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("TerminalColorLevel() returned invalid level: %v", level)
		}
	})

	t.Run("TerminalColorLevelWithNoColor", func(t *testing.T) {
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

		level := core.TerminalColorLevel()
		if level != core.LevelNone {
			t.Errorf("TerminalColorLevel() with NO_COLOR should return LevelNone, got %v", level)
		}
	})

	t.Run("TerminalColorLevelWithForceColor", func(t *testing.T) {
		originalForceColor := os.Getenv("FORCE_COLOR")
		originalNoColor := os.Getenv("NO_COLOR")
		defer func() {
			if originalForceColor == "" {
				os.Unsetenv("FORCE_COLOR")
			} else {
				os.Setenv("FORCE_COLOR", originalForceColor)
			}
			if originalNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", originalNoColor)
			}
		}()

		os.Unsetenv("NO_COLOR")
		os.Setenv("FORCE_COLOR", "1")
		core.ClearCache()

		level := core.TerminalColorLevel()
		if level != core.LevelTrue {
			t.Errorf("TerminalColorLevel() with FORCE_COLOR should return LevelTrue, got %v", level)
		}
	})

	t.Run("TerminalColorLevelWithCustomColors", func(t *testing.T) {
		testCases := []struct {
			env      string
			expected core.Level
		}{
			{"COLOR_24", core.LevelTrue},
			{"COLOR_256", core.Level256},
			{"COLOR_16", core.Level16},
		}

		for _, tc := range testCases {
			t.Run(tc.env, func(t *testing.T) {
				originalValue := os.Getenv(tc.env)
				originalNoColor := os.Getenv("NO_COLOR")
				originalForceColor := os.Getenv("FORCE_COLOR")
				defer func() {
					if originalValue == "" {
						os.Unsetenv(tc.env)
					} else {
						os.Setenv(tc.env, originalValue)
					}
					if originalNoColor == "" {
						os.Unsetenv("NO_COLOR")
					} else {
						os.Setenv("NO_COLOR", originalNoColor)
					}
					if originalForceColor == "" {
						os.Unsetenv("FORCE_COLOR")
					} else {
						os.Setenv("FORCE_COLOR", originalForceColor)
					}
				}()

				os.Unsetenv("NO_COLOR")
				os.Unsetenv("FORCE_COLOR")
				os.Setenv(tc.env, "1")
				core.ClearCache()

				level := core.TerminalColorLevel()
				if level != tc.expected {
					t.Errorf("TerminalColorLevel() with %s should return %v, got %v", tc.env, tc.expected, level)
				}
			})
		}
	})

	t.Run("TerminalColorLevelWithColorTerm", func(t *testing.T) {
		testCases := []struct {
			value    string
			expected core.Level
		}{
			{"truecolor", core.LevelTrue},
			{"24bit", core.LevelTrue},
			{"256color", core.Level256},
		}

		for _, tc := range testCases {
			t.Run(tc.value, func(t *testing.T) {
				// Store ALL relevant environment variables
				envVars := []string{
					"NO_COLOR", "FORCE_COLOR", "COLOR_24", "COLOR_256", "COLOR_16",
					"COLORTERM", "TERM", "WT_SESSION", "WT_PROFILE_ID", "ANSICON",
					"ConEmuANSI", "TERM_PROGRAM", "CI", "TERMUX_VERSION", "WSLENV", "SSH_CONNECTION",
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

				// Set only the COLORTERM variable we want to test
				os.Setenv("COLORTERM", tc.value)

				// Clear cache AFTER setting environment variables
				core.ClearCache()

				level := core.TerminalColorLevel()
				if level != tc.expected {
					t.Errorf("TerminalColorLevel() with COLORTERM=%s should return %v, got %v", tc.value, tc.expected, level)
				}
			})
		}
	})

	t.Run("TerminalColorLevelWithTerm", func(t *testing.T) {
		testCases := []struct {
			value    string
			expected core.Level
		}{
			{"xterm-256color", core.Level256},
			{"screen-256color", core.Level256},
			{"tmux-256color", core.Level256},
			{"rxvt-256color", core.Level256},
			{"xterm", core.Level16},
			{"screen", core.Level16},
			{"tmux", core.Level16},
			{"rxvt", core.Level16},
			{"dumb", core.LevelNone},
		}

		for _, tc := range testCases {
			t.Run(tc.value, func(t *testing.T) {
				// Store ALL relevant environment variables
				envVars := []string{
					"NO_COLOR", "FORCE_COLOR", "COLOR_24", "COLOR_256", "COLOR_16",
					"COLORTERM", "TERM", "WT_SESSION", "WT_PROFILE_ID", "ANSICON",
					"ConEmuANSI", "TERM_PROGRAM", "CI", "TERMUX_VERSION", "WSLENV", "SSH_CONNECTION",
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

				// Set only the TERM variable we want to test
				os.Setenv("TERM", tc.value)

				// Clear cache AFTER setting environment variables
				core.ClearCache()

				level := core.TerminalColorLevel()
				if level != tc.expected {
					t.Errorf("TerminalColorLevel() with TERM=%s should return %v, got %v", tc.value, tc.expected, level)
				}
			})
		}
	})

	t.Run("TerminalColorLevelWithSpecialEnvironments", func(t *testing.T) {
		testCases := []struct {
			env      string
			value    string
			expected core.Level
		}{
			{"WT_SESSION", "test-session", core.LevelTrue},
			{"WT_PROFILE_ID", "test-profile", core.LevelTrue},
			{"ANSICON", "1", core.Level256},
			{"ConEmuANSI", "ON", core.Level256},
			{"TERM_PROGRAM", "iTerm.app", core.LevelTrue},
			{"CI", "true", core.Level16},
			{"TERMUX_VERSION", "0.118", core.Level256},
			{"WSLENV", "PATH/l", core.Level256},
			{"SSH_CONNECTION", "192.168.1.1 22 192.168.1.2 12345", core.Level256},
		}

		for _, tc := range testCases {
			t.Run(tc.env, func(t *testing.T) {
				// Store ALL relevant environment variables
				envVars := []string{
					"NO_COLOR", "FORCE_COLOR", "COLOR_24", "COLOR_256", "COLOR_16",
					"COLORTERM", "TERM", "WT_SESSION", "WT_PROFILE_ID", "ANSICON",
					"ConEmuANSI", "TERM_PROGRAM", "CI", "TERMUX_VERSION", "WSLENV", "SSH_CONNECTION",
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

				// Set only the specific environment variable we want to test
				os.Setenv(tc.env, tc.value)

				// Clear cache AFTER setting environment variables
				core.ClearCache()

				level := core.TerminalColorLevel()
				if level != tc.expected {
					t.Errorf("TerminalColorLevel() with %s=%s should return %v, got %v", tc.env, tc.value, tc.expected, level)
				}
			})
		}
	})

	t.Run("TerminalColorLevelDefault", func(t *testing.T) {
		// Clear all environment variables that affect color detection
		envVars := []string{
			"NO_COLOR", "FORCE_COLOR", "COLOR_24", "COLOR_256", "COLOR_16",
			"COLORTERM", "TERM", "WT_SESSION", "WT_PROFILE_ID", "ANSICON",
			"ConEmuANSI", "TERM_PROGRAM", "CI", "TERMUX_VERSION", "WSLENV", "SSH_CONNECTION",
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

		core.ClearCache()
		level := core.TerminalColorLevel()

		// Default should be Level16
		if level != core.Level16 {
			t.Errorf("TerminalColorLevel() with no environment should return Level16, got %v", level)
		}
	})
}
