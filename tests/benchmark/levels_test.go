package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/glint/internal/core"
)

// BenchmarkTerminalColorLevel benchmarks the TerminalColorLevel function before cache
func BenchmarkTerminalColorLevel(b *testing.B) {
	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelCached benchmarks the TerminalColorLevel function after cache
func BenchmarkTerminalColorLevelCached(b *testing.B) {
	// Prime the cache
	core.TerminalColorLevel()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelWithNoColor benchmarks with NO_COLOR environment variable
func BenchmarkTerminalColorLevelWithNoColor(b *testing.B) {
	originalNoColor := os.Getenv("NO_COLOR")
	os.Setenv("NO_COLOR", "1")
	defer func() {
		if originalNoColor == "" {
			os.Unsetenv("NO_COLOR")
		} else {
			os.Setenv("NO_COLOR", originalNoColor)
		}
	}()

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelWithForceColor benchmarks with FORCE_COLOR environment variable
func BenchmarkTerminalColorLevelWithForceColor(b *testing.B) {
	originalForceColor := os.Getenv("FORCE_COLOR")
	os.Setenv("FORCE_COLOR", "1")
	defer func() {
		if originalForceColor == "" {
			os.Unsetenv("FORCE_COLOR")
		} else {
			os.Setenv("FORCE_COLOR", originalForceColor)
		}
	}()

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelWithTrueColor benchmarks with truecolor support
func BenchmarkTerminalColorLevelWithTrueColor(b *testing.B) {
	originalColorTerm := os.Getenv("COLORTERM")
	os.Setenv("COLORTERM", "truecolor")
	defer func() {
		if originalColorTerm == "" {
			os.Unsetenv("COLORTERM")
		} else {
			os.Setenv("COLORTERM", originalColorTerm)
		}
	}()

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelWithXterm256 benchmarks with xterm-256color
func BenchmarkTerminalColorLevelWithXterm256(b *testing.B) {
	originalTerm := os.Getenv("TERM")
	os.Setenv("TERM", "xterm-256color")
	defer func() {
		if originalTerm == "" {
			os.Unsetenv("TERM")
		} else {
			os.Setenv("TERM", originalTerm)
		}
	}()

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelWithWindowsTerminal benchmarks with Windows Terminal
func BenchmarkTerminalColorLevelWithWindowsTerminal(b *testing.B) {
	originalWTSession := os.Getenv("WT_SESSION")
	os.Setenv("WT_SESSION", "test-session")
	defer func() {
		if originalWTSession == "" {
			os.Unsetenv("WT_SESSION")
		} else {
			os.Setenv("WT_SESSION", originalWTSession)
		}
	}()

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelWithITerm benchmarks with iTerm.app
func BenchmarkTerminalColorLevelWithITerm(b *testing.B) {
	originalTermProgram := os.Getenv("TERM_PROGRAM")
	os.Setenv("TERM_PROGRAM", "iTerm.app")
	defer func() {
		if originalTermProgram == "" {
			os.Unsetenv("TERM_PROGRAM")
		} else {
			os.Setenv("TERM_PROGRAM", originalTermProgram)
		}
	}()

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelWithCI benchmarks with CI environment
func BenchmarkTerminalColorLevelWithCI(b *testing.B) {
	originalCI := os.Getenv("CI")
	os.Setenv("CI", "true")
	defer func() {
		if originalCI == "" {
			os.Unsetenv("CI")
		} else {
			os.Setenv("CI", originalCI)
		}
	}()

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkTerminalColorLevelWithCustomColors benchmarks with custom color environment variables
func BenchmarkTerminalColorLevelWithCustomColors(b *testing.B) {
	testCases := []struct {
		name string
		env  string
	}{
		{"Custom24", "COLOR_24"},
		{"Custom256", "COLOR_256"},
		{"Custom16", "COLOR_16"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			originalValue := os.Getenv(tc.env)
			os.Setenv(tc.env, "1")
			defer func() {
				if originalValue == "" {
					os.Unsetenv(tc.env)
				} else {
					os.Setenv(tc.env, originalValue)
				}
			}()

			core.ClearCache()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				core.TerminalColorLevel()
			}
		})
	}
}

// BenchmarkTerminalColorLevelWithDumbTerm benchmarks with dumb terminal
func BenchmarkTerminalColorLevelWithDumbTerm(b *testing.B) {
	originalTerm := os.Getenv("TERM")
	os.Setenv("TERM", "dumb")
	defer func() {
		if originalTerm == "" {
			os.Unsetenv("TERM")
		} else {
			os.Setenv("TERM", originalTerm)
		}
	}()

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}
