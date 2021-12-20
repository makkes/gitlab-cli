package versions

import (
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/mmcdole/gofeed"
)

func LatestVersion(repo string, currentVersion semver.Version, upgradeMajor, includePreReleases bool) (string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(repo + "/releases.atom")
	if err != nil {
		return "", err
	}
	if len(feed.Items) == 0 {
		return "", fmt.Errorf("no entry in releases feed")
	}

	for _, item := range feed.Items {
		IDElements := strings.Split(item.GUID, "/")
		if len(IDElements) < 3 {
			continue // entry is malformed, skip it
		}
		v, err := semver.ParseTolerant(IDElements[2])
		if err != nil {
			continue // this is not a semver, skip it
		}
		if v.Pre != nil && !includePreReleases {
			continue // this is a pre-release, skip it
		}
		if v.Major > currentVersion.Major && !upgradeMajor {
			continue // this is a new major version, skip ip
		}
		return fmt.Sprintf("v%s", v.String()), nil
	}
	return "", fmt.Errorf("no parseable version found in releases feed")
}
