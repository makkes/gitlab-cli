package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var configFile = ".gitlab-cli.json"

type Config struct {
	Token string
	URL   string
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

func checkPermissions(f *os.File) {
	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("Error checking file permissions: %s\n", err)
	}
	if fi.Mode() != 0600 {
		fmt.Printf("Correcting configuration file permissions from %#o to %#o\n", fi.Mode(), 0600)
		err = f.Chmod(0600)
		if err != nil {
			fmt.Printf("Error correcting configuration file permissions: %s\n", err)
			return
		}
	}
}

func Read() *Config {
	f, err := os.Open(gitlabCLIConf())
	if err != nil {
		fmt.Printf("Error opening configuration file: %s\n", err)
		return &defaultConfig
	}
	defer f.Close()

	checkPermissions(f)

	bytes, err := ioutil.ReadAll(f)
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

	// this converts legacy configurations
	if config.URL == "" {
		fmt.Fprintf(os.Stderr, "Converting legacy configuration file to new format\n")
		config.URL = "https://gitlab.com"
		config.Write()
	}
	return &config
}

func (c *Config) Write() {
	unindented, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
	var indented bytes.Buffer
	if err := json.Indent(&indented, unindented, "", "  "); err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
	}
	err = ioutil.WriteFile(gitlabCLIConf(), indented.Bytes(), 0600)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
}

func gitlabCLIConf() string {
	return path.Join(os.Getenv("HOME"), configFile)
}
