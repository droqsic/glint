package core

import (
	"os"
	"sync"
)

const (
	EnvTerm           = "TERM"                 // Terminal type (e.g., xterm-256color)
	EnvColorTerm      = "COLORTERM"            // Color support hint (e.g., truecolor, 24bit)
	EnvNoColor        = "NO_COLOR"             // Disables color output explicitly
	EnvForceColor     = "FORCE_COLOR"          // Forces color output regardless of other detection
	EnvTermProgram    = "TERM_PROGRAM"         // Terminal program (e.g., iTerm.app, Apple_Terminal)
	EnvTermProgramVer = "TERM_PROGRAM_VERSION" // Terminal program version
	EnvWTSession      = "WT_SESSION"           // Windows Terminal session flag
	EnvWTProfileID    = "WT_PROFILE_ID"        // Windows Terminal profile ID
	EnvANSICON        = "ANSICON"              // Indicates ANSI support in legacy Windows terminals
	EnvConEmuANSI     = "ConEmuANSI"           // ANSI support flag for ConEmu
	EnvCI             = "CI"                   // Continuous integration environment
	EnvSSHConnection  = "SSH_CONNECTION"       // Remote SSH session
	EnvWSLEnv         = "WSLENV"               // Present in WSL (Windows Subsystem for Linux)
	EnvTermuxVersion  = "TERMUX_VERSION"       // Termux shell on Android
	EnvCustomColor16  = "COLOR_16"             // Custom flag to force 16 color mode
	EnvCustomColor256 = "COLOR_256"            // Custom flag to force 256 color mode
	EnvCustomColor24  = "COLOR_24"             // Custom flag to force 24-bit truecolor mode
)

var (
	envCache map[string]string // envCache stores environment variable values to avoid repeated system calls
	envMutex sync.RWMutex      // envMutex protects concurrent access to the environment cache
	envInit  bool              // envInit tracks whether the cache has been initialized

	// knownKeys is the list of environment variables to cache
	knownKeys = []string{
		EnvTerm,
		EnvColorTerm,
		EnvNoColor,
		EnvForceColor,
		EnvTermProgram,
		EnvTermProgramVer,
		EnvWTSession,
		EnvWTProfileID,
		EnvANSICON,
		EnvConEmuANSI,
		EnvCI,
		EnvSSHConnection,
		EnvWSLEnv,
		EnvTermuxVersion,
		EnvCustomColor16,
		EnvCustomColor256,
		EnvCustomColor24,
	}
)

// init initializes the environment variable cache.
func init() {
	envCache = make(map[string]string, len(knownKeys))
}

// SetEnvCache populates the environment variable cache if it hasn't been initialized.
// This function is thread-safe and will only initialize the cache once.
func SetEnvCache() {
	if envInit {
		return
	}

	envMutex.Lock()
	defer envMutex.Unlock()

	if envInit {
		return
	}

	for _, key := range knownKeys {
		envCache[key] = os.Getenv(key)
	}

	envInit = true
}

// GetEnvCache retrieves an environment variable value from the cache.
// If the cache hasn't been initialized, it will initialize it first.
func GetEnvCache(name string) string {
	if !envInit {
		SetEnvCache()
	}

	envMutex.RLock()
	defer envMutex.RUnlock()

	return envCache[name]
}

// ClearCache clears the environment variable cache, forcing a refresh on the next call to GetEnvCache.
// This is useful when environment variables might have changed during program execution.
func ClearCache() {
	envMutex.Lock()
	defer envMutex.Unlock()

	for k := range envCache {
		delete(envCache, k)
	}

	envInit = false
}
