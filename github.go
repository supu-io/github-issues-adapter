package main

import (
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/supu-io/payload"
	"golang.org/x/oauth2"
)

func doMove(p payload.Payload, client *github.Client) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *p.Config.Github.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	if client == nil {
		client = github.NewClient(tc)
	}
	org, repo, number := getIssueDetails(*p.Issue.ID)

	for _, status := range *p.Status {
		_, err := client.Issues.RemoveLabelForIssue(org, repo, number, status)
		if err != nil {
			println(err.Error())
		}
	}

	status := p.Transition.To
	client.Issues.AddLabelsToIssue(org, repo, number, []string{*status})

	return nil
}

func getIssueDetails(id string) (string, string, int) {
	parts := strings.Split(id, "/")
	org := parts[0]
	repo := parts[1]
	number, _ := strconv.Atoi(parts[2])

	return org, repo, number
}
