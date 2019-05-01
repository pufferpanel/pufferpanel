package version

import "fmt"

var (
	Hash    = "unknown"
	Version = "nightly"
	Display string
)

func init() {
	Display = fmt.Sprintf("PufferPanel %s (%s)", Version, Hash)
}
