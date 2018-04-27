package projects

import (
	"github.com/mlosev/gitlab-whisper/internal/git"
	"github.com/mlosev/gitlab-whisper/internal/gitlab"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	gitlabapi "github.com/xanzy/go-gitlab"
	"net/url"
	"os"
	"path/filepath"
)

type SyncFlags struct {
	*RootFlags
	Archived    bool
	All         bool
	Destination string
	Run         bool
}

func NewSyncCommand(rootFlags *RootFlags) *cobra.Command {
	flags := SyncFlags{RootFlags: rootFlags}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync all projects, available to current user",
		Run: func(cmd *cobra.Command, args []string) {
			runSyncCommand(flags)
		},
	}

	cmd.Flags().StringVarP(&flags.Destination, "destination", "d", os.Getenv("GOPATH"), "Destination directory for all projects to sync to")
	cmd.Flags().BoolVarP(&flags.Archived, "with-archived", "", false, "Include archived projects")
	cmd.Flags().BoolVarP(&flags.All, "all", "a", false, "Include all projects with access")
	cmd.Flags().BoolVarP(&flags.Run, "run", "r", false, "Perform real synchronization instead of dry run")

	return cmd
}

func runSyncCommand(flags SyncFlags) {
	logger := logrus.WithField("destination", flags.Destination)

	_, err := os.Stat(flags.Destination)
	if err != nil {
		logger.WithError(err).Fatal("Unable to access destination folder")
	}

	baseDir := filepath.Join(flags.Destination, "src")

	client := gitlab.GetClient()

	projects, err := client.ListProjects(gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PageElements: flags.PageElements,
			PageNumber:   flags.PageNumber,
			PagesLimit:   flags.PagesLimit,
		},
		Membership: !flags.All,
		Archived:   flags.Archived,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Unable to list projects")
	}

	for _, p := range projects {
		logger := logger.WithFields(logrus.Fields{
			"namespace": p.Namespace.Name,
			"project":   p.Name,
		})

		path, err := prepareProjectPath(p, baseDir)
		if err != nil {
			logger.WithError(err).Fatal("Unable to prepare project path")
		}

		repo := git.NewRepo(path)

		var errP error
		_, err = os.Stat(path)
		switch err {
		case nil:
			errP = updateProject(p, repo, !flags.Run)
		default:
			errP = cloneProject(p, repo, !flags.Run)
		}

		if errP != nil {
			logger.WithError(errP).Fatal("Unable to clone or update project")
		}
	}
}

func prepareProjectPath(project *gitlabapi.Project, baseDir string) (string, error) {
	u, err := url.Parse(os.Getenv("GITLAB_API_ENDPOINT"))
	if err != nil {
		return "", err
	}

	parentPath := filepath.Join(baseDir, u.Host, project.Namespace.Path)

	err = os.MkdirAll(parentPath, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(parentPath, project.Path), nil
}

func cloneProject(project *gitlabapi.Project, repo *git.Repo, dryRun bool) error {
	logger := logrus.WithField("projectPath", project.PathWithNamespace)
	logger.Info("Cloning project...")

	if dryRun {
		logger.Info("[DRY RUN] Cloned")
		return nil
	}

	stderr, err := repo.Clone(project.SSHURLToRepo)
	if err != nil {
		return errors.WithMessage(err, stderr)
	}

	return nil
}

func updateProject(project *gitlabapi.Project, repo *git.Repo, dryRun bool) error {
	logger := logrus.WithField("projectPath", project.PathWithNamespace)
	logger.Info("Updating project...")

	if dryRun {
		logger.Info("[DRY RUN] Updated")
		return nil
	}

	stderr, err := repo.Fetch("origin")
	if err != nil {
		return errors.WithMessage(err, stderr)
	}

	hasCommits, err := repo.HasCommits()
	if err != nil {
		return errors.WithMessage(err, "Unable to check commits in repo")
	}

	if !hasCommits {
		logger.Warn("Repo has not commits, nothing to update")
		return nil
	}

	stderr, err = repo.Pull("origin", "")
	if err != nil {
		return errors.WithMessage(err, stderr)
	}

	return nil
}
