package migrations

import "embed"

// FS contains every Goose migration so production images do not need a
// separately mounted migrations directory.
//
//go:embed *.sql
var FS embed.FS
