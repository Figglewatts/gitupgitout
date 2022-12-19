package source

import (
	"context"
)

type RepoSource interface {
	// ListRepos lists the repos available from this source
	ListRepos(ctx context.Context) ([]string, error)

	// GetRepo will fetch the repo where we can mirror the URL from
	GetRepo(ctx context.Context, repo string) (string, error)
}
