package projects

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/template"

	"github.com/makkes/gitlab-cli/v3/api"
	"github.com/makkes/gitlab-cli/v3/config"
	"github.com/makkes/gitlab-cli/v3/table"
	"github.com/spf13/cobra"
)

func projectsCommand(client api.Client, cfg config.Config, quiet bool, format string, page int, membership bool, out io.Writer) error {
	var path string
	if membership {
		path = fmt.Sprintf("/projects?membership=true&page=%d", page)
	} else {
		path = fmt.Sprintf("/users/${user}/projects?page=%d", page)
	}
	resp, status, err := client.Get(path)

	if err != nil {
		if status == 404 {
			return fmt.Errorf("cannot list projects: User %s not found. Please check your configuration", cfg.Get("user"))
		}
		return fmt.Errorf("Cannot list projects: %s", err)
	}
	projects := make([]api.Project, 0)
	err = json.Unmarshal(resp, &projects)
	if err != nil {
		return err
	}

	// put all project name => ID mappings into the cache
	for _, p := range projects {
		cfg.Cache().Put("projects", p.Name, strconv.Itoa(p.ID))
	}
	cfg.Write()

	if format != "" {
		tmpl, err := template.New("").Parse(format)
		if err != nil {
			return fmt.Errorf("template parsing error: %s", err)
		}

		for _, p := range projects {
			err = tmpl.Execute(out, p)
			if err != nil {
				return fmt.Errorf("template parsing error: %s", err)
			}
			fmt.Fprintln(out)
		}
		return nil
	}

	if quiet {
		for _, p := range projects {
			fmt.Fprintln(out, p.ID)
		}
		return nil
	}
	table.PrintProjects(out, projects)
	return nil
}

func NewCommand(client api.Client, cfg config.Config) *cobra.Command {
	var quiet *bool
	var format *string
	var page *int
	var membership *bool

	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List all your projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return projectsCommand(client, cfg, *quiet, *format, *page, *membership, os.Stdout)
		},
	}

	quiet = cmd.Flags().BoolP("quiet", "q", false, "Only display numeric IDs")
	format = cmd.Flags().String("format", "", "Pretty-print projects using a Go template")
	page = cmd.Flags().Int("page", 1, "Page of results to display")
	membership = cmd.Flags().Bool("member", false, "Displays projects you are a member of")

	return cmd
}
