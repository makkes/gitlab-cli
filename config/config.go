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
	Cache *Cache
}

var defaultConfig = Config{
	Token: "",
	User:  "",
	Cache: &Cache{
		data: make(map[string]map[string]string),
	},
}

func Read() *Config {
	bytes, err := ioutil.ReadFile(gitlabCLIConf())
	if err != nil {
		fmt.Printf("WARNING: Error reading configuration: %s\n", err)
		return &defaultConfig
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Printf("WARNING: Error reading configuration: %s\n", err)
		return &defaultConfig
	}
	if config.Cache == nil {
		config.Cache = NewCache()
	}
	return &config
}

func (c *Config) Write() {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
	bytes = append(bytes, '\n')
	err = ioutil.WriteFile(gitlabCLIConf(), bytes, 0644)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
	fmt.Printf("Configuration stored in %s\n", gitlabCLIConf())
}

func gitlabCLIConf() string {
	return path.Join(os.Getenv("HOME"), configFile)
}
