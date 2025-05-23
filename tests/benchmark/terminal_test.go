package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/probe"
)

// BenchmarkIsTerminal benchmarks the probe.IsTerminal function
func BenchmarkIsTerminal(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fd)
	}
}

// BenchmarkIsCygwinTerminal benchmarks the probe.IsCygwinTerminal function
func BenchmarkIsCygwinTerminal(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(fd)
	}
}

// BenchmarkIsTerminalStdout benchmarks IsTerminal specifically for stdout
func BenchmarkIsTerminalStdout(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsTerminal(os.Stdout.Fd())
	}
}

// BenchmarkIsTerminalStderr benchmarks IsTerminal specifically for stderr
func BenchmarkIsTerminalStderr(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsTerminal(os.Stderr.Fd())
	}
}

// BenchmarkIsTerminalStdin benchmarks IsTerminal specifically for stdin
func BenchmarkIsTerminalStdin(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsTerminal(os.Stdin.Fd())
	}
}

// BenchmarkIsCygwinTerminalStdout benchmarks IsCygwinTerminal specifically for stdout
func BenchmarkIsCygwinTerminalStdout(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(os.Stdout.Fd())
	}
}

// BenchmarkIsCygwinTerminalStderr benchmarks IsCygwinTerminal specifically for stderr
func BenchmarkIsCygwinTerminalStderr(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(os.Stderr.Fd())
	}
}

// BenchmarkTerminalDetectionCombined benchmarks combined terminal detection
func BenchmarkTerminalDetectionCombined(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isTerminal := probe.IsTerminal(fd)
		isCygwin := probe.IsCygwinTerminal(fd)
		_ = isTerminal || isCygwin
	}
}

// BenchmarkTerminalDetectionPattern benchmarks the pattern used in glint
func BenchmarkTerminalDetectionPattern(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This mimics the pattern used in glint.ColorSupport()
		if !probe.IsTerminal(fd) && !probe.IsCygwinTerminal(fd) {
			// Not a terminal
			continue
		}
		// Is a terminal
	}
}

// BenchmarkFileDescriptorOperations benchmarks file descriptor operations
func BenchmarkFileDescriptorOperations(b *testing.B) {
	b.Run("StdoutFd", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = os.Stdout.Fd()
		}
	})

	b.Run("StderrFd", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = os.Stderr.Fd()
		}
	})

	b.Run("StdinFd", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = os.Stdin.Fd()
		}
	})
}

// BenchmarkTerminalDetectionAllocs benchmarks memory allocations for terminal detection
func BenchmarkTerminalDetectionAllocs(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsTerminal(fd)
	}
}

// BenchmarkCygwinTerminalDetectionAllocs benchmarks memory allocations for Cygwin terminal detection
func BenchmarkCygwinTerminalDetectionAllocs(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		probe.IsCygwinTerminal(fd)
	}
}

// BenchmarkTerminalDetectionConcurrent benchmarks terminal detection under concurrent access
func BenchmarkTerminalDetectionConcurrent(b *testing.B) {
	fd := os.Stdout.Fd()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			probe.IsTerminal(fd)
		}
	})
}

// BenchmarkCygwinTerminalDetectionConcurrent benchmarks Cygwin terminal detection under concurrent access
func BenchmarkCygwinTerminalDetectionConcurrent(b *testing.B) {
	fd := os.Stdout.Fd()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			probe.IsCygwinTerminal(fd)
		}
	})
}

// BenchmarkTerminalDetectionMixed benchmarks mixed terminal detection operations
func BenchmarkTerminalDetectionMixed(b *testing.B) {
	fds := []uintptr{
		os.Stdout.Fd(),
		os.Stderr.Fd(),
		os.Stdin.Fd(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, fd := range fds {
			probe.IsTerminal(fd)
			probe.IsCygwinTerminal(fd)
		}
	}
}
