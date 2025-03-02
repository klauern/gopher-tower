// Package static provides embedded static files for the web UI
package static

import "embed"

//go:embed all:frontend
var Files embed.FS
