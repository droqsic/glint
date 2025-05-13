package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/feature"
)

// BenchmarkIsColorSupportedLevel benchmarks the IsColorSupportedLevel function before it has been cached
func BenchmarkIsColorSupportedLevel(b *testing.B) {
	resetColorLevelCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.IsColorSupportedLevel()
	}
}

// BenchmarkIsColorSupportedLevelCached benchmarks the IsColorSupportedLevel function after it has been cached
func BenchmarkIsColorSupportedLevelCached(b *testing.B) {
	glint.IsColorSupportedLevel()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.IsColorSupportedLevel()
	}
}

// BenchmarkDetectColorSupport benchmarks the DetectColorSupport function before it has been cached
func BenchmarkDetectColorSupport(b *testing.B) {
	resetDetectColorCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		feature.DetectColorSupport()
	}
}

// BenchmarkDetectColorSupportCached benchmarks the DetectColorSupport function after it has been cached
func BenchmarkDetectColorSupportCached(b *testing.B) {
	feature.DetectColorSupport()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		feature.DetectColorSupport()
	}
}

// resetColorLevelCache resets the cached result of IsColorSupportedLevel
func resetColorLevelCache() {
	os.Setenv("GLINT_RESET_CACHE", "1")
	os.Unsetenv("GLINT_RESET_CACHE")
}

// resetDetectColorCache resets the cached result of DetectColorSupport
func resetDetectColorCache() {
	os.Setenv("GLINT_RESET_CACHE", "1")
	os.Unsetenv("GLINT_RESET_CACHE")
}
