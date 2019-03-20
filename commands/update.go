package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/pkg/errors"

	"github.com/bitrise-io/go-utils/sliceutil"

	"github.com/spf13/cobra"
	"github.com/trapacska/repotool/storage"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update repositories to follow repo URL change.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateRepos()
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

func updateRepos() error {
	var authToken string
	var updatedRepos []string
	var transferredRepos []storage.Repo

	if err := storage.Read(func(table *storage.Table) {
		authToken = table.GithubAccessToken
		updatedRepos = table.UpdatedRepositories
		transferredRepos = table.TransferredRepositories
	}); err != nil {
		return err
	}

	for _, repo := range transferredRepos {
		from := repoOwner(repo.OriginalURL) + "/" + repoName(repo.OriginalURL)
		to := repo.TargetOrg + "/" + repoName(repo.OriginalURL)
		newURL := strings.Replace(repo.OriginalURL, from, to, -1)

		if sliceutil.IsStringInSlice(newURL, updatedRepos) {
			fmt.Println("Skip (already updated):", newURL)
			continue
		}

		tempDir, err := pathutil.NormalizedOSTempDirPath("repo")
		if err != nil {
			log.Errorf("%s", err)
			continue
		}

		if out, err := command.New("git", "clone", newURL, tempDir).RunAndReturnTrimmedCombinedOutput(); err != nil {
			log.Errorf("%s", errors.Wrap(err, out))
			continue
		}

		if out, err := command.New("git", "branch", "git-url-update").SetDir(tempDir).RunAndReturnTrimmedCombinedOutput(); err != nil {
			log.Errorf("%s", errors.Wrap(err, out))
			continue
		}

		if out, err := command.New("git", "checkout", "git-url-update").SetDir(tempDir).RunAndReturnTrimmedCombinedOutput(); err != nil {
			log.Errorf("%s", errors.Wrap(err, out))
			continue
		}

		if err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			c = []byte(strings.Replace(string(c), from, to, -1))

			return ioutil.WriteFile(path, c, info.Mode())
		}); err != nil {
			log.Errorf("%s", err)
			continue
		}

		if out, err := command.New("git", "add", ".").SetDir(tempDir).RunAndReturnTrimmedCombinedOutput(); err != nil {
			log.Errorf("%s", errors.Wrap(err, out))
			continue
		}

		if out, err := command.New("git", "commit", "-m", "updated git urls").SetDir(tempDir).RunAndReturnTrimmedCombinedOutput(); err != nil {
			log.Errorf("%s", errors.Wrap(err, out))
			continue
		}

		if out, err := command.New("git", "push", "origin", "git-url-update").SetDir(tempDir).RunAndReturnTrimmedCombinedOutput(); err != nil {
			log.Errorf("%s", errors.Wrap(err, out))
			continue
		}

		// opening PR
		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(map[string]interface{}{
			"title": "Update repository URLs",
			"head":  "git-url-update",
			"base":  "master",
		}); err != nil {
			log.Errorf("%s", err)
			continue
		}

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", repoOwner(newURL), repoName(newURL))

		req, err := http.NewRequest(http.MethodPost, url, &b)
		if err != nil {
			log.Errorf("%s", err)
			continue
		}

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

		fmt.Println("OK", tempDir)
		fmt.Println()

		if err := storage.Update(func(table *storage.Table) {
			table.UpdatedRepositories = append(table.UpdatedRepositories, newURL)
		}); err != nil {
			log.Errorf("%s", err)
			continue
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
