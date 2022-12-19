package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/figglewatts/gitupgitout/internal/config"
	"github.com/figglewatts/gitupgitout/pkg/git"
	"github.com/figglewatts/gitupgitout/pkg/logging"
	"github.com/figglewatts/gitupgitout/pkg/source"
)

type App struct {
	conf config.Config
	log  logging.Logger

	concurrency int
}

func New(configReader io.Reader, logger logging.Logger, opts ...Option) (
	*App, error,
) {
	conf, err := config.Load(configReader)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	ok, err := git.IsInstalled()
	if err != nil {
		return nil, fmt.Errorf("checking if git was installed: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("git must be installed and on your PATH")
	}

	ok, err = git.IsLfsInstalled()
	if err != nil {
		return nil, fmt.Errorf("checking if git lfs was installed: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("git lfs must be installed")
	}

	app := App{
		conf: *conf,
		log:  logger,
	}
	app.configure(opts...)

	if app.concurrency <= 0 {
		app.concurrency = runtime.NumCPU()
	}

	return &app, nil
}

func (app App) Run(ctx context.Context) error {
	if err := app.runOnce(ctx); err != nil {
		app.log.Error("%v", err)
		return err
	}

	return nil
}

func (app App) runOnce(ctx context.Context) error {
	app.log.Print("processing mirrors")

	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(app.conf.Mirrors))
	mirrorErrors := make(chan error)
	for i, mirror := range app.conf.Mirrors {
		go app.processMirror(appCtx, i, mirror, mirrorErrors, &wg)
	}

	go func() {
		wg.Wait()
		close(mirrorErrors)
		app.log.Verbose("closing error channel")
	}()

	for err := range mirrorErrors {
		if err != nil {
			return fmt.Errorf("error processing mirror: %w", err)
		}
	}

	return nil
}

func (app App) processMirror(
	ctx context.Context,
	idx int, mirror config.Mirror, errChan chan<- error, wg *sync.WaitGroup,
) {
	defer wg.Done()
	app.log.Verbose("processing mirror %d", idx)

	app.log.Verbose("getting mirror %d source", idx)
	src := mirror.GetSource()
	if src == nil {
		errChan <- fmt.Errorf("unable to get mirror %d source", idx)
		return
	}

	app.log.Verbose("listing repos for mirror %d", idx)
	repoNames, err := src.ListRepos(ctx)
	if err != nil {
		errChan <- fmt.Errorf("mirror %d: %w", idx, err)
		return
	}

	semaphore := make(chan struct{}, app.concurrency)
	repoErrs := make(chan error)
	var repoWg sync.WaitGroup
	repoWg.Add(len(repoNames))
	app.log.Print("found %d repo(s)", len(repoNames))
	for i, repoName := range repoNames {
		go func(repoName string, idx int) {
			defer repoWg.Done()

			select {
			case <-ctx.Done():
				return
			case semaphore <- struct{}{}:
			}

			defer func() {
				select {
				case <-ctx.Done():
					return
				case <-semaphore:
				}
			}()

			select {
			case <-ctx.Done():
				return
			case repoErrs <- app.processRepo(
				ctx, mirror.CloneTo, repoName, src,
			):
			}
		}(repoName, i)
	}

	go func() {
		repoWg.Wait()
		close(repoErrs)
		close(semaphore)
		app.log.Verbose("closing mirror channels")
	}()

	for err := range repoErrs {
		if err != nil && !errors.Is(err, context.Canceled) {
			app.log.Error("mirroring repo: %s", err)
		}
	}
}

func (app App) processRepo(
	ctx context.Context, repoPath, repoName string, src source.RepoSource,
) error {
	app.log.Print("mirroring repo '%s'", repoName)

	app.log.Verbose("getting repo '%s' from source", repoName)
	repoUrl, err := src.GetRepo(ctx, repoName)
	if err != nil {
		return fmt.Errorf("repo '%s': %w", repoName, err)
	}

	app.log.Verbose("cloning/fetching repo '%s'", repoName)
	_, err = git.Mirror(ctx, repoPath, repoName, repoUrl, app.log)
	if err != nil {
		return fmt.Errorf("repo '%s': %w", repoName, err)
	}

	app.log.Print("mirrored repo '%s'", repoName)
	return nil
}
