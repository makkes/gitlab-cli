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

const repo = "https://github.com/makkes/gitlab-cli"

func NewCommand() *cobra.Command {
	var dryRun *bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update GitLab CLI to latest version",
		Args:  cobra.MaximumNArgs(1),
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
				if *dryRun {
					fmt.Printf("A new version is available: %s\nSee %s for details\n", latestVersionString, repo+"/releases/"+latestVersionString)
					return nil
				}
				downloadURL := fmt.Sprintf("%s/releases/download/%s/gitlab_%s_%s_%s",
					repo, latestVersionString, latestVersionString, runtime.GOOS, runtime.GOARCH)
				fmt.Printf("Updating to %s\n", latestVersionString)
				resp, err := http.Get(downloadURL) // #nosec G107
				if err != nil {
					return fmt.Errorf("Could not download latest release: %w", err)
				}
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("Could not download latest release: received HTTP status %d", resp.StatusCode)
				}
				exec, err := os.Executable()
				if err != nil {
					return fmt.Errorf("Could not get current executable to update: %w", err)
				}
				dest := exec + ".new"
				stat, err := os.Stat(exec)
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
				err = os.Rename(dest, exec)
				if err != nil {
					return fmt.Errorf("Could not update to new version: %w", err)
				}
			} else if *dryRun {
				fmt.Printf("No update available, yet.\n")
				return nil
			}
			return nil
		},
	}

	dryRun = cmd.Flags().BoolP("dry-run", "d", false, "Only check if an update is available")

	return cmd
}
