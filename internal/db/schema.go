package db

import (
	"embed"
	_ "embed"
)

//go:embed migrations/000001_init_schema.up.sql
var Schema string

//go:embed migrations/*.sql
var Migrations embed.FS
