package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mlosev/gitlab-whisper/cmd/projects"
	"github.com/mlosev/gitlab-whisper/cmd/users"
	"github.com/mlosev/gitlab-whisper/internal/gitlab"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitlab-whisper",
		Short: "Utility for dealing with gitlab projects/users/etc.",
		Long: `Sometimes we need to all projects from gitlab instance
                we have access to local machine
                This tool will help you with that`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			_ = cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initGitlabClient()
		},
	}

	return cmd
}

func AddCommands(cmd *cobra.Command) {
	projects.AddCommands(cmd)
	users.AddCommands(cmd)
}

func Execute() {
	cmd := NewRootCommand()
	AddCommands(cmd)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initGitlabClient() {
	_ = godotenv.Load("gitlab.env")

	token := os.Getenv("GITLAB_API_PRIVATE_TOKEN")
	endpoint := os.Getenv("GITLAB_API_ENDPOINT")

	client, err := gitlab.NewClient(endpoint, token)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to setup gitlab client")
	}

	gitlab.Initialize(client)
}
