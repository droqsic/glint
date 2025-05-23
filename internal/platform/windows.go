//go:build windows
// +build windows

package platform

import (
	"os"
	"sync"

	"golang.org/x/sys/windows"
)

var (
	vtProcessingEnabled     bool      // Caches the result of EnableVirtualTerminal
	vtProcessingEnabledOnce sync.Once // Ensures that EnableVirtualTerminal is called only once
)

// EnableVirtualTerminal enables virtual terminal processing on Windows.
func EnableVirtualTerminal() bool {
	vtProcessingEnabledOnce.Do(vtProcessing)
	return vtProcessingEnabled
}

// This allows ANSI escape sequences to be processed by the Windows console.
// This function is optimized for zero allocations and maximum performance.
// It returns true if virtual terminal processing is enabled, false otherwise.
// The result is cached after the first call for subsequent calls to be extremely fast.
func vtProcessing() {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		vtProcessingEnabled = false
		return
	}

	stdout := windows.Handle(os.Stdout.Fd())

	var mode uint32
	if err := windows.GetConsoleMode(stdout, &mode); err != nil {
		vtProcessingEnabled = false
		return
	}

	if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 {
		vtProcessingEnabled = true
		return
	}

	err := windows.SetConsoleMode(stdout, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	vtProcessingEnabled = (err == nil)
}
