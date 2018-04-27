package gitlab

import (
	gitlabapi "github.com/xanzy/go-gitlab"
	"sort"
)

type ListUsersOptions struct {
	ListOptions
	Active bool
}

func (c *Client) ListUsers(options ListUsersOptions) ([]*gitlabapi.User, error) {
	usersAll := make([]*gitlabapi.User, 0)

	listOptions := gitlabapi.ListUsersOptions{
		ListOptions: gitlabapi.ListOptions{
			PerPage: options.PageElements,
			Page:    options.PageNumber,
		},
		Active: &options.Active,
	}

	count := 0

	for {
		count++

		users, resp, err := c.Client.Users.ListUsers(&listOptions)
		if err != nil {
			return nil, err
		}

		usersAll = append(usersAll, users...)

		if count == options.PagesLimit {
			break
		}

		if resp.NextPage == 0 {
			break
		}

		listOptions.ListOptions.Page = resp.NextPage
	}

	sort.Sort(byUserName(usersAll))

	return usersAll, nil
}
