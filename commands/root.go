package commands

import (
	"io/ioutil"
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

func parseRepoList(file, repo string) (repos []string, err error) {
	if s := strings.TrimSpace(repo); len(s) > 0 {
		repos = append(repos, s)
		return
	}

	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	for _, fileLine := range strings.Split(string(fileData), "\n") {
		if s := strings.TrimSpace(fileLine); len(s) > 0 {
			repos = append(repos, s)
		}
	}
	return repos, nil
}
