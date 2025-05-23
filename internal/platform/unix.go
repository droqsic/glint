//go:build !windows
// +build !windows

package platform

// EnableVirtualTerminal is a no-op implementation for non-Windows systems.
// Unix-like systems typically support ANSI escape sequences by default, so no additional configuration is required.
// It always returns false to indicate it's a no-op.
func EnableVirtualTerminal() bool {
	return false
}
