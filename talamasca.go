package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v53/github"
	"os"
)

type Config struct {
	Repo  string `json:"repo"`
	Owner string `json:"owner"`
}

// Fetch all releases for a given repository
func fetchRepoReleases(owner string, reponame string) ([]*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), owner, reponame, nil)

	return releases, err
}

func loadConfig(filename string) (*Config, error) {
	configFile, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func main() {

	// @TODO deal with errors later when I actually care.
	config, _ := loadConfig("./config.json")

	releases, err := fetchRepoReleases(config.Owner, config.Repo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%v\n", releases[1].GetTagName())

	last := releases[1].GetTagName()
	current := releases[0].GetTagName()

	// Get commits from a tag release
	client := github.NewClient(nil)
	//opt := &github.CommitsListOptions{SHA: "tags/8.1.3"}
	//commits, _, err := client.Repositories.ListCommits(context.Background(), "pantheon-systems", "search_api_pantheon", opt)
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}
	//
	//for i, commit := range commits {
	//	fmt.Printf("%v. %v - %v\n", i+1, commit.GetSHA(), commit.GetCommit().GetMessage())
	//}

	//for i, release := range releases {
	//	fmt.Printf("%v. %v\n", i+1, *release.Body)
	//}

	// Get repository.
	commits, _, err := client.Repositories.CompareCommits(context.Background(), config.Owner, config.Repo, last, current, nil)

	// Print commits
	for i, commit := range commits.Commits {
		fmt.Printf("%v. %v - %v\n", i+1, commit.GetSHA(), commit.GetCommit().GetMessage())
	}

	commit, _, err := client.Repositories.GetCommit(context.Background(), config.Owner, config.Repo, "aa88ed97ad7270d83bc3425fbe9bbe401c7c41f4", nil)

	ghAuthor := commit.GetAuthor()
	fmt.Printf("Github User: Name-%v\nID-%v\nEmail-%v\n", ghAuthor.GetLogin(), commit.GetAuthor().GetID(), commit.GetAuthor().GetEmail())

	gitAuthor := commit.GetCommit().GetAuthor()
	fmt.Printf("Git User: Login-%v\nName-%v\nEmail-%v\n", gitAuthor.GetLogin(), gitAuthor.GetName(), gitAuthor.GetEmail())

}
