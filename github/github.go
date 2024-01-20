// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
package github

import (
	"aczietlow/talamasca/config"
	"context"
	"fmt"
	"github.com/google/go-github/v58/github"
	"time"
)

var client *github.Client
var conf *config.Config

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
	ItemsMap   map[int]*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type CommitData struct {
	Sha            string
	GithubAuthor   string
	GitAuthor      string
	Message        string
	ReleaseTagName string
}

func Setup(c *config.Config) {
	// Setup returns a new GitHub API client. If a nil httpClient is
	// provided, a new http.Client will be used. To use API methods which require
	// authentication, provide an http.Client that will perform the authentication
	// for you (such as that provided by the golang.org/x/oauth2 library).

	conf = c
	client = github.NewClient(nil).WithAuthToken(conf.Token)
}

func FetchRepoReleases() ([]*github.RepositoryRelease, error) {
	releases, _, err := client.Repositories.ListReleases(context.Background(), conf.Owner, conf.Repo, nil)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func FetchCommitsFromMostRecentRelease() (*github.CommitsComparison, error) {
	releases, err := FetchRepoReleases()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	//fmt.Printf("%v\n", releases[1].GetTagName())

	previousReleasedTag := releases[2].GetTagName()
	currentReleasedTag := releases[1].GetTagName()

	commits, _, err := client.Repositories.CompareCommits(context.Background(), conf.Owner, conf.Repo, previousReleasedTag, currentReleasedTag, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	return commits, err
}

// GetCommitBySha Fetches specific commit sha and prints the ghAuthor and gitAuthor (WHICH ARE DIFFERENT)
func GetCommitBySha(sha string) {
	commit, _, err := client.Repositories.GetCommit(context.Background(), conf.Owner, conf.Repo, sha, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ghAuthor := commit.GetAuthor()
	fmt.Printf("Github User: Name-%v\nID-%v\nEmail-%v\n", ghAuthor.GetLogin(), commit.GetAuthor().GetID(), commit.GetAuthor().GetEmail())

	gitAuthor := commit.GetCommit().GetAuthor()
	fmt.Printf("Git User: Login-%v\nName-%v\nEmail-%v\n", gitAuthor.GetLogin(), gitAuthor.GetName(), gitAuthor.GetEmail())
}

func GetCommitsContributors(commits *github.CommitsComparison, releaseTagName string) []CommitData {
	var commitContributions []CommitData

	for i, commit := range commits.Commits {
		fmt.Printf("%v. %v - %v\n", i+1, commit.GetSHA(), commit.GetCommit().GetMessage())
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
		cd := CommitData{sha, githubAuthor, gitAuthor, message, releaseTagName}
		commitContributions = append(commitContributions, cd)
	}
	return commitContributions
}

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
