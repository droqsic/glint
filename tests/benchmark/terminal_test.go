package benchmark

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
)

// BenchmarkIsTerminal benchmarks the IsTerminal function
func BenchmarkIsTerminal(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.IsTerminal(fd)
	}
}

// BenchmarkIsCygwinTerminal benchmarks the IsCygwinTerminal function
func BenchmarkIsCygwinTerminal(b *testing.B) {
	fd := os.Stdout.Fd()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glint.IsCygwinTerminal(fd)
	}
}
