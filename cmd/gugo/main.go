package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/figglewatts/gitupgitout/internal/ctxcron"
	"github.com/figglewatts/gitupgitout/pkg/app"
	"github.com/figglewatts/gitupgitout/pkg/logging"
	"github.com/robfig/cron/v3"
)

func runApp(log logging.StdLogger) (err error) {
	var (
		verboseFlag = flag.Bool("verbose", false, "enable verbose logging")
		configFlag  = flag.String(
			"config", "gugo.yaml", "path to config file",
		)
		concurrencyFlag = flag.Int(
			"concurrency", 0, "number of mirrors to process at once",
		)
		cronFlag          = flag.String("cron", "", "run on a cron schedule")
		runBeforeCronFlag = flag.Bool(
			"run-before-cron", false, "run anyway before first scheduled cron",
		)
	)
	flag.Parse()

	if *verboseFlag {
		log.LogVerbose = true
	}

	// setup graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-signalChan
		cancel()
		close(signalChan)
	}()

	configFile, err := os.Open(*configFlag)
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	defer func(configFile *os.File) {
		deferErr := configFile.Close()
		if err == nil && deferErr != nil {
			err = deferErr
		}
	}(configFile)

	gugoApp, err := app.New(
		configFile, log, app.WithConcurrency(*concurrencyFlag),
	)
	if err != nil {
		return fmt.Errorf("initialising app: %w", err)
	}

	if *cronFlag == "" {
		return gugoApp.Run(ctx)
	} else {
		schedule, err := cron.ParseStandard(*cronFlag)
		if err != nil {
			return fmt.Errorf("parsing cron schedule: %w", err)
		}
		log.Print("running on cron schedule: %v", *cronFlag)

		cronScheduler := ctxcron.NewCronScheduler()
		_ = cronScheduler.AddFunc(
			*cronFlag, func() {
				err := gugoApp.Run(ctx)
				if err != nil {
					log.Error("%s", err)
				}
				nextExecution := schedule.Next(time.Now())
				log.Print(
					"awaiting next cron execution at %s",
					nextExecution.Format(time.RFC1123),
				)
			},
		)

		if *runBeforeCronFlag {
			err := gugoApp.Run(ctx)
			if err != nil {
				log.Error("%s", err)
			}
		}

		nextExecution := schedule.Next(time.Now())
		log.Print(
			"awaiting next cron execution at %s",
			nextExecution.Format(time.RFC1123),
		)
		cronScheduler.Run(ctx)
		return nil
	}
}

func main() {
	log := logging.NewStd()
	if err := runApp(log); err != nil {
		log.Error("%s", err)
		os.Exit(1)
	}
}
