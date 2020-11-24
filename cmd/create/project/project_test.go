package project

import (
	"fmt"
	"strings"
	"testing"

	"github.com/makkes/gitlab-cli/mock"
)

func TestClientError(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Err: fmt.Errorf("Some client error"),
	}
	err := createCommand(client, []string{"new project"}, &out)
	if err == nil {
		t.Error("Expected a non-nil error")
	}
	if err.Error() != "Cannot create project: Some client error" {
		t.Errorf("Unexpected error message '%s'", err)
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}

func TestUnknownResultFromGitLab(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte(fmt.Sprintf("This is not JSON")),
	}
	err := createCommand(client, []string{"new project"}, &out)
	if err == nil {
		t.Error("Expected a non-nil error")
	}
	if err.Error() != "Unexpected result from GitLab API: invalid character 'T' looking for beginning of value" {
		t.Errorf("Unexpected error message '%s'", err)
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}

func TestHappyPath(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte(fmt.Sprintf(`{
			"name": "new project",
			"ssh_url_to_repo": "SSH clone URL",
			"http_url_to_repo": "HTTP clone URL"}`)),
	}
	err := createCommand(client, []string{"new project"}, &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "Project 'new project' created\nClone via SSH: SSH clone URL\nClone via HTTP: HTTP clone URL\n" {
		t.Errorf("Unexpected output '%s'", out.String())
	}
}
