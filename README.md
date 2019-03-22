# Summary
This is a small tool with you can:
1) Setup your GitHub authentication
2) Transfer one or a bunch of repositories (one per line from file) for the selected organisation
3) Update transferred repositories to replace all url chunks that remained from the original organisation, also sends in the PR with the changes.


# Manual - after-transfer - todos

- When the repo transferred the Bitrise apps will remain wired in but their git url still will point to the original URL. (In the app's settings tab)
- The next dep update might requires to re-create the dep descriptor files

# Storage model

> $ cat .storage.json

```json
{
 "GithubAccessToken": "...",
 "TransferredRepositories": [
  {
   "OriginalURL": "https://github.com/bitrise-samples/sample-repo1",
   "TargetOrg": "trapacska"
  }
 ],
 "UpdatedRepositories": [
  "https://github.com/trapacska/sample-repo1"
 ]
}
```

# Commands

> $ repotool -h

```
This tool helps in repository operations.

Usage:
  repotool [command]

Available Commands:
  auth        Set GitHub authentication for this tool.
  help        Help about any command
  transfer    Transfers repository to another organization.
  update      Update repositories to follow repo URL change.

Flags:
  -h, --help   help for repotool

Use "repotool [command] --help" for more information about a command.
```

> $ repotool auth -h

```
Set GitHub authentication for this tool.

Usage:
  repotool auth "github-access-token" [flags]

Flags:
  -h, --help   help for auth
```

> $ repotool transfer -h

```
Transfers repository to another organization.

Usage:
  repotool transfer [flags]

Flags:
  -h, --help                 help for transfer
  -f, --source-file string   Source file which contains one repository url per line
  -r, --source-repo string   Source repository url
  -o, --target-org string    Target organization
```

> $ repotool update -h

```
Update repositories to follow repo URL change.

Usage:
  repotool update [flags]

Flags:
  -h, --help   help for update
```