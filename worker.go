package main

import (
	"log"
	"strings"

	"github.com/google/go-github/github"
)

type Repository struct {
	Owner string
	Name  string
	SHA   string
	Files []string
}

var (
	client *github.Client
	err    error
)

func main() {
	client = github.NewClient(nil)

	repository, err := getRepository("citruspi", "Spotify-Notifications")

	if err != nil {
		log.Fatal(err)
	}
}

func getRepository(owner string, name string) (Repository, error) {
	var repository Repository

	repository.Owner = owner
	repository.Name = name

	latestSHA, err := getLatestSHA(owner, name)

	if err != nil {
		return repository, err
	}

	repository.SHA = latestSHA

	tree, err := getTree(owner, name, latestSHA)

	if err != nil {
		return repository, err
	}

	for _, entry := range tree.Entries {
		repository.Files = append(repository.Files, stringify(entry.Path))
	}

	return repository, nil
}

func getLatestSHA(owner string, repository string) (string, error) {
	commits, _, err := client.Repositories.ListCommits(owner, repository, nil)

	if err != nil {
		return "", err
	}

	latestSHA := stringify(commits[0].SHA)

	return latestSHA, nil
}

func getTree(owner string, repository string, sha string) (*github.Tree, error) {
	tree, _, err := client.Git.GetTree(owner, repository, sha, true)

	if err != nil {
		return nil, err
	}

	return tree, nil
}

func stringify(message interface{}) string {
	var str string

	str = github.Stringify(message)
	str = strings.Replace(str, "\"", "", 2)

	return str
}
