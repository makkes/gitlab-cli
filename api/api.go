package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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

type Issue struct {
	ProjectID int    `json:"project_id"`
	ID        int    `json:"iid"`
	Title     string `json:"title"`
	URL       string `json:"web_url"`
	State     string `json:"state"`
}

type Pipeline struct {
	ID     int
	Status string
}

type PipelineDetails struct {
	ProjectID        int
	ID               int
	Status           string
	URL              string    `json:"web_url"`
	RecordedDuration *int      `json:"duration"`
	StartedAt        time.Time `json:"started_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	FinishedAt       time.Time `json:"finished_at"`
}

type Var struct {
	Key              string
	Value            string
	Protected        bool
	EnvironmentScope string `json:"environment_scope"`
}

func (pd PipelineDetails) Duration(now time.Time) string {
	if pd.Status == "running" {
		started := pd.StartedAt
		if !pd.FinishedAt.IsZero() {
			started = pd.UpdatedAt
		}
		return now.Sub(started).Truncate(time.Second).String()
	}
	if pd.RecordedDuration == nil {
		return "-"
	}
	return time.Duration(int(time.Second) * *pd.RecordedDuration).String()
}

type Pipelines []Pipeline

func (p Pipelines) Filter(cb func(int, Pipeline) bool) Pipelines {
	res := make(Pipelines, 0)
	for idx, pipeline := range p {
		if cb(idx, pipeline) {
			res = append(res, pipeline)
		}
	}
	return res
}

type Client interface {
	Get(path string) ([]byte, int, error)
	Post(path string, body io.Reader) ([]byte, int, error)
	Delete(path string) (int, error)
	FindProject(nameOrID string) (*Project, error)
	FindProjectDetails(nameOrID string) ([]byte, error)
	Login(token, url string) (string, error)
	GetPipelineDetails(projectID, pipelineID string) ([]byte, error)
}

type HTTPClient struct {
	basePath string
	config   config.Config
	client   http.Client
}

func NewAPIClient(cfg config.Config) *HTTPClient {
	client := http.Client{}
	return &HTTPClient{
		basePath: "/api/v4",
		config:   cfg,
		client:   client,
	}
}

func (c HTTPClient) parse(input string) string {
	return strings.Replace(input, "${user}", c.config.Get(config.User), -1)
}

func (c *HTTPClient) Login(token, url string) (string, error) {
	c.config.Set(config.Token, token)
	c.config.Set(config.URL, url)
	res, _, err := c.Get("/user")
	if err != nil {
		return "", err
	}
	var user struct {
		Username string
	}
	err = json.Unmarshal(res, &user)
	if err != nil {
		return "", err
	}
	c.config.Set(config.User, user.Username)
	c.config.Cache().Flush()
	return user.Username, nil
}

func (c HTTPClient) GetPipelineDetails(projectID, pipelineID string) ([]byte, error) {
	resp, _, err := c.Get(fmt.Sprintf("/projects/%s/pipelines/%s", url.PathEscape(projectID), url.PathEscape(pipelineID)))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindProjectDetails searches for a project by its ID or its name,
// with the ID having precedence over the name and returns the
// raw JSON object as byte array.
func (c HTTPClient) FindProjectDetails(nameOrID string) ([]byte, error) {
	// first try to get the project by its cached ID
	if cachedID := c.config.Cache().Get("projects", nameOrID); cachedID != "" {
		resp, _, err := c.Get("/projects/" + url.PathEscape(cachedID))
		if err == nil {
			return resp, nil
		}
	}

	// then try to find the project by its ID
	resp, _, err := c.Get("/projects/" + url.PathEscape(nameOrID))
	if err == nil {
		return resp, nil
	}

	// now try to find the project by name as a last resort
	resp, _, err = c.Get("/users/${user}/projects/?search=" + url.QueryEscape(nameOrID))
	if err != nil {
		return nil, err
	}
	projects := make([]map[string]interface{}, 0)
	err = json.Unmarshal(resp, &projects)
	if err != nil {
		return nil, err
	}
	if len(projects) <= 0 {
		return nil, fmt.Errorf("Project '%s' not found", nameOrID)
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
func (c HTTPClient) FindProject(nameOrID string) (*Project, error) {
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

func (c HTTPClient) isLoggedIn() bool {
	return c.config != nil && c.config.Get(config.URL) != "" && c.config.Get(config.Token) != ""
}

func (c HTTPClient) Get(path string) ([]byte, int, error) {
	if !c.isLoggedIn() {
		return nil, 0, ErrNotLoggedIn
	}
	req, err := http.NewRequest("GET", c.config.Get(config.URL)+c.basePath+c.parse(path), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add("Private-Token", c.config.Get(config.Token))
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("%s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, 0, nil
}

func (c HTTPClient) Post(path string, reqBody io.Reader) ([]byte, int, error) {
	if !c.isLoggedIn() {
		return nil, 0, ErrNotLoggedIn
	}
	req, err := http.NewRequest("POST", c.config.Get(config.URL)+c.basePath+c.parse(path), reqBody)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add("Private-Token", c.config.Get(config.Token))
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp.StatusCode, fmt.Errorf("%s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, 0, nil
}

func (c HTTPClient) Delete(path string) (int, error) {
	if !c.isLoggedIn() {
		return 0, ErrNotLoggedIn
	}
	req, err := http.NewRequest("DELETE", c.config.Get(config.URL)+c.basePath+c.parse(path), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Private-Token", c.config.Get(config.Token))
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp.StatusCode, fmt.Errorf("%s", resp.Status)
	}
	return resp.StatusCode, nil
}
