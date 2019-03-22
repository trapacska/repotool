package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trapacska/repotool/storage"
)

var authCmd = &cobra.Command{
	Use:   "auth \"github-access-token\"",
	Short: "Set GitHub authentication for this tool.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return storage.Update(func(table *storage.Table) {
			table.GithubAccessToken = args[0]
		})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid count of arguments")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
