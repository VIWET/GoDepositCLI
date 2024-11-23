package version

import "fmt"

const (
	Major = 0
	Minor = 0
	Patch = 1
)

func Version() string {
	return fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
}
