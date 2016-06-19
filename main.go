package main

import (
	"os"

	_ "github.com/openshift/origin/pkg/api/install"

	"github.com/feedhenry/openshift-template-tool/cmd"
)

func main() {
	if err := cmd.NewRootCommand(os.Stdin, os.Stdout, os.Stderr).Execute(); err != nil {
		os.Exit(1)
	}
}
