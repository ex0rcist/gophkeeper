// Allows to embed list of migrate files to pass to migrator
// See https://pkg.go.dev/embed
package migrations

import "embed"

//go:embed "*.sql"
var Migrations embed.FS
