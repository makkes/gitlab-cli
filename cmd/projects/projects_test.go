package projects

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/makkes/gitlab-cli/api"
	"github.com/makkes/gitlab-cli/config"
)

type mockClient struct {
	res    []byte
	status int
	err    error
}

func (m mockClient) Get(path string) ([]byte, int, error) {
	return m.res, m.status, m.err
}

func (m mockClient) Post(path string, body io.Reader) ([]byte, int, error) {
	return m.res, 0, m.err
}

func (m mockClient) Delete(path string) (int, error) {
	return 0, m.err
}

func (m mockClient) FindProject(nameOrID string) (*api.Project, error) {
	return nil, nil
}

type mockCache struct {
	calls [][]string
	cache map[string]map[string]string
}

func (c mockCache) Flush() {}

func (c mockCache) Get(cacheName, key string) string {
	return ""
}

func (c *mockCache) Put(cacheName, key, value string) {
	if c.calls == nil {
		c.calls = make([][]string, 0)
	}
	c.calls = append(c.calls, []string{cacheName, key, value})
}

type mockConfig struct {
	cache      config.Cache
	cfg        map[string]string
	writeCalls int
}

func (c mockConfig) Cache() config.Cache {
	return c.cache
}

func (c *mockConfig) Write() {
	c.writeCalls++
}

func (c mockConfig) Get(key string) string {
	return c.cfg[key]
}

func (c mockConfig) Set(key, value string) {}

func TestClientError(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		err: fmt.Errorf("Some client error"),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, true, "", &out)
	if err == nil {
		t.Error("Expected a non-nil error")
	}
	if err.Error() != "Cannot list projects: Some client error" {
		t.Errorf("Unexpected error message '%s'", err)
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}

func TestUnknownProject(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		status: 404,
		err:    fmt.Errorf("Project not found"),
	}
	config := &mockConfig{
		cache: &mockCache{},
		cfg:   map[string]string{"user": "Dilbert"},
	}
	err := projectsCommand(client, config, true, "", &out)
	if err == nil {
		t.Error("Expected a non-nil error")
	}
	if err.Error() != "Cannot list projects: User Dilbert not found. Check your configuration!" {
		t.Errorf("Unexpected error message '%s'", err)
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}

func TestBrokenResponse(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		res: []byte("this is not JSON"),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, true, "", &out)
	if err == nil {
		t.Error("Expected a non-nil error")
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}
func TestEmptyResult(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		res: []byte(`[]`),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, true, "", &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "" {
		t.Errorf("Expected empty output but got '%s'", out.String())
	}
}

func TestQuietOutput(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		res: []byte(`[{"id": 123}, {"id": 456}]`),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, true, "", &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "123\n456\n" {
		t.Errorf("Unexpected output '%s'", out.String())
	}
}

func TestFormattedOutput(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		res: []byte(`[{"id": 123, "name":"broken arrow"}, {"id": 456, "name":"almanac"}]`),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, false, "{{.Name}}", &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "broken arrow\nalmanac\n" {
		t.Errorf("Unexpected output '%s'", out.String())
	}
}
func TestFormattedOutputError(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		res: []byte(`[{"id": 123, "name":"broken arrow"}, {"id": 456, "name":"almanac"}]`),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, false, "{{.Broken}", &out)
	if err == nil {
		t.Error("Expected a non-nil error")
	}
	if out.String() != "" {
		t.Errorf("Expected empty output but got '%s'", out.String())
	}
}

type mockOutput struct {
	n   int
	err error
}

func (m mockOutput) Write(p []byte) (n int, err error) {
	return m.n, m.err
}
func TestTemplateExecutionError(t *testing.T) {
	client := mockClient{
		res: []byte(`[{"id": 123, "name":"broken arrow"}, {"id": 456, "name":"almanac"}]`),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, false, "{{.Name}}", &mockOutput{err: fmt.Errorf("some error")})
	if err == nil {
		t.Error("Expected a non-nil error")
	}
}

func TestTableOutput(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		res: []byte(`[{"id": 123, "name":"broken arrow"}, {"id": 456, "name":"almanac"}]`),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, false, "", &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != `ID              NAME                                     URL                                                CLONE                                             
123             broken arrow                                                                                                                                  
456             almanac                                                                                                                                       
` {
		t.Errorf("Unexpected output '%s'", out.String())
	}
}

func TestEmptyTableOutput(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		res: []byte(`[]`),
	}
	config := &mockConfig{
		cache: &mockCache{},
	}
	err := projectsCommand(client, config, false, "", &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "ID              NAME                                     URL                                                CLONE                                             \n" {
		t.Errorf("Unexpected output '%s'", out.String())
	}
}

func TestCache(t *testing.T) {
	var out strings.Builder
	client := mockClient{
		res: []byte(`[{"id": 123, "name": "p1"}, {"id": 456, "name":"p2"}]`),
	}
	cache := mockCache{}
	config := &mockConfig{
		cache: &cache,
	}
	err := projectsCommand(client, config, true, "", &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if len(cache.calls) != 2 {
		t.Errorf("Expected two values to be cached but got %d", len(cache.calls))
	}
	call := cache.calls[0]
	if call[0] != "projects" || call[1] != "p1" || call[2] != "123" {
		t.Errorf("Unexpected PUT call on cache: %s", call)
	}
	call = cache.calls[1]
	if call[0] != "projects" || call[1] != "p2" || call[2] != "456" {
		t.Errorf("Unexpected PUT call on cache: %s", call)
	}
	if config.writeCalls != 1 {
		t.Errorf("Expected the config to be written once but was %d", config.writeCalls)
	}
}

func TestNewCommand(t *testing.T) {
	cmd := NewCommand(mockClient{}, &mockConfig{})
	flags := cmd.Flags()

	quietFlag := flags.Lookup("quiet")
	if quietFlag == nil {
		t.Errorf("Expected 'quiet' flag to exist")
	}
	if quietFlag.Value.Type() != "bool" {
		t.Errorf("Expected 'quiet' flag to be a bool but is %s", quietFlag.Value.Type())
	}
	if quietFlag.DefValue != "false" {
		t.Errorf("Expected default value of 'quiet' flag to be 'false' but is '%s'", quietFlag.DefValue)
	}

	formatFlag := flags.Lookup("format")
	if formatFlag == nil {
		t.Errorf("Expected 'format' flag to exist")
	}
	if formatFlag.Value.Type() != "string" {
		t.Errorf("Expected 'format' flag to be a string but is %s", formatFlag.Value.Type())
	}
	if formatFlag.DefValue != "" {
		t.Errorf("Expected default value of 'format' flag to be '' but is '%s'", formatFlag.DefValue)
	}
}
