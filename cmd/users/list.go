package users

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
	All bool
}

func NewListCommand(rootFlags *RootFlags) *cobra.Command {
	flags := ListFlags{RootFlags: rootFlags}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(cmd *cobra.Command, args []string) {
			runListCommand(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.All, "all", "a", false, "List all users, not only active ones")

	return cmd
}

func runListCommand(flags ListFlags) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID",
		"Username",
		"Email",
		"State",
	})

	client := gitlab.GetClient()

	users, err := client.ListUsers(gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{
			PageElements: flags.PageElements,
			PageNumber:   flags.PageNumber,
			PagesLimit:   flags.PagesLimit,
		},
		Active: !flags.All,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Unable to list users")
	}

	for _, u := range users {
		table.Append([]string{
			strconv.FormatInt(int64(u.ID), 10),
			u.Username,
			u.Email,
			u.State,
		})
	}

	table.Render()

	logrus.Infof("Total users: %d", len(users))
}
