package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v53/github"
	"os"
)

type Config struct {
	Repo  string `json:"repo"`
	Owner string `json:"owner"`
}

var client *github.Client

func initClient(config *Config) {
	// NewClient returns a new GitHub API client. If a nil httpClient is
	// provided, a new http.Client will be used. To use API methods which require
	// authentication, provide an http.Client that will perform the authentication
	// for you (such as that provided by the golang.org/x/oauth2 library).

	// @todo add auth if present in config file later.
	client = github.NewClient(nil)
	fmt.Println(client.UserAgent)
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

	// Just making in clear that this is PoC and needs to be refactored.
	// Or my code identifies as Spaghetti.
	spaghetti()
}

func spaghetti() {
	// @TODO deal with errors later when I actually care.
	config, _ := loadConfig("./config.json")

	initClient(config)

	releases, _ := fetchRepoReleases(config.Owner, config.Repo)
	commits, _ := fetchCommitsFromMostRecentRelease(config)
	var commitContributorions [][]string

	// Get commits from a tag release
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

	// Print commits
	for i, commit := range commits.Commits {
		fmt.Printf("%v. %v - %v\n", i+1, commit.GetSHA(), commit.GetCommit().GetMessage())
		release := releases[0].GetTagName()
		var gitAuthor string
		var githubAuthor string
		gitUser := commit.GetCommit().GetAuthor()
		if gitUser != nil {
			gitAuthor = gitUser.GetEmail()
		} else {
			gitAuthor = ""
		}
		githubUser := commit.GetAuthor()
		if githubUser != nil {
			githubAuthor = githubUser.GetLogin()
		} else {
			githubAuthor = ""
		}
		sha := commit.GetSHA()
		message := commit.GetCommit().GetMessage()
		commitData := []string{sha, githubAuthor, gitAuthor, message, release}
		commitContributorions = append(commitContributorions, commitData)
	}

	writeToFile(config.Repo+".csv", commitContributorions)
}

// Fetches specific commit sha and prints the ghAuthor and gitAuthor (WHICH ARE DIFFERENT)
func getCommit(config *Config, sha string) {
	commit, _, err := client.Repositories.GetCommit(context.Background(), config.Owner, config.Repo, sha, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ghAuthor := commit.GetAuthor()
	fmt.Printf("Github User: Name-%v\nID-%v\nEmail-%v\n", ghAuthor.GetLogin(), commit.GetAuthor().GetID(), commit.GetAuthor().GetEmail())

	gitAuthor := commit.GetCommit().GetAuthor()
	fmt.Printf("Git User: Login-%v\nName-%v\nEmail-%v\n", gitAuthor.GetLogin(), gitAuthor.GetName(), gitAuthor.GetEmail())
}

// Fetch all releases for a given repository
func fetchRepoReleases(owner string, reponame string) ([]*github.RepositoryRelease, error) {
	releases, _, err := client.Repositories.ListReleases(context.Background(), owner, reponame, nil)
	return releases, err
}

func fetchCommitsFromMostRecentRelease(config *Config) (*github.CommitsComparison, error) {
	releases, err := fetchRepoReleases(config.Owner, config.Repo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	//fmt.Printf("%v\n", releases[1].GetTagName())

	previousReleasedTag := releases[2].GetTagName()
	currentReleasedTag := releases[1].GetTagName()

	commits, _, err := client.Repositories.CompareCommits(context.Background(), config.Owner, config.Repo, previousReleasedTag, currentReleasedTag, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	return commits, err
}

func writeToFile(fileName string, data [][]string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Error opening the file.")
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			fmt.Println("Error writing to CSV:", err)
			return
		}
	}

}
