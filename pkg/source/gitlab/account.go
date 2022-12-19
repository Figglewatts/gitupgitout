package gitlab

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/xanzy/go-gitlab"
)

type Account struct {
	Url     *url.URL
	Account string

	client *gitlab.Client
}

const (
	TokenEnvVar = "GITLAB_TOKEN"
)

func (a *Account) ListRepos(ctx context.Context) ([]string, error) {
	err := a.lazyInitialise(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialising source: %w", err)
	}

	repos, _, err := a.client.Projects.ListUserProjects(
		a.Account,
		nil, gitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}

	var result []string
	for _, repo := range repos {
		result = append(result, repo.PathWithNamespace)
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

	repo, _, err := a.client.Projects.GetProject(
		repoName, nil, gitlab.WithContext(ctx),
	)
	if err != nil {
		return "", fmt.Errorf("getting project '%s': %w", repoName, err)
	}

	return repo.SSHURLToRepo, nil
}

func (a *Account) lazyInitialise(ctx context.Context) error {
	if a.client != nil {
		return nil
	}
	if a.Url == nil {
		a.Url = &url.URL{
			Scheme: "https",
			Host:   "gitlab.com",
			Path:   "api/v4/",
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

	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(a.Url.String()))
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}
	a.client = client

	return nil
}
