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
var Version = "none"

// keys for the configuration.
const (
	User  = "user"
	Token = "token"
	URL   = "url"
)

type Config interface {
	Cache() Cache
	Write()
	Get(key string) string
	Set(key, value string)
}

type inMemoryConfig struct {
	version int
	cfg     map[string]string
	cache   *MapCache
}

func (c *inMemoryConfig) Get(key string) string {
	return c.cfg[key]
}

func (c *inMemoryConfig) Set(key, value string) {
	c.cfg[key] = value
}

func (c *inMemoryConfig) Cache() Cache {
	return c.cache
}

// MarshalJSON makes InMemoryConfig implement json.Marshaler.
func (c *inMemoryConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version int
		Config  map[string]string
		Cache   *MapCache
	}{
		Version: c.version,
		Config:  c.cfg,
		Cache:   c.cache,
	})
}

// UnmarshalJSON makes Cache implement json.Unmarshaler.
func (c *inMemoryConfig) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Version int
		Config  map[string]string
		Cache   *MapCache
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.version = aux.Version
	c.cfg = aux.Config
	c.cache = aux.Cache
	return nil
}

var defaultConfig = inMemoryConfig{
	version: 1,
	cfg:     make(map[string]string),
	cache: &MapCache{
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

func Read() Config {
	f, err := os.Open(gitlabCLIConf())
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("WARNING: Error opening configuration file: %s\n\n", err)
		}
		return &defaultConfig
	}
	defer f.Close()

	checkPermissions(f)

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("WARNING: Error reading configuration: %s\n\n", err)
		return &defaultConfig
	}
	var config inMemoryConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Printf("WARNING: Error reading configuration: %s\n\n", err)
		return &defaultConfig
	}

	writeConfigFile := false

	// We're presented with an old-style config file and need to updated its format
	if config.version == 0 {
		fmt.Fprintf(os.Stderr, "Converting legacy configuration file to new format\n")
		var oldConfig map[string]interface{}
		err = json.Unmarshal(bytes, &oldConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: Error reading legacy configuration file: %s\n", err)
			return &defaultConfig
		}
		config.cfg = map[string]string{
			User:  oldConfig["User"].(string),
			Token: oldConfig["Token"].(string),
		}
		if oldConfig["URL"] == nil {
			config.cfg[URL] = "https://gitlab.com"
		} else {
			config.cfg[URL] = oldConfig["URL"].(string)
		}
		config.version = 1
		writeConfigFile = true
	}

	if config.cache == nil {
		config.cache = NewMapCache()
		writeConfigFile = true
	}

	if writeConfigFile {
		config.Write()
	}

	return &config
}

func (c *inMemoryConfig) Write() {
	unindented, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("Error writing configuration: %s\n", err)
		return
	}
	var indented bytes.Buffer
	if err = json.Indent(&indented, unindented, "", "  "); err != nil {
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
