package main

import (
	"fmt"
	"os"

	generatedcli "github.com/usesapient/cli/internal/cli"
	"github.com/usesapient/cli/internal/customcli"
)

// version and buildTime can be set at build time using Go linker flags.
var version string
var buildTime string

func main() {
	if version != "" {
		generatedcli.Version = version
	}
	if buildTime != "" {
		generatedcli.BuildTime = buildTime
	}

	if err := customcli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
