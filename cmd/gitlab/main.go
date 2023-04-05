package main

import (
	"math/rand"
	"time"

	"github.com/makkes/gitlab-cli/cmd"
	"github.com/makkes/gitlab-cli/config"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute(config.Read())
}
