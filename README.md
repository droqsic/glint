# Glint

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/droqsic/glint.svg)](https://pkg.go.dev/github.com/droqsic/glint)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Workflow](https://github.com/droqsic/glint/actions/workflows/go.yml/badge.svg)](https://github.com/droqsic/glint/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/droqsic/glint?nocache=1)](https://goreportcard.com/report/github.com/droqsic/glint)
[![Latest Release](https://img.shields.io/github/v/release/droqsic/glint)](https://github.com/droqsic/glint/releases)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)

</div>

**Glint** is a lightweight, cross-platform Go library for detecting and enabling terminal color support. It combines performance, accuracy, and simplicity — with benchmarks showing it's up to **1000x faster** than alternatives.

## Features

- ⚡ **High Performance**: Advanced caching makes repeated checks nearly instantaneous
- 🌐 **Cross-Platform**: Works on Windows, macOS, Linux, BSD, and more
- 🧠 **Zero Allocations**: Efficient design ensures no heap allocations
- 🔒 **Thread-Safe**: Safe for concurrent use from multiple goroutines
- 🧼 **Simple API**: Easy to use and integrate
- 📦 **Minimal Dependencies**: Only depends on the Go standard library, `x/sys`, and `probe`

## Installation

```bash
go get github.com/droqsic/glint
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/droqsic/glint"
)

func main() {
    if glint.IsColorSupported() {
        fmt.Println("Terminal supports colors")
    } else {
        fmt.Println("Terminal does not support colors")
    }

    fmt.Printf("Color support level: %s\n", glint.IsColorSupportedLevel())

    // Force-enable color support (use cautiously)
    glint.ForceColorSupport()
}
```

## How It Works

Glint determines terminal color support through:

- 🧾 **Environment Variables**: Inspects `TERM`, `COLORTERM`, `NO_COLOR`
- 🧪 **Terminal Detection**: Uses the `probe` library to check if output is a terminal
- 🪟 **Windows Support**: Enables virtual terminal sequences when necessary
- 🌈 **Color Levels**: Distinguishes between None, 16, 256, and True Color

All results are cached to ensure ultra-fast subsequent checks.

## Performance

Glint is engineered for speed. Here's what benchmarks reveal:

```
BenchmarkIsColorSupported-12                1000000000	         0.25 ns/op
BenchmarkIsColorSupportedLevel-12           1000000000	         0.89 ns/op
BenchmarkForceColorSupport-12               42567984	        30.15 ns/op
```

- 🔁 **Cached Checks**: Almost all operations complete in under 1 ns
- 🪄 **Zero Allocations**: No memory allocations for any operation
- 💯 **High Throughput**: Ideal for performance-critical CLI tools

## Color Support Levels

Glint can detect four levels of color support:

| Level       | Colors          | Example Terminals                               |
| ----------- | --------------- | ----------------------------------------------- |
| `LevelNone` | 0               | Non-interactive shells, logs                    |
| `Level16`   | 16 ANSI colors  | `xterm`, `vt100`, `screen`                      |
| `Level256`  | 256 colors      | `xterm-256color`, `screen-256color`             |
| `LevelTrue` | 16M true colors | Terminals with `COLORTERM=truecolor` or `24bit` |

## Thread Safety

Glint is built with concurrency in mind. It uses synchronization mechanisms such as `sync.Once` and `sync.Map` to manage its internal cache, allowing multiple goroutines to access color support checks safely and efficiently. The design is optimized for read-heavy workloads, ensuring high throughput and low latency even under concurrent access. This makes Glint well-suited for use in modern, parallelized Go applications.

## Contributing

We welcome contributions of all kinds! Bug fixes, new features, test improvements, and docs are all appreciated.

- Read the [Contributing Guide](docs/CONTRIBUTING.md) to get started
- Please follow the [Code of Conduct](docs/CODE_OF_CONDUCT.md)

## License

Glint is released under the MIT License. For the full license text, please see the [LICENSE](LICENSE) file.

## Acknowledgements

This project is inspired by the need for high-performance terminal color detection in real-world Go applications. Special thanks to:

- The Go team for their exceptional language and tooling
- The maintainers of x/sys for low-level system access
- The creators of [Probe](https://github.com/droqsic/probe) for the terminal detection library
- All contributors who help make Glint better
