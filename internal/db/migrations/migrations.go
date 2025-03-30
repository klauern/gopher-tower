package migrations

import "embed"

// Files contains the embedded migration SQL files
//
//go:embed *.sql
var Files embed.FS
