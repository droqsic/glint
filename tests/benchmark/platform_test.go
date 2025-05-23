package benchmark

import (
	"runtime"
	"testing"

	"github.com/droqsic/glint/internal/platform"
)

// BenchmarkEnableVirtualTerminal benchmarks the EnableVirtualTerminal function
func BenchmarkEnableVirtualTerminal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		platform.EnableVirtualTerminal()
	}
}

// BenchmarkEnableVirtualTerminalWindows benchmarks EnableVirtualTerminal specifically on Windows
func BenchmarkEnableVirtualTerminalWindows(b *testing.B) {
	if runtime.GOOS != "windows" {
		b.Skip("Skipping Windows-specific benchmark on non-Windows platform")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		platform.EnableVirtualTerminal()
	}
}

// BenchmarkEnableVirtualTerminalUnix benchmarks EnableVirtualTerminal on Unix-like systems
func BenchmarkEnableVirtualTerminalUnix(b *testing.B) {
	if runtime.GOOS == "windows" {
		b.Skip("Skipping Unix-specific benchmark on Windows platform")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		platform.EnableVirtualTerminal()
	}
}

// BenchmarkEnableVirtualTerminalCached benchmarks the cached behavior of EnableVirtualTerminal
func BenchmarkEnableVirtualTerminalCached(b *testing.B) {
	// Prime the cache by calling it once
	platform.EnableVirtualTerminal()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		platform.EnableVirtualTerminal()
	}
}

// BenchmarkEnableVirtualTerminalConcurrent benchmarks EnableVirtualTerminal under concurrent access
func BenchmarkEnableVirtualTerminalConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			platform.EnableVirtualTerminal()
		}
	})
}

// BenchmarkPlatformSpecificOperations benchmarks platform-specific operations
func BenchmarkPlatformSpecificOperations(b *testing.B) {
	b.Run("EnableVT", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			platform.EnableVirtualTerminal()
		}
	})

	b.Run("EnableVTWithGOOSCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if runtime.GOOS == "windows" {
				platform.EnableVirtualTerminal()
			}
		}
	})
}

// BenchmarkRuntimeGOOSCheck benchmarks the cost of runtime.GOOS checks
func BenchmarkRuntimeGOOSCheck(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runtime.GOOS == "windows"
	}
}

// BenchmarkPlatformDetection benchmarks platform detection patterns
func BenchmarkPlatformDetection(b *testing.B) {
	b.Run("DirectGOOSCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			isWindows := runtime.GOOS == "windows"
			_ = isWindows
		}
	})

	b.Run("SwitchGOOS", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			switch runtime.GOOS {
			case "windows":
				// Windows-specific logic
			case "darwin", "linux", "freebsd", "openbsd", "netbsd":
				// Unix-like logic
			default:
				// Other platforms
			}
		}
	})
}
