package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/danstewart/curlbox/commands"
	"github.com/lmittmann/tint"
)

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCmd.Usage = commands.Help

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runCmd.Usage = commands.Help

	configureLogging()

	if len(os.Args) < 2 {
		commands.Help()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		commands.Create(createCmd)
	case "run":
		commands.Run(runCmd)
	default:
		commands.Help()
		os.Exit(1)
	}
}

func configureLogging() {
	logLevel := slog.LevelInfo
	if os.Getenv("DEBUG") == "1" {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.DateTime,
		}),
	))
}
