package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/makkes/gitlab-cli/config"

	"github.com/stretchr/testify/assert"
)

func TestUpdateHappyPath(t *testing.T) {
	releasesFeed, err := ioutil.ReadFile("releases_test.atom")
	require.NoError(t, err)
	outFile, err := ioutil.TempFile("", "gitlab-update-test")
	outFile.Close()
	expectedUpdatedContent := []byte("this is the updated gitlab binary")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/releases.atom":
			w.Write(releasesFeed)
		case fmt.Sprintf("/releases/download/v3.6.3/gitlab_v3.6.3_%s_%s", runtime.GOOS, runtime.GOARCH):
			w.Write(expectedUpdatedContent)
		default:
			t.Errorf("Unexpected request for %s", r.RequestURI)
		}
	}))
	defer ts.Close()

	repo = ts.URL
	config.Version = "3.6.2-55-ghf448b"
	var out strings.Builder

	err = updateCommand(false, &out, func() (string, error) { return outFile.Name(), nil })

	require.NoError(t, err)
	assert.Equal(t, "Updating to v3.6.3\n", out.String())
	updatedContent, err := ioutil.ReadFile(outFile.Name())
	require.NoError(t, err)
	assert.Equal(t, expectedUpdatedContent, updatedContent)
}
