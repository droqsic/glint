//go:build !windows
// +build !windows

package platform

// enableVirtualTerminal is a no-op implementation for non-Windows systems.
// Unix-like systems typically support ANSI escape sequences by default, so no additional configuration is required.
func EnableVirtualTerminal() {}
