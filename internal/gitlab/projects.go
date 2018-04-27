package gitlab

import (
	gitlabapi "github.com/xanzy/go-gitlab"
	"sort"
)

type ListProjectsOptions struct {
	ListOptions
	Membership bool
	Archived   bool
}

func (c *Client) ListProjects(options ListProjectsOptions) ([]*gitlabapi.Project, error) {
	projectsAll := make([]*gitlabapi.Project, 0)

	listOptions := gitlabapi.ListProjectsOptions{
		ListOptions: gitlabapi.ListOptions{
			PerPage: options.PageElements,
			Page:    options.PageNumber,
		},
		Archived:   &options.Archived,
		Membership: &options.Membership,
	}

	count := 0

	for {
		count++

		projects, resp, err := c.Client.Projects.ListProjects(&listOptions)
		if err != nil {
			return nil, err
		}

		projectsAll = append(projectsAll, projects...)

		if count == options.PagesLimit {
			break
		}

		if resp.NextPage == 0 {
			break
		}

		listOptions.ListOptions.Page = resp.NextPage
	}

	sort.Sort(byProjectPath(projectsAll))

	return projectsAll, nil
}
