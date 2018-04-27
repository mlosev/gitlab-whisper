package gitlab

import (
	gitlabapi "github.com/xanzy/go-gitlab"
	"strings"
)

// projects

type byProjectID []*gitlabapi.Project

func (s byProjectID) Len() int           { return len(s) }
func (s byProjectID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byProjectID) Less(i, j int) bool { return s[i].ID < s[j].ID }

type byProjectPath []*gitlabapi.Project

func (s byProjectPath) Len() int      { return len(s) }
func (s byProjectPath) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byProjectPath) Less(i, j int) bool {
	return strings.Compare(s[i].PathWithNamespace, s[j].PathWithNamespace) == -1
}

// users

type byUserID []*gitlabapi.User

func (s byUserID) Len() int           { return len(s) }
func (s byUserID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byUserID) Less(i, j int) bool { return s[i].ID < s[j].ID }

type byUserName []*gitlabapi.User

func (s byUserName) Len() int           { return len(s) }
func (s byUserName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byUserName) Less(i, j int) bool { return strings.Compare(s[i].Username, s[j].Username) == -1 }
