package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/makkes/gitlab-cli/config"
)

type Project struct {
	ID   int
	Name string
	URL  string `json:"web_url"`
}

type Pipeline struct {
	ID     int
	Status string
}

type PipelineDetails struct {
	ID       int
	Status   string
	URL      string `json:"web_url"`
	Duration int
}

type Pipelines []Pipeline

func (p Pipelines) Filter(cb func(Pipeline) bool) Pipelines {
	res := make(Pipelines, 0)
	for _, pipeline := range p {
		if cb(pipeline) {
			res = append(res, pipeline)
		}
	}
	return res
}

type APIClient struct {
	baseURL string
	config  *config.Config
	client  http.Client
}

func NewAPIClient(cfg *config.Config) APIClient {
	client := http.Client{}
	return APIClient{
		baseURL: "https://gitlab.com/api/v4",
		config:  cfg,
		client:  client,
	}
}

func (c APIClient) parse(input string) string {
	return strings.Replace(input, "${user}", c.config.User, -1)
}

func (c *APIClient) Login(token string) (error, string) {
	c.config.Token = token
	res, err := c.Get("/user")
	if err != nil {
		return err, ""
	}
	var user struct {
		Username string
	}
	err = json.Unmarshal(res, &user)
	if err != nil {
		return err, ""
	}
	c.config.User = user.Username
	return nil, user.Username
}

func (c APIClient) FindProject(name string) ([]Project, error) {
	resp, err := c.Get("/users/${user}/projects/?search=" + url.QueryEscape(name))
	if err != nil {
		return nil, err
	}
	projects := make([]Project, 0)
	err = json.Unmarshal(resp, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil

}

var ErrNotLoggedIn error = errors.New("You are not logged in")

func (c APIClient) Get(path string) ([]byte, error) {
	if c.config == nil {
		return nil, ErrNotLoggedIn
	}
	req, err := http.NewRequest("GET", c.baseURL+c.parse(path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Private-Token", c.config.Token)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error querying GitLab, HTTP status is %d", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
