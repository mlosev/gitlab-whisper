package gitlab

import (
	gitlabapi "github.com/xanzy/go-gitlab"
)

var gitlabClient *Client

type Client struct {
	gitlabapi.Client
}

func NewClient(endpoint, token string) (*Client, error) {
	client := gitlabapi.NewClient(nil, token)

	err := client.SetBaseURL(endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{*client}, nil
}

func Initialize(client *Client) {
	gitlabClient = client
}

func GetClient() *Client {
	if gitlabClient == nil {
		panic("gitlab client is not initialized")
	}
	return gitlabClient
}
