package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type ReleaseInfoEr interface {
	GetLatestReleaseTag(string) (string, error)
}

type GitHubReleaseInfoEr struct{}

func (gh GitHubReleaseInfoEr) GetLatestReleaseTag(repo string) (string, error) {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)
	response, err := http.Get(apiUrl)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	releases := []ReleaseInfo{}

	if err := json.Unmarshal(body, &releases); err != nil {
		return "", err
	}

	tag := releases[0].TagName

	return tag, nil
}

type ReleaseInfo struct {
	ID      uint   `json"id"`
	TagName string `json:"tag_name"`
}

func getReleaseTagMessage(ri ReleaseInfoEr, repo string) (string, error) {
	tag, err := ri.GetLatestReleaseTag(repo)
	if err != nil {
		return "", fmt.Errorf("Error querying GitHub API: %s", err)
	}

	return fmt.Sprintf("The latest release is %s", tag), nil
}

func main() {
	gh := GitHubReleaseInfoEr{}
	msg, err := getReleaseTagMessage(gh, "docker/machine")
	if err != nil {
		fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}

	fmt.Println(msg)
}
