package benchmark

import (
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/internal/core"
)

// BenchmarkColorSupportAllocs benchmarks memory allocations for ColorSupport
func BenchmarkColorSupportAllocs(b *testing.B) {
	// Prime the cache to test cached performance
	glint.ColorSupport()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ColorSupport()
	}
}

// BenchmarkColorLevelAllocs benchmarks memory allocations for ColorLevel
func BenchmarkColorLevelAllocs(b *testing.B) {
	// Prime the cache to test cached performance
	glint.ColorLevel()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ColorLevel()
	}
}

// BenchmarkGetEnvCacheAllocs benchmarks memory allocations for GetEnvCache
func BenchmarkGetEnvCacheAllocs(b *testing.B) {
	// Prime the cache
	core.SetEnvCache()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.GetEnvCache(core.EnvTerm)
	}
}

// BenchmarkTerminalColorLevelAllocs benchmarks memory allocations for TerminalColorLevel
func BenchmarkTerminalColorLevelAllocs(b *testing.B) {
	// Prime the cache
	core.TerminalColorLevel()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.TerminalColorLevel()
	}
}

// BenchmarkForceColorAllocs benchmarks memory allocations for ForceColor
func BenchmarkForceColorAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ForceColor(true)
	}
}

// BenchmarkResetColorAllocs benchmarks memory allocations for ResetColor
func BenchmarkResetColorAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.ResetColor()
	}
}

// BenchmarkSetEnvCacheAllocs benchmarks memory allocations for SetEnvCache
func BenchmarkSetEnvCacheAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ClearCache()
		core.SetEnvCache()
	}
}

// BenchmarkClearCacheAllocs benchmarks memory allocations for ClearCache
func BenchmarkClearCacheAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.SetEnvCache()
		core.ClearCache()
	}
}

// BenchmarkColorSupportUncachedAllocs benchmarks memory allocations for uncached ColorSupport
func BenchmarkColorSupportUncachedAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resetColorCache()
		glint.ColorSupport()
	}
}

// BenchmarkColorLevelUncachedAllocs benchmarks memory allocations for uncached ColorLevel
func BenchmarkColorLevelUncachedAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resetColorCache()
		glint.ColorLevel()
	}
}

// BenchmarkGetEnvCacheUncachedAllocs benchmarks memory allocations for uncached GetEnvCache
func BenchmarkGetEnvCacheUncachedAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ClearCache()
		core.GetEnvCache(core.EnvTerm)
	}
}

// BenchmarkTerminalColorLevelUncachedAllocs benchmarks memory allocations for uncached TerminalColorLevel
func BenchmarkTerminalColorLevelUncachedAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ClearCache()
		core.TerminalColorLevel()
	}
}

// BenchmarkMultipleEnvCacheAllocs benchmarks memory allocations for multiple env cache accesses
func BenchmarkMultipleEnvCacheAllocs(b *testing.B) {
	keys := []string{
		core.EnvTerm,
		core.EnvColorTerm,
		core.EnvNoColor,
		core.EnvForceColor,
		core.EnvTermProgram,
	}

	// Prime the cache
	core.SetEnvCache()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			core.GetEnvCache(key)
		}
	}
}

// BenchmarkFullWorkflowAllocs benchmarks memory allocations for a full workflow
func BenchmarkFullWorkflowAllocs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate a typical usage pattern
		if glint.ColorSupport() {
			level := glint.ColorLevel()
			_ = level
		}
	}
}

// BenchmarkZeroAllocationsCheck verifies zero allocations for cached operations
func BenchmarkZeroAllocationsCheck(b *testing.B) {
	// Prime all caches
	glint.ColorSupport()
	glint.ColorLevel()
	core.SetEnvCache()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// These should all be zero allocations when cached
		glint.ColorSupport()
		glint.ColorLevel()
		core.GetEnvCache(core.EnvTerm)
		core.TerminalColorLevel()
	}
}
