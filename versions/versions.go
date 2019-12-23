package versions

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

func LatestVersion(repo string) (string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(repo + "/releases.atom")
	if err != nil {
		return "", err
	}
	latestEntry := feed.Items[0]
	if latestEntry == nil {
		return "", fmt.Errorf("no entry in releases feed")
	}
	IDElements := strings.Split(latestEntry.GUID, "/")
	if len(IDElements) < 3 {
		return "", fmt.Errorf("ID entry has unexpected format '%s'", latestEntry.GUID)
	}
	return IDElements[2], nil
}
