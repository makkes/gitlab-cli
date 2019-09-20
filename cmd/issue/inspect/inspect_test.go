package inspect

import (
	"fmt"
	"strings"
	"testing"

	"github.com/makkes/gitlab-cli/mock"
)

func TestWrongArgFormat(t *testing.T) {
	in := []struct {
		name string
		args []string
	}{
		{"empty", []string{""}},
		{"no colon", []string{"ab"}},
		{"no 1st arg", []string{":b"}},
		{"no 2nd arg", []string{"a:"}},
	}
	for _, tt := range in {
		t.Run(tt.name, func(t *testing.T) {
			var out strings.Builder
			err := inspectCommand(mock.Client{}, tt.args, &out)
			if out.String() != "" {
				t.Errorf("Expected output to be empty but it is '%s'", out.String())
			}
			if err == nil {
				t.Errorf("Expected a non-nil error")
			}
			if err.Error() != "ID must be of the form PROJECT_ID:ISSUE_ID" {
				t.Errorf("Expected error message to be '' but is '%s'", err.Error())
			}

		})
	}
}

func TestUnknownProject(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Status: 404,
		Err:    fmt.Errorf("Project not found"),
	}
	err := inspectCommand(client, []string{"pid:iid"}, &out)
	if err == nil {
		t.Errorf("Expected non-nil error")
	}
	if err.Error() != "Issue pid:iid not found" {
		t.Errorf("Unexpected error '%s'", err.Error())
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}

func TestUnknownClientError(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Status: 500,
		Err:    fmt.Errorf("Server broken"),
	}
	err := inspectCommand(client, []string{"pid:iid"}, &out)
	if err == nil {
		t.Errorf("Expected non-nil error")
	}
	if err.Error() != "Server broken" {
		t.Errorf("Unexpected error '%s'", err.Error())
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}

func TestHappyPath(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte("\"hello\""),
	}
	err := inspectCommand(client, []string{"pid:iid"}, &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "\"hello\"\n" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}
