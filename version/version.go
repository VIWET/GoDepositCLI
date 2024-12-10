package version

import "fmt"

var (
	Major string
	Minor string
	Patch string
)

func Version() string {
	return fmt.Sprintf("v%s.%s.%s", Major, Minor, Patch)
}
