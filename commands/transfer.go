package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/bitrise-io/go-utils/log"
	"github.com/spf13/cobra"
	"github.com/trapacska/repotool/storage"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfers repository to another organization.",
	RunE: func(cmd *cobra.Command, args []string) error {
		repos, err := parseRepoList(sourceFile, sourceRepo)
		if err != nil {
			return err
		}
		return transferRepos(repos, targetOrg)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(sourceFile) > 0 && len(sourceRepo) > 0 {
			return fmt.Errorf("both --source-file and --source-repo cannot be set")
		}

		var authenticated bool
		if err := storage.Read(func(table *storage.Table) {
			authenticated = len(table.GithubAccessToken) > 0
		}); err != nil {
			return err
		}

		if !authenticated {
			return fmt.Errorf("not authenticated")
		}

		return nil
	},
}

func transferRepos(repos []string, org string) error {
	var authToken string
	var transferredRepos []storage.Repo

	if err := storage.Read(func(table *storage.Table) {
		authToken = table.GithubAccessToken
		transferredRepos = table.TransferredRepositories
	}); err != nil {
		return err
	}

main:
	for _, repoURL := range repos {
		for _, repo := range transferredRepos {
			if repoURL == repo.OriginalURL {
				fmt.Println("Skip (already transferred):", repoURL)
				continue main
			}
		}

		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(map[string]interface{}{"new_owner": org}); err != nil {
			log.Errorf("%s", err)
			continue
		}

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/transfer", repoOwner(repoURL), repoName(repoURL))

		req, err := http.NewRequest(http.MethodPost, url, &b)
		if err != nil {
			log.Errorf("%s", err)
			continue
		}

		req.Header.Add("Accept", "application/vnd.github.nightshade-preview+json")
		req.Header.Add("Authorization", "token "+authToken)

		if dumped, err := httputil.DumpRequest(req, true); err == nil {
			fmt.Println(string(dumped))
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if dumped, err := httputil.DumpResponse(resp, true); err == nil {
				fmt.Println(string(dumped))
			}
			log.Errorf("%s", err)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode > 210 {
			if dumped, err := httputil.DumpResponse(resp, true); err == nil {
				fmt.Println(string(dumped))
			}
			log.Errorf("%s", fmt.Errorf("non successful status code"))
			continue
		}

		fmt.Println("OK")
		fmt.Println()

		if err := storage.Update(func(table *storage.Table) {
			table.TransferredRepositories = append(table.TransferredRepositories, storage.Repo{OriginalURL: repoURL, TargetOrg: org})
		}); err != nil {
			log.Errorf("%s", err)
			continue
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().StringVarP(&targetOrg, "target-org", "o", "", "Target organization")
	transferCmd.Flags().StringVarP(&sourceRepo, "source-repo", "r", "", "Source repository url")
	transferCmd.Flags().StringVarP(&sourceFile, "source-file", "f", "", "Source file which contains one repository url per line")
	transferCmd.MarkFlagRequired("target-org")
}
