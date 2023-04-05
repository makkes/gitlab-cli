package mock

import (
	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
)

type Client struct {
	Res    []byte
	Status int
	Err    error
}

var _ api.Client = Client{}

func (m Client) Get(path string) ([]byte, int, error) {
	return m.Res, m.Status, m.Err
}

func (m Client) GetAccessTokens(pid string) ([]api.ProjectAccessToken, error) {
	return nil, nil
}

func (m Client) Post(path string, body interface{}) ([]byte, int, error) {
	return m.Res, 0, m.Err
}

func (m Client) Delete(path string) (int, error) {
	return 0, m.Err
}

func (m Client) FindProject(nameOrID string) (*api.Project, error) {
	return nil, nil
}

func (m Client) FindProjectDetails(nameOrID string) ([]byte, error) {
	return nil, nil
}

func (m Client) Login(token, url string) (string, error) {
	return "", nil
}

func (m Client) GetPipelineDetails(projectID, pipelineID string) ([]byte, error) {
	return nil, nil
}

type Cache struct {
	Calls [][]string
	Cache map[string]map[string]string
}

func (c Cache) Flush() {}

func (c Cache) Get(cacheName, key string) string {
	return ""
}

func (c *Cache) Put(cacheName, key, value string) {
	if c.Calls == nil {
		c.Calls = make([][]string, 0)
	}
	c.Calls = append(c.Calls, []string{cacheName, key, value})
}

type Config struct {
	CacheData  config.Cache
	Cfg        map[string]string
	WriteCalls int
}

func (c Config) Cache() config.Cache {
	return c.CacheData
}

func (c *Config) Write() {
	c.WriteCalls++
}

func (c Config) Get(key string) string {
	return c.Cfg[key]
}

func (c Config) Set(key, value string) {}
