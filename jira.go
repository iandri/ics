package main

import (
	"context"
	"fmt"
	"ics/requests"
	"net/url"

	"github.com/palantir/stacktrace"
)

func GetTicket(ctx context.Context, basePath, username, password, ticket string) ([]byte, error) {
	loginUrl := fmt.Sprintf("%s/login.jsp", basePath)
	ticketUrl := fmt.Sprintf("%s/si/jira.issueviews:issue-xml/%s/%s.xml", basePath, ticket, ticket)
	client, err := requests.NewClient()
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	if _, err := client.Get(ctx, loginUrl); err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	values := url.Values{}
	values.Add("os_username", username)
	values.Add("os_password", password)
	values.Add("login", "Log+In")
	values.Encode()

	if _, err := client.Post(ctx, loginUrl, values,
		requests.WithHeader("Content-Type", "application/x-www-form-urlencoded")); err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	body, err := client.Get(ctx, ticketUrl)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	return body, nil

}
