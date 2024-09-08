package commands

import (
	"flag"
	"fmt"
	"log"
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
	if info.IsDir() {
		log.Fatal("Directory '" + curlboxDir + "' already exists")
	}

	if err := os.MkdirAll(curlboxDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Create .curlbox-root file to mark the root of the curlbox
	file, err := os.Create(path.Join(curlboxDir + ".curlbox-root"))
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// TODO: Create .gitignore and ignore secrets.toml

	fmt.Println("Created curlbox " + createCmd.Arg(0))
}
