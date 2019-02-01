package main

import (
	"github.com/makkes/gitlab-cli/config"
	"github.com/makkes/gitlab-cli/cmd"
)

func main() {
	cmd.Execute(config.Read())
}
