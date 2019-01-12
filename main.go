package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/genuinetools/img/version"
	"github.com/genuinetools/pkg/cli"
	"github.com/sirupsen/logrus"
)

var (
	debug     bool
	directory string
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

	p.FlagSet.StringVar(&directory, "dir", "downloads", "directory to store the downloaded favorites in")

	p.Before = func(ctx context.Context) error {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		logrus.WithFields(logrus.Fields{
			"dir": directory,
		}).Info("successfully created directory")

		return nil
	}

	p.Action = func(ctx context.Context, args []string) error {
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

		return nil
	}

	// Run our program.
	p.Run()
}
