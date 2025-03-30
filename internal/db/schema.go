package db

import (
	"embed"
	_ "embed"
)

//go:embed schema.sql
var Schema string

//go:embed migrate/migrations/*.sql
var Migrations embed.FS
