package integration

import (
	"os"
	"testing"

	"github.com/droqsic/glint"
)

// TestIsTerminal tests the IsTerminal function
func TestIsTerminal(t *testing.T) {
	t.Run("Stdout", func(t *testing.T) {
		fd := os.Stdout.Fd()
		result := glint.IsTerminal(fd)

		// We can't make strong assertions about the result because it depends on
		// whether the test is running in a terminal or not
		t.Logf("IsTerminal(os.Stdout.Fd()) = %v", result)

		// If we're running in CI, the result is likely to be false
		// If we're running locally in a terminal, the result is likely to be true
	})

	// Test with a file
	t.Run("File", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "terminal_test")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fd := tmpFile.Fd()

		result := glint.IsTerminal(fd)
		if result {
			t.Errorf("Expected IsTerminal(%v) to be false, got true", fd)
		}
	})
}

// TestIsCygwinTerminal tests the IsCygwinTerminal function
func TestIsCygwinTerminal(t *testing.T) {
	t.Run("Stdout", func(t *testing.T) {
		fd := os.Stdout.Fd()
		result := glint.IsCygwinTerminal(fd)

		// We can't make strong assertions about the result because it depends on
		// whether the test is running in a Cygwin/MSYS2 terminal or not
		t.Logf("IsCygwinTerminal(os.Stdout.Fd()) = %v", result)

		// If we're not running in a Cygwin/MSYS2 terminal, the result should be false
		// If we are running in a Cygwin/MSYS2 terminal, the result should be true
	})

	// Test with a file
	t.Run("File", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "terminal_test")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fd := tmpFile.Fd()

		result := glint.IsCygwinTerminal(fd)
		if result {
			t.Errorf("Expected IsCygwinTerminal(%v) to be false, got true", fd)
		}
	})
}

// TestTerminalFunctions tests both IsTerminal and IsCygwinTerminal together
func TestTerminalFunctions(t *testing.T) {
	stdoutFd := os.Stdout.Fd()
	stderrFd := os.Stderr.Fd()
	stdinFd := os.Stdin.Fd()

	t.Logf("IsTerminal(os.Stdout.Fd()) = %v", glint.IsTerminal(stdoutFd))
	t.Logf("IsTerminal(os.Stderr.Fd()) = %v", glint.IsTerminal(stderrFd))
	t.Logf("IsTerminal(os.Stdin.Fd()) = %v", glint.IsTerminal(stdinFd))

	t.Logf("IsCygwinTerminal(os.Stdout.Fd()) = %v", glint.IsCygwinTerminal(stdoutFd))
	t.Logf("IsCygwinTerminal(os.Stderr.Fd()) = %v", glint.IsCygwinTerminal(stderrFd))
	t.Logf("IsCygwinTerminal(os.Stdin.Fd()) = %v", glint.IsCygwinTerminal(stdinFd))

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	rFd := r.Fd()
	wFd := w.Fd()

	if glint.IsTerminal(rFd) {
		t.Errorf("Expected IsTerminal(%v) to be false, got true", rFd)
	}
	if glint.IsTerminal(wFd) {
		t.Errorf("Expected IsTerminal(%v) to be false, got true", wFd)
	}

	if glint.IsCygwinTerminal(rFd) {
		t.Errorf("Expected IsCygwinTerminal(%v) to be false, got true", rFd)
	}
	if glint.IsCygwinTerminal(wFd) {
		t.Errorf("Expected IsCygwinTerminal(%v) to be false, got true", wFd)
	}
}
