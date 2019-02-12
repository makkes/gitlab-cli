package mock

import (
	"io"

	"github.com/makkes/gitlab-cli/config"

	"github.com/makkes/gitlab-cli/api"
)

type MockClient struct {
	Res    []byte
	Status int
	Err    error
}

func (m MockClient) Get(path string) ([]byte, int, error) {
	return m.Res, m.Status, m.Err
}

func (m MockClient) Post(path string, body io.Reader) ([]byte, int, error) {
	return m.Res, 0, m.Err
}

func (m MockClient) Delete(path string) (int, error) {
	return 0, m.Err
}

func (m MockClient) FindProject(nameOrID string) (*api.Project, error) {
	return nil, nil
}

func (m MockClient) FindProjectDetails(nameOrID string) ([]byte, error) {
	return nil, nil
}

type MockCache struct {
	Calls [][]string
	Cache map[string]map[string]string
}

func (c MockCache) Flush() {}

func (c MockCache) Get(cacheName, key string) string {
	return ""
}

func (c *MockCache) Put(cacheName, key, value string) {
	if c.Calls == nil {
		c.Calls = make([][]string, 0)
	}
	c.Calls = append(c.Calls, []string{cacheName, key, value})
}

type MockConfig struct {
	CacheData  config.Cache
	Cfg        map[string]string
	WriteCalls int
}

func (c MockConfig) Cache() config.Cache {
	return c.CacheData
}

func (c *MockConfig) Write() {
	c.WriteCalls++
}

func (c MockConfig) Get(key string) string {
	return c.Cfg[key]
}

func (c MockConfig) Set(key, value string) {}
