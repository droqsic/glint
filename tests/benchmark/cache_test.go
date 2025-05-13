package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/glint/feature"
)

// BenchmarkGetEnvCache benchmarks the GetEnvCache function before it has been cached
func BenchmarkGetEnvCache(b *testing.B) {
	resetEnvCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		feature.GetEnvCache(feature.EnvTerm)
	}
}

// BenchmarkGetEnvCacheCached benchmarks the GetEnvCache function after it has been cached
func BenchmarkGetEnvCacheCached(b *testing.B) {
	feature.GetEnvCache(feature.EnvTerm)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		feature.GetEnvCache(feature.EnvTerm)
	}
}

// resetEnvCache resets the cached environment variables
func resetEnvCache() {
	os.Setenv("GLINT_RESET_CACHE", "1")
	os.Unsetenv("GLINT_RESET_CACHE")
}
