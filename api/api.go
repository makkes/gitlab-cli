package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/makkes/gitlab-cli/config"
)

type Project struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	SSHGitURL      string `json:"ssh_url_to_repo"`
	HTTPGitURL     string `json:"http_url_to_repo"`
	URL            string `json:"web_url"`
	LastActivityAt string `json:"last_activity_at"`
}

type Pipeline struct {
	ID     int
	Status string
}

type PipelineDetails struct {
	ProjectID int
	ID        int
	Status    string
	URL       string `json:"web_url"`
	Duration  int
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

type Client interface {
	Get(path string) ([]byte, error)
}

type APIClient struct {
	basePath string
	config   config.Config
	client   http.Client
}

func NewAPIClient(cfg config.Config) APIClient {
	client := http.Client{}
	return APIClient{
		basePath: "/api/v4",
		config:   cfg,
		client:   client,
	}
}

func (c APIClient) parse(input string) string {
	return strings.Replace(input, "${user}", c.config.Get(config.User), -1)
}

func (c *APIClient) Login(token, url string) (error, string) {
	c.config.Set(config.Token, token)
	c.config.Set(config.URL, url)
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
	c.config.Set(config.User, user.Username)
	c.config.Cache().Flush()
	return nil, user.Username
}

func (c APIClient) GetPipelineDetails(projectID, pipelineID string) ([]byte, error) {
	resp, err := c.Get(fmt.Sprintf("/projects/%s/pipelines/%s", url.PathEscape(projectID), url.PathEscape(pipelineID)))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindProjectDetails searches for a project by its ID or its name,
// with the ID having precedence over the name and returns the
// raw JSON object as byte array.
func (c APIClient) FindProjectDetails(nameOrID string) ([]byte, error) {
	// first try to get the project by its cached ID
	if cachedID := c.config.Cache().Get("projects", nameOrID); cachedID != "" {
		resp, err := c.Get("/projects/" + url.PathEscape(cachedID))
		if err == nil {
			return resp, nil
		}
	}

	// then try to find the project by its ID
	resp, err := c.Get("/projects/" + url.PathEscape(nameOrID))
	if err == nil {
		return resp, nil
	}

	// now try to find the project by name as a last resort
	resp, err = c.Get("/users/${user}/projects/?search=" + url.QueryEscape(nameOrID))
	if err != nil {
		return nil, err
	}
	projects := make([]map[string]interface{}, 0)
	err = json.Unmarshal(resp, &projects)
	if err != nil {
		return nil, err
	}
	if len(projects) <= 0 {
		return nil, errors.New("No project found")
	}
	c.config.Cache().Put("projects", nameOrID, strconv.Itoa(int((projects[0]["id"].(float64)))))
	c.config.Write()
	res, err := json.Marshal(projects[0])
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FindProject searches for a project by its ID or its name,
// with the ID having precedence over the name.
func (c APIClient) FindProject(nameOrID string) (*Project, error) {
	projectBytes, err := c.FindProjectDetails(nameOrID)
	if err != nil {
		return nil, err
	}
	var project Project
	err = json.Unmarshal(projectBytes, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

var ErrNotLoggedIn = errors.New("You are not logged in")

func (c APIClient) Get(path string) ([]byte, error) {
	if c.config == nil {
		return nil, ErrNotLoggedIn
	}
	req, err := http.NewRequest("GET", c.config.Get(config.URL)+c.basePath+c.parse(path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Private-Token", c.config.Get(config.Token))
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error querying GitLab, HTTP status is %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
