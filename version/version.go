package version

var GitVersion string

func Version() string {
	return GitVersion
}
