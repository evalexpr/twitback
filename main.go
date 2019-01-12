package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/W1lkins/twitback/version"
	"github.com/genuinetools/pkg/cli"
	"github.com/sirupsen/logrus"
)

var (
	// Whether or not to enable debug logging
	debug bool

	// Where to store the downloaded favorites
	directory string

	// How often to run
	interval time.Duration

	// Whether we should run once or not
	once bool

	// Twitter consumer key
	cKey string
	// Twitter consumer secret
	cSecret string
	// Twitter access token
	tToken string
	// Twitter access secret
	tSecret string
)

func main() {
	p := cli.NewProgram()
	p.Name = "twitback"
	p.Description = "A bot to backup your favorites from Twitter"
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	p.FlagSet = flag.NewFlagSet("twitback", flag.ExitOnError)
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")

	p.FlagSet.BoolVar(&once, "once", false, "run once and exit")
	p.FlagSet.DurationVar(&interval, "interval", 20*time.Hour, "update interval (ex. 5ms, 10s, 1m, 3h)")

	p.FlagSet.StringVar(&directory, "dir", "downloads", "directory to store the downloaded favorites in")

	p.FlagSet.StringVar(&cKey, "consumer-key", os.Getenv("CONSUMER_KEY"), "twitter consumer key")
	p.FlagSet.StringVar(&cSecret, "consumer-secret", os.Getenv("CONSUMER_SECRET"), "twitter consumer secret")
	p.FlagSet.StringVar(&tToken, "access-token", os.Getenv("ACCESS_TOKEN"), "twitter access token")
	p.FlagSet.StringVar(&tSecret, "access-secret", os.Getenv("ACCESS_SECRET"), "twitter access secret")

	p.Before = func(ctx context.Context) error {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	p.Action = func(ctx context.Context, args []string) error {
		if cKey == "" || cSecret == "" || tToken == "" || tSecret == "" {
			return fmt.Errorf("consumer-key, consumer-secret, access-token, and access-secret are required")
		}

		ticker := time.NewTicker(interval)

		if _, err := os.Stat(directory); os.IsNotExist(err) {
			if err := os.MkdirAll(directory, 0755); err != nil {
				return fmt.Errorf("creating directory %s failed: %v", directory, err)
			}
		}

		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		if once {
			run()
			os.Exit(0)
		}

		logrus.Infof("starting bot to update every %s", interval)
		for range ticker.C {
			run()
		}

		return nil
	}

	// Run our program.
	p.Run()
}

func run() error {
	c := NewClient()
	if err := c.VerifyCredentials(); err != nil {
		return fmt.Errorf("could not authenticate with the Twitter API: %v", err)
	}
	if err := c.DownloadFavorites(); err != nil {
		return err
	}
	return nil
}
