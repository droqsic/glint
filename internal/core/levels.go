package core

type Level int8 // Level represents the color support capability of a terminal.

const (
	LevelNone Level = iota // LevelNone indicates no color support
	Level16                // Level16 indicates basic ANSI color support (16 colors)
	Level256               // Level256 indicates extended ANSI color support (256 colors)
	LevelTrue              // LevelTrue indicates 24-bit RGB color support (TrueColor)
)

// String returns a human-readable description of the color support level.
func (l Level) String() string {
	switch l {
	case LevelNone:
		return "No color support detected"
	case Level16:
		return "ANSI color support (16 colors)"
	case Level256:
		return "ANSI extended color support (256 colors)"
	case LevelTrue:
		return "TrueColor support (24-bit RGB)"
	default:
		return "Unknown terminal color capability"
	}
}

// terminalColorLevel determines the color support level of the terminal based on environment variables and terminal type.
// This is an internal function used by DetectColorLevel.
func TerminalColorLevel() Level {
	SetEnvCache()

	// NO_COLOR environment variable takes precedence over everything else
	if GetEnvCache(EnvNoColor) != "" {
		return LevelNone
	}

	// Check for explicit color forcing environment variables
	if GetEnvCache(EnvForceColor) != "" {
		return LevelTrue
	}

	// Check for custom color level environment variables
	if GetEnvCache(EnvCustomColor24) != "" {
		return LevelTrue
	}
	if GetEnvCache(EnvCustomColor256) != "" {
		return Level256
	}
	if GetEnvCache(EnvCustomColor16) != "" {
		return Level16
	}

	// Check COLORTERM for truecolor or 256 color support
	switch GetEnvCache(EnvColorTerm) {
	case "truecolor", "24bit":
		return LevelTrue
	case "256color":
		return Level256
	}

	// Check TERM for color support information
	switch GetEnvCache(EnvTerm) {
	case "xterm-256color", "screen-256color", "tmux-256color", "rxvt-256color":
		return Level256
	case "xterm", "screen", "tmux", "rxvt":
		return Level16
	case "dumb":
		return LevelNone
	}

	// Check for specific terminal environments
	if GetEnvCache(EnvWTSession) != "" || GetEnvCache(EnvWTProfileID) != "" {
		return LevelTrue
	}

	if GetEnvCache(EnvANSICON) != "" || GetEnvCache(EnvConEmuANSI) == "ON" {
		return Level256
	}

	if GetEnvCache(EnvTermProgram) == "iTerm.app" {
		return LevelTrue
	}

	// CI environments typically support at least basic colors
	if GetEnvCache(EnvCI) != "" {
		return Level16
	}

	// Termux on Android supports 256 colors
	if GetEnvCache(EnvTermuxVersion) != "" {
		return Level256
	}

	// WSL typically supports 256 colors
	if GetEnvCache(EnvWSLEnv) != "" {
		return Level256
	}

	// SSH connections typically support 256 colors
	if GetEnvCache(EnvSSHConnection) != "" {
		return Level256
	}

	// Default to basic 16 colors if we can't determine anything more specific
	return Level16
}
