package commands

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

func Run(runCmd *flag.FlagSet) {
	runCmd.Parse(os.Args[2:])
	if runCmd.NArg() == 0 {
		Help()
		os.Exit(1)
	}

	scriptPath := validateScriptPath(runCmd.Arg(0))

	if !filepath.IsAbs(scriptPath) {
		cwd, _ := os.Getwd()
		scriptPath = filepath.Join(cwd, scriptPath)
	}

	// Get the environment of variables to use
	env := os.Getenv("ENV")
	if env == "" {
		env = "default"
	}

	// Pass any extra arguments to the script
	extraArgs := make([]string, 0)
	for idx, arg := range runCmd.Args() {
		if idx == 0 {
			continue
		}
		extraArgs = append(extraArgs, arg)
	}

	// Send script output right back out to the terminal
	cmd := exec.Command(scriptPath, extraArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Load variables into command environment
	loadVariablesIntoScriptEnv(scriptPath, env, cmd)

	// Run the script in the script directory
	// This makes it easier to chain scripts together by using the relative path
	cmd.Dir = filepath.Dir(scriptPath)

	slog.Debug("Running script", "Script", scriptPath, "Env", env, "Args", extraArgs, "Variables", cmd.Env)

	// Passthrough PATH environment variable
	// so any scripts can use the same PATH as the parent process
	cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))

	runScript(cmd)
}

// validateScriptPath ensures the script path is valid, prints an error and exits if it isn't
func validateScriptPath(scriptPath string) string {
	info, err := os.Stat(scriptPath)
	if os.IsNotExist(err) {
		slog.Error("Script file not found", "Path", scriptPath)
		os.Exit(1)
	}

	if info.IsDir() {
		slog.Error("Script path is a directory", "Path", scriptPath)
		os.Exit(1)
	}

	return scriptPath
}

// loadVariablesIntoScriptEnv loads the variables from the variable files into the script environment
func loadVariablesIntoScriptEnv(scriptFile string, env string, cmd *exec.Cmd) {
	scriptDir := filepath.Dir(scriptFile)

	varFiles, err := findVarFiles(scriptDir)
	if err != nil {
		slog.Error("Error while finding variable files", "Error", err.Error())
		os.Exit(1)
	}

	envVars := make(map[string]string)
	for _, varFile := range varFiles {
		slog.Debug("Parsing variables", "File", varFile)

		// Parse file as toml into a map
		var vars map[string]map[string]any
		if _, err := toml.DecodeFile(varFile, &vars); err != nil {
			slog.Error("Error parsing vars file", "File", varFile, "Error", err)
			os.Exit(1)
		}

		// Prefer the specific environment but fallback to default
		if data, ok := vars[env]; ok {
			for k, v := range data {
				envVars[k] = fmt.Sprintf("%v", v)
			}
		} else if data, ok := vars["default"]; ok {
			slog.Warn("Environment not found, using 'default' instead", "Env", env, "File", varFile)
			for k, v := range data {
				envVars[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	for k, v := range envVars {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
}

// findVarFiles finds all the variable files in the directory tree
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

	// Reverse the order of the files so that the most specific files are loaded first
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}

	return files, nil
}

// runScript runs the script
// If the script is not executable, it will prompt the user to make it executable and rerun
func runScript(cmd *exec.Cmd) {
	err := cmd.Run()

	if err != nil {
		// If permission denied error then prompt user to autofix
		if errors.Is(err, os.ErrPermission) {
			fmt.Println("Script is not executable, do you want to make executable and rerun? (y/N)")

			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				os.Exit(1)
			}

			if strings.ToLower(response) == "y" {
				// Change the file mode to executable
				scriptPath := cmd.Args[0]
				if err := os.Chmod(scriptPath, 0755); err != nil {
					slog.Error("Error changing file mode", "Error", err.Error())
					os.Exit(1)
				}

				// Run the script again
				runScript(cmd)
				return
			}
		}

		slog.Error("Error running script", "Error", err.Error())
		os.Exit(1)
	}
}
