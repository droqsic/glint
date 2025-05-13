package feature

import (
	"os"
	"sync"
)

const (
	EnvTerm      = "TERM"      // The terminal type
	EnvNoColor   = "NO_COLOR"  // Disable color support if set to "1" or "true"
	EnvColorTerm = "COLORTERM" // The color support type
)

var (
	EnvCache    sync.Map  // Caches commonly used environment variables
	envInitOnce sync.Once // Ensures the environment cache is initialized only once
)

// SetEnvCache initializes the environment cache
// This function is called automatically when needed, but can be called manually if needed
// It stores the values in the EnvCache map, thread-safe and will not overwrite existing values
func SetEnvCache() {
	envInitOnce.Do(initEnvCache)
}

// GetEnvCache retrieves a value from the environment cache
// It returns an empty string if the variable is not found, otherwise it returns the value as a string.
// This function is thread-safe and optimized for zero allocations.
func GetEnvCache(name string) string {
	envInitOnce.Do(initEnvCache)

	if val, ok := EnvCache.Load(name); ok && val != nil {
		return val.(string)
	}

	return ""
}

// initEnvCache initializes the environment cache
// It stores the values in the EnvCache map
// This function is called automatically when needed, but can be called manually if needed. It is thread-safe.
func initEnvCache() {
	termVal := os.Getenv(EnvTerm)
	colorTermVal := os.Getenv(EnvColorTerm)
	noColorVal := os.Getenv(EnvNoColor)

	EnvCache.Store(EnvTerm, termVal)
	EnvCache.Store(EnvColorTerm, colorTermVal)
	EnvCache.Store(EnvNoColor, noColorVal)
}

// ResetCache resets the environment cache for testing purposes. It is not meant to be used in production code.
func ResetCache() {
	envInitOnce = sync.Once{}

	EnvCache.Range(func(key, value interface{}) bool {
		EnvCache.Delete(key)
		return true
	})

	SetEnvCache()
}
