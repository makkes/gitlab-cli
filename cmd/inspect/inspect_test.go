package inspect

import (
	"testing"

	"github.com/makkes/gitlab-cli/mock"
)

func TestSubCommands(t *testing.T) {
	cmd := NewCommand(mock.Client{})
	subCmds := cmd.Commands()
	if len(subCmds) != 3 {
		t.Errorf("Expected 1 sub-command but got %d", len(subCmds))
	}
	if cmd.UseLine() != "inspect" {
		t.Errorf("Unexpected usage line '%s'", cmd.UseLine())
	}
}
