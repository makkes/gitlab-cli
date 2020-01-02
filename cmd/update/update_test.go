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
		currentVersion string
		updatedBinary  []byte
		dryRun         bool
		repo           string
		out            *regexp.Regexp
		outErr         bool
	}{
		"update happy path": {
			currentVersion: "3.6.2-55-ghf448b",
			updatedBinary:  []byte("this is the updated gitlab binary"),
			out:            regexp.MustCompile(`^Updating to v3.6.3\n`),
		},
		"no update happy path": {
			currentVersion: "3.6.3",
			updatedBinary:  []byte{},
			out:            regexp.MustCompile(`^No update available, yet.\n`),
		},
		"dry-run no update happy path": {
			currentVersion: "3.6.3",
			updatedBinary:  []byte{},
			dryRun:         true,
			out:            regexp.MustCompile(`^No update available, yet.\n`),
		},
		"dry-run update happy path": {
			currentVersion: "3.6.2-55-ghf448b",
			updatedBinary:  []byte{},
			dryRun:         true,
			out:            regexp.MustCompile(`^A new version is available: v3.6.3\nSee http://127.0.0.1:\d+/releases/v3.6.3 for details\n`),
		},
		"unreachable repo": {
			repo:          "http://doesntexist",
			outErr:        true,
			out:           regexp.MustCompile(`^Could not fetch latest version: `),
			updatedBinary: []byte{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(tt *testing.T) {
			releasesFeed, err := ioutil.ReadFile("releases_test.atom")
			require.NoError(t, err)

			outFile, err := ioutil.TempFile("", "gitlab-update-test")
			require.NoError(t, err)
			outFile.Close()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/releases.atom":
					w.Write(releasesFeed)
				case fmt.Sprintf("/releases/download/v3.6.3/gitlab_v3.6.3_%s_%s", runtime.GOOS, runtime.GOARCH):
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

			err = updateCommand(tc.dryRun, &out, func() (string, error) { return outFile.Name(), nil })

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
