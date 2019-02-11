package main

import (
	"github.com/makkes/gitlab-cli/cmd"
	"github.com/makkes/gitlab-cli/config"
)

func main() {
	cmd.Execute(config.Read())
}
