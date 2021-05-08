package api

import "time"

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

type Job struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Stage     string `json:"stage"`
	Status    string `json:"status"`
}

type Jobs []Job

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
