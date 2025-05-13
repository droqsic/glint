package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
)

// BenchmarkIsColorSupported benchmarks the IsColorSupported function before it has been cached
func BenchmarkIsColorSupported(b *testing.B) {
	resetColorSupportCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.IsColorSupported()
	}
}

// BenchmarkIsColorSupportedCached benchmarks the IsColorSupported function after it has been cached
func BenchmarkIsColorSupportedCached(b *testing.B) {
	glint.IsColorSupported()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.IsColorSupported()
	}
}

// BenchmarkForceColorSupport benchmarks the ForceColorSupport function
func BenchmarkForceColorSupport(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ForceColorSupport()
	}
}

// resetColorSupportCache resets the cached result of IsColorSupported
func resetColorSupportCache() {
	os.Setenv("GLINT_RESET_CACHE", "1")
	os.Unsetenv("GLINT_RESET_CACHE")
}
