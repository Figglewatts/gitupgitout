# gitupgitout (gugo)
A program used to mirror Git-related things to a local filesystem.

## Usage
### Prerequisites
- Git
- Git LFS

Install via:
```
go install github.com/figglewatts/gitupgitout/cmd/gugo
```

Or download a binary from the releases page.

Afterwards, you can run the application like this:
```
$ gugo --help
```

```
Usage of gugo:
  -concurrency int
    	number of mirrors to process at once
  -config string
    	path to config file (default "gugo.yaml")
  -cron string
    	run on a cron schedule
  -run-before-cron
    	run anyway before first scheduled cron
  -verbose
    	enable verbose logging
```

## Config
A config file for `gugo` looks like this:
```yaml
mirrors:
  - source:
      githubAccount:
        account: figglewatts
    cloneTo: ./repos
  - source:
      gitlabAccount:
        account: someone_else
    cloneTo: ./repos
```
You can define multiple mirrors of a different variety of sources, specifying
where repositories will be cloned to.

Each repository clone operation is processed in parallel (control how parallel with
the `-concurrency` flag). If a repository has already been cloned, on subsequent runs
it will only be fetched from (via `git fetch --all`).

Each mirror can only have one source. Sources are outlined in the below section.

### Sources
#### GitHub Account
```yaml
githubAccount:
  account: <account_name>
  url: <url_to_github_api>  # optional
```
- `GITHUB_TOKEN` environment variable required for authentication.

#### GitLab Account
```yaml
gitlabAccount:
  account: <account_name>
  url: <url_to_gitlab_api>  # optional
```
- `GITLAB_TOKEN` environment variable required for authentication.