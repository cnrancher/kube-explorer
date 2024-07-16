package version

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	Version   = "dev"
	GitCommit = "HEAD"

	// K-EXPLORER
	releasePattern = regexp.MustCompile("^v[0-9]")
)

func FriendlyVersion() string {
	return fmt.Sprintf("%s (%s)", Version, GitCommit)
}

func IsRelease() bool {
	return !strings.Contains(Version, "dev") && releasePattern.MatchString(Version)
}
