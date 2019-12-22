package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/makkes/gitlab-cli/config"
	"github.com/makkes/gitlab-cli/semver"

	"github.com/makkes/gitlab-cli/versions"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update GitLab CLI to latest version",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			latestVersionString, err := versions.LatestVersion()
			if err != nil {
				return fmt.Errorf("Could not fetch latest version: %w", err)
			}
			latestVersion, err := semver.NewVersion(latestVersionString)
			if err != nil {
				return fmt.Errorf("Could not parse latest version '%s': %w", latestVersionString, err)
			}
			currentVersion, err := semver.NewVersion(config.Version)
			if err != nil {
				return fmt.Errorf("Could not parse current version '%s': %w", config.Version, err)
			}
			if currentVersion.Compare(*latestVersion) == -1 {
				downloadURL := fmt.Sprintf("https://github.com/makkes/gitlab-cli/releases/download/%s/gitlab_%s_%s_%s",
					latestVersionString, latestVersionString, runtime.GOOS, runtime.GOARCH)
				fmt.Printf("Updating to %s\n", latestVersionString)
				resp, err := http.Get(downloadURL) // #nosec G107
				if err != nil {
					return fmt.Errorf("Could not download latest release: %w", err)
				}
				dest := os.Args[0] + ".new"
				stat, err := os.Stat(os.Args[0])
				if err != nil {
					return fmt.Errorf("Could not stat binary for updating: %w", err)
				}
				binary, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, stat.Mode())
				if err != nil {
					return fmt.Errorf("Could not open binary for updating: %w", err)
				}
				defer binary.Close()
				_, err = io.Copy(binary, resp.Body)
				if err != nil {
					return fmt.Errorf("Could not download new version: %w", err)
				}
				err = os.Rename(dest, os.Args[0])
				if err != nil {
					return fmt.Errorf("Could not update to new version: %w", err)
				}
			}
			return nil
		},
	}
	return cmd
}
