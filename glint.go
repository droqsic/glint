package glint

import (
	"os"
	"runtime"
	"sync"

	"github.com/droqsic/glint/internal/core"
	"github.com/droqsic/glint/internal/platform"
	"github.com/droqsic/probe"
)

var (
	colorSupport           bool       // colorSupport stores the detected terminal color support status
	colorSupportOnce       sync.Once  // colorSupportOnce ensures color support detection runs only once
	colorLevel             core.Level // colorLevel stores the detected terminal color level
	colorLevelOnce         sync.Once  // colorLevelOnce ensures color level detection runs only once
	forceColorSupport      *bool      // forceColorSupport allows overriding automatic color support detection, nil means automatic detection, non-nil means forced value
	forceColorSupportMutex sync.Mutex // forceColorSupportMutex protects concurrent access to forceColorSupport
)

// ColorSupport determines whether the current terminal supports color output.
// It checks if the output is a terminal and if the terminal supports color.
// The result is cached after the first call for performance. This function is thread-safe.
func ColorSupport() bool {
	forceColorSupportMutex.Lock()

	if forceColorSupport != nil {
		defer forceColorSupportMutex.Unlock()
		return *forceColorSupport
	}

	forceColorSupportMutex.Unlock()

	colorSupportOnce.Do(func() {
		if !probe.IsTerminal(os.Stdout.Fd()) && !probe.IsCygwinTerminal(os.Stdout.Fd()) {
			colorSupport = false
			return
		}

		if core.TerminalColorLevel() == core.LevelNone {
			colorSupport = false
			return
		}

		colorSupport = true
	})
	return colorSupport
}

// ColorLevel determines the color support level of the current terminal.
// It checks if the output is a terminal and what level of color it supports.
// The result is cached after the first call for performance. This function is thread-safe.
func ColorLevel() core.Level {
	if ColorSupport() {
		colorLevelOnce.Do(func() {
			if !probe.IsTerminal(os.Stdout.Fd()) && !probe.IsCygwinTerminal(os.Stdout.Fd()) {
				colorLevel = core.LevelNone
				return
			}
			colorLevel = core.TerminalColorLevel()
		})
	}
	return colorLevel
}

// ForceColor overrides automatic color support detection with a fixed value.
// This is useful for applications that want to explicitly enable or disable color regardless of terminal capabilities.
// However, it still respects the NO_COLOR environment variable - if NO_COLOR is set, colors will be disabled regardless.
func ForceColor(value bool) {
	forceColorSupportMutex.Lock()
	defer forceColorSupportMutex.Unlock()

	if value && core.GetEnvCache(core.EnvNoColor) != "" {
		value = false
	}

	if value && runtime.GOOS == "windows" {
		vtEnabled := platform.EnableVirtualTerminal()
		if !vtEnabled {
			colorLevelOnce = sync.Once{}
			colorLevelOnce.Do(func() {
				colorLevel = core.Level16
			})
		}
	}

	forceColorSupport = &value
	colorSupportOnce = sync.Once{}
}

// ResetColor resets color support detection to automatic mode, clearing any previously forced settings.
// This allows the system to detect terminal capabilities again, and also clears any previously forced settings.
func ResetColor() {
	forceColorSupportMutex.Lock()
	defer forceColorSupportMutex.Unlock()

	forceColorSupport = nil
	colorSupportOnce = sync.Once{}
	colorLevelOnce = sync.Once{}

	colorSupport = false
	colorLevel = core.LevelNone
}
