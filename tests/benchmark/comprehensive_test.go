package benchmark

import (
	"os"
	"runtime"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/internal/core"
	"github.com/droqsic/glint/internal/platform"
)

// BenchmarkFullGlintWorkflow benchmarks a complete glint workflow
func BenchmarkFullGlintWorkflow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Typical usage pattern
		if glint.ColorSupport() {
			level := glint.ColorLevel()
			switch level {
			case core.LevelNone:
				// No color
			case core.Level16:
				// 16 colors
			case core.Level256:
				// 256 colors
			case core.LevelTrue:
				// True color
			}
		}
	}
}

// BenchmarkGlintInitialization benchmarks the initialization cost of glint
func BenchmarkGlintInitialization(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resetColorCache()
		// First call triggers initialization
		glint.ColorSupport()
		glint.ColorLevel()
	}
}

// BenchmarkEnvironmentVariableScenarios benchmarks different environment variable scenarios
func BenchmarkEnvironmentVariableScenarios(b *testing.B) {
	scenarios := []struct {
		name string
		envs map[string]string
	}{
		{
			name: "NoColorSet",
			envs: map[string]string{"NO_COLOR": "1"},
		},
		{
			name: "ForceColorSet",
			envs: map[string]string{"FORCE_COLOR": "1"},
		},
		{
			name: "TrueColorTerminal",
			envs: map[string]string{
				"TERM":      "xterm-256color",
				"COLORTERM": "truecolor",
			},
		},
		{
			name: "WindowsTerminal",
			envs: map[string]string{
				"WT_SESSION":    "test-session",
				"WT_PROFILE_ID": "test-profile",
			},
		},
		{
			name: "ITerm",
			envs: map[string]string{
				"TERM_PROGRAM": "iTerm.app",
			},
		},
		{
			name: "CIEnvironment",
			envs: map[string]string{
				"CI": "true",
			},
		},
		{
			name: "CustomColors",
			envs: map[string]string{
				"COLOR_24": "1",
			},
		},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			// Store original values
			originalEnvs := make(map[string]string)
			for key := range scenario.envs {
				originalEnvs[key] = os.Getenv(key)
			}

			// Set test values
			for key, value := range scenario.envs {
				os.Setenv(key, value)
			}

			// Restore original values after benchmark
			defer func() {
				for key, originalValue := range originalEnvs {
					if originalValue == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, originalValue)
					}
				}
			}()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				resetColorCache()
				glint.ColorSupport()
				glint.ColorLevel()
			}
		})
	}
}

// BenchmarkPlatformSpecificBehavior benchmarks platform-specific behavior
func BenchmarkPlatformSpecificBehavior(b *testing.B) {
	b.Run("CurrentPlatform", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if runtime.GOOS == "windows" {
				platform.EnableVirtualTerminal()
			}
			glint.ColorSupport()
		}
	})

	b.Run("ForceColorOnWindows", func(b *testing.B) {
		if runtime.GOOS != "windows" {
			b.Skip("Skipping Windows-specific test on non-Windows platform")
		}

		for i := 0; i < b.N; i++ {
			glint.ForceColor(true)
			glint.ResetColor()
		}
	})

	b.Run("ForceColorOnUnix", func(b *testing.B) {
		if runtime.GOOS == "windows" {
			b.Skip("Skipping Unix-specific test on Windows platform")
		}

		for i := 0; i < b.N; i++ {
			glint.ForceColor(true)
			glint.ResetColor()
		}
	})
}

// BenchmarkCacheEfficiency benchmarks cache efficiency
func BenchmarkCacheEfficiency(b *testing.B) {
	b.Run("ColdCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			core.ClearCache()
			core.GetEnvCache(core.EnvTerm)
		}
	})

	b.Run("WarmCache", func(b *testing.B) {
		core.SetEnvCache()
		for i := 0; i < b.N; i++ {
			core.GetEnvCache(core.EnvTerm)
		}
	})

	b.Run("CacheHitRatio", func(b *testing.B) {
		keys := []string{
			core.EnvTerm,
			core.EnvColorTerm,
			core.EnvNoColor,
			core.EnvForceColor,
		}

		core.SetEnvCache()
		for i := 0; i < b.N; i++ {
			for _, key := range keys {
				core.GetEnvCache(key)
			}
		}
	})
}

// BenchmarkWorstCaseScenarios benchmarks worst-case scenarios
func BenchmarkWorstCaseScenarios(b *testing.B) {
	b.Run("FrequentCacheClear", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			glint.ColorSupport()
			core.ClearCache()
			glint.ColorLevel()
			glint.ResetColor()
		}
	})

	b.Run("FrequentForceToggle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			glint.ForceColor(true)
			glint.ColorSupport()
			glint.ForceColor(false)
			glint.ColorSupport()
			glint.ResetColor()
		}
	})

	b.Run("ManyEnvironmentChecks", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			core.TerminalColorLevel()
		}
	})
}

// BenchmarkRealWorldUsage benchmarks real-world usage patterns
func BenchmarkRealWorldUsage(b *testing.B) {
	b.Run("CLIToolPattern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Typical CLI tool pattern
			if glint.ColorSupport() {
				level := glint.ColorLevel()
				if level >= core.Level256 {
					// Use 256 colors
				} else if level >= core.Level16 {
					// Use 16 colors
				}
			}
		}
	})

	b.Run("LibraryPattern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Library initialization pattern
			colorSupported := glint.ColorSupport()
			if colorSupported {
				colorLevel := glint.ColorLevel()
				_ = colorLevel
			}
		}
	})

	b.Run("ConditionalColoringPattern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Conditional coloring pattern
			if glint.ColorSupport() {
				// Apply colors
			} else {
				// Plain text
			}
		}
	})
}

// BenchmarkEdgeCases benchmarks edge cases
func BenchmarkEdgeCases(b *testing.B) {
	b.Run("EmptyEnvironment", func(b *testing.B) {
		// Clear all relevant environment variables
		envVars := []string{
			"TERM", "COLORTERM", "NO_COLOR", "FORCE_COLOR",
			"TERM_PROGRAM", "WT_SESSION", "CI",
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

		for i := 0; i < b.N; i++ {
			resetColorCache()
			glint.ColorSupport()
			glint.ColorLevel()
		}
	})

	b.Run("InvalidEnvironmentValues", func(b *testing.B) {
		os.Setenv("TERM", "invalid-terminal-type")
		os.Setenv("COLORTERM", "invalid-color-type")
		defer func() {
			os.Unsetenv("TERM")
			os.Unsetenv("COLORTERM")
		}()

		for i := 0; i < b.N; i++ {
			resetColorCache()
			glint.ColorSupport()
			glint.ColorLevel()
		}
	})
}
