package projects

import (
	"fmt"
	"strings"
	"testing"

	"github.com/makkes/gitlab-cli/mock"
)

func TestClientError(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Err: fmt.Errorf("Some client error"),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, true, "", 0, false, &out)
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
	client := mock.Client{
		Status: 404,
		Err:    fmt.Errorf("Project not found"),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
		Cfg:       map[string]string{"user": "Dilbert"},
	}
	err := projectsCommand(client, config, true, "", 0, false, &out)
	if err == nil {
		t.Error("Expected a non-nil error")
	}
	if err.Error() != "cannot list projects: User Dilbert not found. Please check your configuration" {
		t.Errorf("Unexpected error message '%s'", err)
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}

func TestBrokenResponse(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte("this is not JSON"),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, true, "", 0, false, &out)
	if err == nil {
		t.Error("Expected a non-nil error")
	}
	if out.String() != "" {
		t.Errorf("Expected output to be empty but it is '%s'", out.String())
	}
}
func TestEmptyResult(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte(`[]`),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, true, "", 0, false, &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "" {
		t.Errorf("Expected empty output but got '%s'", out.String())
	}
}

func TestQuietOutput(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte(`[{"id": 123}, {"id": 456}]`),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, true, "", 0, false, &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "123\n456\n" {
		t.Errorf("Unexpected output '%s'", out.String())
	}
}

func TestFormattedOutput(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte(`[{"id": 123, "name":"broken arrow"}, {"id": 456, "name":"almanac"}]`),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, false, "{{.Name}}", 0, false, &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "broken arrow\nalmanac\n" {
		t.Errorf("Unexpected output '%s'", out.String())
	}
}
func TestFormattedOutputError(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte(`[{"id": 123, "name":"broken arrow"}, {"id": 456, "name":"almanac"}]`),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, false, "{{.Broken}", 0, false, &out)
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
	client := mock.Client{
		Res: []byte(`[{"id": 123, "name":"broken arrow"}, {"id": 456, "name":"almanac"}]`),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, false, "{{.Name}}", 0, false, &mockOutput{err: fmt.Errorf("some error")})
	if err == nil {
		t.Error("Expected a non-nil error")
	}
}

func TestTableOutput(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte(`[{"id": 123, "name":"broken arrow"}, {"id": 456, "name":"almanac"}]`),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, false, "", 0, false, &out)
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
	client := mock.Client{
		Res: []byte(`[]`),
	}
	config := &mock.Config{
		CacheData: &mock.Cache{},
	}
	err := projectsCommand(client, config, false, "", 0, false, &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if out.String() != "ID              NAME                                     URL                                                CLONE                                             \n" {
		t.Errorf("Unexpected output '%s'", out.String())
	}
}

func TestCache(t *testing.T) {
	var out strings.Builder
	client := mock.Client{
		Res: []byte(`[{"id": 123, "name": "p1"}, {"id": 456, "name":"p2"}]`),
	}
	cache := mock.Cache{}
	config := &mock.Config{
		CacheData: &cache,
	}
	err := projectsCommand(client, config, true, "", 0, false, &out)
	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
	if len(cache.Calls) != 2 {
		t.Errorf("Expected two values to be cached but got %d", len(cache.Calls))
	}
	call := cache.Calls[0]
	if call[0] != "projects" || call[1] != "p1" || call[2] != "123" {
		t.Errorf("Unexpected PUT call on cache: %s", call)
	}
	call = cache.Calls[1]
	if call[0] != "projects" || call[1] != "p2" || call[2] != "456" {
		t.Errorf("Unexpected PUT call on cache: %s", call)
	}
	if config.WriteCalls != 1 {
		t.Errorf("Expected the config to be written once but was %d", config.WriteCalls)
	}
}

func TestNewCommand(t *testing.T) {
	cmd := NewCommand(mock.Client{}, &mock.Config{})
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
