package projects

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/template"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
	"github.com/makkes/gitlab-cli/table"
	"github.com/spf13/cobra"
)

func projectsCommand(client api.Client, cfg config.Config, quiet bool, format string, page int, membership bool, group string, out io.Writer) error {
	var path string
	if group != "" {
		path = fmt.Sprintf("/groups/%s/projects?page=%d", group, page)
	} else if membership {
		path = fmt.Sprintf("/projects?membership=true&page=%d", page)
	} else {
		path = fmt.Sprintf("/users/${user}/projects?page=%d", page)
	}
	resp, status, err := client.Get(path)

	if err != nil {
		if status == 404 {
			if group != "" {
				return fmt.Errorf("cannot list projects: Group %s not found", group)
			}
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
		Use:   "projects [GROUP]",
		Short: "List all your projects or projects within a specific group",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			group := ""
			if len(args) > 0 {
				group = args[0]
			}
			return projectsCommand(client, cfg, *quiet, *format, *page, *membership, group, os.Stdout)
		},
	}

	quiet = cmd.Flags().BoolP("quiet", "q", false, "Only display numeric IDs")
	format = cmd.Flags().String("format", "", "Pretty-print projects using a Go template")
	page = cmd.Flags().Int("page", 1, "Page of results to display")
	membership = cmd.Flags().Bool("member", false, "Displays projects you are a member of")

	return cmd
}
