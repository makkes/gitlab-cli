package update

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"

	"github.com/makkes/gitlab-cli/config"
	"github.com/makkes/gitlab-cli/versions"
)

var repo = "https://github.com/makkes/gitlab-cli"

func updateCommand(dryRun bool, includePreReleases bool, upgradeMajor bool, out io.Writer,
	getExecutable func() (string, error)) error {
	currentVersion, err := semver.ParseTolerant(config.Version)
	if err != nil {
		return fmt.Errorf("could not parse current version '%s': %w", config.Version, err)
	}
	latestVersionString, err := versions.LatestVersion(repo, currentVersion, upgradeMajor, includePreReleases)
	if err != nil {
		return fmt.Errorf("could not fetch latest version: %w", err)
	}
	latestVersion, err := semver.ParseTolerant(latestVersionString)
	if err != nil {
		return fmt.Errorf("could not parse latest version '%s': %w", latestVersionString, err)
	}
	if currentVersion.Compare(latestVersion) == -1 {
		if dryRun {
			fmt.Fprintf(out, "A new version is available: %s\nSee %s for details\n",
				latestVersionString, repo+"/releases/"+latestVersionString)
			return nil
		}
		downloadURL := fmt.Sprintf("%s/releases/download/%s/gitlab_%s_%s_%s.tar.gz",
			repo, latestVersionString, latestVersionString, runtime.GOOS, runtime.GOARCH)
		fmt.Fprintf(out, "Updating to %s\n", latestVersionString)
		resp, err := http.Get(downloadURL) // #nosec G107
		if err != nil {
			return fmt.Errorf("could not download latest release from %s: %w", downloadURL, err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("could not download latest release from %s: received HTTP status %d", downloadURL, resp.StatusCode)
		}

		exec, err := getExecutable()
		if err != nil {
			return fmt.Errorf("could not get current executable to update: %w", err)
		}
		dest := exec + ".new"
		stat, err := os.Stat(exec)
		if err != nil {
			return fmt.Errorf("could not stat binary for updating: %w", err)
		}
		binary, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, stat.Mode())
		if err != nil {
			return fmt.Errorf("could not open binary for updating: %w", err)
		}
		defer func() {
			os.Remove(dest)
		}()

		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed creating gzip reader: %w", err)
		}
		defer gzReader.Close()

		tarReader := tar.NewReader(gzReader)
		if _, err := tarReader.Next(); err != nil {
			return fmt.Errorf("failed advancing in tar archive: %w", err)
		}

		_, err = io.Copy(binary, tarReader)
		binary.Close()
		if err != nil {
			return fmt.Errorf("could not download new version: %w", err)
		}
		err = os.Rename(dest, exec)
		if err != nil {
			return fmt.Errorf("could not update to new version: %w", err)
		}
	} else {
		fmt.Fprintf(out, "You're already on the latest version %s.\n", config.Version)
		return nil
	}
	return nil
}

func NewCommand() *cobra.Command {
	var dryRun *bool
	var upgradeMajor *bool
	var includePre *bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update GitLab CLI to latest version",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateCommand(*dryRun, *includePre, *upgradeMajor, os.Stdout, os.Executable)
		},
	}

	dryRun = cmd.Flags().BoolP("dry-run", "n", false, "Only check if an update is available")
	includePre = cmd.Flags().BoolP("pre", "", false, "Upgrade to next pre-release version, if available")
	upgradeMajor = cmd.Flags().BoolP("major", "", false, "Upgrade major version, if available")

	return cmd
}
