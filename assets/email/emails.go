package email

import "embed"

//go:embed *.html emails.json
var FS embed.FS
