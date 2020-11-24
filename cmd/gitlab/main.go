package main

import (
	"github.com/makkes/gitlab-cli/v3/cmd"
	"github.com/makkes/gitlab-cli/v3/config"
)

func main() {
	cmd.Execute(config.Read())
}
