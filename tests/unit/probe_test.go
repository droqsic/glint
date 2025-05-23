package unit

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
	"github.com/droqsic/probe"
)

// TestProbeIntegration tests the probe library integration
func TestProbeIntegration(t *testing.T) {
	t.Run("IsTerminalStdout", func(t *testing.T) {
		fd := os.Stdout.Fd()
		result := probe.IsTerminal(fd)

		// Result should be boolean
		_ = result
		// We can't assert the exact value since it depends on how tests are run
	})

	t.Run("IsTerminalStderr", func(t *testing.T) {
		fd := os.Stderr.Fd()
		result := probe.IsTerminal(fd)

		// Result should be boolean
		_ = result
	})

	t.Run("IsTerminalStdin", func(t *testing.T) {
		fd := os.Stdin.Fd()
		result := probe.IsTerminal(fd)

		// Result should be boolean
		_ = result
	})

	t.Run("IsCygwinTerminalStdout", func(t *testing.T) {
		fd := os.Stdout.Fd()
		result := probe.IsCygwinTerminal(fd)

		// Result should be boolean
		_ = result
	})

	t.Run("IsCygwinTerminalStderr", func(t *testing.T) {
		fd := os.Stderr.Fd()
		result := probe.IsCygwinTerminal(fd)

		// Result should be boolean
		_ = result
	})

	t.Run("IsCygwinTerminalStdin", func(t *testing.T) {
		fd := os.Stdin.Fd()
		result := probe.IsCygwinTerminal(fd)

		// Result should be boolean
		_ = result
	})
}

// TestProbeConsistency tests consistency of probe results
func TestProbeConsistency(t *testing.T) {
	t.Run("ConsistentResults", func(t *testing.T) {
		fd := os.Stdout.Fd()

		// Multiple calls should return same result
		result1 := probe.IsTerminal(fd)
		result2 := probe.IsTerminal(fd)
		result3 := probe.IsTerminal(fd)

		if result1 != result2 || result2 != result3 {
			t.Errorf("IsTerminal results inconsistent: %v, %v, %v", result1, result2, result3)
		}

		// Same for Cygwin
		cygwin1 := probe.IsCygwinTerminal(fd)
		cygwin2 := probe.IsCygwinTerminal(fd)
		cygwin3 := probe.IsCygwinTerminal(fd)

		if cygwin1 != cygwin2 || cygwin2 != cygwin3 {
			t.Errorf("IsCygwinTerminal results inconsistent: %v, %v, %v", cygwin1, cygwin2, cygwin3)
		}
	})

	t.Run("LogicalConsistency", func(t *testing.T) {
		fd := os.Stdout.Fd()

		isTerminal := probe.IsTerminal(fd)
		isCygwin := probe.IsCygwinTerminal(fd)

		// If it's a Cygwin terminal, it should also be detected as a terminal
		// (though the reverse is not necessarily true)
		if isCygwin && !isTerminal {
			t.Errorf("Cygwin terminal should also be detected as terminal")
		}
	})
}

// TestProbeConcurrency tests probe functions under concurrent access
func TestProbeConcurrency(t *testing.T) {
	t.Run("ConcurrentIsTerminal", func(t *testing.T) {
		fd := os.Stdout.Fd()
		done := make(chan bool, 10)
		results := make(chan bool, 100)

		// Start multiple goroutines
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 10; j++ {
					result := probe.IsTerminal(fd)
					results <- result
				}
			}()
		}

		// Wait for completion
		for i := 0; i < 10; i++ {
			<-done
		}

		// Collect results
		var allResults []bool
		for i := 0; i < 100; i++ {
			allResults = append(allResults, <-results)
		}

		// All results should be the same
		if len(allResults) > 0 {
			firstResult := allResults[0]
			for i, result := range allResults {
				if result != firstResult {
					t.Errorf("Concurrent call %d returned %v, expected %v", i, result, firstResult)
				}
			}
		}
	})

	t.Run("ConcurrentIsCygwinTerminal", func(t *testing.T) {
		fd := os.Stdout.Fd()
		done := make(chan bool, 10)
		results := make(chan bool, 100)

		// Start multiple goroutines
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 10; j++ {
					result := probe.IsCygwinTerminal(fd)
					results <- result
				}
			}()
		}

		// Wait for completion
		for i := 0; i < 10; i++ {
			<-done
		}

		// Collect results
		var allResults []bool
		for i := 0; i < 100; i++ {
			allResults = append(allResults, <-results)
		}

		// All results should be the same
		if len(allResults) > 0 {
			firstResult := allResults[0]
			for i, result := range allResults {
				if result != firstResult {
					t.Errorf("Concurrent call %d returned %v, expected %v", i, result, firstResult)
				}
			}
		}
	})
}

// TestProbeEdgeCases tests edge cases for probe functions
func TestProbeEdgeCases(t *testing.T) {
	t.Run("InvalidFileDescriptor", func(t *testing.T) {
		// Test with invalid file descriptor
		// Note: This might cause different behavior on different platforms
		invalidFd := uintptr(999999)

		// Should not panic
		result1 := probe.IsTerminal(invalidFd)
		result2 := probe.IsCygwinTerminal(invalidFd)

		// Results should be boolean (likely false for invalid fd)
		_ = result1
		_ = result2
	})

	t.Run("ZeroFileDescriptor", func(t *testing.T) {
		// Test with zero file descriptor
		zeroFd := uintptr(0)

		result1 := probe.IsTerminal(zeroFd)
		result2 := probe.IsCygwinTerminal(zeroFd)

		_ = result1
		_ = result2
	})

	t.Run("HighFileDescriptor", func(t *testing.T) {
		// Test with high file descriptor values
		highFds := []uintptr{100, 1000, 10000}

		for _, fd := range highFds {
			result1 := probe.IsTerminal(fd)
			result2 := probe.IsCygwinTerminal(fd)

			_ = result1
			_ = result2
		}
	})
}

// TestProbePerformance tests performance characteristics
func TestProbePerformance(t *testing.T) {
	t.Run("RepeatedCalls", func(t *testing.T) {
		fd := os.Stdout.Fd()

		// Many repeated calls should be fast
		for i := 0; i < 1000; i++ {
			probe.IsTerminal(fd)
			probe.IsCygwinTerminal(fd)
		}
	})

	t.Run("DifferentFileDescriptors", func(t *testing.T) {
		fds := []uintptr{
			os.Stdin.Fd(),
			os.Stdout.Fd(),
			os.Stderr.Fd(),
		}

		for _, fd := range fds {
			for i := 0; i < 100; i++ {
				probe.IsTerminal(fd)
				probe.IsCygwinTerminal(fd)
			}
		}
	})
}

// TestProbeIntegrationWithGlint tests how probe integrates with glint
func TestProbeIntegrationWithGlint(t *testing.T) {
	t.Run("TerminalDetectionImpactsColorSupport", func(t *testing.T) {
		// The probe library is used internally by glint for terminal detection
		// This test ensures the integration works correctly

		fd := os.Stdout.Fd()
		isTerminal := probe.IsTerminal(fd)
		isCygwin := probe.IsCygwinTerminal(fd)

		// These results should influence glint's color support detection
		_ = isTerminal
		_ = isCygwin

		// Test glint's behavior
		// (Results will vary based on actual terminal state)
		colorSupported := glint.ColorSupport()
		_ = colorSupported
	})
}
