package main

import (
	"aczietlow/talamasca/config"
	"aczietlow/talamasca/github"
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	// @TODO deal with errors later when I actually care.
	conf, _ := config.LoadConfig("./config.json")

	github.Setup(conf)

	// Just making in clear that this is PoC and needs to be refactored.
	// Or my code identifies as Spaghetti.
	spaghetti(conf)
}

func spaghetti(conf *config.Config) {

	releases, _ := github.FetchRepoReleases()
	commits, _ := github.FetchCommitsFromMostRecentRelease()
	commitContributions := github.GetCommitsContributors(commits, releases[0].GetTagName())

	writeToFile(conf.Repo+".csv", commitContributions)
}

func writeToFile(fileName string, data []github.CommitData) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Error opening the file.")
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err = writer.Write([]string{value.Sha, value.GithubAuthor, value.GitAuthor, value.Message, value.ReleaseTagName})
		if err != nil {
			fmt.Println("Error writing to CSV:", err)
			return
		}
	}

}
