package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/blang/semver/v4"

	"github.com/makkes/gitlab-cli/config"

	"github.com/makkes/gitlab-cli/versions"
	"github.com/spf13/cobra"
)

var repo = "https://github.com/makkes/gitlab-cli"

func updateCommand(dryRun bool, includePreReleases bool, out io.Writer, getExecutable func() (string, error)) error {
	latestVersionString, err := versions.LatestVersion(repo, includePreReleases)
	if err != nil {
		return fmt.Errorf("Could not fetch latest version: %w", err)
	}
	latestVersion, err := semver.ParseTolerant(latestVersionString)
	if err != nil {
		return fmt.Errorf("Could not parse latest version '%s': %w", latestVersionString, err)
	}
	currentVersion, err := semver.ParseTolerant(config.Version)
	if err != nil {
		return fmt.Errorf("Could not parse current version '%s': %w", config.Version, err)
	}
	if currentVersion.Compare(latestVersion) == -1 {
		if dryRun {
			fmt.Fprintf(out, "A new version is available: %s\nSee %s for details\n", latestVersionString, repo+"/releases/"+latestVersionString)
			return nil
		}
		downloadURL := fmt.Sprintf("%s/releases/download/%s/gitlab_%s_%s_%s",
			repo, latestVersionString, latestVersionString, runtime.GOOS, runtime.GOARCH)
		fmt.Fprintf(out, "Updating to %s\n", latestVersionString)
		resp, err := http.Get(downloadURL) // #nosec G107
		if err != nil {
			return fmt.Errorf("Could not download latest release: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Could not download latest release: received HTTP status %d", resp.StatusCode)
		}
		exec, err := getExecutable()
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
		_, err = io.Copy(binary, resp.Body)
		binary.Close()
		if err != nil {
			return fmt.Errorf("Could not download new version: %w", err)
		}
		err = os.Rename(dest, exec)
		if err != nil {
			return fmt.Errorf("Could not update to new version: %w", err)
		}
	} else {
		fmt.Fprintf(out, "You're already on the latest version %s.\n", config.Version)
		return nil
	}
	return nil
}

func NewCommand() *cobra.Command {
	var dryRun *bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update GitLab CLI to latest version",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateCommand(*dryRun, false, os.Stdout, os.Executable)
		},
	}

	dryRun = cmd.Flags().BoolP("dry-run", "d", false, "Only check if an update is available")

	return cmd
}
