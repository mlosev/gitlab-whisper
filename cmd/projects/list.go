package projects

import (
	"github.com/mlosev/gitlab-whisper/internal/gitlab"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

type ListFlags struct {
	*RootFlags
	Archived bool
	All      bool
}

func NewListCommand(rootFlags *RootFlags) *cobra.Command {
	flags := ListFlags{RootFlags: rootFlags}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects, available to current user",
		Run: func(cmd *cobra.Command, args []string) {
			runListCommand(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.Archived, "with-archived", "", false, "Include archived projects")
	cmd.Flags().BoolVarP(&flags.All, "all", "a", false, "Include all projects with access")

	return cmd
}

func runListCommand(flags ListFlags) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID",
		"Namespace",
		"Name",
		"SSH Url",
	})

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
		table.Append([]string{
			strconv.FormatInt(int64(p.ID), 10),
			p.Namespace.Path,
			p.Path,
			p.SSHURLToRepo,
		})
	}

	table.Render()

	logrus.Infof("Total projects: %d", len(projects))
}
