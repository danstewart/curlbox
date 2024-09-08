package commands

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func Run(runCmd *flag.FlagSet) {
	runCmd.Parse(os.Args[2:])
	if runCmd.NArg() == 0 {
		Help()
		os.Exit(1)
	}

	var scriptPath = runCmd.Arg(0)

	info, err := os.Stat(scriptPath)
	if os.IsNotExist(err) {
		slog.Error("Script file not found", "Path", scriptPath)
		os.Exit(1)
	}
	if info.IsDir() {
		slog.Error("Script path is a directory", "Path", scriptPath)
		os.Exit(1)
	}

	var scriptDir = filepath.Dir(scriptPath)

	varFiles, err := findVarFiles(scriptDir)
	if err != nil {
		slog.Error("Error while finding variable files", "Error", err.Error())
		os.Exit(1)
	}

	// Get the environment of variables to use
	env := os.Getenv("ENV")
	if env == "" {
		env = "default"
	}
	slog.Debug("Running script", "Script", scriptPath, "Env", env)

	cmd := exec.Command(scriptPath)

	// Load variables into command environment
	// TODO: Need to test resolution order of variables
	for _, varFile := range varFiles {
		slog.Debug("Parsing variables", "File", varFile)

		// Parse file as toml into a map
		var vars map[string]map[string]any
		if _, err := toml.DecodeFile(varFile, &vars); err != nil {
			slog.Error("Error parsing vars file", "File", varFile, "Error", err)
			os.Exit(1)
		}

		// Error if the environment doesn't exist
		if _, ok := vars[env]; !ok {
			slog.Warn("Environment not found", "ENV", env, "FILE", varFile)
		}

		for k, v := range vars[env] {
			cmd.Env = append(cmd.Env, k+"="+fmt.Sprintf("%v", v))
		}

		slog.Debug("Loaded variables", "Data", cmd.Env)
	}

	out, err := cmd.Output()
	if err != nil {
		slog.Error("Error running script", "Error", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(out))

	// TODO:
	// Load all vars
	// Run script
}

func findVarFiles(dir string) ([]string, error) {
	var files []string
	var foundRoot = false

	for !foundRoot {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && path != dir {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			if info.Name() == "vars.toml" || info.Name() == "secrets.toml" {
				files = append(files, path)
			}

			if info.Name() == ".curlbox-root" {
				foundRoot = true
			}

			return nil
		})

		// Traverse up the directory tree until we find the root of the curlbox
		if !foundRoot {
			// If we reach the root of the filesystem, we've hit a problem...
			if dir == "/" || dir == "." {
				return nil, errors.New("could not find the root of the curlbox")
			}

			dir = filepath.Dir(dir)
		}

		if err != nil {
			return nil, err
		}
	}

	return files, nil
}
