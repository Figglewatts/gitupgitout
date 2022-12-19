package git

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func IsInstalled() (bool, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsLfsInstalled() (bool, error) {
	cmd := exec.Command("git", "lfs")
	out, err := cmd.CombinedOutput()
	if err != nil {
		output := string(out)
		if strings.Contains(output, "not a git command") {
			return false, nil
		}
		return false, fmt.Errorf(
			"checking if lfs is installed: "+
				"output: %s\nerror: %w", output, err,
		)
	}
	return true, nil
}

func mirrorClone(
	ctx context.Context, workingDir, repoName, repoUrl string,
) error {
	return runGit(
		ctx, workingDir, "clone", "--mirror", repoUrl, repoName,
	)
}

func fetch(ctx context.Context, workingDir, repoName string) error {
	workingDir = path.Join(workingDir, repoName)
	err := runGit(ctx, workingDir, "fetch", "--all")
	if err != nil {
		return fmt.Errorf("fetching: %w", err)
	}

	err = runGit(ctx, workingDir, "lfs", "install")
	if err != nil {
		return fmt.Errorf("running lfs install: %w", err)
	}
	return nil
}

func runGit(ctx context.Context, workingDir string, arg ...string) error {
	cmd := exec.CommandContext(ctx, "git", arg...)
	cmd.Dir = workingDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		commandRun := fmt.Sprintf("git %v", arg)
		return fmt.Errorf("running %s\noutput:%s: %w", commandRun, out, err)
	}
	return nil
}
