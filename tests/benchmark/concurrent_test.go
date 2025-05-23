package benchmark

import (
	"sync"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/glint/internal/core"
)

// BenchmarkColorSupportConcurrent benchmarks ColorSupport under concurrent access
func BenchmarkColorSupportConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			glint.ColorSupport()
		}
	})
}

// BenchmarkColorLevelConcurrent benchmarks ColorLevel under concurrent access
func BenchmarkColorLevelConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			glint.ColorLevel()
		}
	})
}

// BenchmarkGetEnvCacheConcurrent benchmarks GetEnvCache under concurrent access
func BenchmarkGetEnvCacheConcurrent(b *testing.B) {
	// Prime the cache
	core.SetEnvCache()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			core.GetEnvCache(core.EnvTerm)
		}
	})
}

// BenchmarkTerminalColorLevelConcurrent benchmarks TerminalColorLevel under concurrent access
func BenchmarkTerminalColorLevelConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			core.TerminalColorLevel()
		}
	})
}

// BenchmarkForceColorConcurrent benchmarks ForceColor under concurrent access
func BenchmarkForceColorConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			glint.ForceColor(true)
			glint.ForceColor(false)
		}
	})
}

// BenchmarkResetColorConcurrent benchmarks ResetColor under concurrent access
// Note: ResetColor is not designed for concurrent use, so this tests sequential calls
func BenchmarkResetColorConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		glint.ResetColor()
	}
}

// BenchmarkMixedOperationsConcurrent benchmarks mixed operations under concurrent access
func BenchmarkMixedOperationsConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Mix of read and write operations
			glint.ColorSupport()
			glint.ColorLevel()
			core.GetEnvCache(core.EnvTerm)
			core.TerminalColorLevel()
		}
	})
}

// BenchmarkCacheOperationsConcurrent benchmarks cache operations under concurrent access
func BenchmarkCacheOperationsConcurrent(b *testing.B) {
	keys := []string{
		core.EnvTerm,
		core.EnvColorTerm,
		core.EnvNoColor,
		core.EnvForceColor,
		core.EnvTermProgram,
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, key := range keys {
				core.GetEnvCache(key)
			}
		}
	})
}

// BenchmarkHighContentionScenario benchmarks high contention scenarios
func BenchmarkHighContentionScenario(b *testing.B) {
	const numGoroutines = 100
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go func() {
				defer wg.Done()
				glint.ColorSupport()
				glint.ColorLevel()
				core.GetEnvCache(core.EnvTerm)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkCacheInitializationConcurrent benchmarks concurrent cache initialization
// Note: ClearCache is not designed for concurrent use, so this tests sequential calls
func BenchmarkCacheInitializationConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		core.ClearCache()
		core.SetEnvCache()
	}
}

// BenchmarkConcurrentCacheAccess benchmarks concurrent access to different cache keys
func BenchmarkConcurrentCacheAccess(b *testing.B) {
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
	}

	b.RunParallel(func(pb *testing.PB) {
		keyIndex := 0
		for pb.Next() {
			core.GetEnvCache(keys[keyIndex%len(keys)])
			keyIndex++
		}
	})
}

// BenchmarkSyncOncePerformance benchmarks the performance of sync.Once patterns
func BenchmarkSyncOncePerformance(b *testing.B) {
	var once sync.Once
	var result bool

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.Do(func() {
				result = true
			})
			_ = result
		}
	})
}

// BenchmarkSafeReadOperationsConcurrent benchmarks safe read operations under concurrent access
func BenchmarkSafeReadOperationsConcurrent(b *testing.B) {
	// Prime the caches first
	glint.ColorSupport()
	glint.ColorLevel()
	core.SetEnvCache()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Only safe read operations
			glint.ColorSupport()
			glint.ColorLevel()
			core.GetEnvCache(core.EnvTerm)
			core.GetEnvCache(core.EnvColorTerm)
		}
	})
}

// BenchmarkForceColorToggleConcurrent benchmarks ForceColor toggle operations
func BenchmarkForceColorToggleConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Safer than ResetColor - just toggle between true/false
			glint.ForceColor(true)
			glint.ColorSupport()
			glint.ForceColor(false)
			glint.ColorSupport()
		}
	})
}
