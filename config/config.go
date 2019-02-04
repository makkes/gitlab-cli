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

var defaultConfig = Config{}

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
	return &config
}

func (c *Config) Write() {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
	err = ioutil.WriteFile(gitlabCLIConf(), bytes, 0600)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
}

func gitlabCLIConf() string {
	return path.Join(os.Getenv("HOME"), configFile)
}
