package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/glint/internal/core"
)

// BenchmarkSetEnvCache benchmarks the SetEnvCache function
func BenchmarkSetEnvCache(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ClearCache()
		core.SetEnvCache()
	}
}

// BenchmarkGetEnvCache benchmarks the GetEnvCache function before it has been cached
func BenchmarkGetEnvCache(b *testing.B) {
	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.GetEnvCache(core.EnvTerm)
	}
}

// BenchmarkGetEnvCacheCached benchmarks the GetEnvCache function after it has been cached
func BenchmarkGetEnvCacheCached(b *testing.B) {
	// Prime the cache
	core.GetEnvCache(core.EnvTerm)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.GetEnvCache(core.EnvTerm)
	}
}

// BenchmarkGetEnvCacheMultipleKeys benchmarks GetEnvCache with different environment variables
func BenchmarkGetEnvCacheMultipleKeys(b *testing.B) {
	keys := []string{
		core.EnvTerm,
		core.EnvColorTerm,
		core.EnvNoColor,
		core.EnvForceColor,
		core.EnvTermProgram,
		core.EnvWTSession,
		core.EnvANSICON,
	}

	// Prime the cache
	for _, key := range keys {
		core.GetEnvCache(key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			core.GetEnvCache(key)
		}
	}
}

// BenchmarkGetEnvCacheWithValues benchmarks GetEnvCache with actual environment values set
func BenchmarkGetEnvCacheWithValues(b *testing.B) {
	// Set up test environment variables
	testEnvs := map[string]string{
		core.EnvTerm:        "xterm-256color",
		core.EnvColorTerm:   "truecolor",
		core.EnvTermProgram: "iTerm.app",
		core.EnvWTSession:   "test-session",
	}

	// Store original values
	originalEnvs := make(map[string]string)
	for key := range testEnvs {
		originalEnvs[key] = os.Getenv(key)
	}

	// Set test values
	for key, value := range testEnvs {
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

	core.ClearCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for key := range testEnvs {
			core.GetEnvCache(key)
		}
	}
}

// BenchmarkClearCache benchmarks the ClearCache function
func BenchmarkClearCache(b *testing.B) {
	// Prime the cache first
	core.SetEnvCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ClearCache()
		// Re-initialize for next iteration
		core.SetEnvCache()
	}
}

// BenchmarkCacheInitialization benchmarks the full cache initialization process
func BenchmarkCacheInitialization(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ClearCache()
		// This will trigger cache initialization
		core.GetEnvCache(core.EnvTerm)
	}
}

// BenchmarkGetEnvCacheSequential benchmarks sequential access to different env vars
func BenchmarkGetEnvCacheSequential(b *testing.B) {
	// Prime the cache
	core.SetEnvCache()

	keys := []string{
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			core.GetEnvCache(key)
		}
	}
}

// BenchmarkGetEnvCacheEmpty benchmarks GetEnvCache with non-existent keys
func BenchmarkGetEnvCacheEmpty(b *testing.B) {
	// Prime the cache
	core.SetEnvCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.GetEnvCache("NON_EXISTENT_KEY")
	}
}
