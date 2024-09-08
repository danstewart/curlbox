package commands

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
)

func Create(createCmd *flag.FlagSet) {
	createCmd.Parse(os.Args[2:])
	if createCmd.NArg() == 0 {
		Help()
		os.Exit(1)
	}

	var curlboxDir = createCmd.Arg(0)

	info, _ := os.Stat(curlboxDir)
	if info != nil && info.IsDir() {
		slog.Error("Directory already exists", "Path", curlboxDir)
		os.Exit(1)
	}

	if err := os.MkdirAll(curlboxDir, os.ModePerm); err != nil {
		slog.Error("Error creating directory", "Path", curlboxDir, "Error", err)
		os.Exit(1)
	}

	// Create .curlbox-root file to mark the root of the curlbox
	curlboxRootFile := path.Join(curlboxDir, ".curlbox-root")
	file, err := os.Create(curlboxRootFile)
	if err != nil {
		slog.Error("Error creating .curlbox-root file", "Path", curlboxRootFile, "Error", err)
		os.Exit(1)
	}
	file.Close()

	// Create .gitignore and ignore secrets.toml
	gitignoreFile := path.Join(curlboxDir, ".gitignore")
	file, err = os.Create(gitignoreFile)
	if err != nil {
		slog.Error("Error creating .gitignore file", "Path", gitignoreFile, "Error", err)
		os.Exit(1)
	}
	file.WriteString("secrets.toml\n")
	file.Close()

	fmt.Println("Created curlbox " + createCmd.Arg(0))
}
