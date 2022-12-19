package github

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type Account struct {
	Url     *url.URL
	Account string

	client *github.Client
}

const (
	TokenEnvVar = "GITHUB_TOKEN"
)

func (a *Account) ListRepos(ctx context.Context) ([]string, error) {
	err := a.lazyInitialise(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialising source: %w", err)
	}

	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 30}}
	for {
		repos, resp, err := a.client.Repositories.List(ctx, a.Account, opt)
		if err != nil {
			return nil, fmt.Errorf("listing repos: %w", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var result []string
	for _, repo := range allRepos {
		result = append(result, *repo.Name)
	}
	return result, nil
}

func (a *Account) GetRepo(ctx context.Context, repoName string) (
	string, error,
) {
	err := a.lazyInitialise(ctx)
	if err != nil {
		return "", fmt.Errorf("initialising source: %w", err)
	}

	repo, _, err := a.client.Repositories.Get(ctx, a.Account, repoName)
	if err != nil {
		return "", fmt.Errorf("getting repo '%s': %w", repoName, err)
	}

	return repo.GetSSHURL(), nil
}

func (a *Account) lazyInitialise(ctx context.Context) error {
	if a.client != nil {
		return nil
	}
	if a.Url == nil {
		a.Url = &url.URL{
			Scheme: "https",
			Host:   "api.github.com",
		}
	}
	if a.Account == "" {
		return fmt.Errorf("account is empty")
	}

	token, ok := os.LookupEnv(TokenEnvVar)
	if !ok {
		return fmt.Errorf(
			"environment variable %s is required", TokenEnvVar,
		)
	}

	tc := oauth2.NewClient(
		ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	)
	client, err := github.NewEnterpriseClient(
		a.Url.String(), a.Url.String(), tc,
	)
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}
	a.client = client

	return nil
}
