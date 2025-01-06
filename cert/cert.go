package cert

import "embed"

//go:embed "*.pem"
var Cert embed.FS
