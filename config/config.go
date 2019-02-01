package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var configFile = ".gitlab-cli.json"

type Config struct {
	Token string
	User  string
}

func Read() *Config {
	bytes, err := ioutil.ReadFile(gitlabCLIConf())
	if err != nil {
		fmt.Printf("WARNING: Error reading configuration: %s\n", err)
		return nil
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Printf("WARNING: Error reading configuration: %s\n", err)
		return nil
	}
	return &config
}

func gitlabCLIConf() string {
	return path.Join(os.Getenv("HOME"), configFile)
}

func WriteLoginCredentials(token, user string) {
	c := Config{token, user}
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
	err = ioutil.WriteFile(gitlabCLIConf(), bytes, 0644)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
}
