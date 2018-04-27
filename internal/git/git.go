package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// Repo - simple repo struct
type Repo struct {
	localPath string
}

// NewRepo - return wrapper for git repo
func NewRepo(repoDir string) *Repo {
	return &Repo{localPath: repoDir}
}

func (r *Repo) execCommand(args []string, chdir bool) (string, string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("git", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if chdir {
		cmd.Dir = r.localPath
	}

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

// Path - return local path for the repo
func (r *Repo) Path() string {
	return r.localPath
}

// Clone - clone repo from given url
func (r *Repo) Clone(repoURL string) (string, error) {
	command := []string{"clone", repoURL, r.localPath}
	_, stderr, err := r.execCommand(command, false)

	return stderr, err

}

// Pull - pull updates for repo
func (r *Repo) Pull(remote, branch string) (string, error) {
	command := []string{"pull", remote}

	if branch == "" {
		out, stderr, err := r.getCurrentBranch()
		if err != nil {
			return stderr, err
		}
		branch = out
	}

	command = append(command, branch)

	_, stderr, err := r.execCommand(command, true)

	return stderr, err
}

func (r *Repo) getCurrentBranch() (string, string, error) {
	command := []string{"rev-parse", "--abbrev-ref", "HEAD"}
	out, stderr, err := r.execCommand(command, true)
	branch := strings.TrimSpace(out)
	return branch, stderr, err
}

// HasCommits - check if repo has any commits in it
func (r *Repo) HasCommits() (bool, error) {
	command := []string{"rev-list", "-n", "1", "--all"}

	stdout, _, err := r.execCommand(command, true)
	if err != nil {
		return false, err
	}

	if len(stdout) == 0 {
		return false, nil
	}

	return true, err
}

// Fetch - fetch updates from origin
func (r *Repo) Fetch(remote string) (string, error) {
	command := []string{"fetch", remote}

	_, stderr, err := r.execCommand(command, true)

	return stderr, err
}
