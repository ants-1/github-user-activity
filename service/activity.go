package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GitHubEvent struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
	Payload   struct {
		Size int `json:"size"`
	} `json:"payload"`
}

func GetActivity(name string) ([]GitHubEvent, error) {

	baseUrl := `https://api.github.com/users/%s/events`
	url := fmt.Sprintf(baseUrl, name)

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "github-user-activity")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var activities []GitHubEvent
	err = json.Unmarshal(body, &activities)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return activities, nil
}

func FormatEvent(e GitHubEvent) string {
	switch e.Type {
	case "PushEvent":
		return fmt.Sprintf("Pushed %d commits to %s", e.Payload.Size, e.Repo.Name)
	case "IssuesEvent":
		return fmt.Sprintf("Opened a new issue in %s", e.Repo.Name)
	case "WatchEvent":
		return fmt.Sprintf("Starred %s", e.Repo.Name)
	case "ForkEvent":
		return fmt.Sprintf("Forked %s", e.Repo.Name)
	case "PullRequestEvent":
		return fmt.Sprintf("Opened a pull request in %s", e.Repo.Name)
	default:
		return fmt.Sprintf("%s on %s", e.Type, e.Repo.Name)
	}
}
