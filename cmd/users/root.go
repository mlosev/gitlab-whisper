package users

import "github.com/spf13/cobra"

type RootFlags struct {
	PageNumber   int
	PageElements int
	PagesLimit   int
}

func NewRootCmd(flags *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Various actions with users",
	}

	cmd.PersistentFlags().IntVarP(&flags.PageNumber, "page", "p", 1, "Page number to request")
	cmd.PersistentFlags().IntVarP(&flags.PageElements, "size", "s", 10, "Number of elements on page")
	cmd.PersistentFlags().IntVarP(&flags.PagesLimit, "limit", "l", 1<<20, "How many pages to retrieve")

	return cmd
}

func AddCommands(cmd *cobra.Command) {
	flags := &RootFlags{}

	rootCmd := NewRootCmd(flags)

	rootCmd.AddCommand(NewListCommand(flags))

	cmd.AddCommand(rootCmd)
}
