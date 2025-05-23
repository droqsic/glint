package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/internal/core"
)

// BenchmarkColorSupport benchmarks the ColorSupport function before it has been cached
func BenchmarkColorSupport(b *testing.B) {
	resetColorCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ColorSupport()
	}
}

// BenchmarkColorSupportCached benchmarks the ColorSupport function after it has been cached
func BenchmarkColorSupportCached(b *testing.B) {
	// Prime the cache
	glint.ColorSupport()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ColorSupport()
	}
}

// BenchmarkColorLevel benchmarks the ColorLevel function before it has been cached
func BenchmarkColorLevel(b *testing.B) {
	resetColorCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ColorLevel()
	}
}

// BenchmarkColorLevelCached benchmarks the ColorLevel function after it has been cached
func BenchmarkColorLevelCached(b *testing.B) {
	// Prime the cache
	glint.ColorLevel()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ColorLevel()
	}
}

// BenchmarkForceColor benchmarks the ForceColor function
func BenchmarkForceColor(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ForceColor(true)
		glint.ForceColor(false)
	}
}

// BenchmarkForceColorTrue benchmarks the ForceColor function with true value
func BenchmarkForceColorTrue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ForceColor(true)
	}
}

// BenchmarkForceColorFalse benchmarks the ForceColor function with false value
func BenchmarkForceColorFalse(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ForceColor(false)
	}
}

// BenchmarkResetColor benchmarks the ResetColor function
func BenchmarkResetColor(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ResetColor()
	}
}

// BenchmarkColorSupportWithForce benchmarks ColorSupport when force is enabled
func BenchmarkColorSupportWithForce(b *testing.B) {
	glint.ForceColor(true)
	defer glint.ResetColor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ColorSupport()
	}
}

// BenchmarkColorSupportWithNoColor benchmarks ColorSupport with NO_COLOR set
func BenchmarkColorSupportWithNoColor(b *testing.B) {
	originalNoColor := os.Getenv("NO_COLOR")
	os.Setenv("NO_COLOR", "1")
	defer func() {
		if originalNoColor == "" {
			os.Unsetenv("NO_COLOR")
		} else {
			os.Setenv("NO_COLOR", originalNoColor)
		}
	}()

	resetColorCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ColorSupport()
	}
}

// BenchmarkColorLevelWithDifferentTerms benchmarks ColorLevel with different TERM values
func BenchmarkColorLevelWithDifferentTerms(b *testing.B) {
	terms := []string{"xterm-256color", "xterm", "screen-256color", "dumb", ""}

	for _, term := range terms {
		b.Run("TERM="+term, func(b *testing.B) {
			originalTerm := os.Getenv("TERM")
			if term == "" {
				os.Unsetenv("TERM")
			} else {
				os.Setenv("TERM", term)
			}
			defer func() {
				if originalTerm == "" {
					os.Unsetenv("TERM")
				} else {
					os.Setenv("TERM", originalTerm)
				}
			}()

			resetColorCache()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				glint.ColorLevel()
			}
		})
	}
}

// resetColorCache resets the cached color support and level results
func resetColorCache() {
	glint.ResetColor()
	core.ClearCache()
}
