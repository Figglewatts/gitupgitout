package git

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/figglewatts/gitupgitout/pkg/logging"
)

func Mirror(
	ctx context.Context, repoPath, repoName, sshUrl string, log logging.Logger,
) (string, error) {
	if err := ensurePathExists(repoPath); err != nil {
		return "", fmt.Errorf("creating git archive directory: %w", err)
	}

	log.Verbose(
		"cloning/fetching repo %s (%s) to %s", repoName, sshUrl, repoPath,
	)

	// see if the repo exists, as if it doesn't we'll need to clone it first
	repoExists, err := repoExists(repoPath, repoName)
	if err != nil {
		return "", fmt.Errorf(
			"checking if repo '%s/%s' existed: %w", repoPath, repoName, err,
		)
	}

	if !repoExists {
		// if the repo was nil we need to clone it as it hadn't already been
		log.Verbose("repo did not exist, cloning repo %s", sshUrl)
		err = mirrorClone(
			ctx, repoPath, getRepoDirectory(repoName), sshUrl,
		)
		if err != nil {
			return "", fmt.Errorf("cloning repo '%s': %w", sshUrl, err)
		}
	}

	// fetch the latest changes from the repo
	log.Verbose("fetching latest changes from repo %s", sshUrl)
	err = fetch(ctx, repoPath, getRepoDirectory(repoName))
	if err != nil {
		return "", fmt.Errorf("fetching git repo '%s': %w", sshUrl, err)
	}

	return getRepoPath(repoPath, repoName), nil
}

func repoExists(repoPath, repoName string) (bool, error) {
	if _, err := os.Stat(getRepoPath(repoPath, repoName)); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getRepoDirectory(repoName string) string {
	return fmt.Sprintf("%s.git", repoName)
}

func getRepoPath(repoPath, repoName string) string {
	return path.Join(repoPath, getRepoDirectory(repoName))
}

const (
	DirFileMode = 0755
)

func ensurePathExists(repoPath string) error {
	return os.MkdirAll(repoPath, DirFileMode)
}
