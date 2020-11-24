package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/makkes/gitlab-cli/config"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	tests := map[string]struct {
		feed           string
		currentVersion string
		updatedVersion string
		updatedBinary  []byte
		dryRun         bool
		upgradeMajor   bool
		pre            bool
		repo           string
		out            *regexp.Regexp
		outErr         bool
	}{
		"update happy path": {
			currentVersion: "3.6.2-55-ghf448b",
			updatedVersion: "3.6.3",
			updatedBinary:  []byte("this is the updated gitlab binary"),
			out:            regexp.MustCompile(`^Updating to v3.6.3\n`),
		},
		"no update happy path": {
			currentVersion: "3.6.3",
			updatedBinary:  []byte{},
			out:            regexp.MustCompile(`^You're already on the latest version `),
		},
		"dry-run no update happy path": {
			currentVersion: "3.6.3",
			updatedBinary:  []byte{},
			dryRun:         true,
			out:            regexp.MustCompile(`^You're already on the latest version `),
		},
		"dry-run update happy path": {
			currentVersion: "3.6.2-55-ghf448b",
			updatedVersion: "3.6.3",
			updatedBinary:  []byte{},
			dryRun:         true,
			out:            regexp.MustCompile(`^A new version is available: v3.6.3\nSee http://127.0.0.1:\d+/releases/v3.6.3 for details\n`),
		},
		"unreachable repo": {
			repo:           "http://doesntexist",
			currentVersion: "1.3.0",
			outErr:         true,
			out:            regexp.MustCompile(`^Could not fetch latest version: `),
			updatedBinary:  []byte{},
		},
		"dry-run update to pre-release": {
			currentVersion: "3.6.3",
			updatedVersion: "3.8.0-beta.1",
			updatedBinary:  []byte{},
			dryRun:         true,
			pre:            true,
			out:            regexp.MustCompile(`^A new version is available: v3.8.0-beta.1\nSee http://127.0.0.1:\d+/releases/v3.8.0-beta.1 for details\n`),
		},
		"update to pre-release": {
			currentVersion: "3.6.3",
			updatedVersion: "3.8.0-beta.1",
			updatedBinary:  []byte("the updated binary"),
			dryRun:         false,
			pre:            true,
			out:            regexp.MustCompile(`^Updating to v3.8.0-beta.1\n`),
		},
		"no update of major version": {
			feed:           "releases_test_major_upgrade.atom",
			currentVersion: "3.6.3",
			updatedVersion: "4.0.0",
			updatedBinary:  []byte{},
			dryRun:         false,
			pre:            false,
			out:            regexp.MustCompile(`^You're already on the latest version `),
		},
		"update of major version": {
			feed:           "releases_test_major_upgrade.atom",
			currentVersion: "3.6.3",
			updatedVersion: "4.0.0",
			updatedBinary:  []byte("4.0.0"),
			dryRun:         false,
			pre:            false,
			upgradeMajor:   true,
			out:            regexp.MustCompile(`^Updating to v4.0.0\n`),
		},
		"update of major pre-release version": {
			feed:           "releases_test_major_upgrade_pre.atom",
			currentVersion: "3.6.3",
			updatedVersion: "4.0.0-beta.1",
			updatedBinary:  []byte("4.0.0 Beta.1"),
			dryRun:         false,
			pre:            true,
			upgradeMajor:   true,
			out:            regexp.MustCompile(`^Updating to v4.0.0-beta.1\n`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(tt *testing.T) {
			feed := "releases_test.atom"
			if tc.feed != "" {
				feed = tc.feed
			}
			releasesFeed, err := ioutil.ReadFile(feed)
			require.NoError(t, err)

			outFile, err := ioutil.TempFile("", "gitlab-update-test")
			require.NoError(t, err)
			outFile.Close()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/releases.atom":
					w.Write(releasesFeed)
				case fmt.Sprintf("/releases/download/v%s/gitlab_v%s_%s_%s", tc.updatedVersion, tc.updatedVersion, runtime.GOOS, runtime.GOARCH):
					w.Write([]byte(tc.updatedBinary))
				default:
					t.Errorf("Unexpected request for %s", r.RequestURI)
				}
			}))
			defer ts.Close()

			repo = ts.URL
			if tc.repo != "" {
				repo = tc.repo
			}
			config.Version = tc.currentVersion
			var out strings.Builder

			err = updateCommand(tc.dryRun, tc.pre, tc.upgradeMajor, &out, func() (string, error) { return outFile.Name(), nil })

			if tc.outErr {
				require.Error(t, err)
				assert.True(t, tc.out.MatchString(err.Error()), "unexpected error: '%s' does not match '%s'", err.Error(), tc.out.String())
			} else {
				require.NoError(t, err)
				assert.True(t, tc.out.MatchString(out.String()), "unexpected output: '%s' does not match '%s'", out.String(), tc.out.String())
			}
			updatedContent, err := ioutil.ReadFile(outFile.Name())
			require.NoError(t, err)
			assert.Equal(t, tc.updatedBinary, updatedContent)
		})
	}
}
