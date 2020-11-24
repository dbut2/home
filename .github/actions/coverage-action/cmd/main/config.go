package main

import (
	"encoding/json"
	"strings"

	"github.com/google/go-github/github"
)

type githubContext struct {
	RunID string `json:"run_id"`
	Event struct {
		PullRequest struct {
			Number int `json:"number"`
		} `json:"pull_request"`
	} `json:"event"`
}

// EnvDecode is
func (g *githubContext) EnvDecode(val string) error {
	if err := json.Unmarshal([]byte(val), g); err != nil {
		return err
	}
	return nil
}

type config struct {
	client               *github.Client
	Repository           string        `env:"GITHUB_REPOSITORY,required"`
	GitHubToken          string        `env:"INPUT_TOKEN,required"`
	HardTarget           float64       `env:"INPUT_HARD_TARGET,required"`
	SoftTarget           float64       `env:"INPUT_SOFT_TARGET"`
	EnableWarningComment bool          `env:"INPUT_ENABLE_SOFT_TARGET_WARNING"`
	GitHubContext        githubContext `env:"GITHUB_CONTEXT"`
}

func (c config) RepoOwner() string {
	return strings.Split(c.Repository, "/")[0]
}

func (c config) RepoName() string {
	return strings.Split(c.Repository, "/")[1]
}

func (c config) PullRequestNumber() int {
	return c.GitHubContext.Event.PullRequest.Number
}

func (c config) RunID() string {
	return cfg.GitHubContext.RunID
}
