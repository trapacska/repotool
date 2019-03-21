package commands

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	targetOrg  string
	sourceFile string
	sourceRepo string
	rootCmd    = &cobra.Command{
		Use:   "repotool",
		Short: "This tool helps in repository operations.",
	}
)

// Execute ...
func Execute() error {
	return rootCmd.Execute()
}

func repoOwner(repoURL string) string {
	s := filepath.Base(filepath.Dir(repoURL))
	if i := strings.LastIndex(s, ":"); i > -1 {
		return s[i+1:]
	}
	return s
}
func repoName(repoURL string) string {
	return strings.TrimSuffix(filepath.Base(repoURL), filepath.Ext(filepath.Base(repoURL)))
}
